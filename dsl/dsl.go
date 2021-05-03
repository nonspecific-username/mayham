package dsl


import (
    "errors"
    "fmt"

    yaml "gopkg.in/yaml.v2"
)


type DSLConfig struct {
    SpawnNum []SpawnNumMod `yaml:"SpawnNum"`
}


func Load(input *[]byte) (*DSLConfig, error, *[]error) {
    cfg := &DSLConfig{}
    err := yaml.UnmarshalStrict(*input, cfg)

    if err != nil {
        return nil, err, nil
    }

    validationErrors := cfg.Validate()
    if len(*validationErrors) > 0 {
        return nil, errors.New("Failed to validate DSLConfig"), validationErrors
    }

    return cfg, nil, nil
}


func (cfg *DSLConfig) Validate() *[]error {
    var errorsOutput []error
    for i, mod := range(cfg.SpawnNum) {
        err := mod.Validate()
        if err != nil {
            msg := fmt.Sprintf("SpawnNum[%d]: %v", i, err)
            errorsOutput = append(errorsOutput, errors.New(msg))
        }
    }

    return &errorsOutput
}


func (cfg DSLConfig) String() string {
    output, _ := yaml.Marshal(cfg)
    return string(output)
}

