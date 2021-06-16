package main


import (
//    "fmt"
//    "io/ioutil"
    "log"
    "os"
    "time"

    "github.com/nonspecific-username/mayham/persistent"
    "github.com/nonspecific-username/mayham/dsl"
//    "github.com/nonspecific-username/mayham/generator"
//    "github.com/google/uuid"
)


func main() {
    config := os.Args[1]
/*
   data, err := ioutil.ReadFile(config)
   if err != nil {
       log.Fatal("Error reading input file %s: %v", config, err)
   }

   cfg, err, validationErrors := dsl.LoadMulti(&data)
   if err != nil {
       log.Printf("Error parsing input file %s: %v:", config, err)
       if validationErrors != nil && len(*validationErrors) > 0 {
           for _, valErr := range(*validationErrors) {
               log.Printf("\t%v", valErr)
           }
       }
       os.Exit(1)
   }

   log.Printf("%s\n", cfg.String())
*/

    cfg, syncCh, stopCh, err, validationErrors := persistent.WatchFile(config)
    if err != nil {
        log.Printf("Error parsing input file %s: %v:", config, err)
        if validationErrors != nil && len(*validationErrors) > 0 {
            for _, valErr := range(*validationErrors) {
                log.Printf("\t%v", valErr)
            }
        }
        os.Exit(1)
    }

    //(*cfg)["a078e5ce-a08e-4b0a-8a79-0535786df2e6"].SpawnNum[0].Param1 += 5
    (*cfg)["test"] = &dsl.DSLConfig{
        Name: "Test 1",
        Enabled: true,
        Description: "Test Mod",
    }

    (*cfg)["test"].SpawnNum = append((*cfg)["test"].SpawnNum, dsl.SpawnNumMod{
        Enabled: true,
        Spawn: &dsl.SpawnSelector{},
        Mode: dsl.RandomFactor,
        Param1: 1,
        Param2: 5,
        MaxActorsMode: dsl.MAMatch,
    })
    log.Printf("%s", cfg.String())
    syncCh <- true
    time.Sleep(time.Second * 5)
    stopCh <- true

    return
}
