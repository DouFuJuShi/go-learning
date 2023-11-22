# 声明和范围 Declarations and scope

声明将非空标识符绑定到常量、类型、类型参数、变量、函数、标签或包。程序中的每个标识符都必须声明。标识符不能在同一块中声明两次，并且标识符不能同时在文件块和包块中声明。

空白标识符可以像声明中的任何其他标识符一样使用，但它不会引入绑定，因此不会被声明。在 package 块中，标识符 init 只能用于 init 函数声明，并且与空白标识符一样，它不会引入新的绑定。

```ebnf
Declaration   = ConstDecl | TypeDecl | VarDecl .
TopLevelDecl  = Declaration | FunctionDecl | MethodDecl .
```

声明的标识符的范围是源文本的范围，其中标识符表示指定的常量、类型、变量、函数、标签或包。

Go 使用块进行词法作用域：

1. 预声明标识符的范围是 Universe 块。

2. 在顶层（任何函数之外）声明的常量、类型、变量或函数（但不包括方法）的标识符的作用域是包块。

3. 导入包的包名范围是包含导入声明的文件的文件块。

4. 表示方法接收者、函数参数或结果变量的标识符的范围是函数体。

5. 表示函数类型参数或由方法接收器声明的标识符的作用域从函数名之后开始，到函数体结束为止。

6. 表示类型的类型参数的标识符的范围从类型名称之后开始，到 TypeSpec 末尾结束。

7. 函数内部声明的常量或变量标识符的范围从 ConstSpec 或 VarSpec（用于短变量声明的 ShortVarDecl）的末尾开始，到最内层包含块的末尾结束。

8. 函数内部声明的类型标识符的范围从 TypeSpec 中的标识符开始，到最内层包含块的末尾结束。

在代码块中声明的标识符可在内部代码块中重新声明。当内部声明的标识符处于作用域中时，它表示内部声明所声明的实体。

包子句不是声明；包名不出现在任何作用域中。其目的是识别属于同一软件包的文件，并指定导入声明的默认软件包名称。

## 标签范围 Label scopes

