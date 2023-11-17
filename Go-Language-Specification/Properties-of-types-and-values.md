# 类型和值的属性 Properties of types and values

## 底层类型 Underlying types

每个类型 T 都有一个底层类型( *underlying type*)：如果 T 是预声明的布尔、数字或字符串类型之一，或者类型字面量，则相应的底层类型是 T 本身。否则，T 的底层类型是 T 在其声明中引用的类型的底层类型。对于作为其类型约束( [type constraint](https://go.dev/ref/spec#Type_constraints))的底层类型的类型形参，它始终是一个接口。

```go
type (
    A1 = string
    A2 = A1
)

type (
    B1 string
    B2 B1
    B3 []B1
    B4 B3
)

func f[P any](x P) { … }
```

字符串 A1、A2、B1 和 B2 的底层类型是字符串。[]B1、B3 和 B4 的底层类型是 []B1。P 的基础类型是 interface{}。

## 核心类型 Core types

每个非接口类型 T 都有一个核心类型，该核心类型与 T 的底层类型相同。

如果满足以下条件之一，则接口 T 具有核心类型：

1. 存在一个单一类型 U，它是 T 的类型集中所有类型的底层类型；或

2. T的类型集仅包含具有相同元素类型E的通道类型，并且所有有向通道具有相同的方向。

其他接口不具有核心类型。

根据满足的条件，接口的核心类型是：

1. 类型 U；或

2. 如果 T 只包含双向通道，则是 chan E 类型；如果 T 包含定向通道，则是 chan<- E 或 <-chan E 类型。

根据定义，核心类型绝不是已定义的类型、类型参数或接口类型。

具有核心类型的接口示例：

```go
type Celsius float32
type Kelvin  float32

interface{ int }                          // int
interface{ Celsius|Kelvin }               // float32
interface{ ~chan int }                    // chan int
interface{ ~chan int|~chan<- int }        // chan<- int
interface{ ~[]*data; String() string }    // []*data
```

无核心类型的接口示例：

```go
interface{}                               // no single underlying type
interface{ Celsius|float64 }              // no single underlying type
interface{ chan int | chan<- string }     // channels have different element types
interface{ <-chan int | chan<- int }      // directional channels have different directions
```

某些操作（切片表达式、附加和复制）依赖于接受字节切片和字符串的稍微宽松的核心类型形式。具体来说，如果恰好有[]byte和string两种类型，这两种类型是接口T的类型集中所有类型的底层类型，则T的核心类型称为bytestring。

[注意： bytestring 不是真正的类型；它不能用于声明变量由其他类型组成。它的存在只是为了描述某些从字节序列读取的操作的行为，这些字节序列可能是字节片或字符串。]

一些核心类型为bytestring的接口案例：

```go
interface{ int }                          // int (same as ordinary core type)
interface{ []byte | string }              // bytestring
interface{ ~[]byte | myString }           // bytestring
```

## 类型标识 Type identity

<mark>两种类型要么相同，要么不同</mark>。

<mark>已命名的类型总是不同于其他任何类型</mark>。否则，如果两个类型的底层类型字面结构相同，即它们具有相同的字面结构，相应的组件具有相同的类型，那么这两个类型就是相同的。详细说明：

- 如果两个数组类型具有<mark>相同的元素类型和相同的数组长度</mark>，则它们是相同的。

- 如果两个切片类型具有<mark>相同的元素类型</mark>，则它们是相同的。

- 如果两个结构类型具有<mark>相同的字段序列，并且相应的字段具有相同的名称、相同的类型和相同的标记</mark>，则它们是相同的。不同包中的非导出字段名称总是不同的。

- 如果两个指针类型<mark>具有相同的基类型</mark>，则它们是相同的。

- 如果两个<mark>函数类型的参数和结果值数量相同，相应的参数和结果类型相同，并且两个函数都是可变函数或两个函数都不是可变函数</mark>，那么这两个函数类型就是相同的。参数和结果名称不需要匹配。

- 如果两个接口类型定义了<mark>相同的类型集</mark>，那么它们就是相同的。

- 如果两个映射类型具有<mark>相同的键和元素类型</mark>，则它们是相同的。

- 如果两个通道类型具有<mark>相同的元素类型和相同的方向</mark>，则它们是相同的。

- 如果<mark>两个实例化类型的定义类型和所有类型参数都相同</mark>，则它们是相同的。

鉴于声明：

```go
type (
    A0 = []string
    A1 = A0
    A2 = struct{ a, b int }
    A3 = int
    A4 = func(A3, float64) *A0
    A5 = func(x int, _ float64) *[]string

    B0 A0
    B1 []string
    B2 struct{ a, b int }
    B3 struct{ a, c int }
    B4 func(int, float64) *B0
    B5 func(x int, y float64) *A1

    C0 = B0
    D0[P1, P2 any] struct{ x P1; y P2 }
    E0 = D0[int, string]
)
```

这些类型是相同的：

```textile
A0, A1, and []string
A2 and struct{ a, b int }
A3 and int
A4, func(int, float64) *[]string, and A5

B0 and C0
D0[int, string] and E0
[]int and []int
struct{ a, b *B5 } and struct{ a, b *B5 }
func(x int, y float64) *[]string, func(int, float64) (result *[]string), and A5
```

B0 和 B1 是不同的，因为它们是由不同类型定义创建的新​​类型<sup>1</sup>； func(int, float64) *B0 和 func(x int, y float64) *[]string 不同，因为 B0 与 []string 不同<sup>2</sup>； P1和P2是不同的，因为它们是不同的类型参数<sup>3</sup>。 D0[int, string] 和 struct{ x int; y string } 是不同的，因为前者是实例化的定义类型，而后者是类型文字（但它们仍然是可赋值的）<sup>4</sup>。

## 可赋值性 Assignability

如果满足以下条件之一，则 V 类型的值 x 可分配给 T 类型的变量（“x 可分配给 T”）：

- V 和 T 类型相同。

- V 和 T 具有相同的底层类型，但不是类型形参，并且 V 或 T 中至少之一不是命名类型。

- V和T是具有相同元素类型的通道类型，V是双向通道，并且V或T中至少一个不是命名类型。

- T 是接口类型，但不是类型参数，并且 x 实现 T。

- x 是预声明的标识符 nil，T 是指针、函数、切片、映射、通道或接口类型，但不是类型参数。

- x 是一个无类型常量，可由类型 T 的值表示。

此外，如果 x 的类型 V 或 T 是类型参数，并且满足以下条件之一，则 x 可分配给类型 T 的变量：

- x 是预先声明的标识符 nil，T 是类型参数，并且 x 可分配给 T 类型集中的每个类型。

- V 不是命名类型，T 是类型参数，x 可分配给 T 类型集中的每个类型。

- V 是类型参数，T 不是命名类型，V 类型集中的每个类型的值都可以分配给 T。

## 可表示性 Representability

如果满足以下条件之一，则常量 x 可由类型 T 的值表示，其中 T 不是类型参数：

- 值x 位于由类型 T 确定的值集合中。

- T 是浮点类型，x 可以舍入到 T 的精度而不会溢出。舍入使用 IEEE 754 舍入到偶数规则，但 IEEE 负零进一步简化为无符号零。请注意，常量值永远不会导致 IEEE 负零、NaN 或无穷大。

- T 是复数类型，x 的分量 real(x) 和 imag(x) 可用 T 的组件类型（float32 或 float64）的值表示。

如果 T 是类型参数，则 x 可由类型 T 的值表示（如果 x 可由 T 的类型集中的每种类型的值表示）。

```go
x                   T           x is representable by a value of T because

'a'                 byte        97 is in the set of byte values
97                  rune        rune is an alias for int32, and 97 is in the set of 32-bit integers
"foo"               string      "foo" is in the set of string values
1024                int16       1024 is in the set of 16-bit integers
42.0                byte        42 is in the set of unsigned 8-bit integers
1e10                uint64      10000000000 is in the set of unsigned 64-bit integers
2.718281828459045   float32     2.718281828459045 rounds to 2.7182817 which is in the set of float32 values
-1e-1000            float64     -1e-1000 rounds to IEEE -0.0 which is further simplified to 0.0
0i                  int         0 is an integer value
(42 + 0i)           float32     42.0 (with zero imaginary part) is in the set of float32 values
x                   T           x is not representable by a value of T because

0                   bool        0 is not in the set of boolean values
'a'                 string      'a' is a rune, it is not in the set of string values
1024                byte        1024 is not in the set of unsigned 8-bit integers
-1                  uint16      -1 is not in the set of unsigned 16-bit integers
1.1                 int         1.1 is not an integer value
42i                 float32     (0 + 42i) is not in the set of float32 values
1e1000              float64     1e1000 overflows to IEEE +Inf after rounding
```

## 方法集 Method sets

类型的方法集决定了可以对该类型的操作数调用的方法。每个类型都有一个与其关联的（可能是空的）方法集：

- 定义类型 T 的方法集由使用接收者类型 T 声明的所有方法组成。

- 定义类型 T 的指针（T 既不是指针也不是接口）的方法集是以接收器 *T 或 T 声明的所有方法的集合。

- 接口类型的方法集是接口类型集中每个类型的方法集的交集（得到的方法集通常只是接口中已声明的方法集）。

进一步的规则适用于包含嵌入字段的结构体（和结构体指针），详见结构体类型一节。任何其他类型的方法集都是空的。

在方法集中，每个方法都必须有一个唯一的非空白方法名。
