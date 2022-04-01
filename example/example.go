package main

import (
    "fmt"
    "github.com/DGHeroin/ActionJob"
    "time"
)

func main() {
    dispatcher := ActionJob.NewDispatcher(100)
    dispatcher.Start()

    dispatcher.Run(func() {
        fmt.Println("job 1")
    })
    dispatcher.Run(func() {
        time.Sleep(time.Second * 2)
        fmt.Println("job 2")
    })
    dispatcher.Run(func() {
        fmt.Println("job 3")
    })

    dispatcher.Run(func() {
        dispatcher.Stop()
        fmt.Println("job stop")
    })
    time.Sleep(time.Second)
    dispatcher.WaitJobDone()
}
