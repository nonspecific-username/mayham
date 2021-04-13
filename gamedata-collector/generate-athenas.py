import yaml
import sys

with open(sys.argv[1], 'r') as f:
    data = yaml.load(f)


def render_bvc(num):
    tpl = '(ValueType=Int,DisabledValueModes=102,ValueFlags=0,ValueMode=AttributeInitializationData,Range=(Value=1.000000,Variance=0.000000),AttributeInitializer=None,AttributeData=None,AttributeInitializationData=(BaseValueConstant={bvc}.000000,DataTableValue=(DataTable=None,RowName="",ValueName=""),BaseValueAttribute=None,AttributeInitializer=None,BaseValueScale=1.000000),BlackboardKey=(KeyName="",bRuntimeKey=False),Condition=None,Actor=None)'
    return tpl.format(bvc=num)


def render_patch(spawner, attrib, num):
    output = ['SparkEarlyLevelPatchEntry',
              '(1,1,0,Monastery_P)',
              spawner['Path'],
              spawner['AttrBase'] + '.' + attrib,
              '0',
              '',
              render_bvc(int(num))]
    return ','.join(output)


def render_mult(spawner, factor):
    nap = spawner['NumActorsParam'] * factor
    print(render_patch(spawner, 'NumActorsParam', nap))
    print(render_patch(spawner, 'MaxAliveActorsWhenPassive', nap))
    print(render_patch(spawner, 'MaxAliveActorsWhenThreatened', nap))

for sname, spawner in data['Game']['Monastery'].items():
    if spawner['Type'] == 'Single':
        continue
    render_mult(spawner, 10)
