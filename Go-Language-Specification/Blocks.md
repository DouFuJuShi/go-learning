# 代码块 Blocks

块是匹配大括号内的可能为空的声明和语句序列。

```ebnf
Block = "{" StatementList "}" .
StatementList = { Statement ";" } .
```

除了源代码中的显式代码块，还有隐式代码块：

1. Universe 宇宙块包含所有 Go 源文本。

2. 每个包都有一个包块，其中包含该包的所有 Go 源文本。

3. 每个文件都有一个文件块，其中包含该文件中的所有 Go 源文本。

4. 每个“if”、“for”和“switch”语句都被视为位于其自己的隐式块中。

5. switch "或 "select "语句中的每个子句都是一个隐式块。

块嵌套并影响范围界定。
