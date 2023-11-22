# [Generate Go files by processing source](https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source)

```go
go generate [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]
```

生成现有文件中的指令描述的运行命令。这些命令可以运行任何进程，但目的是创建或更新 Go 源文件。

Go generate 永远不会由 gobuild、gotest 等自动运行。它必须显式运行。

Go generate 扫描文件中的指令，这些指令是以下形式的行，

```go
//go:generate command argument...
```

（注意：前导不能有空格，“//go”中不能有空格）其中command是要运行的生成器，对应一个可以在本地运行的可执行文件。它必须位于 shell 路径 (gofmt)、完全限定路径 (/usr/you/bin/mytool) 或命令别名中，如下所述。

请注意，go generate 不会解析文件，因此注释或多行字符串中看起来像指令的行会被视为指令。

该指令的参数是空格分隔的标记或双引号字符串，在运行时作为单独的参数传递给生成器。

带引号的字符串使用 Go 语法并在执行前进行评估；带引号的字符串显示为生成器的单个参数。

为了向人类和机器传达代码已生成，生成的源代码应具有与以下正则表达式匹配的行（采用 Go 语法）：

```go
^// Code generated .* DO NOT EDIT\.$
```

该行必须出现在文件中第一个非注释、非空白文本之前。

Go generate 在运行发电机时会设置几个变量：

```go
$GOARCH
	The execution architecture (arm, amd64, etc.)
$GOOS
	The execution operating system (linux, windows, etc.)
$GOFILE
	The base name of the file.
$GOLINE
	The line number of the directive in the source file.
$GOPACKAGE
	The name of the package of the file containing the directive.
$GOROOT
	The GOROOT directory for the 'go' command that invoked the
	generator, containing the Go toolchain and standard library.
$DOLLAR
	A dollar sign.
$PATH
	The $PATH of the parent process, with $GOROOT/bin
	placed at the beginning. This causes generators
	that execute 'go' commands to use the same 'go'
	as the parent 'go generate' command.
```

除了变量替换和带引号的字符串求值之外，命令行上不会执行“通配符”等特殊处理。

作为运行命令之前的最后一步，对具有字母数字名称的任何环境变量（例如 $GOFILE 或 $HOME）的任何调用都会在整个命令行中展开。在所有操作系统上，变量扩展的语法都是 $NAME。由于计算的顺序，变量甚至在带引号的字符串内也会扩展。如果未设置变量 NAME，$NAME 将扩展为空字符串。

形式的指令，

```go
//go:generate -command xxx args...
```

仅针对该源文件的其余部分指定字符串 xxx 表示由参数标识的命令。这可用于创建别名或处理多字生成器。例如，

```go
//go:generate -command foo go tool foo
```

指定命令“foo”代表生成器“go tool foo”。
按照命令行上给出的顺序生成进程包，一次一个。如果命令行列出单个目录中的 .go 文件，它们将被视为单个包。在包内，generate 按文件名顺序处理包中的源文件，一次一个。在源文件中，generate 按照生成器在文件中出现的顺序运行生成器，一次一个。 gogenerate 工具还设置构建标签“generate”，以便文件可以由 gogenerate 检查但在构建过程中被忽略。
对于具有无效代码的包，生成仅处理具有有效包子句的源文件。
如果任何生成器返回错误退出状态，“gogenerate”将跳过该包的所有进一步处理。
生成器在包的源目录中运行。
Gogenerate 接受两个特定的标志：

```go

```
