package main

import (
   
   "net/http"
   "github.com/prometheus/client_golang/prometheus/promhttp"
   "flag"
   "log"
   "github.com/prometheus/client_golang/prometheus"
)

type metricCollector struct {
   CounterMetric *prometheus.Desc
   GaugeMetric *prometheus.Desc
   HistogramMetric *prometheus.Desc
   SummaryMetric *prometheus.Desc
}


func (collect *metricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.CounterMetric
	ch <- collect.GaugeMetric
	ch <- collect.HistogramMetric
	ch <- collect.SummaryMetric
}


func (collect *metricCollector) Collect(ch chan<- prometheus.Metric) {	
	ch <- prometheus.MustNewConstMetric(collect.CounterMetric,prometheus.CounterValue, 1, "192.168.56.100")
	ch <- prometheus.MustNewConstMetric(collect.GaugeMetric,prometheus.GaugeValue, 1, "192.168.56.100")

	//模拟数据
	ch <- prometheus.MustNewConstSummary(
        collect.SummaryMetric,
        4711, 403.34,
        map[float64]float64{0.5: 42.3, 0.9: 323.3},
        "200",
    )

	//模拟数据
    ch <- prometheus.MustNewConstHistogram(
		collect.HistogramMetric,
		4711, 403.34,
		map[float64]uint64{25: 121, 50: 2403, 100: 3221, 200: 4233},
		"200",
		)  
 }

 

// 新建一个采集器
func newMetricCollector() *metricCollector {
	return &metricCollector{
		CounterMetric: prometheus.NewDesc(
			"aaa_metrics",                           // 指标的名称
			"Show metrics for aaa",                  // 帮助信息
			[]string{"host"},                       // label名称数组
			prometheus.Labels{"api_name": "example"}, // labels
		),

		GaugeMetric: prometheus.NewDesc(
			"bbb_metrics",                           // 指标的名称
			"Show metrics for bbb",                  // 帮助信息
			[]string{"host"},                       // label名称数组
			prometheus.Labels{"api_name": "example"}, // labels
		),

		HistogramMetric: prometheus.NewDesc(
			"ccc_metrics",                           // 指标的名称
			"Show metrics for ccc",                  // 帮助信息
			[]string{"host"},                       // label名称数组
			prometheus.Labels{"api_name": "example"}, // labels
		),

		SummaryMetric: prometheus.NewDesc(
			"ddd_metrics",                           // 指标的名称
			"Show metrics for ddd",                  // 帮助信息
			[]string{"host"},                       // label名称数组
			prometheus.Labels{"api_name": "example"}, // labels
		),
	}
}

var (
	// Set during go build
	// version   string
	// gitCommit string

	// 命令行参数
	listenAddr  = flag.String("web.listen-port", "9000", "An port to listen on for web interface and telemetry.")
	metricsPath = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
	metricsNamespace = flag.String("metric.namespace", "ECSDATA", "Prometheus metrics namespace, as the prefix of metrics name")
)

func main() {
   //注册指标
   foo := newMetricCollector()
   registry := prometheus.NewRegistry()
   registry.MustRegister(foo)

     
   
   http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
   
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>A Prometheus Exporter</title></head>
			<body>
			<h1>A Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
   

	log.Printf("Starting Server at http://localhost:%s%s", *listenAddr, *metricsPath)
	log.Fatal(http.ListenAndServe(":"+*listenAddr, nil))
}





