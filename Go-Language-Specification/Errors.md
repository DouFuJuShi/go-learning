# 错误 Errors

预先声明的类型错误定义为

```go
type error interface {
    Error() string
}
```

这是表示错误条件的常规接口，nil 值表示无错误。例如，可以定义一个从文件中读取数据的函数：

```go
func Read(f *File, b []byte) (n int, err error)
```



## 扩展

### 标准库创建满足error错误值的两种方法：

1. #### 第一种

```go
import (
    "errors"
)

var  err = errors.New("First Error")
```

2. #### 第二种

```go
import (
    "fmt"
)
var errWithCtx = fmt.Errorf("index %d is out of bounds", i)
```

以上两种方法都返回error接口的实例

```go

// $GOROOT/src/errors/errors.go

type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}



// 大多数情况下，使用这两种方法构建的错误值就可以满足我们的需求了。但我们也要看到，
// 虽然这两种构建错误值的方法很方便，但它们给错误处理者提供的错误上下文（Error Context）
// 只限于以字符串形式呈现的信息，也就是 Error 方法返回的信息。

// 错误处理者需要从错误值中提取出更多信息，帮助他选择错误处理路径，显然这两种方法就不能满足了。
// 这个时候，我们可以自定义错误类型来满足这一需求。
// 比如：标准库中的 net 包就定义了一种携带额外错误上下文的错误类型

// $GOROOT/src/net/net.go
type OpError struct {
    Op string
    Net string
    Source Addr
    Addr Addr
    Err error
}

// 这样，错误处理者就可以根据这个类型的错误值提供的额外上下文信息，比如 Op、Net、Source 等，做出错误处理路径的选择，比如下面标准库中的代码：

// $GOROOT/src/net/http/server.go
func isCommonNetReadError(err error) bool {
    if err == io.EOF {
        return true
    }
    if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
        return true
    }
    if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
        return true
    }
    return false
}
```

### Go语言中的几种惯用错误处理策略

#### 策略一：透明错误处理策略

Go 语言中错误处理，是根据函数 / 方法返回的 error 类型变量中携带的错误值信息做决策，并选择后续执行路径的过程。  
这样，最简单的错误策略莫过于完全不关心返回错误值携带的具体上下文信息，只要发生错误就进入唯一的错误处理执行路径，比如下面这段代码：

```go
package main
fun main() {
    err := function()
    if err != nil {
        // 不关心err变量底层错误值所携带的具体上下文信息
        // 执行简单错误处理逻辑并返回
        ... ...
        return err
    }
}

func function(...) error {
    return errors.New("some error occurred")
}
```

#### 策略二：“哨兵”处理策略

当错误处理方不能只根据“透明的错误值”就做出错误处理路径选取的情况下，错误处理方会尝试对返回的错误值进行检视，于是就有可能出现下面代码中的反模式：

```go
data, err := b.Peek(1)
if err != nil {
    switch err.Error() {
    case "bufio: negative count":
        // ... ...
        return
    case "bufio: buffer full":
        // ... ...
        return
    case "bufio: invalid use of UnreadByte":
        // ... ...
        return
    default:
        // ... ...
        return
    }
}
```

简单来说，反模式就是，错误处理方以透明错误值所能提供的唯一上下文信息（描述错误的字符串），作为错误处理路径选择的依据。  
但这种“反模式”会造成严重的隐式耦合。这也就意味着，错误值构造方不经意间的一次错误描述字符串的改动，都会造成错误处理方处理行为的变化，  
并且这种通过字符串比较的方式，对错误值进行检视的性能也很差。

那这有什么办法吗？Go 标准库采用了定义导出的（Exported）“哨兵”错误值的方式，来辅助错误处理方检视（inspect）错误值并做出错误处理分支的决策，  
比如下面的 bufio 包中定义的“哨兵错误”：

