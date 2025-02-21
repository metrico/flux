package table

import (
	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/internal/execute/groupkey"
)

// Sort will read a TableIterator and produce another TableIterator
// where the keys are sorted.
//
// This method will buffer all of the data since it needs to ensure
// all of the tables are read to avoid any deadlocks. Be careful
// using this method in performance sensitive areas.
func Sort(tables flux.TableIterator) (flux.TableIterator, error) {
	groups := groupkey.NewLookup()
	if err := tables.Do(func(table flux.Table) error {
		buffered, err := Copy(table)
		if err != nil {
			return err
		}
		groups.Set(buffered.Key(), buffered)
		return nil
	}); err != nil {
		return nil, err
	}

	var buffered []flux.Table
	groups.Range(func(_ flux.GroupKey, value interface{}) error {
		buffered = append(buffered, value.(flux.Table))
		return nil
	})
	return Iterator(buffered), nil
}
