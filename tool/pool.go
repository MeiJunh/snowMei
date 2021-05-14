package tool

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"runtime"
	"sync/atomic"
	"time"
)

/*
	简单资源池实现
*/

// 资源池
type Task struct {
	ID       int64       // 可以按照id进行分配
	TaskInfo interface{} // 函数入参信息
}

type Pool struct {
	size     int
	work     chan *Task       // 任务
	workFunc func(info *Task) // 任务执行函数
	Ch       []chan *Task
	curidx   uint32
	closeCh  chan struct{} // 协程池结束通知
	option   *option       // 部分配置参数
}

type option struct {
	timeOut      int64  // 单个任务超时时间，单位秒
	callFuncName string // 调用方函数信息。支持传入,默认使用上层函数名
}

type poolOptionSet func(option *option)

func SetPoolTimeOut(timeOut int64) poolOptionSet {
	return func(option *option) {
		option.timeOut = timeOut
	}
}

func SetPoolCallFuncName(callFuncName string) poolOptionSet {
	return func(option *option) {
		option.callFuncName = callFuncName
	}
}

func New(size, queueLen int, workFunc func(info *Task), options ...poolOptionSet) *Pool {
	// 保证线程池的正常运行
	if size <= 0 {
		size = 50
	}

	pool := &Pool{
		size:     size,
		work:     make(chan *Task),
		workFunc: workFunc,
		Ch:       make([]chan *Task, size),
		curidx:   0,
		closeCh:  make(chan struct{}),
		option: &option{
			timeOut:      10,                                                                         // 默认10秒
			callFuncName: GetFuncName(runtime.FuncForPC(reflect.ValueOf(workFunc).Pointer()).Name()), // 默认使用执行函数函数名
		},
	}

	for _, f := range options {
		f(pool.option)
	}

	for i := 0; i < size; i++ {
		if queueLen <= 0 {
			pool.Ch[i] = make(chan *Task)
		} else {
			pool.Ch[i] = make(chan *Task, queueLen)
		}
		idx := i
		go pool.worker(idx)
	}

	return pool
}

func (p *Pool) worker(idx int) {
	for {
		select {
		case <-p.closeCh:
			fmt.Printf("work pool close:%d", idx)
			return
		default:
			p.workerIn(idx)
		}
	}
}

func (p *Pool) workerIn(idx int) {
	for {
		workFinish := make(chan struct{})
		select {
		case task := <-p.Ch[idx]:
			func() {
				defer PanicFunc(func(panic string) {
					fmt.Printf("workerIn panic err:%s", panic)
				})

				tc := time.After(time.Second * time.Duration(p.option.timeOut))

				// 如果workFunc一直不结束就造成内存泄漏
				go func() {
					// 执行任务
					p.workFunc(task)
					select {
					case workFinish <- struct{}{}:
						return
					default:
						return
					}
				}()

				select {
				case <-workFinish:
					close(workFinish)
					return
				case <-tc:
					// 任务超时结束goroutine，防止func一直不返回造成内存泄漏
					taskStr, _ := jsoniter.MarshalToString(task)
					fmt.Printf("call func name:%s,pool work task time out,task:%s", p.option.callFuncName, taskStr)
					return
				}
			}()
		}
	}
}

// Dispatch 根据int64哈希任务分发task
func (p *Pool) Dispatch(task *Task) (err error) {
	idx := task.ID % (int64)(p.size)
	select {
	case p.Ch[idx] <- task:
		return
	}
}

// RandDispatch 顺序随机分发task
func (p *Pool) RandDispatch(task *Task) (err error) {
	idx := atomic.LoadUint32(&p.curidx) % (uint32)(p.size)
	atomic.AddUint32(&p.curidx, 1)
	fmt.Printf("RandDispatch,idx:%d\n", p.curidx)
	select {
	case p.Ch[idx] <- task:
		return
	}
}

func (p *Pool) Close() {
	// 退出协程
	for i := 0; i < p.size; i++ {
		p.closeCh <- struct{}{}
	}
}
