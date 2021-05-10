package hotfix


import (
    "fmt"
    "strings"
)


type HotfixMethod string

const (
    EarlyLevel HotfixMethod = "SparkEarlyLevelPatchEntry"
)


type Hotfix struct {
    lines []string
}


func (hf* Hotfix) AddRegular(method HotfixMethod, notify int, pkg string, object string,
                             attr string, fromLen int, from string, value string) {
    hf.lines = append(hf.lines, fmt.Sprintf("%s,(1,1,%d,%s),%s,%s,%d,%s,%s", method, notify, pkg,
                                            object, attr, fromLen, from, value))
}


func (hf* Hotfix) AddRaw(line string) {
    hf.lines = append(hf.lines, line)
}


func (hf* Hotfix) Render() string {
    return strings.Join(hf.lines, "\n")
}


func RenderBVCOverride(number int) string {
    var template string = "(ValueType=Int,DisabledValueModes=102,ValueFlags=0,ValueMode=AttributeInitializationData,Range=(Value=1.000000,Variance=0.000000),AttributeInitializer=None,AttributeData=None,AttributeInitializationData=(BaseValueConstant=%d.000000,DataTableValue=(DataTable=None,RowName=\"\",ValueName=\"\"),BaseValueAttribute=None,AttributeInitializer=None,BaseValueScale=1.000000),BlackboardKey=(KeyName=\"\",bRuntimeKey=False),Condition=None,Actor=None)"
    return fmt.Sprintf(template, number)
}
