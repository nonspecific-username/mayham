package persistent


import (
    "io/ioutil"
    "time"
    "os"

    "github.com/nonspecific-username/mayham/dsl"
)


var (
    syncCh chan bool
    stopCh chan bool
)


func Open(path string) (*dsl.MultiDSLConfig, error, *[]error) {
    var cfg dsl.MultiDSLConfig
    var errs *[]error
    if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
        cfg = dsl.NewMulti()
    } else {
        data, err := ioutil.ReadFile(path)
        if err != nil {
            return nil, err, nil
        }

        cfg, err, errs = dsl.LoadMulti(&data)
        if err != nil {
            return nil, err, errs
        }
    }

    syncCh = make(chan bool)
    stopCh = make(chan bool)

    ticker := time.NewTicker(1 * time.Second)
    go func () {
        for {
            select {
            case <- stopCh:
                ticker.Stop()
                return
            case <- syncCh:
                newData := cfg.String()
                ioutil.WriteFile(path, []byte(newData), 0644)
            case <-ticker.C:
            }
        }
    }()

    return &cfg, nil, nil
}


func Close() {
    Sync()
    time.Sleep(time.Second)
    stopCh <- true
}


func Sync() {
    syncCh <- true
}
