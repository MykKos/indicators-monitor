package metrics

type Point struct {
	Table  string
	Tags   map[string]string
	Fields map[string]interface{}
}
