package prometheus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const sysHTTPDefaultTimeout = 3 * time.Hour // дефолтный таймаут системных ручек

type Metrics struct {
	s *http.Server
}

func New(address string) *Metrics {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Addr:         address,
		WriteTimeout: sysHTTPDefaultTimeout,
		ReadTimeout:  sysHTTPDefaultTimeout,
		IdleTimeout:  sysHTTPDefaultTimeout,
		Handler:      mux,
	}

	return &Metrics{
		s: s,
	}
}

func (m *Metrics) Run() error {
	err := m.s.ListenAndServe()
	if err != nil {
		return fmt.Errorf("[srv.ListenAndServe]: %w", err)
	}

	return nil
}
