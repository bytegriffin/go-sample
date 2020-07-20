package concurrent

import (
	"fmt"
	"testing"
)

const (
	NUMBER = 10
)

type task struct {
	begin  int
	end    int
	result chan<- int
}

func (t *task) do() {
	sum := 0
	for i := t.begin; i <= t.end; i++ {
		sum += i
	}
	t.result <- sum
}

func InitTask(taskChan chan<- task, r chan int, p int) {
	qu := p / 10
	mod := p % 10
	high := qu * 10
	for j := 0; j < qu; j++ {
		b := 10*j + 1
		e := 10 * (j + 1)
		tsk := task{
			begin:  b,
			end:    e,
			result: r,
		}
		taskChan <- tsk
	}
	if mod != 0 {
		tsk := task{
			begin:  high + 1,
			end:    p,
			result: r,
		}
		taskChan <- tsk
	}
	close(taskChan)
}

func TestTaskPool(t *testing.T) {
	workers := NUMBER
	taskChan := make(chan task, 10)  // 工作通道
	resultChan := make(chan int, 10) // 结果通道
	done := make(chan struct{}, 10)
	go InitTask(taskChan, resultChan, 100)
	DistributeTask(taskChan, workers, done)
	go CloseResult(done, resultChan, workers)
	sum := ProcessResult(resultChan)
	fmt.Println("sum=", sum)
}

func DistributeTask(taskchan <-chan task, workers int, done chan struct{}) {
	for i := 0; i < workers; i++ {
		go ProcessTask(taskchan, done)
	}
}

func ProcessTask(taskchan <-chan task, done chan struct{}) {
	for t := range taskchan {
		t.do()
	}
	done <- struct{}{}
}

func CloseResult(done chan struct{}, resultchan chan int, workers int) {
	for i := 0; i < workers; i++ {
		<-done
	}
	close(done)
	close(resultchan)
}

func ProcessResult(resultchan chan int) int {
	sum := 0
	for r := range resultchan {
		sum += r
	}
	return sum
}
