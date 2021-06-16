package persistent


import (
    "io/ioutil"
    "time"
    "os"

    "github.com/nonspecific-username/mayham/dsl"
)


func WatchFile(path string) (*dsl.MultiDSLConfig, chan bool, chan bool, error, *[]error) {
    var cfg dsl.MultiDSLConfig
    var errs *[]error
    if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
        cfg = dsl.NewMulti()
    } else {
        data, err := ioutil.ReadFile(path)
        if err != nil {
            return nil, nil, nil, err, nil
        }

        cfg, err, errs = dsl.LoadMulti(&data)
        if err != nil {
            return nil, nil, nil, err, errs
        }
    }

    needToSync := make(chan bool)
    shouldStop := make(chan bool)

    ticker := time.NewTicker(1 * time.Second)
    go func () {
        for {
            select {
            case <-shouldStop:
                ticker.Stop()
                return
            case <- needToSync:
                newData := cfg.String()
                ioutil.WriteFile(path, []byte(newData), 0644)
            case <-ticker.C:
            }
        }
    }()

    return &cfg, needToSync, shouldStop, nil, nil
}
