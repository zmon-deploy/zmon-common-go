package database

import (
	"fmt"
	"github.com/zmon-deploy/zmon-common-go/stringutil"
	"github.com/influxdata/influxdb1-client/models"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"
	"time"
)

type IterateQueryFunc func(int, int, models.Row) error

type InfluxdbAdapter struct {
	influxdb.Client
}

func NewInfluxdbAdapter(config influxdb.HTTPConfig) (*InfluxdbAdapter, error) {
	client, err := influxdb.NewHTTPClient(config)
	if err != nil {
		return nil, err
	}
	if _, _, err := client.Ping(10 * time.Second); err != nil {
		return nil, errors.Wrap(err, "failed to ping influxdb")
	}

	return &InfluxdbAdapter{client}, nil
}

func (a *InfluxdbAdapter) QueryAndIterate(query, database string, fn IterateQueryFunc) error {
	resp, err := a.Query(influxdb.NewQuery(query, database, ""))
	if err != nil {
		return errors.Wrapf(err, "failed to query influxdb, database: %s, query: %s", database, query)
	}
	if resp.Error() != nil {
		return errors.Wrapf(resp.Error(), "failed to get response of influxdb query, database: %s, query: %s", database, query)
	}

	for i, result := range resp.Results {
		for ii, row := range result.Series {
			if err := fn(i, ii, row); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *InfluxdbAdapter) ExtractTagFromQuery(query, database, tagKey string) ([]string, error) {
	var tagValues []string

	err := a.QueryAndIterate(query, database, func(_, _ int, row models.Row) error {
		if val, ok := row.Tags[tagKey]; ok {
			tagValues = append(tagValues, val)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query")
	}

	return tagValues, nil
}

func (a *InfluxdbAdapter) FindMeasurements(database string) ([]string, error) {
	return a.findMeasurements(database, "show measurements")
}

func (a *InfluxdbAdapter) FindMeasurementsLike(database string, regex string) ([]string, error) {
	return a.findMeasurements(database, fmt.Sprintf("show measurements with measurement =~ /%s/", regex))
}

func (a *InfluxdbAdapter) findMeasurements(database string, query string) ([]string, error) {
	measurements := stringutil.NewStringSet()

	err := a.QueryAndIterate(query, database, func(_, _ int, row models.Row) error {
		for _, value := range row.Values {
			measurements.Add(value[0].(string))
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query result")
	}

	return measurements.Values(), nil
}

func (a *InfluxdbAdapter) WritePoints(database string, points ...*influxdb.Point) error {
	batchPoints, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{Database: database})
	if err != nil {
		return errors.Wrap(err, "failed to create a influxdb batch points")
	}

	for _, point := range points {
		batchPoints.AddPoint(point)
	}

	if err := a.Write(batchPoints); err != nil {
		return errors.Wrap(err, "failed to write batch into influxdb")
	}
	return nil
}
