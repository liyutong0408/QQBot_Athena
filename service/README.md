# service

事件处理

如果有想添加的新事件请在此编写并在api内注册

## example

``` go

type ServiceName struct {
    // 为了格式整齐美观
    // (虽然我自己写的也不漂亮
}

func (service serviceName) FuncName(ch chan bool,framework model.Framework) {
    // todo 判断条件,事件处理


     // 如果处理事件 请向ch中传入true
     // 如果掠过事件 请向ch中传入false 并结束函数
     // ch用来判断是否有服务对消息进行处理,然而只有复读功能用到了这个，毕竟不能复读自己的命令
}

```
