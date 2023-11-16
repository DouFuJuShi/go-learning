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

十进制浮点数字面量由整数部分（小数位）、小数点、小数部分（小数位）和指数部分（e 或 E 后跟可选符号和小数位）组成。整数部分或小数部分之一可以省略；小数点或指数部分之一可以被省略。指数值 exp 将尾数（整数和小数部分）缩放 10<sup>exp</sup> 。

十六进制浮点文字由 0x 或 0X 前缀、整数部分（十六进制数字）、小数点、小数部分（十六进制数字）和指数部分（p 或 P 后跟可选符号和十进制数字）组成）。整数部分或小数部分可以省略；小数点也可以省略，但指数部分是必需的。 （此语法与 IEEE 754-2008 §5.12.3 中给出的语法匹配。）指数值 exp 将尾数（整数和小数部分）缩放 2<sup>exp</sup>。

为了便于阅读，下划线字符 _ 可以出现在基本前缀之后或连续数字之间；这样的下划线不会改变字面值。

```ebnf
float_lit         = decimal_float_lit | hex_float_lit .

decimal_float_lit = decimal_digits "." [ decimal_digits ] [ decimal_exponent ] |
                    decimal_digits decimal_exponent |
                    "." decimal_digits [ decimal_exponent ] .
decimal_exponent  = ( "e" | "E" ) [ "+" | "-" ] decimal_digits .

hex_float_lit     = "0" ( "x" | "X" ) hex_mantissa hex_exponent .
hex_mantissa      = [ "_" ] hex_digits "." [ hex_digits ] |
                    [ "_" ] hex_digits |
                    "." hex_digits .
hex_exponent      = ( "p" | "P" ) [ "+" | "-" ] decimal_digits .
```

```go
0.
72.40
072.40       // == 72.40
2.71828
1.e+0
6.67428e-11
1E6
.25
.12345E+5
1_5.         // == 15.0
0.15e+0_2    // == 15.0

0x1p-2       // == 0.25
0x2.p10      // == 2048.0
0x1.Fp+0     // == 1.9375
0X.8p-0      // == 0.5
0X_1FFFP-16  // == 0.1249847412109375
0x15e-2      // == 0x15e - 2 (integer subtraction)

0x.p1        // invalid: mantissa has no digits
1p-2         // invalid: p exponent requires hexadecimal mantissa
0x1.5e-2     // invalid: hexadecimal mantissa requires p exponent
1_.5         // invalid: _ must separate successive digits
1._5         // invalid: _ must separate successive digits
1.5_e1       // invalid: _ must separate successive digits
1.5e_1       // invalid: _ must separate successive digits
1.5e1_       // invalid: _ must separate successive digits
```

## 虚数字面量 Imaginary Literals

