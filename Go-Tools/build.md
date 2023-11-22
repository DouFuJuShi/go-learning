# [Compile packages and dependencies](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)

```shell
go build [-o output] [build flags] [packages]
```

Build 编译由导入路径命名的包及其依赖项，但不安装编译结果。

如果 build 的参数是来自单个目录的 .go 文件列表，则 build 将它们视为指定单个包的源文件列表。

编译包时，build 会忽略以“_test.go”结尾的文件。

编译单个main包时，build 将生成的可执行文件写入以第一个源文件命名的输出文件（“go build ed.go rx.go”写入“ed”或“ed.exe”）或源代码目录（ “go build unix/sam”写入“sam”或“sam.exe”）。编写 Windows 可执行文件时会添加“.exe”后缀。

当编译多个包或单个非主包时， build 会编译这些包，但会丢弃生成的对象，仅用于检查是否可以构建包。

-o 标志强制 build 将生成的可执行文件或对象写入指定的输出文件或目录，而不是最后两段中描述的默认行为。如果指定的输出是现有目录或以斜杠或反斜杠结尾，则任何生成的可执行文件都将写入该目录。

build flags 由 build、clean、get、install、list、run 和 test 命令共享：

```
-C dir
    Change to dir before running the command.
    Any files named on the command line are interpreted after
    changing directories.
    If used, this flag must be the first one in the command line.
-a
    force rebuilding of packages that are already up-to-date.
-n
    print the commands but do not run them.
-p n
    the number of programs, such as build commands or
    test binaries, that can be run in parallel.
    The default is GOMAXPROCS, normally the number of CPUs available.
-race
    enable data race detection.
    Supported only on linux/amd64, freebsd/amd64, darwin/amd64, darwin/arm64, windows/amd64,
    linux/ppc64le and linux/arm64 (only for 48-bit VMA).
-msan
    enable interoperation with memory sanitizer.
    Supported only on linux/amd64, linux/arm64, freebsd/amd64
    and only with Clang/LLVM as the host C compiler.
    PIE build mode will be used on all platforms except linux/amd64.
-asan
    enable interoperation with address sanitizer.
    Supported only on linux/arm64, linux/amd64.
    Supported only on linux/amd64 or linux/arm64 and only with GCC 7 and higher
    or Clang/LLVM 9 and higher.
-cover
    enable code coverage instrumentation.
-covermode set,count,atomic
    set the mode for coverage analysis.
    The default is "set" unless -race is enabled,
    in which case it is "atomic".
    The values:
    set: bool: does this statement run?
    count: int: how many times does this statement run?
    atomic: int: count, but correct in multithreaded tests;
        significantly more expensive.
    Sets -cover.
-coverpkg pattern1,pattern2,pattern3
    For a build that targets package 'main' (e.g. building a Go
    executable), apply coverage analysis to each package matching
    the patterns. The default is to apply coverage analysis to
    packages in the main Go module. See 'go help packages' for a
    description of package patterns.  Sets -cover.
-v
    print the names of packages as they are compiled.
-work
    print the name of the temporary work directory and
    do not delete it when exiting.
-x
    print the commands.
-asmflags '[pattern=]arg list'
    arguments to pass on each go tool asm invocation.
-buildmode mode
    build mode to use. See 'go help buildmode' for more.
-buildvcs
    Whether to stamp binaries with version control information
    ("true", "false", or "auto"). By default ("auto"), version control
    information is stamped into a binary if the main package, the main module
    containing it, and the current directory are all in the same repository.
    Use -buildvcs=false to always omit version control information, or
    -buildvcs=true to error out if version control information is available but
    cannot be included due to a missing tool or ambiguous directory structure.
-compiler name
    name of compiler to use, as in runtime.Compiler (gccgo or gc).
-gccgoflags '[pattern=]arg list'
    arguments to pass on each gccgo compiler/linker invocation.
-gcflags '[pattern=]arg list'
    arguments to pass on each go tool compile invocation.
    “pattern” 指的是包模式 通过go help packages可以进行查看
-installsuffix suffix
    a suffix to use in the name of the package installation directory,
    in order to keep output separate from default builds.
    If using the -race flag, the install suffix is automatically set to race
    or, if set explicitly, has _race appended to it. Likewise for the -msan
    and -asan flags. Using a -buildmode option that requires non-default compile
    flags has a similar effect.
-ldflags '[pattern=]arg list'
    arguments to pass on each go tool link invocation.
-linkshared
    build code that will be linked against shared libraries previously
    created with -buildmode=shared.
-mod mode
    module download mode to use: readonly, vendor, or mod.
    By default, if a vendor directory is present and the go version in go.mod
    is 1.14 or higher, the go command acts as if -mod=vendor were set.
    Otherwise, the go command acts as if -mod=readonly were set.
    See https://golang.org/ref/mod#build-commands for details.
-modcacherw
    leave newly-created directories in the module cache read-write
    instead of making them read-only.
-modfile file
    in module aware mode, read (and possibly write) an alternate go.mod
    file instead of the one in the module root directory. A file named
    "go.mod" must still be present in order to determine the module root
    directory, but it is not accessed. When -modfile is specified, an
    alternate go.sum file is also used: its path is derived from the
    -modfile flag by trimming the ".mod" extension and appending ".sum".
-overlay file
    read a JSON config file that provides an overlay for build operations.
    The file is a JSON struct with a single field, named 'Replace', that
    maps each disk file path (a string) to its backing file path, so that
    a build will run as if the disk file path exists with the contents
    given by the backing file paths, or as if the disk file path does not
    exist if its backing file path is empty. Support for the -overlay flag
    has some limitations: importantly, cgo files included from outside the
    include path must be in the same directory as the Go package they are
    included from, and overlays will not appear when binaries and tests are
    run through go run and go test respectively.
-pgo file
    specify the file path of a profile for profile-guided optimization (PGO).
    When the special name "auto" is specified, for each main package in the
    build, the go command selects a file named "default.pgo" in the package's
    directory if that file exists, and applies it to the (transitive)
    dependencies of the main package (other packages are not affected).
    Special name "off" turns off PGO. The default is "auto".
-pkgdir dir
    install and load all packages from dir instead of the usual locations.
    For example, when building with a non-standard configuration,
    use -pkgdir to keep generated packages in a separate location.
-tags tag,list
    a comma-separated list of additional build tags to consider satisfied
    during the build. For more information about build tags, see
    'go help buildconstraint'. (Earlier versions of Go used a
    space-separated list, and that form is deprecated but still recognized.)
-trimpath
    remove all file system paths from the resulting executable.
    Instead of absolute file system paths, the recorded file names
    will begin either a module path@version (when using modules),
    or a plain import path (when using the standard library, or GOPATH).
-toolexec 'cmd args'
    a program to use to invoke toolchain programs like vet and asm.
    For example, instead of running asm, the go command will run
    'cmd args /path/to/asm <arguments for asm>'.
    The TOOLEXEC_IMPORTPATH environment variable will be set,
    matching 'go list -f {{.ImportPath}}' for the package being built.
```

