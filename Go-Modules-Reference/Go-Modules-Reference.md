# Go Modules Reference<sub>[https://go.dev/ref/mod]([Go Modules Reference - The Go Programming Language](https://go.dev/ref/mod))</sub>

## 介绍

模块是 Go 管理依赖关系的方式。

本文档是 Go 模块系统的详细参考手册。有关创建 Go 项目的介绍，请参阅如何编写 Go 代码[How to Write Go Code](https://go.dev/doc/code.html).。有关使用模块、将项目迁移到模块以及其他主题的信息，请参阅从[Using Go Modules](https://blog.golang.org/using-go-modules)。

## Module 模块、Package包和 Version 版本

module 模块是一起发布、版本控制和分发的包的集合。模块可以直接从版本控制存储库或模块代理服务器下载。

模块由[模块路径](https://go.dev/ref/mod#glos-module-path)（在 go.mod 文件中声明）以及有关模块依赖项的信息来标识。模块根目录是包含 go.mod 文件的目录。主模块是包含调用 go 命令的目录的模块。

模块中的每个包都是同一目录中编译在一起的源文件的集合。包路径是与包含包的子目录连接的模块路径（相对于模块根）。例如，模块“golang.org/x/net”包含目录“html”中的包。该包的路径是“golang.org/x/net/html”。

### Module paths 某块路径

模块路径module path是模块的规范名称，在模块的 go.mod 文件中使用 module 指令声明。模块的路径是模块内包路径的前缀。

模块路径应该描述模块的功能以及在哪里可以找到它。通常，模块路径由存储库根路径、存储库中的目录（通常为空）和主要版本后缀（仅适用于主要版本 2 或更高版本）组成。

- *存储库根路径*是模块路径的一部分，对应于开发模块的版本控制存储库的根目录。大多数模块都是在其存储库的根目录中定义的，因此这通常是整个路径。例如，golang.org/x/net 是同名模块的存储库根路径。有关 go 命令如何使用从模块路径派生的 HTTP 请求来查找存储库的信息，请参阅查找模块路径的存储库[Finding a repository for a module path](https://go.dev/ref/mod#vcs-find)。

- 如果模块未在存储库的根目录中定义，则*模块子目录*是命名该目录的模块路径的一部分，不包括主版本后缀。这也用作语义版本标签的前缀。例如，模块 golang.org/x/tools/gopls 位于根路径为 golang.org/x/tools 的存储库的 gopls 子目录中，因此它具有模块子目录 gopls。请参阅将版本映射到存储库中的提交和模块目录。

- 如果模块以主版本 2 或更高版本发布，则模块路径必须以主版本后缀（如 /v2）结尾。这可能是也可能不是子目录名称的一部分。例如，路径为 golang.org/x/repo/sub/v2 的模块可能位于存储库 golang.org/x/repo 的 /sub 或 /sub/v2 子目录中。

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

```
