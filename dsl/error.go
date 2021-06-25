package dsl


import (
    yaml "gopkg.in/yaml.v2"
)


type ValidationError struct {
    Field string `yaml:"field" json:"field"`
    Description string `yaml:"desc" json:"desc"`
    Value interface{} `yaml:"value" json:"value"`
}


func NewValidationError(field string, desc string, value interface{}) *ValidationError {
    return &ValidationError{Field: field,
                            Description: desc,
                            Value: value}
}


func (ve *ValidationError) String() string {
    bytes, _ := yaml.Marshal(ve)
    return string(bytes)
}
