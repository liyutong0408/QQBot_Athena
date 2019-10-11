# service

事件处理

如果有想添加的新事件请在此编写并在api内注册

## example

``` go
type ServiceName struct {
}

func (service serviceName) FuncName(ch chan bool,framework model.Framework) {
  // todo 进行判断
  // 如果处理事件 请向ch中传入true
  // 如果掠过事件 请向ch中传入false 并结束函数
}

```
