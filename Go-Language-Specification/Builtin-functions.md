# 内置函数 Built-in functions

内置函数是预先声明的。它们的调用方式与任何其他函数一样，但其中一些函数接受类型而不是表达式作为第一个参数。
内置函数没有标准的 Go 类型，因此只能出现在调用表达式中；它们不能用作函数值。

## 追加和复制切片 Appending to and Copying slices

内置函数 append 和 copy 可协助进行常见的切片操作。这两个函数的结果与参数引用的内存是否重叠无关。

可变参数函数append将零个或多个值x附加到切片s并返回与s类型相同的结果切片。 s 的核心类型必须是 []E 类型的切片。值 x 被传递给 ...E 类型的参数，并应用相应的参数传递规则。作为一种特殊情况，如果 s 的核心类型是 []byte，append 还接受第二个参数，核心类型为 bytestring，后跟 ...。这种形式附加字节切片或字符串的字节。

```go
append(s S, x ...E) S  // core type of S is []E
```

如果 s 的容量不足以容纳附加值，append 会分配一个新的、足够大的底层数组，该数组既适合现有切片元素又适合附加值。否则，追加会重新使用底层数组。

```go
s0 := []int{0, 0}
s1 := append(s0, 2)                // append a single element     s1 is []int{0, 0, 2}
s2 := append(s1, 3, 5, 7)          // append multiple elements    s2 is []int{0, 0, 2, 3, 5, 7}
s3 := append(s2, s0...)            // append a slice              s3 is []int{0, 0, 2, 3, 5, 7, 0, 0}
s4 := append(s3[3:6], s3[2:]...)   // append overlapping slice    s4 is []int{3, 5, 7, 2, 3, 5, 7, 0, 0}

var t []interface{}
t = append(t, 42, 3.1415, "foo")   //                             t is []interface{}{42, 3.1415, "foo"}

var b []byte
b = append(b, "bar"...)            // append string contents      b is []byte{'b', 'a', 'r' }
```

函数 copy 将片元从源 src 复制到目标 dst，并返回复制的片元个数。两个参数的核心类型必须是元素类型相同的切片。复制的元素个数是 len(src) 和 len(dst) 的最小值。作为特例，如果目的地的核心类型是 []字节，copy 也接受核心类型为字节的源参数。这种形式会将字节片或字符串中的字节复制到字节片中。

```go
copy(dst, src []T) int
copy(dst []byte, src string) int
```

栗子：

```go
var a = [...]int{0, 1, 2, 3, 4, 5, 6, 7}
var s = make([]int, 6)
var b = make([]byte, 5)
n1 := copy(s, a[0:])            // n1 == 6, s is []int{0, 1, 2, 3, 4, 5}
n2 := copy(s, s[2:])            // n2 == 4, s is []int{2, 3, 4, 5, 4, 5}
n3 := copy(b, "Hello, World!")  // n3 == 5, b is []byte("Hello")
```

## Clear

内置函数clear接受map、slice或type参数类型的参数，并删除或清零所有元素。

```text
Call        Argument type     Result

clear(m)    map[K]T           deletes all entries, resulting in an
                              empty map (len(m) == 0)

clear(s)    []T               sets all elements up to the length of
                              s to the zero value of T

clear(t)    type parameter    see below
```

如果要清除的参数类型是类型参数，则其类型集中的所有类型都必须是映射或切片，并且清除执行与实际类型参数相对应的操作。

如果映射或切片为零，则清除是无操作。

## Close

对于核心类型为通道的参数 ch，内置函数 close 会记录通道上不再发送的值。如果 ch 是只接收通道，则会出错。向关闭的通道发送或关闭通道会导致运行时恐慌。关闭 nil 通道也会导致运行时恐慌。调用关闭后，在接收到任何先前发送的值后，接收操作将返回通道类型的零值，而不会阻塞。多值接收操作会返回一个接收值以及通道是否已关闭的指示。

## 复数操作 Manipulating complex numbers

三个函数组合和分解复数。内置函数 complex 从浮点实部和虚部构造复数值，而 real 和 imag 则提取复数值的实部和虚部。

```
complex(realPart, imaginaryPart floatT) complexT
real(complexT) floatT
imag(complexT) floatT
```

参数类型和返回值相对应。对于复数，两个参数必须是相同的浮点类型，返回类型是具有相应浮点成分的复数类型：float32 参数为 complex64，float64 参数为 complex128。如果其中一个参数的值为非类型常量，则首先隐式转换为另一个参数的类型。如果两个参数都求值为未键入常量，则它们必须是非复数或其虚部必须为零，函数的返回值是一个未键入的复数常量。

