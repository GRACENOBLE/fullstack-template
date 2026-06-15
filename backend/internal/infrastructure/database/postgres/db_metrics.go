package postgres

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type dbStatsCollector struct {
	db *sql.DB

	maxOpenConns      *prometheus.Desc
	openConns         *prometheus.Desc
	inUse             *prometheus.Desc
	idle              *prometheus.Desc
	waitCount         *prometheus.Desc
	waitDuration      *prometheus.Desc
	maxIdleClosed     *prometheus.Desc
	maxIdleTimeClosed *prometheus.Desc
	maxLifetimeClosed *prometheus.Desc
}

// NewDBStatsCollector returns a prometheus.Collector that exports sql.DBStats.
func NewDBStatsCollector(db *sql.DB) prometheus.Collector {
	return &dbStatsCollector{
		db: db,
		maxOpenConns: prometheus.NewDesc(
			"db_pool_max_open_connections",
			"Maximum number of open connections to the database.",
			nil, nil,
		),
		openConns: prometheus.NewDesc(
			"db_pool_open_connections",
			"Current number of open connections including in-use and idle.",
			nil, nil,
		),
		inUse: prometheus.NewDesc(
			"db_pool_in_use_connections",
			"Current number of connections in use.",
			nil, nil,
		),
		idle: prometheus.NewDesc(
			"db_pool_idle_connections",
			"Current number of idle connections.",
			nil, nil,
		),
		waitCount: prometheus.NewDesc(
			"db_pool_wait_count_total",
			"Total number of connections waited for.",
			nil, nil,
		),
		waitDuration: prometheus.NewDesc(
			"db_pool_wait_duration_seconds_total",
			"Total time blocked waiting for a new connection.",
			nil, nil,
		),
		maxIdleClosed: prometheus.NewDesc(
			"db_pool_max_idle_closed_total",
			"Total connections closed due to SetMaxIdleConns.",
			nil, nil,
		),
		maxIdleTimeClosed: prometheus.NewDesc(
			"db_pool_max_idle_time_closed_total",
			"Total connections closed due to SetConnMaxIdleTime.",
			nil, nil,
		),
		maxLifetimeClosed: prometheus.NewDesc(
			"db_pool_max_lifetime_closed_total",
			"Total connections closed due to SetConnMaxLifetime.",
			nil, nil,
		),
	}
}

func (c *dbStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenConns
	ch <- c.openConns
	ch <- c.inUse
	ch <- c.idle
	ch <- c.waitCount
	ch <- c.waitDuration
	ch <- c.maxIdleClosed
	ch <- c.maxIdleTimeClosed
	ch <- c.maxLifetimeClosed
}

func (c *dbStatsCollector) Collect(ch chan<- prometheus.Metric) {
	s := c.db.Stats()
	ch <- prometheus.MustNewConstMetric(c.maxOpenConns, prometheus.GaugeValue, float64(s.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(c.openConns, prometheus.GaugeValue, float64(s.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.inUse, prometheus.GaugeValue, float64(s.InUse))
	ch <- prometheus.MustNewConstMetric(c.idle, prometheus.GaugeValue, float64(s.Idle))
	ch <- prometheus.MustNewConstMetric(c.waitCount, prometheus.CounterValue, float64(s.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.waitDuration, prometheus.CounterValue, s.WaitDuration.Seconds())
	ch <- prometheus.MustNewConstMetric(c.maxIdleClosed, prometheus.CounterValue, float64(s.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(c.maxIdleTimeClosed, prometheus.CounterValue, float64(s.MaxIdleTimeClosed))
	ch <- prometheus.MustNewConstMetric(c.maxLifetimeClosed, prometheus.CounterValue, float64(s.MaxLifetimeClosed))
}
