package main


import (
    "io/ioutil"
    "log"
    "os"

    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/generator"
)


func main() {
    configs := os.Args[1:]

    for _, confPath := range(configs) {
        data, err := ioutil.ReadFile(confPath)
        if err != nil {
            log.Fatal("Error reading input file " + confPath)
            log.Println(err)
            return
        }

        cfg, err := dsl.Load(&data)
        if err != nil {
            log.Fatal("Error parsing input file " + confPath)
            log.Println(err)
            return
        }

        mod, err := generator.Generate(cfg)
        if err != nil {
            log.Fatal("Error generating hotfixes from input file " + confPath)
            log.Println(err)
            return
        }
        log.Println(mod)
    }
}
