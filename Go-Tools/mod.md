# 模块维护 Module maintenance

Go mod 提供对模块操作的访问。

请注意，对模块的支持内置于所有 go 命令中，而不仅仅是“go mod”。例如，日常添加、删除、升级和降级依赖项应该使用“go get”来完成。有关模块功能的概述，请参阅“转到帮助模块”。

用法

```go
go mod <command> [arguments]
```

命令是：

```go
download    download modules to local cache
edit        edit go.mod from tools or scripts
graph       print module requirement graph
init        initialize new module in current directory
tidy        add missing and remove unused modules
vendor      make vendored copy of dependencies
verify      verify dependencies have expected content
why         explain why packages or modules are needed
```

使用“go help mod <command>”获取有关命令的更多信息。

## 下载模块到本地缓存 Download modules to local cache

```go
go mod download [-x] [-json] [-reuse=old.json] [modules]
```

Dowload下载指定的模块，这些模块可以是选择主模块依赖关系的模块模式，也可以是 path@version 形式的模块查询。

如果不带参数，下载适用于在主模块中构建和测试包所需的模块：主模块明确需要的模块（如果是“go 1.17”或更高版本），或者所有传递所需的模块（如果是“go”） 1.16' 或更低。

go命令在普通执行过程中会根据需要自动下载模块。 “go mod download”命令主要用于预填充本地缓存或计算 Go 模块代理的答案。

默认情况下，下载不会向标准输出写入任何内容。它可以将进度消息和错误打印到标准错误。

<mark>-json 标志</mark>使 download 将一系列 JSON 对象打印到标准输出，描述每个下载的模块（或失败），对应于以下 Go 结构：

```go
type Module struct {
    Path     string // module path
    Query    string // version query corresponding to this version
    Version  string // module version
    Error    string // error loading module
    Info     string // absolute path to cached .info file
    GoMod    string // absolute path to cached .mod file
    Zip      string // absolute path to cached .zip file
    Dir      string // absolute path to cached source root directory
    Sum      string // checksum for path, version (as in go.sum)
    GoModSum string // checksum for go.mod (as in go.sum)
    Origin   any    // provenance of module
    Reuse    bool   // reuse of old module info is safe
}
```

<mark>-reuse 标志</mark>接受包含先前“go mod download -json”调用的 JSON 输出的文件名。 go 命令可以使用此文件来确定模块自上次调用以来未发生更改，并避免重新下载它。通过将 Reuse 字段设置为 true，将在新输出中标记未重新下载的模块。通常模块缓存会自动提供这种重用； -reuse 标志对于不保留模块缓存的系统很有用。

-x 标志使 download 打印 download 执行的命令。

## Edit go.mod from tools or scripts

```go
go mod edit [editing flags] [-fmt|-print|-json] [go.mod]
```

Edit 提供了用于编辑 go.mod 的命令行界面，主要供工具或脚本使用。它只读取 go.mod；它不会查找有关所涉及模块的信息。默认情况下，edit读取和写入主模块的go.mod文件，但可以在编辑标志后指定不同的目标文件。

编辑标志指定编辑操作的序列。

<mark>-fmt 标志</mark>重新格式化 go.mod 文件而不进行其他更改。使用或重写 go.mod 文件的任何其他修改也暗示了这种重新格式化。唯一需要此标志的情况是没有指定其他标志，如“go mod edit -fmt”中。

<mark>-module</mark> 标志更改模块的路径（go.mod 文件中的模块行）。

<mark></mark>-require=path@version 和 -droprequire=path 标志添加和删除对给定模块路径和版本的要求。请注意，-require 会覆盖路径上的任何现有要求。这些标志主要用于理解模块图的工具。用户应该更喜欢“go get path@version”或“go get path@none”，这会根据需要进行其他 go.mod 调整以满足其他模块施加的约束。

<mark>-exclude=path@version 和 -dropexclude=path@version 标志</mark>添加和删除给定模块路径和版本的排除项。请注意，如果该排除已存在，则 -exclude=path@version 是无操作。

<mark>-replace=old[@v]=new[@v] </mark>标志添加给定模块路径和版本对的替换。如果old@v中的@v被省略，则会添加左侧没有版本的替换，适用于旧模块路径的所有版本。如果new@v中的@v被省略，则新路径应该是本地模块根目录，而不是模块路径。请注意，-replace 会覆盖旧[@v] 的任何冗余替换，因此省略 @v 将删除特定版本的现有替换。

<mark>-dropreplace=old[@v] 标志</mark>删除给定模块路径和版本对的替换。如果省略@v，则左侧没有版本的替换将被删除。

-retract=version 和 -dropretract=version 标志添加和删除给定版本的撤消。该版本可以是单个版本（如“v1.2.3”）或闭区间（如“[v1.1.0，v1.1.9]”）。请注意，如果该撤消已经存在，则 -retract=version 是无操作。

-require、-droprequire、-exclude、-dropexclude、-replace、-dropreplace、-retract 和 -dropretract 编辑标志可以重复，并且更改将按给定的顺序应用。

-go=version 标志设置预期的 Go 语言版本。

-toolchain=name 标志设置要使用的 Go 工具链。

