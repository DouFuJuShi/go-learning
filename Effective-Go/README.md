# 介绍 Introduction

Go 是一门新语言。尽管它借鉴了现有语言的思想，但它具有不寻常的特性，使得有效的 Go 程序在性质上不同于用其他同类语言编写的程序。将 C++ 或 Java 程序直接翻译成 Go 不太可能产生令人满意的结果——Java 程序是用 Java 编写的，而不是 Go。另一方面，从 Go 的角度思考问题可能会产生一个成功但完全不同的程序。换句话说，要写好Go，理解它的属性和习惯用法很重要。了解 Go 编程的既定约定也很重要，例如命名、格式、程序构造等，以便其他 Go 程序员能够轻松理解您编写的程序。

本文档提供了编写清晰、惯用的 Go 代码的技巧。它补充了语言规范[language specification](https://go.dev/ref/spec)、Go 之旅和如何编写 Go 代码，所有这些您都应该首先阅读。

2022 年 1 月添加的注释：本文档是为 2009 年发布的 Go 编写的，此后没有进行过重大更新。虽然它是了解如何使用该语言本身的一个很好的指南，但由于该语言的稳定性，它很少介绍库，也没有介绍自编写以来 Go 生态系统的重大变化，例如构建系统、测试、模块和多态性。没有更新它的计划，因为已经发生了很多事情，并且大量且不断增长的文档、博客和书籍很好地描述了现代 Go 的用法。 《Effective Go》仍然很有用，但读者应该明白它远不是一个完整的指南。有关上下文，请参阅问题 [28782](https://github.com/golang/go/issues/28782)。

## 例子 Examples

[Go 包源]([- The Go Programming Language](https://go.dev/src/))不仅用作核心库，还用作如何使用该语言的示例。此外，许多包包含可工作的、独立的可执行示例，您可以直接从 golang.org 网站运行，例如这个（如有必要，单击“示例”一词将其打开）。如果您对如何解决问题或如何实现某些内容有疑问，库中的文档、代码和示例可以提供答案、想法和背景。

# 格式 Formatting

格式问题是最具争议性但最不重要的问题。人们可以适应不同的格式样式，但如果没有必要的话最好，如果每个人都遵循相同的样式，那么花在该主题上的时间就会更少。问题是如何在没有长篇规范性风格指南的情况下接近这个乌托邦。

对于 Go，我们采取了一种不寻常的方法，让机器处理大多数格式问题。 gofmt 程序（也可称为 go fmt，它在包级别而不是源文件级别运行）读取 Go 程序并以标准缩进和垂直对齐方式发出源代码，保留注释并在必要时重新格式化注释。如果你想知道如何处理一些新的布局情况，请运行 gofmt;如果答案看起来不正确，请重新安排您的程序（或提交有关 gofmt 的错误），不要解决它。

例如，无需花时间排列结构字段上的注释。 Gofmt 将为您做到这一点。鉴于声明

```go
type T struct {
    name string // name of the object
    value int // its value
}
```

gofmt 会将列排成一列：

```go
type T struct {
    name    string // name of the object
    value   int    // its value
}
```

标准包中的所有 Go 代码均已使用 gofmt 进行格式化。

一些格式细节仍然保留。非常简单地说：

    缩进：我们使用制表符进行缩进，gofmt 默认情况下会发出它们。仅在必要时才使用空格。

    行长：Go 没有行长度限制。不用担心打孔卡溢出。如果感觉某行太长，请将其换行并用一个额外的制表符缩进。

    圆括号：Go 需要的括号比 C 和 Java 少：控制结构（if、for、switch）的语法中没有括号。此外，运算符优先级层次结构更短、更清晰，因此

```go
x<<8 + y<<16
```

与其他语言不同的是，空格意味着什么。

# 注释 Commentary

Go 提供了 C 风格的 /* */ 块注释和 C++ 风格的 // 行注释。行注释是常态；块注释主要显示为包注释，但在表达式中很有用或可禁用大量代码。

出现在顶级声明之前的注释（中间没有换行符）被视为记录声明本身。这些“文档注释”是给定 Go 包或命令的主要文档。有关文档评论的更多信息，请参阅“Go Doc Comments”。

# 命名 Names

在 Go 语言中，名称和其他语言一样重要。它们甚至有语义上的影响：名字在包外的可见性取决于它的第一个字符是否大写。因此，我们值得花一点时间来谈谈 Go 程序中的命名规则。

## 包名 package names

导入包时，包名称将成为内容的访问器。后

```go
import "bytes"
```

导入包可以谈论bytes.Buffer。如果使用该包的每个人都可以使用相同的名称来引用其内容，这将很有帮助，这意味着包名称应该很好：简短、简洁、令人印象深刻。按照惯例，包的名称都是小写的、单字的。不需要下划线或混合大写字母。为了简洁起见，因为每个使用你的包的人都会输入这个名字。并且不必担心先验的冲突。包名只是导入时的默认名称；它不需要在所有源代码中都是唯一的，并且在极少数发生冲突的情况下，导入包可以选择不同的名称在本地使用。无论如何，混淆很少发生，因为导入中的文件名决定了正在使用哪个包。

另一个约定是包名称是其源目录的基本名称； src/encoding/base64 中的包作为“encoding/base64”导入，但名称为base64，而不是encoding_base64，也不是encodingBase64。

包的导入者将使用名称来引用其内容，因此包中的导出名称可以使用该事实来避免重复。 （不要使用 import . 表示法，它可以简化必须在正在测试的包外部运行的测试，但应该避免。）例如，bufio 包中的缓冲读取器类型称为 Reader，而不是 BufReader，因为用户将其视为 bufio.Reader，这是一个清晰、简洁的名称。此外，由于导入的实体始终使用其包名称进行寻址，因此 bufio.Reader 不会与 io.Reader 冲突。类似地，创建ring.Ring新实例的函数（Go中构造函数的定义）通常被称为NewRing，但由于Ring是包导出的唯一类型，并且由于包被称为ring，所以它是称为“New”，包的客户端将其视为“ring.New”。使用包结构来帮助您选择好的名称。

另一个简短的例子是once.Do; Once.Do(setup) 读起来很好，并且不会通过编写 Once.DoOrWaitUntilDone(setup) 来改进。长名称不会自动使内容更具可读性。有用的文档注释通常比超长的名称更有价值。

## Getters

Go 不提供对 getter 和 setter 的自动支持。自己提供 getter 和 setter 并没有什么问题，而且这样做通常是合适的，但将 Get 放入 getter 的名称中既不符合习惯，也没有必要。如果您有一个名为owner（小写，未导出）的字段，则getter方法应称为Owner（大写，已导出），而不是GetOwner。使用大写名称进行导出提供了区分字段和方法的钩子。如果需要，setter 函数可能会被称为 SetOwner。这两个名字在实践中读起来都很好：

```go
owner := obj.Owner()
if owner != user {
    obj.SetOwner(user)
}
```

## 接口名 Interface names

按照惯例，单方法接口由方法名称加上 -er 后缀或类似修饰来命名，以构造代理名词：Reader、Writer、Formatter、CloseNotifier 等。

这样的名称有很多，尊重它们以及它们捕获的函数名称是很有成效的。 Read、Write、Close、Flush、String 等都有规范的签名和含义。为了避免混淆，请勿为您的方法指定这些名称之一，除非它具有相同的签名和含义。相反，如果您的类型实现的方法与众所周知的类型上的方法具有相同含义，请为其指定相同的名称和签名；调用字符串转换器方法 String 而不是 ToString。

## 混合大小写 MixedCaps

最后，Go 中的约定是使用 MixedCaps 或 mixCaps 而不是下划线来书写多单词名称。

# 分号 Semicolons

与 C 一样，Go 的形式语法使用分号来终止语句，但与 C 不同的是，这些分号不会出现在源代码中。相反，词法分析器使用一个简单的规则在扫描时自动插入分号，因此输入文本基本上不含分号。

规则是这样的。如果换行符之前的最后一个标记是标识符（包括 int 和 float64 等字）、数字或字符串常量等基本字面量，或以下标记之一

```go
break continue fallthrough return ++ -- ) }
```

词法分析器总是在标记后面插入分号。这可以概括为，“如果换行符出现在可以结束语句的标记之后，则插入分号”。

分号也可以在紧靠结束括号之前省略，因此语句如

```go
go func() { for { dst <- <-src } }()
```

不需要分号。惯用的 Go 程序仅在 for 循环子句等位置使用分号，以分隔初始值设定项、条件和延续元素。如果您以这种方式编写代码，它们还需要分隔一行中的多个语句。

分号插入规则的一个后果是您不能将控制结构（if、for、switch 或 select）的左大括号放在下一行。如果这样做，将在大括号之前插入分号，这可能会导致不必要的效果。像这样写它们

```go
if i < f() {
    g()
}
```

而不是这样的：

```go
if i < f()  // wrong!
{           // wrong!
    g()
}
```

## 控制结构 Control structures

Go 的控制结构与 C 的控制结构相关，但在一些重要方面有所不同。没有 do 或 while 循环，只有稍微概括的 for；切换更灵活； if 和 switch 接受类似于 for 的可选初始化语句； break 和 continue 语句采用可选标签来标识要中断或继续的内容；并且有新的控制结构，包括类型开关和多路通信多路复用器、选择。语法也略有不同：没有括号，并且正文必须始终用大括号分隔。

## If

在 Go 中，一个简单的 if 看起来像这样：

```go
if x > 0 {
    return y
}
```

强制大括号鼓励在多行上编写简单的 if 语句。无论如何，这样做是一种很好的风格，特别是当主体包含诸如 return 或 break 之类的控制语句时。

由于 if 和 switch 接受初始化语句，因此通常会看到它们用于设置局部变量。

```go
if err := file.Chmod(0664); err != nil {
    log.Print(err)
    return err
}
```

在 Go 库中，您会发现当 if 语句不流入下一个语句时（即主体以 break、continue、goto 或 return 结尾），不必要的 else 就会被省略。

```go
f, err := os.Open(name)
if err != nil {
    return err
}
codeUsing(f)
```

这是一个常见情况的示例，其中代码必须防范一系列错误条件。如果控制流成功地沿着页面运行，则代码可读性良好，从而消除出现的错误情况。由于错误情况往往以 return 语句结束，因此生成的代码不需要 else 语句。

```go
f, err := os.Open(name)
if err != nil {
    return err
}
d, err := f.Stat()
if err != nil {
    f.Close()
    return err
}
codeUsing(f, d)
```

## 重声明和重赋值 Redeclaration and reassignment

题外话：上一节中的最后一个示例演示了 := 短声明形式如何工作的细节。调用 os.Open 的声明如下：

```go
f, err := os.Open(name)
```

该语句声明了两个变量，f 和 err。几行之后，对 f.Stat 的调用显示：

```go
d, err := f.Stat()
```

看起来好像它声明了 d 和 err。但请注意，两个语句中都出现了 err。这种重复是合法的： err 由第一个语句声明，但仅在第二个语句中重新分配。这意味着对 f.Stat 的调用使用上面声明的现有 err 变量，并且只是给它一个新值。

在 := 声明中，即使已经声明了变量 v，也可能会出现，前提是：

- 此声明与 v 的现有声明处于同一范围内（如果 v 已在外部范围中声明，则该声明将创建一个新变量§），

- 初始化中的相应值可分配给 v，并且

- 该声明至少创建了一个其他变量。

这种不寻常的属性是纯粹的实用主义，使得使用单个 err 值变得很容易，例如，在长 if-else 链中。你会看到它经常被使用。

§ 这里值得注意的是，在 Go 中，函数参数和返回值的范围与函数体相同，即使它们在词法上出现在包围函数体的大括号之外。

## For

Go 的 for 循环与 C 的相似，但并不相同。它统一了 for 和 while，没有 do-while。它有三种形式，其中只有一种有分号。

```go
// Like a C for
for init; condition; post { }

// Like a C while
for condition { }

// Like a C for(;;)
for { }
```

简短的声明使得在循环中声明索引变量变得容易。

```go
sum := 0
for i := 0; i < 10; i++ {
    sum += i
}
```

如果要循环数组、切片、字符串或映射，或者从通道读取，则范围子句可以管理循环。

```go
for key, value := range oldMap {
    newMap[key] = value
}
```

如果您只需要范围中的第一项（键或索引），请删除第二项：

```go
for key := range m {
    if key.expired() {
        delete(m, key)
    }
}
```

如果您只需要范围中的第二项（值），请使用空白标识符（下划线）来丢弃第一项：

```go
sum := 0
for _, value := range array {
    sum += value
}
```

空白标识符有很多用途，将在后文介绍。

对于字符串，该范围可以为您完成更多工作，通过解析 UTF-8 来分解各个 Unicode 代码点。错误的编码会消耗一个字节并产生替换符文 U+FFFD。 （名称（以及关联的内置类型）符文是单个 Unicode 代码点的 Go 术语。有关详细信息，请参阅语言规范。）

```go
for pos, char := range "日本\x80語" { // \x80 is an illegal UTF-8 encoding
    fmt.Printf("character %#U starts at byte position %d\n", char, pos)
}
```

打印

```go
character U+65E5 '日' starts at byte position 0
character U+672C '本' starts at byte position 3
character U+FFFD '�' starts at byte position 6
character U+8A9E '語' starts at byte position 7
```

最后，Go 没有逗号运算符，++ 和 -- 是语句而不是表达式。因此，如果你想在 for 中运行多个变量，你应该使用并行赋值（尽管这排除了 ++ 和 --）。

```go
// Reverse a
for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
    a[i], a[j] = a[j], a[i]
}
```

## Switch

Go 的 switch 比 C 的更通用。表达式不一定是常量，甚至也不一定是整数，从上到下对情况进行评估，直到找到匹配为止，如果switch没有表达式，则切换到 true。因此，将 if-else-if-else 链写成 switch 是可能的，也是习以为常的。

```go
func unhex(c byte) byte {
    switch {
    case '0' <= c && c <= '9':
        return c - '0'
    case 'a' <= c && c <= 'f':
        return c - 'a' + 10
    case 'A' <= c && c <= 'F':
        return c - 'A' + 10
    }
    return 0
}
```

不会自动失败，但案例可以以逗号分隔的列表形式呈现。

```go
func shouldEscape(c byte) bool {
    switch c {
    case ' ', '?', '&', '=', '#', '+', '%':
        return true
    }
    return false
}
```

尽管它们在 Go 中不像其他一些类似 C 的语言那样常见，但 Break 语句可用于提前终止 switch。但有时，有必要打破周围的循环，而不是switch，在 Go 中，可以通过在循环上放置标签并“breaking”该标签来完成。此示例显示了两种用途。

```go
Loop:
    for n := 0; n < len(src); n += size {
        switch {
        case src[n] < sizeOne:
            if validateOnly {
                break
            }
            size = 1
            update(src[n])

        case src[n] < sizeTwo:
            if n+1 >= len(src) {
                err = errShortInput
                break Loop
            }
            if validateOnly {
                break
            }
            size = 2
            update(src[n] + src[n+1]<<shift)
        }
    }
```

当然， continue 语句也接受可选标签，但它仅适用于循环。

为了结束本节，这里有一个使用两个 switch 语句的字节片比较例程：

```go
// Compare returns an integer comparing the two byte slices,
// lexicographically.
// The result will be 0 if a == b, -1 if a < b, and +1 if a > b
func Compare(a, b []byte) int {
    for i := 0; i < len(a) && i < len(b); i++ {
        switch {
        case a[i] > b[i]:
            return 1
        case a[i] < b[i]:
            return -1
        }
    }
    switch {
    case len(a) > len(b):
        return 1
    case len(a) < len(b):
        return -1
    }
    return 0
}
```

## Type switch

switch 还可用于发现接口变量的动态类型。这样的类型开关使用类型断言的语法，并将关键字 type 放在括号内。如果 switch 在表达式中声明变量，则该变量将在每个子句中具有相应的类型。在这种情况下重用名称也是惯用的做法，实际上在每种情况下都声明具有相同名称但类型不同的新变量。

# Functions

## 多函数返回值 Multiple return values

Go 的不寻常特性之一是函数和方法可以返回多个值。这种形式可用于改进 C 程序中的一些笨拙的习惯用法：带内错误返回，例如 EOF 的 -1 以及修改通过地址传递的参数。

在 C 语言中，写入错误通过负计数来表示，错误代码隐藏在易失性位置中。在 Go 中，Write 可以返回一个计数和一个错误：“是的，你写了一些字节，但不是全部，因为你填满了设备”。 os 包中文件的 Write 方法的签名是：

```go
func (file *File) Write(b []byte) (n int, err error)
```

正如文档所述，当 n != len(b) 时，它返回写入的字节数和非零错误。这是一种常见的风格；有关更多示例，请参阅有关错误处理的部分。

类似的方法不需要传递指向返回值的指针来模拟引用参数。这是一个简单的函数，用于从字节切片中的某个位置抓取数字，返回该数字和下一个位置。

```go
func nextInt(b []byte, i int) (int, int) {
    for ; i < len(b) && !isDigit(b[i]); i++ {
    }
    x := 0
    for ; i < len(b) && isDigit(b[i]); i++ {
        x = x*10 + int(b[i]) - '0'
    }
    return x, i
}
```

你可以用它来扫描输入片段 b 中的数字，就像这样：

```go
   for i := 0; i < len(b); {
        x, i = nextInt(b, i)
        fmt.Println(x)
    }
```

## 命名结果参数 Named result parameters

Go 函数的返回或结果“参数”可以被命名并用作常规变量，就像传入参数一样。当命名时，它们在函数开始时被初始化为其类型的零值；如果函数执行不带参数的 return 语句，则结果参数的当前值将用作返回值。

名称不是强制性的，但它们可以使代码更短、更清晰：它们是文档。如果我们命名 nextInt 的结果，那么返回的 int 是哪个就很明显了。

```go
func nextInt(b []byte, pos int) (value, nextPos int) {
```

由于命名结果已初始化并与未修饰的返回绑定在一起，因此它们可以简化并澄清。这是 io.ReadFull 的一个版本，很好地使用了它们：

```go
func ReadFull(r Reader, buf []byte) (n int, err error) {
    for len(buf) > 0 && err == nil {
        var nr int
        nr, err = r.Read(buf)
        n += nr
        buf = buf[nr:]
    }
    return
}
```

## 延迟执行 Defer

Go 的 defer 语句安排函数调用（延迟函数）在执行延迟的函数返回之前立即运行。这是一种不寻常但有效的方法来处理诸如无论函数返回哪条路径都必须释放资源等情况。典型的例子是解锁互斥体或关闭文件。

```go
// Contents returns the file's contents as a string.
func Contents(filename string) (string, error) {
    f, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()  // f.Close will run when we're finished.

    var result []byte
    buf := make([]byte, 100)
    for {
        n, err := f.Read(buf[0:])
        result = append(result, buf[0:n]...) // append is discussed later.
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err  // f will be closed if we return here.
        }
    }
    return string(result), nil // f will be closed if we return here.
}
```

推迟对 Close 等函数的调用有两个优点。首先，它保证您永远不会忘记关闭文件，如果您稍后编辑该函数以添加新的返回路径，则很容易犯这个错误。其次，这意味着收盘价位于开盘价附近，这比将其放在函数末尾要清晰得多。

延迟函数的参数（如果函数是方法，则包括接收者）在延迟执行时计算，而不是在调用执行时计算。除了避免担心函数执行时变量值发生变化之外，这意味着单个延迟调用站点可以延迟多个函数执行。这是一个愚蠢的例子。

```go
for i := 0; i < 5; i++ {
    defer fmt.Printf("%d ", i)
}
```

延迟函数按 LIFO 顺序执行，因此此代码将导致函数返回时打印 4 3 2 1 0。一个更合理的例子是通过程序跟踪函数执行的简单方法。我们可以编写几个简单的跟踪例程，如下所示：

```go
func trace(s string)   { fmt.Println("entering:", s) }
func untrace(s string) { fmt.Println("leaving:", s) }

// Use them like this:
func a() {
    trace("a")
    defer untrace("a")
    // do something....
}
```

我们可以通过利用延迟执行时评估延迟函数的参数这一事实来做得更好。跟踪例程可以设置非跟踪例程的参数。这个例子：

```go
func trace(s string) string {
    fmt.Println("entering:", s)
    return s
}

func un(s string) {
    fmt.Println("leaving:", s)
}

func a() {
    defer un(trace("a"))
    fmt.Println("in a")
}

func b() {
    defer un(trace("b"))
    fmt.Println("in b")
    a()
}

func main() {
    b()
}
```

打印

```go
entering: b
in b
entering: a
in a
leaving: a
leaving: b
```

对于习惯于其他语言的块级资源管理的程序员来说，defer 可能看起来很奇怪，但它最有趣和最强大的应用恰恰来自于它不是基于块而是基于函数的事实。在有关恐慌和恢复的部分中，我们将看到其可能性的另一个示例。

# 数据 Data

## 使用new进行分配

Go 有两个分配原语，内置函数 new 和 make。它们做不同的事情并适用于不同的类型，这可能会令人困惑，但规则很简单。我们先来说说新的。它是一个分配内存的内置函数，但与其他一些语言中的同名函数不同，它不会初始化内存，而只是将其清零。也就是说，new(T) 为 T 类型的新项分配归零存储并返回其地址，即 *T 类型的值。在 Go 术语中，它返回一个指向新分配的 T 类型零值的指针。

由于 new 返回的内存被清零，因此在设计数据结构时安排可以使用每种类型的零值而无需进一步初始化是很有帮助的。这意味着该数据结构的用户可以创建新的数据结构并开始工作。例如，bytes.Buffer 的文档指出“Buffer 的零值是一个可供使用的空缓冲区”。同样，sync.Mutex 没有显式构造函数或 Init 方法。相反，sync.Mutex 的零值被定义为未锁定的互斥体。

零值即有用 "属性是临时性的。请看下面的类型声明

```go
type SyncedBuffer struct {
    lock    sync.Mutex
    buffer  bytes.Buffer
}
```

SyncedBuffer 类型的值也可以在分配或声明后立即使用。在下一个代码段中，p 和 v 都可以正常工作，无需进一步安排。

```go
p := new(SyncedBuffer)  // type *SyncedBuffer
var v SyncedBuffer      // type  SyncedBuffer
```

## 构造函数和复合字面量

有时零值不够好，需要初始化构造函数，如本例中派生自包 os.

```go
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := new(File)
    f.fd = fd
    f.name = name
    f.dirinfo = nil
    f.nepipe = 0
    return f
}
```

里面有很多样板。我们可以使用复合文字来简化它，这是一个每次计算时都会创建一个新实例的表达式。

```go
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := File{fd, name, nil, 0}
    return &f
}
```

请注意，与 C 语言不同，返回局部变量的地址是完全没有问题的；与变量相关的存储空间在函数返回后仍然存在。事实上，获取复合字面量的地址会在每次求值时分配一个新的实例，因此我们可以将最后两行合并起来。

```go
return &File{fd, name, nil, 0}
```

复合字面量的字段是按顺序排列的，而且必须全部存在。但是，通过将元素明确标注为字段:值对，初始化符可以以任何顺序出现，缺失的初始化符则保留为各自的零值。因此，我们可以说

```go
return &File{fd: fd, name: name}
```

作为一种限制情况，如果复合文字根本不包含任何字段，它会为该类型创建零值。表达式 new(File) 和 &File{} 是等效的。

还可以为数组、切片和映射创建复合文字，其中字段标签为索引或映射键（视情况而定）。在这些示例中，无论 Enone、Eio 和 Einval 的值如何，只要它们不同，初始化都会起作用。

```go
a := [...]string   {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
s := []string      {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
```

## 使用make进行分配

回到分配。内置函数 make(T, args) 的用途与 new(T) 不同。它仅创建切片、映射和通道，并返回 T 类型（而非 *T）的已初始化（未归零）值。区别的原因是这三种类型在幕后表示对使用前必须初始化的数据结构的引用。例如，切片是一个三项描述符，包含指向数据（在数组内）的指针、长度和容量，并且在这些项初始化之前，切片为零。对于切片、映射和通道，make 会初始化内部数据结构并准备要使用的值。例如，

```go
make([]int, 10, 100)
```

分配一个包含 100 个整数的数组，然后创建一个长度为 10、容量为 100 的切片结构，指向数组的前 10 个元素。 （创建切片时，可以省略容量；有关更多信息，请参阅有关切片的部分。）相反，new([]int) 返回一个指向新分配的、清零切片结构的指针，即指向零切片值。

这些例子说明了 new 和 make 之间的区别。

```go
var p *[]int = new([]int)       // allocates slice structure; *p == nil; rarely useful
var v  []int = make([]int, 100) // the slice v now refers to a new array of 100 ints

// Unnecessarily complex:
var p *[]int = new([]int)
*p = make([]int, 100, 100)

// Idiomatic:
v := make([]int, 100)
```

请记住，make 仅适用于映射、切片和通道，并且不返回指针。要获得显式指针，请使用 new 分配或显式获取变量的地址。

## Arrays

数组在规划内存的详细布局时非常有用，有时还能帮助避免分配，但数组主要是切片的构件，也就是下一节的主题。为了给下一节的主题打下基础，这里先介绍一下数组。

数组在 Go 和 C 中的工作方式存在重大差异。在 Go 中，

- 数组是值。将一个数组分配给另一个数组会复制所有元素。

- 特别是，如果将数组传递给函数，它将收到数组的副本，而不是指向它的指针。

- 数组的大小是其类型的一部分。类型 [10]int 和 [20]int 是不同的。

value 属性可能有用，但也很昂贵；如果你想要类似 C 的行为和效率，你可以传递一个指向数组的指针。

```go
func Sum(a *[3]float64) (sum float64) {
    for _, v := range *a {
        sum += v
    }
    return
}

array := [...]float64{7.0, 8.5, 9.1}
x := Sum(&array)  // Note the explicit address-of operator
```

但即使这种风格也不是 Go 的惯用风格。使用切片代替。

## Slices

切片封装了数组，为数据序列提供了一个更通用、更强大、更方便的接口。除了具有明确维度的项目（如变换矩阵）外，Go 中的大多数数组编程都是通过切片而非简单数组完成的。

切片保存对基础数组的引用，如果将一个切片分配给另一个切片，则两个切片都引用同一个数组。如果函数采用切片参数，则调用者可以看到它对切片元素所做的更改，类似于将指针传递给底层数组。因此，Read 函数可以接受切片参数，而不是指针和计数；切片内的长度设置了要读取的数据量的上限。以下是 os 包中 File 类型的 Read 方法的签名：

```go
func (f *File) Read(buf []byte) (n int, err error)
```

该方法返回读取的字节数和错误值（如果有）。要读入较大缓冲区 buf 的前 32 个字节，请对缓冲区进行切片（此处用作动词）。

```go
 n, err := f.Read(buf[0:32])
```

这种分片方法既常见又高效。事实上，抛开效率不谈，下面的代码段也可以读取缓冲区的前 32 个字节。

```go
    var n int
    var err error
    for i := 0; i < 32; i++ {
        nbytes, e := f.Read(buf[i:i+1])  // Read one byte.
        n += nbytes
        if nbytes == 0 || e != nil {
            err = e
            break
        }
    }

```

切片的长度可以更改，只要它仍然符合底层数组的限制即可；只需将其分配给其自身的一部分即可。切片的容量可通过内置函数 cap 访问，报告切片可以采用的最大长度。这是一个将数据附加到切片的函数。如果数据超出容量，则重新分配切片。返回结果切片。该函数利用了 len 和 cap 在应用于 nil 切片时合法的事实，并返回 0。

```go
func Append(slice, data []byte) []byte {
    l := len(slice)
    if l + len(data) > cap(slice) {  // reallocate
        // Allocate double what's needed, for future growth.
        newSlice := make([]byte, (l+len(data))*2)
        // The copy function is predeclared and works for any slice type.
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:l+len(data)]
    copy(slice[l:], data)
    return slice
}
```

之后我们必须返回切片，因为虽然 Append 可以修改切片的元素，但切片本身（保存指针、长度和容量的运行时数据结构）是按值传递的。
附加到切片的想法非常有用，它被附加内置函数捕获。不过，为了理解该函数的设计，我们需要更多信息，因此我们稍后会再讨论它。

## Two-dimensional slices 二维切片

Go 的数组和切片是一维的。要创建 等价的二维 数组或切片，需要定义数组的数组或切片的切片，如下所示：

```go
type Transform [3][3]float64  // A 3x3 array, really an array of arrays.
type LinesOfText [][]byte     // A slice of byte slices.
```

因为切片是可变长度的，所以每个内部切片可以具有不同的长度。这可能是一种常见的情况，就像我们的 LinesOfText 示例中一样：每行都有独立的长度。

```go
text := LinesOfText{
    []byte("Now is the time"),
    []byte("for all good gophers"),
    []byte("to bring some fun to the party."),
}
```

有时需要分配 2D 切片，例如在处理像素扫描线时可能会出现这种情况。有两种方法可以实现这一目标。一是独立分配每个切片；另一种是分配一个数组并将各个切片指向其中。使用哪个取决于您的应用程序。如果切片可能增大或缩小，则应独立分配它们以避免覆盖下一行；如果不是，则使用单次分配来构造对象可能会更有效。作为参考，这里是这两种方法的草图。首先，一次一行：

```go
// Allocate the top-level slice.
picture := make([][]uint8, YSize) // One row per unit of y.
// Loop over the rows, allocating the slice for each row.
for i := range picture {
    picture[i] = make([]uint8, XSize)
}
```

现在作为一个分配，分成几行：

```go
// Allocate the top-level slice, the same as before.
picture := make([][]uint8, YSize) // One row per unit of y.
// Allocate one large slice to hold all the pixels.
pixels := make([]uint8, XSize*YSize) // Has type []uint8 even though picture is [][]uint8.
// Loop over the rows, slicing each row from the front of the remaining pixels slice.
for i := range picture {
    picture[i], pixels = pixels[:XSize], pixels[XSize:]
}
```

## Maps

映射是一种方便而强大的内置数据结构，它将一种类型的值（键）与另一种类型的值（元素或值）相关联。键可以是定义了相等运算符的任何类型，例如整数、浮点和复数、字符串、指针、接口（只要动态类型支持相等）、结构和数组。切片不能用作映射键，因为它们上没有定义相等性。与切片一样，映射保存对底层数据结构的引用。如果将映射传递给更改映射内容的函数，则更改将在调用者中可见。

可以使用通常的复合文字语法和冒号分隔的键值对来构建映射，因此在初始化期间构建它们很容易。

```go
var timeZone = map[string]int{
    "UTC":  0*60*60,
    "EST": -5*60*60,
    "CST": -6*60*60,
    "MST": -7*60*60,
    "PST": -8*60*60,
}
```

分配和获取映射值在语法上看起来就像对数组和切片执行相同操作，只是索引不需要是整数。

```go
offset := timeZone["EST"]
```

尝试使用映射中不存在的键获取映射值将返回映射中条目类型的零值。例如，如果映射包含整数，则查找不存在的键将返回 0。集合可以实现为值类型 bool 的映射。将映射条目设置为true以将值放入集合中，然后通过简单的索引进行测试。

```go
attended := map[string]bool{
    "Ann": true,
    "Joe": true,
    ...
}

if attended[person] { // will be false if person is not in the map
    fmt.Println(person, "was at the meeting")
}
```

有时您需要区分缺失条目和零值。是否有“UTC”条目，或者是 0 因为它根本不在地图中？您可以通过多重赋值的形式进行区分。

```go
var seconds int
var ok bool
seconds, ok = timeZone[tz]
```

出于显而易见的原因，这被称为“, ok”习语。在此示例中，如果 tz 存在，则秒数将被适当设置，并且 ok 将为 true；如果不是，秒数将设置为零，并且 ok 将为 false。这是一个将其与漂亮的错误报告结合在一起的函数：

```go
func offset(tz string) int {
    if seconds, ok := timeZone[tz]; ok {
        return seconds
    }
    log.Println("unknown time zone:", tz)
    return 0
}
```

要测试Map中是否存在而不关心实际值，可以使用空白标识符 (_) 代替值的常用变量。

```go
_, present := timeZone[tz]
```

要删除映射条目，请使用删除内置函数，其参数是映射和要删除的键。即使钥匙已经不在地图上，这样做也是安全的。

```go
delete(timeZone, "PDT")  // Now on Standard Time
```

## Printing 格式化打印

Go 中的格式化打印与 C 的 printf 系列风格相似，但内容更丰富、更通用。这些函数位于 fmt 包中，名称大写：fmt.Printf、fmt.Fprintf、fmt.Sprintf 等。字符串函数（Sprintf 等）返回一个字符串，而不是填入提供的缓冲区。

您不需要提供格式字符串。对于 Printf、Fprintf 和 Sprintf 中的每一个，都有另一对函数，例如 Print 和 Println。这些函数不采用格式字符串，而是为每个参数生成默认格式。 Println 版本还在参数之间插入空格，并向输出附加换行符，而 Print 版本仅当两侧的操作数都不是字符串时才添加空格。在此示例中，每行都会产生相同的输出。

```go
fmt.Printf("Hello %d\n", 23)
fmt.Fprint(os.Stdout, "Hello ", 23, "\n")
fmt.Println("Hello", 23)
fmt.Println(fmt.Sprint("Hello ", 23))
```

格式化打印函数 fmt.Fprint 和类似函数将任何实现 io.Writer 接口的对象作为第一个参数；变量 os.Stdout 和 os.Stderr 是熟悉的实例。

这里事情开始与 C 有所不同。首先，诸如 %d 之类的数字格式不采用符号或大小标志；相反，打印例程使用参数的类型来决定这些属性。

```go
var x uint64 = 1<<64 - 1
fmt.Printf("%d %x; %d %x\n", x, x, int64(x), int64(x))
```

 打印

```go
18446744073709551615 ffffffffffffffff; -1 -1
```

如果只需要默认的转换，例如整数转换为十进制，可以使用通用格式 %v（表示 "值"）；其结果与 Print 和 Println 产生的结果完全相同。此外，这种格式可以打印任何值，甚至包括数组、片、结构和映射。下面是上一节中定义的时区映射的打印语句。

```go
fmt.Printf("%v\n", timeZone)  // or just fmt.Println(timeZone)
```

给出输出：

```go
map[CST:-21600 EST:-18000 MST:-25200 PST:-28800 UTC:0]
```

对于Map，Printf 和朋友按键按字典顺序对输出进行排序。

在打印结构体时，修改后的格式 %+v 会用结构体的名称注释字段，而对于任何值，备用格式 %#v 会以完整的 Go 语法打印该值。

```go
type T struct {
    a int
    b float64
    c string
}
t := &T{ 7, -2.35, "abc\tdef" }
fmt.Printf("%v\n", t)
fmt.Printf("%+v\n", t)
fmt.Printf("%#v\n", t)
fmt.Printf("%#v\n", timeZone)
```

打印

```go
&{7 -2.35 abc   def}
&{a:7 b:-2.35 c:abc     def}
&main.T{a:7, b:-2.35, c:"abc\tdef"}
map[string]int{"CST":-21600, "EST":-18000, "MST":-25200, "PST":-28800, "UTC":0}
```

（注意 & 符号。）当应用于 string 或 []byte 类型的值时，带引号的字符串格式也可以通过 %q 获得。如果可能的话，备用格式 %#q 将使用反引号。 （%q 格式也适用于整数和符文，生成单引号符文常量。）此外，%x 适用于字符串、字节数组和字节切片以及整数，生成长十六进制字符串，并带有空格在格式 (% x) 中，它在字节之间放置空格。
另一种方便的格式是 %T，它打印值的类型。

```go
fmt.Printf("%T\n", timeZone)
```

打印

```go
map[string]int
```

如果要控制自定义类型的默认格式，只需在类型上定义一个签名为 String() string 的方法即可。对于我们的简单类型 T，可以这样定义

```go
func (t *T) String() string {
    return fmt.Sprintf("%d/%g/%q", t.a, t.b, t.c)
}
fmt.Printf("%v\n", t)
```

以格式打印

```go
7/-2.35/"abc\tdef"
```

（如果您需要打印 T 类型的值以及指向 T 的指针，则 String 的接收器必须是值类型；此示例使用指针，因为这对于结构类型来说更高效且更惯用。请参阅下面有关指针与指针的部分。价值接收者了解更多信息。）

我们的 String 方法能够调用 Sprintf，因为打印例程是完全可重入的并且可以通过这种方式包装。然而，关于这种方法有一个重要的细节需要理解：不要通过调用 Sprintf 来构造 String 方法，这种方式会无限期地重复出现在 String 方法中。如果 Sprintf 调用尝试将接收者直接打印为字符串，从而再次调用该方法，就会发生这种情况。正如本例所示，这是一个常见且容易犯的错误。

```go
type MyString string

func (m MyString) String() string {
    return fmt.Sprintf("MyString=%s", m) // Error: will recur forever.
}
```

这也很容易解决：将参数转换为基本字符串类型，而基本字符串类型没有该方法。

```go
type MyString string
func (m MyString) String() string {
    return fmt.Sprintf("MyString=%s", string(m)) // OK: note conversion.
}
```

在初始化部分，我们将看到另一种避免这种递归的技术。

另一种打印技术是将打印例程的参数直接传递给另一个这样的例程。 Printf 的签名使用 ...interface{} 类型作为其最终参数，以指定任意数量的参数（任意类型）可以出现在格式之后。

```go
func Printf(format string, v ...interface{}) (n int, err error) {
```

在函数 Printf 中，v 的作用类似于 []interface{} 类型的变量，但如果将其传递给另一个可变参数函数，则它的作用类似于常规参数列表。这是我们上面使用的函数 log.Println 的实现。它将其参数直接传递给 fmt.Sprintln 进行实际格式化。

```go
// Println prints to the standard logger in the manner of fmt.Println.
func Println(v ...interface{}) {
    std.Output(2, fmt.Sprintln(v...))  // Output takes parameters (int, string)
}
```

我们在 Sprintln 的嵌套调用中的 v 之后写上 ... 来告诉编译器将 v 视为参数列表；否则它只会将 v 作为单个切片参数传递。

打印的内容比我们在这里介绍的还要多。有关详细信息，请参阅 fmt 包的 godoc 文档。
顺便说一句， ... 参数可以是特定类型，例如 ...int ，用于选择整数列表中最小的一个 min 函数：

```go
func Min(a ...int) int {
    min := int(^uint(0) >> 1)  // largest int
    for _, i := range a {
        if i < min {
            min = i
        }
    }
    return min
}
```

## Append

现在我们有了解释append内置函数的设计所需的缺失部分。 append 的签名与上面我们自定义的 Append 函数不同。概括地说，它是这样的：

```go
func append(slice []T, elements ...T) []T
```

其中 T 是任何给定类型的占位符。实际上，您无法在 Go 中编写类型 T 由调用者确定的函数。这就是为什么内置了append：它需要编译器的支持。

append 的作用是将元素追加到切片的末尾并返回结果。需要返回结果，因为与我们手写的 Append 一样，底层数组可能会发生变化。这个简单的例子

```go
x := []int{1,2,3}
x = append(x, 4, 5, 6)
fmt.Println(x)
```

打印 [1 2 3 4 5 6]。因此，append 的工作方式有点像 Printf，收集任意数量的参数。

但是，如果我们想做我们的 Append 所做的事情，将一个切片附加到一个切片上，该怎么办？简单：在调用站点使用 ...，就像我们在上面调用 Output 中所做的那样。此代码片段产生与上面的输出相同的输出。

```go
x := []int{1,2,3}
y := []int{4,5,6}
x = append(x, y...)
fmt.Println(x)
```

如果没有这个......，它就不会编译，因为类型会错误； y 不是 int 类型。

# Initialization 初始化

虽然从表面上看它与 C 或 C++ 中的初始化没有太大区别，但 Go 中的初始化更强大。可以在初始化期间构建复杂的结构，并且可以正确处理初始化对象之间（甚至不同包之间）的排序问题。

## Constants 常量
