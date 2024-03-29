# Go Modules Reference Go模块手册

## 介绍

模块是 Go 管理依赖关系的方式。

本文档是 Go 模块系统的详细参考手册。有关创建 Go 项目的介绍，请参阅如何编写 Go 代码[How to Write Go Code](https://go.dev/doc/code.html).。有关使用模块、将项目迁移到模块以及其他主题的信息，请参阅从[Using Go Modules](https://blog.golang.org/using-go-modules)。

## Module 模块、Package包和 Version 版本

module 模块是一起发布、版本控制和分发的包的集合。模块可以直接从版本控制存储库或模块代理服务器下载。

模块由[模块路径](https://go.dev/ref/mod#glos-module-path)（在 go.mod 文件中声明）以及有关模块依赖项的信息来标识。模块根目录是包含 go.mod 文件的目录。主模块是包含调用 go 命令的目录的模块。

模块中的每个包都是同一目录中编译在一起的源文件的集合。包路径是与包含包的子目录连接的模块路径（相对于模块根）。例如，模块“golang.org/x/net”包含目录“html”中的包。该包的路径是“golang.org/x/net/html”。

### Module paths 模块路径

模块路径module path是模块的规范名称，在模块的 go.mod 文件中使用 module 指令声明。<mark>模块的路径是模块内包路径的前缀</mark>。

模块路径应该描述<mark>模块的功能</mark>以及在<mark>哪里可以找到它</mark>。通常，模块路径由存储库根路径、存储库中的目录（通常为空）和主版本后缀（仅适用于主要版本 2 或更高版本）组成。

- *存储库根路径*是模块路径的一部分，对应于开发模块的版本控制存储库的根目录。大多数模块都是在其存储库的根目录中定义的，因此这通常是整个路径。例如，golang.org/x/net 是同名模块的存储库根路径。有关 go 命令如何使用从模块路径派生的 HTTP 请求来查找存储库的信息，请参阅查找模块路径的存储库[Finding a repository for a module path](https://go.dev/ref/mod#vcs-find)。

- 如果模块未在存储库的根目录中定义，则*模块子目录*是命名该目录的模块路径的一部分，不包括主版本后缀。这也用作语义版本标签的前缀。例如，模块 golang.org/x/tools/gopls 位于根路径为 golang.org/x/tools 的存储库的 gopls 子目录中，因此它具有模块子目录 gopls。请参阅将版本映射到存储库中的提交([Mapping versions to commits](https://go.dev/ref/mod#vcs-version))和模块目录([Module directories within a repository](https://go.dev/ref/mod#vcs-dir))。

- 如果模块以主版本 2 或更高版本发布，则模块路径必须以<mark>主版本后缀（如 /v2）</mark>结尾。这可能是也可能不是子目录名称的一部分。例如，路径为 golang.org/x/repo/sub/v2 的模块可能位于存储库 golang.org/x/repo 的 /sub 或 /sub/v2 子目录中。

如果某个模块可能被其他模块依赖，则必须遵循这些规则，以便 go 命令可以找到并下载该模块。模块路径中允许的字符也有一些词法限制[lexical restrictions](https://go.dev/ref/mod#go-mod-file-ident)。

### Versions 版本

*版本version*标识*模块的不可变快照*，可以是发行版 [release](https://go.dev/ref/mod#glos-release-version)或预发行版[pre-release](https://go.dev/ref/mod#glos-pre-release-version)。每个版本都以字母 v 开头，后跟语义版本。有关如何格式化、解释和比较版本的详细信息，请参阅[Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)。
概括地说，语义版本由三个非负整数（从左到右依次为主要版本、次要版本和补丁版本）组成，中间用点隔开。补丁版本后面可以是一个可选的预发布字符串，以连字符开头。在预发布字符串或补丁版本后，还可以是以加号开头的构建元数据字符串。例如，v0.0.0、v1.12.134、v8.0.5-pre 和 v2.0.9+meta 都是有效的版本。

版本的每个部分都表明该版本是否稳定，是否与以前的版本兼容。

- 在对模块的公共接口或记录的功能进行向后不兼容的更改后（例如，删除包后），必须递增主要版本[major version](https://go.dev/ref/mod#glos-major-version)，并且必须将次要版本和补丁版本设置为零。

- 在向后兼容的更改之后（例如，在添加新功能之后），必须递增次要版本[minor version](https://go.dev/ref/mod#glos-minor-version)并将补丁版本设置为零。

- 在不影响模块公共接口的更改（例如错误修复或优化）之后，必须增加补丁版本[patch version](https://go.dev/ref/mod#glos-patch-version)。

- 预发布后缀[pre-release](https://go.dev/ref/mod#glos-pre-release-version)表示版本是预发布的。预发布版本排序在相应的发布版本之前。例如，v1.2.3-pre 出现在 v1.2.3 之前。

- 为了比较版本，构建元数据后缀将被忽略。版本控制存储库中会忽略带有构建元数据的标签，但构建元数据会保留在 go.mod 文件中指定的版本中。后缀+不兼容表示迁移到模块版本主要版本2或更高版本之前发布的版本[Go Modules Reference - The Go Programming Language](https://go.dev/ref/mod#non-module-compat)

如果版本的主要版本为 0 或具有预发布后缀，则该版本被视为不稳定。不稳定版本不受兼容性要求的约束。例如，v0.2.0 可能与 v0.1.0 不兼容，v1.5.0-beta 可能与 v1.5.0 不兼容。

Go 可以使用不遵循这些约定的标签、分支或修订来访问版本控制系统中的模块。但是，在主模块中，go 命令会自动将不遵循此标准的修订名称转换为规范版本。作为此过程的一部分，go 命令还将删除构建元数据后缀（+incompatible）。这可能会产生伪版本[pseudo-version](https://go.dev/ref/mod#glos-pseudo-version)、对修订标识符（例如Git commit hash）和来自版本控制系统的时间戳进行编码的预发布版本。例如，命令 go get golang.org/x/net@daa7c041 会将提交哈希 daa7c041 转换为伪版本 v0.0.0-20191109021931-daa7c04131f5。主模块之外需要规范版本，如果 go.mod 文件中出现像 master 这样的非规范版本，go 命令将报告错误。

### Pseudo-versions 伪版本

伪版本*pseudo-version*是一种特殊格式的预发布版本[pre-release](https://go.dev/ref/mod#glos-pre-release-version) [version](https://go.dev/ref/mod#glos-version)，它对版本控制存储库中特定修订的信息进行编码。例如，v0.0.0-20191109021931-daa7c04131f5是一个伪版本。

伪版本可以指没有可用的语义版本标签的修订。它们可用于在创建版本标签之前测试提交，例如在开发分支上。

每个伪版本由三部分组成：

- 基本版本前缀（vX.0.0 或 vX.Y.Z-0），源自修订版之前的语义版本标签，或者源自 vX.0.0（如果没有此类标签）。

- 时间戳 (yyyymmddhhmmss)，这是创建修订版的 UTC 时间。在 Git 中，这是提交时间，而不是作者时间。

- 修订标识符 (abcdefabcdef)，它是提交哈希的 12 个字符前缀，或者在 Subversion 中是一个零填充的修订号。

每个伪版本可能采用三种形式之一，具体取决于基本版本。这些形式确保伪版本比其基本版本更高，但比下一个标记版本低。

- 当没有已知的基本版本时，使用 vX.0.0-yyyymmddhhmmss-abcdefabcdef。与所有版本一样，主版本 X 必须与模块的主版本后缀匹配。

- 当基础版本是像 vX.Y.Z-pre 这样的预发布版本时，使用 vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef。

- 当基础版本是像 vX.Y.Z 这样的发布版本时，使用 vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef。例如，如果基本版本是 v1.2.3，则伪版本可能是 v1.2.4-0.20191109021931-daa7c04131f5。

多个伪版本可能通过使用不同的基本版本来引用相同的提交。当写入伪版本后标记较低版本时，这种情况自然会发生。

这些形式为伪版本提供了两个有用的属性：

- 具有已知基本版本的伪版本排序高于这些版本，但低于后续版本的其他预发布版本。

- 具有相同基本版本前缀的伪版本按时间顺序排序。

go 命令执行多项检查，以确保模块作者可以控制如何将伪版本与其他版本进行比较，并且伪版本引用实际上是模块提交历史记录一部分的修订。

- 如果指定了基本版本，则必须有一个相应的语义版本标记，该标记是伪版本所描述的修订版本的祖先。这可以防止开发人员使用比所有标记版本（例如 v1.999.999-99999999999999-daa7c04131f5）更高的伪版本来绕过[minimal version selection](https://go.dev/ref/mod#glos-minimal-version-selection)。

- 时间戳必须与修订版本的时间戳相匹配。这可以防止攻击者用无限数量的相同伪版本淹没模块代理。这也可以防止模块使用者更改版本的相对顺序。

- 修订版本必须是模块存储库的分支或标签之一的祖先。这可以防止攻击者引用未经批准的更改或拉取请求。

伪版本永远不需要手动输入。许多命令接受提交哈希或分支名称，并自动将其转换为伪版本（或标记版本，如果可用）。例如：

```shell
go get example.com/mod@master
go list -m -json example.com/mod@abcd1234
```

### Major version suffixes 主要版本后缀

从主要版本 2 开始，模块路径必须具有与主要版本匹配的主要版本后缀，例如 /v2。例如，如果模块在 v1.0.0 处具有路径 example.com/mod，则在版本 v2.0.0 处它必须具有路径 example.com/mod/v2。

主要版本后缀实现导入兼容性规则[import compatibility rule](https://research.swtch.com/vgo-import)：

> 如果旧包和新包具有相同的导入路径，则新包必须向后兼容旧包。

根据定义，模块的新主要版本中的包不向后兼容先前主要版本中的相应包。因此，从 v2 开始，包需要新的导入路径。这是通过向模块路径添加主要版本后缀来完成的。由于模块路径是模块内每个包的导入路径的前缀，因此将主要版本后缀添加到模块路径可为每个不兼容的版本提供不同的导入路径。

主要版本 v0 或 v1 不允许使用主要版本后缀。 v0和v1之间不需要更改模块路径，因为v0版本不稳定并且没有兼容性保证。此外，对于大多数模块，v1 向后兼容最新的 v0 版本； v1 版本充当对兼容性的承诺，而不是表示与 v0 相比不兼容的更改。

作为一种特殊情况，以 gopkg.in/ 开头的模块路径必须始终具有主版本后缀，即使是 v0 和 v1。后缀必须以点而不是斜线开头（例如 gopkg.in/yaml.v2）。

主要版本后缀允许模块的多个主要版本共存于同一构建中。由于钻石依赖性问题([diamond dependency problem](https://research.swtch.com/vgo-import#dependency_story).)，这可能是必要的。通常，如果通过传递依赖关系需要两个不同版本的模块，则将使用较高版本。但是，如果两个版本不兼容，则两个版本都无法满足所有客户的要求。由于不兼容的版本必须具有不同的主版本号，因此由于主版本后缀，它们也必须具有不同的模块路径。这解决了冲突：具有不同后缀的模块被视为单独的模块，并且它们的包（甚至相对于其模块根位于同一子目录中的包）是不同的。

许多 Go 项目在迁移到模块之前（可能甚至在引入模块之前）就发布了 v2 或更高版本，而不使用主要版本后缀。这些版本用+不兼容的构建标签进行注释（例如，v2.0.0+不兼容）。有关更多信息，请参阅与非模块存储库的兼容性([Compatibility with non-module repositories](https://go.dev/ref/mod#non-module-compat))。

### Resolving a package to a module 将包解析为模块

当 go 命令使用包路径加载包时，它需要确定哪个模块提供该包。

go 命令首先在构建列表( [build list](https://go.dev/ref/mod#glos-build-list))中搜索路径为包路径前缀的模块。例如，如果导入了包 example.com/a/b，并且模块 example.com/a 在构建列表中，则 go 命令将检查目录 b 中的 example.com/a 是否包含该包。目录中必须至少存在一个扩展名为 .go 的文件，才能将其视为包。构建约束([Build constraints](https://go.dev/pkg/go/build/#hdr-Build_Constraints))不适用于此目的。如果构建列表中恰好有一个模块提供了该包，则使用该模块。如果没有模块提供该包，或者有两个或更多模块提供该包，则 go 命令会报告错误。 -mod=mod 标志指示 go 命令尝试查找提供缺失软件包的新模块并更新 go.mod 和 go.sum。 go get 和 go mod tidy 命令会自动执行此操作。

当 go 命令查找包路径的新模块时，它会检查 GOPROXY 环境变量，该变量是代理 URL 或关键字 direct 或 off 的逗号分隔列表。代理 URL 指示 go 命令应使用 GOPROXY 协议联系模块代理。 direct 指示 go 命令应与版本控制系统通信。 off 表示不应尝试进行任何通信。 GOPRIVATE 和 GONOPROXY 环境变量也可用于控制此行为。

对于 GOPROXY 列表中的每个条目，go 命令请求可能提供包的每个模块路径的最新版本（即包路径的每个前缀）。对于每个成功请求的模块路径，go 命令都会下载最新版本的模块并检查该模块是否包含所请求的包。如果一个或多个模块包含所请求的包，则使用具有最长路径的模块。如果找到一个或多个模块，但没有一个包含所请求的包，则会报告错误。如果没有找到模块，go 命令会尝试 GOPROXY 列表中的下一个条目。如果没有留下任何条目，则会报告错误。

例如，假设 go 命令正在查找提供包 golang.org/x/net/html 的模块，并且 GOPROXY 设置为 https://corp.example.com,https://proxy.golang.org 。 go 命令可能会发出以下请求：

- 前往 https://corp.example.com/（并行）：
  
  - 请求最新版本的 golang.org/x/net/html
  
  - 请求最新版本的 golang.org/x/net
  
  - 请求最新版本的 golang.org/x
  
  - 请求最新版本的 golang.org

- 对于 https://proxy.golang.org/，如果对 https://corp.example.com/ 的所有请求都失败并返回 404 或 410：
  
  - 请求最新版本的 golang.org/x/net/html
  
  - 请求最新版本的 golang.org/x/net
  
  - 请求最新版本的 golang.org/x
  
  - 请求最新版本的 golang.org

找到合适的模块后，go 命令会将新的需求([requirement](https://go.dev/ref/mod#go-mod-file-require))以及新模块的路径和版本添加到主模块的 go.mod 文件中。这保证了以后加载相同的包时，相同的版本会使用相同的模块。如果解析的包不是由主模块中的包导入的，则新需求将有 ***// indirect***。

## go.mod 文件

一个模块由utf-8编码的文本文件定义，其根目录中名为go.mod。 go.mod文件是面向行的。每行都有一个指令，由关键字组成，然后是参数。例如：

```go-mod
module example.com/my/thing

go 1.12

require example.com/other/thing v1.0.2
require example.com/new/thing/v2 v2.3.4
exclude example.com/old/thing v1.2.3
replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
retract [v1.9.0, v1.9.5]
```

前导关键字可以从相邻行中分解出来以创建一个块，就像在 Go 导入中一样。

```go-mod
require (
    example.com/new/thing/v2 v2.3.4
    example.com/old/thing v1.2.3
)
```

go.mod 文件被设计为人类可读且机器可写。 go 命令提供了几个更改 go.mod 文件的子命令。例如，go get 可以升级或降级特定的依赖项。加载模块图的命令将在需要时自动更新 go.mod。 go mod edit 可以执行低级编辑。 Go 程序可以使用 golang.org/x/mod/modfile 包以编程方式进行相同的更改。

主模块([main module](https://go.dev/ref/mod#glos-main-module))以及使用本地文件路径指定的任何替换模块([replacement module](https://go.dev/ref/mod#go-mod-file-replace))都需要 go.mod 文件。但是，缺少显式 go.mod 文件的模块可能仍需要作为依赖项( [required](https://go.dev/ref/mod#go-mod-file-require))，或用作使用模块路径和版本指定的替换；请参阅与非模块存储库的兼容性([Compatibility with non-module repositories](https://go.dev/ref/mod#non-module-compat).)。

### Lexical elements 元素

当解析 go.mod 文件时，其内容被分解为一系列标记。有多种标记：空格、注释、标点符号、关键字、标识符和字符串。

空白由空格 (U+0020)、制表符 (U+0009)、回车符 (U+000D) 和换行符 (U+000A) 组成。除了换行符之外的空白字符没有任何作用，除非可以分隔本来要组合的标记。换行符是重要的标记。

注释以 // 开始，一直到行尾。 /* */ 不允许注释。

标点符号包括 (、) 和 =>。

关键字区分 go.mod 文件中不同类型的指令。允许的关键字包括 module、go、require、replace、exclusion 和撤回。

标识符是非空白字符的序列，例如模块路径或语义版本。

字符串是带引号的字符序列。有两种字符串：以引号开头和结尾的解释字符串 (", U+0022) 和以重音符号开头和结尾的原始字符串 (`, U+0060)。解释字符串可能包含由反斜杠 ( \, U+005C) 后跟另一个字符。转义引号 (\") 不会终止解释的字符串。解释字符串的不带引号的值是引号之间的字符序列，每个转义序列都替换为反斜杠后面的字符（例如，\" 替换为 "，\n 替换为 n）。相反，原始字符串的未加引号的值只是重音符号之间的字符序列；反斜杠在原始字符串中没有特殊含义。

标识符和字符串在 go.mod 语法中是可以互换的。

### Module paths and versions 模块路径和版本

go.mod 文件中的大多数标识符和字符串都是模块路径或版本。

模块路径必须满足以下要求：

- 该路径必须由一个或多个以斜杠（/、U+002F）分隔的路径元素组成。它不能以斜线开头或结尾。

- 每个路径元素都是一个由 ASCII 字母、ASCII 数字和有限的 ASCII 标点符号（-、.、_ 和 ~）组成的非空字符串。

- 路径元素不能以点 (., U+002E) 开头或结尾。

- 第一个点之前的元素前缀不得是 Windows 上的保留文件名，无论大小写（CON、com1、NuL 等）。

- 第一个点之前的元素前缀不得以波浪号后跟一个或多个数字结尾（如 EXAMPL~1.COM）。

如果模块路径出现在 require 指令中并且未被替换，或者模块路径出现在 replace 指令的右侧，则 go 命令可能需要使用该路径下载模块，并且必须满足一些附加要求。

- 按照惯例，前导路径元素（直到第一个斜杠，如果有）必须仅包含小写 ASCII 字母、ASCII 数字、点 (.、U+002E) 和破折号 (-、U+002D) ）；它必须至少包含一个点并且不能以破折号开头。

- 对于 /vN 形式的最终路径元素，其中 N 看起来是数字（ASCII 数字和点），N 不能以前导零开头，不能是 /v1，并且不能包含任何点。
  
  - 对于以 gopkg.in/ 开头的路径，此要求被路径遵循 gopkg.in 服务约定的要求所取代。

go.mod 文件中的版本可能是规范的([canonical](https://go.dev/ref/mod#glos-canonical-version))或非规范的。

规范版本以字母 v 开头，后跟遵循语义版本 2.0.0 规范的语义版本。有关详细信息，请参阅版本。

大多数其他标识符和字符串可以用作非规范版本，但有一些限制以避免文件系统、存储库和模块代理出现问题。非规范版本仅允许在主模块的 go.mod 文件中使用。当 go 命令自动更新 go.mod 文件时，它会尝试将每个非规范版本替换为等效的规范版本。

在模块路径与版本关联的地方（如 require、replace 和 except 指令），最终路径元素必须与版本一致。请参阅主要版本后缀([Major version suffixes](https://go.dev/ref/mod#major-version-suffixes))。

### Grammar 语法

go.mod 语法在下面使用扩展巴科斯-诺尔范式 (EBNF) 指定。有关 EBNF 语法的详细信息，请参阅 Go 语言规范中的符号部分。

```ebnf
GoMod = { Directive } .
Directive = ModuleDirective |
            GoDirective |
            RequireDirective |
            ExcludeDirective |
            ReplaceDirective |
            RetractDirective .
```

换行符、标识符和字符串分别用 newline、ident 和 string 表示。

模块路径和版本用 ModulePath 和 Version 表示。

```ebnf
ModulePath = ident | string . /* see restrictions above */
Version = ident | string .    /* see restrictions above */
```

### module directive 模块指令

模块指令定义主模块的路径([path](https://go.dev/ref/mod#glos-module-path))。 go.mod 文件必须恰好包含一个模块指令。

```ebnf
ModuleDirective = "module" ( ModulePath | "(" newline ModulePath newline ")" ) newline .
```

例如：

```go-mod
module golang.org/x/net
```

### Deprecation 弃用

可以在段落开头包含字符串 Deprecated:（区分大小写）的注释块中将模块标记为已弃用。弃用消息在冒号之后开始，一直到段落末尾。注释可能出现在模块指令之前，也可能出现在同一行之后。

例：

```go-mod
// Deprecated: use example.com/mod/v2 instead.
module example.com/mod
```

从 Go 1.17 开始， go list -m -u 检查构建列表中所有已弃用模块的信息。 go get 检查构建在命令行上命名的包所需的已弃用模块。

当 go 命令检索模块的弃用信息时，它会从与 @latest 版本查询匹配的版本加载 go.mod 文件，而不考虑撤回或排除。 go 命令从同一 go.mod 文件加载收回版本的列表。

要弃用模块，作者可以添加 // Deprecated: 注释并标记新版本。作者可能会在更高版本中更改或删除弃用消息。

弃用适用于模块的所有次要版本。为此，高于 v2 的主要版本被视为单独的模块，因为它们的主要版本后缀为它们提供了不同的模块路径。

弃用消息旨在通知用户该模块不再受支持，并提供迁移说明，例如迁移到最新的主要版本。个别次要版本和补丁版本不能被弃用；撤回可能更适合于此。

### go directive

go 指令指示模块是在假定给定 Go 版本的语义的情况下编写的。该版本必须是有效的 Go 版本，例如 1.9、1.14 或 1.21rc1。

go 指令设置使用该模块所需的最低 Go 版本。在 Go 1.21 之前，该指令仅是建议性的；现在这是一个强制性要求：Go 工具链拒绝使用声明较新 Go 版本的模块。

go 指令是选择要运行的 Go 工具链的输入。有关详细信息，请参阅“Go 工具链([Go toolchains](https://go.dev/doc/toolchain))”。

go 指令影响新语言功能的使用：

- 对于模块内的包，编译器拒绝使用 go 指令指定的版本之后引入的语言功能。例如，如果模块具有指令 go 1.12，则其包可能不会使用 Go 1.13 中引入的数字文字，例如 1_000_000。

- 如果较旧的 Go 版本构建该模块的某个包并遇到编译错误，该错误会指出该模块是为较新的 Go 版本编写的。例如，假设模块的版本为 1.13，而包使用数字文字 1_000_000。如果该包是使用 Go 1.12 构建的，编译器会注意到该代码是为 Go 1.13 编写的。

go 指令还会影响 go 命令的行为：

- 在 go 1.14 或更高版本中，可以启用自动[vendoring](https://go.dev/ref/mod#vendoring)。如果文件vendor/modules.txt存在并且与go.mod一致，则无需显式使用-mod=vendor标志。

- 在 go 1.16 或更高版本中，所有包模式仅匹配由主模块中的包和测试传递导入的包。这是自引入模块以来 go mod vendor保留的同一组软件包。在较低版本中，all 还包括对主模块中的包导入的包的测试、对这些包的测试等等。

- 在 go 1.17 或更高版本中：
  
  - go.mod 文件包含每个模块的显式[`require` directive](https://go.dev/ref/mod#go-mod-file-require)，该指令提供由主模块中的包或测试传递导入的任何包。 （在 go 1.16 及更低版本中，仅当最小版本选择([minimal version selection](https://go.dev/ref/mod#minimal-version-selection))会选择不同版本时才包含间接依赖项[indirect dependency](https://go.dev/ref/mod#glos-direct-dependency)。）此额外信息支持模块图修剪[module graph pruning](https://go.dev/ref/mod#graph-pruning)和延迟模块加载[lazy module loading](https://go.dev/ref/mod#lazy-loading)。
  
  - 因为可能比以前的 go 版本有更多 // 间接依赖关系，所以间接依赖关系被记录在 go.mod 文件中的单独块中。
  
  - go mod vendor 会省略 go.mod 和 go.sum 文件中的 vendored 依赖项。(这使得在 vendor 的子目录中调用 go 命令时能识别正确的主模块）。
  
  - go mod供应商在vendor/modules.txt中记录每个依赖项的go.mod文件中的go版本。

- go 1.21 或更高版本
  
  - go 行声明了与此模块一起使用所需的最低 Go 版本。
  
  - go行必须大于或等于所有依赖的go行。
  
  - go 命令不再尝试保持与之前旧版本 Go 的兼容性。
  
  - go 命令更加小心地将 go.mod 文件的校验和保存在 go.sum 文件中。

一个 go.mod 文件最多可以包含一个 go 指令。如果不存在，大多数命令都会添加包含当前 Go 版本的 go 指令。

如果缺少 go 指令，则假定为 go 1.16。

```ebnf
GoDirective = "go" GoVersion newline .
GoVersion = string | ident .  /* valid release version; see above */
```

例如：

```go-mod
go 1.14
```

### toolchain directive 工具链指令

工具链指令声明了与模块一起使用的建议 Go 工具链。建议的 Go 工具链版本不能低于 go 指令中声明的所需 Go 版本。仅当模块是主模块并且默认工具链的版本低于建议工具链的版本时，工具链指令才有效。

为了重现性，go 命令在更新 go.mod 文件中的 go 版本时（通常在 go get 期间），会在工具链行中写入自己的工具链名称。

有关详细信息，请参阅“[Go toolchains](https://go.dev/doc/toolchain)”。

```ebnf
ToolchainDirective = "toolchain" ToolchainName newline .
ToolchainName = string | ident .  /* valid toolchain name; see “Go toolchains” */
```

例如

```go-mod
toolchain go1.21.0
```

### exclude directive

排除指令可防止 go 命令加载模块版本。

从 Go 1.16 开始，如果任何 go.mod 文件中的 require 指令引用的版本被主模块的 go.mod 文件中的 except 指令排除，则该要求将被忽略。这可能会导致像 go get 和 go mod tidy 这样的命令向 go.mod 添加对更高版本的新要求，并在适当的情况下使用 // 间接注释。

在 Go 1.16 之前，如果 require 指令引用了排除版本，则 go 命令会列出该模块的可用版本（如 go list -m -versions 所示）并加载下一个更高的非排除版本。这可能会导致不确定的版本选择，因为下一个更高版本可能会随着时间而改变。发布版本和预发布版本都被考虑用于此目的，但伪版本则不然。如果没有更高的版本，go命令会报错。

except 指令仅适用于主模块的 go.mod 文件，并在其他模块中被忽略。有关详细信息，请参阅最小版本选择。

```ebnf
ExcludeDirective = "exclude" ( ExcludeSpec | "(" newline { ExcludeSpec } ")" newline ) .
ExcludeSpec = ModulePath Version newline .
```

例如：

```go-mod
exclude golang.org/x/net v1.2.3

exclude (
    golang.org/x/crypto v1.4.5
    golang.org/x/text v1.6.7
)
```

### require directive

require 指令声明给定模块依赖项的最低所需版本。对于每个所需的模块版本，go 命令会加载该版本的 go.mod 文件并合并该文件中的要求。加载所有需求后，go 命令会使用最小版本选择 (MVS [minimal version selection (MVS)](https://go.dev/ref/mod#minimal-version-selection)) 来解决它们以生成构建列表。

go 命令会针对某些需求自动添加 // 间接注释。 // 间接注释表示主模块中的任何包都不会直接导入所需模块中的包。

如果 go 指令指定 go 1.16 或更低版本，则当选定的模块版本高于主模块的其他依赖项已经暗示（传递）的版本时，go 命令会添加间接要求。发生这种情况的原因可能是显式升级 (go get -u ./...)、删除了先前施加要求的某些其他依赖项 (go mod tidy)，或者导入了自身没有相应要求的包的依赖项go.mod 文件（例如完全缺少 go.mod 文件的依赖项）。

在 go 1.17 及更高版本中，go 命令为每个模块添加了间接要求，该模块提供由主模块中的包或测试导入（甚至间接）或作为参数传递给 go get 的任何包。这些更全面的要求支持模块图修剪和延迟模块加载。

```ebnf
RequireDirective = "require" ( RequireSpec | "(" newline { RequireSpec } ")" newline ) .
RequireSpec = ModulePath Version newline .
```

例如：

```go-mod
require golang.org/x/net v1.2.3

require (
    golang.org/x/crypto v1.4.5 // indirect
    golang.org/x/text v1.6.7
)
```

### replace directive

replace指令用其他地方找到的内容替换模块的特定版本或模块的所有版本的内容。可以使用另一个模块路径和版本或特定于平台的文件路径来指定替换。

如果箭头左侧 (=>) 存在版本，则仅替换该模块的特定版本；其他版本都可以正常访问。如果省略左侧版本，则替换该模块的所有版本。

如果箭头右侧的路径是绝对或相对路径（以./或../开头），则将其解释为替换模块根目录的本地文件路径，其中必须包含go.mod文件。在这种情况下，必须省略替换版本。

如果右侧的路径不是本地路径，则它必须是有效的模块路径。在这种情况下，需要一个版本。相同的模块版本不得同时出现在构建列表中。

无论替换是使用本地路径还是模块路径指定，如果替换模块具有 go.mod 文件，则其模块指令必须与其替换的模块路径匹配。

替换指令仅适用于主模块的 go.mod 文件，并在其他模块中被忽略。有关详细信息，请参阅 [Minimal version selection](https://go.dev/ref/mod#minimal-version-selection)。

如果有多个主模块，则所有主模块的 go.mod 文件都适用。不允许主模块之间存在冲突的替换指令，并且必须在 go.work 文件的替换中删除或覆盖。

请注意，单独的替换指令不会将模块添加到模块图中。还需要在主模块的 go.mod 文件或依赖项的 go.mod 文件中引用替换模块版本的 require 指令。如果不需要左侧的模块版本，则替换指令无效。

```ebnf
ReplaceDirective = "replace" ( ReplaceSpec | "(" newline { ReplaceSpec } ")" newline ) .
ReplaceSpec = ModulePath [ Version ] "=>" FilePath newline
            | ModulePath [ Version ] "=>" ModulePath Version newline .
FilePath = /* platform-specific relative or absolute file path */
```

例如：

```go-mod
replace golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5

replace (
    golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5
    golang.org/x/net => example.com/fork/net v1.4.5
    golang.org/x/net v1.2.3 => ./fork/net
    golang.org/x/net => ./fork/net
)
```

### retract directive 撤回指令

撤回指令指示不应依赖 go.mod 定义的模块的版本或版本范围。当版本过早发布或版本发布后发现严重问题时，撤回指令非常有用。撤回的版本应在版本控制存储库和模块代理中保持可用，以确保依赖于它们的构建不会被破坏。 “撤回”这个词借用自学术文献：撤回的研究论文仍然可用，但它有问题，不应该成为未来工作的基础。

当模块版本被收回时，用户将不会使用 go get、go mod tidy 或其他命令自动升级到该模块版本。依赖于撤回版本的构建应该继续工作，但是当用户使用 go list -m -u 检查更新或使用 go get 更新相关模块时，用户将收到撤回通知。

要收回版本，模块作者应向 go.mod 添加收回指令，然后发布包含该指令的新版本。新版本必须高于其他发布或预发布版本；也就是说，@latest [version query](https://go.dev/ref/mod#version-queries)应在考虑撤回之前解析为新版本。 go 命令加载并应用 go list -m -retracted $modpath@latest 显示的版本的撤回（其中 $modpath 是模块路径）。

除非使用 -retracted 标志，否则收回的版本将从 go list -m -versions 打印的版本列表中隐藏。解析 @>=v1.2.3 或 @latest 等版本查询时，将排除已收回的版本。

包含撤回的版本可能会自行撤回。如果模块的最高版本或预发布版本自行收回，则在排除收回的版本后，@latest 查询将解析为较低版本。

例如，考虑模块 example.com/m 的作者意外发布版本 v1.0.0 的情况。为了防止用户升级到 v1.0.0，作者可以在 go.mod 中添加两个撤回指令，然后用撤回标记 v1.0.1。

```go-mod
retract (
    v1.0.0 // Published accidentally.
    v1.0.1 // Contains retractions only.
)
```

当用户运行 go get example.com/m@latest 时，go 命令会读取 v1.0.1（现在是最高版本）的撤回内容。 v1.0.0 和 v1.0.1 都已撤销，因此 go 命令将升级（或降级！）到下一个最高版本，可能是 v0.9.5。

撤消指令可以使用单个版本（如 v1.0.0）或具有上限和下限的封闭版本间隔编写，由 [ 和 ] 分隔（如 [v1.1.0, v1.2.0]）。单一版本相当于上限和下限相同的区间。与其他指令一样，多个撤回指令可以组合在一个块中，该块由 ( 在行尾和 ) 分隔在其自己的行上。

每个撤回指令都应有一条注释，解释撤回的理由，但这不是强制性的。 go 命令可能会在有关撤回版本的警告和 go 列表输出中显示基本原理注释。基本原理注释可以直接写在撤消指令上方（中间没有空行），也可以写在同一行之后。如果注释出现在块上方，则它适用于该块内没有自己注释的所有撤回指令。理由注释可以跨越多行。

```ebnf
RetractDirective = "retract" ( RetractSpec | "(" newline { RetractSpec } ")" newline ) .
RetractSpec = ( Version | "[" Version "," Version "]" ) newline .
```

例如：

- 收回 v1.0.0 和 v1.9.9 之间的所有版本：

```go-mod
retract v1.0.0
retract [v1.0.0, v1.9.9]
retract (
    v1.0.0
    [v1.0.0, v1.9.9]
)
```

- 在过早发布版本 v1.0.0 后返回到无版本控制：

```go-mod
retract [v0.0.0, v1.0.1] // assuming v1.0.1 contains this retraction.
```

- 清除包含所有伪版本和标记版本的模块：

```go-mod
retract [v0.0.0-0, v0.15.2]  // assuming v0.15.2 contains this retraction.
```

撤销指令是在<mark> Go 1.16 中添加</mark>的。 Go 1.15及更低版本如果在主模块的go.mod文件中写入retract指令会报告错误，并且会忽略依赖项的go.mod文件中的retract指令。

### Automatic updates

如果 go.mod 缺少信息或没有准确反映现实，大多数命令都会报告错误。 go get 和 go mod tidy 命令可用于解决大多数此类问题。此外，**-mod=mod 标志可以与大多数模块感知命​​令（go build、go test 等）一起使用**，以指示 go 命令自动修复 go.mod 和 go.sum 中的问题。

例如，考虑这个 go.mod 文件：

```go-mod
module example.com/M

go 1.16

require (
    example.com/A v1
    example.com/B v1.0.0
    example.com/C v1.0.0
    example.com/D v1.2.3
    example.com/E dev
)

exclude example.com/D v1.2.3
```

使用 -mod=mod 触发的更新将非规范版本标识符重写为规范 semver 形式，因此 example.com/A 的 v1 变为 v1.0.0，example.com/E 的 dev 变为 dev 上最新提交的伪版本分支，可能是 v0.0.0-20180523231146-b3f5c0f6e5f1。

此更新修改了要求以尊重排除，因此对排除的 example.com/D v1.2.3 的要求更新为使用 example.com/D 的下一个可用版本，可能是 v1.2.4 或 v1.3.0。

该更新删除了多余或误导性的要求。例如，如果 example.com/A v1.0.0 本身需要 example.com/B v1.2.0 和 example.com/C v1.0.0，那么 go.mod 对 example.com/B v1.0.0 的要求会产生误导（已被取代） example.com/A 需要 v1.2.0），而 example.com/C v1.0.0 的要求是多余的（暗示 example.com/A 需要相同版本），因此两者都将被删除。如果主模块包含直接从 example.com/B 或 example.com/C 导入包的包，则需求将保留，但更新为实际使用的版本。

最后，更新以规范格式重新格式化 go.mod，以便未来的机械更改将导致最小的差异。如果只需要更改格式，go 命令不会更新 go.mod。

因为模块图定义了 import 语句的含义，所以任何加载包的命令也使用 go.mod，因此可以更新它，包括 go build、go get、go install、go list、go test、go mod tidy。

在 Go 1.15 及更低版本中，默认启用 -mod=mod 标志，因此会自动执行更新。从 Go 1.16 开始，go 命令的行为就像设置了 -mod=readonly 一样：如果需要对 go.mod 进行任何更改，go 命令会报告错误并建议修复。

## Minimal version selection (MVS) 最小版本选择

Go 使用一种称为最小版本选择 (MVS) 的算法来选择构建包时要使用的一组模块版本。 MVS 在 Russ Cox 的《[Minimal Version Selection](https://research.swtch.com/vgo-mvs)》中进行了详细描述。

从概念上讲，MVS 在由 go.mod 文件指定的模块有向图上运行。图中的每个顶点代表一个模块版本。每条边代表依赖项的最低所需版本，使用 require 指令指定。该图可以通过主模块的 go.mod 文件中的排除和替换指令以及 go.work 文件中的替换指令进行修改。

MVS 生成构建列表([build list](https://go.dev/ref/mod#glos-build-list))作为输出，即用于构建的模块版本列表。

MVS 从主模块（图中没有版本的特殊顶点）开始并遍历图，跟踪每个模块所需的最高版本。在遍历结束时，最高要求的版本构成构建列表：它们是满足所有要求的最低版本。

可以使用命令 go list -m all 检查构建列表。与其他依赖管理系统不同，构建列表不保存在“锁”文件中。 MVS 是确定性的，当新版本的依赖项发布时，构建列表不会改变，因此 MVS 用于在每个模块感知命​​令的开始处计算它。

考虑下图中的示例。主模块需要1.2或更高版本的模块A和1.2或更高版本的模块B。 A 1.2 和 B 1.2 分别需要 C 1.3 和 C 1.4。 C 1.3 和 C 1.4 都需要 D 1.2。
![](images/buildlist.svg)

MVS 访问并加载每个以蓝色突出显示的模块版本的 go.mod 文件。在图遍历结束时，MVS 返回包含粗体版本的构建列表：A 1.2、B 1.2、C 1.4 和 D 1.2。请注意，可以使用更高版本的 B 和 D，但 MVS 不会选择它们，因为没有什么需要它们。

### Replacement 替换

模块的内容（包括其 go.mod 文件）可以使用主模块的 go.mod 文件或工作区的 go.work 文件中的替换指令进行替换。替换指令可以应用于模块的特定版本或模块的所有版本。

替换会更改模块图，因为替换模块可能具有与替换版本不同的依赖项。

考虑下面的示例，其中 C 1.4 已替换为 R。R 依赖于 D 1.3 而不是 D 1.2，因此 MVS 返回包含 A 1.2、B 1.2、C 1.4（替换为 R）和 D 1.3 的构建列表。

![](images/replace.svg)

### Exclusion 排除

还可以使用主模块的 go.mod 文件中的排除指令在特定版本中排除模块。

排除也会改变模块图。当某个版本被排除时，它将从模块图中删除，并且对其的要求将被重定向到下一个更高的版本。

考虑下面的例子。 C 1.3 已被排除。 MVS 将表现为 A 1.2 需要 C 1.4（下一个更高版本）而不是 C 1.3。

![](images/exclude.svg)

### Upgrades 升级版本

go get 命令可用于升级一组模块。要执行升级，go 命令会在运行 MVS 之前通过将已访问版本的边添加到升级版本来更改模块图。

考虑下面的例子。模块B可以从1.2升级到1.3，C可以从1.3升级到1.4，模块D可以从1.2升级到1.3。

![](images/upgrade.svg)

升级（和降级）可能会添加或删除间接依赖项。在这种情况下，升级后，E 1.1 和 F 1.1 会出现在构建列表中，因为 B 1.3 需要 E 1.1。

为了保留升级，go 命令更新了 go.mod 中的要求。它将把 B 的要求更改为 1.3 版本。它还将通过 // 间接注释添加对 C 1.4 和 D 1.3 的要求，因为否则不会选择这些版本。

### Downgrade 降级版本

go get 命令也可用于降级一组模块。要执行降级，go 命令通过删除降级版本之上的版本来更改模块图。它还删除依赖于已删除版本的其他模块的版本，因为它们可能与其依赖项的降级版本不兼容。如果主模块需要通过降级删除的模块版本，则要求更改为尚未删除的先前版本。如果没有可用的先前版本，则删除该要求。

考虑下面的例子。假设发现 C 1.4 存在问题，因此我们降级到 C 1.3。 C 1.4 已从模块图中删除。 B 1.2 也被删除，因为它需要 C 1.4 或更高版本。主模块对B的要求改为1.1。

![](images/downgrade.svg)

go get 还可以完全删除依赖项，在参数后使用 @none 后缀。这与降级类似。指定模块的所有版本都将从模块图中删除。

## Module graph pruning 模块图修剪

如果主模块为 go 1.17 或更高版本，则用于最小版本选择的模块图仅包含在其自己的 go.mod 文件中指定 go 1.17 或更高版本的每个模块依赖项的直接要求，除非该版本的模块也是（传递地）go 1.16 或更低版本的某些其他依赖项需要。 （go 1.17 依赖项的传递依赖项已从模块图中删除。）

由于 go 1.17 go.mod 文件包含在该模块中构建任何包或测试所需的每个依赖项的 require 指令，因此修剪后的模块图包含 go build 或 go test 明确要求的任何依赖项中的包所需的所有依赖项主模块。在给定模块中构建任何包或测试不需要的模块不会影响其包的运行时行为，因此从模块图中删除的依赖项只会导致其他不相关的模块之间的干扰。

其需求已被删除的模块仍然出现在模块图中，并且仍然由 go list -m all 报告：它们选择的版本是已知的并且定义良好，并且可以从这些模块加载包（例如，作为传递依赖项从其他模块加载的测试）。然而，由于 go 命令无法轻松识别满足这些模块的哪些依赖关系，因此 go build 和 go test 的参数不能包含来自其需求已被删除的模块的包。 go get 将包含每个命名包的模块提升为显式依赖项，允许在该包上调用 go build 或 go test。

由于 Go 1.16 及更早版本不支持模块图修剪，因此指定 go 1.16 或更低版本的每个模块仍包含依赖关系的完整传递闭包（包括传递 go 1.17 依赖关系）。 （在 go 1.16 及更低版本中，go.mod 文件仅包含直接依赖项，因此必须加载更大的图以确保包含所有间接依赖项。）

默认情况下，go mod tidy 为模块记录的 go.sum 文件包含低于其 go 指令中指定版本的 Go 版本所需的校验和。因此，go 1.17 模块包含 Go 1.16 加载的完整模块图所需的校验和，但 go 1.18 模块将仅包含 Go 1.17 加载的修剪模块图所需的校验和。 -compat 标志可用于覆盖默认版本（例如，在 go 1.17 模块中更积极地修剪 go.sum 文件）。

更多细节请参见设计文档( [the design document](https://go.googlesource.com/proposal/+/master/design/36460-lazy-module-loading.md))。

### Lazy module loading

为模块图修剪添加的更全面的要求还可以在模块内工作时实现另一种优化。如果主模块是 go 1.17 或更高版本，则 go 命令会避免加载完整的模块图，直到（除非）需要它。相反，它仅加载主模块的 go.mod 文件，然后尝试加载仅使用这些要求构建的包。如果在这些需求中没有找到要导入的包（例如，对主模块外部的包的测试的依赖项），则根据需要加载模块图的其余部分。

如果可以在不加载模块图的情况下找到所有导入的包，则 go 命令将仅加载包含这些包的模块的 go.mod 文件，并根据主模块的要求检查它们的要求，以确保它们在本地一致。 （由于版本控制合并、手动编辑以及使用本地文件系统路径替换的模块的更改，可能会出现不一致。）

## Workspaces 工作空间

工作区是磁盘上的模块集合，在运行最小版本选择 (MVS) 时用作主模块。

可以在 go.work 文件中声明工作空间，该文件指定工作空间中每个模块的模块目录的相对路径。当不存在 go.work 文件时，工作区由包含当前目录的单个模块组成。

大多数与模块一起使用的 go 子命令都对当前工作区确定的模块集进行操作。 go mod init、go modwhy、gomodedit、gomodtidy、gomodvendor 和 goget 始终在单个主模块上运行。

命令通过首先检查 GOWORK 环境变量来确定它是否位于工作区上下文中。如果 **GOWORK** 设置为关闭，该命令将位于单模块上下文中。如果为空或未提供，该命令将搜索当前工作目录，然后是连续的父目录，以查找文件 go.work。如果找到文件，该命令将在它定义的工作空间中运行；否则，工作区将仅包含包含工作目录的模块。如果 GOWORK 命名以 .work 结尾的现有文件的路径，则将启用工作空间模式。任何其他值都是错误的。您可以使用 go env GOWORK 命令来确定 go 命令正在使用哪个 go.work 文件。如果 go 命令不是工作空间模式，则 go env GOWORK 将为空。

### go.work files

工作空间由名为 go.work 的 UTF-8 编码文本文件定义。 go.work 文件是面向行的。每行包含一个指令，由关键字和参数组成。例如：

```go-module
go 1.18

use ./my/first/thing
use ./my/second/thing

replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
```

与 go.mod 文件中一样，可以从相邻行中分解出前导关键字来创建块。

```go-module
use (
    ./my/first/thing
    ./my/second/thing
)
```

go 命令提供了几个用于操作 go.work 文件的子命令。 go work init 创建新的 go.work 文件。 go work use 将模块目录添加到 go.work 文件中。 go work edit 执行低级编辑。 Go 程序可以使用 golang.org/x/mod/modfile 包以编程方式进行相同的更改。

### Lexical elements 词法元素

go.work 文件中的词汇元素的定义方式与 go.mod 文件中的定义方式完全相同。

### Grammar 语法

go.work 语法在下面使用扩展巴科斯-诺尔范式 (EBNF) 指定。有关 EBNF 语法的详细信息，请参阅 Go 语言规范中的符号部分。

```ebnf
GoWork = { Directive } .
Directive = GoDirective |
            ToolchainDirective |
            UseDirective |
            ReplaceDirective .
```

换行符、标识符和字符串分别用 newline、ident 和 string 表示。

模块路径和版本用 ModulePath 和 Version 表示。模块路径和版本的指定方式与 go.mod 文件完全相同。

```ebnf
ModulePath = ident | string . /* see restrictions above */
Version = ident | string .    /* see restrictions above */
```

### go directive

有效的 go.work 文件中需要有 go 指令。版本必须是有效的 Go 发行版本：正整数后跟一个点和一个非负整数（例如 1.18、1.19）。

go 指令指示 go.work 文件要使用的 go 工具链版本。如果对 go.work 文件格式进行更改，工具链的未来版本将根据其指示的版本解释该文件。

一个 go.work 文件最多可以包含一个 go 指令。

```ebnf
GoDirective = "go" GoVersion newline .
GoVersion = string | ident .  /* valid release version; see above */
```

例如：

```go-module
go 1.18
```

### toolchain directive

工具链指令声明要在工作区中使用的建议 Go 工具链。仅当默认工具链早于建议的工具链时，它才有效。
有关详细信息，请参阅“Go 工具链”。

```go-module
ToolchainDirective = "toolchain" ToolchainName newline .
ToolchainName = string | ident .  /* valid toolchain name; see “Go toolchains” */
```

例如：

```go-module
toolchain go1.21.0
```

### use directive

用户将磁盘上的模块添加到工作区中的主模块集中。它的参数是包含模块的 go.mod 文件的目录的相对路径。 use 指令不会添加包含在其参数目录的子目录中的模块。这些模块可以通过包含其 go.mod 文件的目录添加到单独的 use 指令中。

```go-module
UseDirective = "use" ( UseSpec | "(" newline { UseSpec } ")" newline ) .
UseSpec = FilePath newline .
FilePath = /* platform-specific relative or absolute file path */
```

例如：

```go-module
use ./mymod  // example.com/mymod

use (
    ../othermod
    ./subdir/thirdmod
)
```

### replace directive

与 go.mod 文件中的替换指令类似，go.work 文件中的替换指令将模块的特定版本或模块的所有版本的内容替换为在其他地方找到的内容。 go.work 中的通配符替换会覆盖 go.mod 文件中特定于版本的替换。

go.work 文件中的替换指令会覆盖工作区模块中相同模块或模块版本的任何替换。

```go
ReplaceDirective = "replace" ( ReplaceSpec | "(" newline { ReplaceSpec } ")" newline ) .
ReplaceSpec = ModulePath [ Version ] "=>" FilePath newline
            | ModulePath [ Version ] "=>" ModulePath Version newline .
FilePath = /* platform-specific relative or absolute file path */
```

例如：

```go-module
replace golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5

replace (
    golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5
    golang.org/x/net => example.com/fork/net v1.4.5
    golang.org/x/net v1.2.3 => ./fork/net
    golang.org/x/net => ./fork/net
)
```

## Compatibility with non-module repositories 与非模块存储库的兼容性

为了<mark>确保从 GOPATH 到模块的顺利过渡</mark>，go 命令可以通过添加 go.mod 文件，以模块感知模式(module-aware mode)从尚未迁移到模块的仓库中下载和编译软件包。

当 go 命令直接从存储库下载给定版本的模块时，它会查找模块路径的存储库 URL，将版本映射到存储库中的修订版，然后提取该修订版存储库的存档。<mark>如果模块的路径等于存储库根路径，并且存储库根目录不包含 go.mod 文件，则 go 命令会在模块缓存(module cache)中合成一个 go.mod 文件，其中包含模块指令而不包含其他内容。</mark>由于合成 go.mod 文件不包含其依赖项的 require 指令，因此依赖于它们的其他模块可能需要额外的 require 指令（带有 // indirect 注释），以确保在每个构建上以相同版本获取每个依赖项。

当 go 命令从代理下载模块时，它会与模块内容的其余部分分开下载 go.mod 文件。如果原始模块没有合成的 go.mod 文件，代理预计会提供合成的 go.mod 文件。

### +incompatible versions +incompatible版本

主版本以 2 或更高版本发布的模块必须在其模块路径上具有匹配的<mark>主版本后缀</mark>。例如，如果模块在 `v2.0.0` 版本发布，则其路径必须具有 /v2 后缀。这允许 go 命令将项目的多个主要版本视为不同的模块，即使它们是在同一存储库中开发的。

当 go 命令增加模块支持时，引入了主版本后缀要求，并且**许多存储库在此之前已经标记了主版本 2 或更高版本的版本**。**为了保持与这些存储库的兼容性**，**go 命令会在<mark>没有 go.mod 文件</mark>的情况下向<mark>主版本 2 或着更高</mark>的版本<mark>添加 +incompatible</mark> 的后缀**。 +incompatible表示该版本与主版本号较低的版本属于同一模块；因此，go 命令可能会自动升级到更高的不兼容版本，即使它可能会破坏构建。

考虑下面的示例要求：

```go-module
require example.com/m v4.1.2+incompatible
```

版本v4.1.2+incompatible指的是提供模块example.com/m的存储库中的语义版本标签([semantic version tag](https://go.dev/ref/mod#glos-semantic-version-tag))v4.1.2。该模块必须位于存储库根目录中（即存储库根路径也必须是 example.com/m），**并且不得存在 go.mod 文件**。该模块可能具有较低主版本号的版本，例如 v1.5.2，并且 go 命令可能会自动从这些版本升级到v4.1.2+incompatible（有关升级如何工作的信息，请参阅最小版本选择 (MVS)）。

在版本 v2.0.0 被标记后迁移到模块的存储库通常应该发布新的主要版本。在上面的示例中，作者应创建一个路径为 example.com/m/v5 的模块，并应发布版本 v5.0.0。作者还应该更新模块中包的导入，以使用前缀 example.com/m/v5 而不是 example.com/m。有关更详细的示例，请参阅 Go 模块：v2 及更高版本([Go Modules: v2 and Beyond](https://blog.golang.org/v2-go-modules))。

请注意，+incompatible 后缀不应出现在版本库的标签上；类似 v4.1.2+incompatible 这样的标签将被忽略。后缀只会出现在 go 命令使用的版本中。有关版本和标签之间区别的详情，请参阅将版本映射到提交( [Mapping versions to commits](https://go.dev/ref/mod#vcs-version))。

另请注意，+incompatible 后缀可能会出现在伪版本中。例如，v2.0.1-20200722182040-012345abcdef+incompatible 可能是一个有效的伪版本。

### Minimal module compatibility

以主要版本 2 或更高版本发布的模块需要在其模块路径上具有主要版本后缀。该模块可能会也可能不会在其存储库内的主要版本子目录中开发。这对于构建 GOPATH 模式时在模块内导入包的包有影响。

通常在 GOPATH 模式下，包存储在与其存储库根路径匹配的目录中，并与其在存储库中的目录相连接。例如，存储库中子目录 sub 中根路径为 example.com/repo 的包将存储在 $GOPATH/src/example.com/repo/sub 中，并将作为 example.com/repo/sub 导入。

对于带有主版本后缀的模块，人们可能期望在目录 $GOPATH/src/example.com/repo/v2/sub 中找到包 example.com/repo/v2/sub。这需要在其存储库的 v2 子目录中开发该模块。 go 命令支持这一点，但不要求它（请参阅将版本映射到提交）。

如果模块不是在主版本子目录中开发的，那么它在 GOPATH 中的目录将不包含主版本后缀，并且其包可能会在没有主版本后缀的情况下导入。在上面的示例中，该包将在目录 $GOPATH/src/example.com/repo/sub 中找到，并将作为 example.com/repo/sub 导入。

这给打算在模块模式和 GOPATH 模式下构建的包带来了一个问题：模块模式需要后缀，而 GOPATH 模式则不需要。

为了解决这个问题，Go 1.11 中添加了最小模块兼容性，并向后移植到 Go 1.9.7 和 1.10.3。当导入路径解析为 GOPATH 模式下的目录时：

- 解析 `$modpath/$vn/$dir` 形式的导入时，其中：
  
  - $modpath 是有效的module path,
  
  - $vn 是主版本后缀,
  
  - $dir 可能是一个空子目录，

- 如果以下所有条件均为真：
  
  - 包` $modpath/$vn/$dir` 不存在于任何相关的vendor目录中。
  
  - go.mod 文件与导入文件位于同一目录中，或者位于 $GOPATH/src 根目录之前的任何父目录中，
  
  - 不存在 `$GOPATH[i]/src/$modpath/$vn/$suffix` 目录（对于任何根 $GOPATH[i]），
  
  - 文件 `$GOPATH[d]/src/$modpath/go.mod` 存在（对于某些根 `$GOPATH[d]`）并将模块路径声明为 `$modpath/$vn`，

- 然后`$modpath/$vn/$dir`的导入被解析到目录`$GOPATH[d]/src/$modpath/$dir`。

此规则允许已迁移到模块的包导入在 GOPATH 模式下构建时已迁移到模块的其他包，即使未使用主版本子目录也是如此。

## Module-aware commands 模块命令

大多数 go 命令可以在模块感知模式或 GOPATH 模式下运行。在模块感知模式下，go 命令使用 go.mod 文件来查找版本化依赖项，它通常从模块缓存中加载包，如果模块丢失则下载模块。在GOPATH模式下，go命令会忽略模块；它会在供应商目录和 GOPATH 中查找依赖项。

从 Go 1.16 开始，默认启用模块感知模式，无论 go.mod 文件是否存在。在较低版本中，当当前目录或任何父目录中存在 go.mod 文件时，会启用模块感知模式。

模块感知模式可以使用 GO111MODULE 环境变量进行控制，该变量可以设置为on、off 或auto。

- 如果 GO111MODULE=off，go 命令会忽略 go.mod 文件并在 GOPATH 模式下运行。

- 如果 GO111MODULE=on 或未设置，则即使不存在 go.mod 文件，go 命令也会以模块感知模式运行。并非所有命令都可以在没有 go.mod 文件的情况下工作：请参阅模块外部的模块命令( [Module commands outside a module](https://go.dev/ref/mod#commands-outside))。

- 如果 GO111MODULE=auto，则当当前目录或任何父目录中存在 go.mod 文件时，go 命令将以模块感知模式运行。在 Go 1.15 及更低版本中，这是默认行为。即使不存在 go.mod 文件，go mod 子命令和 go install 也会在模块感知模式下运行版本查询。

在模块感知模式下，GOPATH 不再定义构建期间导入的含义，但它仍然存储下载的依赖项（在 GOPATH/pkg/mod 中；请参阅模块缓存[Module cache](https://go.dev/ref/mod#module-cache)）和安装的命令（在 GOPATH/bin 中，除非设置了 GOBIN） 。

### Build commands Build命令

所有加载有关包的信息的命令都是模块感知的。这包括：

- `go build`
- `go fix`
- `go generate`
- `go install`
- `go list`
- `go run`
- `go test`
- `go vet`

当在模块感知模式下运行时，这些命令使用 go.mod 文件来解释命令行上列出的或 Go 源文件中写入的导入路径。这些命令接受所有模块命令通用的以下标志。

- -mod 标志控制 go.mod 是否可以自动更新以及是否使用vendor目录。
  
  - -mod=mod 告诉 go 命令忽略vendor目录并自动更新 go.mod，例如，当任何已知模块未提供导入的包时。[automatically update](https://go.dev/ref/mod#go-mod-file-updates)
  
  - -mod=readonly 告诉 go 命令忽略vendor目录，并在需要更新 go.mod 时报告错误。
  
  - -mod=vendor 告诉 go 命令使用vendor目录。在这种模式下，go命令不会使用网络或模块缓存。
  
  - <mark>默认情况下</mark>，如果 go.mod 中的 go 版本为 1.14 或更高版本并且存在vendor目录，则 go 命令的行为就像使用了 -mod=vendor 一样。否则，go 命令的行为就像使用了 -mod=readonly 一样。
  
  - go get 拒绝接受这个标志，因为命令的目的是修改依赖关系，只有 -mod=mod 才允许修改依赖关系。

- -modcacherw 标志指示 go 命令在模块缓存中创建具有读写权限的新目录，而不是将其设置为只读。当一致使用此标志时（通常通过在环境中设置 GOFLAGS=-modcacherw 或运行 go env -w GOFLAGS=-modcacherw），可以使用 rm -r 等命令删除模块缓存，而无需先更改权限。 go clean -modcache 命令可用于删除模块缓存，无论是否使用 -modcacherw。

- -modfile=file.mod 标志指示 go 命令读取（也可能写入）模块根目录中的备用文件而不是 go.mod。文件名必须以 .mod 结尾。名为 go.mod 的文件必须仍然存在才能确定模块根目录，但无法访问该文件。当指定 -modfile 时，还会使用备用 go.sum 文件：其路径是通过修剪 .mod 扩展名并附加 .sum 从 -modfile 标志派生的。

### Vendoring

使用模块时，go 命令通常通过将模块从源下载到模块缓存中来满足依赖关系，然后从这些下载的副本中加载包。 Vendoring 可用于允许与旧版本的 Go 进行互操作，或确保用于构建的所有文件都存储在单个文件树中。

go mod vendor 命令会在主模块根目录下创建一个名为 vendor 的目录，其中包含构建和测试主模块软件包所需的所有软件包副本。仅在主模块之外的软件包测试中导入的软件包不包括在内。与 go mod tidy 和其他模块命令一样，在构建 vendor 目录时，除了忽略（ignore）之外，不考虑其他构建约束。

go mod vendor还创建文件vendor/modules.txt，其中包含vendored包的列表以及从中复制它们的模块版本。启用vendoring后，此清单将用作模块版本信息的来源，如 go list -m 和 go version -m 所报告的那样。当go命令读取vendor/modules.txt时，它会检查模块版本是否与go.mod一致。如果自生成vendor/modules.txt以来go.mod发生了变化，go命令将报告错误。应再次运行 go modvendor 以更新vendor目录。

如果vendor目录存在于主模块的根目录中，并且主模块的 go.mod 文件中的 go 版本为 1.14 或更高版本，则会自动使用该目录。要显式启用vendor，请使用标志 -mod=vendor 调用 go 命令。要禁用供应商，请使用标志 -mod=readonly 或 -mod=mod。

启用vendoring后，go build 和 go test 等构建命令会从vendor目录加载包，而不是访问网络或本地模块缓存。 go list -m 命令仅打印有关 go.mod 中列出的模块的信息。当启用vendor时，go mod 命令（例如 go mod download 和 go mod tidy）的工作方式没有什么不同，并且仍然会下载模块并访问模块缓存。当启用vendor时，go get 的工作方式也没有什么不同。

与 GOPATH 模式下的vendor不同，go 命令忽略主模块根目录以外位置的vendor目录。此外，由于未使用其他模块中的vendor目录，因此 go 命令在构建模块 zip 文件时不包含vendor目录（但请参阅已知错误 [#31562](https://go.dev/issue/31562) 和[#37397](https://go.dev/issue/37397)）。

### go get

```shell
go get [-d] [-t] [-u] [build flags] [packages]
```

例子：

```shell
# Upgrade a specific module.
$ go get golang.org/x/net

# Upgrade modules that provide packages imported by packages in the main module.
$ go get -u ./...

# Upgrade or downgrade to a specific version of a module.
$ go get golang.org/x/text@v0.3.2

# Update to the commit on the module's master branch.
$ go get golang.org/x/text@master

# Remove a dependency on a module and downgrade modules that require it
# to versions that don't require it.
$ go get golang.org/x/text@none

# Upgrade the minimum required Go version for the main module.
$ go get go

# Upgrade the suggested Go toolchain, leaving the minimum Go version alone.
$ go get toolchain

# Upgrade to the latest patch release of the suggested Go toolchain.
$ go get toolchain@patch
```

go get 命令更新主模块的 go.mod 文件中的模块依赖项，然后构建并安装命令行上列出的包。

第一步是确定要更新哪些模块。 go get 接受包列表、包模式和模块路径作为参数。如果指定了包参数，则 go get 更新提供该包的模块。如果指定了包模式（例如，全部或带有 ... 通配符的路径），go get 会将模式扩展为一组包，然后更新提供包的模块。如果参数命名的是模块而不是包（例如，模块 golang.org/x/net 的根目录中没有包），则 go get 将更新模块但不会构建包。如果没有指定参数， go get 的行为就像 if 。被指定（当前目录中的包）；这可以与 -u 标志一起使用来更新提供导入包的模块。

每个参数可能包含一个版本查询后缀，指示所需的版本，如 go get golang.org/x/text@v0.3.0 中。版本查询后缀由 @ 符号后跟版本查询组成，它可能表示特定版本 (v0.3.0)、版本前缀 (v0.3)、分支或标记名称 (master)、修订版本 (1234abcd) ，或最新、升级、补丁或无特殊查询之一。如果没有给出版本，go get 使用@upgrade 查询。

一旦 go get 将其参数解析为特定模块和版本，go get 将在主模块的 go.mod 文件中添加、更改或删除 require 指令，以确保模块将来保持所需的版本。请注意，go.mod 文件中所需的版本是最低版本，并且可能会随着新依赖项的添加而自动增加。有关如何选择版本以及如何通过模块感知命​​令解决冲突的详细信息，请参阅最小版本选择 (MVS)。

当添加、升级或降级命令行上指定的模块时，如果指定模块的新版本需要更高版本的其他模块，则可能会升级其他模块。例如，假设模块 example.com/a 升级到版本 v1.5.0，并且该版本需要版本 v1.2.0 的模块 example.com/b。如果当前版本 v1.1.0 需要模块 example.com/b，则 go get example.com/a@v1.5.0 也会将 example.com/b 升级到 v1.2.0。

![](images/get-upgrade.svg)

当命令行上指定的模块被降级或删除时，其他模块可能会被降级。为了继续上面的示例，假设模块 example.com/b 降级到 v1.1.0。模块 example.com/a 也将降级到需要 v1.1.0 或更低版本的 example.com/b 的版本。

![](images/get-downgrade.svg)

可以使用版本后缀@none 删除模块要求。这是一种特殊的降级。依赖于已删除模块的模块将根据需要降级或删除。即使主模块中的包导入了一个或多个模块需求，也可以删除该模块需求。在这种情况下，下一个构建命令可能会添加新的模块要求。

如果一个模块需要两个不同版本（在命令行参数中明确指定或满足升级和降级），go get 将报告错误。

go get 选择一组新版本后，它会检查任何新选择的模块版本或任何提供在命令行上命名的包的模块是否被撤回或弃用。 go get 为它找到的每个撤回版本或已弃用的模块打印一条警告。 go list -m -u all 可用于检查所有依赖项中的撤回和弃用。

go get 更新 go.mod 文件后，它会构建在命令行上命名的包。可执行文件将安装在 GOBIN 环境变量指定的目录中，如果未设置 GOPATH 环境变量，则默认为 `$GOPATH/bin` 或 `$HOME/go/bin`。

go get 支持以下标志：

- <mark>-d 标志</mark>告诉 go get 不要构建或安装包。当使用 -d 时，go get 只会管理 go.mod 中的依赖项。不推荐使用不带 -d 的 go get 来构建和安装软件包（从 Go 1.17 开始）。在 Go 1.18 中，-d 将始终启用。

- -u 标志告诉 go get 升级提供由命令行上命名的包直接或间接导入的包的模块。 -u 选择的每个模块都将升级到其最新版本，除非更高版本（预发行版）已需要该模块。

- -u=patch 标志（不是 -u patch）还告诉 go get 升级依赖项，但 go get 会将每个依赖项升级到最新的补丁版本（类似于 @patch 版本查询）。

- -t 标志告诉 go get 考虑构建命令行上命名的包的测试所需的模块。当 -t 和 -u 一起使用时，go get 也会更新测试依赖项。

- 不应再使用 -insecure 标志。它允许 go get 解析自定义导入路径并使用不安全的方案（例如 HTTP）从存储库和模块代理中获取。 GOINSECURE 环境变量提供了更细粒度的控制，应该改为使用。

从 Go 1.16 开始，推荐使用 go install 命令来构建和安装程序。当与版本后缀（如 @latest 或 @v1.4.6）一起使用时， go install 以模块感知模式构建包，忽略当前目录或任何父目录（如果有）中的 go.mod 文件。

go get 更专注于管理 go.mod 中的需求。 -d 标志已被弃用，在 Go 1.18 中，它将始终启用。

### go install

```go
go install [build flags] [packages]
```

例子：

```shell
# Install the latest version of a program,
# ignoring go.mod in the current directory (if any).
$ go install golang.org/x/tools/gopls@latest

# Install a specific version of a program.
$ go install golang.org/x/tools/gopls@v0.6.4

# Install a program at the version selected by the module in the current directory.
$ go install golang.org/x/tools/gopls

# Install all programs in a directory.
$ go install ./cmd/...
```

go install 命令构建并安装由命令行上的路径命名的包。可执行文件（主包）安装到 GOBIN 环境变量指定的目录中，如果未设置 GOPATH 环境变量，则默认为 \$GOPATH/bin 或 \$HOME/go/bin。 $GOROOT 中的可执行文件安装在 \$GOROOT/bin 或 \$GOTOOLDIR 而不是 \$GOBIN 中。不可执行的包已构建并缓存，但未安装。

从 Go 1.16 开始，如果参数具有版本后缀（如 @latest 或 @v1.0.0）， go install 将在模块感知模式下构建包，<mark>忽略当前目录或任何父目录（如果有）中的 go.mod 文件。这对于安装可执行文件而不影响主模块的依赖关系很有用。</mark>

为了消除构建过程中使用的模块版本的歧义，参数必须满足以下限制条件：

- 参数必须是包路径或包模式（带有“...”通配符）。它们不能是标准包（如 fmt）、元模式（std、cmd、all）或相对或绝对文件路径。

- 所有参数必须具有相同的版本后缀。不允许不同的查询，即使它们引用相同的版本。

- 所有参数必须引用同一版本的同一模块中的包。

- 包路径参数必须引用主包。模式参数仅匹配主包。

- 没有模块被视为主模块。
  
  - 如果包含在命令行上命名的包的模块有一个 go.mod 文件，则它不能包含指令（替换和排除），否则如果它是主模块，则会导致它被不同地解释。
  
  - 该模块不能要求其自身的更高版本。
  
  - 任何模块中都不使用vendor目录。 （vendor目录不包含在模块 zip 文件中，因此 go install 不会下载它们。）

有关支持的版本查询语法，请参阅版本查询( [Version queries](https://go.dev/ref/mod#version-queries))。 Go 1.15 及更低版本不支持在 go install 中使用版本查询。

如果参数没有版本后缀，go install 可能会在模块感知模式或 GOPATH 模式下运行，具体取决于 GO111MODULE 环境变量和 go.mod 文件是否存在。有关详细信息，请参阅模块感知命​​令。如果启用了模块感知模式，则 go install 在主模块的上下文中运行，该模块可能与包含正在安装的包的模块不同。

### go list -m

```shell
go list -m [-u] [-retracted] [-versions] [list flags] [modules]
```

```shell
$ go list -m all
$ go list -m -versions example.com/m
$ go list -m -json example.com/m@latest
```

-m 标志使 go list 列出模块而不是包。在这种模式下，go list 的参数可以是模块、模块模式（包含 ... 通配符）、版本查询或特殊模式 all，它匹配构建列表中的所有模块。如果未指定参数，则列出主模块。

列出模块时， -f 标志仍然指定应用于 Go 结构的格式模板，但现在是 Module 结构：

```go
type Module struct {
    Path       string        // module path
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
}

type ModuleError struct {
    Err string // the error itself
}
```

默认输出是打印模块路径，然后打印有关版本和替换（如果有）的信息。例如， go list -m all 可能会打印：

```shell
example.com/main/module
golang.org/x/net v0.1.0
golang.org/x/text v0.3.0 => /tmp/text
rsc.io/pdf v0.1.1
```

Module 结构有一个 String 方法，用于格式化此行输出，因此默认格式相当于 -f '{{.String}}'。

请注意，当模块被替换时，其 Replace 字段描述替换模块，并且其 Dir 字段设置为替换模块的源代码（如果存在）。 （也就是说，如果 Replace 为非零，则 Dir 设置为 Replace.Dir，无法访问替换的源代码。）

-u 标志添加有关可用升级的信息。当给定模块的最新版本比当前版本更新时，list -u 将模块的 Update 字段设置为有关较新模块的信息。 list -u 还打印当前选定的版本是否已撤回以及该模块是否已弃用。该模块的 String 方法通过在当前版本后面的括号中格式化新版本来指示可用的升级。例如， go list -m -u all 可能会打印：

```shell
example.com/main/module
golang.org/x/old v1.9.9 (deprecated)
golang.org/x/net v0.1.0 (retracted) [v0.2.0]
golang.org/x/text v0.3.0 [v0.4.0] => /tmp/text
rsc.io/pdf v0.1.1 [v0.1.2]
```

（对于工具， go list -m -u -json all 可能更方便解析。）

**-versions 标志**使 list 将模块的 Versions 字段设置为该模块的所有已知版本的列表，根据语义版本控制从最低到最高排序。该标志还更改默认输出格式以显示模块路径，后跟空格分隔的版本列表。除非还指定了 -retracted 标志，否则此列表中将省略已收回的版本。

-retracted 标志指示 list 在使用 -versions 标志打印的列表中显示收回的版本，并在解析版本查询时考虑收回的版本。例如， go list -m -retracted example.com/m@latest 显示模块 example.com/m 的最高版本或预发布版本，即使该版本已撤回。在此版本中，retract 指令和弃用是从 go.mod 文件加载的。 -retracted 标志是在 Go 1.16 中添加的。

模板函数 module 接受一个字符串参数，该参数必须是模块路径或查询，并以 Module 结构体的形式返回指定的模块。如果发生错误，结果将是一个带有非零错误字段的模块结构。

### go mod download

```shell
go mod download [-x] [-json] [-reuse=old.json] [modules]
```

```shell
$ go mod download
$ go mod download golang.org/x/mod@v0.2.0
```

go mod download 命令将指定模块下载到模块缓存中。参数可以是模块路径或模块模式，选择主模块的依赖项或path@version形式的版本查询。如果没有参数，下载适用于主模块的所有依赖项。

go命令在普通执行过程中会根据需要自动下载模块。 go mod download 命令主要用于预填充模块缓存或加载由模块代理提供的数据。

默认情况下，下载不会向标准输出写入任何内容。它将进度消息和错误打印到标准错误。

-json 标志使 download 将一系列 JSON 对象打印到标准输出，描述每个下载的模块（或失败），对应于以下 Go 结构：

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

-<mark>x 标志</mark>导致 download 将 download 执行的命令打印到标准错误。

-reuse 标志接受包含先前“go mod download -json”调用的 JSON 输出的文件名。 go 命令可以使用此文件来确定模块自上次调用以来未发生更改，并避免重新下载它。通过将 Reuse 字段设置为 true，将在新输出中标记未重新下载的模块。通常模块缓存会自动提供这种重用； -reuse 标志对于不保留模块缓存的系统很有用。

### go mod edit

```shell
go mod edit [editing flags] [-fmt|-print|-json] [go.mod]
```

```shell
# Add a replace directive.
$ go mod edit -replace example.com/a@v1.0.0=./a

# Remove a replace directive.
$ go mod edit -dropreplace example.com/a@v1.0.0

# Set the go version, add a requirement, and print the file
# instead of writing it to disk.
$ go mod edit -go=1.14 -require=example.com/m@v1.0.0 -print

# Format the go.mod file.
$ go mod edit -fmt

# Format and print a different .mod file.
$ go mod edit -print tools.mod

# Print a JSON representation of the go.mod file.
$ go mod edit -json
```

go mod edit 命令提供了一个用于编辑和格式化 go.mod 文件的命令行界面，主要供工具和脚本使用。 go mod edit 只读取一个 go.mod 文件；它不会查找有关其他模块的信息。默认情况下，go mod edit 读取和写入主模块的 go.mod 文件，但可以在编辑标志后指定不同的目标文件。

编辑标志指定了编辑操作的顺序。

- -module 标志更改模块的路径（go.mod 文件的模块行）。

- -go=version 标志设置预期的 Go 语言版本。

- -require=path@version 和 -droprequire=path 标志添加和删除对给定模块路径和版本的要求。请注意，-require 会覆盖路径上的任何现有要求。这些标志主要用于理解模块图的工具。用户应该更喜欢 go get path@version 或 go get path@none，它们根据需要进行其他 go.mod 调整以满足其他模块施加的约束。查看[`go get`](https://go.dev/ref/mod#go-get)。

- -exclude=path@version 和 -dropexclude=path@version 标志添加和删除给定模块路径和版本的排除项。请注意，如果该排除已存在，则 -exclude=path@version 是无操作。

- -replace=old[@v]=new[@v] 标志添加给定模块路径和版本对的替换。如果old@v中的@v被省略，则会添加左侧没有版本的替换，适用于旧模块路径的所有版本。如果new@v中的@v被省略，则新路径应该是本地模块根目录，而不是模块路径。请注意，-replace 会覆盖旧[@v] 的任何冗余替换，因此省略 @v 将删除特定版本的替换。

- -dropreplace=old[@v] 标志删除给定模块路径和版本对的替换。如果提供了@v，则删除给定版本的替换。左侧没有版本的现有替换件仍可以替换该模块。如果省略@v，则删除没有版本的替换。

- -retract=version 和 -dropretract=version 标志添加和删除给定版本的撤回，该版本可以是单个版本（如 v1.2.3）或一个间隔（如 [v1.1.0,v1.2.0]）。请注意，-retract 标志无法为撤回指令添加基本原理注释。建议提供基本原理注释，并且可以通过 go list -m -u 和其他命令显示。

编辑标志可以重复。更改将按给定的顺序应用。

go mod edit 有额外的标志来控制其输出。

- <mark>-fmt 标志</mark>重新格式化 go.mod 文件而不进行其他更改。使用或重写 go.mod 文件的任何其他修改也暗示了这种重新格式化。唯一需要此标志的情况是没有指定其他标志，如 go mod edit -fmt 中。

- -print 标志以文本格式打印最终的 go.mod，而不是将其写回磁盘。

- -json 标志以 JSON 格式打印最终的 go.mod，而不是以文本格式将其写回到磁盘。 JSON 输出对应于以下 Go 类型：

```go
type Module struct {
    Path    string
    Version string
}

type GoMod struct {
    Module  ModPath
    Go      string
    Require []Require
    Exclude []Module
    Replace []Replace
    Retract []Retract
}

type ModPath struct {
    Path       string
    Deprecated string
}

type Require struct {
    Path     string
    Version  string
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

请注意，这仅描述了 go.mod 文件本身，而不是间接引用的其他模块。对于可用于构建的完整模块集，请使用 go list -m -json all。请参阅 go list -m。

例如，工具可以通过解析 go mod edit -json 的输出来获取 go.mod 文件作为数据结构，然后可以通过使用 -require、-exclude 等调用 go mod edit 来进行更改。

工具还可以使用包 golang.org/x/mod/modfile 来解析、编辑和格式化 go.mod 文件。

### go mod graph

```shell
go mod graph [-go=version]
```

go mod graph 命令以文本形式打印模块需求图（已应用替换）。例如：

```shell
example.com/main example.com/a@v1.1.0
example.com/main example.com/b@v1.2.0
example.com/a@v1.1.0 example.com/b@v1.1.1
example.com/a@v1.1.0 example.com/c@v1.3.0
example.com/b@v1.1.0 example.com/c@v1.1.0
example.com/b@v1.2.0 example.com/c@v1.2.0
```

模块图中的每个顶点代表模块的特定版本。图中的每条边代表对依赖项的最低版本的要求。

go mod graph 打印图表的边缘，每行一个。每行都有两个以空格分隔的字段：模块版本及其依赖项之一。每个模块版本都被标识为路径@版本形式的字符串。主模块没有 @version 后缀，因为它没有版本。

-go 标志使 go mod graph 报告给定 Go 版本加载的模块图，而不是 go.mod 文件中 go 指令指示的版本。

有关如何选择版本的更多信息，请参阅最小版本选择 (MVS)。另请参阅 go list -m 来打印选定的版本，并参阅 go mod Why 来了解为什么需要模块。

### go mod init

```shell
go mod init [module-path]
```

```shell
go mod init
go mod init example.com/m
```

go mod init 命令初始化并在当前目录中写入一个新的 go.mod 文件，实际上创建了一个以当前目录为根的新模块。 go.mod 文件不得已存在。

init 接受一个可选参数，即新模块的模块路径。有关选择模块路径的说明，请参阅模块路径。如果省略模块路径参数，init 将尝试使用 .go 文件中的导入注释、供应商工具配置文件和当前目录（如果在 GOPATH 中）来推断模块路径。

如果存在vendor工具的配置文件，init 将尝试从中导入模块需求。 init 支持以下配置文件。

- `GLOCKFILE` (Glock)
- `Godeps/Godeps.json` (Godeps)
- `Gopkg.lock` (dep)
- `dependencies.tsv` (godeps)
- `glide.lock` (glide)
- `vendor.conf` (trash)
- `vendor.yml` (govend)
- `vendor/manifest` (gvt)
- `vendor/vendor.json` (govendor)

Vendoring 工具配置文件并不总是能够以完美的保真度进行翻译。例如，如果在同一存储库中导入多个不同版本的包，并且该存储库仅包含一个模块，则导入的 go.mod 只能需要一个版本的模块。您可能希望运行 go list -m all 来检查构建列表中的所有版本，并运行 go mod tidy 来添加缺少的需求并删除未使用的需求。

### go mod tidy

```shell
go mod tidy [-e] [-v] [-go=version] [-compat=version]
```

go mod tidy 确保 go.mod 文件与模块中的源代码匹配。它添加了构建当前模块的包和依赖项所需的任何缺失的模块要求，并删除了对不提供任何相关包的模块的要求。它还将所有缺失的条目添加到 go.sum 并删除不必要的条目。

-e 标志（在 Go 1.16 中添加）会导致 go mod tidy 尝试继续，尽管在加载包时遇到错误。

-v 标志使 go mod tidy 将有关已删除模块的信息打印到标准错误。

go mod tidy 的工作原理是递归地加载主模块中的所有包以及它们导入的所有包。这包括测试导入的包（包括其他模块中的测试）。 go mod tidy 的作用就好像所有构建标签都已启用，因此它会考虑特定于平台的源文件和需要自定义构建标签的文件，即使这些源文件通常不会构建。有一个例外：忽略构建标记未启用，因此具有构建约束 // +buildignore 的文件将不会被考虑。请注意，go mod tidy 不会考虑主模块中名为 testdata 的目录或名称以 .或 _ 除非这些包是由其他包显式导入的。

一旦 go mod tidy 加载了这组包，它就会确保提供一个或多个包的每个模块在主模块的 go.mod 文件中都有一个 require 指令，或者（如果主模块是 go 1.16 或更低版本）需要另一个必需的模块。 go mod tidy 将添加对每个缺失模块的最新版本的要求（有关最新版本的定义，请参阅版本查询）。 go mod tidy 将删除不提供上述集合中任何包的模块的 require 指令。

go mod tidy 还可以添加或删除require 指令上 // indirect 的注释。 // indirect 注释表示模块不提供由主模块中的包导入的包。 （有关何时添加// indirect 依赖项和注释的更多详细信息，请参阅 require 指令。）

如果设置了 -go 标志，go mod tidy 会将 go 指令更新为指定版本，根据该版本启用或禁用模块图修剪和延迟模块加载（并根据需要添加或删除间接需求）。

默认情况下，当 go 指令中指示的版本之前的 Go 版本加载模块图时，go mod tidy 将检查所选模块版本是否不会更改。还可以通过 -compat 标志显式指定检查兼容性的版本。

### go mod vendor

```shell
go mod vendor [-e] [-v] [-o]
```

go mod vendor 命令在主模块的根目录中构造一个名为vendor的目录，其中包含支持主模块中包的构建和测试所需的所有包的副本。不包括仅通过主模块外部的包测试导入的包。与 go mod tidy 和其他模块命令一样，构建vendor目录时不考虑除忽略之外的构建约束[build constraints](https://go.dev/ref/mod#glos-build-constraint)。

启用vendoring后，go 命令将从vendor目录加载包，而不是将模块从其源下载到模块缓存中并使用这些下载的副本。有关更多信息，请参阅[Vendoring](https://go.dev/ref/mod#vendoring) 。

go mod vendor 还创建文件vendor/modules.txt，其中包含vendor包的列表以及从中复制它们的模块版本。启用vendoring后，此清单将用作模块版本信息的来源，如 go list -m 和 go version -m 所报告的那样。当go命令读取vendor/modules.txt时，它会检查模块版本是否与go.mod一致。如果 go.mod 在生成了 vendor/modules.txt 后发生了变化，则应该再次运行 go mod vendor。

请注意，在重新构建之前，go mod vendor 会删除vendor 目录（如果存在）。不应对vendor的包进行本地更改。 go 命令不会检查vendor目录中的包是否未被修改，但可以通过运行 go mod vendor 并检查是否没有进行任何更改来验证vendor目录的完整性。

-e 标志（在 Go 1.16 中添加）会导致 go mod 供应商尝试继续，尽管在加载包时遇到错误。

-v 标志使 go mod供应商将供应的模块和包的名称打印到标准错误。

-o 标志（在 Go 1.18 中添加）使 go modvendor 在指定目录而不是供应商处输出供应商树。参数可以是绝对路径或相对于模块根的路径。

### go mod verify

```shell
go mod verify
```

go mod verify 检查存储在模块缓存中的主模块的依赖关系自下载以来尚未被修改。要执行此检查， go mod verify 对每个下载的模块 .zip 文件和提取的目录进行哈希处理，然后将这些哈希值与首次下载模块时记录的哈希值进行比较。 go mod verify 检查构建列表中的每个模块（可以使用 go list -m all 打印）。

如果所有模块均未修改，则 go mod verify 打印“所有模块已验证”。否则，它会报告哪些模块已更改并以非零状态退出。

请注意，所有模块感知命​​令都会验证主模块的 go.sum 文件中的哈希值是否与下载到模块缓存中的模块记录的哈希值相匹配。如果 go.sum 中缺少哈希值（例如，因为该模块是第一次使用），则 go 命令将使用校验和数据库验证其哈希值（除非模块路径与 GOPRIVATE 或 GONOSUMDB 匹配）。有关详细信息，请参阅验证模块。

相反，go mod verify 检查模块 .zip 文件及其提取的目录是否具有与首次下载时模块缓存中记录的哈希值相匹配的哈希值。这对于在下载并验证模块后检测模块缓存中文件的更改非常有用。 go mod verify 不会下载不在缓存中的模块的内容，并且它不会使用 go.sum 文件来验证模块内容。但是，go mod verify 可能会下载 go.mod 文件以执行最小版本选择。它将使用 go.sum 来验证这些文件，并且可能会为丢失的哈希值添加 go.sum 条目。

### go mod why

```shell
go mod why [-m] [-vendor] packages...
```

go mod why 显示导入图中从主模块到每个列出的包的最短路径。

输出是一系列节，每个节对应命令行上命名的每个包或模块，以空行分隔。每个节都以注释行开头，以 # 开头，给出目标包或模块。后续行给出了通过导入图的路径，每行一个包。如果主模块未引用包或模块，则该节将显示一个带括号的注释来指示这一事实。

例如：

```shell
$ go mod why golang.org/x/text/language golang.org/x/text/encoding
# golang.org/x/text/language
rsc.io/quote
rsc.io/sampler
golang.org/x/text/language

# golang.org/x/text/encoding
(main module does not need package golang.org/x/text/encoding)
```

-m 标志导致 go mod 为什么将其参数视为模块列表。 go mod Why 将打印每个模块中任何包的路径。请注意，即使使用 -m，go mod Why 也会查询包图，而不是 go mod graph 打印的模块图。

-vendor 标志导致 go mod 为什么忽略主模块之外的包测试中的导入（就像 go mod供应商所做的那样）。默认情况下，go mod Why 会考虑与 all 模式匹配的包图。在声明 go 1.16 或更高版本的模块中（使用 go.mod 中的 go 指令），此标志在 Go 1.16 之后无效，因为 all 的含义已更改以匹配 go mod 供应商匹配的包集。

### go version -m

```shell
go version [-m] [-v] [file ...]
```

例如：

```shell
# Print Go version used to build go.
$ go version

# Print Go version used to build a specific executable.
$ go version ~/go/bin/gopls

# Print Go version and module versions used to build a specific executable.
$ go version -m ~/go/bin/gopls

# Print Go version and module versions used to build executables in a directory.
$ go version -m ~/go/bin/
```

go version 报告用于构建命令行上指定的每个可执行文件的 Go 版本。

如果命令行上没有命名文件，go version 会打印自己的版本信息。

如果指定了目录，go version 会递归地遍历该目录，查找可识别的 Go 二进制文件并报告其版本。默认情况下，go 版本不会报告目录扫描期间发现的无法识别的文件。 -v 标志导致它报告无法识别的文件。

-m 标志使 go version 打印每个可执行文件的嵌入模块版本信息（如果可用）。对于每个可执行文件，go version -m 都会打印一张包含制表符分隔列的表格，如下所示。

```shell
$ go version -m ~/go/bin/goimports
/home/jrgopher/go/bin/goimports: go1.14.3
        path    golang.org/x/tools/cmd/goimports
        mod     golang.org/x/tools      v0.0.0-20200518203908-8018eb2c26ba      h1:0Lcy64USfQQL6GAJma8BdHCgeofcchQj+Z7j0SXYAzU=
        dep     golang.org/x/mod        v0.2.0          h1:KU7oHjnv3XNWfa5COkzUifxZmxp1TyI7ImMXqFxLwvQ=
        dep     golang.org/x/xerrors    v0.0.0-20191204190536-9bdfabe68543      h1:E7g+9GITq07hpfrRu66IVDexMakfv52eLZ2CXBWiKr4=
```

表格的格式将来可能会发生变化。可以从runtime/debug.ReadBuildInfo 获得相同的信息。

表中每一行的含义由第一列中的单词确定。

- **path**: 用于构建可执行文件的主包的路径。

- **mod**:包含主包的模块。这些列分别是模块路径、版本和总和。主模块有版本（devel），没有sum。

- **dep**: 提供链接到可执行文件的一个或多个包的模块。与 mod 格式相同。

- **=>**: 替换前一行的模块。如果替换是本地目录，则仅列出目录路径（无版本或总和）。如果替换是模块版本，则会列出路径、版本和总和，就像 mod 和 dep 一样。被替换的模块没有总和。

### go clean -modcache

```go
go clean [-modcache]
```

-modcache 标志导致 go clean 删除整个模块缓存，包括版本化依赖项的解压源代码。

这通常是删除模块缓存的最佳方法。默认情况下，模块缓存中的大多数文件和目录都是只读的，以防止测试和编辑者在经过身份验证后无意中更改文件。不幸的是，这会导致像 rm -r 这样的命令失败，因为如果不先使其父目录可写，则无法删除文件。

-modcacherw 标志（被 go build 和其他模块感知命​​令接受）导致模块缓存中的新目录可写。要将 -modcacherw 传递给所有模块感知命​​令，请将其添加到 GOFLAGS 变量中。 GOGFLAGS 可以在环境中设置或使用 go env -w 设置。例如，以下命令将其永久设置：

```shell
go env -w GOFLAGS=-modcacherw
```

-modcacherw 应谨慎使用；开发人员应小心不要更改模块缓存中的文件。 go mod verify 可用于检查缓存中的文件是否与主模块的 go.sum 文件中的哈希值匹配。

### Version queries

有多个命令允许您使用版本查询指定模块的版本，版本查询出现在命令行上模块或包路径后面的 @ 字符之后。

例如：

```shell
go get example.com/m@latest
go mod download example.com/m@master
go list -m -json example.com/m@e3702bed2
```

版本查询可能是以下之一：

- 完全指定的语义版本，例如 v1.2.3，它选择特定版本。请参阅版本了解语法。

- 语义版本前缀，例如 v1 或 v1.2，它选择具有该前缀的最高可用版本。

- 语义版本比较，例如 <v1.2.3 或 >=v1.5.6，它选择与比较目标最接近的可用版本（> 和 >= 为最低版本，< 和 <= 为最高版本）。

- 底层源存储库的修订标识符，例如提交哈希前缀、修订标签或分支名称。如果修订版标记有语义版本，则此查询将选择该版本。否则，此查询将为基础提交选择伪版本。请注意，无法通过这种方式选择名称与其他版本查询匹配的分支和标签。例如，查询 v2 选择以 v2 开头的最新版本，而不是名为 v2 的分支。

- 字符串latest，选择最高的可用发行版本。如果没有发布版本，latest 将选择最高的预发布版本。如果没有标记版本，latest 会在存储库默认分支的顶端选择一个伪版本进行提交。

- 字符串升级，与最新版本类似，但如果当前需要模块的版本高于最新版本选择的版本（例如预发布版本），则升级将选择当前版本。

- 字符串补丁，选择与当前所需版本具有相同主要版本号和次要版本号的最新可用版本。如果当前不需要版本，则 patch 相当于最新版本。从 Go 1.16 开始，go get 在使用 patch 时需要当前版本（但 -u=patch 标志没有此要求）。

除了针对特定命名版本或修订的查询之外，所有查询都会考虑 go list -m -versions 报告的可用版本（请参阅 go list -m）。此列表仅包含标记版本，而不包含伪版本。不考虑主模块的 go.mod 文件中的排除指令不允许的模块版本。同一模块最新版本的 go.mod 文件中的撤回指令所覆盖的版本也会被忽略，除非 -retracted 标志与 go list -m 一起使用以及加载撤回指令时除外。

发布版本优先于预发布版本。例如，如果版本 v1.2.2 和 v1.2.3-pre 可用，则最新查询将选择 v1.2.2，即使 v1.2.3-pre 更高。 <v1.2.4 查询也会选择 v1.2.2，即使 v1.2.3-pre 更接近 v1.2.4。如果没有可用的发布或预发布版本，则最新、升级和补丁查询将为存储库默认分支顶部的提交选择一个伪版本。其他查询会报错。

### Module commands outside a module 模块外的模块命令

模块感知的 Go 命令通常在工作目录或父目录中的 go.mod 文件定义的主模块的上下文中运行。某些命令可能在没有 go.mod 文件的情况下以模块感知模式运行，但大多数命令的工作方式不同，或者在不存在 go.mod 文件时报告错误。

有关启用和禁用模块感知模式的信息，请参阅模块感知命​​令。

| Command         | Behavior                                                              |
| --------------- | --------------------------------------------------------------------- |
| go build        | 只能加载、导入和构建标准库中的包和在命令行中指定为 .go 文件的包。无法构建来自其他模块的包，因为没有地方记录模块需求并确保确定性构建。 |
| go doc          |                                                                       |
| go fix          |                                                                       |
| go fmt          |                                                                       |
| go generate     |                                                                       |
| go install      |                                                                       |
| go list         |                                                                       |
| go run          |                                                                       |
| go test         |                                                                       |
| go vet          |                                                                       |
| go get          | 包和可执行文件可以像平常一样构建和安装。请注意，当没有 go.mod 文件运行 go get 时，没有主模块，因此不应用替换和排除指令。  |
| go list -m      | 大多数参数都需要显式版本查询，除非使用 -versions 标志。                                     |
| go mod download | 大多数参数都需要显式版本查询。                                                       |
| go mod edit     | 需要显式文件参数。                                                             |
| go mod graph    | 这些命令需要 go.mod 文件，如果不存在，则会报告错误。                                        |
| go mod tidy     |                                                                       |
| go mod vendor   |                                                                       |
| go mod verify   |                                                                       |
| go mod why      |                                                                       |

### go work init

```go
go work init [moddirs]
```

Init 初始化并在当前目录中写入一个新的 go.work 文件，实际上是在当前目录中创建一个新的工作空间。

go work init 可以选择接受工作区模块的路径作为参数。如果省略该参数，将创建一个没有模块的空工作区。

每个参数路径都会添加到 go.work 文件中的 use 指令中。当前的 go 版本也将在 go.work 文件中列出。

### go work edit

```shell
go work edit [editing flags] [go.work]
```

go work edit 命令提供了用于编辑 go.work 的命令行界面，主要供工具或脚本使用。它只读取 go.work；它不会查找有关所涉及模块的信息。如果未指定文件，Edit 会在当前目录及其父目录中查找 go.work 文件

编辑标志指定编辑操作的序列。

- -fmt 标志重新格式化 go.work 文件而不进行其他更改。使用或重写 go.work 文件的任何其他修改也暗示了这种重新格式化。唯一需要此标志的情况是没有指定其他标志，如“go work edit -fmt”中。

- -use=path 和 -dropuse=path 标志从 go.work 文件的模块目录集中添加和删除 use 指令。

- -replace=old[@v]=new[@v] 标志添加给定模块路径和版本对的替换。如果old@v中的@v被省略，则会添加左侧没有版本的替换，适用于旧模块路径的所有版本。如果new@v中的@v被省略，则新路径应该是本地模块根目录，而不是模块路径。请注意，-replace 会覆盖旧[@v] 的任何冗余替换，因此省略 @v 将删除特定版本的现有替换。

- -dropreplace=old[@v] 标志删除给定模块路径和版本对的替换。如果省略@v，则左侧没有版本的替换将被删除。

- -go=version 标志设置预期的 Go 语言版本。

编辑标志可以重复。更改将按给定的顺序应用。

go work edit 有额外的标志来控制其输出

- -print 标志以文本格式打印最终的 go.work，而不是将其写回 go.mod。

- -json 标志以 JSON 格式打印最终的 go.work 文件，而不是将其写回 go.mod。 JSON 输出对应于以下 Go 类型：

```go
type Module struct {
    Path    string
    Version string
}

type GoWork struct {
    Go        string
    Directory []Directory
    Replace   []Replace
}

type Use struct {
    Path       string
    ModulePath string
}

type Replace struct {
    Old Module
    New Module
}
```

### go work use

```shell
go work use [-r] [moddirs]
```

go work use 命令提供了一个命令行界面，用于将目录（可选）递归添加到 go.work 文件。

如果磁盘上存在命令行 go.work 文件中列出的每个参数目录，则 use 指令将添加到 go.work 文件中；如果磁盘上不存在，则将其从 go.work 文件中删除。

-r 标志在参数目录中递归搜索模块，并且 use 命令的操作就像将每个目录指定为参数一样：即，将为存在的目录添加 use 指令，并为不存在的目录删除 use 指令。

### go work sync

```shell
go work sync
```

go work sync 命令将工作区的构建列表同步回工作区的模块。

工作区的构建列表是用于在工作区中进行构建的所有（传递）依赖模块的版本集。 go work sync 使用最小版本选择 (MVS) 算法生成构建列表，然后将这些版本同步回工作区中指定的每个模块（使用 use 指令）。

计算出工作区构建列表后，工作区中每个模块的 go.mod 文件都会被重写，并使用与该模块相关的依赖项进行升级以匹配工作区构建列表。请注意，最小版本选择可保证每个模块的构建列表版本始终与每个工作区模块中的版本相同或更高。

## Module proxies

### GOPROXY 协议

模块代理是一个 HTTP 服务器，可以响应下面指定路径的 GET 请求。这些请求没有查询参数，也不需要特定的标头，因此即使是从固定文件系统（包括 file:// URL）提供服务的站点也可以是模块代理。

成功的 HTTP 响应必须具有状态代码 200（正常）。遵循重定向 (3xx)。状态代码为 4xx 和 5xx 的响应被视为错误。错误代码 404（未找到）和 410（已消失）表示请求的模块或版本在代理上不可用，但可能在其他地方找到。错误响应的内容类型应为 text/plain，字符集为 utf-8 或 us-ascii。

go 命令可以配置为使用 GOPROXY 环境变量联系代理或源控制服务器，该变量接受代理 URL 列表。该列表可能包含关键字 direct 或 off（有关详细信息，请参阅环境变量）。列表元素可以用逗号 (,) 或竖线 (|) 分隔，这决定了错误回退行为。当 URL 后跟逗号时，只有在 404（未找到）或 410（已消失）响应后，go 命令才会回退到后面的源。当 URL 后跟管道时，go 命令会在发生任何错误（包括超时等非 HTTP 错误）后回退到后面的源。这种错误处理行为让代理充当未知模块的看门人。例如，对于不在批准列表中的模块，代理可能会响应错误 403（禁止）（请参阅服务私有模块的私有代理）。

下表指定了模块代理必须响应的查询。对于每个路径，$base 是代理 URL 的路径部分，$module 是模块路径，$version 是版本。例如，如果代理 URL 是 https://example.com/mod，并且客户端正在请求版本 v0.3.2 的模块 golang.org/x/text 的 go.mod 文件，则客户端将发送 GET请求 https://example.com/mod/golang.org/x/text/@v/v0.3.2.mod。

为了避免在不区分大小写的文件系统中提供服务时出现歧义，通过将每个大写字母替换为感叹号，后跟相应的小写字母，对 $module 和 $version 元素进行大小写编码。这允许模块 example.com/M 和 example.com/m 都存储在磁盘上，因为前者被编码为 example.com/!m。

| Path                              | Description                                                                                                                                                                                                                        |
| --------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| \$base/\$module/@v/list           | 以纯文本形式返回给定模块的已知版本列表，每行一个。此列表不应包含伪版本。<br/>返回有关模块的特定版本的 JSON 格式的元数据。响应必须是与以下 Go 数据结构相对应的 JSON 对象：<br/>type Info struct {<br/>       Version string    // version string<br/>        Time    time.Time // commit time<br/>}           |
| \$base/\$module/@v/\$version.info | 版本字段是必需的，并且必须包含有效的规范版本（请参阅版本）。请求路径中的$version不需要是相同的版本，甚至不需要是有效的版本；该端点可用于查找分支名称或修订标识符的版本。但是，如果 $version 是规范版本，且主版本与 $module 兼容，则成功响应中的 Version 字段必须相同。<br/>时间字段是可选的。如果存在，则它必须是 RFC 3339 格式的字符串。表示版本创建的时间。<br/>将来可能会添加更多字段，因此保留其他名称。 |
| \$base/\$module/@v/\$version.mod  | 返回特定版本模块的 go.mod 文件。如果模块没有所请求版本的 go.mod 文件，则必须返回仅包含具有所请求模块路径的模块语句的文件。否则，必须返回原始的、未修改的 go.mod 文件。                                                                                                                                    |
| \$base/\$module/@v/\$version.zip  | 返回包含特定版本模块内容的 zip 文件。有关如何格式化此 zip 文件的详细信息，请参阅模块 zip 文件。                                                                                                                                                                            |
| \$base/\$module/@latest           | 返回有关模块最新已知版本的 JSON 格式元数据，格式与 \$base/\$module/@v/\$version.info 相同。如果 \$base/\$module/@v/list 为空或没有列出合适的版本，则最新版本应该是 go 命令应使用的模块版本。该端点是可选的，并且不需要模块代理来实现它。                                                                            |

当解析模块的最新版本时，go命令将请求\$base/\$module/@v/list，如果没有找到合适的版本，则请求\$base/\$module/@latest。 go 命令按顺序优先选择：语义上最高的发布版本、语义上最高的预发布版本和按时间顺序排列的最新伪版本。在 Go 1.12 及更早版本中，go 命令将 \$base/\$module/@v/list 中的伪版本视为预发布版本，但自 Go 1.13 以来，情况不再如此。

模块代理必须始终为 \$base/\$module/\$version.mod 和 \$base/\$module/\$version.zip 查询提供成功响应的相同内容。此内容使用 go.sum 文件以及默认情况下的校验和数据库进行加密验证。

go 命令将从模块代理下载的大部分内容缓存在\$GOPATH/pkg/mod/cache/download 的模块缓存中。即使直接从版本控制系统下载，go 命令也会合成显式信息、mod 和 zip 文件并将它们存储在该目录中，就像直接从代理下载它们一样。缓存布局与代理 URL 空间相同，因此在（或将其复制到）https://example.com/proxy 提供 \$GOPATH/pkg/mod/cache/download 将允许用户通过设置 GOPROXY 访问缓存的模块版本到 https://example.com/proxy。

### Communicating with proxies 与proxy交互

go命令可以从模块代理下载模块源代码和元数据。 GOPROXY 环境变量可用于配置 go 命令可以连接到哪些代理以及是否可以直接与版本控制系统([version control systems](https://go.dev/ref/mod#vcs))通信。下载的模块数据保存在模块缓存中。 go 命令仅在需要缓存中尚未存在的信息时才会联系代理。

GOPROXY 协议部分描述了可以发送到 GOPROXY 服务器的请求。然而，了解 go 命令何时发出这些请求也很有帮助。例如，go build 遵循以下过程：

- 通过读取 go.mod 文件并执行最小版本选择 (MVS) 来计算构建列表。

- 读取命令行上指定的包以及它们导入的包。

- 如果构建列表中的任何模块都没有提供某个包，请查找提供该包的模块。在 go.mod 的最新版本中添加模块要求，然后重新开始。

- 在加载所有内容后构建包。

当 go 命令计算构建列表时，它会加载模块图中每个模块的 go.mod 文件。如果 go.mod 文件不在缓存中，go 命令将使用 \$module/@v/\$version.mod 请求从代理下载它（其中 \$module 是模块路径，\$version 是版本）。这些请求可以使用像curl这样的工具进行测试。例如，以下命令下载 golang.org/x/mod v0.2.0 版本的 go.mod 文件：

```shell
$ curl https://proxy.golang.org/golang.org/x/mod/@v/v0.2.0.mod
module golang.org/x/mod

go 1.12

require (
    golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
    golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e
    golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898
)
```

为了加载包，go 命令需要提供它的模块的源代码。模块源代码分布在 .zip 文件中，这些文件被提取到模块缓存中。如果模块 .zip 不在缓存中，go 命令将使用 \$module/@v/\$version.zip 请求下载它。

```shell
$ curl -O https://proxy.golang.org/golang.org/x/mod/@v/v0.2.0.zip
$ unzip -l v0.2.0.zip | head
Archive:  v0.2.0.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
     1479  00-00-1980 00:00   golang.org/x/mod@v0.2.0/LICENSE
     1303  00-00-1980 00:00   golang.org/x/mod@v0.2.0/PATENTS
      559  00-00-1980 00:00   golang.org/x/mod@v0.2.0/README
       21  00-00-1980 00:00   golang.org/x/mod@v0.2.0/codereview.cfg
      214  00-00-1980 00:00   golang.org/x/mod@v0.2.0/go.mod
     1476  00-00-1980 00:00   golang.org/x/mod@v0.2.0/go.sum
     5224  00-00-1980 00:00   golang.org/x/mod@v0.2.0/gosumcheck/main.go
```

请注意，.mod 和 .zip 请求是分开的，尽管 go.mod 文件通常包含在 .zip 文件中。 go命令可能需要下载许多不同模块的go.mod文件，而.mod文件比.zip文件小得多。此外，如果 Go 项目没有 go.mod 文件，代理将提供一个仅包含模块指令的合成 go.mod 文件。从版本控制系统下载时，go 命令会生成合成的 go.mod 文件。

如果 go 命令需要加载构建列表中任何模块都没有提供的包，它将尝试查找提供它的新模块。将包解析为模块部分描述了此过程。总之，go 命令请求有关可能包含该包的每个模块路径的最新版本的信息。例如，对于包 golang.org/x/net/html，go 命令将尝试查找模块 golang.org/x/net/html、golang.org/x/net、golang.org 的最新版本/x/ 和 golang.org。只有 golang.org/x/net 实际上存在并提供该包，因此 go 命令使用该模块的最新版本。如果多个模块提供该包，go 命令将使用路径最长的模块。

当 go 命令请求模块的最新版本时，它首先发送对 \$module/@v/list 的请求。如果列表为空或者返回的版本均不可使用，则会发送 \$module/@latest 请求。选择版本后，go 命令会发送 \$module/@v/\$version.info 元数据请求。然后，它可能会发送 \$module/@v/\$version.mod 和 \$module/@v/\$version.zip 请求来加载 go.mod 文件和源代码。

```shell
$ curl https://proxy.golang.org/golang.org/x/mod/@v/list
v0.1.0
v0.2.0

$ curl https://proxy.golang.org/golang.org/x/mod/@v/v0.2.0.info
{"Version":"v0.2.0","Time":"2020-01-02T17:33:45Z"}
```

下载 .mod 或 .zip 文件后，go 命令计算加密哈希并检查它是否与主模块的 go.sum 文件中的哈希匹配。如果 go.sum 中不存在哈希值，则默认情况下，go 命令会从校验和数据库中检索它。如果计算出的哈希值不匹配，go 命令会报告安全错误，并且不会在模块缓存中安装该文件。 GOPRIVATE 和 GONOSUMDB 环境变量可用于禁用对特定模块的校验和数据库的请求。 GOSUMDB 环境变量也可以设置为 off 以完全禁用对校验和数据库的请求。有关详细信息，请参阅验证模块。请注意，为 .info 请求返回的版本列表和版本元数据未经身份验证，并且可能会随着时间的推移而发生变化。

### Serving modules directly from a proxy 直接从代理提供模块服务

大多数模块都是从版本控制存储库开发和提供的。在直接模式下，go 命令使用版本控制工具下载此类模块（请参阅版本控制系统）。也可以直接从模块代理提供模块。这对于希望在不公开其版本控制服务器的情况下提供模块服务的组织以及使用 go 命令不支持的版本控制工具的组织非常有用。

当 go 命令以直接模式下载模块时，它首先根据模块路径使用 HTTP GET 请求查找模块服务器的 URL。它在 HTML 响应中查找名为 go-import 的 <meta> 标记。标签的内容必须包含存储库根路径、版本控制系统和 URL，并以空格分隔。有关详细信息，请参阅查找模块路径的存储库。

如果版本控制系统是 mod，go 命令会使用 GOPROXY 协议从给定的 URL 下载模块。

例如，假设 go 命令正在尝试下载版本 v1.0.0 的模块 example.com/gopher。它向 https://example.com/gopher?go-get=1 发送请求。服务器响应一个包含标签的 HTML 文档：

```
<meta name="go-import" content="example.com/gopher mod https://modproxy.example.com">
```

根据此响应，go 命令通过发送 https://modproxy.example.com/example.com/gopher/@v/v1.0.0.info、v1.0.0.mod 和 v1.0.0 的请求来下载模块。压缩。

请注意，直接从代理提供的模块无法在 GOPATH 模式下使用 go get 下载。

## Version control systems 版本控制系统

go命令可以直接从版本控制存储库下载模块源代码和元数据。从代理下载模块通常更快，但如果代理不可用或代理无法访问模块的存储库（对于私有存储库通常是这样），则需要直接连接到存储库。支持 Git、Subversion、Mercurial、Bazaar 和 Fossil。版本控制工具必须安装在 PATH 中的目录中才能让 go 命令使用它。

要从源存储库而不是代理下载特定模块，请设置 GOPRIVATE 或 GONOPROXY 环境变量。要将 go 命令配置为直接从源存储库下载所有模块，请将 GOPROXY 设置为 direct。有关详细信息，请参阅环境变量[Environment variables](https://go.dev/ref/mod#environment-variables)。

### Finding a repository for a module path 查找模块路径的存储库

当 go 命令以直接模式下载模块时，它首先查找包含该模块的存储库。

如果模块路径的路径组件末尾有 VCS 限定符（.bzr、.fossil、.git、.hg、.svn 之一），则 go 命令将使用该路径限定符之前的所有内容作为存储库 URL。例如，对于模块 example.com/foo.git/bar，go 命令使用 git 下载位于 example.com/foo.git 的存储库，期望在 bar 子目录中找到该模块。 go命令会根据版本控制工具支持的协议猜测要使用的协议。

如果模块路径没有限定符，go 命令将使用 ?go-get=1 查询字符串向从模块路径派生的 URL 发送 HTTP GET 请求。例如，对于模块 golang.org/x/mod，go 命令可能会发送以下请求：

```shell
https://golang.org/x/mod?go-get=1 (preferred)
http://golang.org/x/mod?go-get=1  (fallback, only with GOINSECURE)
```

go 命令遵循重定向，但会忽略响应状态代码，因此服务器可能会响应 404 或任何其他错误状态。 GOINSECURE 环境变量可以设置为允许回退并重定向到特定模块的未加密 HTTP。

服务器必须响应 HTML 文档，该文档的 <head> 中包含 <meta> 标记。 <meta> 标签应该出现在文档的早期，以避免混淆 go 命令的受限解析器。特别是，它应该出现在任何原始 JavaScript 或 CSS 之前。 <meta> 标记必须采用以下形式：

```html
<meta name="go-import" content="root-path vcs repo-url">
```

root-path 是存储库根路径，即与存储库根目录对应的模块路径部分。它必须是请求的模块路径的前缀或完全匹配。如果不完全匹配，则会对前缀发出另一个请求以验证 <meta> 标记是否匹配。

vcs是版本控制系统。它必须是下表列出的工具之一或关键字 mod，它指示 go 命令使用 GOPROXY 协议从给定 URL 下载模块。有关详细信息，请参阅直接从代理提供服务模块。

repo-url 是存储库的 URL。如果 URL 不包含方案（因为模块路径具有 VCS 限定符或因为 <meta> 标记缺少方案），则 go 命令将尝试版本控制系统支持的每个协议。例如，对于 Git，go 命令将尝试 https://，然后尝试 git+ssh://。仅当模块路径与 GOINSECURE 环境变量匹配时，才可以使用不安全协议（例如 http:// 和 git://）。

| Name       | Command | GoVCS default      | Secure schemes      |
| ---------- | ------- | ------------------ | ------------------- |
| Bazaar     | bzr     | Private only       | https, bzr+ssh      |
| Fossil     | fossil  | Private only       | https               |
| Git        | git     | Public and private | https, git+ssh, ssh |
| Mercurial  | hg      | Public and private | https, ssh          |
| Subversion | svn     | Private onl        | https, svn+ssh      |

作为一个例子，再次考虑 golang.org/x/mod。 go 命令向 https://golang.org/x/mod?go-get=1 发送请求。服务器响应一个包含标签的 HTML 文档：

```html
<meta name="go-import" content="golang.org/x/mod git https://go.googlesource.com/mod">
```

根据此响应，go 命令将使用远程 URL https://go.googlesource.com/mod 处的 Git 存储库。

GitHub 和其他流行的托管服务响应所有存储库的 ?go-get=1 查询，因此通常不需要为这些站点托管的模块进行服务器配置。

找到存储库 URL 后，go 命令会将存储库克隆到模块缓存中。一般来说，go 命令会尝试避免从存储库中获取不需要的数据。但是，实际使用的命令因版本控制系统而异，并且可能随时间而变化。对于 Git，go 命令可以列出大多数可用版本，而无需下载提交。它通常会获取提交而不下载祖先提交，但有时这样做是必要的。

### Mapping versions to commits 将版本映射到commit

go 命令可以检查存储库中特定规范版本（例如 v1.2.3、v2.4.0-beta 或 v3.0.0+不兼容）的模块。每个模块版本在存储库中都应该有一个语义版本标签，指示应该为给定版本签出哪个修订版本。

如果模块定义在存储库根目录或根目录的主版本子目录中，则每个版本标记名称等于相应的版本。例如，模块 golang.org/x/text 定义在其存储库的根目录中，因此版本 v0.3.2 在该存储库中具有标签 v0.3.2。对于大多数模块来说都是如此。

如果模块定义在存储库内的子目录中，即模块路径的模块子目录部分不为空，则每个标记名称必须以模块子目录为前缀，后跟斜杠。例如，模块 golang.org/x/tools/gopls 定义在存储库的 gopls 子目录中，根路径为 golang.org/x/tools。该模块的 v0.4.0 版本必须在该存储库中具有名为 gopls/v0.4.0 的标签。

语义版本标签的主版本号必须与模块路径的主版本后缀（如果有）一致。例如，标签 v1.0.0 可能属于模块 example.com/mod，但不属于 example.com/mod/v2，后者将具有类似 v2.0.0 的标签。

如果不存在 go.mod 文件，并且该模块位于存储库根目录中，则主版本 v2 或更高版本的标签可能属于没有主版本后缀的模块。这种版本用后缀+不兼容表示。版本标签本身不能有后缀。请参阅与非模块存储库的兼容性。 [Compatibility with non-module repositories](https://go.dev/ref/mod#non-module-compat).

一旦创建标签，就不应将其删除或更改为其他版本。版本经过验证以确保安全、可重复的构建。如果标签被修改，客户端在下载时可能会看到安全错误。即使删除标签后，其内容也可能在模块代理上仍然可用。

### Mapping pseudo-versions to commits

go 命令可以检查存储库中特定修订版的模块，编码为伪版本，如 v1.3.2-0.20191109021931-daa7c04131f5。

伪版本的最后 12 个字符（上例中的 daa7c04131f5）表示存储库中要签出的修订版本。这的含义取决于版本控制系统。对于 Git 和 Mercurial，这是提交哈希的前缀。对于 Subversion，这是一个用零填充的修订号。

在签出提交之前，go 命令会验证时间戳（上面的 20191109021931）是否与提交日期匹配。它还验证基本版本（v1.3.1，上例中 v1.3.2 之前的版本）是否对应于作为提交祖先的语义版本标记。这些检查确保模块作者可以完全控制伪版本与其他发布版本的比较。

有关更多信息，请参阅伪版本。 [Pseudo-versions](https://go.dev/ref/mod#pseudo-versions)

### Mapping branches and commits to versions

可以使用版本查询在特定分支、标签或修订版本处签出模块。

```shell
go get example.com/mod@master
```

go 命令将这些名称转换为可以通过最小版本选择（MVS）使用的规范版本。 MVS 取决于明确订购版本的能力。随着时间的推移，分支名称和修订版本无法可靠地进行比较，因为它们依赖于可能会发生变化的存储库结构。

如果修订版标有一个或多个语义版本标签（例如 v1.2.3），则将使用最高有效版本的标签。 go命令只考虑可能属于目标模块的语义版本标签；例如，example.com/mod/v2 不会考虑标签 v1.5.2，因为主要版本与模块路径的后缀不匹配。

如果修订版本未标记有效的语义版本标签，则 go 命令将生成伪版本。如果修订版具有带有有效语义版本标签的祖先，则最高祖先版本将用作伪版本基础。请参阅伪版本。

### Module directories within a repository

一旦在特定版本中签出了模块的存储库，go 命令必须找到包含模块的 go.mod 文件的目录（模块的根目录）。
回想一下，模块路径由三部分组成：存储库根路径（对应于存储库根目录）、模块子目录和主版本后缀（仅适用于 v2 或更高版本发布的模块）。
对于大多数模块，模块路径等于存储库根路径，因此模块的根目录就是存储库的根目录。
模块有时在存储库子目录中定义。这通常是针对具有多个需要独立发布和版本控制的组件的大型存储库。这样的模块应该在与存储库根路径之后的模块路径部分匹配的子目录中找到。例如，假设模块 example.com/monorepo/foo/bar 位于根路径为 example.com/monorepo 的存储库中。它的 go.mod 文件必须位于 foo/bar 子目录中。
如果模块以主版本 v2 或更高版本发布，则其路径必须具有主版本后缀。带有主版本后缀的模块可以在两个子目录之一中定义：一个带有后缀，一个不带有后缀。例如，假设上述模块的新版本已发布，路径为 example.com/monorepo/foo/bar/v2。它的 go.mod 文件可能位于 foo/bar 或 foo/bar/v2 中。
带有主版本后缀的子目录是主版本子目录。它们可用于在单个分支上开发模块的多个主要版本。当多个主要版本的开发在不同的分支上进行时，这可能是不必要的。然而，大版本子目录有一个重要的属性：在GOPATH模式下，包导入路径与GOPATH/src下的目录完全匹配。 go 命令在 GOPATH 模式下提供最小的模块兼容性（请参阅与非模块存储库的兼容性），因此主要版本子目录对于与 GOPATH 模式下构建的项目的兼容性并不总是必需的。不过，不支持最低模块兼容性的旧工具可能会存在问题。
一旦 go 命令找到模块根目录，它就会创建该目录内容的 .zip 文件，然后将该 .zip 文件提取到模块缓存中。有关 .zip 文件中可能包含哪些文件的详细信息，请参阅文件路径和大小限制。 .zip 文件的内容在提取到模块缓存之前会进行身份验证，就像从代理下载 .zip 文件一样。
模块 zip 文件不包含供应商目录或任何嵌套模块（包含 go.mod 文件的子目录）的内容。这意味着模块必须注意不要引用其目录之外或其他模块中的文件。例如， //go:embed 模式不得与嵌套模块中的文件匹配。在文件不应包含在模块中的情况下，此行为可能是一种有用的解决方法。例如，如果存储库将大文件签入 testdata 目录，模块作者可以在 testdata 中添加一个空的 go.mod 文件，这样他们的用户就不需要下载这些文件。当然，这可能会减少用户测试其依赖项的覆盖范围。

### Special case for LICENSE files

当 go 命令为不在存储库根目录中的模块创建 .zip 文件时，如果该模块的根目录中没有名为 LICENSE 的文件（与 go.mod 一起），则 go 命令将复制名为 LICENSE 的文件如果同一版本中存在许可证，则来自存储库根目录的许可证。
这种特殊情况允许相同的许可证文件应用于存储库中的所有模块。这仅适用于专门名为 LICENSE 的文件，不带 .txt 等扩展名。不幸的是，如果不破坏现有模块的加密总和，就无法扩展它；请参阅验证模块。其他工具和网站（例如 pkg.go.dev）可能会识别具有其他名称的文件。
另请注意，go 命令在创建模块 .zip 文件时不包含符号链接；请参阅文件路径和大小限制。因此，如果存储库的根目录中没有 LICENSE 文件，作者可以在子目录中定义的模块中创建其许可证文件的副本，以确保这些文件包含在模块 .zip 文件中。

### Controlling version control tools with GOVCS 使用 GOVCS 控制版本控制工具

go 命令能够使用 git 等版本控制命令下载模块，这对于去中心化的包生态系统至关重要，在该生态系统中可以从任何服务器导入代码。如果恶意服务器找到一种方法导致调用的版本控制命令运行非预期代码，这也是一个潜在的安全问题。

为了平衡功能和安全问题，go命令默认只使用git和hg从公共服务器下载代码。它将使用任何已知的版本控制系统从私有服务器下载代码，私有服务器定义为与 GOPRIVATE 环境变量匹配的托管包。仅允许 Git 和 Mercurial 的理由是，这两个系统最关注作为不受信任服务器的客户端运行的问题。相比之下，Bazaar、Fossil 和 Subversion 主要用于受信任、经过身份验证的环境，并且没有像攻击面那样受到严格审查。

版本控制命令限制仅在使用直接版本控制访问下载代码时适用。从代理下载模块时，go 命令使用 GOPROXY 协议，该协议始终是允许的。默认情况下，go 命令对公共模块使用 Go 模块镜像 (proxy.golang.org)，并且仅在私有模块或镜像拒绝提供公共包（通常出于法律原因）时才回退到版本控制。因此，默认情况下，客户端仍然可以访问 Bazaar、Fossil 或 Subversion 存储库提供的公共代码，因为这些下载使用 Go 模块镜像，这会带来使用自定义沙箱运行版本控制命令的安全风险。

版本控制命令限制仅在使用直接版本控制访问下载代码时适用。从代理下载模块时，go 命令使用 GOPROXY 协议，该协议始终是允许的。默认情况下，go 命令对公共模块使用 Go 模块镜像 (proxy.golang.org)，并且仅在私有模块或镜像拒绝提供公共包（通常出于法律原因）时才回退到版本控制。因此，默认情况下，客户端仍然可以访问 Bazaar、Fossil 或 Subversion 存储库提供的公共代码，因为这些下载使用 Go 模块镜像，这会带来使用自定义沙箱运行版本控制命令的安全风险。

例如，考虑：

```shell
GOVCS=github.com:git,evil.com:off,*:git|hg
```

通过此设置，模块或导入路径以 github.com/ 开头的代码只能使用 git； evil.com 上的路径不能使用任何版本控制命令，所有其他路径（* 匹配所有内容）只能使用 git 或 hg。

特殊模式 public 和 private 与公共和私有模块或导入路径相匹配。如果路径与 GOPRIVATE 变量匹配，则该路径是私有的；否则它是公开的。

如果 GOVCS 变量中没有规则与特定模块或导入路径匹配，则 go 命令将应用其默认规则，该规则现在可以用 GOVCS 表示法概括为 public:git|hg,private:all。

要允许对任何包不受限制地使用任何版本控制系统，请使用：

```shell
GOVCS=*:all
```

要禁用版本控制的所有使用，请使用：

```shell
GOVCS=*:off
```

go env -w 命令可用于设置 GOVCS 变量以供将来的 go 命令调用。
GOVCS 在 Go 1.16 中引入。 Go 的早期版本可以对任何模块使用任何已知的版本控制工具。

## Module zip files

模块版本以 .zip 文件形式分发。很少需要直接与这些文件交互，因为 go 命令会自动从模块代理和版本控制存储库中创建、下载和提取它们。但是，了解这些文件对于了解跨平台兼容性约束或实现模块代理仍然很有用。

go mod download 命令下载一个或多个模块的 zip 文件，然后将这些文件提取到模块缓存中。根据 GOPROXY 和其他环境变量，go 命令可以从代理下载 zip 文件，也可以克隆源代码控制存储库并从中创建 zip 文件。 -json 标志可用于查找下载 zip 文件的位置及其在模块缓存中提取的内容。

golang.org/x/mod/zip 包可用于以编程方式创建、提取或检查 zip 文件的内容。

### File path and size constraints

模块 zip 文件的内容有许多限制。这些限制确保可以在各种平台上安全、一致地提取 zip 文件。

- 模块 zip 文件的大小最多为 500 MiB。其文件的未压缩总大小也限制为 500 MiB。 go.mod 文件限制为 16 MiB。许可证文件也限制为 16 MiB。这些限制的存在是为了减轻对用户、代理和模块生态系统其他部分的拒绝服务攻击。模块目录树中包含超过 500 MiB 文件的存储库应在提交时标记模块版本，仅包含构建模块包所需的文件；构建通常不需要视频、模型和其他大型资源。

- 模块 zip 文件中的每个文件都必须以前缀 \$module@\$version/ 开头，其中 $module 是模块路径，$version 是版本，例如 golang.org/x/mod@v0.3.0/。模块路径必须有效，版本必须有效且规范，并且版本必须与模块路径的主版本后缀匹配。有关特定定义和限制，请参阅模块路径和版本。

- 文件模式、时间戳和其他元数据将被忽略。

- 空目录（路径以斜线结尾的条目）可能包含在模块 zip 文件中，但不会被提取。 go 命令在它创建的 zip 文件中不包含空目录。

- 创建 zip 文件时，符号链接和其他不规则文件将被忽略，因为它们不可跨操作系统和文件系统移植，并且没有可移植的方式以 zip 文件格式表示它们。

- 创建 zip 文件时，名为供应商的目录中的文件将被忽略，因为从不使用主模块之外的供应商目录。

- 创建 zip 文件时，包含 go.mod 文件的目录（模块根目录除外）中的文件将被忽略，因为它们不是模块的一部分。 go 命令在提取 zip 文件时会忽略包含 go.mod 文件的子目录。

- 在 Unicode 大小写折叠下，zip 文件中的两个文件不能具有相同的路径（请参阅 strings.EqualFold）。这确保了可以在不区分大小写的文件系统上提取 zip 文件而不会发生冲突。

- go.mod 文件可能会也可能不会出现在顶级目录 (\$module@\$version/go.mod) 中。如果存在，它必须具有名称 go.mod （全部小写）。任何其他目录中都不允许名为 go.mod 的文件。

- 模块内的文件和目录名称可以由 Unicode 字母、ASCII 数字、ASCII 空格字符 (U+0020) 和 ASCII 标点字符 !#$%&()+,-.=@[]^_{} 组成～。请注意，包路径可能不包含所有这些字符。请参阅 module.CheckFilePath 和 module.CheckImportPath 了解差异。

- 第一个点之前的文件或目录名不得是 Windows 上的保留文件名，无论大小写（CON、com1、NuL 等）。

## Private modules 私有模块

Go 模块经常在版本控制服务器和模块代理上开发和分发，而这些服务器和模块代理在公共互联网上不可用。 go 命令可以从私有源下载和构建模块，尽管它通常需要一些配置。

下面的环境变量可用于配置对私有模块的访问。有关详细信息，请参阅环境变量。另请参阅隐私以了解有关控制发送到公共服务器的信息的信息。

- GOPROXY - 模块代理 URL 列表。 go 命令将尝试按顺序从每个服务器下载模块。关键字 direct 指示 go 命令从开发模块的版本控制存储库下载模块，而不是使用代理。

- GOPRIVATE - 应被视为私有的模块路径前缀的全局模式列表。充当 GONOPROXY 和 GONOSUMDB 的默认值。

- GONOPROXY - 不应从代理下载的模块路径前缀的全局模式列表。 go 命令将从开发模块的版本控制存储库下载匹配的模块，无论 GOPROXY 是什么。

- GONOSUMDB - 不应使用公共校验和数据库 sum.golang.org 检查的模块路径前缀的全局模式列表。

- GOINSECURE - 可以通过 HTTP 和其他不安全协议检索的模块路径前缀的全局模式列表。

这些变量可以在开发环境中设置（例如，在 .profile 文件中），也可以使用 go env -w 永久设置。

本节的其余部分描述了提供对私有模块代理和版本控制存储库的访问的常见模式。

### Private proxy serving all modules

为所有模块（公共和私有）提供服务的中央私有代理服务器为管理员提供了最大的控制权，并且对单个开发人员来说需要最少的配置。

要将 go 命令配置为使用此类服务​​器，请设置以下环境变量，将 https://proxy.corp.example.com 替换为您的代理 URL，将 corp.example.com 替换为您的模块前缀：

```shell
GOPROXY=https://proxy.corp.example.com
GONOSUMDB=corp.example.com
```

GOPROXY 设置指示 go 命令仅从 https://proxy.corp.example.com 下载模块； go 命令不会连接到其他代理或版本控制存储库。

GONOSUMDB 设置指示 go 命令不要使用公共校验和数据库来验证路径以 corp.example.com 开头的模块。

在此配置中运行的代理可能需要对私有版本控制服务器的读取访问权限。它还需要访问公共互联网来下载公共模块的新版本。

有几种现有的 GOPROXY 服务器实现可以通过这种方式使用。最小的实现将从模块缓存目录中提供文件，并使用 go mod download （具有适当的配置）来检索丢失的模块。

### Direct access to private modules

go 命令可以配置为绕过公共代理并直接从版本控制服务器下载私有模块。当运行私人代理服务器不可行时，这非常有用。

要将 go 命令配置为以这种方式工作，请设置 GOPRIVATE，替换 corp.example.com 私有模块前缀：

```shell
GOPRIVATE=corp.example.com
```

在这种情况下不需要更改 GOPROXY 变​​量。它默认为 https://proxy.golang.org,direct，它指示 go 命令首先尝试从 https://proxy.golang.org 下载模块，然后如果该代理响应 404，则回退到直接连接（未找到）或 410（已消失）。

GOPRIVATE 设置指示 go 命令不要连接到以 corp.example.com 开头的模块的代理或校验和数据库。

可能仍然需要内部 HTTP 服务器来解析存储库 URL 的模块路径。例如，当 go 命令下载模块 corp.example.com/mod 时，它将向 https://corp.example.com/mod?go-get=1 发送 GET 请求，并查找存储库响应中的 URL。为了避免这种要求，请确保每个私有模块路径都有一个 VCS 后缀（如 .git），标记存储库根前缀。例如，当 go 命令下载模块 corp.example.com/repo.git/mod 时，它将克隆 Git 存储库 https://corp.example.com/repo.git 或 ssh://corp.example .com/repo.git，无需提出额外请求。

开发人员需要对包含私有模块的存储库进行读取访问。这可以在全局 VCS 配置文件（如 .gitconfig）中进行配置。最好将 VCS 工具配置为不需要交互式身份验证提示。默认情况下，调用 Git 时，go 命令通过设置 GIT_TERMINAL_PROMPT=0 来禁用交互式提示，但它遵循显式设置。

### Passing credentials to private proxies

go 命令在与代理服务器通信时支持 HTTP 基本身份验证。

凭证可以在 .netrc 文件中指定。例如，包含以下行的 .netrc 文件将配置 go 命令以使用给定的用户名和密码连接到计算机 proxy.corp.example.com。

```shell
machine proxy.corp.example.com
login jrgopher
password hunter2
```

文件的位置可以使用 NETRC 环境变量设置。如果未设置 NETRC，go 命令将在类 UNIX 平台上读取 $HOME/.netrc 或在 Windows 上读取 %USERPROFILE%\_netrc。

.netrc 中的字段用空格、制表符和换行符分隔。不幸的是，这些字符不能在用户名或密码中使用。另请注意，计算机名称不能是完整的 URL，因此不可能为同一计算机上的不同路径指定不同的用户名和密码。

或者，可以直接在 GOPROXY URL 中指定凭证。例如：

```shell
GOPROXY=https://jrgopher:hunter2@proxy.corp.example.com
```

采用此方法时请务必小心：环境变量可能会出现在 shell 历史记录和日志中。

### Passing credentials to private repositories

go 命令可以直接从版本控制存储库下载模块。如果不使用私有代理，这对于私有模块是必需的。请参阅直接访问私有模块进行配置。

go 命令在直接下载模块时运行 git 等版本控制工具。这些工具执行自己的身份验证，因此您可能需要在特定于工具的配置文件（如 .gitconfig）中配置凭据。

为了确保此操作顺利进行，请确保 go 命令使用正确的存储库 URL，并且版本控制工具不需要交互式输入密码。 go 命令更喜欢 https:// URL 而不是其他方案（如 ssh://），除非在查找存储库 URL 时指定了该方案。特别是对于 GitHub 存储库，go 命令假定为 https://。

对于大多数服务器，您可以将客户端配置为通过 HTTP 进行身份验证。例如，GitHub 支持使用 OAuth 个人访问令牌作为 HTTP 密码。您可以将 HTTP 密码存储在 .netrc 文件中，就像将凭据传递给私有代理时一样。

或者，您可以将 https:// URL 重写为另一个方案。例如，在 .gitconfig 中：

```shell
[url "git@github.com:"]
    insteadOf = https://github.com/
```

### Privacy   隐私

go 命令可以从模块代理服务器和版本控制系统下载模块和元数据。环境变量 GOPROXY 控制使用哪些服务器。环境变量 GOPRIVATE 和 GONOPROXY 控制从代理获取哪些模块。

GOPROXY 的默认值为：

```shell
https://proxy.golang.org,direct
```

通过此设置，当 go 命令下载模块或模块元数据时，它会首先向 proxy.golang.org 发送请求，这是 Google 运营的公共模块代理（隐私政策）。有关每个请求中发送哪些信息的详细信息，请参阅 GOPROXY 协议。 go 命令不传输个人身份信息，但它确实传输所请求的完整模块路径。如果代理响应 404（未找到）或 410（已消失）状态，则 go 命令将尝试直接连接到提供该模块的版本控制系统。有关详细信息，请参阅版本控制系统。

GOPRIVATE 或 GONOPROXY 环境变量可以设置为与私有模块前缀匹配的 glob 模式列表，不应从任何代理请求。例如：

```shell
GOPRIVATE=*.corp.example.com,*.research.example.com
```

GOPRIVATE 只是充当 GONOPROXY 和 GONOSUMDB 的默认值，因此无需设置 GONOPROXY，除非 GONOSUMDB 应具有不同的值。当模块路径与 GONOPROXY 匹配时，go 命令会忽略该模块的 GOPROXY 并直接从其版本控制存储库中获取它。当没有代理服务私有模块时，这非常有用。请参阅直接访问私有模块。

如果有一个可信代理为所有模块提供服务，则不应设置 GONOPROXY。例如，如果 GOPROXY 设置为一个源，则 go 命令将不会从其他源下载模块。在这种情况下仍应设置 GONOSUMDB。

```shell
GOPROXY=https://proxy.corp.example.com
GONOSUMDB=*.corp.example.com,*.research.example.com
```

如果存在仅服务私有模块的受信任代理，则不应设置 GONOPROXY，但必须注意确保代理以正确的状态代码进行响应。例如，考虑以下配置：

```shell
GOPROXY=https://proxy.corp.example.com,https://proxy.golang.org
GONOSUMDB=*.corp.example.com,*.research.example.com
```

假设由于拼写错误，开发人员尝试下载不存在的模块。

```shell
go mod download corp.example.com/secret-product/typo@latest
```

go 命令首先从 proxy.corp.example.com 请求此模块。如果该代理响应 404（未找到）或 410（已消失），go 命令将回退到 proxy.golang.org，在请求 URL 中传输秘密产品路径。如果私有代理以任何其他错误代码响应，go 命令将打印错误并且不会回退到其他源。

除了代理之外，go 命令还可以连接到校验和数据库以验证 go.sum 中未列出的模块的加密哈希值。 GOSUMDB 环境变量设置校验和数据库的名称、URL 和公钥。 GOSUMDB 的默认值是 sum.golang.org，由 Google 运营的公共校验和数据库（隐私政策）。有关每个请求传输内容的详细信息，请参阅校验和数据库。与代理一样，go 命令不会传输个人身份信息，但它会传输所请求的完整模块路径，并且校验和数据库无法计算非公共模块的校验和。

GONOSUMDB 环境变量可以设置为指示哪些模块是私有的并且不应从校验和数据库请求的模式。 GOPRIVATE 作为 GONOSUMDB 和 GONOPROXY 的默认值，因此没有必要设置 GONOSUMDB，除非 GONOPROXY 应具有不同的值。

代理可以镜像校验和数据库。如果 GOPROXY 中的代理执行此操作，则 go 命令将不会直接连接到校验和数据库。

GOSUMDB 可以设置为 off 以完全禁用校验和数据库。通过此设置，go 命令将不会验证下载的模块，除非它们已经在 go.sum 中。请参阅验证模块。

## Module cache

模块缓存是 go 命令存储下载的模块文件的目录。模块缓存与构建缓存不同，构建缓存包含已编译的包和其他构建工件。

模块缓存的默认位置是 \$GOPATH/pkg/mod。要使用不同的位置，请设置 GOMODCACHE 环境变量。

模块缓存没有最大大小，并且 go 命令不会自动删除其内容。

缓存可能被同一台机器上开发的多个Go项目共享。无论主模块的位置如何，go 命令都将使用相同的缓存。 go 命令的多个实例可以同时安全地访问同一模块缓存。

go 命令在缓存中创建具有只读权限的模块源文件和目录，以防止下载模块后对其进行意外更改。这有一个不幸的副作用，即使缓存难以使用 rm -rf 等命令删除。可以使用 go clean -modcache 删除缓存。或者，当使用 -modcacherw 标志时，go 命令将创建具有读写权限的新目录。这增加了编辑器、测试和其他程序修改模块缓存中的文件的风险。 go mod verify 命令可用于检测对主模块的依赖项的修改。它扫描每个模块依赖项的提取内容并确认它们与 go.sum 中的预期哈希匹配。

下表解释了模块缓存中大多数文件的用途。一些临时文件（锁定文件、临时目录）被省略。对于每个路径，$module 是模块路径，$version 是版本。以斜杠 (/) 结尾的路径是目录。模块路径和版本中的大写字母使用感叹号进行转义（Azure 转义为 !azure），以避免在不区分大小写的文件系统上发生冲突。

| Path                                         | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| -------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `$module@$version/`                          | Directory containing extracted contents of a module `.zip` file. This serves as a module root directory for a downloaded module. It won't contain a `go.mod` file if the original module didn't have one.                                                                                                                                                                                                                                                                               |
| `cache/download/`                            | Directory containing files downloaded from module proxies and files derived from [version control systems](https://go.dev/ref/mod#vcs). The layout of this directory follows the [`GOPROXY` protocol](https://go.dev/ref/mod#goproxy-protocol), so this directory may be used as a proxy when served by an HTTP file server or when referenced with a `file://` URL.                                                                                                                    |
| `cache/download/$module/@v/list`             | List of known versions (see [`GOPROXY` protocol](https://go.dev/ref/mod#goproxy-protocol)). This may change over time, so the `go` command usually fetches a new copy instead of re-using this file.                                                                                                                                                                                                                                                                                    |
| `cache/download/$module/@v/$version.info`    | JSON metadata about the version. (see [`GOPROXY` protocol](https://go.dev/ref/mod#goproxy-protocol)). This may change over time, so the `go` command usually fetches a new copy instead of re-using this file.                                                                                                                                                                                                                                                                          |
| `cache/download/$module/@v/$version.mod`     | The `go.mod` file for this version (see [`GOPROXY` protocol](https://go.dev/ref/mod#goproxy-protocol)). If the original module did not have a `go.mod` file, this is a synthesized file with no requirements.                                                                                                                                                                                                                                                                           |
| `cache/download/$module/@v/$version.zip`     | The zipped contents of the module (see [`GOPROXY` protocol](https://go.dev/ref/mod#goproxy-protocol) and [Module zip files](https://go.dev/ref/mod#zip-files)).                                                                                                                                                                                                                                                                                                                         |
| `cache/download/$module/@v/$version.ziphash` | A cryptographic hash of the files in the `.zip` file. Note that the `.zip` file itself is not hashed, so file order, compression, alignment, and metadata don't affect the hash. When using a module, the `go` command verifies this hash matches the corresponding line in [`go.sum`](https://go.dev/ref/go-sum-files). The [`go mod verify`](https://go.dev/ref/mod#go-mod-verify) command checks that the hashes of module `.zip` files and extracted directories match these files. |
| `cache/download/sumdb/`                      | Directory containing files downloaded from a [checksum database](https://go.dev/ref/mod#checksum-database) (typically `sum.golang.org`).                                                                                                                                                                                                                                                                                                                                                |
| `cache/vcs/`                                 | Contains cloned version control repositories for modules fetched directly from their sources. Directory names are hex-encoded hashes derived from the repository type and URL. Repositories are optimized for size on disk. For example, cloned Git repositories are bare and shallow when possible.                                                                                                                                                                                    |

## Authenticating modules

当 go 命令将模块 zip 文件或 go.mod 文件下载到模块缓存中时，它会计算加密哈希并将其与已知值进行比较，以验证文件自首次下载以来没有更改。如果下载的文件没有正确的哈希值，go 命令会报告安全错误。
对于 go.mod 文件，go 命令根据文件内容计算哈希值。对于模块 zip 文件，go 命令按确定的顺序根据存档中文件的名称和内容计算哈希值。哈希不受文件顺序、压缩、对齐和其他元数据的影响。有关哈希实现的详细信息，请参阅 golang.org/x/mod/sumdb/dirhash。
go 命令将每个哈希与主模块的 go.sum 文件中的相应行进行比较。如果该哈希值与 go.sum 中的哈希值不同，则 go 命令会报告安全错误并删除下载的文件，而不将其添加到模块缓存中。
如果 go.sum 文件不存在，或者它不包含下载文件的哈希值，则 go 命令可以使用校验和数据库（公共可用模块的哈希值的全局源）来验证哈希值。一旦验证了哈希值，go 命令会将其添加到 go.sum 并将下载的文件添加到模块缓存中。如果模块是私有的（由 GOPRIVATE 或 GONOSUMDB 环境变量匹配），或者如果禁用校验和数据库（通过设置 GOSUMDB=off），则 go 命令接受哈希并将文件添加到模块缓存而不验证它。
模块缓存通常由系统上的所有 Go 项目共享，每个模块可能有自己的 go.sum 文件，其哈希值可能不同。为了避免信任其他模块，go 命令每当访问模块缓存中的文件时都会使用主模块的 go.sum 来验证哈希值。 Zip 文件哈希值的计算成本很高，因此 go 命令会检查与 zip 文件一起存储的预先计算的哈希值，而不是重新哈希文件。 go mod verify 命令可用于检查 zip 文件和提取的目录自添加到模块缓存以来是否未被修改。

### go.sum files

模块的根目录中可能有一个名为 go.sum 的文本文件，以及它的 go.mod 文件。 go.sum 文件包含模块的直接和间接依赖项的加密哈希值。当 go 命令将模块 .mod 或 .zip 文件下载到模块缓存中时，它会计算哈希并检查该哈希是否与主模块的 go.sum 文件中的相应哈希匹配。如果模块没有依赖项或者使用替换指令将所有依赖项替换为本地目录，则 go.sum 可能为空或不存在。

go.sum 中的每一行都有三个由空格分隔的字段：模块路径、版本（可能以 /go.mod 结尾）和哈希值。

- 模块路径是哈希所属模块的名称。

- 版本是哈希所属模块的版本。如果版本以 /go.mod 结尾，则哈希值仅适用于模块的 go.mod 文件；否则，哈希值适用于模块的 .zip 文件中的文件。

- 哈希列由算法名称（如 h1）和 Base64 编码的加密哈希组成，并用冒号 (:) 分隔。目前，SHA-256 (h1) 是唯一支持的哈希算法。如果将来发现 SHA-256 中的漏洞，将添加对另一种算法（名为 h2 等）的支持。

go.sum 文件可能包含模块的多个版本的哈希值。 go 命令可能需要从依赖项的多个版本加载 go.mod 文件，以便执行最小版本选择。 go.sum 还可能包含不再需要的模块版本的哈希值（例如，升级后）。 go mod tidy 将添加缺失的哈希值，并从 go.sum 中删除不必要的哈希值。

### Checksum database

校验和数据库是 go.sum 行的全局源。 go 命令可以在许多情况下使用它来检测代理或源服务器的不当行为。
校验和数据库允许所有公开可用的模块版本的全局一致性和可靠性。它使得不受信任的代理成为可能，因为它们无法在不被注意到的情况下提供错误的代码。它还确保与特定版本相关的位不会日复一日地发生变化，即使模块的作者随后更改了其存储库中的标签。
校验和数据库由 Google 运营的 sum.golang.org 提供服务。它是 go.sum 行哈希的透明日志（或“默克尔树”），由 Trillian 支持。 Merkle 树的主要优点是独立审计员可以验证它没有被篡改，因此它比简单的数据库更值得信赖。
go 命令使用最初在提案：保护公共 Go 模块生态系统中概述的协议与校验和数据库进行交互。
下表指定了校验和数据库必须响应的查询。对于每个路径，$base 是校验和数据库 URL 的路径部分，$module 是模块路径，\$version 是版本。例如，如果校验和数据库 URL 是 https://sum.golang.org，并且客户端正在请求版本 v0.3.2 的模块 golang.org/x/text 的记录，则客户端将发送 GET 请求https://sum.golang.org/lookup/golang.org/x/text@v0.3.2。
为了避免在不区分大小写的文件系统中提供服务时出现歧义，通过将每个大写字母替换为感叹号，后跟相应的小写字母，对 $module 和 $version 元素进行大小写编码。这允许模块 example.com/M 和 example.com/m 都存储在磁盘上，因为前者被编码为 example.com/!m。
路径中用方括号括起来的部分（例如 [.p/\$W]）表示可选值。

| Path                            | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| ------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `$base/latest`                  | Returns a signed, encoded tree description for the latest log. This signed description is in the form of a [note](https://pkg.go.dev/golang.org/x/mod/sumdb/note), which is text that has been signed by one or more server keys and can be verified using the server's public key. The tree description provides the size of the tree and the hash of the tree head at that size. This encoding is described in `[golang.org/x/mod/sumdb/tlog#FormatTree](https://pkg.go.dev/golang.org/x/mod/sumdb/tlog#FormatTree)`. |
| `$base/lookup/$module@$version` | Returns the log record number for the entry about `$module` at `$version`, followed by the data for the record (that is, the `go.sum` lines for `$module` at `$version`) and a signed, encoded tree description that contains the record.                                                                                                                                                                                                                                                                               |
| `$base/tile/$H/$L/$K[.p/$W]`    | Returns a [log tile](https://research.swtch.com/tlog#serving_tiles), which is a set of hashes that make up a section of the log. Each tile is defined in a two-dimensional coordinate at tile level `$L`, `$K`th from the left, with a tile height of `$H`. The optional `.p/$W` suffix indicates a partial log tile with only `$W` hashes. Clients must fall back to fetching the full tile if a partial tile is not found.                                                                                            |
| `$base/tile/$H/data/$K[.p/$W]`  | Returns the record data for the leaf hashes in `/tile/$H/0/$K[.p/$W]` (with a literal `data` path element).                                                                                                                                                                                                                                                                                                                                                                                                             |

如果 go 命令查询校验和数据库，那么第一步是通过 /lookup 端点检索记录数据。如果日志中尚未记录模块版本，则校验和数据库将在回复之前尝试从源服务器获取它。此 /lookup 数据提供了此模块版本的总和及其在日志中的位置，这通知客户端应获取哪些图块来执行证明。 go 命令在向主模块的 go.sum 文件添加新的 go.sum 行之前执行“包含”证明（日志中存在特定记录）和“一致性”证明（树未被篡改）。重要的是，如果没有首先根据签名树哈希对其进行身份验证，并根据客户端的签名树哈希时间线对签名树哈希进行身份验证，则永远不要使用来自 /lookup 的数据。
签名的树哈希和校验和数据库提供的新切片存储在模块缓存中，因此 go 命令只需要获取丢失的切片。
go命令不需要直接连接到校验和数据库。它可以通过镜像校验和数据库并支持上述协议的模块代理请求模块和。这对于阻止组织外部请求的私人企业代理特别有帮助。
GOSUMDB 环境变量标识要使用的校验和数据库的名称以及可选的公钥和 URL，如下所示：

```shell
GOSUMDB="sum.golang.org"
GOSUMDB="sum.golang.org+<publickey>"
GOSUMDB="sum.golang.org+<publickey> https://sum.golang.org"
```

go 命令知道 sum.golang.org 的公钥，并且知道名称 sum.golang.google.cn （在中国大陆可用）连接到 sum.golang.org 校验和数据库；使用任何其他数据库都需要明确提供公钥。 URL 默认为 https:// 后跟数据库名称。
GOSUMDB 默认为 sum.golang.org，由 Google 运行的 Go 校验和数据库。请参阅 https://sum.golang.org/privacy 了解该服务的隐私政策。
如果 GOSUMDB 设置为 off，或者使用 -insecure 标志调用 go get，则不会查阅校验和数据库，并且接受所有无法识别的模块，但代价是放弃所有模块经过验证的可重复下载的安全保证。绕过特定模块的校验和数据库的更好方法是使用 GOPRIVATE 或 GONOSUMDB 环境变量。有关详细信息，请参阅私有模块。
go env -w 命令可用于设置这些变量以供将来的 go 命令调用。

## Environment variables

go 命令中的模块行为可以使用下面列出的环境变量进行配置。该列表仅包含与模块相关的环境变量。请参阅 go help 环境以获取 go 命令识别的所有环境变量的列表。

| Variable      | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `GO111MODULE` | Controls whether the `go` command runs in module-aware mode or `GOPATH` mode. Three values are recognized:<br><br>- `off`: the `go` command ignores `go.mod` files and runs in `GOPATH` mode.<br>- `on` (or unset): the `go` command runs in module-aware mode, even when no `go.mod` file is present.<br>- `auto`: the `go` command runs in module-aware mode if a `go.mod` file is present in the current directory or any parent directory. In Go 1.15 and lower, this was the default.<br><br>See [Module-aware commands](https://go.dev/ref/mod#mod-commands) for more information.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| `GOMODCACHE`  | The directory where the `go` command will store downloaded modules and related files. See [Module cache](https://go.dev/ref/mod#module-cache) for details on the structure of this directory.If `GOMODCACHE` is not set, it defaults to `$GOPATH/pkg/mod`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `GOINSECURE`  | Comma-separated list of glob patterns (in the syntax of Go's [`path.Match`](https://go.dev/pkg/path/#Match)) of module path prefixes that may always be fetched in an insecure manner. Only applies to dependencies that are being fetched directly.Unlike the `-insecure` flag on `go get`, `GOINSECURE` does not disable module checksum database validation. `GOPRIVATE` or `GONOSUMDB` may be used to achieve that.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| `GONOPROXY`   | Comma-separated list of glob patterns (in the syntax of Go's [`path.Match`](https://go.dev/pkg/path/#Match)) of module path prefixes that should always be fetched directly from version control repositories, not from module proxies.If `GONOPROXY` is not set, it defaults to `GOPRIVATE`. See [Privacy](https://go.dev/ref/mod#private-module-privacy).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| `GONOSUMDB`   | Comma-separated list of glob patterns (in the syntax of Go's [`path.Match`](https://go.dev/pkg/path/#Match)) of module path prefixes for which the `go` should not verify checksums using the checksum database.If `GONOSUMDB` is not set, it defaults to `GOPRIVATE`. See [Privacy](https://go.dev/ref/mod#private-module-privacy).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| `GOPATH`      | In `GOPATH` mode, the `GOPATH` variable is a list of directories that may contain Go code.In module-aware mode, the [module cache](https://go.dev/ref/mod#glos-module-cache) is stored in the `pkg/mod` subdirectory of the first `GOPATH` directory. Module source code outside the cache may be stored in any directory.If `GOPATH` is not set, it defaults to the `go` subdirectory of the user's home directory.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| `GOPRIVATE`   | Comma-separated list of glob patterns (in the syntax of Go's [`path.Match`](https://go.dev/pkg/path/#Match)) of module path prefixes that should be considered private. `GOPRIVATE` is a default value for `GONOPROXY` and `GONOSUMDB`. See [Privacy](https://go.dev/ref/mod#private-module-privacy). `GOPRIVATE` also determines whether a module is considered private for `GOVCS`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| `GOPROXY`     | List of module proxy URLs, separated by commas (`,`) or pipes (`\|`). When the `go` command looks up information about a module, it contacts each proxy in the list in sequence until it receives a successful response or a terminal error. A proxy may respond with a 404 (Not Found) or 410 (Gone) status to indicate the module is not available on that server.The `go` command's error fallback behavior is determined by the separator characters between URLs. If a proxy URL is followed by a comma, the `go` command falls back to the next URL after a 404 or 410 error; all other errors are considered terminal. If the proxy URL is followed by a pipe, the `go` command falls back to the next source after any error, including non-HTTP errors like timeouts.`GOPROXY` URLs may have the schemes `https`, `http`, or `file`. If a URL has no scheme, `https` is assumed. A module cache may be used directly as a file proxy:<br><br>GOPROXY=file://$(go env GOMODCACHE)/cache/download<br><br>Two keywords may be used in place of proxy URLs:<br><br>- `off`: disallows downloading modules from any source.<br>- `direct`: download directly from version control repositories instead of using a module proxy.<br><br>`GOPROXY` defaults to `https://proxy.golang.org,direct`. Under that configuration, the `go` command first contacts the Go module mirror run by Google, then falls back to a direct connection if the mirror does not have the module. See https://proxy.golang.org/privacy for the mirror's privacy policy. The `GOPRIVATE` and `GONOPROXY` environment variables may be set to prevent specific modules from being downloaded using proxies. See [Privacy](https://go.dev/ref/mod#private-module-privacy) for information on private proxy configuration.See [Module proxies](https://go.dev/ref/mod#module-proxy) and [Resolving a package to a module](https://go.dev/ref/mod#resolve-pkg-mod) for more information on how proxies are used. |
| `GOSUMDB`     | Identifies the name of the checksum database to use and optionally its public key and URL. For example:<br><br>GOSUMDB="sum.golang.org"<br>GOSUMDB="sum.golang.org+<publickey>"<br>GOSUMDB="sum.golang.org+<publickey> https://sum.golang.org"<br><br>The `go` command knows the public key of `sum.golang.org` and also that the name `sum.golang.google.cn` (available inside mainland China) connects to the `sum.golang.org` database; use of any other database requires giving the public key explicitly. The URL defaults to `https://` followed by the database name.`GOSUMDB` defaults to `sum.golang.org`, the Go checksum database run by Google. See [Privacy: Go modules services](https://sum.golang.org/privacy) for the service's privacy policy.If `GOSUMDB` is set to `off` or if `go get` is invoked with the `-insecure` flag, the checksum database is not consulted, and all unrecognized modules are accepted, at the cost of giving up the security guarantee of verified repeatable downloads for all modules. A better way to bypass the checksum database for specific modules is to use the `GOPRIVATE` or `GONOSUMDB` environment variables.See [Authenticating modules](https://go.dev/ref/mod#authenticating) and [Privacy](https://go.dev/ref/mod#private-module-privacy) for more information.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `GOVCS`       | Controls the set of version control tools the `go` command may use to download public and private modules (defined by whether their paths match a pattern in `GOPRIVATE`) or other modules matching a glob pattern.If `GOVCS` is not set, or if a module does not match any pattern in `GOVCS`, the `go` command may use `git` and `hg` for a public module, or any known version control tool for a private module. Concretely, the `go` command acts as if `GOVCS` were set to:<br><br>public:git\|hg,private:all<br><br>See [Controlling version control tools with `GOVCS`](https://go.dev/ref/mod#vcs-govcs) for a complete explanation.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| `GOWORK`      | The `GOWORK` environment variable instructs the `go` command to enter workspace mode using the provided [`go.work` file](#go-work-file) to define the workspace. If `GOWORK` is set to `off` workspace mode is disabled. This can be used to run the `go` command in single module mode: for example, `GOWORK=off go build .` builds the `.` package in single-module mode.`If `GOWORK` is empty, the `go` command will search for a `go.work` file as described in the [Workspaces](#workspaces) section.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |

## Glossary

**build constraint:** 编译包时判断是否使用Go源文件的条件。构建约束可以用文件名后缀（例如，foo_linux_amd64.go）或构建约束注释（例如，//+build linux,amd64）来表示。请参阅[Build Constraints](https://go.dev/pkg/go/build/#hdr-Build_Constraints)。

**build list:** 将用于构建命令（例如 go build、go list 或 go test）的模块版本列表。构建列表是使用最小版本选择从主模块的 go.mod 文件和传递所需模块中的 go.mod 文件确定的。构建列表包含模块图中所有模块的版本，而不仅仅是与特定命令相关的版本。

**canonical version:** 格式正确的版本，除了+不兼容之外，没有构建元数据后缀。例如，v1.2.3 是规范版本，但 v1.2.3+meta 则不是。

**current module:** 主模块的同义词。

**deprecated module:** 其作者不再支持的模块（尽管主要版本为此目的被视为不同的模块）。已弃用的模块在最新版本的 go.mod 文件中标有弃用注释。

**direct dependency:** 其路径出现在主模块中的包或测试的 .go 源文件的导入声明中或包含此类包的模块中的包。 （比较间接依赖。）

**direct mode:** 环境变量的设置，使 go 命令直接从版本控制系统（而不是模块代理）下载模块。 GOPROXY=direct 对所有模块执行此操作。 GOPRIVATE 和 GONOPROXY 对匹配模式列表的模块执行此操作。

**`go.mod` file:** 定义模块的路径、要求和其他元数据的文件。出现在模块的根目录中。请参阅有关 go.mod 文件的部分。

**`go.work` file** 定义要在工作区中使用的模块集的文件。请参阅有关 go.work 文件的部分

**import path:** 用于导入 Go 源文件中的包的字符串。与包路径同义。

**indirect dependency:** 由主模块中的包或测试传递导入的包，但其路径未出现在主模块中的任何导入声明中；或者出现在模块图中但不提供任何由主模块直接导入的包的模块。 （比较直接依赖。）

**lazy module loading:** Go 1.17 中的一项更改是，避免在指定 go 1.17 或更高版本的模块中加载不需要模块图的命令。请参阅延迟模块加载。

**main module:** 调用 go 命令的模块。主模块由当前目录或父目录中的 go.mod 文件定义。请参阅模块、包和版本。

**major version:** 语义版本中的第一个数字（v1.2.3 中为 1）。在具有不兼容更改的版本中，主版本必须递增，并且次要版本和补丁版本必须设置为 0。主版本为 0 的语义版本被认为是不稳定的。

**major version subdirectory:** 版本控制存储库中与模块的主要版本后缀匹配的子目录，可以在其中定义模块。例如，根路径为 example.com/mod 的存储库中的模块 example.com/mod/v2 可以定义在存储库根目录或主版本子目录 v2 中。请参阅存储库中的模块目录。

**major version suffix:** 与主版本号匹配的模块路径后缀。例如，example.com/mod/v2 中的 /v2 v2.0.0 及更高版本需要主版本后缀，而早期版本不允许使用主版本后缀。请参阅有关主要版本后缀的部分。

**minimal version selection (MVS):** 用于确定将在构建中使用的所有模块的版本的算法。有关详细信息，请参阅最小版本选择部分。

**minor version:** 语义版本中的第二个数字（v1.2.3 中为 2）。在具有新的向后兼容功能的版本中，次要版本必须增加，并且补丁版本必须设置为 0。

**module:** 一起发布、版本化和分发的软件包的集合。

**module cache:** 存储下载模块的本地目录，位于 GOPATH/pkg/mod 中。请参阅模块缓存。

**module graph:** 模块需求的有向图，植根于主模块。图中的每个顶点都是一个模块；每个边缘都是 go.mod 文件中 require 语句的一个版本（受主模块的 go.mod 文件中的替换和排除语句的影响）。

**module graph pruning:** Go 1.17 中的一项更改是通过省略指定 go 1.17 或更高版本的模块的传递依赖关系来减小模块图的大小。请参阅模块图修剪。

**module path:** 标识模块并充当模块内包导入路径的前缀的路径。例如，“golang.org/x/net”。

**module proxy:** 实现 GOPROXY 协议的 Web 服务器。 go 命令从模块代理下载版本信息、go.mod 文件和模块 zip 文件。

**module root directory:** 包含定义模块的 go.mod 文件的目录。

**module subdirectory:** 存储库根路径之后的模块路径部分，指示定义模块的子目录。当非空时，模块子目录也是语义版本标签的前缀。模块子目录不包含主版本后缀（如果有），即使该模块位于主版本子目录中。请参阅模块路径。

**package:** 同一目录中编译在一起的源文件的集合。请参阅 Go 语言规范中的包部分。

**package path:** 唯一标识包的路径。包路径是与模块内的子目录连接的模块路径。例如“golang.org/x/net/html”是“html”子目录中模块“golang.org/x/net”中包的包路径。导入路径的同义词。

**patch version:** 语义版本中的第三个数字（v1.2.3 中的 3）。在模块公共接口没有变化的版本中，补丁版本必须增加。

**pre-release version:** 带有破折号的版本，后跟一系列点分隔的标识符，紧跟在补丁版本之后，例如 v1.2.3-beta4。预发布版本被认为不稳定，并且不假定与其他版本兼容。预发布版本排序在相应的发布版本之前：v1.2.3-pre 位于 v1.2.3 之前。另请参阅发布版本。

**pseudo-version:** 对修订标识符（例如 Git 提交哈希）和来自版本控制系统的时间戳进行编码的版本。例如，v0.0.0-20191109021931-daa7c04131f5。用于与非模块存储库兼容以及标记版本不可用的其他情况。

**release version:** 没有预发布后缀的版本。例如，v1.2.3，而不是 v1.2.3-pre。另请参阅预发行版本。

**repository root path:** 模块路径中与版本控制存储库的根目录相对应的部分。请参阅模块路径。

**retracted version:** 一个不应该依赖的版本，要么因为它发布过早，要么因为它发布后发现了严重的问题。请参阅撤消指令。

**semantic version tag:** 版本控制存储库中将版本映射到特定修订版的标记。请参阅将版本映射到提交。

**selected version:** 通过最小版本选择选择的给定模块的版本。所选版本是在模块图中找到的模块路径的最高版本。

**vendor directory:** 名为vendor的目录，包含在主模块中构建包所需的其他模块的包。由 go mod 供应商维护。请参阅“[Vendoring](https://go.dev/ref/mod#vendoring)”。

**version:** 模块的不可变快照的标识符，写为字母 v 后跟语义版本。请参阅版本部分。

**workspace:** 磁盘上的模块集合，在运行最小版本选择 (MVS) 时用作主模块。请参阅有关工作区的部分

## 其他

[Go module机制下升级major版本号的实践](https://tonybai.com/2019/06/03/the-practice-of-upgrading-major-version-under-go-module/)

[Go Modules Reference - The Go Programming Language](https://go.dev/ref/mod)

[通过一个例子让你彻底掌握 Go 工作区模式](https://polarisxu.studygolang.com/posts/go/workspace/)
