# 词法元素 Lexical elements

## 注释 Comments

注释作为程序文档。有两种形式：

1. 行注释以字符序列 // 开始，到行尾停止。

2. 一般注释以字符序列 /* 开始，以第一个后续字符序列 */ 结束。

注释不能从符码或字符串文字开始，也不能从注释开始。不含换行符的一般注释就像空格。任何其他注释都与换行符一样。

## 标记 Tokens

标记构成了 Go 语言的词汇。它分为四类：标识符、关键字、运算符和标点符号以及字面量。由空格（U+0020）、水平制表符（U+0009）、回车符（U+000D）和换行符（U+000A）组成的空白区域将被忽略，除非它能分隔本来会合并成一个标记的标记。此外，换行符或文件结束符也可能触发分号的插入。在将输入分割成标记时，下一个标记是构成有效标记的最长字符序列。

## 分号 Semicolons

正式语法使用分号"; "作为许多语句的结束符。Go 程序可以使用以下两条规则省略大部分分号：

1. 当输入被分割成标记符时，如果一行的最后一个标记符是以下的情况，则分号会被马上自动插入到标记流中：
   
   - [identifier](https://go.dev/ref/spec#Identifiers)
   
   - [integer](https://go.dev/ref/spec#Integer_literals), [floating-point](https://go.dev/ref/spec#Floating-point_literals), [imaginary](https://go.dev/ref/spec#Imaginary_literals), [rune](https://go.dev/ref/spec#Rune_literals), [string](https://go.dev/ref/spec#String_literals) 字面量
   
   - `break`, `continue`, `fallthrough`,  `return` 关键词([keywords](https://go.dev/ref/spec#Keywords))之一
   
   - `++`, `--`, `)`, `]`,  `}`运算符和标点符号( [operators and punctuation](https://go.dev/ref/spec#Operators_and_punctuation))之一 

2. 为了允许复杂的语句占据一行，可以在结束的“)”或“}”之前省略分号。

为反映习惯用法，本文档中的代码示例使用这些规则省略了分号。

## 标识符 Identifiers

标识符用于命名变量和类型等程序实体。标识符是由一个或多个字母和数字组成的序列。标识符的第一个字符必须是字母。

```ebnf
identifier = letter { letter | unicode_digit } .
```

```go
a
_x9
ThisVariableIsExported
αβ
```

一些标识符是预先声明([predeclared](https://go.dev/ref/spec#Predeclared_identifiers))的。

## 关键词 Keywords

以下关键字为保留字，不得用作标识符。

```go
break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
```

## 运算符和标点符号 Operators and punctuation

以下字符序列表示运算符[operators](https://go.dev/ref/spec#Operators)（包括赋值运算符[assignment operators](https://go.dev/ref/spec#Assignment_statements)）和标点符号：

```go
+    &     +=    &=     &&    ==    !=    (    )
-    |     -=    |=     ||    <     <=    [    ]
*    ^     *=    ^=     <-    >     >=    {    }
/    <<    /=    <<=    ++    =     :=    ,    ;
%    >>    %=    >>=    --    !     ...   .    :
     &^          &^=          ~
```

## 整型字面量 Integer Literals

整数字面量是代表整数常量([integer constant](https://go.dev/ref/spec#Constants))的数字序列。可选前缀设置非十进制基数： 二进制为 0b 或 0B，八进制为 0、0o 或 0O，十六进制为 0x 或 0X。单个 0 被视为十进制 0。在十六进制字面中，字母 a 至 f 和 A 至 F 代表数值 10 至 15。

为便于阅读，可在基数前缀后或连续数字之间使用下划线字符 _；此类下划线不会改变字面量的值。

```ebnf
int_lit        = decimal_lit | binary_lit | octal_lit | hex_lit .
decimal_lit    = "0" | ( "1" … "9" ) [ [ "_" ] decimal_digits ] .
binary_lit     = "0" ( "b" | "B" ) [ "_" ] binary_digits .
octal_lit      = "0" [ "o" | "O" ] [ "_" ] octal_digits .
hex_lit        = "0" ( "x" | "X" ) [ "_" ] hex_digits .

decimal_digits = decimal_digit { [ "_" ] decimal_digit } .
binary_digits  = binary_digit { [ "_" ] binary_digit } .
octal_digits   = octal_digit { [ "_" ] octal_digit } .
hex_digits     = hex_digit { [ "_" ] hex_digit } .
```

```go
42
4_2
0600
0_600
0o600
0O600       // second character is capital letter 'O'
0xBadFace
0xBad_Face
0x_67_7a_2f_cc_40_c6
170141183460469231731687303715884105727
170_141183_460469_231731_687303_715884_105727

_42         // an identifier, not an integer literal
42_         // invalid: _ must separate successive digits
4__2        // invalid: only one _ at a time
0_xBadFace  // invalid: _ must separate successive digits
```

## 浮点数字面量 Floating-point literals

浮点字面量是浮点常量([floating-point constant](https://go.dev/ref/spec#Constants).)的十进制或十六进制表示形式。

十进制浮点文字由整数部分（小数位）、小数点、小数部分（小数位）和指数部分（e 或 E 后跟可选符号和小数位）组成。整数部分或小数部分之一可以省略；小数点或指数部分之一可以被省略。指数值 exp 将尾数（整数和小数部分）缩放 10 ~exp~。


