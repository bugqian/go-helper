package parallel

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Func func() error
type job struct {
	name         string // 任务名称
	task         []Func // 具体执行任务
	parallelSize int    // 并发任务个数
	taskChan     chan taskDetail
	failTaskChan chan failTask
}
type taskDetail struct {
	index    int
	execTask Func
}

type failTask struct {
	index int
	err   error
}

func (j *job) run() error {
	start := time.Now()
	log.Println(fmt.Sprintf("start parallel job{%s},size=%d,parallel num %d", j.name, len(j.task), j.parallelSize))
	j.taskChan = make(chan taskDetail, j.parallelSize) // 确认并发管道数量
	j.failTaskChan = make(chan failTask, len(j.task))  // 最多任务量错误
	workDoneChan := make(chan struct{}, j.parallelSize)
	// 分发任务
	go func() {
		defer close(j.taskChan)
		for i, task := range j.task {
			j.taskChan <- taskDetail{index: i, execTask: task}
		}
	}()

	// 根据并发数量建立协程
	for i := 0; i < j.parallelSize; i++ {
		go func() {
			for task := range j.taskChan {
				err := task.execTask()
				if err != nil {
					j.failTaskChan <- failTask{err: fmt.Errorf("task %d execute failed: %w", task.index, err), index: task.index}
				}
			}
			workDoneChan <- struct{}{} // 完成计数
		}()
	}

	go func() {
		// 等待所有通道结束任务执行
		for i := 0; i < j.parallelSize; i++ {
			<-workDoneChan
		}
		close(workDoneChan)   // 关闭工作完成通道
		close(j.failTaskChan) // 关闭错误通道
	}()

	var errMsg []string
	for ch := range j.failTaskChan {
		errMsg = append(errMsg, ch.err.Error())
	}

	if len(errMsg) > 0 {
		return errors.New(fmt.Sprintf("parallel job {%s} fail ,time=%v,err:%s",
			j.name, time.Since(start), strings.Join(errMsg, ",")))
	}

	log.Println(fmt.Sprintf("parallel job {%s} finished,time=%v", j.name, time.Since(start)))
	return nil
}
