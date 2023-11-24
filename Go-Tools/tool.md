# [Run specified go tool](https://pkg.go.dev/cmd/go#hdr-Run_specified_go_tool)

```shell
go tool [-n] command [args...]
```

Tool 运行由参数标识的 go tool 命令。不带任何参数，它会打印已知工具的列表。
<mark>-n</mark> 标志使工具打印将要执行的命令，但不执行它。
有关每个工具命令的更多信息，请参阅“go doc cmd/\<command\>”。

## 查看所有command

```shell
go tool
// or 
go tool -n
```

## 查看某个command的更多信息

```shell
go doc cmd/<command> 
```
