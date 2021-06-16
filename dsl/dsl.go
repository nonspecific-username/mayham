package dsl


import (
    "errors"
    "fmt"

    yaml "gopkg.in/yaml.v2"
)


type DSLConfig struct {
    Name string `yaml:"name",omitempty`
    Description string `yaml:"description"`
    Enabled bool `yaml:"enabled"`
    SpawnNum []SpawnNumMod `yaml:"spawn_num"`
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

    if cfg.Name == "" {
        msg := "\"name\" is required"
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    for i, mod := range(cfg.SpawnNum) {
        errs := mod.Validate()
        if errs != nil {
            for _, e := range(*errs) {
                msg := fmt.Sprintf("SpawnNum[%d]: %v", i, e)
                errorsOutput = append(errorsOutput, errors.New(msg))
            }
        }
    }

    return &errorsOutput
}


func (cfg DSLConfig) String() string {
    output, _ := yaml.Marshal(cfg)
    return string(output)
}

