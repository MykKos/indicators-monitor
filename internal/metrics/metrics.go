package metrics

type MetricsProvider interface {
	HitSave(Point)
}
