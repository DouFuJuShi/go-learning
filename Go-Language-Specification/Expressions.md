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

如果表达式 x 具有静态类型 T 并且 M 在类型 T 的方法集中，则 x.M 称为方法值。方法值 x.M 是一个函数值，可以使用与 x.M 的方法调用相同的参数进行调用。在计算方法值期间计算并保存表达式 x；然后，保存的副本将用作任何调用中的接收者，该调用可能稍后执行。

```go
type S struct { *T }
type T int
func (t T) M() { print(t) }

t := new(T)
s := S{T: t}
f := t.M                    // receiver *t is evaluated and stored in f
g := s.M                    // receiver *(s.T) is evaluated and stored in g
*t = 42                     // does not affect stored receivers in f and g
```

类型T可以是接口类型或非接口类型。

正如上面对方法表达式的讨论一样，考虑一个具有两个方法的结构体类型 T：Mv（其接收者的类型为 T）和 Mp（其接收者的类型为 *T）。

```go
type T struct {
    a int
}
func (tv  T) Mv(a int) int         { return 0 }  // value receiver
func (tp *T) Mp(f float32) float32 { return 1 }  // pointer receiver

var t T
var pt *T
func makeT() T
```

  表达式：

```go
t.Mv
```

产生类型的函数值

```go
func(int) int
```

这两个调用是等效的：

```go
t.Mv(7)
f := t.Mv; f(7)
```

类似地，表达式

```go
pt.Mp
```

产生类型的函数值

```go
func(float32) float32
```

与选择器一样，使用指针对具有值接收器的非接口方法的引用将自动取消对该指针的引用：pt.Mv 相当于 (*pt).Mv。

与方法调用一样，使用可寻址值对具有指针接收器的非接口方法的引用将自动采用该值的地址：t.Mp 相当于 (&t).Mp。

```go
f := t.Mv; f(7)   // like t.Mv(7)
f := pt.Mp; f(7)  // like pt.Mp(7)
f := pt.Mv; f(7)  // like (*pt).Mv(7)
f := t.Mp; f(7)   // like (&t).Mp(7)
f := makeT().Mp   // invalid: result of makeT() is not addressable
```

尽管上面的示例使用非接口类型，但从接口类型的值创建方法值也是合法的。

```go
var i interface { M(int) } = myVal
f := i.M; f(7)  // like i.M(7)
```

## 索引表达式 Index expressions

索引表达式主要表达形式：

```go
a[x]
```

表示数组的元素、指向由 x 索引的数组、切片、字符串或映射的指针。值 x 分别称为索引或映射键。以下规则适用：

如果 a 既不是映射也不是类型参数：

- 索引 x 必须是无类型常量或其核心类型必须是整数

- 常量索引必须是非负的并且可以用 int 类型的值表示

- 无类型常量索引的指定类型为 int

- 如果 0 <= x < len(a)，则索引 x 在范围内，否则超出范围

对于数组类型 A 的 a：

- 常量索引必须在范围内

- 如果x 在运行时超出范围，发生运行时恐慌

- a[x] 是索引 x 处的数组元素，a[x] 的类型是 A 的元素类型

对于指向数组类型的指针：

- a[x] 是 (*a)[x] 的简写 

对于切片类型 S 的：

- 如果 x 在运行时超出范围，则会发生运行时恐慌

- a[x] 是索引 x 处的切片元素，a[x] 的类型是 S 的元素类型

- 

对于字符串类型：

- 如果字符串 a 也是常量，则常量索引必须在范围内

- 如果 x 在运行时超出范围，则会发生运行时恐慌

- a[x] 是索引 x 处的非常量字节值，a[x] 的类型是 byte

- a[x] 不能进行赋值

对于地图类型为 M 的：

- x 的类型必须可分配给 M 的键类型

- 如果映射包含键为 x 的条目，则 a[x] 是键为 x 的映射元素，a[x] 的类型是 M 的元素类型

- 如果映射为 nil 或不包含这样的条目，则 a[x] 是 M 元素类型的零值

对于类型参数类型 P：

- 索引表达式 a[x] 必须对 P 类型集中所有类型的值都有效。

- P的类型集中所有类型的元素类型必须相同。在这种情况下，字符串类型的元素类型是字节。

- 如果P的类型集中存在映射类型，则该类型集中的所有类型都必须是映射类型，并且各自的键类型必须全部相同。

- a[x] 是索引 x 处的数组、切片或字符串元素，或者是带有 P 实例化类型参数的键 x 的映射元素，并且 a[x] 的类型是（相同）的类型元素类型。

