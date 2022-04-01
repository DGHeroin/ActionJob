package ActionJob

import "sync"

type Job func()

type worker struct {
    workerPool chan chan Job
    ch         chan Job
    quit       chan bool
    wg         *sync.WaitGroup
}

func newWorker(workerPool chan chan Job, wg *sync.WaitGroup) *worker {
    return &worker{
        workerPool: workerPool,
        ch:         make(chan Job),
        quit:       make(chan bool),
        wg:         wg,
    }
}

func (w *worker) start() {
    defer w.wg.Done()
    for {
        // 把自己的chan放去worker pool竞争job
        w.workerPool <- w.ch
        select {
        case job := <-w.ch: // 拿到job
            go func(job Job) {
                defer func() {
                    recover()
                }()
                job()
            }(job)
        case <-w.quit:
            return
        }
    }
}

func (w *worker) stop() {
    close(w.quit)
}
