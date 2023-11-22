# Compile and run Go program

```shell
go run [build flags] [-exec xprog] package [arguments...]
```

Run 编译并运行指定的 Go main 包。通常，包被指定为来自单个目录的 .go 源文件列表，但它也可能是导入路径、文件系统路径或与单个已知包匹配的模式，如“go run”。或“运行我的/cmd”。

如果包参数有版本后缀（如@latest或@v1.0.0），“go run”将以模块感知模式构建程序，忽略当前目录或任何父目录中的go.mod文件（如果有）一。这对于运行程序而不影响主模块的依赖关系很有用。

如果包参数没有版本后缀，“go run”可能会在模块感知模式或 GOPATH 模式下运行，具体取决于 GO111MODULE 环境变量和 go.mod 文件是否存在。有关详细信息，请参阅“转到帮助模块”。如果启用模块感知模式，“go run”将在主模块的上下文中运行。

默认情况下，“go run”直接运行编译后的二进制文件：“a.out argument...”。如果给出了 -exec 标志，“go run”将使用以下命令调用二进制文件 xprog:

```shell
'xprog a.out arguments...'.
```

如果未给出 -exec 标志，则 GOOS 或 GOARCH 与系统默认值不同，并且可以在当前搜索路径上找到名为 go_$GOOS_$GOARCH_exec 的程序，“go run”使用该程序调用二进制文件，例如'go_js_wasm_exec a.out 参数...'。当模拟器或其他执行方法可用时，这允许执行交叉编译的程序。
