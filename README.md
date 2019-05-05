

# metrics 客户端

数据采集使用[go-metrics](https://github.com/rcrowley/go-metrics)

传输使用UDP, 仿StatsD上传采集数据, InfluxDB进行数据存储, Grafana进行展示。


### 数据封装


```
//挂载配置文件，已修改statsd模版
docker run --ulimit nofile=66000:66000  -v /root/telegraf.conf:/etc/telegraf/telegraf.conf   -d   --name docker-statsd-influxdb-grafana   -p 3003:3003   -p 3004:8888   -p 8086:8086   -p 8125:8125/udp   samuelebistoletti/docker-statsd-influxdb-grafana:latest

```

### register

register 使用的name 必须是不同的


### telegraf 配置修改

将 ` [[inputs.statsd]]` 部分配置打开, 修改templates为:

```
   templates = [
      "* measurement.measurement.field"
   ]
```
表示传值prefix.name.field 最好表示为` prefix_name  field`



## 参考

[multiple field values issues](https://github.com/influxdata/telegraf/issues/2913)