-asmflags、-gccgoflags、-gcflags 和 -ldflags 标志接受以空格分隔的参数列表，以在构建期间传递给基础工具。要在列表中的元素中嵌入空格，请用单引号或双引号将其引起来。参数列表前面可以有一个包模式和一个等号，这将该参数列表的使用限制为构建与该模式匹配的包（有关包模式package pattern的描述，请参阅“go help packages”）。如果没有模式<mark>pattern</mark>，参数列表仅适用于命令行上命名的包。这些标志可以以不同的模式重复，以便为不同的包集指定不同的参数。如果一个包与多个标志中给出的模式匹配，则命令行上的最新匹配获胜。例如，“go build -gcflags=-S fmt”仅打印软件包 fmt 的反汇编，而“go build -gcflags=all=-S fmt”则打印 fmt 及其所有依赖项的反汇编。

有关指定包的更多信息，请参阅 "go help packages"。有关包和二进制文件安装位置的更多信息，请运行 "go help gopath"。有关 Go 和 C/C++ 之间调用的更多信息，请运行 "go help c"。

注意：Build 遵循某些约定，例如 "go help gopath "所描述的约定。但并非所有项目都能遵循这些约定。有自己的约定或使用独立软件编译系统的安装程序可能会选择使用低级调用，如 "go tool compile "和 "go tool link"，以避免编译工具的一些开销和设计决策。

