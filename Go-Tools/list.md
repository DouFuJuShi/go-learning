# List packages or modules

```shell
go list [-f format] [-json] [-m] [list flags] [build flags] [packages]
```

列表列出了已命名的包，每行一个。最常用的标志是 -f 和 -json，它们控制为每个包打印的输出的形式。下面记录的其他列表标志控制更具体的细节。

默认输出显示包导入路径：

```
bytes
encoding/json
github.com/gorilla/mux
golang.org/x/net/html
```

<mark>-f 标志</mark>使用包模板的语法指定列表的替代格式。默认输出相当于 -f '{{.ImportPath}}'。传递给模板的结构是：

```go
type Package struct {
    Dir            string   // directory containing package sources
    ImportPath     string   // import path of package in dir
    ImportComment  string   // path in import comment on package statement
    Name           string   // package name
    Doc            string   // package documentation string
    Target         string   // install path
    Shlib          string   // the shared library that contains this package (only set when -linkshared)
    Goroot         bool     // is this package in the Go root?
    Standard       bool     // is this package part of the standard Go library?
    Stale          bool     // would 'go install' do anything for this package?
    StaleReason    string   // explanation for Stale==true
    Root           string   // Go root or Go path dir containing this package
    ConflictDir    string   // this directory shadows Dir in $GOPATH
    BinaryOnly     bool     // binary-only package (no longer supported)
    ForTest        string   // package is only for use in named test
    Export         string   // file containing export data (when using -export)
    BuildID        string   // build ID of the compiled package (when using -export)
    Module         *Module  // info about package's containing module, if any (can be nil)
    Match          []string // command-line patterns matching this package
    DepOnly        bool     // package is only a dependency, not explicitly listed
    DefaultGODEBUG string  // default GODEBUG setting, for main packages

    // Source files
    GoFiles           []string   // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
    CgoFiles          []string   // .go source files that import "C"
    CompiledGoFiles   []string   // .go files presented to compiler (when using -compiled)
    IgnoredGoFiles    []string   // .go source files ignored due to build constraints
    IgnoredOtherFiles []string // non-.go source files ignored due to build constraints
    CFiles            []string   // .c source files
    CXXFiles          []string   // .cc, .cxx and .cpp source files
    MFiles            []string   // .m source files
    HFiles            []string   // .h, .hh, .hpp and .hxx source files
    FFiles            []string   // .f, .F, .for and .f90 Fortran source files
    SFiles            []string   // .s source files
    SwigFiles         []string   // .swig files
    SwigCXXFiles      []string   // .swigcxx files
    SysoFiles         []string   // .syso object files to add to archive
    TestGoFiles       []string   // _test.go files in package
    XTestGoFiles      []string   // _test.go files outside package

    // Embedded files
    EmbedPatterns      []string // //go:embed patterns
    EmbedFiles         []string // files matched by EmbedPatterns
    TestEmbedPatterns  []string // //go:embed patterns in TestGoFiles
    TestEmbedFiles     []string // files matched by TestEmbedPatterns
    XTestEmbedPatterns []string // //go:embed patterns in XTestGoFiles
    XTestEmbedFiles    []string // files matched by XTestEmbedPatterns

    // Cgo directives
    CgoCFLAGS    []string // cgo: flags for C compiler
    CgoCPPFLAGS  []string // cgo: flags for C preprocessor
    CgoCXXFLAGS  []string // cgo: flags for C++ compiler
    CgoFFLAGS    []string // cgo: flags for Fortran compiler
    CgoLDFLAGS   []string // cgo: flags for linker
    CgoPkgConfig []string // cgo: pkg-config names

    // Dependency information
    Imports      []string          // import paths used by this package
    ImportMap    map[string]string // map from source import to ImportPath (identity entries omitted)
    Deps         []string          // all (recursively) imported dependencies
    TestImports  []string          // imports from TestGoFiles
    XTestImports []string          // imports from XTestGoFiles

    // Error information
    Incomplete bool            // this package or a dependency has an error
    Error      *PackageError   // error loading package
    DepsErrors []*PackageError // errors loading dependencies
}
```

存储在vendor目录中的包报告包含vendor目录路径的 ImportPath（例如，“d/vendor/p”而不是“p”），以便 ImportPath 唯一标识包的给定副本。 Imports、Deps、TestImports 和 XTestImports 列表也包含这些扩展的导入路径。有关供应商的更多信息，请参阅 golang.org/s/go15vendor。