- 如果 P 的类型集包括字符串类型，则不能赋值给 a[x]。

否则 a[x] 是非法的。

map[K]V 类型的映射 a 上的索引表达式，用于赋值语句或特殊形式的初始化

```go
v, ok = a[x]
v, ok := a[x]
var v, ok = a[x]
```

产生一个额外的无类型布尔值。如果键 x 存在于映射中，则 ok 的值为 true，否则为 false。

为 nil 映射的元素赋值会导致运行时恐慌。

## 切片表达式 Slice expressions

切片表达式从字符串、数组、数组指针或切片构造子字符串或切片。有两种变体：指定下限和上限的简单形式，以及还指定容量上限的完整形式。

### 简单切片表达式 Simple slice expressions

主要表达方式

```go
a[low : high]
```

构造一个子字符串或切片。 a 的核心类型必须是字符串、数组、数组指针、切片或字节串。索引 low 和 high 选择操作数 a 的哪些元素出现在结果中。结果的索引从 0 开始，长度等于高 - 低。切片数组 a 后

```go
a := [5]int{1, 2, 3, 4, 5}
s := a[1:4]
```

切片 s 的类型为 []int，长度为 3，容量为 4，元素数

```go
s[0] == 2
s[1] == 3
s[2] == 4
```

为了方便起见，可以省略任何索引。缺失的低指数默认为零；缺失的高索引默认为切片操作数的长度：

```go
a[2:]  // same as a[2 : len(a)]
a[:3]  // same as a[0 : 3]
a[:]   // same as a[0 : len(a)]
```

如果 a 是指向数组的指针，则 a[low : high] 是 (*a)[low : high] 的简写。

对于数组或字符串，如果 0 <= low <= high <= len(a)，则索引在范围内，否则索引超出范围。对于切片，索引上限是切片容量 cap(a) 而不是长度。常量索引必须是非负的并且可以用 int 类型的值表示；对于数组或常量字符串，常量索引也必须在范围内。如果两个指数都是常数，则它们必须满足 low <= high。如果索引在运行时超出范围，则会发生运行时恐慌。

除无类型字符串外，如果切片操作数是字符串或切片，则切片操作的结果是与操作数类型相同的非常量值。对于无类型字符串操作数，结果是字符串类型的非常量值。如果切片操作数是数组，则它必须是可寻址的，并且切片操作的结果是与数组具有相同元素类型的切片。

如果有效切片表达式的切片操作数是 nil 切片，则结果也是 nil 切片。否则，如果结果是切片，则它与操作数共享其基础数组。

```go
var a [10]int
s1 := a[3:7]   // underlying array of s1 is array a; &s1[2] == &a[5]
s2 := s1[1:4]  // underlying array of s2 is underlying array of s1 which is array a; &s2[1] == &a[5]
s2[1] = 42     // s2[1] == s1[2] == a[5] == 42; they all refer to the same underlying array element

var s []int
s3 := s[:0]    // s3 == nil
```

### 完整切片表达式 Full slice expressions

主要表达式

```go
a[low : high : max]
```

构造一个与简单切片表达式 a[low : high] 相同类型、相同长度和元素的切片。此外，它还通过将结果设置为 max-low 来控制结果切片的容量。仅第一个索引可以省略；默认为 0。 a 的核心类型必须是数组、数组指针或切片（但不是字符串）。切片数组 a 后

```go
a := [5]int{1, 2, 3, 4, 5}
t := a[1:3:5]
```

切片 t 具有类型 []int、长度 2、容量 4 和元素

```go
t[0] == 2
t[1] == 3
```

对于简单的切片表达式，如果 a 是指向数组的指针，则 a[low : high : max] 是 (*a)[low : high : max] 的简写。如果切片操作数是数组，则它必须是可寻址的。

如果 0 <= low <= high <= max <= cap(a)，则索引在范围内，否则索引超出范围。常量索引必须是非负的并且可以用 int 类型的值表示；对于数组，常量索引也必须在范围内。如果多个索引是常数，则存在的常数必须在相对于彼此的范围内。如果索引在运行时超出范围，则会发生运行时恐慌。

## 类型断言 Type assertions

对于接口类型的表达式 x，但不是类型参数，且类型为 T，主表达式

```go
x.(T)
```

断言 x 不为 nil，并且存储在 x 中的值的类型为 T。符号 x.(T) 称为类型断言。

