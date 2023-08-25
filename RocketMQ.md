# RocketMQ

## offical site
https://rocketmq.apache.org/    

https://github.com/apache/rocketmq

## Go SDK

RocketMQ 4.X https://github.com/apache/rocketmq-client-go   

RocketMQ 5.X https://github.com/apache/rocketmq-clients/golang


## Installtion


### Proxy Setting
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
