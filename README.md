### ActionJob

一个高性能的任务分发器

```golang
import "github.com/DGHeroin/ActionJob"

func foo() {
    dispatcher := ActionJob.NewDispatcher(100)
    dispatcher.Start()

    dispatcher.Run(func() {
        fmt.Println("job 1")
    })

    dispatcher.WaitJobDone()
}
    
    
```
