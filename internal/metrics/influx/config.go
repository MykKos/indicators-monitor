package influx

type InfluxConfig struct {
	Url       string `toml:"url"`
	Database  string `toml:"database"`
	Precision string `toml:"precision"`
	Debug     bool   `toml:"debug"`
}
