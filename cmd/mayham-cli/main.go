package main


import (
    "fmt"

    "github.com/nonspecific-username/mayham/hotfix"
)


func main() {
    var hf hotfix.Hotfix = hotfix.Regular{Method: "SparkEarlyLevelPatchEntry",
                                          Notify: 0,
                                          Pkg: "MatchAll",
                                          Object: "ObjectPath",
                                          Attr: "Attr",
                                          FromLen: 0,
                                          From: "",
                                          Value: hotfix.RenderBVCOverride(100)}
    fmt.Printf(hf.Render() + "\n")
}
