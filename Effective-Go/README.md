# 介绍 Introduction

Go 是一门新语言。尽管它借鉴了现有语言的思想，但它具有不寻常的特性，使得有效的 Go 程序在性质上不同于用其他同类语言编写的程序。将 C++ 或 Java 程序直接翻译成 Go 不太可能产生令人满意的结果——Java 程序是用 Java 编写的，而不是 Go。另一方面，从 Go 的角度思考问题可能会产生一个成功但完全不同的程序。换句话说，要写好Go，理解它的属性和习惯用法很重要。了解 Go 编程的既定约定也很重要，例如命名、格式、程序构造等，以便其他 Go 程序员能够轻松理解您编写的程序。

本文档提供了编写清晰、惯用的 Go 代码的技巧。它补充了语言规范[language specification](https://go.dev/ref/spec)、Go 之旅和如何编写 Go 代码，所有这些您都应该首先阅读。

2022 年 1 月添加的注释：本文档是为 2009 年发布的 Go 编写的，此后没有进行过重大更新。虽然它是了解如何使用该语言本身的一个很好的指南，但由于该语言的稳定性，它很少介绍库，也没有介绍自编写以来 Go 生态系统的重大变化，例如构建系统、测试、模块和多态性。没有更新它的计划，因为已经发生了很多事情，并且大量且不断增长的文档、博客和书籍很好地描述了现代 Go 的用法。 《Effective Go》仍然很有用，但读者应该明白它远不是一个完整的指南。有关上下文，请参阅问题 [28782](https://github.com/golang/go/issues/28782)。

## 例子 Examples

Go 包源([The Go Programming Language](https://go.dev/src/))不仅用作核心库，还用作如何使用该语言的示例。此外，许多包包含可工作的、独立的可执行示例，您可以直接从 golang.org 网站运行，例如这个（如有必要，单击“示例”一词将其打开）。如果您对如何解决问题或如何实现某些内容有疑问，库中的文档、代码和示例可以提供答案、想法和背景。

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

## Constants

Go 中的常量就是常量。它们是在编译时创建的，即使在函数中定义为局部变量也是如此，并且只能是数字、字符（Runes）、字符串或布尔值。由于编译时限制，定义它们的表达式必须是常量表达式，可由编译器计算。例如， 1<<3 是常量表达式，而 math.Sin(math.Pi/4) 不是，因为对 math.Sin 的函数调用需要在运行时发生。

在 Go 中，枚举常量是使用 iota 枚举器创建的。由于 iota 可以是表达式的一部分，并且表达式可以隐式重复，因此很容易构建复杂的值集。

```go
type ByteSize float64

const (
    _           = iota // ignore first value by assigning to blank identifier
    KB ByteSize = 1 << (10 * iota)
    MB
    GB
    TB
    PB
    EB
    ZB
    YB
)
```

将诸如 String 之类的方法附加到任何用户定义的类型的能力使得任意值可以自动格式化自己以进行打印。尽管您会看到它最常应用于结构，但此技术对于标量类型（例如 ByteSize 等浮点类型）也很有用。

```go
func (b ByteSize) String() string {
    switch {
    case b >= YB:
        return fmt.Sprintf("%.2fYB", b/YB)
    case b >= ZB:
        return fmt.Sprintf("%.2fZB", b/ZB)
    case b >= EB:
        return fmt.Sprintf("%.2fEB", b/EB)
    case b >= PB:
        return fmt.Sprintf("%.2fPB", b/PB)
    case b >= TB:
        return fmt.Sprintf("%.2fTB", b/TB)
    case b >= GB:
        return fmt.Sprintf("%.2fGB", b/GB)
    case b >= MB:
        return fmt.Sprintf("%.2fMB", b/MB)
    case b >= KB:
        return fmt.Sprintf("%.2fKB", b/KB)
    }
    return fmt.Sprintf("%.2fB", b)
}
```

表达式 YB 打印为 1.00YB，而 ByteSize(1e13) 打印为 9.09TB。

这里使用 Sprintf 来实现 ByteSize 的 String 方法是安全的（避免无限期地重复），不是因为转换，而是因为它使用 %f 调用 Sprintf，这不是字符串格式：Sprintf 仅在需要字符串时才会调用 String 方法，%f 想要一个浮点值。

## Variables

变量可以像常量一样初始化，但初始化器可以是运行时计算的通用表达式。

```go
var (
    home   = os.Getenv("HOME")
    user   = os.Getenv("USER")
    gopath = os.Getenv("GOPATH")
)
```

## The init function

最后，每个源文件都可以定义自己的 niladic init 函数来设置所需的任何状态。 （实际上每个文件可以有多个 init 函数。）finally 的意思是finally：init 在包中的所有变量声明都评估了它们的初始值设定项之后被调用，并且只有在所有导入的包都被初始化之后才评估这些变量。

除了不能用声明表示的初始化之外，init 函数的常见用途是在真正执行开始之前验证或修复程序状态的正确性。

```go
func init() {
    if user == "" {
        log.Fatal("$USER not set")
    }
    if home == "" {
        home = "/home/" + user
    }
    if gopath == "" {
        gopath = home + "/go"
    }
    // gopath may be overridden by --gopath flag on command line.
    flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}
```

# Methods

## Pointers vs. Values

正如我们在 ByteSize 中看到的，可以为任何命名类型（指针或接口除外）定义方法；接收者不必是结构体。

在上面对切片的讨论中，我们编写了一个 Append 函数。我们可以将其定义为切片上的方法。为此，我们首先声明一个可以将方法绑定到的命名类型，然后使该方法的接收者成为该类型的值。

```go
type ByteSlice []byte

func (slice ByteSlice) Append(data []byte) []byte {
    // Body exactly the same as the Append function defined above.
}
```

这仍然需要返回更新后的切片的方法。我们可以通过重新定义该方法以将指向 ByteSlice 的指针作为其接收者来消除这种笨拙，以便该方法可以覆盖调用者的切片。

```go
func (p *ByteSlice) Append(data []byte) {
    slice := *p
    // Body as above, without the return.
    *p = slice
}
```

事实上，我们可以做得更好。如果我们修改我们的函数，使其看起来像标准的 Write 方法，如下所示，

```go
func (p *ByteSlice) Write(data []byte) (n int, err error) {
    slice := *p
    // Again as above.
    *p = slice
    return len(data), nil
}
```

那么 *ByteSlice 类型满足标准接口 io.Writer，这很方便。例如，我们可以打印成一张。

```go
  var b ByteSlice
  fmt.Fprintf(&b, "This hour has %d days\n", 7)
```

我们传递 ByteSlice 的地址，因为只有 *ByteSlice 满足 io.Writer。关于接收者的指针与值的规则是，可以在指针和值上调用值方法，但只能在指针上调用指针方法。

出现这条规则是因为指针方法可以修改接收者；对值调用它们将导致该方法接收该值的副本，因此任何修改都将被丢弃。因此，该语言不允许出现这种错误。不过，有一个方便的例外。当值是可寻址的时，语言会通过自动插入地址运算符来处理对值调用指针方法的常见情况。在我们的示例中，变量 b 是可寻址的，因此我们可以仅使用 b.Write 调用其 Write 方法。编译器会将其重写为 (&b).Write 。

顺便说一句，在字节片上使用 Write 的想法是 bytes.Buffer 实现的核心。

# Interfaces and other types

## Interfaces

Go 中的接口提供了一种指定对象行为的方法：如果某个东西可以做到这一点，那么它就可以在这里使用。我们已经看到了几个简单的例子；自定义打印机可以通过 String 方法实现，而 Fprintf 可以使用 Write 方法生成任何内容的输出。仅具有一两个方法的接口在 Go 代码中很常见，并且通常会根据方法指定一个派生名称，例如 io.Writer 表示实现 Write 的东西。

一个类型可以实现多个接口。例如，如果集合实现了 sort.Interface，则可以通过包 sort 中的例程进行排序，其中包含 Len()、Less(i, j int) bool 和 Swap(i, j int)，并且它还可以具有自定义格式化程序。在这个人为的例子中，序列满足了两者。

```go
type Sequence []int

// Methods required by sort.Interface.
func (s Sequence) Len() int {
    return len(s)
}
func (s Sequence) Less(i, j int) bool {
    return s[i] < s[j]
}
func (s Sequence) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

// Copy returns a copy of the Sequence.
func (s Sequence) Copy() Sequence {
    copy := make(Sequence, 0, len(s))
    return append(copy, s...)
}

// Method for printing - sorts the elements before printing.
func (s Sequence) String() string {
    s = s.Copy() // Make a copy; don't overwrite argument.
    sort.Sort(s)
    str := "["
    for i, elem := range s { // Loop is O(N²); will fix that in next example.
        if i > 0 {
            str += " "
        }
        str += fmt.Sprint(elem)
    }
    return str + "]"
}
```

## Conversions

Sequence 的 String 方法正在重新创建 Sprint 已经为切片所做的工作。 （它的复杂度也为 O(N²)，这很差。）如果我们在调用 Sprint 之前将 Sequence 转换为普通的 []int，我们就可以分担工作量（并且还可以加快速度）。

```go
func (s Sequence) String() string {
    s = s.Copy()
    sort.Sort(s)
    return fmt.Sprint([]int(s))
}
```

此方法是从 String 方法安全调用 Sprintf 的转换技术的另一个示例。因为如果忽略类型名称，这两种类型（Sequence 和 []int）是相同的，因此它们之间的转换是合法的。转换不会创建新值，它只是暂时表现为现有值具有新类型。 （还有其他合法的转换，例如从整数到浮点数，确实会创建新值。）

在 Go 程序中，转换表达式的类型以访问一组不同的方法是一种惯用手法。举例来说，我们可以使用现有的 sort.IntSlice 类型将整个示例简化为这样：

```go
type Sequence []int

// Method for printing - sorts the elements before printing
func (s Sequence) String() string {
    s = s.Copy()
    sort.IntSlice(s).Sort()
    return fmt.Sprint([]int(s))
}
```

现在，我们不再让 Sequence 实现多个接口（排序和打印），而是使用将数据项转换为多种类型（Sequence、sort.IntSlice 和 []int）的能力，每种类型都执行以下操作的一部分工作。这在实践中比较不寻常，但可能很有效。

## Interface conversions and type assertions 接口转换和类型断言

类型开关是一种转换形式：它们采用一个接口，并且对于开关中的每种情况，在某种意义上将其转换为该情况的类型。下面是 fmt.Printf 下的代码如何使用类型开关将值转换为字符串的简化版本。如果它已经是一个字符串，我们需要接口保存的实际字符串值，而如果它有一个 String 方法，我们需要调用该方法的结果。

```go
type Stringer interface {
    String() string
}

var value interface{} // Value provided by caller.
switch str := value.(type) {
case string:
    return str
case Stringer:
    return str.String()
}
```

第一种情况找到了具体的值；第二个将接口转换为另一个接口。以这种方式混合类型是完全可以的。

如果我们只关心一种类型怎么办？如果我们知道该值包含一个字符串而我们只想提取它怎么办？单一情况类型开关可以，但类型断言也可以。类型断言采用接口值并从中提取指定显式类型的值。该语法借用了打开类型开关的子句，但使用显式类型而不是 type 关键字：

```go
value.(typeName)
```

结果是一个静态类型为 typeName 的新值。该类型必须是接口所持有的具体类型，或者是该值可以转换为的第二个接口类型。要提取我们知道值中的字符串，我们可以编写：

```go
str := value.(string)
```

但如果结果发现该值不包含字符串，则程序将因运行时错误而崩溃。为了防止这种情况，请使用“comma, ok”习惯用法来安全地测试该值是否是字符串：

```go
str, ok := value.(string)
if ok {
    fmt.Printf("string value is: %q\n", str)
} else {
    fmt.Printf("value is not a string\n")
}
```

如果类型断言失败，str 仍然存在并且是字符串类型，但它将具有零值，即空字符串。

作为功​​能的说明，这里有一个 if-else 语句，它相当于打开本节的类型开关。

```go
if str, ok := value.(string); ok {
    return str
} else if str, ok := value.(Stringer); ok {
    return str.String()
}
```

## Generality 通用性

如果某个类型的存在只是为了实现一个接口，并且永远不会导出该接口之外的方法，则无需导出该类型本身。仅导出接口可以清楚地表明该值没有超出接口中描述的行为。它还避免了对公共方法的每个实例重复文档的需要。

在这种情况下，构造函数应该返回接口值而不是实现类型。例如，在哈希库中，crc32.NewIEEE 和 adler32.New 都返回接口类型 hash.Hash32。在Go程序中用CRC-32算法替换Adler-32只需要更改构造函数调用即可；其余代码不受算法更改的影响。

类似的方法允许将各种加密包中的流式密码算法与它们链接在一起的块密码分开。 crypto/cipher 包中的 Block 接口指定块密码的行为，它提供单个数据块的加密。然后，与 bufio 包类比，实现该接口的密码包可用于构造由 Stream 接口表示的流密码，而无需知道块加密的细节。

加密/密码接口如下所示：

```go
type Block interface {
    BlockSize() int
    Encrypt(dst, src []byte)
    Decrypt(dst, src []byte)
}

type Stream interface {
    XORKeyStream(dst, src []byte)
}
```

这是计数器模式（CTR）流的定义，它将分组密码转换为流密码；请注意，分组密码的详细信息已被抽象掉：

```go
// NewCTR returns a Stream that encrypts/decrypts using the given Block in
// counter mode. The length of iv must be the same as the Block's block size.
func NewCTR(block Block, iv []byte) Stream
```

NewCTR不仅适用于一种特定的加密算法和数据源，而且适用于Block接口和任何Stream的任何实现。由于它们返回接口值，因此用其他加密模式替换 CTR 加密是局部更改。必须编辑构造函数调用，但由于周围的代码必须仅将结果视为 Stream，因此它不会注意到差异。

## Interfaces and methods

由于几乎任何东西都可以附加方法，因此几乎任何东西都可以满足接口。 http 包中就有一个说明性示例，它定义了 Handler 接口。任何实现 Handler 的对象都可以处理 HTTP 请求。

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

ResponseWriter 本身是一个接口，提供对将响应返回给客户端所需的方法的访问。这些方法包括标准 Write 方法，因此只要可以使用 io.Writer，就可以使用 http.ResponseWriter。 Request 是一个结构体，包含来自客户端的请求的解析表示。

为了简洁起见，我们忽略 POST 并假设 HTTP 请求始终是 GET；这种简化不会影响处理程序的设置方式。这是一个处理程序的简单实现，用于计算页面被访问的次数。

```go
// Simple counter server.
type Counter struct {
    n int
}

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ctr.n++
    fmt.Fprintf(w, "counter = %d\n", ctr.n)
}
```

（与我们的主题保持一致，请注意 Fprintf 如何打印到 http.ResponseWriter。）在真实服务器中，对 ctr.n 的访问需要防止并发访问。请参阅sync 和atomic 包以获取建议。

以下是如何将这样的服务器附加到 URL 树上的一个节点上，以供参考。

```go
import "net/http"
...
ctr := new(Counter)
http.Handle("/counter", ctr)
```

但为什么要让 Counter 成为一个结构体呢？只需要一个整数即可。 （接收者必须是一个指针，以便调用者可以看到增量。）

```go
// Simpler counter server.
type Counter int

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    *ctr++
    fmt.Fprintf(w, "counter = %d\n", *ctr)
}
```

如果您的程序有一些内部状态需要通知页面已被访问怎么办？将频道绑定到网页。

```go
// A channel that sends a notification on each visit.
// (Probably want the channel to be buffered.)
type Chan chan *http.Request

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ch <- req
    fmt.Fprint(w, "notification sent")
}
```

最后，假设我们想在 /args 上显示调用服务器二进制文件时使用的参数。编写一个函数来打印参数很容易。

```go
func ArgServer() {
    fmt.Println(os.Args)
}
```

我们如何将其转变为 HTTP 服务器？我们可以使 ArgServer 成为某种类型的方法，我们忽略其值，但有一种更简洁的方法。由于我们可以为除指针和接口之外的任何类型定义方法，因此我们可以为函数编写方法。 http 包中包含以下代码：

```go
// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, req).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, req *Request) {
    f(w, req)
}
```

HandlerFunc 是一种带有 ServeHTTP 方法的类型，因此该类型的值可以为 HTTP 请求提供服务。看一下方法的实现：接收者是一个函数f，方法调用f。这可能看起来很奇怪，但它与接收器是一个通道以及在通道上发送的方法没有什么不同。
为了使 ArgServer 成为 HTTP 服务器，我们首先修改它以具有正确的签名。

```go
// Argument server.
func ArgServer(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(w, os.Args)
}
```

ArgServer 现在与 HandlerFunc 具有相同的签名，因此可以将其转换为该类型以访问其方法，就像我们将 Sequence 转换为 IntSlice 以访问 IntSlice.Sort 一样。设置它的代码很简洁：

```go
http.Handle("/args", http.HandlerFunc(ArgServer))
```

当有人访问页面 /args 时，安装在该页面的处理程序的值为 ArgServer 且类型为 HandlerFunc。 HTTP 服务器将调用该类型的 ServeHTTP 方法，以 ArgServer 作为接收者，后者将依次调用 ArgServer（通过 HandlerFunc.ServeHTTP 内的调用 f(w, req)）。然后将显示参数。

在本节中，我们从一个结构体、一个整数、一个通道和一个函数创建了一个 HTTP 服务器，所有这些都是因为接口只是方法集，可以为（几乎）任何类型定义。

# blank identifier

我们已经在 for range 循环和映射的上下文中多次提到了空白标识符。空白标识符可以分配或声明为任何类型的任何值，并且该值将被无害地丢弃。这有点像写入 Unix /dev/null 文件：它表示一个只写值，用作需要变量但实际值无关的占位符。它的用途超出了我们已经见过的用途。

## 多重赋值中的空白标识符

在 for range 循环中使用空白标识符是一般情况的特殊情况：多重赋值。

如果一项赋值操作需要在左侧有多个值，但程序不会使用其中一个值，则赋值操作左侧的空白标识符可以避免创建虚拟变量，并清楚地表明：该值将被丢弃。例如，当调用一个返回值和错误的函数时，但只有错误很重要，请使用空白标识符来丢弃不相关的值。

```go
if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Printf("%s does not exist\n", path)
}
```

有时您会看到丢弃错误值以忽略错误的代码；这是可怕的做法。经常检查错误返回；提供它们是有原因的。

```go
// Bad! This code will crash if path does not exist.
fi, _ := os.Stat(path)
if fi.IsDir() {
    fmt.Printf("%s is a directory\n", path)
}
```

## 未使用的导入和变量

导入包或声明变量而不使用它是错误的。未使用的导入会使程序膨胀并减慢编译速度，而已初始化但未使用的变量至少是浪费的计算，并且可能表明存在更大的错误。然而，当程序处于积极开发状态时，经常会出现未使用的导入和变量，并且为了继续编译而删除它们，然后又需要它们，这可能很烦人。空白标识符提供了一种解决方法。

这个写了一半的程序有两个未使用的导入（fmt 和 io）和一个未使用的变量（fd），因此它不会编译，但很高兴看看到目前为止的代码是否正确。

```go
package main

import (
    "fmt"
    "io"
    "log"
    "os"
)

func main() {
    fd, err := os.Open("test.go")
    if err != nil {
        log.Fatal(err)
    }
    // TODO: use fd.
}
```

要消除有关未使用导入的错误，请使用空白标识符来引用导入包中的符号。类似地，将未使用的变量 fd 分配给空白标识符将消除未使用的变量错误。该版本的程序确实可以编译。

```go
package main

import (
    "fmt"
    "io"
    "log"
    "os"
)

var _ = fmt.Printf // For debugging; delete when done.
var _ io.Reader    // For debugging; delete when done.

func main() {
    fd, err := os.Open("test.go")
    if err != nil {
        log.Fatal(err)
    }
    // TODO: use fd.
    _ = fd
}
```

按照惯例，用于消除导入错误的全局声明应该在导入之后立即出现并进行注释，这样既可以使它们易于查找，又可以提醒以后进行清理。

## Import for side effect

上一个示例中未使用的导入（如 fmt 或 io）最终应被使用或删除：空白分配将代码标识为正在进行的工作。但有时导入一个包只是为了它的副作用，而不需要任何显式的使用，这是有用的。例如，在其 init 函数期间，net/http/pprof 包会注册提供调试信息的 HTTP 处理程序。它有一个导出的 API，但大多数客户端只需要注册处理程序并通过网页访问数据。要仅导入包的副作用，请将包重命名为空白标识符：

```go
import _ "net/http/pprof"
```

这种导入形式清楚地表明，导入软件包是为了它的副作用，因为软件包没有其他可能的用途：在这个文件中，它没有名字。(如果它有名字，而我们不使用这个名字，编译器就会拒绝该程序）。

## 接口检查 Interface checks

正如我们在上面对接口的讨论中看到的，类型不需要显式声明它实现了接口。相反，类型仅通过实现接口的方法来实现接口。实际上，大多数接口转换都是静态的，因此在编译时进行检查。例如，将 *os.File 传递给需要 io.Reader 的函数将无法编译，除非 *os.File 实现 io.Reader 接口。

不过，一些接口检查确实在运行时发生。其中一个实例位于encoding/json 包中，它定义了Marshaler 接口。当 JSON 编码器接收到实现该接口的值时，编码器会调用该值的封送处理方法将其转换为 JSON，而不是执行标准转换。编码器在运行时使用如下类型断言检查此属性：

```go
m, ok := val.(json.Marshaler)
```

如果只需要询问类型是否实现接口，而不实际使用接口本身，也许作为错误检查的一部分，请使用空白标识符来忽略类型断言的值：

```go
if _, ok := val.(json.Marshaler); ok {
    fmt.Printf("value %v of type %T implements json.Marshaler\n", val, val)
}
```

出现这种情况的一个地方是当需要保证在实现该类型的包内它实际上满足接口时。如果某个类型（例如 json.RawMessage）需要自定义 JSON 表示形式，则它应该实现 json.Marshaler，但没有静态转换会导致编译器自动验证这一点。如果类型无意中无法满足接口，JSON 编码器仍然可以工作，但不会使用自定义实现。为了保证实现正确，可以在包中使用使用空白标识符的全局声明：

```go
var _ json.Marshaler = (*RawMessage)(nil)
```

在此声明中，涉及将 *RawMessage 转换为 Marshaler 的赋值要求 *RawMessage 实现 Marshaler，而该属性将在编译时进行检查。如果 json.Marshaler 接口发生变化，这个软件包将不再能编译，我们也会收到需要更新的通知。

此构造中出现空白标识符表明该声明仅用于类型检查，而不是用于创建变量。不过，不要对每种满足接口的类型都这样做。按照惯例，只有当代码中不存在静态转换时才使用此类声明，这种情况很少见。

# 嵌入 Embedding

Go 并不提供典型的、类型驱动的子类化概念，但它可以通过在结构或接口中嵌入类型来 "借用 "实现的片段。

接口嵌入非常简单。我们之前已经提到过 io.Reader 和 io.Writer 接口；这是它们的定义。

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

io 包还导出几个其他接口，这些接口指定可以实现多个此类方法的对象。例如，io.ReadWriter，一个包含读取和写入的接口。我们可以通过显式列出这两个方法来指定 io.ReadWriter，但嵌入这两个接口来形成新的接口会更容易、更容易引起共鸣，如下所示：

```go
// ReadWriter is the interface that combines the Reader and Writer interfaces.
type ReadWriter interface {
    Reader
    Writer
}
```

这就是它的样子：ReadWriter 可以做 Reader 和 Writer 所做的事情；它是嵌入式接口的联合。只有接口才能嵌入接口内。

同样的基本思想也适用于结构，但具有更深远的影响。 bufio 包有两种结构类型，bufio.Reader 和 bufio.Writer，当然每种类型都实现了 io 包中的类似接口。 bufio 还实现了一种缓冲读取器/写入器，它通过使用嵌入将读取器和写入器组合到一个结构中来实现：它列出了结构中的类型，但不给出字段名称。

```go
// ReadWriter stores pointers to a Reader and a Writer.
// It implements io.ReadWriter.
type ReadWriter struct {
    *Reader  // *bufio.Reader
    *Writer  // *bufio.Writer
}
```

内嵌元素是指向结构体的指针，当然在使用前必须初始化为指向有效的结构体。ReadWriter 结构可以写成

```go
type ReadWriter struct {
    reader *Reader
    writer *Writer
}
```

但是为了提升字段的方法并满足 io 接口，我们还需要提供转发方法，如下所示：

```go
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
    return rw.reader.Read(p)
}
```

通过直接嵌入结构，我们避免了这种簿记。嵌入类型的方法是免费提供的，这意味着 bufio.ReadWriter 不仅具有 bufio.Reader 和 bufio.Writer 的方法，还满足所有三个接口：io.Reader、io.Writer 和 io.ReadWriter。

嵌入与子类化有一个重要的区别。当我们嵌入一个类型时，该类型的方法成为外部类型的方法，但是当它们被调用时，该方法的接收者是内部类型，而不是外部类型。在我们的示例中，当调用 bufio.ReadWriter 的 Read 方法时，它与上面写的转发方法具有完全相同的效果；接收者是 ReadWriter 的读者字段，而不是 ReadWriter 本身。

嵌入也可以带来简单的便利。此示例显示了一个嵌入字段和一个常规命名字段。

```go
type Job struct {
    Command string
    *log.Logger
}
```

Job类型现在有*log.Logger的Print、Printf、Println等方法。当然，我们可以给 Logger 一个字段名称，但没有必要这样做。现在，初始化后，我们可以登录到作业：

```go
job.Println("starting now...")
```

Logger 是 Job 结构体的常规字段，因此我们可以在 Job 的构造函数中以通常的方式初始化它，如下所示：

```go
func NewJob(command string, logger *log.Logger) *Job {
    return &Job{command, logger}
}
```

或使用复合文字，

```go
job := &Job{command, log.New(os.Stderr, "Job: ", log.Ldate)}
```

如果我们需要直接引用嵌入字段，则该字段的类型名称（忽略包限定符）将用作字段名称，就像在 ReadWriter 结构体的 Read 方法中所做的那样。在这里，如果我们需要访问Job变量job的*log.Logger，我们会写job.Logger，如果我们想改进Logger的方法，这将很有用。

```go
func (job *Job) Printf(format string, args ...interface{}) {
    job.Logger.Printf("%q: %s", job.Command, fmt.Sprintf(format, args...))
}
```

嵌入类型引入了名称冲突的问题，但解决它们的规则很简单。首先，字段或方法 X 将任何其他项 X 隐藏在类型的更深层嵌套部分中。如果 log.Logger 包含一个名为 Command 的字段或方法，则 Job 的 Command 字段将支配它。

其次，如果相同的名称出现在同一嵌套级别，通常是一个错误；如果 Job 结构包含另一个名为 Logger 的字段或方法，则嵌入 log.Logger 将是错误的。但是，如果在类型定义之外的程序中从未提及过重复名称，则可以。此资格提供了一些保护，防止对外部嵌入的类型进行更改；如果添加的字段与另一个子类型中的另一个字段冲突（如果两个字段都没有使用过），那么没有问题。

# 并发 Concurrency

## Share by communicating

并发编程是一个很大的主题，这里仅介绍一些特定于 Go 的重点内容。

由于实现对共享变量的正确访问所需的微妙之处，许多环境中的并发编程变得困难。 Go 鼓励采用不同的方法，其中共享值在通道上传递，事实上，永远不会被单独的执行线程主动共享。在任何给定时间只有一个 goroutine 可以访问该值。按照设计，数据竞争不会发生。为了鼓励这种思维方式，我们将其简化为一句口号：

`不要通过共享内存进行通信；相反，通过通信来共享内存。`

这种方法可能太过分了。例如，最好通过在整数变量周围放置互斥体来完成引用计数。但作为一种高级方法，使用通道来控制访问可以更轻松地编写清晰、正确的程序。

考虑此模型的一种方法是考虑在一个 CPU 上运行的典型单线程程序。它不需要同步原语。现在运行另一个这样的实例；它也不需要同步。现在让这两个人交流；如果通信的是同步器，则仍然不需要其他同步。例如，Unix 管道就非常适合这个模型。尽管 Go 的并发方法起源于 Hoare 的通信顺序进程（CSP），但它也可以被视为 Unix 管道的类型安全泛化。

## Goroutines

它们之所以被称为 goroutines，是因为现有的术语--线程、coroutines、进程等--传达了不准确的内涵。goroutine 的模型很简单：它是一个在同一地址空间与其他 goroutine 并行执行的函数。它是轻量级的，只需分配堆栈空间。堆栈一开始很小，所以很便宜，然后通过分配（和释放）堆存储空间来增长。

Goroutines 被多路复用到多个操作系统线程上，因此如果其中一个线程发生阻塞（例如在等待 I/O 时），其他线程会继续运行。它们的设计隐藏了线程创建和管理的许多复杂性。

使用 go 关键字作为函数或方法调用的前缀，以在新的 goroutine 中运行该调用。当调用完成时，goroutine 会默默地退出。 （其效果类似于 Unix shell 中用于在后台运行命令的 & 表示法。）

```go
go list.Sort()  // run list.Sort concurrently; don't wait for it.
```

函数字面量在 goroutine 调用中非常方便。

```go
func Announce(message string, delay time.Duration) {
    go func() {
        time.Sleep(delay)
        fmt.Println(message)
    }()  // Note the parentheses - must call the function.
}
```

在 Go 中，函数字面是闭包：只要函数处于活动状态，其实现就会确保函数引用的变量继续存在。

这些示例不太实用，因为这些函数无法发出完成信号。为此，我们需要渠道。

## Channels

与映射一样，通道是通过 make 分配的，结果值充当对底层数据结构的引用。如果提供了可选的整数参数，它将设置通道的缓冲区大小。对于无缓冲或同步通道，默认值为零。

```go
ci := make(chan int)            // unbuffered channel of integers
cj := make(chan int, 0)         // unbuffered channel of integers
cs := make(chan *os.File, 100)  // buffered channel of pointers to Files
```

无缓冲通道将通信（值的交换）与同步相结合，保证两个计算（goroutine）处于已知状态。

有很多使用通道的好用语。这是让我们开始的一个。在上一节中，我们在后台启动了排序。通道可以允许启动 goroutine 等待排序完成。

```go
c := make(chan int)  // Allocate a channel.
// Start the sort in a goroutine; when it completes, signal on the channel.
go func() {
    list.Sort()
    c <- 1  // Send a signal; value does not matter.
}()
doSomethingForAWhile()
<-c   // Wait for sort to finish; discard sent value.
```

接收器总是阻塞，直到有数据要接收。如果通道未缓冲，则发送方会阻塞，直到接收方收到该值。如果通道有缓冲区，则发送方只会阻塞，直到该值被复制到缓冲区为止；如果缓冲区已满，则意味着等待某个接收器检索到值。

缓冲通道可以像信号量一样使用，例如限制吞吐量。在此示例中，传入请求被传递给句柄，该句柄将一个值发送到通道中，处理请求，然后从通道接收一个值，为下一个消费者准备好“信号量”。通道缓冲区的容量限制了同时调用处理的数量。

```go
var sem = make(chan int, MaxOutstanding)

func handle(r *Request) {
    sem <- 1    // Wait for active queue to drain.
    process(r)  // May take a long time.
    <-sem       // Done; enable next request to run.
}

func Serve(queue chan *Request) {
    for {
        req := <-queue
        go handle(req)  // Don't wait for handle to finish.
    }
}
```

一旦 MaxOutstanding 处理程序正在执行进程，更多处理程序将阻止尝试发送到已填充的通道缓冲区，直到现有处理程序之一完成并从缓冲区接收。

不过，这种设计有一个问题：Serve 为每个传入请求创建一个新的 goroutine，尽管其中只有 MaxOutstanding 可以随时运行。因此，如果请求传入得太快，程序可能会消耗无限的资源。我们可以通过改变 Serve 来控制 goroutine 的创建来解决这个缺陷。这是一个明显的解决方案，但请注意它有一个错误，我们随后将修复：

```go
func Serve(queue chan *Request) {
    for req := range queue {
        sem <- 1
        go func() {
            process(req) // Buggy; see explanation below.
            <-sem
        }()
    }
}
```

该错误在于，在 Go for 循环中，每次迭代都会重用循环变量，因此 req 变量在所有 goroutine 之间共享。那不是我们想要的。我们需要确保每个 goroutine 的 req 都是唯一的。这是一种方法，将 req 的值作为参数传递给 goroutine 中的闭包：

```go
func Serve(queue chan *Request) {
    for req := range queue {
        sem <- 1
        go func(req *Request) {
            process(req)
            <-sem
        }(req)
    }
}
```

将此版本与之前的版本进行比较，看看声明和运行闭包的方式有何不同。另一个解决方案是创建一个具有相同名称的新变量，如下例所示：

```go
func Serve(queue chan *Request) {
    for req := range queue {
        req := req // Create new instance of req for the goroutine.
        sem <- 1
        go func() {
            process(req)
            <-sem
        }()
    }
}
```

写起来可能看起来很奇怪

```go
req := req
```

但在 Go 中这样做是合法且惯用的。您将获得具有相同名称的变量的新版本，故意在本地隐藏循环变量，但对于每个 goroutine 都是唯一的。

回到编写服务器的一般问题，另一种很好地管理资源的方法是启动固定数量的句柄 goroutine，所有这些都从请求通道读取。 Goroutine 的数量限制了同时调用进程的数量。此服务功能还接受一个通道，在该通道上它将被告知退出；启动 goroutine 后，它会阻止从该通道接收数据。

```go
func handle(queue chan *Request) {
    for r := range queue {
        process(r)
    }
}

func Serve(clientRequests chan *Request, quit chan bool) {
    // Start handlers
    for i := 0; i < MaxOutstanding; i++ {
        go handle(clientRequests)
    }
    <-quit  // Wait to be told to exit.
}
```

## Channels of channels

Go 最重要的属性之一是通道是一流的值，可以像其他值一样进行分配和传递。此属性的常见用途是实现安全的并行解复用。

在上一节的示例中，handle 是请求的理想处理程序，但我们没有定义它正在处理的类型。如果该类型包含用于回复的通道，则每个客户端都可以提供自己的答案路径。这是 Request 类型的示意性定义。

```go
type Request struct {
    args        []int
    f           func([]int) int
    resultChan  chan int
}
```

客户端提供一个函数及其参数，以及请求对象内用于接收答案的通道。

```go
func sum(a []int) (s int) {
    for _, v := range a {
        s += v
    }
    return
}

request := &Request{[]int{3, 4, 5}, sum, make(chan int)}
// Send request
clientRequests <- request
// Wait for response.
fmt.Printf("answer: %d\n", <-request.resultChan)
```

在服务器端，处理程序函数是唯一发生变化的东西。

```go
func handle(queue chan *Request) {
    for req := range queue {
        req.resultChan <- req.f(req.args)
    }
}
```

显然，要使其现实，还有很多工作要做，但此代码是速率受限、并行、非阻塞 RPC 系统的框架，并且看不到互斥体。

## Parallelization

这些想法的另一个应用是跨多个 CPU 核心并行计算。如果计算可以分解为可以独立执行的单独部分，那么它就可以并行化，并在每个部分完成时通过一个通道发出信号。

假设我们要对项目向量执行一项昂贵的操作，并且每个项目的操作值是独立的，如这个理想化示例所示。

```go
type Vector []float64

// Apply the operation to v[i], v[i+1] ... up to v[n-1].
func (v Vector) DoSome(i, n int, u Vector, c chan int) {
    for ; i < n; i++ {
        v[i] += u.Op(v[i])
    }
    c <- 1    // signal that this piece is done
}
```

我们在循环中独立启动各个部分，每个 CPU 一个。它们可以按任何顺序完成，但这并不重要；我们只是在启动所有 goroutine 后清空通道来计算完成信号。

```go
const numCPU = 4 // number of CPU cores

func (v Vector) DoAll(u Vector) {
    c := make(chan int, numCPU)  // Buffering optional but sensible.
    for i := 0; i < numCPU; i++ {
        go v.DoSome(i*len(v)/numCPU, (i+1)*len(v)/numCPU, u, c)
    }
    // Drain the channel.
    for i := 0; i < numCPU; i++ {
        <-c    // wait for one task to complete
    }
    // All done.
}
```

与其为 numCPU 创建一个常量值，我们还不如询问运行时哪个值合适。函数 runtime.NumCPU 返回机器中硬件 CPU 内核的数量，因此我们可以写道

```go
var numCPU = runtime.NumCPU()
```

还有一个函数runtime.GOMAXPROCS，它报告（或设置）Go程序可以同时运行的用户指定的内核数量。它默认为runtime.NumCPU 的值，但可以通过设置类似命名的shell 环境变量或使用正数调用函数来覆盖。用零调用它只是查询该值。因此，如果我们想满足用户的资源请求，我们应该写

```go
var numCPU = runtime.GOMAXPROCS(0)
```

请务必不要混淆并发性（将程序构建为独立执行的组件）和并行性（在多个 CPU 上并行执行计算以提高效率）的概念。虽然Go的并发特性可以使一些问题很容易构建为并行计算，但Go是一种并发语言，而不是并行语言，并且并非所有并行化问题都适合Go的模型。有关区别的讨论，请参阅本博客文章 [this blog post](https://blog.golang.org/2013/01/concurrency-is-not-parallelism.html)中引用的谈话。

## A leaky buffer

并发编程的工具甚至可以使非并发的想法更容易表达。下面是一个从 RPC 包中抽象出来的示例。客户端 goroutine 循环从某个源（可能是网络）接收数据。为了避免分配和释放缓冲区，它保留一个空闲列表，并使用缓冲通道来表示它。如果通道为空，则会分配新的缓冲区。一旦消息缓冲区准备就绪，它就会通过 serverChan 发送到服务器。

```go
var freeList = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)

func client() {
    for {
        var b *Buffer
        // Grab a buffer if available; allocate if not.
        select {
        case b = <-freeList:
            // Got one; nothing more to do.
        default:
            // None free, so allocate a new one.
            b = new(Buffer)
        }
        load(b)              // Read next message from the net.
        serverChan <- b      // Send to server.
    }
}
```

服务器循环接收来自客户端的每条消息，对其进行处理，并将缓冲区返回到空闲列表。

```go
func server() {
    for {
        b := <-serverChan    // Wait for work.
        process(b)
        // Reuse buffer if there's room.
        select {
        case freeList <- b:
            // Buffer on free list; nothing more to do.
        default:
            // Free list full, just carry on.
        }
    }
}
```

客户端尝试从 freeList 中检索缓冲区；如果没有可用的，它会分配一个新的。服务器发送到 freeList 会将 b 放回到空闲列表中，除非列表已满，在这种情况下，缓冲区将被丢弃到地板上以供垃圾收集器回收。 （当没有其他情况准备好时，select 语句中的默认子句就会执行，这意味着 select 永远不会阻塞。）此实现仅用几行就构建了一个漏桶空闲列表，依赖于缓冲通道和垃圾收集器进行簿记。

# Errors

库例程通常必须向调用者返回某种错误指示。如前所述，Go 的多值返回可以轻松地在正常返回值旁边返回详细的错误描述。使用此功能提供详细的错误信息是一种很好的风格。例如，正如我们将看到的，os.Open 不仅在失败时返回一个 nil 指针，它还返回一个描述错误的错误值。

按照惯例，错误具有错误类型，一个简单的内置接口。

```go
type error interface {
    Error() string
}
```

库编写者可以自由地使用更丰富的模型来实现此接口，这样不仅可以看到错误，还可以提供一些上下文。如前所述，除了通常的 *os.File 返回值之外，os.Open 还返回一个错误值。如果文件打开成功，错误将为nil，但是当出现问题时，它会持有一个os.PathError：

```go
// PathError records an error and the operation and
// file path that caused it.
type PathError struct {
    Op string    // "open", "unlink", etc.
    Path string  // The associated file.
    Err error    // Returned by the system call.
}

func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
}
```

PathError 的 Error 生成如下字符串：

```go
open /etc/passwx: no such file or directory
```

此类错误（包括有问题的文件名、操作及其触发的操作系统错误）很有用，即使打印的内容与导致该错误的调用相距甚远；它比简单的“没有这样的文件或目录”提供更多信息。
如果可行，错误字符串应标识其来源，例如通过使用前缀来命名生成错误的操作或包。例如，在包图像中，由于未知格式而导致的解码错误的字符串表示为“image:unknown format”。
关心精确错误详细信息的调用者可以使用类型开关或类型断言来查找特定错误并提取详细信息。对于 PathErrors，这可能包括检查内部 Err 字段是否存在可恢复的故障。

```go
for try := 0; try < 2; try++ {
    file, err = os.Create(filename)
    if err == nil {
        return
    }
    if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
        deleteTempFiles()  // Recover some space.
        continue
    }
    return
}
```

这里的第二个 if 语句是另一种类型断言。如果失败，ok 将为 false，e 将为 nil。如果成功，ok 将为 true，这意味着错误的类型为 *os.PathError，然后 e 也是类型，我们可以检查它以获取有关错误的更多信息。

## Panic

向调用者报告错误的常用方法是将错误作为额外的返回值返回。规范的 Read 方法是一个众所周知的实例；它返回字节计数和错误。但如果错误无法恢复怎么办？有时程序根本无法继续。

为此，有一个内置函数恐慌，它实际上会创建一个运行时错误，从而停止程序（但请参阅下一节）。该函数采用任意类型的单个参数（通常是字符串）在程序终止时打印。这也是表明发生了不可能的事情的一种方式，例如退出无限循环。

```go
// A toy implementation of cube root using Newton's method.
func CubeRoot(x float64) float64 {
    z := x/3   // Arbitrary initial value
    for i := 0; i < 1e6; i++ {
        prevz := z
        z -= (z*z*z-x) / (3*z*z)
        if veryClose(z, prevz) {
            return z
        }
    }
    // A million iterations has not converged; something is wrong.
    panic(fmt.Sprintf("CubeRoot(%g) did not converge", x))
}
```

这只是一个示例，但真正的库函数应该避免恐慌。如果问题可以被掩盖或解决，那么让事情继续运行总是比取消整个程序更好。一个可能的反例是在初始化期间：如果库确实无法自行设置，可以这么说，恐慌可能是合理的。

```go
var user = os.Getenv("USER")

func init() {
    if user == "" {
        panic("no value for $USER")
    }
}
```

## Recover

当恐慌被调用时，包括隐式的运行时错误，例如索引切片越界或类型断言失败，它会立即停止当前函数的执行，并开始展开 goroutine 的堆栈，同时运行任何延迟函数。如果展开到达 goroutine 堆栈的顶部，程序就会终止。但是，可以使用内置函数recover来重新获得对goroutine的控制并恢复正常执行。

对recover的调用会停止展开并返回传递给panic的参数。由于展开时运行的唯一代码位于延迟函数内部，因此恢复仅在延迟函数内部有用。

恢复的一个应用是关闭服务器内发生故障的 goroutine，而不杀死其他正在执行的 goroutine。

```go
func server(workChan <-chan *Work) {
    for work := range workChan {
        go safelyDo(work)
    }
}

func safelyDo(work *Work) {
    defer func() {
        if err := recover(); err != nil {
            log.Println("work failed:", err)
        }
    }()
    do(work)
}
```

在这个例子中，如果 do(work) 出现恐慌，结果将被记录下来，并且 goroutine 将干净地退出，而不会干扰其他 goroutine。延迟关闭中无需执行任何其他操作；调用recover 可以完全处理该情况。

因为除非直接从延迟函数调用，否则recover总是返回nil，所以延迟代码可以调用本身使用panic和recover的库例程而不会失败。例如，safetyDo 中的延迟函数可能会在调用恢复之前调用日志记录函数，并且该日志记录代码将不受恐慌状态影响而运行。

有了我们的恢复模式，do 函数（以及它调用的任何函数）就可以通过调用panic 来彻底摆脱任何糟糕的情况。我们可以利用这个想法来简化复杂软件中的错误处理。让我们看一下 regexp 包的理想化版本，它通过使用本地错误类型调用恐慌来报告解析错误。下面是 Error 的定义、错误方法和 Compile 函数。

```go
// Error is the type of a parse error; it satisfies the error interface.
type Error string
func (e Error) Error() string {
    return string(e)
}

// error is a method of *Regexp that reports parsing errors by
// panicking with an Error.
func (regexp *Regexp) error(err string) {
    panic(Error(err))
}

// Compile returns a parsed representation of the regular expression.
func Compile(str string) (regexp *Regexp, err error) {
    regexp = new(Regexp)
    // doParse will panic if there is a parse error.
    defer func() {
        if e := recover(); e != nil {
            regexp = nil    // Clear return value.
            err = e.(Error) // Will re-panic if not a parse error.
        }
    }()
    return regexp.doParse(str), nil
}
```

如果 doParse 发生混乱，恢复块会将返回值设置为 nil——延迟函数可以修改命名返回值。然后，它会在对 err 的赋值中通过断言它具有本地类型 Error 来检查问题是否是解析错误。如果没有，类型断言将失败，导致运行时错误，继续堆栈展开，就好像没有任何东西中断它一样。此检查意味着如果发生意外情况，例如索引越界，即使我们使用恐慌和恢复来处理解析错误，代码也会失败。

错误处理到位后，错误方法（因为它是绑定到类型的方法，所以它与内置错误类型具有相同的名称很好，甚至很自然）可以轻松报告解析错误，而无需担心展开手动解析堆栈：

```go
if pos == 0 {
    re.error("'*' illegal at start of expression")
}
```

尽管此模式很有用，但它应该仅在包内使用。 Parse 将其内部恐慌调用转换为错误值；它不会向客户端暴露恐慌。这是一个值得遵循的好规则。

顺便说一句，如果发生实际错误，这种重新恐慌习惯用法会更改恐慌值。但是，原始故障和新故障都会出现在崩溃报告中，因此问题的根本原因仍然可见。因此，这种简单的重新恐慌方法通常就足够了——毕竟这是一次崩溃——但如果您只想显示原始值，则可以编写更多代码来过滤意外问题并使用原始错误重新恐慌。这留给读者作为练习。

# A web server

让我们看一个完整的 Go 程序（一个 Web 服务器）。这实际上是一种网络重新服务器。 Google 在 Chart.apis.google.com 上提供了一项服务，可以自动将数据格式化为图表和图形。不过，它很难以交互方式使用，因为您需要将数据作为查询放入 URL 中。这里的程序为一种数据形式提供了一个更好的界面：给定一小段文本，它调用图表服务器生成一个 QR 码，即对文本进行编码的框矩阵。该图像可以用手机的摄像头抓取并解释为 URL，这样您就无需在手机的小键盘中输入 URL。

下面是完整的程序。以下是解释。

```go
package main

import (
    "flag"
    "html/template"
    "log"
    "net/http"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
    flag.Parse()
    http.Handle("/", http.HandlerFunc(QR))
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func QR(w http.ResponseWriter, req *http.Request) {
    templ.Execute(w, req.FormValue("s"))
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET">
    <input maxLength=1024 size=70 name=s value="" title="Text to QR Encode">
    <input type=submit value="Show QR" name=qr>
</form>
</body>
</html>
```

直到主要部分应该很容易理解。 one 标志为我们的服务器设置默认 HTTP 端口。模板变量 templ 是有趣的地方。它构建一个 HTML 模板，服务器将执行该模板来显示页面；稍后会详细介绍。
main 函数解析标志，并使用我们上面讨论的机制将函数 QR 绑定到服务器的根路径。然后调用http.ListenAndServe来启动服务器；它在服务器运行时阻塞。
QR 只是接收包含表单数据的请求，并对名为 s 的表单值中的数据执行模板。
模板包html/template功能强大；这个程序只是触及了它的功能。本质上，它通过替换从传递给 templ.Execute 的数据项派生的元素（在本例中为表单值）来动态重写一段 HTML 文本。在模板文本 (templateStr) 中，双括号分隔的部分表示模板操作。仅当当前数据项的值称为 时，从 {{if .}} 到 {{end}} 的部分才会执行。 （点），非空。也就是说，当字符串为空时，这一段模板被抑制。
这两个片段 {{.}} 表示在网页上显示提供给模板的数据（查询字符串）。 HTML 模板包自动提供适当的转义，以便文本可以安全显示。
模板字符串的其余部分只是页面加载时显示的 HTML。如果这个解释太快，请参阅模板包的文档以进行更彻底的讨论。
现在您已经得到了：一个有用的 Web 服务器，由几行代码加上一些数据驱动的 HTML 文本组成。 Go 足够强大，可以用几行代码完成很多事情。