## pattern

```shell
go help packages
```

[packages]导入路径的一个列表，一个导入路径

导入路径如果是 根路径 或以 . 或 ..元素开头的导入路径被解释为文件系统路径，并
表示该目录中的软件包。

否则，导入路径 P 表示在 GOPATH 环境变量中列出的某些 DIR 的目录 DIR/src/P 中找到的包（有关更多详细信息，请参阅：'go help gopath'）。

如果没有导入路径指定，Action只申请当前路径的下包

这里有四个保留的包名：

- main  “main”表示独立可执行文件中的顶级包。

- all  “all”扩展到所有 GOPATH 中找到的所有包树。例如，“go list all”列出本地上的所有包系统。使用模块时，“all”扩展到所有包主模块及其依赖项，包括依赖项任何这些测试都需要。

- std “std”与 all 类似，但仅扩展为标准中的go package。

- cmd “cmd” 扩展为 Go 存储库的命令及其内部库。

- ...  "..."是一个通配符，能匹配任何字符串，包含空字符串以及包含斜杠的字符串，这样的模式扩展到 GOPATH 树中找到的名称与模式匹配的所有包目录。为了使常见模式更加方便，有两种特殊情况。首先，模式末尾的 /... 可以匹配空字符串，因此 net/... 可以匹配 net 及其子目录中的包，例如 net/http。其次，<mark>任何包含通配符的斜杠分隔模式元素永远不会参与vendored包路径中“vendor”元素的匹配</mark>，因此 ./... 不匹配 ./vendor 或 ./ 子目录中的包mycode/vendor，但 ./vendor/... 和 ./mycode/vendor/... 可以。但请注意，本身包含代码的名为供应商的目录不是vendored包：cmd/vendor 将是名为vendor的命令，并且模式 cmd/... 与它匹配。有关供应商的更多信息，请参阅 golang.org/s/go15vendor。

## Build constraints

构建约束，也称为构建标签，是文件应包含在包中的条件。构建约束由开始的行注释给出

```go
//go:build
```

约束可以出现在任何类型的源文件中（不仅仅是 Go），但它们必须出现在文件顶部附近，前面只能有空行和其他行注释。这些规则意味着在 Go 文件中，构建约束必须出现在 package 子句之前。

为了区分构建约束和包文档，构建约束后面应该跟一个空行。

构建约束注释被评估为包含由 ||、&& 和 ! 组合的构建标签的表达式。运算符和括号。运算符的含义与 Go 中相同。

例如，以下构建约束会在满足“linux”和“386”约束时，或者在满足“darwin”但不满足“cgo”时约束要构建的文件：

```go
//go:build (linux && 386) || (darwin && !cgo)
```

一个文件包含多个 //go:build 行是错误的。

在特定的构建过程中，要满足以下构建标记的要求：

- 目标操作系统，由runtime.GOOS拼写，使用GOOS环境变量设置。

- 目标架构，由runtime.GOARCH拼写，使用GOARCH环境变量设置。

- 任何架构功能，采用 GOARCH.feature 形式（例如“amd64.v2”），如下详述。

- “unix”，如果 GOOS 是 Unix 或类 Unix 系统。

