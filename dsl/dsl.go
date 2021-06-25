package dsl


import (
    "encoding/json"
    "errors"
    "fmt"

    yaml "gopkg.in/yaml.v2"
)


type ModConfig struct {
    Name string `yaml:"name" json:"name"`
    Description string `yaml:"description" json:"description"`
    Author string `yaml:"author" json:"author"`
    Enabled bool `yaml:"enabled" json:"enabled"`
    NumActors []NumActorsMod `yaml:"num_actors" json:"num_actors" singleUpdate:"no"`
}


func NewModConfig() *ModConfig {
    return &ModConfig{}
}


func (cfg *ModConfig) FromYAML(input *[]byte) (*ModConfig, error, *ValidationError) {
    err := yaml.UnmarshalStrict(*input, cfg)

    if err != nil {
        return nil, err, nil
    }

    validationError := cfg.Validate()
    if validationError != nil {
        return nil, errors.New("Failed to validate ModConfig"), validationError
    }

    return cfg, nil, nil
}


func (cfg *ModConfig) FromJSON(input *[]byte) (*ModConfig, error, *ValidationError) {
    err := json.Unmarshal(*input, cfg)

    if err != nil {
        return nil, err, nil
    }

    validationError := cfg.Validate()
    if validationError != nil {
        return nil, errors.New("Failed to validate ModConfig"), validationError
    }

    return cfg, nil, nil
}


func (cfg *ModConfig) Validate() *ValidationError {
    var errorsOutput []*ValidationError
    var numActorErrors []*ValidationError

    if cfg.Name == "" {
        errorsOutput = append(errorsOutput, NewValidationError("name", "missing param", nil))
    }

    for i, mod := range(cfg.NumActors) {
        validationError := mod.Validate()
        if validationError != nil {
            validationError.Field = fmt.Sprintf("%d", i)
            numActorErrors = append(numActorErrors, validationError)
        }
    }

    if len(numActorErrors) > 0 {
        numActorError := NewValidationError("num_actors", "NumActors validation failed", numActorErrors)
        errorsOutput = append(errorsOutput, numActorError)
    }

    if len(errorsOutput) > 0 {
        return NewValidationError("", "ModConfig validation failed", errorsOutput)
    } else {
        return nil
    }
}


func (cfg ModConfig) YAML() string {
    str, _ := yaml.Marshal(cfg)

    return string(str)
}


func (cfg ModConfig) JSON() string {
    str, _ := json.Marshal(cfg)

    return string(str)
}

