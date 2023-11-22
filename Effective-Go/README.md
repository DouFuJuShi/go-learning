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

```
