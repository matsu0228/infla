package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/matsu0228/infla/rabbitmq/repository"
	"github.com/matsu0228/infla/rabbitmqHandling/controller"
)

// errorExit :エラー終了時の共通処理
func errorExit(err error) {
	log.Fatal("[ERROR] ", err)
}

func envLoad() error {
	if err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("APP_ENV"))); err != nil {
		return err
	}
	return nil
}

func enqueue(mq *repository.RabbitMQ, qn string) {
	for i := 0; i < 1000; i++ {
		v := fmt.Sprintf("%v: test value", i)
		err := mq.Enqueue(qn, v)
		if err != nil {
			errorExit(err)
		}

		time.Sleep(500 * time.Millisecond)
	}
	err := mq.ChannelClose()
	if err != nil {
		errorExit(err)
	}
}

type mqRunner struct {
	MQ *repository.RabbitMQ
}

func newMqRunner(mq *repository.RabbitMQ) *mqRunner {
	return &mqRunner{
		MQ: mq,
	}
}

// Run is task main
func (m *mqRunner) Run() error {
	log.Printf("[INFO] Run() with %#v", *m.MQ)

	defer m.MQ.ChannelClose()
	mqs, err := m.MQ.Dequeue()
	if err != nil {
		errorExit(err)
	}
	// pp.Print(mqs)

	// Ack()
	for _, q := range mqs {
		err := m.MQ.Ack(q.DeliveryTag)
		// err := mq.Nack(q.DeliveryTag)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

	if err := envLoad(); err != nil {
		errorExit(fmt.Errorf("can not load .env.%s :err=%v", os.Getenv("APP_ENV"), err))
	}
	mqUser := os.Getenv("MQ_USER")
	mqPass := os.Getenv("MQ_PASSWORD") //環境変数
	mqHost := os.Getenv("MQ_HOST")
	mqVhost := os.Getenv("MQ_VHOST")
	mqPORT := os.Getenv("MQ_PORT")
	mqQueueName := os.Getenv("MQ_QUEUE_NAME")
	// requestURL := os.Getenv("REQUEST_URL")
	mqDequeueUpper := 50

	// log.Printf("MQ setting: u:%v, p:%v, %v, %v :%v, qName:%v", mqUser, mqPass, mqHost, mqVhost, mqPORT, mqQueueName)
	// pp.Print("requestURL setting: ", requestURL)

	// // MQ初期化
	mq, err := repository.NewRabbitMQ(mqUser, mqPass, mqHost, mqVhost, mqPORT, mqQueueName, mqDequeueUpper)
	if err != nil {
		errorExit(err)
	}

	// pp.Print(mq)

	// create test data
	// enqueue(mq, mqQueueName)

	mqUsecase := newMqRunner(mq)
	ctx := context.Background()
	clockTime := 20 * time.Millisecond
	clockController := controller.NewController(clockTime, mqUsecase, true)
	clockController.Exec(ctx)
}
