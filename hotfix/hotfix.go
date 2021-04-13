package hotfix


import (
    "fmt"
)


type Hotfix interface {
    Render() string
}


type Regular struct {
    Method string
    Notify int
    Pkg string
    Object string
    Attr string
    FromLen int
    From string
    Value string
}


type DataTable struct {
    Method string
    Notify int
    Pkg string
    Object string
    Row string
    Attr string
    FromLen int
    From string
    Value string
}


func (hf Regular) Render() string {
    return fmt.Sprintf("%s,(1,1,%d,%s),%s,%s,%d,%s,%s", hf.Method, hf.Notify, hf.Pkg, hf.Object, hf.Attr, hf.FromLen, hf.From, hf.Value)
}


func RenderBVCOverride(number int) string {
    var template string = "(ValueType=Int,DisabledValueModes=102,ValueFlags=0,ValueMode=AttributeInitializationData,Range=(Value=1.000000,Variance=0.000000),AttributeInitializer=None,AttributeData=None,AttributeInitializationData=(BaseValueConstant=%d.000000,DataTableValue=(DataTable=None,RowName=\"\",ValueName=\"\"),BaseValueAttribute=None,AttributeInitializer=None,BaseValueScale=1.000000),BlackboardKey=(KeyName=\"\",bRuntimeKey=False),Condition=None,Actor=None)"
    return fmt.Sprintf(template, number)
}
