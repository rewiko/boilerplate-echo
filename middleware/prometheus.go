// Inspired by https://github.com/0neSe7en/echo-prometheus/blob/master/monitor.go
// MIT License

// Copyright (c) 2018 WangSiyuan

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// PrometheusConfig contains the configuation for the echo-prometheus
	// middleware.
	PrometheusConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Namespace is single-word prefix relevant to the domain the metric
		// belongs to. For metrics specific to an application, the prefix is
		// usually the application name itself.
		Namespace string
	}
)

var (
	// DefaultPrometheusConfig supplies Prometheus client with the default
	// skipper and the 'echo' namespace.
	DefaultPrometheusConfig = PrometheusConfig{
		Skipper:   middleware.DefaultSkipper,
		Namespace: "http",
	}
)

var (
	httpQueryStatusCode *prometheus.CounterVec
	httpQueryDuration   *prometheus.HistogramVec
)

func initCollector(namespace string) {

	httpQueryStatusCode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "status_code",
			Help:      "Give http response status code",
		}, []string{"code", "method", "host", "path"},
	)
	httpQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "duration_seconds",
			Help:      "Duration of http request",
			Buckets:   []float64{0.025, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 1, 1.5, 2.5},
		}, []string{"code", "method", "host", "path"},
	)

	prometheus.MustRegister(httpQueryDuration, httpQueryStatusCode)
}

// NewMetric returns an echo middleware with the default configuration.
func NewMetric() echo.MiddlewareFunc {
	return NewMetricWithConfig(DefaultPrometheusConfig)
}

// NewMetricWithConfig returns an echo middleware with a custom configuration.
func NewMetricWithConfig(config PrometheusConfig) echo.MiddlewareFunc {
	initCollector(config.Namespace)
	if config.Skipper == nil {
		config.Skipper = DefaultPrometheusConfig.Skipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}
			path := c.Path()
			status := strconv.Itoa(res.Status)
			elapsed := time.Since(start).Seconds()
			httpQueryStatusCode.WithLabelValues(status, req.Method, req.Host, path).Inc()
			httpQueryDuration.WithLabelValues(status, req.Method, req.Host, path).Observe(elapsed)
			return nil
		}
	}
}
