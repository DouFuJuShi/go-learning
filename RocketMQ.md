# RocketMQ

## offical site
https://rocketmq.apache.org/    

https://github.com/apache/rocketmq

## Go SDK

RocketMQ 4.X https://github.com/apache/rocketmq-client-go   

RocketMQ 5.X https://github.com/apache/rocketmq-clients/golang


## Installtion
### namesrv Setting
```property
rocketmqHome=/Users/aqi/Downloads/rocketmq-all-5.1.3-bin-release
kvConfigPath=/Users/aqi/namesrv/kvConfig.json
configStorePath=/Users/aqi/namesrv/namesrv.properties
productEnvName=center
clusterTest=false
orderMessageEnable=false
returnOrderTopicConfigToBroker=true
clientRequestThreadPoolNums=8
defaultThreadPoolNums=16
clientRequestThreadPoolQueueCapacity=50000
defaultThreadPoolQueueCapacity=10000
scanNotActiveBrokerInterval=5000
unRegisterBrokerQueueCapacity=3000
supportActingMaster=false
enableAllTopicList=true
enableTopicList=true
notifyMinBrokerIdChanged=false
enableControllerInNamesrv=false
needWaitForService=false
waitSecondsForService=45
bindAddress=0.0.0.0
listenPort=9876
serverWorkerThreads=8
serverCallbackExecutorThreads=0
serverSelectorThreads=3
serverOnewaySemaphoreValue=256
serverAsyncSemaphoreValue=64
serverChannelMaxIdleTimeSeconds=120
serverSocketSndBufSize=0
serverSocketRcvBufSize=0
writeBufferHighWaterMark=0
writeBufferLowWaterMark=0
serverSocketBacklog=1024
serverPooledByteBufAllocatorEnable=true
useEpollNativeSelector=false
clientWorkerThreads=4
clientCallbackExecutorThreads=8
clientOnewaySemaphoreValue=65535
clientAsyncSemaphoreValue=65535
connectTimeoutMillis=3000
channelNotActiveInterval=60000
clientChannelMaxIdleTimeSeconds=120
clientSocketSndBufSize=0
clientSocketRcvBufSize=0
clientPooledByteBufAllocatorEnable=false
clientCloseSocketIfTimeout=true
useTLS=false
socksProxyConfig={}
writeBufferHighWaterMark=0
writeBufferLowWaterMark=0
disableCallbackExecutor=false
disableNettyWorkerGroup=false
```


### Proxy Setting

Config Path: /path/to/rocketmq-all-5.1.3-bin-release/conf/rmq-proxy.json   


```json
{
  "rocketMQClusterName": "orderCluster",
  "namesrvAddr": "127.0.0.1:19876;127.0.0.1:29876;127.0.0.1:39876",
  "grpcServerPort": 8081, 
  "remotingListenPort": 8080
}
```


- rocketMQClusterName 设置集群名
- namesrvAddr 设置namesrv地址
- grpcServerPort 设置GRPC监听端口并开启GRPC协议
- remotingListenPort 设置remoting监听端口并且开启remoting监听
