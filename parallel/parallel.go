package parallel

import (
	"errors"
	"fmt"
	"runtime"
)

var defaultParallelSize = runtime.NumCPU()

func Run(name string, tasks []Func, parallelSize ...int) error {
	if len(name) == 0 {
		return errors.New("任务名称为空")
	}
	if len(tasks) == 0 {
		return errors.New("任务数量为空")
	}
	// 确定并行数量
	var useParallelSize = defaultParallelSize
	if len(tasks) < useParallelSize {
		useParallelSize = len(tasks)
	}
	if len(parallelSize) > 0 {
		if parallelSize[0] < useParallelSize {
			useParallelSize = parallelSize[0]
		}
	}

	parallelJob := job{
		name:         name,
		task:         tasks,
		parallelSize: useParallelSize,
		taskChan:     make(chan taskDetail, useParallelSize),
		failTaskChan: make(chan failTask, len(tasks)),
	}
	return parallelJob.run()
}

// Foreach 并发执行 只返回错误
func Foreach[T any](tasks []T, fun func(T) error, name ...string) error {
	var taskName string = "任务"
	if len(name) > 0 {
		taskName = name[0]
	}
	// 构造并行任务
	taskFunc := make([]Func, len(tasks))
	for i, v := range tasks {
		taskFunc[i] = func() error {
			return fun(v)
		}
	}
	// 执行任务
	return Run(taskName, taskFunc, defaultParallelSize)
}

// Map 并行返回值
func Map[V any, T any](tasks []T, fun func(T) (V, error), name ...string) (res []V, err error) {
	var taskName string
	if len(name) > 0 {
		taskName = name[0]
	}
	respChan := make(chan V, len(tasks)) // 建立返回数据通道
	collected := make(chan bool)         // 执行完成通道标记
	// 构造并行任务
	taskFunc := make([]Func, len(tasks))

	for i, v := range tasks {
		taskFunc[i] = func() error { // 闭包函数
			resp, e := fun(v)
			if e != nil {
				return e
			}
			respChan <- resp // 返回值放回通道
			return nil
		}
	}

	go func() {
		for ch := range respChan {
			res = append(res, ch)
		}
		collected <- true // 返回值收集结束 主线程结束
	}()

	go func() {
		defer close(respChan) // 关闭返回信息通道
		err = Run(fmt.Sprintf("并行处理 %T %s", fun, taskName), taskFunc, defaultParallelSize)
	}()

	<-collected

	return
}
