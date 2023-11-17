# 运行时Panics

执行错误（例如尝试对数组进行索引越界）会触发运行时恐慌，相当于使用实现定义的接口类型runtime.Error 的值调用内置函数恐慌。该类型满足预先声明的接口类型错误。表示不同运行时错误条件的确切错误值未指定。

```go
package runtime

type Error interface {
    error
    // and perhaps other methods
}
```
