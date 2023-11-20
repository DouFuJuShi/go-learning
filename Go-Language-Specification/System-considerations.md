# 系统注意事项 System considerations

## unsafe包 Package `unsafe`

内置软件包 unsafe 是编译器已知的，可通过导入路径 "unsafe "访问，它为低级编程提供了便利，包括违反类型系统的操作。使用 unsafe 的软件包必须经过手动类型安全审查，而且可能无法移植。软件包提供以下接口：

```go
package unsafe

type ArbitraryType int  // shorthand for an arbitrary Go type; it is not a real type
type Pointer *ArbitraryType

func Alignof(variable ArbitraryType) uintptr
func Offsetof(selector ArbitraryType) uintptr
func Sizeof(variable ArbitraryType) uintptr

type IntegerType int  // shorthand for an integer type; it is not a real type
func Add(ptr Pointer, len IntegerType) Pointer
func Slice(ptr *ArbitraryType, len IntegerType) []ArbitraryType
func SliceData(slice []ArbitraryType) *ArbitraryType
func String(ptr *byte, len IntegerType) string
func StringData(str string) *byte
```

指针是一种指针类型，但指针值不能被取消引用。任何底层类型为 uintptr 的指针或值都可以转换为底层类型为 Pointer 的类型，反之亦然。在 Pointer 和 uintptr 之间转换的效果由实现定义。

```go
var f float64
bits = *(*uint64)(unsafe.Pointer(&f))

type ptr unsafe.Pointer
bits = *(*uint64)(ptr(&f))

var p ptr = nil
```

函数 Alignof 和 Sizeof 接受任何类型的表达式 x 并分别返回假设变量 v 的对齐方式或大小，就好像 v 是通过 var v = x 声明的一样。

函数 Offsetof 接受一个（可能带括号的）选择器 s.f，表示由 s 或 *s 表示的结构体的字段 f，并返回相对于结构体地址的字段偏移量（以字节为单位）。如果 f 是嵌入字段，则必须无需通过结构体字段进行指针间接访问即可访问它。对于具有字段 f 的结构体 s：

```go
uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f) == uintptr(unsafe.Pointer(&s.f))
```

计算机体系结构可能需要对齐内存地址；也就是说，变量的地址是因子的倍数，即变量的类型对齐方式。 Alignof 函数采用表示任意类型变量的表达式，并返回变量（类型）的对齐方式（以字节为单位）。对于变量 x：

```go
uintptr(unsafe.Pointer(&x)) % unsafe.Alignof(x) == 0
```

函数 Add 将 len 添加到 ptr 并返回更新后的指针 unsafe.Pointer(uintptr(ptr) + uintptr(len))。 len 参数必须是整数类型或无类型常量。常量 len 参数必须可由 int 类型的值表示；如果它是无类型常量，则其类型为 int。有效使用指针的规则仍然适用。

函数 Slice 返回一个切片，其底层数组从 ptr 开始，长度和容量为 len。 Slice(ptr, len) 相当于

```go
(*[len]ArbitraryType)(unsafe.Pointer(ptr))[:]
```

但作为一种特殊情况，如果 ptr 为 nil 并且 len 为零，则 Slice 返回 nil。

len 参数必须是整数类型或无类型常量。常量 len 参数必须是非负的并且可以用 int 类型的值表示；如果它是无类型常量，则其类型为 int。在运行时，如果 len 为负，或者 ptr 为 nil 并且 len 不为零，则会发生运行时恐慌。
函数 SliceData 返回一个指向切片参数的底层数组的指针。如果切片的容量上限(slice)不为零，则该指针为&slice[:1][0]。如果 slice 为零，则结果为零。否则它是一个指向未指定内存地址的非零指针。
函数 String 返回一个字符串值，其底层字节从 ptr 开始，长度为 len。与函数 Slice 中一样，同样的要求也适用于 ptr 和 len 参数。如果 len 为零，则结果是空字符串“”。由于 Go 字符串是不可变的，因此传递给 String 的字节之后不得修改。
函数 StringData 返回一个指向 str 参数的底层字节的指针。对于空字符串，返回值未指定，并且可能为 nil。由于 Go 字符串是不可变的，因此 StringData 返回的字节不得修改。

## 尺寸和对齐保证 Size and alignment guarantees

对于数字类型，保证以下大小：

```go
type                                 size in bytes

byte, uint8, int8                     1
uint16, int16                         2
uint32, int32, float32                4
uint64, int64, float64, complex64     8
complex128                           16
```

保证以下最小对齐属性：

1. 对于任何类型的变量 x：unsafe.Alignof(x) 至少为 1。

2. 对于结构体类型的变量 x：unsafe.Alignof(x) 是 x 的每个字段 f 的所有值 unsafe.Alignof(x.f) 中最大的，但至少为 1。

3. 对于数组类型的变量 x： unsafe.Alignof(x) 与数组元素类型的变量的对齐方式相同。

如果结构体或数组类型不包含大小大于零的字段（或元素），则其大小为零。两个不同的零大小变量在内存中可能具有相同的地址。
