# gim
gim
### 简介
1. comet, 以直接部署多个节点, 每个节点保证serverID 唯一, 在配置文件comet.toml
2. logic(业务逻辑层), 无状态, 各层通过rpc通讯, 容易扩展, 支持http接口来接收消息
3. job(任务推送层)redis 订阅发布功能进行推送到comet层。

### 架构图
![image](https://github.com/Cluas/static/blob/master/%E6%9E%B6%E6%9E%84.png?raw=true)

### 时序图
以下Comet 层, Logic 层, Job层都可以灵活扩展机器
![image](https://github.com/Cluas/static/blob/master/%E6%97%B6%E5%BA%8F.png?raw=true)

### 特性
1. 分布式, 可拓扑的架构
2. 支持单个, 房间推送
3. 心跳支持(gorilla/websocket内置)
4. 基于redis 做消息推送
5. 轻量级

### 部署
```
// build
make build
// run
make run
// stop
make stop
```
### 依赖环境
#### 语言
* go1.13
#### 第三方包
* log: github.com/uber-go/zap
* rpc: github.com/smallnest/rpcx
* websocket: github.com/gorilla/websocket
* config：github.com/spf13/viper
