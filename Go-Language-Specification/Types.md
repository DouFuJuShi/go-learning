# 类型 Types

类型确定一组值，以及针对于这些值的操作和方法。类型可以由类型名称（如果有的话）表示，如果泛型类型的，则后面必须跟有类型参数[type arguments](https://go.dev/ref/spec#Instantiations) 。也可以使用类型字面量指定类型，该类型字面量由现有类型组成。

```ebnf
Type      = TypeName [ TypeArgs ] | TypeLit | "(" Type ")" .
TypeName  = identifier | QualifiedIdent .
TypeArgs  = "[" TypeList [ "," ] "]" .
TypeList  = Type { "," Type } .
TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
            SliceType | MapType | ChannelType .
```

该语言预先声明了某些类型名称。其他的则通过类型声明([type declarations](https://go.dev/ref/spec#Type_declarations))或类型参数列表([type parameter lists](https://go.dev/ref/spec#Type_parameter_declarations))来引入。复合类型 <mark>*Composite types*</mark>（数组、结构体、指针、函数、接口、切片、映射和通道类型）可以使用类型文字<mark>*type literal*</mark>来构造。

<mark>预声明类型、定义类型和类型参数称为命名类型(Named Types)。如果别名声明中给定的类型是命名类型，则别名表示命名类型。</mark>

## Boolean types

布尔类型表示由预先声明的常量 true 和 false 表示的布尔真值集。预先声明的布尔类型为bool；它是一个定义的类型。

## Numeric types

整数、浮点或复数类型分别表示整数、浮点或复数值的集合。它们统称为数字类型。预先声明的独立于体系结构的数字类型是：

```ebnf
uint8       the set of all unsigned  8-bit integers (0 to 255)
uint16      the set of all unsigned 16-bit integers (0 to 65535)
uint32      the set of all unsigned 32-bit integers (0 to 4294967295)
uint64      the set of all unsigned 64-bit integers (0 to 18446744073709551615)

int8        the set of all signed  8-bit integers (-128 to 127)
int16       the set of all signed 16-bit integers (-32768 to 32767)
int32       the set of all signed 32-bit integers (-2147483648 to 2147483647)
int64       the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807)

float32     the set of all IEEE-754 32-bit floating-point numbers
float64     the set of all IEEE-754 64-bit floating-point numbers

complex64   the set of all complex numbers with float32 real and imaginary parts
complex128  the set of all complex numbers with float64 real and imaginary parts

byte        alias for uint8
rune        alias for int32
```

n 位整数的值是 n 位宽，并使用二进制补码算术([two's complement arithmetic](https://en.wikipedia.org/wiki/Two's_complement))表示。
还有一组具有特定于实现的大小的预先声明的整数类型：

```textile
uint     either 32 or 64 bits
int      same size as uint
uintptr  an unsigned integer large enough to store the uninterpreted bits of a pointer value
```

为了避免可移植性问题，所有数字类型都是已定义的类型，因此除了 byte（uint8 的别名）和 rune（int32 的别名）之外都是不同的。当不同的数值类型混合在表达式或赋值中时，需要显式转换。例如，int32 和 int 不是同一类型，尽管它们在特定体系结构上可能具有相同的大小。

## String types

字符串类型表示字符串值的集合。字符串值是一个（可能为空）字节序列。字节数称为字符串的长度，并且永远不会是负数。字符串是不可变的：一旦创建，就不可能更改字符串的内容。预声明的字符串类型为string；它是一个定义的类型。
可以使用内置函数 len 计算字符串 s 的长度。如果字符串是常量，则长度是编译时常量。字符串的字节可以通过整数索引 0 到 len(s)-1 来访问。获取此类元素的地址是非法的；如果 s[i] 是字符串的第 i 个字节，则 &s[i] 无效。

## Array types

数组是单一类型元素的编号序列，称为元素类型。元素的数量称为数组的长度，并且永远不会是负数。

```ebnf
ArrayType   = "[" ArrayLength "]" ElementType .
ArrayLength = Expression .
ElementType = Type .
```

长度是数组类型的一部分；它必须计算为由 int 类型值表示的非负常量。可以使用内置函数 len 来确定数组 a 的长度。这些元素可以通过整数索引 0 到 len(a)-1 来寻址。数组类型始终是一维的，但可以组合形成多维类型。

```go
[32]byte
[2*N] struct { x, y int32 }
[1000]*float64
[3][5]int
[2][2][2]float64  // same as [2]([2]([2]float64))
```

数组类型 T 可能不具有类型 T 的元素，或者直接或间接包含 T 作为组件的类型的元素（如果包含类型的元素仅为数组或结构类型）。

```go
// invalid array types
type (
    T1 [10]T1                 // element type of T1 is T1
    T2 [10]struct{ f T2 }     // T2 contains T2 as component of a struct
    T3 [10]T4                 // T3 contains T3 as component of a struct in T4
    T4 struct{ f T3 }         // T4 contains T4 as component of array T3 in a struct
)

// valid array types
type (
    T5 [10]*T5                // T5 contains T5 as component of a pointer
    T6 [10]func() T6          // T6 contains T6 as component of a function type
    T7 [10]struct{ f []T7 }   // T7 contains T7 as component of a slice in a struct
)
```

## Slice types

切片是底层数组中连续片段的描述符，用于访问该数组中编号的元素序列。切片类型表示其元素类型的所有数组切片的集合。元素的个数称为片段的长度，永远不会是负数。未初始化的片段的值为 nil。

```ebnf
SliceType = "[" "]" ElementType .
```

切片 s 的长度可以通过内置函数 len 得知；与数组不同的是，它在执行过程中可能会发生变化。这些元素可以通过整数索引 0 到 len(s)-1 来寻址。给定元素的切片索引可能小于底层数组中相同元素的索引。

切片一旦初始化，就始终与保存其元素的底层数组相关联。因此，切片与其数组以及同一数组的其他切片共享存储；相比之下，不同的数组总是代表不同的存储。

切片下面的数组可以延伸超过切片的末尾。容量是该范围的度量：它是切片长度与切片之外的数组长度之和；可以通过从原始切片中切出新切片来创建长度达到该容量的切片。可以使用内置函数 cap(a) 来发现切片 a 的容量。
可以使用内置函数 make 为给定元素类型 T 生成新的初始化切片值，该函数采用切片类型和指定长度和可选容量的参数。使用 make 创建的切片总是分配一个新的隐藏数组，返回的切片值引用该数组。

```go
make([]T, length, capacity)
```

生成与分配数组并对其进行切片相同的切片，因此这两个表达式是等效的：

```go
make([]int, 50, 100)
new([100]int)[0:50]
```

与数组一样，切片始终是一维的，但可以组成更高维的对象。在数组的数组中，内部数组的长度在构造上总是相同的；但在切片的切片（或片的数组）中，内部长度可能会动态变化。此外，内部切片必须单独初始化。

## Struct Types

结构体是命名元素（称为字段）的序列，每个元素都有名称和类型。字段名称可以显式指定（IdentifierList）或隐式指定（EmbeddedField）。在结构体中，非空白字段名称必须是唯一的。

```ebnf
StructType    = "struct" "{" { FieldDecl ";" } "}" .
FieldDecl     = (IdentifierList Type | EmbeddedField) [ Tag ] .
EmbeddedField = [ "*" ] TypeName [ TypeArgs ] .
Tag           = string_lit .
```

```go
// An empty struct.
struct {}

// A struct with 6 fields.
struct {
    x, y int
    u float32
    _ float32  // padding
    A *[]int
    F func()
}
```

<mark>使用类型声明但没有显式字段名称的字段称为嵌入字段</mark>。嵌入字段必须指定为类型名称 T<sup>1</sup> 或指向非接口类型名称 *T 的指针<sup>2</sup>，并且 T 本身可能不是指针类型。非限定类型名称充当字段名称。

```go
struct {
    T1        // field name is T1
    *T2       // field name is T2
    P.T3      // field name is T3
    *P.T4     // field name is T4
    x, y int  // field names are x and y
}

// 以下声明是非法的，因为字段名称在结构类型中必须是唯一
struct {
    T     // conflicts with embedded field *T and *P.T
    *T    // conflicts with embedded field T and *P.T
    *P.T  // conflicts with embedded field T and *T
}
```

如果 x.f 是表示字段或方法 f 的合法选择器，则结构体 x 中嵌入字段的字段或方法 f 被称为提升。

提升字段的作用与结构的普通字段类似，只是它们不能用作结构的复合文字中的字段名称。

给定一个结构体类型 S 和一个命名类型 T，提升的方法将包含在该结构体的方法集中，如下所示：

- 如果S包含嵌入字段T，则S和*S的方法集都包括接收者T的提升方法。*S的方法集还包括接收者*T的提升方法。

- 如果 S 包含嵌入字段 *T，则 S 和 *S 的方法集都包含接收者为 T 或 *T 的提升方法。

字段声明后可跟一个可选的字面字符串标签，它将成为相应字段声明中所有字段的属性。空标签字符串等同于无标签。标签通过反射接口可见，并参与结构体的类型标识，否则将被忽略。

```go
struct {
    x, y float64 ""  // an empty tag string is like an absent tag
    name string  "any string is permitted as a tag"
    _    [4]byte "ceci n'est pas un champ de structure"
}

// A struct corresponding to a TimeStamp protocol buffer.
// The tag strings define the protocol buffer field numbers;
// they follow the convention outlined by the reflect package.
struct {
    microsec  uint64 `protobuf:"1"`
    serverIP6 uint64 `protobuf:"2"`
}
```

如果包含 T 的类型只是数组或结构类型，则结构类型 T 不得直接或间接包含 T 类型的字段或包含 T 作为组件的类型的字段。

```go
// invalid struct types
type (
    T1 struct{ T1 }            // T1 contains a field of T1
    T2 struct{ f [10]T2 }      // T2 contains T2 as component of an array
    T3 struct{ T4 }            // T3 contains T3 as component of an array in struct T4
    T4 struct{ f [10]T3 }      // T4 contains T4 as component of struct T3 in an array
)

// valid struct types
type (
    T5 struct{ f *T5 }         // T5 contains T5 as component of a pointer
    T6 struct{ f func() T6 }   // T6 contains T6 as component of a function type
    T7 struct{ f [10][]T7 }    // T7 contains T7 as component of a slice in an array
)
```

## Pointer types

指针类型表示指向给定类型变量的所有指针的集合，称为指针的基本类型。未初始化的指针的值为 nil。

```ebnf
PointerType = "*" BaseType .
BaseType    = Type .
```

```go
*Point
*[4]int
```

## Function Types

函数类型表示具有相同参数和结果类型的所有函数的集合。函数类型的未初始化变量的值为 nil。

```ebnf
FunctionType   = "func" Signature .
Signature      = Parameters [ Result ] .
Result         = Parameters | Type .
Parameters     = "(" [ ParameterList [ "," ] ] ")" .
ParameterList  = ParameterDecl { "," ParameterDecl } .
ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
```

在参数或结果列表中，名称 (IdentifierList) 必须全部存在或全部不存在。如果存在，则每个名称代表指定类型的一项（参数或结果），并且签名中的所有非空白名称必须是唯一的。如果不存在，则每种类型代表该类型的一项。参数和结果列表总是带括号的，除非只有一个未命名的结果，它可以写为不带括号的类型。

函数签名中的最终传入参数可能具有以 ... 为前缀的类型。具有此类参数的函数称为可变参数，并且可以使用该参数的零个或多个参数来调用。

```go
func()
func(x int) int
func(a, _ int, z float32) bool
func(a, b int, z float32) (bool)
func(prefix string, values ...int)
func(a, b int, z float64, opt ...interface{}) (success bool)
func(int, int, float64) (float64, *[]int)
func(n int) func(p *T)
```

```go
// demo 1
package main 
import (
    "fmt"
)
type funcTest func(int, int) error

func (f funcTest) exec() error {
    return f(1, 2)
}

func a(b int, c int) error {
    fmt.Println(b, c)
    return errors.New("a")
}

func main() {
    aFunc := funcTest(a) // 显示类型转换 a -> funcTest
    fmt.Println(aFunc.exec())
}
```

```go
// demo 2
package main 
import (
    "fmt"
)

type A struct {
}

func (i *A) A(a, b int) {
    fmt.Println(a, b)
}

type cmdable func(int, int)

func (c cmdable) A1() {
    c(1, 2)
}

func (c cmdable) A2() {
    c(2, 3)
}

type U struct {
    *A
    cmdable
}

type UU struct {
    ia Cmdable
}

func main() {
    as := &A{}
    uu := &UU{ia: &U{A: as, cmdable: as.A}}
    uu.ia.A1()
    uu.ia.A2()
}
```

## Interface types

接口类型定义类型集。接口类型的变量可以存储接口类型集中的任何类型的值。这样的类型被称为实现接口([implement the interface](https://go.dev/ref/spec#Implementing_an_interface))。接口类型的未初始化变量的值为 nil。

```go
InterfaceType  = "interface" "{" { InterfaceElem ";" } "}" .
InterfaceElem  = MethodElem | TypeElem .
MethodElem     = MethodName Signature .
MethodName     = identifier .
TypeElem       = TypeTerm { "|" TypeTerm } .
TypeTerm       = Type | UnderlyingType .
UnderlyingType = "~" Type .
```

接口类型由接口元素列表指定。接口元素既可以是方法，也可以是类型元素，其中类型元素是一个或多个类型项的组合。类型项可以是单个类型，也可以是单个底层类型。

### 基本接口 Basic interfaces

接口最基本的形式指定了一个（可能是空的）方法列表。这种接口定义的类型集是实现所有这些方法的类型集，并且相应的方法集恰好由该接口指定的方法组成。其类型集可以完全由方法列表定义的接口称为基本接口。

```go
// A simple File interface.
interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
}
```

每个显式指定的方法名必须是唯一且不为空的名称。

```go
interface {
    String() string
    String() string  // illegal: String not unique
    _(x int)         // illegal: method must have non-blank name
}
```

不止一种类型可以实现一个接口。例如，如果两个类型 S1 和 S2 具有方法集

```go
func (p T) Read(p []byte) (n int, err error)
func (p T) Write(p []byte) (n int, err error)
func (p T) Close() error
```

(其中 T 代表 S1 或 S2），那么 S1 和 S2 都实现了File接口，而不管 S1 和 S2 可能拥有或共享什么其他方法。

作为接口类型集成员的每个类型都实现该接口。任何给定类型都可以实现多个不同的接口。例如，所有类型都实现空接口，它代表所有（非接口）类型的集合：

```go
interface{}
```

为了方便起见，预先声明的类型 any 是空接口的别名。

同样，请看这个接口规范，它出现在一个类型声明中，定义了一个名为 Locker 的接口：

```go
type Locker interface {
    Lock()
    Unlock()
}
```

如果S1和S2也实现:

```go
func (p T) Lock() { … }
func (p T) Unlock() { … }
```

它们实现了 Locker 接口和 File 接口。

### 嵌入接口 Embedded interfaces

在稍微更一般的形式中，接口 T 可以使用（可能限定的）接口类型名称 E 作为接口元素。这称为在 T 中嵌入接口 E。T 的类型集是 T 显式声明的方法定义的类型集与 T 嵌入接口的类型集的交集。换句话说，T 的类型集是实现 T 的所有显式声明的方法以及 E 的所有方法的所有类型的集合。

```go
type Reader interface {
    Read(p []byte) (n int, err error)
    Close() error
}

type Writer interface {
    Write(p []byte) (n int, err error)
    Close() error
}

// ReadWriter's methods are Read, Write, and Close.
type ReadWriter interface {
    Reader  // includes methods of Reader in ReadWriter's method set
    Writer  // includes methods of Writer in ReadWriter's method set
}
```

嵌入接口时，同名的方法必须具有相同的签名。

```go
type ReadCloser interface {
    Reader   // includes methods of Reader in ReadCloser's method set
    Close()  // illegal: signatures of Reader.Close and Close are different
}
```

### 通用接口 General Interfaces

在最一般的形式中，接口元素也可以是任意的类型项 T，或者是指定底层类型 T 的 ~T 形式的项，或者是项 t1|t2|...|tn 的联合。这些元素与方法规范一起，可以如下精确定义接口的类型集：

- 空接口的类型集合是所有<mark>非接口类型</mark>的集合。

- 非空接口的类型集是<mark>其接口元素的类型集的交集</mark>。

- 方法规范的类型集是所有方法集中包含（该方法规范中定义的方法）的非接口类型的集合。

- 非接口类型项的类型集是<mark>仅由该类型组成的集合</mark>。

- ~T 类型项所表达的类型集是<mark>底层类型为 T 的类型</mark>集合。

- 并集 t<sub>1</sub>|t<sub>2</sub>|…|t<sub>n</sub> 的类型集是所有<mark>类型项t<sub>n</sub>类型集的并集</mark>。

量化“所有非接口类型的集合”不仅指当前程序中声明的所有（非接口）类型，还指所有可能程序中的所有可能类型，因此是无限的。类似地，给定实现特定方法的所有非接口类型的集合，这些类型的方法集的交集将恰好包含该方法，即使手头程序中的所有类型总是将该方法与另一个方法配对。

根据构造，<mark>接口的类型集从不包含接口类型</mark>。

```go
// An interface representing only the type int.
interface {
    int
}

// An interface representing all types with underlying type int.
interface {
    ~int
}

// An interface representing all types with underlying type int that implement the String method.
interface {
    ~int
    String() string
}

// An interface representing an empty type set: there is no type that is both an int and a string.
interface {
    int
    string
}
```

<mark>在 ~T 形式的术语中，T 的底层类型必须是自身，T 不能是接口</mark>。

```go
type MyInt int

interface {
    ~[]byte  // the underlying type of []byte is itself
    ~MyInt   // illegal: the underlying type of MyInt is not MyInt
    ~error   // illegal: error is an interface
}
```

联合元素表示类型集的联合：

```go
// The Float interface represents all floating-point types
// (including any named types whose underlying types are
// either float32 or float64).
type Float interface {
    ~float32 | ~float64
}
```

<mark>T 或 ~T 形式的项中的类型 T 不能是类型参数</mark>([type parameter](https://go.dev/ref/spec#Type_parameter_declarations))，并且所有非接口项的类型集必须是成对不相交的（类型集的成对交集必须为空）。给定一个类型参数 P：

<mark>非基本接口只能用作类型约束，或用作其他接口的约束元素</mark>。<mark>它们不能作为值或变量的类型，也不能作为其他非接口类型的组成部分。</mark>

```go
var x Float                     // illegal: Float is not a basic interface

var x interface{} = Float(nil)  // illegal

type Floatish struct {
    f Float                 // illegal
}
```

接口类型 T 不能直接或间接嵌入属于、包含或嵌入 T 的类型元素。

```go
// illegal: Bad may not embed itself
type Bad interface {
    Bad
}

// illegal: Bad1 may not embed itself using Bad2
type Bad1 interface {
    Bad2
}
type Bad2 interface {
    Bad1
}

// illegal: Bad3 may not embed a union containing Bad3
type Bad3 interface {
    ~int | ~string | Bad3
}

// illegal: Bad4 may not embed an array containing Bad4 as element type
type Bad4 interface {
    [10]Bad4
}
```

### 实现接口 Implementing an interface

类型 T 实现接口 I，如果

- T 不是接口，而是 I 类型集的元素；或者

- T 是一个接口，T 的类型集是 I 类型集的子集。

如果 T 实现了一个接口，则 T 类型的值就实现了该接口。

## 映射类型 Map Types

映射是一种类型（称为元素类型）的无序元素组，由另一种类型（称为键类型）的一组唯一键进行索引。未初始化的映射的值为 nil。

```go
MapType     = "map" "[" KeyType "]" ElementType .
KeyType     = Type .

```

比较运算符([comparison operators](https://go.dev/ref/spec#Comparison_operators)) == 和 != 必须为键类型的操作数完全定义；因此键类型不能是函数、映射或切片。如果键类型是接口类型，则必须为动态键值定义这些比较运算符；失败将导致运行时恐慌。

```go
map[string]int
map[*T]struct{ x, y float64 }
map[string]interface{}
```

映射元素的数量称为映射长度。对于映射 m 来说，可以使用内置函数 len 来发现它的长度，并可能在执行过程中发生变化。在执行过程中，可以使用赋值添加元素，也可以使用索引表达式检索元素；还可以使用删除和清除内置函数删除元素。

使用内置函数 make 创建一个新的空映射值，该函数将映射类型和可选的容量提示作为参数：

```go
make(map[string]int)
make(map[string]int, 100)
```

初始容量不限制其大小：映射会增长以容纳其中存储的项目数量，但 nil 映射除外。 nil 映射相当于空映射，只不过不能添加任何元素。

### 通道类型 Channel types

通道提供了一种并发执行函数([concurrently executing functions](https://go.dev/ref/spec#Go_statements))的机制，通过发送([sending](https://go.dev/ref/spec#Send_statements))和接收([receiving](https://go.dev/ref/spec#Receive_operator))指定元素类型的值来进行通信。未初始化的通道的值为 nil。

```ebnf
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
```

可选的 <- 运算符指定通道方向、发送或接收。如果给定方向，则通道是有方向的，否则是双向的。通过分配或显式转换，通道可以被限制为仅发送或仅接收。

```go
chan T          // can be used to send and receive values of type T
chan<- float64  // can only be used to send float64s
<-chan int      // can only be used to receive ints
```

<- 运算符与可能的最左边的 chan 关联：

```go
chan<- chan int    // same as chan<- (chan int)
chan<- <-chan int  // same as chan<- (<-chan int)
<-chan <-chan int  // same as <-chan (<-chan int)
chan (<-chan int)
```

新的初始化通道值可通过内置函数 make 生成，该函数的参数包括通道类型和可选的容量：

```go
make(chan int, 100)
```

容量（以元素数量表示）设置通道中缓冲区的大小。如果容量为零或不存在，则通道无缓冲，并且仅当发送方和接收方都准备好时通信才会成功。否则，如果缓冲区未满（发送）或不为空（接收），则通道将被缓冲，并且通信会成功而不会阻塞。零通道永远不会准备好进行通信。

可以使用内置函数 close 来关闭通道。接收操作符的多值分配形式报告在通道关闭之前是否发送了接收到的值。

单个通道可用于发送语句、接收操作、调用内置函数 cap 和 len，而无需进一步同步。通道的作用类似于先入先出队列。例如，如果一个程序在通道上发送值，而第二个程序接收这些值，那么接收值的顺序就是发送的顺序。
