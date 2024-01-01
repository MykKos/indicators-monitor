package influx

import (
	"fmt"
	"time"

	"indicators-monitor/internal/metrics"

	influx "github.com/influxdata/influxdb/client/v2"
)

type Client struct {
	InfluxClient influx.Client
	influxPoints influx.BatchPoints
	database     string
	precision    string
	debugMode    bool
	points       Queue
}

func NewFromConfig(config InfluxConfig) (*Client, error) {
	influxClient := &Client{
		database:  config.Database,
		precision: config.Precision,
		debugMode: config.Debug,
	}
	ic, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: config.Url,
	})
	if err != nil {
		return nil, fmt.Errorf("influx was not inited: %s", err)
	}
	influxClient.InfluxClient = ic

	go func() {
		influxClient.SetPoints()
		for {
			influxClient.ProcessPoints()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return influxClient, nil
}

func (client *Client) HitSave(point metrics.Point) {
	pt, err := influx.NewPoint(point.Table, point.Tags, point.Fields, time.Now())
	if err != nil {
		fmt.Printf("точки в influx не были созданы: %s\n", err)
		return
	}
	client.points.Push(pt)
}

func (client *Client) Close() {
	client.InfluxClient.Close()
}

func (client *Client) SetPoints() {
	influxPoints, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  client.database,
		Precision: client.precision,
	})
	if err != nil {
		fmt.Printf("batch points не были созданы: %s\n", err)
		panic(err)
	}
	client.influxPoints = influxPoints
}

func (client *Client) ProcessPoints() {
	pointsCounter := len(client.influxPoints.Points())
	// можно попробовать заменить на client.points.PopAll()
	for ; pointsCounter < 100; pointsCounter++ {
		pointInt := client.points.Pop()
		if pointInt == nil {
			break
		}
		point := pointInt.(*influx.Point)
		client.influxPoints.AddPoint(point)
	}
	if pointsCounter != 0 {
		client.Write()
	}
}

func (client *Client) Write() {
	err := client.InfluxClient.Write(client.influxPoints)
	if err != nil {
		fmt.Printf("не удалось сохранить точки в influx: %s\n", err)
		return
	}

	client.SetPoints()
}
