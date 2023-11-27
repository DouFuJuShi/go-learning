# Delve 使用介绍

---

## 安装

### 方案一：克隆git仓库进行编译

```sh-session
$ git clone https://github.com/go-delve/delve
$ cd delve
$ go install github.com/go-delve/delve/cmd/dlv
```

### 方案二：go版本 >= 1.16

```sh-session
# Install the latest release:
$ go install github.com/go-delve/delve/cmd/dlv@latest

# Install at tree head:
$ go install github.com/go-delve/delve/cmd/dlv@master

# Install at a specific version or pseudo-version:
$ go install github.com/go-delve/delve/cmd/dlv@v1.7.3
$ go install github.com/go-delve/delve/cmd/dlv@v1.7.4-0.20211208103735-2f13672765fe
```

如果在安装步骤中收到类似以下的错误：
```shell
found packages native (proc.go) and your_operating_system_and_architecture_combination_is_not_supported_by_delve (support_sentinel.go) in /home/pi/go/src/github.com/go-delve/delve/pkg/proc/native
```

这意味着你的操作系统和CPU架构的组合不被支持，检查go版本的输出。   

