# Compile and install packages and dependencies

```go
go install [build flags] [packages]
```

Install 编译并安装由导入路径命名的包。

可执行文件安装在 GOBIN 环境变量指定的目录中，如果未设置 GOPATH 环境变量，则默认为 \$GOPATH/bin 或 $HOME/go/bin。 $\$GOROOT 中的可执行文件安装在 \$GOROOT/bin 或 \$GOTOOLDIR 而不是 $GOBIN 中

如果参数具有版本后缀（如 @latest 或 @v1.0.0），“go install”将以模块感知模式构建软件包，忽略当前目录或任何父目录（如果有）中的 go.mod 文件。这对于安装可执行文件而不影响主模块的依赖关系很有用。为了消除构建中使用哪些模块版本的歧义，参数必须满足以下约束：

- 参数必须是包路径或包模式（带有“...”通配符）。它们不能是标准包（如 fmt）、元模式（std、cmd、all）或相对或绝对文件路径。

- 所有参数必须具有相同的版本后缀。不允许不同的查询，即使它们引用相同的版本。

- 所有参数必须引用同一版本的同一模块中的包。

- 包路径参数必须引用主包。模式参数仅匹配主包。

- 没有模块被视为“主”模块。如果包含在命令行上命名的包的模块具有 go.mod 文件，则它不得包含会导致其解释方式与主模块不同的指令（替换和排除）。该模块不能要求其自身的更高版本。

- 任何模块中都不使用供应商目录。 （供应商目录不包含在“go install”下载的模块 zip 文件中。）

如果参数没有版本后缀，“go install”可能会在模块感知模式或 GOPATH 模式下运行，具体取决于 GO111MODULE 环境变量和 go.mod 文件是否存在。有关详细信息，请参阅“转到帮助模块”。如果启用模块感知模式，“go install”将在主模块的上下文中运行。

当禁用模块感知模式时，非主包将安装在目录 $GOPATH/pkg/$GOOS_$GOARCH 中。启用模块感知模式后，将构建并缓存非主包，但不会安装。

在Go 1.20之前，标准库安装到$GOROOT/pkg/$GOOS_\$GOARCH。从 Go 1.20 开始，标准库已构建并缓存，但未安装。设置 GODEBUG=installgoroot=all 可恢复 $GOROOT/pkg/$GOOS_$GOARCH 的使用。
