# 语句 statement

语句控制执行。

```ebnf
Statement =
	Declaration | LabeledStmt | SimpleStmt |
	GoStmt | ReturnStmt | BreakStmt | ContinueStmt | GotoStmt |
	FallthroughStmt | Block | IfStmt | SwitchStmt | SelectStmt | ForStmt |
	DeferStmt .

SimpleStmt = EmptyStmt | ExpressionStmt | SendStmt | IncDecStmt | Assignment | ShortVarDecl .
```

## 中断语句 Terminating statements

终止语句中断块中的常规控制流。以下语句终止：

1. “return”或“goto”语句。

2. 调用内置函数 panic。

3. 语句列表以终止语句结束的块。

4. “if”语句，其中：
   
   - “else”分支存在，并且
   
   - 两个分支都是终止语句。

5. “for”语句，其中：
   
   - 没有引用 "for "语句的 "break "语句，以及
   
   - 循环条件不存在，并且
   
   - “for”语句不使用范围子句。

6. “switch”语句，其中：
   
   - 在 "switch "语句中没有 "break "语句、
   
   - 有一个默认情况，并且
   
   - 该语句列出了每种情况，包括默认情况、以终止语句结束或可能标记为“fallthrough”的语句。

7. “select”语句，其中：
   
   - 没有引用“select”语句的“break”语句，并且
   
   - 该语句列出了每种情况，包括默认情况（如果存在），以终止语句结束。

8. 标记语句标记终止语句。

所有其他语句都不会终止。
如果语句列表不为空，且其最后一条非空语句为终止语句，则该语句列表以终止语句结束。

## 空语句 Empty statements

空语句什么也不做。

```ebnf
EmptyStmt = .
```

## 标签语句 labeled statements

标签语句可以是 goto、break 或 continue 语句的目标。

```ebnf
LabeledStmt = Label ":" Statement .
Label       = identifier .
```

```go
Error: log.Panic("error encountered")
```

## 表达式语句 Expression statements

除了特定的内置函数外，函数和方法调用以及接收操作可以出现在语句上下文中。此类语句可加括号。

```ebnf
ExpressionStmt = Expression .
```

语句上下文中不允许使用以下内置函数：

```go
append cap complex imag len make new real
unsafe.Add unsafe.Alignof unsafe.Offsetof unsafe.Sizeof unsafe.Slice unsafe.SliceData unsafe.String unsafe.StringData
```

```go
h(x+y)
f.Close()
<-ch
(<-ch)
len("foo")  // illegal if len is the built-in function
```

## 发送语句 Send statements

发送语句在通道上发送一个值。通道表达式的核心类型必须是通道，通道方向必须允许发送操作，并且要发送的值的类型必须可分配给通道的元素类型。

```ebnf
SendStmt = Channel "<-" Expression .
Channel  = Expression .
```

在通信开始之前，将对通道和值表达式进行评估。通信将阻塞，直到发送可以继续。如果接收器准备就绪，则可以在无缓冲通道上进行发送。如果缓冲区中有空间，则可以继续在缓冲通道上发送。关闭通道上的发送会导致运行时恐慌。零通道上的发送将永远阻塞。

```go
ch <- 3  // send value 3 to channel ch
```

## 自增自减 IncDec statements

“++”和“--”语句通过无类型常量 1 递增或递减其操作数。与赋值一样，操作数必须是可寻址的或映射索引表达式。

```go
IncDecStmt = Expression ( "++" | "--" ) .
```

以下赋值语句在语义上是等效的：

```go
IncDec statement    Assignment
x++                 x += 1
x--                 x -= 1
```

## 赋值语句 Assignment statements

赋值将变量中存储的当前值替换为表达式指定的新值。赋值语句可以将单个值分配给单个变量，或者将多个值分配给匹配数量的变量。

```ebnf
Assignment = ExpressionList assign_op ExpressionList .

assign_op = [ add_op | mul_op ] "=" .
```

