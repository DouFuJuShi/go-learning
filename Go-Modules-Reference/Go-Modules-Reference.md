# Go模块手册<sub>---[Go Modules Reference](https://go.dev/ref/mod)</sub>

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
