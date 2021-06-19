package dsl


import (
    "encoding/json"
    "errors"
    "fmt"

    yaml "gopkg.in/yaml.v2"
)


type DSLConfig struct {
    Name string `yaml:"name",omitempty json:"name"`
    Description string `yaml:"description" json:"description"`
    Enabled bool `yaml:"enabled" json:"enabled"`
    SpawnNum []SpawnNumMod `yaml:"spawn_num" json:"spawn_num"`
}


func NewDSLConfig() *DSLConfig {
    cfg := &DSLConfig{}
    return cfg
}


func (cfg *DSLConfig) FromYAML(input *[]byte) (*DSLConfig, error, *[]error) {
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


func (cfg *DSLConfig) FromJSON(input *[]byte) (*DSLConfig, error, *[]error) {
    err := json.Unmarshal(*input, cfg)

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


func (cfg DSLConfig) YAML() string {
    str, _ := yaml.Marshal(cfg)

    return string(str)
}


func (cfg DSLConfig) JSON() string {
    str, _ := json.Marshal(cfg)

    return string(str)
}