虚数文字表示复数常量([complex constant](https://go.dev/ref/spec#Constants))的虚部。它由一个整数([integer](https://go.dev/ref/spec#Integer_literals))或浮点数([floating-point](https://go.dev/ref/spec#Floating-point_literals))后跟小写字母 i 组成。虚数文字的值是相应整数或浮点数的值乘以虚数单位 i。

```ebnf
imaginary_lit = (decimal_digits | int_lit | float_lit) "i" .
```

为了向后兼容，虚字面量的整数部分完全由十进制数字（可能还有下划线）组成，即使以 0 开头，也被视为十进制整数。

```go
0i
0123i         // == 123i for backward-compatibility
0o123i        // == 0o123 * 1i == 83i
0xabci        // == 0xabc * 1i == 2748i
0.i
2.71828i
1.e+0i
6.67428e-11i
1E6i
.25i
.12345E+5i
0x1p-2i       // == 0x1p-2 * 1i == 0.25i
```

## Rune字面量 Rune Literals

rune 字面量表示 rune 常量([rune constant](https://go.dev/ref/spec#Constants))，即标识 Unicode 代码点的整数值。符文文字表示为用单引号括起来的一个或多个字符，如“x”或“\n”。引号内可以出现除换行符和未转义单引号之外的任何字符。单引号字符表示字符本身的 Unicode 值，而以反斜杠开头的多字符序列则以各种格式编码值。
最简单的形式表示引号内的单个字符；由于 Go 源文本是以 UTF-8 编码的 Unicode 字符，因此多个 UTF-8 编码的字节可以表示单个整数值。例如，文字“a”保存一个字节，表示文字 a，Unicode U+0061，值 0x61，而“ä”保存两个字节 (0xc3 0xa4)，表示文字 a-分音符，U+00E4，值 0xe4。

通过几种反斜杠转义，可以将任意数值编码为 ASCII 文本。有四种方法可以将整数值表示为数字常量： \x 后跟两个十六进制数字；\u 后跟四个十六进制数字；\U 后跟八个十六进制数字，以及一个纯反斜杠 \ 后跟三个八进制数字。在每种情况下，字面量的值都是相应基数的数字所代表的值。

尽管这些表示都产生一个整数，但它们具有不同的有效范围。八进制转义符必须表示 0 到 255 之间的值（含 0 和 255）。十六进制转义符通过构造满足这个条件。转义符 \u 和 \U 代表 Unicode 代码点，因此其中的某些值是非法的，特别是那些高于 0x10FFFF 和代理一半的值。

在反斜杠之后，某些单字符转义符代表特殊值：

```ebnf
\a   U+0007 alert or bell
\b   U+0008 backspace
\f   U+000C form feed
\n   U+000A line feed or newline
\r   U+000D carriage return
\t   U+0009 horizontal tab
\v   U+000B vertical tab
\\   U+005C backslash
\'   U+0027 single quote  (valid escape only within rune literals)
\"   U+0022 double quote  (valid escape only within string literals)
```

Rune字面中反斜线后的未识别字符是非法字符。

```ebnf
rune_lit         = "'" ( unicode_value | byte_value ) "'" .
unicode_value    = unicode_char | little_u_value | big_u_value | escaped_char .
byte_value       = octal_byte_value | hex_byte_value .
octal_byte_value = `\` octal_digit octal_digit octal_digit .
hex_byte_value   = `\` "x" hex_digit hex_digit .
little_u_value   = `\` "u" hex_digit hex_digit hex_digit hex_digit .
big_u_value      = `\` "U" hex_digit hex_digit hex_digit hex_digit
                           hex_digit hex_digit hex_digit hex_digit .
escaped_char     = `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) .
```

```go
'a'
'ä'
'本'
'\t'
'\000'
'\007'
'\377'
'\x07'
'\xff'
'\u12e4'
'\U00101234'
'\''         // rune literal containing single quote character
'aa'         // illegal: too many characters
'\k'         // illegal: k is not recognized after a backslash
'\xa'        // illegal: too few hexadecimal digits
'\0'         // illegal: too few octal digits
'\400'       // illegal: octal value over 255
'\uDFFF'     // illegal: surrogate half
'\U00110000' // illegal: invalid Unicode code point
```

## String字面量 String Literals

字符串字面量表示通过连接字符获得的字符串常量。有两种形式：原始字符串字面量和解释字符串字面量。

原始字符串字面量是反引号之间的字符序列，如 `foo`。在引号内，除了反引号外，可以出现任何字符。原始字符串字面量的值是由引号之间未解释（隐含 UTF-8 编码）的字符组成的字符串；特别是，反斜线没有特殊含义，字符串可能包含换行符。原始字符串字面内的回车符（'\r'）会从原始字符串值中舍弃。

解释字符串字面量是双引号之间的字符序列，如 "bar"。在引号内，除了换行符和未转义的双引号外，可以出现任何字符。引号之间的文本构成字面量的值，反斜线转义的解释与rune字面量中的解释相同（除了 \' 是非法的，\"是合法的），并有相同的限制。三位八进制（\nnn）和两位十六进制（\xnn）转义代表结果字符串的单个字节；所有其他转义代表单个字符的（可能是多字节的）UTF-8 编码。因此，在字符串字面 \377 和 \xFF 表示值为 0xFF=255 的单字节，而 ÿ、\u00FF、\U000000FF 和 \xc3\xbf 表示字符 U+00FF 的 UTF-8 编码的两个字节 0xc3 0xbf。

```ebnf
string_lit             = raw_string_lit | interpreted_string_lit .
raw_string_lit         = "`" { unicode_char | newline } "`" .
interpreted_string_lit = `"` { unicode_value | byte_value } `"` .
```

```go
`abc`                // same as "abc"
`\n
\n`                  // same as "\\n\n\\n"
"\n"
"\""                 // same as `"`
"Hello, world!\n"
"日本語"
"\u65e5本\U00008a9e"
"\xff\u00FF"
"\uD800"             // illegal: surrogate half
"\U00110000"         // illegal: invalid Unicode code point
```

这些示例都代表相同的字符串：

```go
"日本語"                                 // UTF-8 input text
`日本語`                                 // UTF-8 input text as a raw literal
"\u65e5\u672c\u8a9e"                    // the explicit Unicode code points
"\U000065e5\U0000672c\U00008a9e"        // the explicit Unicode code points
"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"  // the explicit UTF-8 bytes
```

如果源代码将一个字符表示为两个代码点，例如涉及重音符号和字母的组合形式，则如果将其放入 rune 文字（它不是单个代码点）中，结果将是错误，并且将显示为如果放置在字符串文字中，则有两个代码点。
