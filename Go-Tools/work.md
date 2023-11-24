# Workspace maintenance

Work 提供对工作区上的操作的访问。

请注意，对工作区的支持内置于许多其他命令中，而不仅仅是“go work”。

有关工作区所属的 Go 模块系统的信息，请参阅“go 帮助模块”。

有关工作区的深入参考，请参阅 https://go.dev/ref/mod#workspaces。

有关工作区的介绍性教程，请参阅 https://go.dev/doc/tutorial/workspaces。

工作空间由 go.work 文件指定，该文件使用“use”指令指定一组模块目录。这些模块被 go 命令用作根模块来进行构建和相关操作。未指定要使用的模块的工作区不能用于从本地模块进行构建。

go.work 文件是面向行的。每行包含一个指令，由关键字和参数组成。例如：

```go-module
go 1.18

use ../foo/bar
use ./baz

replace example.com/foo v1.2.3 => example.com/bar v1.4.5
```

前导关键字可以从相邻行中分解出来以创建一个块，就像在 Go 导入中一样。

```go-module
use (
  ../foo/bar
  ./baz
)
```

use 指令指定要包含在工作区的主模块集中的模块。 use 指令的参数是包含模块的 go.mod 文件的目录。

go 指令指定写入文件的 Go 版本。此版本可以控制的工作空间的语义将来可能会发生变化，但目前指定的版本没有效果。

替换指令与 go.mod 文件中的替换指令具有相同的语法，并且优先于 go.mod 文件中的替换。它的主要目的是覆盖不同工作区模块中的冲突替换。

要确定 go 命令是否在工作区模式下运行，请使用“go env GOWORK”命令。这将指定正在使用的工作区文件。

```go
go work <command> [arguments]
```

```go
edit        edit go.work from tools or scripts
init        initialize workspace file
sync        sync workspace build list to modules
use         add modules to workspace file
```

## Edit go.work from tools or scripts

```go
go work edit [editing flags] [go.work]
```

Edit 提供了用于编辑 go.work 的命令行界面，主要供工具或脚本使用。它只读取 go.work；它不会查找有关所涉及模块的信息。如果未指定文件，Edit 会在当前目录及其父目录中查找 go.work 文件

编辑标志指定编辑操作的序列。

-fmt 标志重新格式化 go.work 文件而不进行其他更改。使用或重写 go.mod 文件的任何其他修改也暗示了这种重新格式化。唯一需要此标志的情况是没有指定其他标志，如“go work edit -fmt”中。

-use=path 和 -dropuse=path 标志从 go.work 文件的模块目录集中添加和删除 use 指令。

-replace=old[@v]=new[@v] 标志添加给定模块路径和版本对的替换。如果old@v中的@v被省略，则会添加左侧没有版本的替换，适用于旧模块路径的所有版本。如果new@v中的@v被省略，则新路径应该是本地模块根目录，而不是模块路径。请注意，-replace 会覆盖旧[@v] 的任何冗余替换，因此省略 @v 将删除特定版本的现有替换。

-dropreplace=old[@v] 标志删除给定模块路径和版本对的替换。如果省略@v，则左侧没有版本的替换将被删除。

-use、-dropuse、-replace 和 -dropreplace 编辑标志可以重复，并且更改将按给定的顺序应用。

-go=version 标志设置预期的 Go 语言版本。

-toolchain=name 标志设置要使用的 Go 工具链。

-print 标志以文本格式打印最终的 go.work，而不是将其写回 go.mod。

-json 标志以 JSON 格式打印最终的 go.work 文件，而不是将其写回 go.mod。 JSON 输出对应于以下 Go 类型：

```go
type GoWork struct {
	Go        string
	Toolchain string
	Use       []Use
	Replace   []Replace
}

type Use struct {
	DiskPath   string
	ModulePath string
}

type Replace struct {
	Old Module
	New Module
}

type Module struct {
	Path    string
	Version string
}
```

#### [Initialize workspace file](https://pkg.go.dev/cmd/go#hdr-Initialize_workspace_file)

```go
go work init [moddirs]
```

Init 初始化并在当前目录中写入一个新的 go.work 文件，实际上是在当前目录中创建一个新的工作空间。

go work init 可以选择接受工作区模块的路径作为参数。如果省略该参数，将创建一个没有模块的空工作区。

每个参数路径都会添加到 go.work 文件中的 use 指令中。当前的 go 版本也将在 go.work 文件中列出。

## 将模块添加到工作区文件  Add modules to workspace file