标签由带标签的语句( [labeled statements](https://go.dev/ref/spec#Labeled_statements))声明，并用在“break”、“continue”和“goto”语句中。定义从未使用过的标签是非法的。与其他标识符相比，标签不是块作用域的，并且不会与非标签的标识符发生冲突。标签的作用域是声明它的函数体，不包括任何嵌套函数的体。

## 空白标识符 Blank identifier

空白标识符由下划线字符 _ 表示。它用作匿名占位符，而不是常规（非空白）标识符，并且在声明、操作数和赋值语句中具有特殊含义。

## 预声明标识符 Predeclared identifiers

以下标识符在 Universe 块([universe block](https://go.dev/ref/spec#Blocks))中隐式声明：

```
Types:
    any bool byte comparable
    complex64 complex128 error float32 float64
    int int8 int16 int32 int64 rune string
    uint uint8 uint16 uint32 uint64 uintptr

Constants:
    true false iota

Zero value:
    nil

Functions:
    append cap clear close complex copy delete imag len
    make max min new panic print println real recover
```

## 导出标识符 Exported identifiers

一个标识符可以导出，以便在另一个包中访问。如果同时具备以下条件，则标识符被导出：

1. 标识符名称的第一个字符是 Unicode 大写字母（Unicode 字符类别 Lu）；和

2. 标识符在包块中声明，或者是字段名称或方法名称。

其他标识符不会导出。

## 标识符的唯一性 Uniqueness of identifiers

在一组标识符中，如果一个标识符与其他标识符都不同，那么这个标识符就被称为唯一标识符。如果两个标识符的拼写不同，或者出现在不同的软件包中且没有导出，那么它们就是不同的。否则，它们就是相同的。

## 常量定义 Constant declarations

常量声明将标识符列表（常量名称）绑定到常量表达式列表的值。标识符的数量必须等于表达式的数量，并且左边的第n个标识符绑定到右边第n个表达式的值。

```ebnf
ConstDecl      = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
ConstSpec      = IdentifierList [ [ Type ] "=" ExpressionList ] .

IdentifierList = identifier { "," identifier } .
ExpressionList = Expression { "," Expression } .
```

如果存在类型，则所有常量均采用指定的类型，并且表达式必须可分配给该类型，该类型不能是类型参数。如果省略类型，则常量采用相应表达式的各个类型。如果表达式值是无类型常量，则声明的常量保持无类型并且常量标识符表示常量值。例如，如果表达式是浮点文字，则常量标识符表示浮点常量，即使文字的小数部分为零。

```go
const Pi float64 = 3.14159265358979323846
const zero = 0.0         // untyped floating-point constant
const (
    size int64 = 1024
    eof        = -1  // untyped integer constant
)
const a, b, c = 3, 4, "foo"  // a = 3, b = 4, c = "foo", untyped integer and string constants
const u, v float32 = 0, 3    // u = 0.0, v = 3.0
```

在带括号的 const 声明列表中，<mark>除第一个 ConstSpec 之外的任何表达式列表都可以省略</mark>。<mark>这样的空列表相当于前面第一个非空表达式列表及其类型（如果有）的文本替换</mark>。因此，省略表达式列表相当于重复前面的列表。标识符的数量必须等于前面列表中的表达式的数量。与 iota 常量生成器一起，该机制允许顺序值的轻量级声明：

```go
const (
    Sunday = iota
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Partyday
    numberOfDays  // this constant is not exported
)
```

## Iota

在常量声明中，预声明的标识符 iota 表示连续的无类型整型常量。它的值是该常量声明中相应 ConstSpec 的索引，从零开始。它可用于构造一组相关常量：

```go
const (
    c0 = iota  // c0 == 0
    c1 = iota  // c1 == 1
    c2 = iota  // c2 == 2
)

const (
    a = 1 << iota  // a == 1  (iota == 0)
    b = 1 << iota  // b == 2  (iota == 1)
    c = 3          // c == 3  (iota == 2, unused)
    d = 1 << iota  // d == 8  (iota == 3)
)

const (
    u         = iota * 42  // u == 0     (untyped integer constant)
    v float64 = iota * 42  // v == 42.0  (float64 constant)
    w         = iota * 42  // w == 84    (untyped integer constant)
)

const x = iota  // x == 0
const y = iota  // y == 0
```

根据定义，在同一个 ConstSpec 中多次使用 iota 都具有相同的值：

```go
const (
    bit0, mask0 = 1 << iota, 1<<iota - 1  // bit0 == 1, mask0 == 0  (iota == 0)
    bit1, mask1                           // bit1 == 2, mask1 == 1  (iota == 1)
    _, _                                  //                        (iota == 2, unused)
    bit3, mask3                           // bit3 == 8, mask3 == 7  (iota == 3)
)
```

最后一个示例利用了最后一个非空表达式列表的隐式重复。

## 类型声明 Type declarations

类型声明将标识符（即类型名）与类型绑定。类型声明有两种形式：别名声明和类型定义。

```ebnf
TypeDecl = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
TypeSpec = AliasDecl | TypeDef .
```

### 别名声明 Alias declarations

别名声明将标识符绑定到给定类型。

```ebnf
AliasDecl = identifier "=" Type .
```

在标识符的范围内，它充当类型的别名。

```go
type (
    nodeList = []*Node  // nodeList and []*Node are identical types
    Polar    = polar    // Polar and polar denote identical types
)
```

### 类型定义 Type definitions

类型定义创建了一个新的、与众不同的类型，新类型的底层类型和操作与声明中给定(引用、使用)的类型相同，并绑定了一个标识符（类型名）。

```go
TypeDef = identifier [ TypeParameters ] Type .
```

新类型称为定义类型。它不同于任何其他类型，包括创建它的类型。

```go
type (
    Point struct{ x, y float64 }  // Point and struct{ x, y float64 } are different types
    polar Point                   // polar and Point denote different types
)

type TreeNode struct {
    left, right *TreeNode
    value any
}

type Block interface {
    BlockSize() int
    Encrypt(src, dst []byte)
    Decrypt(src, dst []byte)
}
```

定义的类型可能有与之相关的方法。它不会继承与给定类型绑定的任何方法，但接口类型或复合类型元素的方法集保持不变：

```go
// A Mutex is a data type with two methods, Lock and Unlock.
type Mutex struct         { /* Mutex fields */ }
func (m *Mutex) Lock()    { /* Lock implementation */ }
func (m *Mutex) Unlock()  { /* Unlock implementation */ }

// NewMutex has the same composition as Mutex but its method set is empty.
type NewMutex Mutex

// The method set of PtrMutex's underlying type *Mutex remains unchanged,
// but the method set of PtrMutex is empty.
type PtrMutex *Mutex

// The method set of *PrintableMutex contains the methods
// Lock and Unlock bound to its embedded field Mutex.
type PrintableMutex struct {
    Mutex
}

// MyBlock is an interface type that has the same method set as Block.
type MyBlock Block
```

类型定义可用于定义不同的布尔、数字或字符串类型，并将方法与之关联：

```go
type TimeZone int

const (
    EST TimeZone = -(5 + iota)
    CST
    MST
    PST
)

func (tz TimeZone) String() string {
    return fmt.Sprintf("GMT%+dh", tz)
}
```

如果类型定义指定类型参数，则类型名称表示泛型类型。泛型类型在使用时必须实例化。

```go
type List[T any] struct {
    next  *List[T]
    value T
}
```

在类型定义中，给定类型不能是类型参数。

```go
type T[P any] P    // illegal: P is a type parameter

func f[T any]() {
    type L T   // illegal: T is a type parameter declared by the enclosing function
}
```

泛型类型也可能具有与其关联的方法。在这种情况下，方法接收者必须声明与泛型类型定义中存在的相同数量的类型参数。

```go
// The method Len returns the number of elements in the linked list l.
func (l *List[T]) Len() int  { … }
```

## 类型形参声明 Type parameter declarations

类型参数列表声明泛型函数或类型声明的类型参数。类型参数列表看起来像普通的函数参数列表([function parameter list](https://go.dev/ref/spec#Function_types))，只是类型参数名称必须全部存在并且列表用方括号而不是圆括号括起来。

```ebnf
TypeParameters  = "[" TypeParamList [ "," ] "]" .
TypeParamList   = TypeParamDecl { "," TypeParamDecl } .
TypeParamDecl   = IdentifierList TypeConstraint .
```

列表中的所有非空白名称必须是唯一的。每个名称都声明一个类型参数，它是一个新的、不同的命名类型，充当声明中（迄今为止）未知类型的占位符。在泛型函数或类型实例化时，类型形参将替换为类型实参。

```go
[P any]
[S interface{ ~[]byte|string }]
[S ~[]E, E any]
[P Constraint[int]]
[_ any]
```

正如每个普通函数参数都有一个参数类型一样，每个类型参数都有一个相应的（元）类型，称为类型约束([*type constraint*](https://go.dev/ref/spec#Type_constraints))。

当一个泛型的类型参数列表声明了带有约束 C 的单个类型参数 P，而文本 P C 构成了一个有效的表达式时，就会出现解析歧义：

```go
type T[P interface{*C}] …
type T[P *C,] …
```

类型参数也可以由与泛型类型相关联的方法声明的接收者规范来声明。

在泛型类型 T 的类型参数列表中，类型约束不能（直接或通过另一个泛型类型的类型参数列表间接）引用 T。

```go
type T1[P T1[P]] …                    // illegal: T1 refers to itself
type T2[P interface{ T2[int] }] …     // illegal: T2 refers to itself
type T3[P interface{ m(T3[int])}] …   // illegal: T3 refers to itself
type T4[P T5[P]] …                    // illegal: T4 refers to T5 and
type T5[P T4[P]] …                    //          T5 refers to T4

type T6[P int] struct{ f *T6[P] }     // ok: reference to T6 is not in type parameter list
```

### 类型约束 Type constraints

类型约束是一个接口，它定义了相应类型参数的允许类型参数集，并控制该类型参数值所支持的操作。

```ebnf
TypeConstraint = TypeElem .
```

如果约束是 interface{E} 形式的接口文字，其中 E 是嵌入的类型元素（不是方法），那么在类型参数列表中，为了方便起见，可以省略外层 interface{ ... }：

```go
[T []P]                      // = [T interface{[]P}]
[T ~int]                     // = [T interface{~int}]
[T int|string]               // = [T interface{int|string}]
type Constraint ~int         // illegal: ~int is not in a type parameter list
```

预先声明的接口类型可比较表示严格可比较的所有非接口类型的集合。

尽管不是类型参数的接口是可比较[comparable](https://go.dev/ref/spec#Comparison_operators)的，但它们并不严格可比较，因此它们不实现可比较。然而，它们满足 [satisfy](https://go.dev/ref/spec#Satisfying_a_type_constraint)可比性。

```go
int                          // implements comparable (int is strictly comparable)
[]byte                       // does not implement comparable (slices cannot be compared)
interface{}                  // does not implement comparable (see above)
interface{ ~int | ~string }  // type parameter only: implements comparable (int, string types are strictly comparable)
interface{ comparable }      // type parameter only: implements comparable (comparable implements itself)
interface{ ~int | ~[]byte }  // type parameter only: does not implement comparable (slices are not comparable)
interface{ ~struct{ any } }  // type parameter only: does not implement comparable (field any is not strictly comparable)
```

可比较的接口和（直接或间接）嵌入可比较的接口<mark>只能用作类型约束</mark>。它们不能是值或变量的类型，也不能是其他非接口类型的组件。

### 满足类型约束 Satisfying a type constraint

如果类型 T 是 类型约束接口 C 所定义类型集的元素，则类型参数 T 满足类型约束 C；即，如果 T 实现 C。作为例外，严格可比较类型约束也可以由可比较（不一定严格可比较）类型参数来满足。更确切地说：

如果符合以下条件，则类型 T 满足约束条件 C

- T实现了C；或者

- C 可以写成 interface{ comparable; E } 的形式，其中 E 是一个基本接口，T 是可比较的并实现了 E。

```go
type argument      type constraint                // constraint satisfaction

int                interface{ ~int }              // satisfied: int implements interface{ ~int }
string             comparable                     // satisfied: string implements comparable (string is strictly comparable)
[]byte             comparable                     // not satisfied: slices are not comparable
any                interface{ comparable; int }   // not satisfied: any does not implement interface{ int }
any                comparable                     // satisfied: any is comparable and implements the basic interface any
struct{f any}      comparable                     // satisfied: struct{f any} is comparable and implements the basic interface any
any                interface{ comparable; m() }   // not satisfied: any does not implement the basic interface interface{ m() }
interface{ m() }   interface{ comparable; m() }   // satisfied: interface{ m() } is comparable and implements the basic interface interface{ m() }
```

由于约束满足规则中的例外，比较类型参数类型的操作数可能会在运行时出现恐慌（即使可比较的类型参数始终是严格可比较的）。

## 变量声明 Variable declarations

变量声明创建一个或多个变量，将相应的标识符绑定到它们，并为每个变量赋予一个类型和一个初始值。

```ebnf
VarDecl     = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
VarSpec     = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
```

```go
var i int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
    i       int
    u, v, s = 2.0, 3.0, "bar"
)
var re, im = complexSqrt(-1)
```

如果给出了表达式列表，则变量将使用遵循赋值语句规则的表达式进行初始化。否则，每个变量都将初始化为零值。

如果存在类型，则每个变量都会被赋予该类型。否则，每个变量都会在赋值中被赋予相应初始化值的类型。如果该值是无类型常量，则首先将其隐式转换为其默认类型；如果它是无类型布尔值，则首先将其隐式转换为 bool 类型。预声明值 nil 不能用于初始化没有显式类型的变量。

```go
var d = math.Sin(0.5)  // d is float64
var i = 42             // i is int
var t, ok = x.(T)      // t is T, ok is bool
var n = nil            // illegal
```

实现限制：<mark>如果从未使用过变量，则编译器可能会认为在函数体内声明该变量是非法的。</mark>

## 简短变量声明 Short variable declarations

简短的变量声明使用以下语法：

```go
ShortVarDecl = IdentifierList ":=" ExpressionList .
```

它是带有初始化表达式但没有类型的常规变量声明的简写：

```ebnf
"var" IdentifierList "=" ExpressionList .
```

```go
i, j := 0, 10
f := func() int { return 7 }
ch := make(chan int)
r, w, _ := os.Pipe()  // os.Pipe() returns a connected pair of Files and an error, if any
_, y, _ := coord(p)   // coord() returns three values; only interested in y coordinate
```

与常规变量声明不同，短变量声明可以重新声明变量，前提是这些变量最初在同一块（或参数列表，如果该块是函数体）中以相同类型声明，并且至少有一个非空变量是新的。因此，重新声明只能出现在多变量短声明中。重新声明不会引入新变量；它只是为原始值分配一个新值。 := 左侧的非空白变量名称必须是唯一的。

```go
field1, offset := nextField(str, 0)
field2, offset := nextField(str, offset)  // redeclares offset
x, y, x := 1, 2, 3                        // illegal: x repeated on left side of :=
```

短变量声明只能出现在函数内部。在某些上下文中，例如“if”、“for”或“switch”语句的初始值设定项，它们可用于声明局部临时变量。

## 函数声明 Function declarations

函数声明将一个标识符（即函数名）绑定到一个函数上。

```ebnf
FunctionDecl = "func" FunctionName [ TypeParameters ] Signature [ FunctionBody ] .
FunctionName = identifier .
FunctionBody = Block .
```

如果函数的签名声明了结果参数，则函数体的语句列表必须以终止语句结束。

```go
func IndexRune(s string, r rune) int {
    for i, c := range s {
        if c == r {
            return i
        }
    }
    // invalid: missing return statement
}
```

如果函数声明中指定了类型参数，则函数名称表示泛函。在调用泛型函数或将其用作值之前，必须先将其实例化。

```go
func min[T ~int|~float64](x, y T) T {
    if x < y {
        return x
    }
    return y
}
```

没有类型参数的函数声明可以省略主体。这种声明为在 Go 之外实现的函数（如汇编例程）提供了签名。

```go
func flushICache(begin, end uintptr)  // implemented externally
```

## 方法声明 Method declarations

方法是带有接收器的函数。方法声明将一个标识符（即方法名）与方法绑定，并将方法与接收器的基本类型关联起来。

```ebnf
MethodDecl = "func" Receiver MethodName Signature [ FunctionBody ] .
Receiver   = Parameters .
```

接收者是通过方法名称前面的额外参数部分指定的。该参数部分必须声明一个非可变参数，即接收者。它的类型必须是定义的类型 T 或指向定义的类型 T 的指针，后面可能跟着方括号括起来的类型参数名称 [P1, P2, …] 列表。 T称为接收者基本类型。接收者基类型不能是指针或接口类型，并且必须在与方法相同的包中定义。据说该方法绑定到其接收者基类型，并且方法名称仅在类型 T 或 *T 的选择器中可见。
非空接收者标识符在方法签名中必须是唯一的。如果方法体内未引用接收者的值，则声明中可以省略其标识符。这同样适用于函数和方法的参数。
对于基类型，绑定到它的方法的非空名称必须是唯一的。如果基类型是结构类型，则非空方法名称和字段名称必须不同。
给定定义类型 Point 声明

```go
func (p *Point) Length() float64 {
    return math.Sqrt(p.x * p.x + p.y * p.y)
}

func (p *Point) Scale(factor float64) {
    p.x *= factor
    p.y *= factor
}
```

将接收器类型为 *Point 的 Length 和 Scale 方法绑定到基本类型 Point。

如果接收者基类型是泛型类型，则接收者规范必须为要使用的方法声明相应的类型参数。这使得接收者类型参数可供该方法使用。从语法上讲，此类型参数声明看起来像接收者基类型的实例化：类型参数必须是表示所声明的类型参数的标识符，一个对应于接收者基类型的每个类型参数。类型参数名称不需要与接收者基本类型定义中相应的参数名称相匹配，并且所有非空白参数名称在接收者参数部分和方法签名中必须是唯一的。接收者类型参数约束由接收者基本类型定义隐含：相应的类型参数具有相应的约束。

```go
type Pair[A, B any] struct {
    a A
    b B
}

func (p Pair[A, B]) Swap() Pair[B, A]  { … }  // receiver declares A, B
func (p Pair[First, _]) First() First  { … }  // receiver declares First, corresponds to A in Pair
```
