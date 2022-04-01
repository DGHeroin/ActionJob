package ActionJob

import (
    "sync"
    "sync/atomic"
)

type Dispatcher struct {
    pool     chan chan Job
    workers  []*worker
    num      int
    jobQueue chan Job
    flag     int32
    wg       sync.WaitGroup
}

func NewDispatcher(maxWorkers int) *Dispatcher {
    pool := make(chan chan Job, maxWorkers)
    return &Dispatcher{
        pool:     pool,
        num:      maxWorkers,
        jobQueue: make(chan Job),
    }
}
func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <-d.jobQueue: // 收到job
            go func(job Job) {
                jobChannel := <-d.pool // 取出参与竞争worker的channel
                jobChannel <- job      // 给对应channel派发任务
            }(job)
        }
    }
}

func (d *Dispatcher) Start() {
    d.workers = make([]*worker, d.num)
    d.wg.Add(d.num)

    for i := 0; i < d.num; i++ {
        worker := newWorker(d.pool, &d.wg)
        d.workers[i] = worker
        go worker.start()
    }

    go d.dispatch()
}

func (d *Dispatcher) Run(fn func()) bool {
    if atomic.LoadInt32(&d.flag) == 1 {
        return false // closed
    }
    d.jobQueue <- fn
    return true
}
func (d *Dispatcher) Stop() {
    // 不再接收新任务
    if !atomic.CompareAndSwapInt32(&d.flag, 0, 1) {
        return // closed
    }
    close(d.jobQueue)
    for _, w := range d.workers {
        go w.stop()
    }
}

func (d *Dispatcher) WaitJobDone() {
    d.wg.Wait()
}
