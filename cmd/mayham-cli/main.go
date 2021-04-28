package main


import (
    "fmt"

    "github.com/nonspecific-username/mayham/hotfix"
    "github.com/nonspecific-username/mayham/dsl"
)


var test_data string = `
SpawnNum:
- spawn:
    pkg: Game
    map: ".*"
    spawn: ".*"
  mode: factor
  param1: 5
  max_actors_mode: match
`


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
    cfg, _ := dsl.Load(test_data)
    fmt.Println(*cfg)
}
