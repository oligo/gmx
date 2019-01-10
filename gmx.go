package gmx

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

const GMX_VERSION = 0


var (
	r = &registry{
		entries: make(map[string]Metric),
	}
)

func init() {
	s, err := setupSocket()

	if err != nil {
		log.Printf("gmx: unable to open socket: %v", err)
		return
	}

	// register the registries keys for discovery
	Publish("keys", (MetricFunc)(func() interface{} {
		return r.keys()
	}))
	go serve(s, r)
}



// Publish registers the metric with the supplied key.
func Publish(key string, metric Metric) {
	r.register(key, metric)
}

func serve(l net.Listener, r *registry) {
	// if listener is a unix socket, try to delete it on shutdown
	if l, ok := l.(*net.UnixListener); ok {
		if a, ok := l.Addr().(*net.UnixAddr); ok {
			defer os.Remove(a.Name)
		}
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			return
		}
		go handle(c, r)
	}
}

func handle(nc net.Conn, reg *registry) {
	// conn makes it easier to send and receive json
	type conn struct {
		net.Conn
		*json.Encoder
		*json.Decoder
	}
	c := conn{
		nc,
		json.NewEncoder(nc),
		json.NewDecoder(nc),
	}
	defer c.Close()
	for {
		var keys []string
		if err := c.Decode(&keys); err != nil {
			if err != io.EOF {
				log.Printf("gmx: client %v sent invalid json request: %v", c.RemoteAddr(), err)
			}
			return
		}
		var result = make(map[string]interface{})
		for _, key := range keys {
			if m, ok := reg.value(key); ok {
				//var metricType string
				//switch m.(type) {
				//case MetricFunc:
				//	metricType = MetricRaw
				//case Counter:
				//	metricType = MetricCounter
				//case Gauge:
				//	metricType = MetricGauge
				//}
				//result[key] = response{Type: metricType, Value: m.Value()}

				// invoke the function for key and store the result

				result[key] = m.Value()

			}
		}
		if err := c.Encode(result); err != nil {
			log.Printf("gmx: could not send response to client %v: %v", c.RemoteAddr(), err)
			return
		}
	}
}

/*
const (
	MetricRaw = "string"
	MetricCounter = "counter"
	MetricGauge = "gauge"
)

type response struct {
	Type string	`json:"type"`
	Value interface{}	`json:"value"`
}

*/

type registry struct {
	sync.Mutex // protects entries from concurrent mutation
	entries    map[string]Metric
}

func (r *registry) register(key string, metric Metric) {
	r.Lock()
	defer r.Unlock()
	r.entries[key] = metric
}

func (r *registry) value(key string) (Metric, bool) {
	r.Lock()
	defer r.Unlock()
	m, ok := r.entries[key]
	return m, ok
}

func (r *registry) keys() (k []string) {
	r.Lock()
	defer r.Unlock()
	for e := range r.entries {
		k = append(k, e)
	}
	return
}