- 正在使用的编译器，“gc”或“gccgo”

- “cgo”，如果支持 cgo 命令（请参阅“go 帮助环境”中的 CGO_ENABLED）。

- 每个 Go 主要版本的术语，直到当前版本：从 Go 版本 1.1 开始为“go1.1”，从 Go 1.12 开始为“go1.12”，依此类推。

- -tags 标志给出的任何附加标签（请参阅“go help build”）。

测试版或次要版本没有单独的构建标签。

如果文件名在去掉扩展名和可能的 _test 后缀后，符合以下任何一种模式：

```go
*_GOOS
*_GOARCH
*_GOOS_GOARCH
```

（示例：source_windows_amd64.go）其中 GOOS 和 GOARCH 分别表示任何已知的操作系统和体系结构值，则该文件被视为具有需要这些术语的隐式构建约束（除了文件中的任何显式约束之外）。

使用 GOOS=android 除了 android 标签和文件之外，还与 GOOS=linux 匹配构建标签和文件。

除了 illumos 标签和文件之外，使用 GOOS=illumos 还可以与 GOOS=solaris 一样匹配构建标签和文件。

除了 ios 标签和文件之外，使用 GOOS=ios 还可以与 GOOS=darwin 一样匹配构建标签和文件。

定义的架构功能构建标签是：

- 对于 GOARCH=386、GO386=387 和 GO386=sse2 分别设置 386.387 和 386.sse2 构建标记。

- 对于 GOARCH=amd64、GOAMD64=v1、v2 和 v3 对应于 amd64.v1、amd64.v2 和 amd64.v3 功能构建标签。

- 对于 GOARCH=arm，GOARM=5、6 和 7 对应于arm.5、arm.6 和arm.7 功能构建标签。

- 对于 GOARCH=mips 或 mipsle，GOMIPS=hardfloat 和 softfloat 对应于 mips.hardfloat 和 mips.softfloat（或 mipsle.hardfloat 和 mipsle.softfloat）功能构建标记。

- 对于 GOARCH=mips64 或 mips64le，GOMIPS64=hardfloat 和 softfloat 对应于 mips64.hardfloat 和 mips64.softfloat（或 mips64le.hardfloat 和 mips64le.softfloat）功能构建标记。

- 对于 GOARCH=ppc64 或 ppc64le，GOPPC64=power8、power9 和 power10 对应于 ppc64.power8、ppc64.power9 和 ppc64.power10（或 ppc64le.power8、ppc64le.power9 和 ppc64le.power10）功能构建标签。

- 对于 GOARCH=wasm，GOWASM=satconv 和signext 对应于 wasm.satconv 和 wasm.signext 功能构建标签。

对于 GOARCH=amd64、arm、ppc64 和 ppc64le，特定功能级别也会为所有先前级别设置功能构建标签。例如，GOAMD64=v2 设置 amd64.v1 和 amd64.v2 功能标志。这可以确保在引入 GOAMD64=v4 时，使用 v2 功能的代码可以继续编译。处理缺少特定功能级别的代码应使用否定：

```shell
//go:build !amd64.v2
```

要防止某个文件被考虑用于任何构建：

```go
//go:build ignore
```

（任何其他不满意的词也可以，但“忽略”是约定俗成的。）

仅在使用 cgo 时且仅在 Linux 和 OS X 上构建文件：

```go
//go:build cgo && (linux || darwin)
```

这样的文件通常与实现其他系统默认功能的另一个文件配对，在这种情况下，该文件将带有约束：

```go
//go:build !(cgo && (linux || darwin))
```

将文件命名为 dns_windows.go 将导致仅在为 Windows 构建包时才包含该文件；同样，仅在构建 32 位 x86 包时才会包含 math_386.s。

Go 版本 1.16 及更早版本使用不同的语法来构建约束，并带有“// +build”前缀。当遇到旧语法时，gofmt 命令将添加等效的 //go:build 约束。
