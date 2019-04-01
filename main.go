package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
	"github.com/chenleji/nautilus/helper"
	_ "github.com/chenleji/nautilus/routers"
	"github.com/chenleji/nautilus/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

const (
	AppHealthCheckURL = "/health"
)

func init() {
	beego.SetLogFuncCall(true)

	// fetch config
	registerConfigCallBack()
}

func main() {
	// consul
	registerService()

	// swagger
	bootSwaggerDoc()

	// healthCheck
	toolbox.AddHealthCheck("health", &helper.HealthCheck{})

	// prometheus
	startMetricMonitor()

	// leaderElection
	bootLeaderElection()

	beego.Run()
}

func registerConfigCallBack() {
	confFile := &utils.ConfigFile{}
	conf := helper.GetConfInst(confFile)

	// register config there
	conf.Register("test", func(delta *helper.ConfDelta) {
		logs.Info("delta new data:", delta.NData)
		logs.Info("delta old data:", delta.OData)
	})

	// run
	conf.Run()
}

func registerService() {
	if beego.BConfig.RunMode == beego.PROD {
		err := helper.Consul{}.New().RegistryService(
			helper.Utils{}.GetAppName(),
			helper.Utils{}.GetAppPort(),
			AppHealthCheckURL)
		if err != nil {
			panic(err)
		}
	}
}

func bootSwaggerDoc() {
	if beego.BConfig.RunMode == beego.DEV {
		beego.SetStaticPath("/swagger", "swagger")
	}
}

func bootLeaderElection() {
	if beego.BConfig.RunMode == beego.PROD {
		le := &helper.LeaderElection{
			Consul: helper.Consul{}.New(),
			TTL:    time.Second,
			Callback: func(leader bool) {
				logs.Info("is leader:", leader)
			},
		}

		go le.Run()
	}
}

func startMetricMonitor() {
	var DriverJob = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "driver",
			Name:      "total_jobs",
			Help:      "The total number of job running on this node.",
		},
		[]string{"node_name"},
	)

	prometheus.MustRegister(DriverJob)
	DriverJob.WithLabelValues("node-1").Inc()

	beego.Handler("/metrics", promhttp.Handler())
}