-print 标志以文本格式打印最终的 go.mod，而不是将其写回 go.mod。

-json 标志以 JSON 格式打印最终的 go.mod 文件，而不是将其写回 go.mod。 JSON 输出对应于以下 Go 类型：

```go
type Module struct {
	Path    string
	Version string
}

type GoMod struct {
	Module    ModPath
	Go        string
	Toolchain string
	Require   []Require
	Exclude   []Module
	Replace   []Replace
	Retract   []Retract
}

type ModPath struct {
	Path       string
	Deprecated string
}

type Require struct {
	Path string
	Version string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}

type Retract struct {
	Low       string
	High      string
	Rationale string
}
```

代表单个版本（不是间隔）的撤回条目将“低”和“高”字段设置为相同的值。
请注意，这仅描述了 go.mod 文件本身，而不是间接引用的其他模块。对于可用于构建的完整模块集，请使用“go list -m -json all”。



## 打印模块需求图  Print module requirement graph

```go
go mod graph [-go=version] [-x]
```

Graph 以文本形式打印模块需求图（已应用替换）。输出中的每一行都有两个以空格分隔的字段：一个模块及其要求之一。每个模块都被标识为路径@版本形式的字符串，但主模块除外，它没有@版本后缀。

-go 标志使 graph 报告给定 Go 版本加载的模块图，而不是 go.mod 文件中“go”指令指示的版本。

-x 标志使 graph 打印 graph 执行的命令。

## 在当前目录下新建一个模块 Initialize new module in current directory

```go
go mod init [module-path]
```

Init 初始化并在当前目录中写入一个新的 go.mod 文件，实际上创建了一个以当前目录为根的新模块。 go.mod 文件不得已存在。

Init 接受一个可选参数，即新模块的模块路径。如果省略模块路径参数，init 将尝试使用 .go 文件中的导入注释、供应商工具配置文件（如 Gopkg.lock）和当前目录（如果在 GOPATH 中）来推断模块路径。

## 添加缺失的模块并删除未使用的模块  Add missing and remove unused modules

```go
go mod tidy [-e] [-v] [-x] [-go=version] [-compat=version]
```

Tidy 确保 go.mod 与模块中的源代码匹配。它添加构建当前模块的包和依赖项所需的任何缺失模块，并删除不提供任何相关包的未使用模块。它还将所有缺失的条目添加到 go.sum 并删除任何不必要的条目。

-v 标志使 tidy 将有关已删除模块的信息打印到标准错误。

-e 标志会导致 tidy 尝试继续，尽管在加载包时遇到错误。

-go 标志导致 tidy 将 go.mod 文件中的“go”指令更新为给定版本，这可能会更改在 go.mod 文件中保留为显式要求的模块依赖项。 （Go 版本 1.17 及更高版本保留更多要求以支持延迟模块加载。）

-compat 标志保留来自指定主要 Go 版本的“go”命令所需的任何附加校验和，以成功加载模块图，并且如果该版本的“go”命令将从某个版本加载任何导入的包，则会导致 tidy 出错。不同的模块版本。默认情况下，tidy 的行为就好像 -compat 标志设置为 go.mod 文件中“go”指令指示的版本之前的版本。

-x 标志使 tidy 打印下载执行的命令。

## 制作依赖项的附带副本 Make vendored copy of dependencies

```go
go mod vendor [-e] [-v] [-o outdir]
```

供应商重置主模块的供应商目录以包含构建和测试所有主模块的包所需的所有包。它不包括供应包的测试代码。

-v 标志使供应商将供应的模块和包的名称打印到标准错误。

-e 标志导致供应商尝试继续，尽管在加载包时遇到错误。

-o 标志使供应商在给定路径而不是“vendor”处创建供应商目录。 go 命令只能使用模块根目录中名为“vendor”的供应商目录，因此该标志主要对其他工具有用。

## 验证依赖项是否具有预期内容  Verify dependencies have expected content

```go
go mod verify
```

Verify 检查存储在本地下载源缓存中的当前模块的依赖项在下载后是否被修改。如果所有模块都未修改，verify 会打印 "所有模块已验证"。否则，它会报告哪些模块被修改，并导致 "go mod "以非零状态退出。

## 解释为什么需要包或模块  why packages or modules are needed

```go
go mod why [-m] [-vendor] packages...
```

why 在导入图中显示从主模块到每个列出软件包的最短路径。如果给定了 -m 标志，why 会将参数视为模块列表，并查找每个模块中任何软件包的路径。

默认情况下，why 查询与“go list all”匹配的包的图表，其中包括对可达包的测试。 -vendor 标志导致为什么要排除依赖项测试。

输出是一系列节，每个节对应命令行上的每个包或模块名称，由空行分隔。每个节以注释行“# package”或“# module”开头，给出目标包或模块。后续行给出了通过导入图的路径，每行一个包。如果主模块未引用包或模块，则该节将显示一个带括号的注释来指示这一事实。

案例：

```shell
$ go mod why golang.org/x/text/language golang.org/x/text/encoding
# golang.org/x/text/language
rsc.io/quote
rsc.io/sampler
golang.org/x/text/language

# golang.org/x/text/encoding
(main module does not need package golang.org/x/text/encoding)
$
```
