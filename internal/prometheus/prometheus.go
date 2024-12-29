package prometheus

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/gorilla/mux"
	"github.com/mikuta0407/switchbot-contact-sensor_exporter/internal/btle"
	"github.com/mikuta0407/switchbot-contact-sensor_exporter/internal/config"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mutex = sync.RWMutex{}
)

func btScan() {
	d, err := dev.NewDevice(config.Device)
	if err != nil {
		log.Fatalf("Can't create device: %s", err)
	}
	ble.SetDefaultDevice(d)
	defer ble.Stop()
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), config.BTScanDuration))
	chkErr(ble.Scan(ctx, true, btle.Handler, nil))
}

func recordMetrics() {
	go func() {
		i := 0
		for {
			fmt.Println("count", i)
			btScan()
			time.Sleep(config.BTScanInterval)
			i++
		}
	}()
}

func staticPage(w http.ResponseWriter, req *http.Request) {
	page := `<html>
    <head><title>wosensortho exporter</title></head>
    <body>
    <h1>wosensortho exporter</h1>
    <p><a href='metrics'>Metrics</a></p>
    </body>
    </html>`
	fmt.Fprintln(w, page)
}

type sensorCollector struct {
}

func newSensorCollector() *sensorCollector {
	return &sensorCollector{}
}

func (collector *sensorCollector) Describe(ch chan<- *prometheus.Desc) {

}

func buildPromDesc(name string, description string, labels map[string]string) *prometheus.Desc {
	return prometheus.NewDesc(
		name,
		description,
		nil,
		labels,
	)
}

func (collector *sensorCollector) Collect(ch chan<- prometheus.Metric) {

	var desc *prometheus.Desc
	mutex.Lock()
	defer mutex.Unlock()

	for mac, data := range btle.BTDevice {
		labels := make(map[string]string)
		labels["mac"] = mac
		for _, l := range config.Config.Sensors[mac].Labels {
			labels[l.Name] = l.Value
		}
		if len(data.ServiceData) > 0 {

			desc = buildPromDesc("open", "Open", labels)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64((data.ServiceData[3]&0b00000010)>>1))

			desc = buildPromDesc("leave_open", "Leave open", labels)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64((data.ServiceData[3]&0b00000100)>>2))

			desc = buildPromDesc("time", "Time", labels)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(data.ServiceData[7]))

			desc = buildPromDesc("battery", "Battery level", labels)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, data.Battery)
		}
	}
}

func Start(httpListenAddress string) {
	recordMetrics()
	sensor := newSensorCollector()
	prometheus.MustRegister(sensor)
	router := mux.NewRouter()
	router.HandleFunc("/", staticPage)
	http.Handle("/", router)
	router.Path("/metrics").Handler(promhttp.Handler())
	err := http.ListenAndServe(httpListenAddress, router)
	log.Fatal(err)
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
	case context.Canceled:
		log.Printf("Canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
