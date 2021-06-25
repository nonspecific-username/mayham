package state


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


func PersistentState(path string) (*MultiModConfig, error, *dsl.ValidationError) {
    var cfg MultiModConfig
    var validationError *dsl.ValidationError
    if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
        cfg = NewMulti()
    } else {
        data, err := ioutil.ReadFile(path)
        if err != nil {
            return nil, err, nil
        }

        cfg, err, validationError = LoadMultiYAML(&data)
        if err != nil {
            return nil, err, validationError
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
                newData := cfg.YAML()
                ioutil.WriteFile(path, []byte(newData), 0644)
            case <-ticker.C:
            }
        }
    }()

    return &cfg, nil, nil
}


func ClosePersistentState() {
    Sync()
    time.Sleep(time.Second)
    stopCh <- true
}


func Sync() {
    syncCh <- true
}
