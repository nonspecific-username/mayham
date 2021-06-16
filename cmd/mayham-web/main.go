package main


import (
    "log"

    "github.com/nonspecific-username/mayham/web"
)


func main() {
    if err := web.Init(); err != nil {
        log.Fatal(err)
    }

    return
}
