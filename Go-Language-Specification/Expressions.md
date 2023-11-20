# 表达式 Expressions

表达式通过将运算符和函数应用于操作数来指定值的计算。

## 操作数 Operands

操作数表示表达式中的基本值。操作数可以是文字、表示常量、变量或函数的（可能是限定的）非空白标识符，或者是带括号的表达式。

```go
Operand     = Literal | OperandName [ TypeArgs ] | "(" Expression ")" .
Literal     = BasicLit | CompositeLit | FunctionLit .
BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
OperandName = identifier | QualifiedIdent .
```

表示泛型函数的操作数名称后面可以跟着类型实参列表；结果操作数是一个实例化函数。

空白标识符只能作为操作数出现在赋值语句的左侧。

实现限制：如果操作数的类型是具有空类型集的类型参数，则编译器不需要报告错误。具有此类类型参数的函数无法实例化；任何尝试都会导致实例化站点出现错误。

## 限定标志符 Qualified identifiers

限定标识符是用包名称前缀限定的标识符。包名称和标识符都不能为空。

```ebnf
QualifiedIdent = PackageName "." identifier .
```

限定标识符访问必须导入的不同包中的标识符。标识符必须在该包的包块中导出和声明。

```go
math.Sin // denotes the Sin function in package math
```

## 复合字面量 Composite literals

复合文字在每次求值时都会构造新的复合值。它们由文字类型和后跟大括号元素列表组成。每个元素前面可以有一个相应的键（可选）。

```ebnf
CompositeLit  = LiteralType LiteralValue .
LiteralType   = StructType | ArrayType | "[" "..." "]" ElementType |
                SliceType | MapType | TypeName [ TypeArgs ] .
LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
ElementList   = KeyedElement { "," KeyedElement } .
KeyedElement  = [ Key ":" ] Element .
Key           = FieldName | Expression | LiteralValue .
FieldName     = identifier .
Element       = Expression | LiteralValue .
```

LiteralType 的核心类型 T 必须是<mark>结构体、数组、切片或映射类型</mark>（语法强制执行此约束，除非该类型作为 TypeName 给出）。元素和键的类型必须可分配给类型 T 的相应字段、元素和键类型；没有额外的转换。该键被解释为结构体文字的字段名称、数组和切片文字的索引以及映射文字的键。对于映射文字，所有元素都必须有一个键。指定具有相同字段名称或常量键值的多个元素是错误的。对于非常量映射键，请参阅有关评估顺序的部分。

对于结构体文字，适用以下规则：

- 键必须是结构类型中声明的字段名称。

- 不包含任何键的元素列表必须按照声明字段的顺序列出每个结构体字段的元素。

- 如果任何元素都有一个键，那么每个元素都必须有一个键。

- 包含键的元素列表不需要每个结构体字段都有一个元素。省略的字段获得该字段的零值。

- 文字可以省略元素列表；这样的文字的计算结果是其类型的零值。

- 为属于不同包的结构的非导出字段指定元素是错误的。

鉴于声明

```go
type Point3D struct { x, y, z float64 }
type Line struct { p, q Point3D }
```

```go
// 可以这样写
origin := Point3D{}                            // zero value for Point3D
line := Line{origin, Point3D{y: -4, z: 12.3}}  // zero value for line.q.x
```

对于数组和切片文字，适用以下规则：

- 每个元素都有一个关联的整数索引，标记其在数组中的位置。

- 带有键的元素使用键作为其索引。键必须是由 int 类型值表示的非负常量；如果它是类型化的，那么它必须是整数类型。

- 没有键的元素使用前一个元素的索引加一。如果第一个元素没有键，则其索引为零。

获取复合文字的地址会生成一个指向用该文字的值初始化的唯一变量的指针。

```go
var pointer *Point3D = &Point3D{y: 1000}
```

请注意，切片或映射类型的零值与同一类型的已初始化但为空的值不同。因此，获取空切片或映射复合文字的地址并不具有与使用 new 分配新切片或映射值相同的效果。

```go
p1 := &[]int{}    // p1 points to an initialized, empty slice with value []int{} and length 0
p2 := new([]int)  // p2 points to an uninitialized slice with value nil and length 0
```

数组文字的长度是文字类型中指定的长度。如果文字中提供的元素少于长度，则缺失的元素将设置为数组元素类型的零值。为元素提供索引值超出数组索引范围的元素是错误的。符号 ... 指定数组长度等于最大元素索引加一。

```go
buffer := [10]string{}             // len(buffer) == 10
intSet := [6]int{1, 2, 3, 5}       // len(intSet) == 6
days := [...]string{"Sat", "Sun"}  // len(days) == 2
```

