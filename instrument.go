package gmx

import (
	"golang.org/x/sync/singleflight"
	"log"
	"sync/atomic"
)

type Metric interface {
	Value() interface{}
}

// Non-numeric metrics is reported via MetricFunc
type MetricFunc func() interface{}

func (mf MetricFunc) Value() interface{} {
	return mf()
}

type Counter interface {
	// All Counter implements Metric interface
	Metric

	// Decrease the counter value
	Dec(int64)

	// Increase the counter value
	Inc(int64)

	// Reset the counter value
	Clear()

}

type GenericCounter struct {
	value int64
}

func (c *GenericCounter) Inc(i int64) {
	atomic.AddInt64(&c.value, i)
}

func (c *GenericCounter) Dec(i int64) {
	atomic.AddInt64(&c.value, -i)
}

func (c *GenericCounter) Clear() {
	atomic.StoreInt64(&c.value, 0)
}

func (c *GenericCounter) Value() interface{} {
	return atomic.LoadInt64(&c.value)
}

// NewCounter publishes and returns a new counter instance
func NewCounter(name string) Counter {
	c := &GenericCounter{}
	Publish(name, c)
	return c
}

// Gauge is a gauge metric
type Gauge interface {
	// All Gauge implements Metric interface
	Metric

	// Update set the gauge with a new value
	Update(int64)
}

type GenericGauge struct {
	name string	// optional
	value int64
	valueFunc func() int64	// for lazy value retrieving, optional
}

func (g *GenericGauge) Update(i int64) {
	atomic.StoreInt64(&g.value, i)
}

func (g *GenericGauge) Value() interface{}{

	if g.valueFunc != nil {
		var callGroup singleflight.Group

		v, err, shared := callGroup.Do(g.name, func() (interface{}, error) {
			 return g.valueFunc(), nil

		})
		if err != nil {
			g.value = 0
		}
		log.Printf("Guage return shared value: %t", shared)

		g.value = v.(int64)

	}

	return atomic.LoadInt64(&g.value)
}


// NewGauge publishes and returns a new Gauge instance
func NewGauge(name string) Gauge {
	g := &GenericGauge{name: name}
	Publish(name, g)
	return g
}

// NewGauge publishes and returns a new GenericGauge instance
func NewGaugeWithCallback(name string, valueFunc func() int64) Gauge {
	g := &GenericGauge{valueFunc: valueFunc, name: name}
	Publish(name, g)
	return g
}
