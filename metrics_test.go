package client

import (
	"errors"
	"github.com/rcrowley/go-metrics"
	"math/rand"
	"testing"
	"time"
)

const fanout = 10

func TestStatsD(t *testing.T) {
	r := metrics.NewRegistry()

	c := metrics.NewCounter()
	r.Register("counter", c)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				c.Dec(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				c.Inc(47)
				time.Sleep(400e6)
			}
		}()
	}

	g := metrics.NewGauge()
	r.Register("gauge", g)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				g.Update(19)
				time.Sleep(600e6)
			}
		}()
		go func() {
			for {
				g.Update(47)
				time.Sleep(700e6)
			}
		}()
	}

	gf := metrics.NewGaugeFloat64()
	r.Register("gaugefloat", gf)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				gf.Update(19.2)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				gf.Update(47.3)
				time.Sleep(400e6)
			}
		}()
	}

	hc := metrics.NewHealthcheck(func(h metrics.Healthcheck) {
		if 0 < rand.Intn(2) {
			h.Healthy()
		} else {
			h.Unhealthy(errors.New("baz"))
		}
	})
	r.Register("healthcheck", hc)

	s := metrics.NewExpDecaySample(1028, 0.015)
	//s := metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	r.Register("histogram", h)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				h.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				h.Update(47)
				time.Sleep(400e6)
			}
		}()
	}

	m := metrics.NewMeter()
	r.Register("meter", m)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				m.Mark(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				m.Mark(47)
				time.Sleep(400e6)
			}
		}()
	}

	tt := metrics.NewTimer()
	r.Register("timer", tt)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				tt.Time(func() { time.Sleep(300e6) })
			}
		}()
		go func() {
			for {
				tt.Time(func() { time.Sleep(400e6) })
			}
		}()
	}

	metrics.RegisterDebugGCStats(r)
	go metrics.CaptureDebugGCStats(r, 5e9)

	metrics.RegisterRuntimeMemStats(r)
	go metrics.CaptureRuntimeMemStats(r, 5e9)

	StatsD(r, 1*time.Second, "metrics", "udp", "127.0.0.1:8125", 0)
}