```go

// $GOROOT/src/bufio/bufio.go
var (
    ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
    ErrInvalidUnreadRune = errors.New("bufio: invalid use of UnreadRune")
    ErrBufferFull        = errors.New("bufio: buffer full")
    ErrNegativeCount     = errors.New("bufio: negative count")
)
// 一般“哨兵”错误值变量以 ErrXXX 格式命名。
// 和透明错误策略相比，“哨兵”策略让错误处理方在有检视错误值的需求时候，可以“有的放矢”。

// 暴露“哨兵”错误值也意味着这些错误值和包的公共函数/方法一起成为了 API 的一部分。
// 一旦发布出去，开发者就要对它进行很好的维护。
// 而“哨兵”错误值也让使用这些值的错误处理方对它产生了依赖。
```

```go
// 利用了上面的哨兵错误，进行错误处理分支的决策
data, err := b.Peek(1)
if err != nil {
    switch err {
    case bufio.ErrNegativeCount:
        // ... ...
        return
    case bufio.ErrBufferFull:
        // ... ...
        return
    case bufio.ErrInvalidUnreadByte:
        // ... ...
        return
    default:
        // ... ...
        return
    }
}

// 从 Go 1.13 版本开始，标准库 errors 包提供了 Is 函数用于错误处理方对错误值的检视.
// Is函数类似于把一个 error 类型变量与“哨兵”错误值进行比较，比如下面代码：
// 类似 if err == ErrOutOfBounds{ … }
if errors.Is(err, ErrOutOfBounds) {
	// 越界的错误处理
}

// 不同的是，如果 error 类型变量的底层错误值是一个包装错误（Wrapped Error）， 
// errors.Is 方法会沿着该包装错误所在错误链（Error Chain)，
// 与链上所有被包装的错误（Wrapped Error）进行比较，直至找到一个匹配的错误为止。
// 下面是 Is 函数应用的一个例子：

var ErrSentinel = errors.New("the underlying sentinel error")

func main() {
   err1 := fmt.Errorf("wrap sentinel: %w", ErrSentinel)
   err2 := fmt.Errorf("wrap err1: %w", err1)
   println(err2 == ErrSentinel) // false
   if errors.Is(err2, ErrSentinel) {
      println("err2 is ErrSentinel")
      return
    }

    println("err2 is not ErrSentinel")
}

// 在这个例子中，我们通过 fmt.Errorf 函数，并且使用 %w 创建包装错误变量 err1 和 err2，
// 其中 err1 实现了对 ErrSentinel 这个“哨兵错误值”的包装，而 err2 又对 err1 进行了包装，这样就形成了一条错误链。
// 位于错误链最上层的是 err2，位于最底层的是 ErrSentinel。之后，我们再分别通过值比较和 errors.Is 这两种方法，
// 判断 err2 与 ErrSentinel 的关系。运行上述代码，我们会看到如下结果：
// false
// err2 is ErrSentinel
```

#### 策略三：错误值类型检视策略

上面我们看到，基于 Go 标准库提供的错误值构造方法构造的“哨兵”错误值，除了让错误处理方可以“有的放矢”的进行值比较之外，  
并没有提供其他有效的错误上下文信息。那如果遇到错误处理方需要错误值提供更多的“错误上下文”的情况，上面这些错误处理策略和错误值构造方式都无法满足。

这种情况下，我们需要通过自定义错误类型的构造错误值的方式，来提供更多的“错误上下文”信息。并且，由于错误值都通过 error 接口变量统一呈现，  
要得到底层错误类型携带的错误上下文信息，错误处理方需要使用 Go 提供的类型断言机制（Type Assertion）或类型选择机制（Type Switch），  
这种错误处理方式，我称之为错误值类型检视策略。

