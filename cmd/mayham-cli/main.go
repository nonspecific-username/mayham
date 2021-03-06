package main


import (
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/generator"
)


func main() {
    configs := os.Args[1:]
    var dslConfigs []*dsl.ModConfig
    var haveErrors bool = false

    for _, confPath := range(configs) {
        data, err := ioutil.ReadFile(confPath)
        if err != nil {
            haveErrors = true
            log.Printf("Error reading input file %s: %v", confPath, err)
            continue
        }

        cfg, err, validationError := dsl.NewModConfig().FromYAML(&data)
        if err != nil {
            haveErrors = true
            log.Printf("Error parsing input file %s: %v", confPath, err)
            if validationError != nil {
                log.Printf(validationError.String())
            }
        }
        dslConfigs = append(dslConfigs, cfg)
    }

    if haveErrors {
        os.Exit(1)
    }

    for _, cfg := range(dslConfigs) {
        fmt.Println(generator.Generate(cfg))
    }

}
