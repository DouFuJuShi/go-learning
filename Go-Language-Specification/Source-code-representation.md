# 源代码表示 Source code representation

源代码是以 UTF-8 编码的 Unicode 文本。文本未进行规范化处理，因此单个重音符号码点有别于由重音符号和字母组合而成的同一字符；后者被视为两个码点。为简便起见，本文件将使用未限定的术语 "字符 "来指代源文本中的 Unicode 代码点。
每个码位都是不同的；例如，大写字母和小写字母就是不同的字符。
执行限制： 为了与其他工具兼容，编译器可能不允许在源文本中使用 NUL 字符 (U+0000)。
执行限制： 为了与其他工具兼容，如果UTF-8 编码的字节序号 (U+FEFF) 是源文本中的第一个 Unicode 代码点，则编译器可以忽略该字节序号。在源文本的其他任何地方，字节序号都是不允许的。

## 字符 Characters

以下术语用于表示特定的 Unicode 字符类别：

```ebnf
newline        = /* the Unicode code point U+000A */ .
unicode_char   = /* an arbitrary Unicode code point except newline */ .
unicode_letter = /* a Unicode code point categorized as "Letter" */ .
unicode_digit  = /* a Unicode code point categorized as "Number, decimal digit" */ .
```

在 Unicode 标准 8.0 ([The Unicode Standard 8.0](https://www.unicode.org/versions/Unicode8.0.0/))中，第 4.5 节“常规类别”定义了一组字符类别。 Go 将任何字母类别 Lu、Ll、Lt、Lm 或 Lo 中的所有字符视为 Unicode 字母，将数字类别 Nd 中的所有字符视为 Unicode 数字。

## 字母和数字 Letters and digits

下划线字符 _ (U+005F) 被视为小写字母。

```ebnf
letter        = unicode_letter | "_" .
decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
```


