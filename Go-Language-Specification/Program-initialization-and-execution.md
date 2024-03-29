# 程序初始化与执行 Program initialization and execution

## 零值 The zero value

当通过声明或调用 new 为变量分配存储空间时，或者当通过复合字面或调用 make 创建新值时，如果没有提供显式初始化，变量或值将被赋予默认值。这种变量或值的每个元素都被设置为其类型的零值：布尔类型为 false，数字类型为 0，字符串为""，指针、函数、接口、片段、通道和映射为 nil。这种初始化是递归进行的，例如，如果没有指定值，结构体数组的每个元素都将被清零。

这两个简单的声明是等效的：

```go
var i int
var i int = 0
```

然后

```go
type T struct { i int; f float64; next *T }
t := new(T)
```

以下成立：

```go
t.i == 0
t.f == 0.0
t.next == nil
```

之后也是如此

```go
var t T
```

## 包初始化 Package initialzation

在包内，包级变量初始化逐步进行，每一步都选择声明顺序中最早的变量，该变量不依赖于未初始化的变量。

更准确地说，如果包级变量尚未初始化并且没有初始化表达式或其初始化表达式不依赖于未初始化的变量，则认为该变量已准备好初始化。通过重复初始化声明顺序中最早并准备初始化的下一个包级变量来进行初始化，直到没有准备好初始化的变量为止。

如果在此过程结束时任何变量仍未初始化，则这些变量是一个或多个初始化周期的一部分，并且程序无效。

由右侧单个（多值）表达式初始化的变量声明左侧的多个变量会一起初始化：如果左侧的任何变量被初始化，则所有这些变量都会被初始化在同一步骤中。

```go
var x = a
var a, b = f() // a and b are initialized together, before x is initialized
```

出于包初始化的目的，空白变量被视为与声明中的任何其他变量一样。

在多个文件中声明的变量的声明顺序由文件呈现给编译器的顺序决定：第一个文件中声明的变量先于第二个文件中声明的任何变量声明，依此类推。为了确保可重复的初始化行为，鼓励构建系统以词法文件名顺序向编译器呈现属于同一包的多个文件。

依赖性分析不依赖于变量的实际值，仅依赖于源中对它们的词汇引用，并进行传递分析。例如，如果变量 x 的初始化表达式引用了一个函数，而该函数的主体引用了变量 y，则 x 依赖于 y。具体来说：

- 对变量或函数的引用是表示该变量或函数的标识符。

- 对方法 m 的引用是 t.m 形式的方法值或方法表达式，其中 t 的（静态）类型不是接口类型，并且方法 m 在 t 的方法集中。是否调用结果函数值 t.m 并不重要。

- 如果 x 的初始化表达式或主体（对于函数和方法）包含对 y 或依赖于 y 的函数或方法的引用，则变量、函数或方法 x 依赖于变量 y。

例如，给定声明

```go
var (
    a = c + b  // == 9
    b = f()    // == 4
    c = f()    // == 5
    d = 3      // == 5 after initialization has finished
)

func f() int {
    d++
    return d
}
```

初始化顺序为d、b、c、a。请注意，初始化表达式中子表达式的顺序无关：在此示例中，a = c + b 和 a = b + c 会产生相同的初始化顺序。

依赖性分析是按软件包进行的；只有当前软件包中声明的变量、函数和（非接口）方法的引用才会被考虑。如果变量之间存在其他隐藏的数据依赖关系，则不指定这些变量之间的初始化顺序。

例如，给定声明

```go
var x = I(T{}).ab()   // x has an undetected, hidden dependency on a and b
var _ = sideEffect()  // unrelated to x, a, or b
var a = b
var b = 42

type I interface      { ab() []int }
type T struct{}
func (T) ab() []int   { return []int{a, b} }
```

变量 a 将在 b 之后初始化，但 x 是在 b 之前、b 和 a 之间还是 a 之后初始化，因此也未指定调用 sideEffect() 的时刻（在 x 初始化之前或之后）。

变量也可以使用在包块中声明的名为 init 的函数来初始化，不带任何参数，也没有结果参数。

每个包可以定义多个这样的函数，甚至在单个源文件中也是如此。在 package 块中，init 标识符只能用于声明 init 函数，而标识符本身并未声明。因此，不能从程序中的任何地方引用 init 函数。
通过为所有包级变量分配初始值来初始化整个包，然后按照它们在源中（可能在多个文件中）出现的顺序调用所有 init 函数，如呈现给编译器的那样。

## 程序初始化 Program initialization

一个完整程序的包是逐步初始化的，一次一个包。如果包具有导入，则导入的包会在初始化包本身之前初始化。如果多个包导入一个包，导入的包只会被初始化一次。通过构造导入包可以保证不存在循环初始化依赖关系。更确切地说：


给定按导入路径排序的所有包的列表，在每个步骤中，所有导入的包（如果有）都已初始化的列表中的第一个未初始化的包将被初始化。重复此步骤直到所有包都被初始化。


包初始化（变量初始化和 init 函数的调用）发生在单个 goroutine 中，按顺序一次一个包。 init 函数可以启动其他 goroutine，这些 goroutine 可以与初始化代码同时运行。然而，初始化总是对 init 函数进行排序：在前一个函数返回之前，它不会调用下一个函数。

## 程序执行 Program execution

一个完整的程序是通过将一个称为main包的未导入的包与它导入的所有包以传递方式链接起来创建的。主包必须具有包名称 main 并声明一个不带参数且不返回值的函数 main。

```go
func main() { … }
```

程序执行首先初始化程序，然后调用 main 包中的 main 函数。当该函数调用返回时，程序退出。它不会等待其他（非主）goroutine 完成。