对于 real 和 imag，参数必须是复数类型，并且返回类型是相应的浮点类型：对于complex64 参数为float32，对于complex128 参数为float64。如果参数计算结果为无类型常量，则它必须是数字，并且函数的返回值是无类型浮点常量。

real 和 imag 函数共同构成了complex 的反函数，因此对于复数类型Z 的值z，z === Z(complex(real(z), imag(z)))。

如果这些函数的操作数都是常量，则返回值也是常量。

```go
var a = complex(2, -2)             // complex128
const b = complex(1.0, -1.4)       // untyped complex constant 1 - 1.4i
x := float32(math.Cos(math.Pi/2))  // float32
var c64 = complex(5, -x)           // complex64
var s int = complex(1, 0)          // untyped complex constant 1 + 0i can be converted to int
_ = complex(1, 2<<s)               // illegal: 2 assumes floating-point type, cannot shift
var rl = real(c64)                 // float32
var im = imag(a)                   // float64
const c = imag(b)                  // untyped constant -1.4
_ = imag(3 << s)                   // illegal: 3 assumes complex type, cannot shift
```

不允许使用类型形参类型的实参。

## 删除map元素 Deletion of map elements

内置函数delete从映射m中删除带有键k的元素。值 k 必须可分配给 m 的键类型。

```go
delete(m, k)  // remove element m[k] from map m
```

如果 m 的类型是类型参数，则该类型集中的所有类型都必须是映射，并且它们都必须具有相同的键类型。

如果映射 m 为 nil 或元素 m[k] 不存在，则删除是无操作的。



## 长度与容量 Length and capacity

内置函数 len 和 cap 接受各种类型的参数并返回 int 类型的结果。该实现保证结果始终适合 int。

```textile
Call      Argument type    Result

len(s)    string type      string length in bytes
          [n]T, *[n]T      array length (== n)
          []T              slice length
          map[K]T          map length (number of defined keys)
          chan T           number of elements queued in channel buffer
          type parameter   see below

cap(s)    [n]T, *[n]T      array length (== n)
          []T              slice capacity
          chan T           channel buffer capacity
          type parameter   see below
```

如果实参类型是类型参数 P，则调用 len(e)（或分别是 cap(e)）对于 P 类型集中的每种类型都必须有效。结果是参数的长度（或容量），该参数的类型对应于实例化 P 的类型参数。

切片的容量是底层数组中为其分配空间的元素数量。在任何时候，以下关系都成立：

```go
0 <= len(s) <= cap(s)
```

nil 切片、映射或通道的长度为 0。nil 切片或通道的容量为 0。

如果 s 是字符串常量，表达式 len(s) 就是常量。如果 s 的类型是数组或数组指针，且表达式 s 不包含通道接收或（非常数）函数调用，则表达式 len(s) 和 cap(s) 为常数；在这种情况下，s 不会被求值。否则，len 和 cap 的调用不是常量，s 将被求值。

```go
const (
	c1 = imag(2i)                    // imag(2i) = 2.0 is a constant
	c2 = len([10]float64{2})         // [10]float64{2} contains no function calls
	c3 = len([10]float64{c1})        // [10]float64{c1} contains no function calls
	c4 = len([10]float64{imag(2i)})  // imag(2i) is a constant and no function call is issued
	c5 = len([10]float64{imag(z)})   // invalid: imag(z) is a (non-constant) function call
)
var z complex128
```

## Making slices, maps and channels

内置函数 make 接受一个 T 类型，可选择在该类型后接一个特定类型的表达式列表。T 的核心类型必须是片段、映射或通道。它返回一个 T 类型的值（不是 *T）。内存的初始化如初始值一节所述。

```
Call             Core type    Result

make(T, n)       slice        slice of type T with length n and capacity n
make(T, n, m)    slice        slice of type T with length n and capacity m

make(T)          map          map of type T
make(T, n)       map          map of type T with initial space for approximately n elements

make(T)          channel      unbuffered channel of type T
make(T, n)       channel      buffered channel of type T, buffer size n
```

每个大小参数 n 和 m 必须是整数类型，具有仅包含整数类型的类型集，或者是无类型常量。常量大小参数必须是非负的并且可以用 int 类型的值表示；如果它是无类型常量，则其类型为 int。如果 n 和 m 都提供并且为常数，则 n 必须不大于 m。对于切片和通道，如果运行时 n 为负或大于 m，则会发生运行时恐慌。

```go
s := make([]int, 10, 100)       // slice with len(s) == 10, cap(s) == 100
s := make([]int, 1e3)           // slice with len(s) == cap(s) == 1000
s := make([]int, 1<<63)         // illegal: len(s) is not representable by a value of type int
s := make([]int, 10, 0)         // illegal: len(s) > cap(s)
c := make(chan int, 10)         // channel with a buffer size of 10
m := make(map[string]int, 100)  // map with initial space for approximately 100 elements
```