如果有错误信息，则为

```go
type PackageError struct {
    ImportStack   []string // shortest path from package named on command line to this one
    Pos           string   // position of error (if present, file:line:col)
    Err           string   // the error itself
}
```

模块信息是一个 Module 结构，定义见下文关于 list -m 的讨论。

模板函数“join”调用 strings.Join。

模板函数“context”返回构建上下文，定义为：

```go
type Context struct {
    GOARCH        string   // target architecture
    GOOS          string   // target operating system
    GOROOT        string   // Go root
    GOPATH        string   // Go path
    CgoEnabled    bool     // whether cgo can be used
    UseAllFiles   bool     // use files regardless of //go:build lines, file names
    Compiler      string   // compiler to assume when computing target paths
    BuildTags     []string // build constraints to match in //go:build lines
    ToolTags      []string // toolchain-specific build constraints
    ReleaseTags   []string // releases the current release is compatible with
    InstallSuffix string   // suffix to use in the name of the install dir
}
```

有关这些字段含义的更多信息，请参阅 go/build 包的 Context 类型的文档。

<mark>-json</mark> 标志导致包数据以 JSON 格式打印，而不是使用模板格式。可以选择为 JSON 标志提供一组要输出的以逗号分隔的必需字段名称。如果是这样，这些必填字段将始终出现在 JSON 输出中，但其他字段可能会被省略以节省计算 JSON 结构的工作。

<mark>-compiled</mark> 标志导致 list 将 CompiledGoFiles 设置为提供给编译器的 Go 源文件。通常，这意味着它会重复 GoFiles 中列出的文件，然后还添加通过处理 CgoFiles 和 SwigFiles 生成的 Go 代码。 Imports 列表包含来自 GoFiles 和 CompiledGoFiles 的所有导入的并集。

<mark>-deps</mark> 标志导致 list 不仅迭代指定的包，还迭代它们的所有依赖项。它以深度优先的后序遍历方式访问它们，以便仅在其所有依赖项之后列出包。命令行上未明确列出的包会将 DepOnly 字段设置为 true。

<mark>-e</mark> 标志更改了对错误包（无法找到或格式错误的包）的处理。默认情况下，list 命令将每个错误包的错误打印到标准错误，并在通常的打印过程中忽略这些包。使用 -e 标志，list 命令永远不会将错误打印到标准错误，而是使用通常的打印来处理错误的包。错误的包将有一个非空的 ImportPath 和一个非零的 Error 字段；其他信息可能会或可能不会丢失（归零）。

<mark>-export</mark> 标志使 list 将 Export 字段设置为包含给定包的最新导出信息的文件的名称，并将 BuildID 字段设置为已编译包的构建 ID。

<mark>-find</mark> 标志使 list 识别指定的包，但不解析它们的依赖关系：Imports 和 Deps 列表将为空。使用 -find 标志时，无法使用 -deps、-test 和 -export 命令。

<mark>-test</mark> 标志使 list 不仅报告命名的包，还报告它们的测试二进制文件（对于带有测试的包），以准确地向源代码分析工具传达测试二进制文件的构造方式。测试二进制文件的报告导入路径是包的导入路径，后跟“.test”后缀，如“math/rand.test”。构建测试时，有时需要专门为该测试重建某些依赖项（最常见的是被测试的包本身）。为特定测试二进制文件重新编译的包的报告导入路径后跟一个空格和括号中的测试二进制文件的名称，如“math/rand math/rand.test”或“regexp [sort.test]”。 ForTest 字段还设置为正在测试的包的名称（前面示例中的“math/rand”或“sort”）。

Dir、Target、Shlib、Root、ConflictDir 和 Export 文件路径都是绝对路径。

默认情况下，GoFiles、CgoFiles 等列表保存 Dir 中的文件名称（即相对于 Dir 的路径，而不是绝对路径）。使用 -compiled 和 -test 标志时添加的生成文件是引用生成的 Go 源文件的缓存副本的绝对路径。虽然它们是Go源文件，但路径可能不以“.go”结尾。

<mark>-m</mark> 标志使 list 列出模块而不是包。

列出模块时， -f 标志仍然指定应用于 Go 结构的格式模板，但现在是 Module 结构：