每个左侧操作数必须是可寻址的、映射索引表达式或（仅适用于 = 赋值）空白标识符。操作数可以用括号括起来。

```go
x = 1
*p = f()
a[i] = 23
(k) = <-ch  // same as: k = <-ch
```

赋值运算 x op= y（其中 op 是二元算术运算符）等效于 x = x op (y)，但仅对 x 求值一次。 op= 构造是单个标记。在赋值运算中，左侧表达式列表和右侧表达式列表都必须恰好包含一个单值表达式，并且左侧表达式不能是空白标识符。

```go
a[i] <<= 2
i &^= 1<<n
```

元组赋值将多值运算的各个元素分配给变量列表。有两种形式。在第一个中，右侧操作数是单个多值表达式，例如函数调用、通道或映射操作或类型断言。左侧操作数的数量必须与值的数量匹配。例如，如果 f 是一个返回两个值的函数，

```go
x, y = f()
```

将第一个值赋给 x，将第二个值赋给 y。在第二种形式中，左侧操作数的数量必须等于右侧表达式的数量，每个表达式都必须是单值，并且右侧的第 n 个表达式被分配给左侧的第 n 个操作数：

```go
one, two, three = '一', '二', '三'
```

空白标识符提供了一种忽略赋值中右侧值的方法：

```go
_ = x       // evaluate x but ignore it
x, _ = f()  // evaluate f() but ignore second result value
```

任务分两个阶段进行。首先，左边的索引表达式和指针间接（包括选择器中的隐式指针间接）的操作数和右边的表达式都按通常的顺序求值。其次，作业按从左到右的顺序进行。

```go
a, b = b, a  // exchange a and b

x := []int{1, 2, 3}
i := 0
i, x[i] = 1, 2  // set i = 1, x[0] = 2

i = 0
x[i], i = 2, 1  // set x[0] = 2, i = 1

x[0], x[0] = 1, 2  // set x[0] = 1, then x[0] = 2 (so x[0] == 2 at end)

x[1], x[3] = 4, 5  // set x[1] = 4, then panic setting x[3] = 5.

type Point struct { x, y int }
var p *Point
x[2], p.x = 6, 7  // set x[2] = 6, then panic setting p.x = 7

i = 2
x = []int{3, 5, 7}
for i, x[i] = range x {  // set i, x[2] = 0, x[0]
	break
}
// after this loop, i == 0 and x is []int{3, 5, 3}
```

在赋值中，每个值都必须可分配给其所分配到的操作数的类型，但有以下特殊情况：

1. 任何键入的值都可以分配给空白标识符。

2. 如果将无类型常量分配给接口类型或空白标识符的变量，则该常量首先会隐式转换为其默认类型。

3. 如果将无类型布尔值分配给接口类型或空白标识符的变量，则它首先会隐式转换为 bool 类型。

## if 语句 if statements

“If”语句根据布尔表达式的值指定两个分支的条件执行。如果表达式计算结果为真，则执行“if”分支，否则，如果存在，则执行“else”分支。

```ebnf
IfStmt = "if" [ SimpleStmt ";" ] Expression Block [ "else" ( IfStmt | Block ) ] .
```

```go
if x > max {
	x = max
}
```

表达式前面可以有一个简单的语句，该语句在计算表达式之前执行。

```go
if x := f(); x < y {
	return x
} else if x > z {
	return z
} else {
	return y
}
```

## Switch语句 Switch statements

“Switch”语句提供多路执行。将表达式或类型与“switch”内的“cases”进行比较，以确定要执行哪个分支。

```ebnf
SwitchStmt = ExprSwitchStmt | TypeSwitchStmt .
```

有两种形式：表达式switch和类型switch。在表达式 switch 中，case 包含与 switch 表达式的值进行比较的表达式。在类型switch中，cases 包含与特殊注释的switc表达式的类型进行比较的类型。 switch 表达式在 switch 语句中只计算一次。

### 表达式switch Expression switches

