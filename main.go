package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"net/http"
)

func main() {
	customCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hello_counter",
			Help: "This is my custom counter",
		},
	)

	prometheus.MustRegister(customCounter)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		customCounter.Inc()

		if err := push.New("http://pushgateway:9091", "my_counter_job").
			Collector(customCounter).Push(); err != nil {
			println("Could not push to Pushgateway:", err)
			return c.String(http.StatusInternalServerError, "Could not push to Pushgateway")
		}

		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
