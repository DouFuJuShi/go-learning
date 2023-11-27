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

### Homebrew

```shell
brew install delve
```

### macOS 注意事项

在 macOS 上，确保同时安装了命令行开发工具：

```shell
$ xcode-select --install
```

如果您没有使用 Xcode 启用开发者模式，则每次使用调试器时都会被要求授权。要启用 "开发者模式 "并只需在每个会话中授权一次，请使用

```shell
sudo /usr/sbin/DevToolsSecurity -enable
```

您可能还需要将您的用户添加到开发人员组：

```shell
sudo dscl . append /Groups/_developer GroupMembership $(whoami)
```

### 编译 macOS 原生后端

您不需要 macOS 本机后端，它存在已知问题。如果您仍然想构建它：

1. 运行 xcode-select --install

2. 在 macOS 10.14 上，通过运行 /Library/Developer/CommandLineTools/Packages/macOS_SDK_headers_for_macOS_10.14.pkg 手动安装旧版包含标头

3. 将 repo 克隆到 $GOPATH/src/github.com/go-delve/delve 中

4. 在该目录中运行 make install（在某些版本的 macOS 上，首次运行时需要 root 权限才能安装新证书）

makefile 将自动创建和安装自签名证书。

## 使用调试器附加到正在运行的 Go 进程

在 GoLand 中，可以将调试器附加到本地计算机、远程计算机或 Docker 容器中正在运行的 Go 进程。

### Attach 本地进程

您可以调试从命令行启动的应用程序。在这种情况下，应用程序在 IDE 外部运行，但在同一台本地计算机上运行。要调试应用程序，您需要在 IDE 中打开项目并将调试器附加到正在运行的进程。

1. 安装 gops 包

```shell
go get -t github.com/google/gops/
```

点击 Run | Attach to Process (`⌥Opt``⇧Shift``F5`) 在通知窗口中，单击调用“go get gops”链接。![](images/go_invoke_go_get_gops.png)
