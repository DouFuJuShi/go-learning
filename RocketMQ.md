# RocketMQ

## Offical Site
https://rocketmq.apache.org/    

https://github.com/apache/rocketmq

## Go SDK

RocketMQ 4.X https://github.com/apache/rocketmq-client-go   

RocketMQ 5.X https://github.com/apache/rocketmq-clients/golang


## Installtion

```shell
# rocketmq-all-5.1.3-bin-release tree
.
├── LICENSE
├── NOTICE
├── README.md
├── benchmark
│   ├── batchproducer.sh
│   ├── consumer.sh
│   ├── producer.sh
│   ├── runclass.sh
│   ├── shutdown.sh
│   └── tproducer.sh
├── bin
│   ├── README.md
│   ├── cachedog.sh
│   ├── cleancache.sh
│   ├── cleancache.v1.sh
│   ├── controller
│   │   ├── fast-try-independent-deployment.cmd
│   │   ├── fast-try-independent-deployment.sh
│   │   ├── fast-try-namesrv-plugin.cmd
│   │   ├── fast-try-namesrv-plugin.sh
│   │   ├── fast-try.cmd
│   │   └── fast-try.sh
│   ├── dledger
│   │   └── fast-try.sh
│   ├── export.sh
│   ├── mqadmin # admin tool
│   ├── mqadmin.cmd
│   ├── mqbroker # 部署Broker服务器
│   ├── mqbroker.cmd
│   ├── mqbroker.numanode0
│   ├── mqbroker.numanode1
│   ├── mqbroker.numanode2
│   ├── mqbroker.numanode3
│   ├── mqbrokercontainer
│   ├── mqcontroller
│   ├── mqcontroller.cmd
│   ├── mqnamesrv # name server
│   ├── mqnamesrv.cmd
│   ├── mqproxy 
│   ├── mqproxy.cmd
│   ├── mqshutdown # 关闭运行中的 broker proxy namesrv controller
│   ├── mqshutdown.cmd
│   ├── os.sh # 在部署Broker服务器之前，强烈建议运行**os.sh**，这是为了优化您的操作系统以获得更好的性能。os.sh参数设置仅供参考。您可以根据目标主机配置调整它们。
│   ├── play.cmd
│   ├── play.sh
│   ├── runbroker.cmd
│   ├── runbroker.sh
│   ├── runserver.cmd
│   ├── runserver.sh
│   ├── setcache.sh
│   ├── startfsrv.sh
│   ├── tools.cmd
│   └── tools.sh
├── conf
│   ├── 2m-2s-async
│   │   ├── broker-a-s.properties
│   │   ├── broker-a.properties
│   │   ├── broker-b-s.properties
│   │   └── broker-b.properties
│   ├── 2m-2s-sync
│   │   ├── broker-a-s.properties
│   │   ├── broker-a.properties
│   │   ├── broker-b-s.properties
│   │   └── broker-b.properties
│   ├── 2m-noslave
│   │   ├── broker-a.properties
│   │   ├── broker-b.properties
│   │   └── broker-trace.properties
│   ├── broker.conf
│   ├── container
│   │   └── 2container-2m-2s
│   │       ├── broker-a-in-container1.conf
│   │       ├── broker-a-in-container2.conf
│   │       ├── broker-b-in-container1.conf
│   │       ├── broker-b-in-container2.conf
│   │       ├── broker-container1.conf
│   │       ├── broker-container2.conf
│   │       └── nameserver.conf
│   ├── controller
│   │   ├── cluster-3n-independent
│   │   │   ├── controller-n0.conf
│   │   │   ├── controller-n1.conf
│   │   │   └── controller-n2.conf
│   │   ├── cluster-3n-namesrv-plugin
│   │   │   ├── namesrv-n0.conf
│   │   │   ├── namesrv-n1.conf
│   │   │   └── namesrv-n2.conf
│   │   ├── controller-standalone.conf
│   │   └── quick-start
│   │       ├── broker-n0.conf
│   │       ├── broker-n1.conf
│   │       └── namesrv.conf
│   ├── dledger
│   │   ├── broker-n0.conf
│   │   ├── broker-n1.conf
│   │   └── broker-n2.conf
│   ├── plain_acl.yml
│   ├── rmq-proxy.json # proxy的配置文件
│   ├── rmq.broker.logback.xml
│   ├── rmq.client.logback.xml
│   ├── rmq.controller.logback.xml
│   ├── rmq.namesrv.logback.xml
│   ├── rmq.proxy.logback.xml
│   ├── rmq.tools.logback.xml
│   └── tools.yml
└── lib
    ├── animal-sniffer-annotations-1.21.jar
    ├── annotations-13.0.jar
    ├── annotations-4.1.1.4.jar
    ├── annotations-api-6.0.53.jar
    ├── awaitility-4.1.0.jar
    ├── bcpkix-jdk15on-1.69.jar
    ├── bcprov-jdk15on-1.69.jar
    ├── bcutil-jdk15on-1.69.jar
    ├── caffeine-2.9.3.jar
    ├── checker-qual-3.12.0.jar
    ├── commons-beanutils-1.9.4.jar
    ├── commons-cli-1.5.0.jar
    ├── commons-codec-1.13.jar
    ├── commons-collections-3.2.2.jar
    ├── commons-digester-2.1.jar
    ├── commons-io-2.7.jar
    ├── commons-lang3-3.12.0.jar
    ├── commons-logging-1.2.jar
    ├── commons-validator-1.7.jar
    ├── concurrentlinkedhashmap-lru-1.4.2.jar
    ├── disruptor-1.2.10.jar
    ├── dledger-0.3.1.2.jar
    ├── error_prone_annotations-2.14.0.jar
    ├── failureaccess-1.0.1.jar
    ├── fastjson-1.2.83.jar
    ├── grpc-api-1.50.0.jar
    ├── grpc-context-1.50.0.jar
    ├── grpc-core-1.50.0.jar
    ├── grpc-netty-shaded-1.50.0.jar
    ├── grpc-protobuf-1.50.0.jar
    ├── grpc-protobuf-lite-1.50.0.jar
    ├── grpc-services-1.50.0.jar
    ├── grpc-stub-1.50.0.jar
    ├── gson-2.9.0.jar
    ├── guava-31.1-jre.jar
    ├── hamcrest-2.1.jar
    ├── j2objc-annotations-1.3.jar
    ├── jaeger-thrift-1.6.0.jar
    ├── jaeger-tracerresolver-1.6.0.jar
    ├── javassist-3.20.0-GA.jar
    ├── javax.annotation-api-1.3.2.jar
    ├── jna-4.2.2.jar
    ├── jsr305-3.0.2.jar
    ├── jul-to-slf4j-2.0.6.jar
    ├── kotlin-stdlib-1.6.20.jar
    ├── kotlin-stdlib-common-1.6.20.jar
    ├── kotlin-stdlib-jdk7-1.6.20.jar
    ├── kotlin-stdlib-jdk8-1.6.20.jar
    ├── libthrift-0.14.1.jar
    ├── listenablefuture-9999.0-empty-to-avoid-conflict-with-guava.jar
    ├── lz4-java-1.8.0.jar
    ├── netty-all-4.1.65.Final.jar
    ├── netty-tcnative-boringssl-static-2.0.53.Final-linux-aarch_64.jar
    ├── netty-tcnative-boringssl-static-2.0.53.Final-linux-x86_64.jar
    ├── netty-tcnative-boringssl-static-2.0.53.Final-osx-aarch_64.jar
    ├── netty-tcnative-boringssl-static-2.0.53.Final-osx-x86_64.jar
    ├── netty-tcnative-boringssl-static-2.0.53.Final-windows-x86_64.jar
    ├── netty-tcnative-boringssl-static-2.0.53.Final.jar
    ├── netty-tcnative-classes-2.0.53.Final.jar
    ├── okhttp-4.11.0.jar
    ├── okio-3.2.0.jar
    ├── okio-jvm-3.0.0.jar
    ├── openmessaging-api-0.3.1-alpha.jar
    ├── opentelemetry-api-1.26.0.jar
    ├── opentelemetry-api-events-1.26.0-alpha.jar
    ├── opentelemetry-api-logs-1.26.0-alpha.jar
    ├── opentelemetry-context-1.26.0.jar
    ├── opentelemetry-exporter-common-1.26.0.jar
    ├── opentelemetry-exporter-logging-1.26.0.jar
    ├── opentelemetry-exporter-otlp-1.26.0.jar
    ├── opentelemetry-exporter-otlp-common-1.26.0.jar
    ├── opentelemetry-exporter-prometheus-1.26.0-alpha.jar
    ├── opentelemetry-extension-incubator-1.26.0-alpha.jar
    ├── opentelemetry-sdk-1.26.0.jar
    ├── opentelemetry-sdk-common-1.26.0.jar
    ├── opentelemetry-sdk-extension-autoconfigure-spi-1.26.0.jar
    ├── opentelemetry-sdk-logs-1.26.0-alpha.jar
    ├── opentelemetry-sdk-metrics-1.26.0.jar
    ├── opentelemetry-sdk-trace-1.26.0.jar
    ├── opentelemetry-semconv-1.26.0-alpha.jar
    ├── opentracing-noop-0.33.0.jar
    ├── opentracing-tracerresolver-0.1.8.jar
    ├── opentracing-util-0.33.0.jar
    ├── perfmark-api-0.25.0.jar
    ├── proto-google-common-protos-2.9.0.jar
    ├── protobuf-java-3.20.1.jar
    ├── protobuf-java-util-3.20.1.jar
    ├── rocketmq-acl-5.1.3.jar
    ├── rocketmq-broker-5.1.3.jar
    ├── rocketmq-client-5.1.3.jar
    ├── rocketmq-common-5.1.3.jar
    ├── rocketmq-container-5.1.3.jar
    ├── rocketmq-controller-5.1.3.jar
    ├── rocketmq-example-5.1.3.jar
    ├── rocketmq-filter-5.1.3.jar
    ├── rocketmq-logback-classic-1.0.1.jar
    ├── rocketmq-namesrv-5.1.3.jar
    ├── rocketmq-openmessaging-5.1.3.jar
    ├── rocketmq-proto-2.0.2.jar
    ├── rocketmq-proxy-5.1.3.jar
    ├── rocketmq-remoting-5.1.3.jar
    ├── rocketmq-shaded-slf4j-api-bridge-1.0.0.jar
    ├── rocketmq-slf4j-api-1.0.1.jar
    ├── rocketmq-srvutil-5.1.3.jar
    ├── rocketmq-store-5.1.3.jar
    ├── rocketmq-tiered-store-5.1.3.jar
    ├── rocketmq-tools-5.1.3.jar
    ├── slf4j-api-2.0.3.jar
    ├── snakeyaml-1.32.jar
    ├── tomcat-annotations-api-8.5.46.jar
    ├── tomcat-embed-core-8.5.46.jar
    └── zstd-jni-1.5.2-2.jar
```
### namesrv Setting
```properties
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

### broker Setting
```properties
brokerConfigPath=
rocketmqHome=/Users/aqi/Downloads/rocketmq-all-5.1.3-bin-release
namesrvAddr=
listenPort=6888
brokerIP1=192.168.20.140
brokerIP2=192.168.20.140
recoverConcurrently=false
brokerPermission=6
defaultTopicQueueNums=8
autoCreateTopicEnable=true # 自动创建Topic 建议线上设置为false
clusterTopicEnable=true
brokerTopicEnable=true
autoCreateSubscriptionGroup=true
messageStorePlugIn=
msgTraceTopicName=RMQ_SYS_TRACE_TOPIC
traceTopicEnable=false
sendMessageThreadPoolNums=4
putMessageFutureThreadPoolNums=4
pullMessageThreadPoolNums=32
litePullMessageThreadPoolNums=32
ackMessageThreadPoolNums=3
processReplyMessageThreadPoolNums=32
queryMessageThreadPoolNums=16
adminBrokerThreadPoolNums=16
clientManageThreadPoolNums=32
consumerManageThreadPoolNums=32
loadBalanceProcessorThreadPoolNums=32
heartbeatThreadPoolNums=8
recoverThreadPoolNums=32
endTransactionThreadPoolNums=24
flushConsumerOffsetInterval=5000
flushConsumerOffsetHistoryInterval=60000
rejectTransactionMessage=false
fetchNameSrvAddrByDnsLookup=false
fetchNamesrvAddrByAddressServer=false
sendThreadPoolQueueCapacity=10000
putThreadPoolQueueCapacity=10000
pullThreadPoolQueueCapacity=100000
litePullThreadPoolQueueCapacity=100000
ackThreadPoolQueueCapacity=100000
replyThreadPoolQueueCapacity=10000
queryThreadPoolQueueCapacity=20000
clientManagerThreadPoolQueueCapacity=1000000
consumerManagerThreadPoolQueueCapacity=1000000
heartbeatThreadPoolQueueCapacity=50000
endTransactionPoolQueueCapacity=100000
adminBrokerThreadPoolQueueCapacity=10000
loadBalanceThreadPoolQueueCapacity=100000
longPollingEnable=true
shortPollingTimeMills=1000
notifyConsumerIdsChangedEnable=true
highSpeedMode=false
commercialBaseCount=1
commercialSizePerMsg=4096
accountStatsEnable=true
accountStatsPrintZeroValues=true
transferMsgByHeap=true
regionId=DefaultRegion
registerBrokerTimeoutMills=24000
sendHeartbeatTimeoutMillis=1000
slaveReadEnable=false
disableConsumeIfConsumerReadSlowly=false
consumerFallbehindThreshold=17179869184
brokerFastFailureEnable=true
waitTimeMillsInSendQueue=200
waitTimeMillsInPullQueue=5000
waitTimeMillsInLitePullQueue=5000
waitTimeMillsInHeartbeatQueue=31000
waitTimeMillsInTransactionQueue=3000
waitTimeMillsInAckQueue=3000
startAcceptSendRequestTimeStamp=0
traceOn=true
enableCalcFilterBitMap=false
rejectPullConsumerEnable=false
expectConsumerNumUseFilter=32
maxErrorRateOfBloomFilter=20
filterDataCleanTimeSpan=86400000
filterSupportRetry=false
enablePropertyFilter=false
compressedRegister=false
forceRegister=true
registerNameServerPeriod=30000
brokerHeartbeatInterval=1000
brokerNotActiveTimeoutMillis=10000
enableNetWorkFlowControl=false
enableBroadcastOffsetStore=true
broadcastOffsetExpireSecond=120
broadcastOffsetExpireMaxSecond=300
popPollingSize=1024
popPollingMapSize=100000
maxPopPollingSize=100000
reviveQueueNum=8
reviveInterval=1000
reviveMaxSlow=3
reviveScanTime=10000
enableSkipLongAwaitingAck=false
reviveAckWaitMs=180000
enablePopLog=false
enablePopBufferMerge=false
popCkStayBufferTime=10000
popCkStayBufferTimeOut=3000
popCkMaxBufferSize=200000
popCkOffsetMaxQueueSize=20000
enablePopBatchAck=false
enableNotifyAfterPopOrderLockRelease=true
realTimeNotifyConsumerChange=true
litePullMessageEnable=true
syncBrokerMemberGroupPeriod=1000
loadBalancePollNameServerInterval=30000
cleanOfflineBrokerInterval=30000
serverLoadBalancerEnable=true
defaultMessageRequestMode=PULL
defaultPopShareQueueNum=-1
transactionTimeOut=6000
transactionCheckMax=15
transactionCheckInterval=30000
transactionOpMsgMaxSize=4096
transactionOpBatchInterval=3000
aclEnable=false # 开启 ACL 权限控制  需要结合 /conf/plain_acl.yml 进行设置
storeReplyMessageEnable=true
enableDetailStat=true
autoDeleteUnusedStats=false
isolateLogEnable=false
forwardTimeout=3000
enableSlaveActingMaster=false
enableRemoteEscape=false
skipPreOnline=false
asyncSendEnable=true
useServerSideResetOffset=false
consumerOffsetUpdateVersionStep=500
delayOffsetUpdateVersionStep=200
lockInStrictMode=false
compatibleWithOldNameSrv=true
enableControllerMode=false
controllerAddr=
fetchControllerAddrByDnsLookup=false
syncBrokerMetadataPeriod=5000
checkSyncStateSetPeriod=5000
syncControllerMetadataPeriod=10000
controllerHeartBeatTimeoutMills=10000
validateSystemTopicWhenUpdateTopic=true
brokerElectionPriority=2147483647
useStaticSubscription=false
metricsExporterType=DISABLE
metricsGrpcExporterTarget=
metricsGrpcExporterHeader=
metricGrpcExporterTimeOutInMills=3000
metricGrpcExporterIntervalInMills=60000
metricLoggingExporterIntervalInMills=10000
metricsPromExporterPort=5557
metricsPromExporterHost=
metricsLabel=
metricsInDelta=false
channelExpiredTimeout=120000
subscriptionExpiredTimeout=600000
estimateAccumulation=true
coldCtrStrategyEnable=false
usePIDColdCtrStrategy=true
cgColdReadThreshold=3145728
globalColdReadThreshold=104857600
fetchNamesrvAddrInterval=10000
bindAddress=0.0.0.0
listenPort=10911
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
storePathRootDir=/Users/aqi/store
storePathCommitLog=
storePathDLedgerCommitLog=
storePathEpochFile=
storePathBrokerIdentity=
readOnlyCommitLogStorePaths=
mappedFileSizeCommitLog=1073741824
compactionMappedFileSize=104857600
compactionCqMappedFileSize=10485760
compactionScheduleInternal=900000
maxOffsetMapSize=104857600
compactionThreadNum=6
enableCompaction=true
mappedFileSizeTimerLog=104857600
timerPrecisionMs=1000
timerRollWindowSlot=172800
timerFlushIntervalMs=1000
timerGetMessageThreadNum=3
timerPutMessageThreadNum=3
timerEnableDisruptor=false
timerEnableCheckMetrics=true
timerInterceptDelayLevel=false
timerMaxDelaySec=259200
timerWheelEnable=true
disappearTimeAfterStart=-1
timerStopEnqueue=false
timerCheckMetricsWhen=05
timerSkipUnknownError=false
timerWarmEnable=false
timerStopDequeue=false
timerCongestNumEachSlot=2147483647
timerMetricSmallThreshold=1000000
timerProgressLogIntervalMs=10000
mappedFileSizeConsumeQueue=6000000
enableConsumeQueueExt=false
mappedFileSizeConsumeQueueExt=50331648
mapperFileSizeBatchConsumeQueue=13800000
bitMapLengthConsumeQueueExt=64
flushIntervalCommitLog=500
commitIntervalCommitLog=200
maxRecoveryCommitlogFiles=30
diskSpaceWarningLevelRatio=90
diskSpaceCleanForciblyRatio=85
useReentrantLockWhenPutMessage=true
flushCommitLogTimed=true
flushIntervalConsumeQueue=1000
cleanResourceInterval=10000
deleteCommitLogFilesInterval=100
deleteConsumeQueueFilesInterval=100
destroyMapedFileIntervalForcibly=120000
redeleteHangedFileInterval=120000
deleteWhen=04
diskMaxUsedSpaceRatio=75
fileReservedTime=72
deleteFileBatchMax=10
putMsgIndexHightWater=600000
maxMessageSize=4194304
checkCRCOnRecover=true
flushCommitLogLeastPages=4
commitCommitLogLeastPages=4
flushLeastPagesWhenWarmMapedFile=4096
flushConsumeQueueLeastPages=2
flushCommitLogThoroughInterval=10000
commitCommitLogThoroughInterval=200
flushConsumeQueueThoroughInterval=60000
maxTransferBytesOnMessageInMemory=262144
maxTransferCountOnMessageInMemory=32
maxTransferBytesOnMessageInDisk=65536
maxTransferCountOnMessageInDisk=8
accessMessageInMemoryMaxRatio=40
messageIndexEnable=true
maxHashSlotNum=5000000
maxIndexNum=20000000
maxMsgsNumBatch=64
messageIndexSafe=false
haListenPort=10912
haSendHeartbeatInterval=5000
haHousekeepingInterval=20000
haTransferBatchSize=32768
haMasterAddress=
haMaxGapNotInSync=268435456
brokerRole=ASYNC_MASTER
flushDiskType=ASYNC_FLUSH
syncFlushTimeout=5000
putMessageTimeout=8000
slaveTimeout=3000
messageDelayLevel=1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h # 延迟消息级别
flushDelayOffsetInterval=10000
cleanFileForciblyEnable=true
warmMapedFileEnable=false
offsetCheckInSlave=false
debugLockEnable=false
duplicationEnable=false
diskFallRecorded=true
osPageCacheBusyTimeOutMills=1000
defaultQueryMaxNum=32
transientStorePoolEnable=false
transientStorePoolSize=5
fastFailIfNoBufferInStorePool=false
enableDLegerCommitLog=false
dLegerGroup=
dLegerPeers=
dLegerSelfId=
preferredLeaderId=
enableBatchPush=false
enableScheduleMessageStats=true
enableLmq=false
enableMultiDispatch=false
maxLmqConsumeQueueNum=20000
enableScheduleAsyncDeliver=false
scheduleAsyncDeliverMaxPendingLimit=2000
scheduleAsyncDeliverMaxResendNum2Blocked=3
maxBatchDeleteFilesNum=50
dispatchCqThreads=10
dispatchCqCacheNum=4096
enableAsyncReput=true
recheckReputOffsetFromCq=false
maxTopicLength=127
autoMessageVersionOnTopicLen=true
travelCqFileNumWhenGetMessage=1
correctLogicMinOffsetSleepInterval=1
correctLogicMinOffsetForceInterval=300000
mappedFileSwapEnable=true
commitLogForceSwapMapInterval=43200000
commitLogSwapMapInterval=3600000
commitLogSwapMapReserveFileNum=100
logicQueueForceSwapMapInterval=43200000
logicQueueSwapMapInterval=3600000
cleanSwapedMapInterval=300000
logicQueueSwapMapReserveFileNum=20
searchBcqByCacheEnable=true
dispatchFromSenderThread=false
wakeCommitWhenPutMessage=true
wakeFlushWhenPutMessage=false
enableCleanExpiredOffset=false
maxAsyncPutMessageRequests=5000
pullBatchMaxMessageCount=160
totalReplicas=1
inSyncReplicas=1
minInSyncReplicas=1
allAckInSyncStateSet=false
enableAutoInSyncReplicas=false
haFlowControlEnable=false
maxHaTransferByteInSecond=104857600
haMaxTimeSlaveNotCatchup=15000
syncMasterFlushOffsetWhenStartup=false
maxChecksumRange=1073741824
replicasPerDiskPartition=1
logicalDiskSpaceCleanForciblyThreshold=0.8
maxSlaveResendLength=268435456
syncFromLastFile=false
asyncLearner=false
maxConsumeQueueScan=20000
sampleCountThreshold=5000
coldDataFlowControlEnable=false
coldDataScanEnable=false
dataReadAheadEnable=false
timerColdDataCheckIntervalMs=60000
sampleSteps=32
accessMessageInMemoryHotRatio=26
enableBuildConsumeQueueConcurrently=false
batchDispatchRequestThreadPoolNums=16
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


