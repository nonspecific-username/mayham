package dsl


import (
    "encoding/json"
    "errors"
    "fmt"

    yaml "gopkg.in/yaml.v2"
)


type ModConfig struct {
    Name string `yaml:"name",omitempty json:"name"`
    Description string `yaml:"description" json:"description"`
    Enabled bool `yaml:"enabled" json:"enabled"`
    NumActors []NumActorsMod `yaml:"num_actors" json:"num_actors"`
}


func NewModConfig() *ModConfig {
    return &ModConfig{}
}


func (cfg *ModConfig) FromYAML(input *[]byte) (*ModConfig, error, *[]error) {
    err := yaml.UnmarshalStrict(*input, cfg)

    if err != nil {
        return nil, err, nil
    }

    validationErrors := cfg.Validate()
    if len(*validationErrors) > 0 {
        return nil, errors.New("Failed to validate ModConfig"), validationErrors
    }

    return cfg, nil, nil
}


func (cfg *ModConfig) FromJSON(input *[]byte) (*ModConfig, error, *[]error) {
    err := json.Unmarshal(*input, cfg)

    if err != nil {
        return nil, err, nil
    }

    validationErrors := cfg.Validate()
    if len(*validationErrors) > 0 {
        return nil, errors.New("Failed to validate ModConfig"), validationErrors
    }

    return cfg, nil, nil
}


func (cfg *ModConfig) Validate() *[]error {
    var errorsOutput []error

    if cfg.Name == "" {
        msg := "\"name\" is required"
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    for i, mod := range(cfg.NumActors) {
        errs := mod.Validate()
        if errs != nil {
            for _, e := range(*errs) {
                msg := fmt.Sprintf("NumActors[%d]: %v", i, e)
                errorsOutput = append(errorsOutput, errors.New(msg))
            }
        }
    }

    return &errorsOutput
}


func (cfg ModConfig) YAML() string {
    str, _ := yaml.Marshal(cfg)

    return string(str)
}


func (cfg ModConfig) JSON() string {
    str, _ := json.Marshal(cfg)

    return string(str)
}