```go
// 这个 json 包中自定义了一个UnmarshalTypeError的错误类型
// $GOROOT/src/encoding/json/decode.go
type UnmarshalTypeError struct {
    Value  string       
    Type   reflect.Type 
    Offset int64        
    Struct string       
    Field  string      
}

// 错误处理方可以通过错误类型检视策略，获得更多错误值的错误上下文信息，下面就是利用这一策略的 json 包的一个方法的实现
// 这段代码通过类型 switch 语句得到了 err 变量代表的动态类型和值，然后在匹配的 case 分支中利用错误上下文信息进行处理
// $GOROOT/src/encoding/json/decode.go
func (d *decodeState) addErrorContext(err error) error {
   if d.errorContext.Struct != nil || len(d.errorContext.FieldStack) > 0 {
      switch err := err.(type) {
         case *UnmarshalTypeError:
            err.Struct = d.errorContext.Struct.Name()
            err.Field = strings.Join(d.errorContext.FieldStack, ".")
            return err
      }
   }
   return err
}

// 从 Go 1.13 版本开始，标准库 errors 包提供了As函数给错误处理方检视错误值。
// As函数类似于通过类型断言判断一个 error 类型变量是否为特定的自定义错误类型，如下面代码所示：

// 类似 if e, ok := err.(*MyError); ok { … }
var e *MyError
if errors.As(err, &e) {
    // 如果err类型为*MyError，变量e将被设置为对应的错误值
}

// 不同的是，如果 error 类型变量的动态错误值是一个包装错误，errors.As函数会沿着该包装错误所在错误链，
// 与链上所有被包装的错误的类型进行比较，直至找到一个匹配的错误类型，就像 errors.Is 函数那样。下面是As函数应用的一个例子

type MyError struct {
    e string
}

func (e *MyError) Error() string {
    return e.e
}

func main() {
   var err = &MyError{"MyError error demo"}
   err1 := fmt.Errorf("wrap err: %w", err)
   err2 := fmt.Errorf("wrap err1: %w", err1)
   var e *MyError
   if errors.As(err2, &e) {
      println("MyError is on the chain of err2")
      println(e == err)
      return
   }
    println("MyError is not on the chain of err2")
} 

// output:
// MyError is on the chain of err2
// true
```

#### 策略四：错误行为特征检视策略

将某个包中的错误类型归类，统一提取出一些公共的错误行为特征，并将这些错误行为特征放入一个公开的接口类型中。  
以标准库中的net包为例，它将包内的所有错误类型的公共行为特征抽象并放入net.Error这个接口中，如下面代码：

```go
// $GOROOT/src/net/net.go
type Error interface {
    error
    Timeout() bool  
    Temporary() bool
}

// net.Error 接口包含两个用于判断错误行为特征的方法：
// Timeout 用来判断是否是超时（Timeout）错误，Temporary 用于判断是否是临时（Temporary）错误

// 而错误处理方只需要依赖这个公共接口，就可以检视具体错误值的错误行为特征信息，并根据这些信息做出后续错误处理分支选择的决策。
// 这里，我们再看一个 http 包使用错误行为特征检视策略进行错误处理的例子，加深下理解：

// $GOROOT/src/net/http/server.go
func (srv *Server) Serve(l net.Listener) error {
   ... ...
   for {
      rw, e := l.Accept()
      if e != nil {
         select {
            case <-srv.getDoneChan():
                return ErrServerClosed
            default:
         }
         if ne, ok := e.(net.Error); ok && ne.Temporary() {
            // 注：这里对临时性(temporary)错误进行处理
            ... ...
            time.Sleep(tempDelay)
            continue
         }
         return e
      }
      ...
   }
   ... ...
}

// Accept 方法实际上返回的错误类型为*OpError，它是 net 包中的一个自定义错误类型，它实现了错误公共特征接口net.Error，如下代码所示
// 
// $GOROOT/src/net/net.go
type OpError struct {
   ... ...
   // Err is the error that occurred during the operation.
   Err error
}

type temporary interface {
    Temporary() bool
}

func (e *OpError) Temporary() bool {
   if ne, ok := e.Err.(*os.SyscallError); ok {
      t, ok := ne.Err.(temporary)
      return ok && t.Temporary()
   }
   t, ok := e.Err.(temporary)
   return ok && t.Temporary()
}
// 因此，OpError 实例可以被错误处理方通过net.Error接口的方法，判断它的行为是否满足 Temporary 或 Timeout 特征。
```