### ACL 权限控制
1. Broker 配置
   ```properties
   aclEnable = true # 默认 false，开启 ACL 需要设置为 true
   ```
2. plain_acl.yml 配置
   默认 /path/to/rocketmq-all-5.1.3-bin-release/conf/plain_acl.yml，可以通过 -Drocketmq.acl.plain.file 指定 ACL 文件名称
```yaml
    globalWhiteRemoteAddresses:
      - 10.10.103.*
      - 192.168.0.*
    
    accounts:
      - accessKey: RocketMQ
        secretKey: 12345678
        whiteRemoteAddress:
        admin: false
        defaultTopicPerm: DENY
        defaultGroupPerm: SUB
        topicPerms:
          - topicA=DENY
          - topicB=PUB|SUB
          - topicC=SUB
        groupPerms:
          # the group should convert to retry topic
          - groupA=DENY
          - groupB=PUB|SUB
          - groupC=SUB
    
      - accessKey: rocketmq2
        secretKey: 12345678
        whiteRemoteAddress: 192.168.1.*
        # if it is admin, it could access all resources
        admin: true
```




| 参数名 | 说明 |
| :--- | :--- |
| globalWhiteRemoteAddresses | 全局白名单配置，策略如下： 空：忽略白名单，继续执行下面校验 全匹配模式：全部放行不会执行后面校验 例如：* 或 ... 或 ::::::: 多 IP 模式：表示白名单 IP 在设置区间段的放行 例如：192.168.0.{1,2} 或 192.168.1.1,192.168.1.2 或 192.168.*. 或 192.168.1-10.5-50 |
| accessKey | 用户唯一标识 |
| secretKey | 访问密码 |
| whiteRemoteAddress | 用户级白名单，格式同 globalWhiteRemoteAddresses |
| admin | 是否为管理员，管理员拥有所有资源访问权限 true or false |
| secretKey | 访问密码 |
| defaultTopicPerm | 默认主题权限，默认值 DENY |
| defaultGroupPerm | 默认消费组权限，默认值 DENY |
| topicPerms | 详细的主题权限 |
| groupPerms | 详细的消费组权限 |
