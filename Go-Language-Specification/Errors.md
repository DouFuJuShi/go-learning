# Errors

预先声明的类型错误定义为

```go
type error interface {
    Error() string
}
```

这是表示错误条件的常规接口，nil 值表示无错误。例如，可以定义一个从文件中读取数据的函数：

```go

```
