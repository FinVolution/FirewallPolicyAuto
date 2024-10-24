# 项目依赖

本项目由`Golang + Iris`构建，运行需要 `Golang` 版本 1.22+。

## 使用方式

### Development

> 配置文件 `config/config.toml`

```toml
listen_addr = ":8080"                         # 监听地址端口

[log]
path = "./logs/api.log"                       # 日志文件路径
level = "DEBUG"                               # 日志级别[DEBUG(default) INFO WARN ERROR FATAL]
max_size = 100                                # 最大日志文件大小(MB)
backups = 10                                  # 最大备份数
max_age = 7                                   # 最大保存天数

[[firewall]]
name = "飞塔防火墙"                             # 防火墙名称
brand = "fortinet"                            # 品牌 当前支持可选[h3c fortinet]
address = "10.4.xxx.xxx"                      # 防火墙地址
protocol = "https"                            # 请求协议 根据防火墙开启服务配置[http https]
version = "v1"                                # API版本 当前仅支持[v1]
token = "fortinet_rest_api_token"             # 认证token[当是飞塔防火墙时必须]
[[firewall.virtual_zone]]
name = "名称1"                                 # 虚拟墙域名称
code = "code1"                                # 虚拟墙域编码
[[firewall.virtual_zone]]
name = "名称2"
code = "code2"

[[firewall]]
name = "H3C防火墙"
brand = "h3c"
address = "172.20.xxx.xxx"
protocol = "https"
version = "v1"
username = "username"                          # 用户名[当是H3C防火墙时必须]
password = "password"                          # 密码[当是H3C防火墙时必须]
```

```bash
# 克隆项目
git clone http://git.ppdaicorp.com/fte/FirewallPolicyAuto.git
# 进入service文件夹
cd service
# 安装依赖包
go mod tidy
# 运行
go run .
```

### Build

```bash
# 构建
make dist
# 产物
bin/service_api
```

> 使用dockerfile构建镜像，在docker run时请注意配置一下config.toml配置文件的挂载路径 /app/firewall-policy-auto/config/config.toml
> 或者修改dockerfile的ENTRYPOINT启动方式，使用 -c 参数指定配置文件路径

## 其他依赖项
### 线上依赖
```json
  github.com/go-playground/validator/v10 v10.22.0  # 参数校验
  github.com/kataras/iris/v12 v12.2.11             # web框架
  github.com/spf13/viper v1.19.0                   # 配置管理
  github.com/stretchr/testify v1.9.0               # 测试
  go.uber.org/zap v1.27.0                          # 日志                       
  gopkg.in/natefinch/lumberjack.v2 v2.2.1          # 日志截断
```

# 目录结构

```
├─config        
│ ├─config.go           # 配置项管理及初始化方法
│ ├─config.toml         # 配置文件（如果需要）
│ └─default.go          # 默认的配置示例(也可以修改此文件，修改启动配置后使用默认配置)
├─pkg          
│ ├─firewall            # 防火墙操作包
│ │  ├─dto              # 数据model
│ │  ├─fortinet_v1      # 飞塔防火墙API封装v1版本
│ │  ├─h3c_v1           # H3C防火墙API封装v1版本
│ │  ├─requests         # API请求封装
│ │  └─factory.go       # 防火墙api统一入口
│ └─logger              # 日志包
├─router             
│ ├─v1                  # v1版本路由
│ ├─cors.go             # cors插件
│ └─router.go           # 路由
├─utils             
│ ├─common.go           # 公共方法
│ ├─response.go         # 标准化返回
│ └─validator.go        # 参数校验包
├─version             
│ └─version.go          # 版本号
├─Dockerfile            # dockerfile配置
├─go.mod                # gomod配置
├─go.sum                # gosum
├─main.go               # 入口文件
├─Makefile              # make配置
└─README.md             # readme
```