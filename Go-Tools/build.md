#### [Compile packages and dependencies](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)

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

-asmflags、-gccgoflags、-gcflags 和 -ldflags 标志接受以空格分隔的参数列表，以在构建期间传递给基础工具。要在列表中的元素中嵌入空格，请用单引号或双引号将其引起来。参数列表前面可以有一个包模式和一个等号，这将该参数列表的使用限制为构建与该模式匹配的包（有关包模式的描述，请参阅“go help packages”）。如果没有模式，参数列表仅适用于命令行上命名的包。这些标志可以以不同的模式重复，以便为不同的包集指定不同的参数。如果一个包与多个标志中给出的模式匹配，则命令行上的最新匹配获胜。例如，“go build -gcflags=-S fmt”仅打印软件包 fmt 的反汇编，而“go build -gcflags=all=-S fmt”则打印 fmt 及其所有依赖项的反汇编。

有关指定包的更多信息，请参阅 "go help packages"。有关包和二进制文件安装位置的更多信息，请运行 "go help gopath"。有关 Go 和 C/C++ 之间调用的更多信息，请运行 "go help c"。

注意：Build 遵循某些约定，例如 "go help gopath "所描述的约定。但并非所有项目都能遵循这些约定。有自己的约定或使用独立软件编译系统的安装程序可能会选择使用低级调用，如 "go tool compile "和 "go tool link"，以避免编译工具的一些开销和设计决策。