```go
type Module struct {
    Path       string        // module path
    Query      string        // version query corresponding to this version
    Version    string        // module version
    Versions   []string      // available module versions
    Replace    *Module       // replaced by this module
    Time       *time.Time    // time version was created
    Update     *Module       // available update (with -u)
    Main       bool          // is this the main module?
    Indirect   bool          // module is only indirectly needed by main module
    Dir        string        // directory holding local copy of files, if any
    GoMod      string        // path to go.mod file describing module, if any
    GoVersion  string        // go version used in module
    Retracted  []string      // retraction information, if any (with -retracted or -u)
    Deprecated string        // deprecation message, if any (with -u)
    Error      *ModuleError  // error loading module
    Origin     any           // provenance of module
    Reuse      bool          // reuse of old module info is safe
}

type ModuleError struct {
    Err string // the error itself
}
```

如果模块位于模块缓存中或者使用了 -modfile 标志，则 GoMod 引用的文件可能位于模块目录之外。
默认输出是打印模块路径，然后打印有关版本和替换（如果有）的信息。例如，“go list -m all”可能会打印：

```go
my/main/module
golang.org/x/text v0.3.0 => /tmp/text
rsc.io/pdf v0.1.1
```

Module 结构有一个 String 方法，用于格式化此行输出，因此默认格式相当于 -f '{{.String}}'。

请注意，当模块被替换时，其 Replace 字段描述替换模块，并且其 Dir 字段设置为替换的源代码（如果存在）。 （也就是说，如果 Replace 为非零，则 Dir 设置为 Replace.Dir，无法访问替换的源代码。）

<mark>-u </mark>标志添加有关可用升级的信息。当给定模块的最新版本比当前版本更新时，list -u 将模块的更新字段设置为有关较新模块的信息。如果当前版本被撤回，list -u 还将设置模块的“Retracted”字段。模块的 String 方法通过在当前版本后面的括号中格式化较新的版本来指示可用的升级。如果版本被撤回，则字符串“(retracted)”将跟随在其后面。例如，“go list -m -u all”可能会打印：

```go
my/main/module
golang.org/x/text v0.3.0 [v0.4.0] => /tmp/text
rsc.io/pdf v0.1.1 (retracted) [v0.1.2]
```

（对于工具，“go list -m -u -json all”可能更方便解析。）

<mark>-versions</mark> 标志使 list 将模块的版本字段设置为该模块的所有已知版本的列表，根据语义版本控制从最早到最新排序。该标志还更改默认输出格式以显示模块路径，后跟空格分隔的版本列表。

<mark>-retracted</mark> 标志使 list 报告有关收回的模块版本的信息。当 -retracted 与 -f 或 -json 一起使用时，Retracted 字段将设置为一个字符串，解释版本被撤回的原因。该字符串取自模块的 go.mod 文件中的撤回指令的注释。当 -retracted 与 -versions 一起使用时，收回的版本将与未收回的版本一起列出。 -retracted 标志可以与 -m 一起使用，也可以不与 -m 一起使用。

list -m 的参数被解释为模块列表，而不是包列表。主模块是包含当前目录的模块。活动模块是主模块及其依赖项。如果没有参数，list -m 显示主模块。使用参数时，list -m 显示参数指定的模块。任何活动模块都可以通过其模块路径来指定。特殊模式“all”指定所有活动模块，首先是主模块，然后是按模块路径排序的依赖项。包含“...”的模式指定其模块路径与该模式匹配的活动模块。 path@version 形式的查询指定该查询的结果，该结果不限于活动模块。有关模块查询的更多信息，请参阅“转到帮助模块”。

模板函数“module”采用单个字符串参数，该参数必须是模块路径或查询，并将指定的模块作为 Module 结构返回。如果发生错误，结果将是一个带有非零错误字段的模块结构。

使用 -m 时，-reuse=old.json 标志接受包含先前“go list -m -json”调用的 JSON 输出的文件名，该调用具有相同的一组修饰符标志（例如 -u、-retracted、和-版本）。 go 命令可以使用此文件来确定模块自上次调用以来未发生更改，并避免重新下载有关该模块的信息。通过将 Reuse 字段设置为 true，将在新输出中标记未重新下载的模块。通常模块缓存会自动提供这种重用； -reuse 标志对于不保留模块缓存的系统很有用。

有关构建标志的更多信息，请参阅“帮助构建”。

有关指定包的更多信息，请参阅“转到帮助包”。
有关模块的更多信息，请参阅 https://golang.org/ref/mod。

使用案例：

1. 列出某个库所有的版本

```shell
go list -m -versions github.com/gin-gonic/gin
```