使用映射类型和大小提示 n 调用 make 将创建一个具有可容纳 n 个映射元素的初始空间的映射。精确的行为取决于实现。

## Min and max

内置函数 min 和 max 分别计算固定数量的有序类型参数的最小值或最大值。必须至少有一个参数。

与运算符相同的类型规则适用：对于有序参数 x 和 y，如果 x + y 有效，则 min(x, y) 有效，并且 min(x, y) 的类型是 x + y 的类型（并且类似地对于最大值）。如果所有参数都是常数，则结果也是常数。

```go
var x, y int
m := min(x)                 // m == x
m := min(x, y)              // m is the smaller of x and y
m := max(x, y, 10)          // m is the larger of x and y but at least 10
c := max(1, 2.0, 10)        // c == 10.0 (floating-point kind)
f := max(0, float32(x))     // type of f is float32
var s []string
_ = min(s...)               // invalid: slice arguments are not permitted
t := max("", "foo", "bar")  // t == "foo" (string kind)
```

对于数字参数，假设所有 NaN 都相等，则 min 和 max 是可交换和结合的：

```go
min(x, y)    == min(y, x)
min(x, y, z) == min(min(x, y), z) == min(x, min(y, z))
```

对于浮点参数负零、NaN 和无穷大，适用以下规则：

```go
   x        y    min(x, y)    max(x, y)

  -0.0    0.0         -0.0          0.0    // negative zero is smaller than (non-negative) zero
  -Inf      y         -Inf            y    // negative infinity is smaller than any other number
  +Inf      y            y         +Inf    // positive infinity is larger than any other number
   NaN      y          NaN          NaN    // if any argument is a NaN, the result is a NaN
```

对于字符串参数，min 的结果是具有最小（或 max，最大）值的第一个参数，按词法字节进行比较：

```go
min(x, y)    == if x <= y then x else y
min(x, y, z) == min(min(x, y), z)
```

## Allocation

内置函数 new 采用类型 T，在运行时为该类型的变量分配存储空间，并返回指向它的 *T 类型的值。该变量按照初始值部分中的描述进行初始化。

```go
new(T)
```

例如:

```go
type S struct { a int; b float64 }
new(S)
```

为 S 类型的变量分配存储空间，对其进行初始化（a=0，b=0.0），并返回一个包含该位置地址的 *S 类型值。

## Handling panics

两个内置函数，即 panic 和 recover，可帮助报告和处理运行时 panic 和程序定义的错误条件。

```go
func panic(interface{})
func recover() interface{}
```

在执行函数 F 时，显式调用 panic 或运行时 panic 会终止 F 的执行。接下来，由 F 的调用者运行的任何延迟函数都会被运行，依此类推，直到执行中的 goroutine 的顶层函数所延迟的任何函数。此时，程序被终止并报告错误条件，包括 panic 参数的值。这个终止序列称为 panic。

```go
panic(42)
panic("unreachable")
panic(Error("cannot parse"))
```

recover函数允许程序管理恐慌 goroutine 的行为。假设函数 G 推迟了调用恢复的函数 D，并且与 G 执行的同一 goroutine 上的函数中发生了恐慌。当延迟函数运行到D时，D调用recover的返回值将是传递给panic调用的值。如果 D 正常返回，而没有引发新的恐慌，则恐慌序列将停止。在这种情况下，G 和恐慌调用之间调用的函数的状态将被丢弃，并恢复正常执行。然后，G 在 D 之前推迟的任何函数都会运行，并且 G 的执行将通过返回到其调用者而终止。

当goroutine没有panic或者recover没有被延迟函数直接调用时，recover的返回值为nil。相反，如果 goroutine 发生恐慌并且由延迟函数直接调用recover，则recover的返回值保证不会为nil。为了确保这一点，使用 nil 接口值（或无类型 nil）调用panic会导致运行时panic。

下面示例中的 protected 函数调用函数参数 g 并保护调用者免受 g 引发的运行时恐慌的影响。

```go
func protect(g func()) {
	defer func() {
		log.Println("done")  // Println executes normally even if there is a panic
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
		}
	}()
	log.Println("start")
	g()
}
```

## Bootstrapping

当前的实现提供了几个在引导期间有用的内置函数。这些函数的记录是为了完整性，但不保证保留在该语言中。他们不返回结果。

```go
Function   Behavior

print      prints all arguments; formatting of arguments is implementation-specific
println    like print but prints spaces between arguments and a newline at the end
```

实施限制：print 和 println 不必接受任意参数类型，但必须支持布尔、数字和字符串类型的打印。
