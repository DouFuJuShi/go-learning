# 符号表示法 Notation

---

通过使用扩展巴科斯-瑙尔形式 (EBNF) 来表达：

```ebnf
Production  = production_name "=" [ Expression ] "." .
Expression  = Alternative { "|" Alternative } .
Alternative = Term { Term } .
Term        = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .
```

EBNF的基本语法形式如下， 这个形式也被叫做production：

> 左式(LeftHandSid) = 右边(RightHandSide)

左式也被叫做 非终端符号 non-terminal symbol，而右式则描述了左式的组成。

production是由术语及以下运算符所构成的表达式，优先级越来越高：

```ebnf
|   alternation
()  grouping 分组
[]  option (0 or 1 times)  可选项 0 或者 1 次
{}  repetition (0 to n times) 重复0...n次
```

用小写的production names来命名区分【词法标记】。非终端符号用CamelCase形式书写。词法标记【终端符号串】通过双引号、反引号包裹；

a … b 形式表示从 a 到 b 的字符集作为替代。水平省略号...也在规范的其他地方使用，非正式地表示未进一步指定的各种枚举或代码片段。字符 ...（与三个字符 ... 相对）不是 Go 语言的标记。



说明：

| 记号       | 定义    | 说明                                                                                   |
|:--------:|:-----:|:------------------------------------------------------------------------------------:|
| =        | 定义    |                                                                                      |
| ,        | 链接符   |                                                                                      |
| ;        | 结束符号  |                                                                                      |
| .        | 结束符号  |                                                                                      |
| |        | 或     |                                                                                      |
| [...]    | 可选    | 可选	0 或者 1 次                                                                          |
| {...}    | 重复    | 重复	重复0...n次                                                                          |
| (...)    | 分组    |                                                                                      |
| "..."    | 终端字符串 | 形成所描述的语言的最基本符号。所描述语言的标点符号(不是EBNF自己的)会被左右加引号(它们也是终端符号)，而其他终端符号会用粗体(这边因不方便加粗，就不加粗了)打印。 |
| '....'   | 终端字符串 | 终端字符串	go这里用了反引号``                                                                    |
| (*...\*) | 注释    |                                                                                      |
| ?...?    | 特殊序列  |                                                                                      |
| -        | 除外    |                                                                                      |

参考：

https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form#Advantages_over_BNF
