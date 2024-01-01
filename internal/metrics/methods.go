package metrics

func NewPoint() Point {
	return Point{Tags: map[string]string{}, Fields: map[string]interface{}{}}
}
