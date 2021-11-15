package main

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type Record struct {
	ConsumerId   int
	ConsumedData int
}

func NewRecordFromString(s string) (rec *Record, err error) {
	if s == "" {
		return nil, errors.New("cannot parse an empty string")
	}
	splitted := strings.Split(s, " ")
	if splitted[0] == "" || splitted[1] == "" {
		return nil, errors.New("string is in unsupported format")
	}
	rec = &Record{}
	rec.ConsumerId, err = strconv.Atoi(splitted[0])
	if err != nil {
		return nil, err
	}

	rec.ConsumedData, err = strconv.Atoi(splitted[1])
	if err != nil {
		return nil, err
	}

	return
}

func SaveUser(db *sql.DB, stats *Record) error {
	_, err := db.Exec(
		SaveStatsQuery,
		stats.ConsumerId,
		stats.ConsumedData,
	)
	return err
}

const SaveStatsQuery = `INSERT INTO statistics(consumer_id, consumed_data) VALUES ($1, $2)
ON CONFLICT (consumer_id) DO UPDATE
SET consumed_data = EXCLUDED.consumed_data + statistics.consumed_data`