在表达式 switch 中，对 switch 表达式进行求值，而 case 表达式（不必是常量）按从左到右、从上到下的顺序求值；第一个等于 switch 表达式的触发执行关联 case 的语句；其他情况被跳过。如果没有 case 匹配并且存在“default”case，则执行其语句。最多可以有一种默认情况，它可以出现在“switch”语句中的任何位置。缺少 switch 表达式相当于布尔值 true。

```go
ExprSwitchStmt = "switch" [ SimpleStmt ";" ] [ Expression ] "{" { ExprCaseClause } "}" .
ExprCaseClause = ExprSwitchCase ":" StatementList .
ExprSwitchCase = "case" ExpressionList | "default" .
```

如果 switch 表达式的值是一个未类型化的常量，那么首先会将其隐式转换[converted](https://go.dev/ref/spec#Conversions)为默认类型([default type](https://go.dev/ref/spec#Constants))。预先声明的非类型值 nil 不能用作 switch 表达式。开关表达式的类型必须是可比较([comparable](https://go.dev/ref/spec#Comparison_operators))的。

如果 case 表达式是无类型的，则它首先会隐式转换为 switch 表达式的类型。对于每个（可能已转换的）case 表达式 x 和 switch 表达式的值 t，x == t 必须是有效的比较。

换句话说，switch 表达式被视为用于声明和初始化临时变量 t，而无需显式类型；它是 t 的值，每个 case 表达式 x 都根据该值进行相等性测试。

在 case 或 default 子句中，最后一个非空语句可能是（可能标记为）“fallthrough”语句，以指示控制应从本子句的末尾流到下一个子句的第一个语句。否则控制流到“switch”语句的末尾。 “fallthrough”语句可能会作为表达式 switch 的最后一个子句之外的所有语句的最后一个语句出现。

switch 表达式前面可以有一个简单的语句，该语句在表达式求值之前执行。

```go
switch tag {
default: s3()
case 0, 1, 2, 3: s1()
case 4, 5, 6, 7: s2()
}

switch x := f(); {  // missing switch expression means "true"
case x < 0: return -x
default: return x
}

switch {
case x < y: f1()
case x < z: f2()
case x == 4: f3()
}
```

实现限制：编译器可能不允许多个 case 表达式求值相同的常量。例如，当前的编译器不允许在 case 表达式中使用重复的整数、浮点或字符串常量。

### 类型switch Type switch

类型开关比较类型而不是值。它在其他方面类似于表达式开关。它由一个特殊的 switch 表达式标记，该表达式具有使用关键字 type 而不是实际类型的类型断言([type assertion](https://go.dev/ref/spec#Type_assertions))的形式：

```go
switch x.(type) {
// cases
}
```

然后，案例将实际类型 T 与表达式 x 的动态类型进行匹配。与类型断言一样，x 必须是接口类型，但不是类型参数，并且案例中列出的每个非接口类型 T 必须实现 x 的类型。类型开关中列出的类型必须全部不同。

```ebnf
TypeSwitchStmt  = "switch" [ SimpleStmt ";" ] TypeSwitchGuard "{" { TypeCaseClause } "}" .
TypeSwitchGuard = [ identifier ":=" ] PrimaryExpr "." "(" "type" ")" .
TypeCaseClause  = TypeSwitchCase ":" StatementList .
TypeSwitchCase  = "case" TypeList | "default" .
```

TypeSwitchGuard 可以包括简短的变量声明。使用该形式时，变量在每个子句的隐式块中的 TypeSwitchCase 末尾声明。在 case 只列出一种类型的子句中，变量具有该类型；否则，变量具有 TypeSwitchGuard 中表达式的类型。

case 可以使用预先声明的标识符 nil 来代替类型；当 TypeSwitchGuard 中的表达式为 nil 接口值时，会选择这种情况。最多可能存在一种零情况。

给定接口类型为interface{}的表达式 x，可进行以下类型转换：

```go
switch i := x.(type) {
case nil:
	printString("x is nil")                // type of i is type of x (interface{})
case int:
	printInt(i)                            // type of i is int
case float64:
	printFloat64(i)                        // type of i is float64
case func(int) float64:
	printFunction(i)                       // type of i is func(int) float64
case bool, string:
	printString("type is bool or string")  // type of i is type of x (interface{})
default:
	printString("don't know the type")     // type of i is type of x (interface{})
}
```

可以重写：

```go
v := x  // x is evaluated exactly once
if v == nil {
	i := v                                 // type of i is type of x (interface{})
	printString("x is nil")
} else if i, isInt := v.(int); isInt {
	printInt(i)                            // type of i is int
} else if i, isFloat64 := v.(float64); isFloat64 {
	printFloat64(i)                        // type of i is float64
} else if i, isFunc := v.(func(int) float64); isFunc {
	printFunction(i)                       // type of i is func(int) float64
} else {
	_, isBool := v.(bool)
	_, isString := v.(string)
	if isBool || isString {
		i := v                         // type of i is type of x (interface{})
		printString("type is bool or string")
	} else {
		i := v                         // type of i is type of x (interface{})
		printString("don't know the type")
	}
}
```

类型参数或泛型类型可以用作情况中的类型。如果在实例化时该类型与交换机中的另一个条目重复，则选择第一个匹配的情况。

```go
func f[P any](x any) int {
	switch x.(type) {
	case P:
		return 0
	case string:
		return 1
	case []P:
		return 2
	case []byte:
		return 3
	default:
		return 4
	}
}

var v1 = f[string]("foo")   // v1 == 0
var v2 = f[byte]([]byte{})  // v2 == 2
```

类型开关防护前面可以有一个简单的语句，该语句在评估防护之前执行。

类型开关中不允许使用“fallthrough”语句。

## For语句 For statements

“for”语句指定重复执行块。有三种形式： 迭代可以由单个条件、“for”子句或“range”子句控制。

```go
ForStmt = "for" [ Condition | ForClause | RangeClause ] Block .
Condition = Expression .
```

### 对于具有单一条件的语句 For statements with single condition

在最简单的形式中，"for "语句指定只要布尔条件求值为 "true"，就重复执行一个程序块。每次迭代前都要对该条件进行评估。如果条件不存在，则相当于布尔值为 true。

```go
for a < b {
	a *= 2
}
```

### 对于带有 for 子句的语句 For statements with for clause

带有 ForClause 的 "for "语句也受其条件控制，但它还可以指定 init 和 post 语句，如赋值、递增或递减语句。init 语句可以是简短的变量声明，但 post 语句则不可以。init 语句声明的变量会在每次迭代中重复使用。

```ebnf
ForClause = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
InitStmt = SimpleStmt .
PostStmt = SimpleStmt .
```

```go
for i := 0; i < 10; i++ {
	f(i)
}
```

如果非空，则在评估第一次迭代的条件之前执行一次 init 语句； post 语句在每次执行块后执行（并且仅当块被执行时）。 ForClause 的任何元素都可以为空，但除非只有一个条件，否则需要分号。如果条件不存在，则相当于布尔值 true。

```go
for cond { S() }  //  is the same as    for ; cond ; { S() }
for      { S() }  //  is the same as    for true     { S() }
```

### 对于带有range子句的语句 For statements with range clause

带有“range”子句的“for”语句会迭代数组、切片、字符串或映射的所有条目，或者通道上接收到的值。对于每个条目，它会将迭代值分配给相应的迭代变量（如果存在），然后执行该块。

```ebnf
RangeClause = [ ExpressionList "=" | IdentifierList ":=" ] "range" Expression .
```

“range”子句右侧的表达式称为范围表达式，其核心类型必须是数组、指向数组的指针、切片、字符串、映射或允许接收操作的通道。与赋值一样，如果存在，左侧的操作数必须是可寻址的或映射索引表达式；它们表示迭代变量。如果范围表达式是通道，则最多允许一个迭代变量，否则最多允许两个。如果最后一个迭代变量是空白标识符，则范围子句相当于没有该标识符的同一子句。

范围表达式 x 在开始循环之前计算一次，但有一个例外：如果最多存在一个迭代变量并且 len(x) 为常量，则不计算范围表达式。

左侧的函数调用每次迭代都会计算一次。对于每次迭代，如果存在相应的迭代变量，则按如下方式生成迭代值：

```
Range expression                          1st value          2nd value

array or slice  a  [n]E, *[n]E, or []E    index    i  int    a[i]       E
string          s  string type            index    i  int    see below  rune
map             m  map[K]V                key      k  K      m[k]       V
channel         c  chan E, <-chan E       element  e  E
```

1. 对于数组、指向数组的指针或切片值 a，索引迭代值从元素索引 0 开始按递增顺序生成。如果最多存在一个迭代变量，则范围循环会生成从 0 到 len( a)-1 并且不索引到数组或切片本身。对于 nil 切片，迭代次数为 0。

2. 对于字符串值，“range”子句从字节索引 0 开始迭代字符串中的 Unicode 代码点。在连续迭代中，索引值将是连续 UTF-8 编码代码点的第一个字节的索引字符串和 rune 类型的第二个值将是相应代码点的值。如果迭代遇到无效的 UTF-8 序列，则第二个值将为 0xFFFD（Unicode 替换字符），并且下一次迭代将在字符串中前进一个字节。

3. 映射的迭代顺序未指定，并且不保证从一次迭代到下一次迭代的顺序相同。如果在迭代过程中删除了尚未到达的映射条目，则不会产生相应的迭代值。如果在迭代期间创建映射条目，则该条目可以在迭代期间产生或者可以被跳过。对于创建的每个条目以及从一次迭代到下一次迭代，选择可能会有所不同。如果映射为零，则迭代次数为 0。

4. 对于通道，生成的迭代值是在通道关闭之前在通道上发送的连续值。如果通道为零，则范围表达式将永远阻塞。

迭代值被分配给各自的迭代变量，如赋值语句中一样。

迭代变量可以通过“range”子句使用短变量声明（:=）的形式来声明。在这种情况下，它们的类型设置为各自迭代值的类型，并且它们的范围是“for”语句的块；它们在每次迭代中都会被重复使用。如果迭代变量在“for”语句之外声明，则执行后它们的值将是最后一次迭代的值。

```go
var testdata *struct {
	a *[7]int
}
for i, _ := range testdata.a {
	// testdata.a is never evaluated; len(testdata.a) is constant
	// i ranges from 0 to 6
	f(i)
}

var a [10]string
for i, s := range a {
	// type of i is int
	// type of s is string
	// s == a[i]
	g(i, s)
}

var key string
var val interface{}  // element type of m is assignable to val
m := map[string]int{"mon":0, "tue":1, "wed":2, "thu":3, "fri":4, "sat":5, "sun":6}
for key, val = range m {
	h(key, val)
}
// key == last map key encountered in iteration
// val == map[key]

var ch chan Work = producer()
for w := range ch {
	doWork(w)
}

// empty a channel
for range ch {}
```

## Go 语句 Go statements

“go”语句在同一地址空间内作为独立的并发控制线程或 goroutine 启动函数调用的执行。

```ebnf
GoStmt = "go" Expression .
```

表达式必须是函数或方法调用；它不能被括号括起来。对于表达式语句，内置函数的调用受到限制。

函数值和参数在调用 goroutine 中照常计算，但与常规调用不同，程序执行不会等待调用的函数完成。相反，该函数开始在新的 goroutine 中独立执行。当函数终止时，它的 goroutine 也会终止。如果函数有任何返回值，它们将在函数完成时被丢弃。

```go
go Server()
go func(ch chan<- bool) { for { sleep(10); ch <- true }} (c)
```

## Select 语句 Select statements

“选择”语句选择将进行一组可能的发送或接收操作中的哪一个。它看起来类似于“switch”语句，但所有情况都指的是通信操作。

```go
SelectStmt = "select" "{" { CommClause } "}" .
CommClause = CommCase ":" StatementList .
CommCase   = "case" ( SendStmt | RecvStmt ) | "default" .
RecvStmt   = [ ExpressionList "=" | IdentifierList ":=" ] RecvExpr .
RecvExpr   = Expression .
```

具有 RecvStmt 的情况可以将 RecvExpr 的结果分配给一个或两个变量，这些变量可以使用短变量声明来声明。 RecvExpr 必须是一个（可能带括号的）接收操作。最多可以有一个默认案例，它可能出现在案例列表中的任何位置。

“select”语句的执行分几个步骤进行：

1. 对于语句中的所有情况，接收操作的通道操作数以及发送语句的通道和右侧表达式在输入“select”语句时按源顺序仅计算一次。结果是一组要接收或发送到的通道，以及要发送的相应值。无论选择哪个（如果有）通信操作来进行，该评估中的任何副作用都会发生。 RecvStmt 左侧带有短变量声明或赋值的表达式尚未计算。

2. 如果一个或多个通信可以继续进行，则通过统一的伪随机选择选择一个可以继续进行的通信。否则，如果存在默认情况，则选择该情况。如果不存在默认情况，则“select”语句将阻塞，直到至少其中一个通信可以继续进行。

3. 除非所选情况为默认情况，否则将执行相应的通信操作。

4. 如果所选情况是带有短变量声明或赋值的 RecvStmt，则会对左侧表达式进行求值，并对接收的值（或多个值）进行赋值。

5. 执行所选case的语句列表。

由于 nil 通道上的通信永远无法继续，因此只有 nil 通道且没有默认情况的 select 会永远阻塞。

```go
var a []int
var c, c1, c2, c3, c4 chan int
var i1, i2 int
select {
case i1 = <-c1:
	print("received ", i1, " from c1\n")
case c2 <- i2:
	print("sent ", i2, " to c2\n")
case i3, ok := (<-c3):  // same as: i3, ok := <-c3
	if ok {
		print("received ", i3, " from c3\n")
	} else {
		print("c3 is closed\n")
	}
case a[f()] = <-c4:
	// same as:
	// case t := <-c4
	//	a[f()] = t
default:
	print("no communication\n")
}

for {  // send random sequence of bits to c
	select {
	case c <- 0:  // note: no statement, no fallthrough, no folding of cases
	case c <- 1:
	}
}

select {}  // block forever
```

## Return 语句 Return statements

函数F中的“return”语句终止F的执行，并且可选地提供一个或多个结果值。 F 推迟的任何函数都会在 F 返回其调用者之前执行。

```ebnf
ReturnStmt = "return" [ ExpressionList ] .
```

在没有结果类型的函数中，“return”语句不得指定任何结果值。

```go
func noResult() {
	return
}
```

有三种方法可以从具有结果类型的函数返回值：

1. 返回值可以在“return”语句中明确列出。每个表达式必须是单值并且可分配给函数结果类型的相应元素。

```go
func simpleF() int {
	return 2
}

func complexF1() (re float64, im float64) {
	return -7.0, -4.0
}
```

2. “return”语句中的表达式列表可以是对多值函数的单个调用。效果就好像从该函数返回的每个值都被分配给具有相应值类型的临时变量，后跟列出这些变量的“return”语句，此时应用前一种情况的规则。

```go
func complexF2() (re float64, im float64) {
	return complexF1()
}
```

3. 如果函数的结果类型指定了其结果参数的名称，则表达式列表可能为空。结果参数充当普通的局部变量，函数可以根据需要为其赋值。 “return”语句返回这些变量的值。

```go
func complexF3() (re float64, im float64) {
	re = 7.0
	im = 4.0
	return
}

func (devnull) Write(p []byte) (n int, _ error) {
	n = len(p)
	return
}
```

不管它们是如何声明的，所有结果值在进入函数时都会被初始化为其类型的零值。指定结果的“return”语句在执行任何延迟函数之前设置结果参数。

实现限制：如果返回位置的范围内存在与结果参数同名的不同实体（常量、类型或变量），则编译器可能不允许“返回”语句中使用空表达式列表。

```go
func f(n int) (res int, err error) {
	if _, err := f(n-1); err != nil {
		return  // invalid return statement: err is shadowed
	}
	return
}
```

## Break 语句 Break statements

“break”语句终止同一函数内最里面的“for”、“switch”或“select”语句的执行。

```ebnf
BreakStmt = "break" [ Label ] .
```

如果有标签，它一定是封闭的“for”、“switch”或“select”语句的标签，并且是执行终止的标签。

```go
OuterLoop:
	for i = 0; i < n; i++ {
		for j = 0; j < m; j++ {
			switch a[i][j] {
			case nil:
				state = Error
				break OuterLoop
			case item:
				state = Found
				break OuterLoop
			}
		}
	}
```

## Continue语句 Continue statements

“继续”语句通过将控制推进到循环块的末尾来开始最内层封闭“for”循环的下一次迭代。 “for”循环必须位于同一函数内。

```go
ContinueStmt = "continue" [ Label ] .
```

如果有标签，那么它一定是外层 "for "语句的标签，也就是执行前进的标签。

```go
RowLoop:
	for y, row := range rows {
		for x, data := range row {
			if data == endOfRow {
				continue RowLoop
			}
			row[x] = data + bias(x, y)
		}
	}
```

## Goto语句 Goto statements

“goto”语句将控制转移到同一函数内具有相应标签的语句。

```ebnf
GotoStmt = "goto" Label .
```

```go
goto Error
```

执行“goto”语句不得导致任何变量进入不在 goto 点范围内的范围。例如这个例子：

```go
	goto L  // BAD
	v := 3
L:
```

是错误的，因为跳转到标签 L 时跳过了 v 的创建。

块外部的“goto”语句无法跳转到该块内部的标签。例如这个例子：

```go
if n%2 == 1 {
	goto L1
}
for n > 0 {
	f()
	n--
L1:
	f()
	n--
}

```

是错误的，因为标签 L1 位于“for”语句的块内，但 goto 不在。

## Fallthrough 语句 Fallthrough statements

“fallthrough”语句将控制转移到表达式“switch”语句中下一个 case 子句的第一个语句。它只能用作此类子句中的最后一个非空语句。

```ebnf
FallthroughStmt = "fallthrough" .
```

## Defer 语句 Defer statements

“defer”语句调用一个函数，该函数的执行被推迟到周围函数返回的那一刻，要么是因为周围函数执行了 return 语句，到达了其函数体的末尾，要么是因为相应的 goroutine 正在恐慌。

```ebnf
DeferStmt = "defer" Expression .
```

表达式必须是函数或方法调用；它不能被括号括起来。对于表达式语句，内置函数的调用受到限制。

每次执行“defer”语句时，调用的函数值和参数都会照常评估并重新保存，但不会调用实际函数。相反，延迟函数会在周围函数返回之前立即调用，调用顺序与延迟函数相反。也就是说，如果周围函数通过显式 return 语句返回，则延迟函数将在该 return 语句设置任何结果参数之后但在函数返回到其调用者之前执行。如果延迟函数值求值为 nil，则在调用该函数时（而不是执行“defer”语句时）会发生执行混乱。

例如，如果延迟函数是函数文字，并且周围函数具有在文字范围内的命名结果参数([named result parameters](https://go.dev/ref/spec#Function_types))，则延迟函数可以在返回结果参数之前访问和修改结果参数。如果延迟函数有任何返回值，它们将在函数完成时被丢弃。 （另请参阅有关处理恐慌的部分。）

```go
lock(l)
defer unlock(l)  // unlocking happens before surrounding function returns

// prints 3 2 1 0 before surrounding function returns
for i := 0; i <= 3; i++ {
	defer fmt.Print(i)
}

// f returns 42
func f() (result int) {
	defer func() {
		// result is accessed after it was set to 6 by the return statement
		result *= 7
	}()
	return 6
}
```
