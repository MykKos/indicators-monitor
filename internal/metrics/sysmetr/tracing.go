package sysmetr

import (
	"runtime"
	"time"

	"indicators-monitor/internal/metrics"
)

// EnableTracing запускает сбор статистики по работе сервиса,
// помогает отследить утечки памяти, незакрытые горутины и тд
// worker — имя entrypoint-а
// appName — имя приложения(напр. имя проекта из репозитория)
func EnableTracing(client metrics.MetricsProvider, appName, worker string) {
	var m runtime.MemStats
	for {
		runtime.ReadMemStats(&m)
		pt := metrics.Point{
			Table: "monitoring",
			Tags: map[string]string{
				"worker": worker,
				"app":    appName,
			},
			Fields: map[string]interface{}{
				"gorutines":   runtime.NumGoroutine(),
				"alloc":       bToKb(m.Alloc),
				"total_alloc": bToKb(m.TotalAlloc),
				"sys":         bToKb(m.TotalAlloc),
				"num_gc":      bToKb(m.TotalAlloc),
			},
		}
		client.HitSave(pt)
		time.Sleep(10 * time.Second)
	}
}

func bToMb(b uint64) int64 {
	return int64(b / 1024 / 1024)
}
func bToKb(b uint64) int64 {
	return int64(b / 1024)
}
