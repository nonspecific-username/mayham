package dsl


import (
    "errors"
    "fmt"

    yaml "gopkg.in/yaml.v2"
)


type MultiDSLConfig map[string]*DSLConfig


func NewMulti() MultiDSLConfig {
    cfg := make(map[string]*DSLConfig)
    return cfg
}

func LoadMulti(input *[]byte) (MultiDSLConfig, error, *[]error) {
    var errorsOutput []error
    cfg := make(map[string]*DSLConfig)
    err := yaml.UnmarshalStrict(*input, cfg)

    if err != nil {
        return nil, err, nil
    }

    for id, subCfg := range cfg {
        errs := subCfg.Validate()
        if len(*errs) > 0 {
            for _, e := range(*errs) {
                msg := fmt.Sprintf("DSLConfig[%s]: %v", id, e)
                errorsOutput = append(errorsOutput, errors.New(msg))
            }

        }
    }
    if len(errorsOutput) > 0 {
        return nil, errors.New("Failed to validate MultiDSLConfig"), &errorsOutput
    }

    return cfg, nil, nil
}


func (cfg MultiDSLConfig) String() string {
    str, _ := yaml.Marshal(cfg)

    return string(str)
}
