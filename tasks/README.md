# tasks

添加定时任务请先了解cron语法

添加方法

```go
Cron.AddFunc(yourSpec,func() {Run(yourFunc)})
```

```go
func yourFunc() error{
    // todo
}
```
