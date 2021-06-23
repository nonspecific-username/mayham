package errors

type badRequest struct {
    Object string `yaml:"obj" json:"obj"`
    Value string `yaml:"value" json:"value"`
    Description string `yaml:"desc" json:"desc"`
}


type notFound struct {
    Object string `yaml:"obj" json:"obj"`
    Id string `yaml:"id" json:"id"`
    Description string `yaml:"desc" json:"desc"`
}


func InvalidValue(obj string, value string) *badRequest {
    return &badRequest{Object: obj,
                       Value: value,
                       Description: "Invalid value"}
}


func NotFound(obj string, objId string) *notFound {
    return &notFound{Object: obj,
                     Id: objId,
                     Description: "Not found"}
}


func ParseError(obj string, value string) *badRequest {
    return &badRequest{Object: obj,
                       Value: value,
                       Description: "Parse error"}
}


func UnsupportedContentType(value string) *badRequest {
    return &badRequest{Object: "Content-Type",
                       Value: value,
                       Description: "Unsupported content-type"}
}


func NoSuchField(obj string, value string) *badRequest {
    return &badRequest{Object: obj,
                       Value: value,
                       Description: "No such field"}
}
