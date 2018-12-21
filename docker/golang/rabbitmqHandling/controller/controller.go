package controller

import (
	"context"
	"log"
	"time"
)

// Runner : 許容しているタスクのインターフェース
type Runner interface {
	Run() error
}

// Controller is struct of clocking controller for someTask
type Controller struct {
	waitTime time.Duration
	// if you change someTask func() iterface, you must change following interface
	task    Runner
	counter uint
	verbos  bool
}

// NewController is contructor
func NewController(waitTime time.Duration, runner Runner, isVerbos bool) Controller {
	clockController := Controller{
		waitTime: waitTime,
		counter:  0,
		verbos:   isVerbos,
		task:     runner,
	}
	return clockController
}

// Exec is main function
func (c Controller) Exec(ctx context.Context) {
	//ticker := time.NewTicker(time.Duration(c.waitSecond) * time.Millisecond)
	ticker := time.NewTicker(c.waitTime)
	defer ticker.Stop()
	child, childCancel := context.WithCancel(ctx)
	defer childCancel()

	for { // deamon化するため無限実行
		select {
		case t := <-ticker.C:
			c.counter++
			// 値を適切に引き渡すため、requestIDを再定義
			// - ref: https://qiita.com/niconegoto/items/3952d3c53d00fccc363b
			requestID := c.counter

			if c.verbos {
				log.Println("[DEBUG] Controller START requestNo=", requestID, "t=", t)
			}

			errCh := make(chan error, 1)
			go func() { // 登録したタスクをブロックせずに実行
				errCh <- c.task.Run()
			}()

			go func() {
				// error channelにリクエストの結果が返ってくるのを待つ
				select {
				case err := <-errCh:
					if err != nil {
						// Deamonの強制終了はしない
						log.Println("[ERROR] Controller.Exec.Run()", err)
					}
					log.Println("[DEBUG] Controller END requestNo=", requestID)
				}
			}()
		case <-child.Done():
			// contextによるExec()の終了
			if c.verbos {
				log.Println("[DEBUG] child cancelled")
			}
			return
		}
	}
}
