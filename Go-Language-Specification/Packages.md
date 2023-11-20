# 包 Packages

Go 程序是通过将软件包连接起来而构建的。一个包又是由一个或多个源文件构建而成，这些源文件共同声明了属于该包的常量、类型、变量和函数，这些常量、类型、变量和函数可在同一包的所有文件中访问。这些元素可以导出并在另一个包中使用。

## 源码文件组织 Source file organization

每个源文件都由一个包子句（package clause）组成，该子句定义了源文件所属的包；接着是一组可能为空的导入声明（import declarations），这些声明声明了源文件希望使用的包的内容；然后是一组可能为空的函数、类型、变量和常量声明。

```ebnf
SourceFile       = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .
```

## 包条款 Package clause

package 子句开始每个源文件并定义该文件所属的包。

```go
PackageClause  = "package" PackageName .
PackageName    = identifier .
```

包名不能是空白标识符。

```go
package math
```

一组共享相同 PackageName 的文件构成一个软件包的实现。实现可能要求软件包的所有源文件都位于同一目录下。

## 导入声明 Import declarations

导入声明指出，包含声明的源文件取决于导入包的功能（§程序初始化和执行），并允许访问该包的导出标识符。导入命名一个用于访问的标识符 (PackageName) 和一个指定要导入的包的 ImportPath。

```go
ImportDecl       = "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) .
ImportSpec       = [ "." | PackageName ] ImportPath .
ImportPath       = string_lit .
```

PackageName 用于限定标识符，以访问导入源文件中软件包的导出标识符。它在文件块中声明。如果省略了 PackageName，则默认为导入包的包子句中指定的标识符。如果出现明确的句号（.）而不是名称，则在该软件包的软件包块中声明的所有软件包导出标识符都将在导入源文件的文件块中声明，并且必须在不使用限定符的情况下进行访问。

ImportPath 的解释取决于实现，但它通常是编译软件包完整文件名的子串，也可能是相对于已安装软件包的存储库。

实现限制：编译器可以将 ImportPaths 限制为仅使用属于 Unicode L、M、N、P 和 S 一般类别的字符（不带空格的图形字符）的非空字符串，并且还可以排除字符 !"#$%& '()*,:;<=>?[\]^`{|} 和 Unicode 替换字符 U+FFFD。

考虑一个包含 package 子句 package math 的已编译包，该子句导出函数 Sin，并将编译后的包安装在由“lib/math”标识的文件中。此表说明了如何在各种类型的导入声明之后导入包的文件中访问 Sin。

```go
Import declaration          Local name of Sin

import   "lib/math"         math.Sin
import m "lib/math"         m.Sin
import . "lib/math"         Sin
```

导入声明声明了导入包和被导入包之间的依赖关系。软件包直接或间接导入自身，或直接导入一个软件包而不引用它的任何导出标识符，都是非法的。如果只是为了包的副作用（初始化）而导入包，则应使用空白标识符作为显式包名：

```go
import _ "lib/math"
```

## 示例包 An example package

下面是一个实现并发素数筛的完整 Go 包。

```go
package main

import "fmt"

// Send the sequence 2, 3, 4, … to channel 'ch'.
func generate(ch chan<- int) {
    for i := 2; ; i++ {
        ch <- i  // Send 'i' to channel 'ch'.
    }
}

// Copy the values from channel 'src' to channel 'dst',
// removing those divisible by 'prime'.
func filter(src <-chan int, dst chan<- int, prime int) {
    for i := range src {  // Loop over values received from 'src'.
        if i%prime != 0 {
            dst <- i  // Send 'i' to channel 'dst'.
        }
    }
}

// The prime sieve: Daisy-chain filter processes together.
func sieve() {
    ch := make(chan int)  // Create a new channel.
    go generate(ch)       // Start generate() as a subprocess.
    for {
        prime := <-ch
        fmt.Print(prime, "\n")
        ch1 := make(chan int)
        go filter(ch, ch1, prime)
        ch = ch1
    }
}

func main() {
    sieve()
}
```
