package client

import (
	"bufio"
	"bytes"
	"github.com/rcrowley/go-metrics"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	// UDP packet limit, see
	// https://en.wikipedia.org/wiki/User_Datagram_Protocol#Packet_structure
	UDP_MAX_PACKET_SIZE int = 64 * 1024
)

// Config provides a container with configuration parameters for
// the StatsD exporter
type Config struct {
	Network       string           // Network: tcp, udp.
	Addr          string           // Network address to connect to | 地址
	Registry      metrics.Registry // Registry to be exported | metrics注册
	FlushInterval time.Duration    // Flush interval | 刷新间隔时间
	Prefix        string           // Prefix to be prepended to metric names | 前缀名字
	Rate          float32          // Rate
	Tags          string           // tag //TODO

	conn net.Conn
}

func StatsD(r metrics.Registry, d time.Duration, prefix string, network string, addr string, rate float32) {

	conn, err := net.Dial(network, addr)
	if err != nil {
		panic("conn remote err!")
	}

	StatsDWithConfig(Config{
		Network:       network,
		Addr:          addr,
		Registry:      r,
		FlushInterval: d,
		Prefix:        prefix,
		Rate:          rate,
		conn:          conn,
	})

}

// WithConfig is a blocking exporter function
func StatsDWithConfig(c Config) {
	for _ = range time.Tick(c.FlushInterval) {
		if err := statsd(&c); err != nil {
			log.Println(err)
			c.conn.Close()
		}
	}
}

func statsd(c *Config) (err error) {

	w := bufio.NewWriter(c.conn)

	c.Registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			ms := metric.Snapshot()
			w.Write(statsdLine(c.Prefix, name, "", ms.Count(), "|c", c.Rate))
		case metrics.Gauge:
			ms := metric.Snapshot()
			w.Write(statsdLine(c.Prefix, name, "", ms.Value(), "|g", c.Rate))
		case metrics.GaugeFloat64:
			ms := metric.Snapshot()
			w.Write(statsdLine(c.Prefix, name, "", ms.Value(), "|g", c.Rate))
		case metrics.Histogram:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})

			fields := make([][]byte, 12)
			fields = append(fields, statsdLine(c.Prefix, name, "count", ms.Count(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "max", ms.Max(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "mean", ms.Mean(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "min", ms.Min(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "stddev", ms.StdDev(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "variance", ms.Variance(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p50", ps[0], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p75", ps[1], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p95", ps[2], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p99", ps[3], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p999", ps[4], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p9999", ps[5], "|g", c.Rate))

			buf := bytes.Join(fields, []byte{})
			w.Write(buf)
		case metrics.Meter:
			ms := metric.Snapshot()
			fields := make([][]byte, 5)
			fields = append(fields, statsdLine(c.Prefix, name, "count", ms.Count(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "m1", ms.Rate1(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "m5", ms.Rate5(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "m15", ms.Rate15(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "mean", ms.RateMean(), "|g", c.Rate))

			buf := bytes.Join(fields, []byte{})
			w.Write(buf)

		case metrics.Timer:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})

			fields := make([][]byte, 12)
			fields = append(fields, statsdLine(c.Prefix, name, "count", ms.Count(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "max", ms.Max(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "mean", ms.Mean(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "min", ms.Min(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "stddev", ms.StdDev(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "variance", ms.Variance(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p50", ps[0], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p75", ps[1], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p95", ps[2], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p99", ps[3], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p999", ps[4], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "p9999", ps[5], "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "m1", ms.Rate1(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "m5", ms.Rate5(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "m15", ms.Rate15(), "|g", c.Rate))
			fields = append(fields, statsdLine(c.Prefix, name, "mean", ms.RateMean(), "|g", c.Rate))

			buf := bytes.Join(fields, []byte{})

			w.Write(buf)

			//case metrics.Healthcheck:
			//	metric.Check()
			//	log.Printf("healthcheck %s\n", name)
			//	log.Printf("  error:       %v\n", metric.Error())
			//
			//case metrics.EWMA:
			//case metrics.Sample:

		}

		w.Flush()
	})

	return
}

//构造发送line
func statsdLine(prefix, name, field string, value interface{}, suffix string, rate float32) []byte {

	//<metricname>:<value>|<type>|@<rate>
	var buffer bytes.Buffer

	//buf := make([]byte, UDP_MAX_PACKET_SIZE)

	//添加前缀
	if prefix != "" {
		//buf = append(buf, prefix...)
		//buf = append(buf, '.')
		buffer.WriteString(prefix)
		buffer.WriteString(".")
	} else {
		buffer.WriteString("statsd")
		buffer.WriteString(".")
	}

	////添加名称
	//buf = append(buf, name...)
	//buf = append(buf, ':')

	//将name注册中的'.'替换成'_', 配合telegraf修改模版,防止将数据库名字改为属性
	if strings.Contains(name, ".") {
		name = strings.ReplaceAll(name, ".", "_")
	}
	//添加名称
	buffer.WriteString(name)

	if field != "" {
		buffer.WriteString(".")
		buffer.WriteString(field)
	}
	buffer.WriteString(":")

	buf := buffer.Bytes()

	switch v := value.(type) {
	case string:
		buf = append(buf, v...)
	case int64:
		buf = strconv.AppendInt(buf, v, 10)
	case float64:
		buf = strconv.AppendFloat(buf, v, 'f', -1, 64)
	default:
		return nil
	}

	if suffix != "" {
		buf = append(buf, suffix...)
	}

	if rate != 0 && rate < 1 {
		buf = append(buf, "|@"...)
		buf = strconv.AppendFloat(buf, float64(rate), 'f', 6, 32)
	}

	buf = append(buf, "\n"...) //每一行打一个回车，telegraf 是使用回车进行读取的
	return buf

}