切片文字描述了整个底层数组文字。因此，切片文字的长度和容量是最大元素索引加一。切片文字具有以下形式

```go
[]T{x1, x2, … xn}
```

是应用于数组的切片操作的简写：

```go
tmp := [n]T{x1, x2, … xn}
tmp[0 : n]
```

在数组、切片或映射类型 T 的复合文字中，如果本身是复​​合文字的元素或映射键与 T 的元素或键类型相同，则可以省略相应的文字类型。类似地，作为地址的元素或键当元素或键类型为 *T 时，复合文字的 可以省略 &T。

```go
[...]Point{{1.5, -3.5}, {0, 0}}     // same as [...]Point{Point{1.5, -3.5}, Point{0, 0}}
[][]int{{1, 2, 3}, {4, 5}}          // same as [][]int{[]int{1, 2, 3}, []int{4, 5}}
[][]Point{{{0, 1}, {1, 2}}}         // same as [][]Point{[]Point{Point{0, 1}, Point{1, 2}}}
map[string]Point{"orig": {0, 0}}    // same as map[string]Point{"orig": Point{0, 0}}
map[Point]string{{0, 0}: "orig"}    // same as map[Point]string{Point{0, 0}: "orig"}

type PPoint *Point
[2]*Point{{1.5, -3.5}, {}}          // same as [2]*Point{&Point{1.5, -3.5}, &Point{}}
[2]PPoint{{1.5, -3.5}, {}}          // same as [2]PPoint{PPoint(&Point{1.5, -3.5}), PPoint(&Point{})}
```

当使用 LiteralType 的 TypeName 形式的复合文字作为“if”、“for”或“switch”语句块的关键字和左大括号之间的操作数出现时，会出现解析歧义，并且复合文字为不包含在圆括号、方括号或花括号中。在这种罕见的情况下，文字的左大括号被错误地解析为引入语句块的大括号。为了解决歧义，复合文字必须出现在括号内。

```go
if x == (T{a,b,c}[i]) { … }
if (x == T{a,b,c}[i]) { … }
```

有效数组、切片和映射文字的示例：

```go
// list of prime numbers
primes := []int{2, 3, 5, 7, 9, 2147483647}

// vowels[ch] is true if ch is a vowel
vowels := [128]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true}

// the array [10]float32{-1, 0, 0, 0, -0.1, -0.1, 0, 0, 0, -1}
filter := [10]float32{-1, 4: -0.1, -0.1, 9: -1}

// frequencies in Hz for equal-tempered scale (A4 = 440Hz)
noteFrequency := map[string]float32{
    "C0": 16.35, "D0": 18.35, "E0": 20.60, "F0": 21.83,
    "G0": 24.50, "A0": 27.50, "B0": 30.87,
}
```

## 函数字面量 Function literals

函数字面量表示匿名函数。函数文字不能声明类型参数。

```ebnf
FunctionLit = "func" Signature FunctionBody .
```

```go
func(a, b int, z float64) bool { return a*b < int(z) }
```

函数文字可以分配给变量或直接调用。

```go
f := func(x, y int) int { return x + y }
func(ch chan int) { ch <- ACK }(replyChan)
```

函数字面量是闭包：它们可以引用周围函数中定义的变量。然后，这些变量在周围的函数和函数文字之间共享，并且只要可访问，它们就会一直存在。

## 主表达式 Primary expressions

主表达式是一元表达式和二元表达式的操作数。

```ebnf
PrimaryExpr =
    Operand |
    Conversion |
    MethodExpr |
    PrimaryExpr Selector |
    PrimaryExpr Index |
    PrimaryExpr Slice |
    PrimaryExpr TypeAssertion |
    PrimaryExpr Arguments .

Selector       = "." identifier .
Index          = "[" Expression [ "," ] "]" .
Slice          = "[" [ Expression ] ":" [ Expression ] "]" |
                 "[" [ Expression ] ":" Expression ":" Expression "]" .
TypeAssertion  = "." "(" Type ")" .
Arguments      = "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")" .
```

```go
x
2
(s + ".txt")
f(3.1415, true)
Point{1, 2}
m["foo"]
s[i : j + 1]
obj.color
f.p[i].x()
```

## 选择器 Selectors

对于不是包名称的主表达式 x，选择器表达式

```ebnf
x.f
```

表示值 x（或有时 *x；见下文）的字段或方法 f。标识符f称为（字段或方法）选择器；它不能是空白标识符。选择器表达式的类型是 f 的类型。如果 x 是包名称，请参阅有关限定标识符的部分。

