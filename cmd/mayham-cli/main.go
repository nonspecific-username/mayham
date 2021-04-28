package main


import (
    "fmt"

    "github.com/nonspecific-username/mayham/hotfix"
)


func main() {
    var hf hotfix.Hotfix
    hf.AddRegular(hotfix.EarlyLevel,
                   0,
                   "MatchAll",
                   "ObjectPath",
                   "Attr",
                   0,
                   "",
                   hotfix.RenderBVCOverride(100))
    hf.AddRegular(hotfix.EarlyLevel,
                   0,
                   "MatchAll",
                   "ObjectPath",
                   "Attr",
                   0,
                   "",
                   hotfix.RenderBVCOverride(300))
    fmt.Printf(hf.Render() + "\n")
}
