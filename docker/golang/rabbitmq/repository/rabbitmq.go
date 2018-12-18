package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// QueueHandler is queue output struct
type QueueHandler struct {
	DeliveryTag uint64
	Body        string
}

// RabbitMQ structure has information for connecting rabbitMQ
type RabbitMQ struct {
	URI          string
	Name         string
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	Queue        amqp.Queue
	DequeueUpper int
}

// NewRabbitMQ :init RabbitMQ
func NewRabbitMQ(user, password, host, vhost, port, qName string, dequeueUpper int) (*RabbitMQ, error) {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost)
	log.Println("[INFO] NewRabbitMQ() ", amqpURI)

	if dequeueUpper <= 0 {
		return &RabbitMQ{}, fmt.Errorf("should use [upper number] grater than 1. got:%v", dequeueUpper)
	}
	mq := &RabbitMQ{
		URI:          amqpURI,
		Name:         qName,
		DequeueUpper: dequeueUpper,
	}

	if err := mq.newConnection(); err != nil {
		return &RabbitMQ{}, err
	}
	return mq, nil
}

func (r *RabbitMQ) failOnError(err error, msg string) error {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
		return err
	}
	return nil
}

func (r *RabbitMQ) newConnection() error {
	log.Printf("[DEBUG] newConnection() with %#v", r)

	conn, err := amqp.Dial(r.URI)
	if err := r.failOnError(err, "Failed to connect to RabbitMQ"); err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err := r.failOnError(err, "Failed to open a channel"); err != nil {
		return err
	}
	//que, err := ch.QueueDeclare(
	//	"test", //name
	//	true,   // durable
	//	false,  // delete when unused
	//	false,  // exclusive
	//	false,  // no-wait
	//	nil,    // argument
	//)

	que, err := ch.QueueInspect(r.Name)
	if err := r.failOnError(err, "Failed to declare queue"); err != nil {
		return err
	}
	r.Connection = conn
	r.Channel = ch
	r.Queue = que
	return nil
}

// Enqueue function execute enqueueing to rabbitMQ
func (r *RabbitMQ) Enqueue(name, payload string) error {
	log.Printf("[DEBUG] Enqueue() to %v with %v", name, payload)
	err := r.Channel.Publish(
		"",
		name,
		false,
		false,
		amqp.Publishing{
			ContentType: "test/plain",
			Body:        []byte(payload),
		},
	)
	if err := r.failOnError(err, "Faild to publish a message"); err != nil {
		return err
	}
	return nil
}

// Dequeue function execute dequeueing from rabbitMQ
// Ackを返すとキューからメッセージが削除されるためAuto-Ackはfalseにする
func (r *RabbitMQ) Dequeue() ([]QueueHandler, error) {

	// キュー内のメッセージ数を取得.最大件数以下となるようにエンキューする
	messageCount := r.Queue.Messages
	dequeueUpper := messageCount
	if messageCount == 0 {
		return []QueueHandler{}, errors.New("nothing Queue from: queueName[" + r.Name)
	}
	if messageCount >= r.DequeueUpper {
		dequeueUpper = r.DequeueUpper
	}

	msgs, err := r.Channel.Consume(
		r.Queue.Name,
		"",
		false, //Auto-Ack: true / 明示的なAckが必要になる: false
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		return []QueueHandler{}, err
	}

	var mqs []QueueHandler //output
	forever := make(chan int)
	ch := make(chan int, 2) //groutine上限の設定

	go func() {
		for msg := range msgs {
			ch <- 1
			log.Printf("[DEBUG] load msgs: ch=%v /%v", ch, r.Queue.Messages)
			// mq := streamToString(msg.Body)
			// mq := string(msg.Body)
			mq := QueueHandler{
				DeliveryTag: msg.DeliveryTag,
				Body:        string(msg.Body),
			}
			mqs = append(mqs, mq)
			forever <- 1
			<-ch
		}
	}()

	// キュー内のメッセージをすべて受け取るまで待機
	// wg.Wait() // キュー内のメッセージをすべて受け取るまで待機
	for i := 0; i < dequeueUpper; i++ {
		<-forever
	}

	return mqs, nil
}

// Ack function execute ack.
// 価格更新が成功したメッセージにAckを返し、キュー内から削除する
func (r *RabbitMQ) Ack(deliveryTag uint64) error {
	// tag, err := strconv.ParseUint(id, 10, 64)
	// if err != nil {
	// 	return err
	// }

	// WARNING
	// if err := r.newConnection(); err != nil {
	// 	return err
	// }
	log.Printf("[DEBUG] Ack() with tag=%v", deliveryTag)
	// 第二引数[multiple]をtrueにしてしまうと、同一チャネル内のキューがすべてAckされてしまうので注意
	// ref: https://godoc.org/github.com/streadway/amqp#Channel.Ack
	// - When multiple is true, this delivery and all prior unacknowledged deliveries on the same channel will be acknowledged. This is useful for batch processing of deliveries.
	if err := r.Channel.Ack(deliveryTag, false); err != nil {
		return err
	}
	return nil
}

// Nack function execute nack and reject.
// 価格更新が失敗したメッセージにNackを返す。(=再度キューが配信される状態とする)
// NOTE: 原因不明であるが、NackしたたとにRejectしないと次の更新時に再度メッセージを取得することができなくなる
func (r *RabbitMQ) Nack(deliveryTag uint64) error {
	if err := r.Channel.Nack(deliveryTag, false, true); err != nil {
		return err
	}
	if err := r.Channel.Reject(deliveryTag, true); err != nil { //第二引数=trueで、キューを再度送信される状態とする
		return err
	}
	return nil
}
