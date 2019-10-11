# tasks

添加定时任务请先了解cron语法

添加方法

```
Cron.AddFunc(yourSpec,func() {Run(yourFunc)})
```

```
func yourFunc() error{
    // todo
}
```
