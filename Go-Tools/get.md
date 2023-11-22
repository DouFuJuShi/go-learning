# Add dependencies to current module and install them

```shell
go get [-t] [-u] [-v] [build flags] [packages]
```

Get 将其命令行参数解析为特定模块版本的包，更新 go.mod 以要求这些版本，并将源代码下载到模块缓存中。

要添加包的依赖项或将其升级到最新版本：

```go
go get example.com/pkg
```

要将软件包升级或降级到特定版本：

```shell
go get example.com/pkg@v1.2.3
```

要将最低所需的 Go 版本升级到最新发布的 Go 版本：

```shell
go get go@latest
```

将 Go 工具链升级到当前 Go 工具链的最新补丁版本：

```shell
go get toolchain@patch
```

有关详细信息，请参阅 https://golang.org/ref/mod#go-get。

在 Go 的早期版本中，“go get”用于构建和安装包。现在，“go get”专门用于调整 go.mod 中的依赖关系。 “go install”可用于构建和安装命令。指定版本后，“go install”将以模块感知模式运行并忽略当前目录中的 go.mod 文件。例如：

```shell
go install example.com/pkg@v1.2.3
go install example.com/pkg@latest
```

有关详细信息，请参阅“go help install”或 https://golang.org/ref/mod#go-install。

“go get”接受以下标志。

-t 标志指示 get 考虑构建命令行上指定的包的测试所需的模块。

-u 标志指示 get 更新提供命令行上命名的包的依赖项的模块，以使用更新的次要版本或补丁版本（如果可用）。

-u=patch 标志（不是 -u patch）还指示 get 更新依赖项，但更改默认值以选择补丁版本。

当 -t 和 -u 标志一起使用时， get 也会更新测试依赖项。

-x 标志在执行命令时打印命令。当直接从存储库下载模块时，这对于调试版本控制命令非常有用。

有关模块的更多信息，请参阅 https://golang.org/ref/mod。

有关使用“go get”更新最低 Go 版本和建议的 Go 工具链的更多信息，请参阅 https://go.dev/doc/toolchain。

有关指定包的更多信息，请参阅“转到帮助包”。

本文描述了 get 使用模块来管理源代码和依赖项的行为。相反，如果 go 命令在 GOPATH 模式下运行，则 get 的标志和效果的详细信息会发生变化，“go help get”也会发生变化。请参阅“go help gopath-get”。
另请参阅：go build、go install、go clean、go mod。