更准确地说，如果 T 不是接口类型，则 x.(T) 断言 x 的动态类型与类型 T 相同。在这种情况下，T 必须实现 x 的（接口）类型；否则，类型断言无效，因为 x 不可能存储类型 T 的值。如果 T 是接口类型，则 x.(T) 断言 x 的动态类型实现接口 T。

如果类型断言成立，则表达式的值是存储在 x 中的值，并且其类型为 T。如果类型断言为 false，则会发生运行时恐慌。换句话说，即使 x 的动态类型仅在运行时已知，但在正确的程序中 x.(T) 的类型已知为 T。

```go
var x interface{} = 7          // x has dynamic type int and value 7
i := x.(int)                   // i has type int and value 7

type I interface { m() }

func f(y I) {
    s := y.(string)        // illegal: string does not implement I (missing method m)
    r := y.(io.Reader)     // r has type io.Reader and the dynamic type of y must implement both I and io.Reader
    …
}
```

在赋值语句或特殊形式的初始化中使用的类型断言

```go
v, ok = x.(T)
v, ok := x.(T)
var v, ok = x.(T)
var v, ok interface{} = x.(T) // dynamic types of v and ok are T and bool
```

产生一个额外的无类型布尔值。如果断言成立，则 ok 的值为 true。否则为 false，并且 v 的值是类型 T 的零值。在这种情况下不会发生运行时恐慌。

## 调用 Calls

给定一个具有函数类型核心类型 F 的表达式 f，

```go
f(a1, a2, … an)
```

使用参数 a1、a2、... an 调用 f。除一种特殊情况外，参数必须是可分配给 F 的参数类型的单值表达式，并在调用函数之前进行计算。表达式的类型是 F 的结果类型。方法调用类似，但方法本身被指定为该方法的接收器类型值的选择器。

```go
math.Atan2(x, y)  // function call
var pt *Point
pt.Scale(3.5)     // method call with receiver pt
```

如果 f 表示泛型函数，则必须先实例化它，然后才能调用它或将其用作函数值。

在函数调用中，函数值和参数按通常的顺序求值。在评估它们之后，调用的参数将按值传递给函数，并且被调用的函数开始执行。当函数返回时，函数的返回参数按值传递回调用者。

调用 nil 函数值会导致运行时恐慌。

作为一种特殊情况，如果函数或方法 g 的返回值数量相等并且可单独分配给另一个函数或方法 f 的参数，则调用 f(g(parameters_of_g)) 将在绑定返回值后调用 f g 的参数按顺序传递给 f 的参数。 f 的调用除了 g 的调用之外不能包含任何参数，并且 g 必须至少有一个返回值。如果 f 有一个最终的 ... 参数，则为它分配常规参数分配后保留的 g 的返回值。

```go
func Split(s string, pos int) (string, string) {
    return s[0:pos], s[pos:]
}

func Join(s, t string) string {
    return s + t
}

if Join(Split(value, len(value)/2)) != value {
    log.Panic("test fails")
}
```

如果 x（的类型）的方法集包含 m 并且参数列表可以分配给 m 的参数列表，则方法调用 x.m() 是有效的。如果 x 是可寻址的并且 &x 的方法集包含 m，则 x.m() 是 (&x).m() 的简写：

```go
var p Point
p.Scale(3.5)
```

没有明显的方法类型，也没有方法文字。

## 将参数传递给... Passing arguments to `...` parameters

如果 f 是可变参数，最终参数 p 类型为 ...T，则在 f 中 p 的类型相当于 []T 类型。如果调用 f 时没有 p 的实际参数，则传递给 p 的值为 nil。否则，传递的值是一个 []T 类型的新切片，带有一个新的基础数组，其连续元素是实际参数，所有参数都必须可分配给 T。因此，切片的长度和容量是绑定到的参数数量p 并且对于每个调用站点可能有所不同。

给定函数和调用

```go
func Greeting(prefix string, who ...string)
Greeting("nobody")
Greeting("hello:", "Joe", "Anna", "Eileen")
```

在 Greeting 中，who 在第一次调用中的值为 nil，在第二次调用中的值为 []string{"Joe", "Anna", "Eileen"}。

如果最后一个参数可分配给切片类型 []T 并且后跟 ...，则它会作为 ...T 参数的值原封不动地传递。在这种情况下，不会创建新切片。

给定切片 s 并调用

```go
s := []string{"James", "Jasmine"}
Greeting("goodbye:", s...)
```

在 Greeting 中，who 将具有与 s 相同的值，并具有相同的底层数组。