选择器 f 可以表示类型 T 的字段或方法 f，也可以指 T 的嵌套嵌入字段的字段或方法 f。遍历到达 f 的嵌入字段的数量称为它在 T 中的深度。 T 中声明的字段或方法 f 的深度为零。 T 中嵌入字段 A 中声明的字段或方法 f 的深度是 A 中 f 的深度加一。

以下规则适用于选择器：

1. 对于 T 或 *T 类型的值 x（其中 T 不是指针或接口类型），x.f 表示 T 中最浅深度处存在 f 的字段或方法。如果不存在恰好一个深度最浅的 f，则选择器表达式是非法的。

2. 对于类型 I 的值 x（其中 I 是接口类型），x.f 表示 x 的动态值的名称为 f 的实际方法。如果 I 的方法集中不存在名称为 f 的方法，则选择器表达式非法。

3. 作为例外，如果 x 的类型是定义的指针类型并且 (*x).f 是表示字段（但不是方法）的有效选择器表达式，则 x.f 是 (*x).f 的简写。

4. 在所有其他情况下，x.f 都是非法的。

5. 如果 x 是指针类型且值为 nil 并且 x.f 表示结构体字段，则对 x.f 进行赋值或求值会导致运行时恐慌。

6. 如果 x 是接口类型且值为 nil，则调用或计算方法 x.f 会导致运行时恐慌。

例如，给定声明：

```go
type T0 struct {
    x int
}

func (*T0) M0()

type T1 struct {
    y int
}

func (T1) M1()

type T2 struct {
    z int
    T1
    *T0
}

func (*T2) M2()

type Q *T2

var t T2     // with t.T0 != nil
var p *T2    // with p != nil and (*p).T0 != nil
var q Q = p
```

有人可能会写：

```go
t.z          // t.z
t.y          // t.T1.y
t.x          // (*t.T0).x

p.z          // (*p).z
p.y          // (*p).T1.y
p.x          // (*(*p).T0).x

q.x          // (*(*q).T0).x        (*q).x is a valid field selector

p.M0()       // ((*p).T0).M0()      M0 expects *T0 receiver
p.M1()       // ((*p).T1).M1()      M1 expects T1 receiver
p.M2()       // p.M2()              M2 expects *T2 receiver
t.M2()       // (&t).M2()           M2 expects *T2 receiver, see section on Calls
```

但以下内容无效：

```go
q.M0()       // (*q).M0 is valid but not a field selector
```

## 方法表达式 Method expressions

如果 M 位于类型 T 的方法集中，则 T.M 是一个可作为常规函数调用的函数，其参数与 M 相同，并以作为方法接收者的附加参数为前缀。

```ebnf
MethodExpr    = ReceiverType "." MethodName .
ReceiverType  = Type .
```

考虑一个具有两个方法的结构体类型 T：Mv（其接收者的类型为 T）和 Mp（其接收者的类型为 *T）。

```go
type T struct {
    a int
}
func (tv  T) Mv(a int) int         { return 0 }  // value receiver
func (tp *T) Mp(f float32) float32 { return 1 }  // pointer receiver

var t T
```

表达方式

```go
T.Mv
```

产生一个与 Mv 等效的函数，但具有显式接收器作为其第一个参数；它有签名

```go
func(tv T, a int) int
```

该函数可以通过显式接收器正常调用，因此这五个调用是等效的：

```go
t.Mv(7)
T.Mv(t, 7)
(T).Mv(t, 7)
f1 := T.Mv; f1(t, 7)
f2 := (T).Mv; f2(t, 7)
```

类似地，表达式

```go
(*T).Mp
```

产生一个代表带有签名的 Mp 的函数值

```go
func(tp *T, f float32) float32
```

对于具有值接收器的方法，可以派生出具有显式指针接收器的函数，因此

```go
(*T).Mv
```

产生一个代表带有签名的 Mv 的函数值

```go
func(tv *T, a int) int
```

这样的函数间接通过接收者创建一个值作为接收者传递给底层方法；该方法不会覆盖在函数调用中传递的地址的值。
最后一种情况，即指针接收器方法的值接收器函数是非法的，因为指针接收器方法不在值类型的方法集中。
从方法派生的函数值通过函数调用语法进行调用；接收者作为调用的第一个参数提供。也就是说，给定 f := T.Mv，f 被调用为 f(t, 7) 而不是 t.f(7)。要构造绑定接收者的函数，请使用函数文字或方法值。
从接口类型的方法派生函数值是合法的。生成的函数采用该接口类型的显式接收器。

## 方法值 Method values

## 索引表达式 Index expressions

## 切片表达式 Slice expressions
