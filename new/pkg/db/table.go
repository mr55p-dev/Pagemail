package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Table struct {
	_name   string
	_driver *sql.DB
}

func (t *Table) CreateRecord(record any) error {
	data, err := getStructData(record)
	if err != nil {
		return err
	}

	stmnt := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		t._name,
		strings.Join(data.Fields, ","),
		GetPlaceholder(data.Fields),
	)
	res, err := t._driver.Exec(stmnt, data.Values...)
	if err != nil {
		return err
	}
	res.RowsAffected()
	return nil
}

func (t *Table) CreateRecords(records []any) error {
	recordData := make([]*StructData, len(records))
	for i, v := range records {
		val, err := getStructData(v)
		if err != nil {
			return err
		}

		recordData[i] = val
	}

	values := make([]string, len(records))
	for i, v := range recordData {
		values[i] = fmt.Sprintf("(%s)", GetPlaceholder(v.Fields))
	}
	valStmnt := strings.Join(values, ",")

	stmnt := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s),",
		t._name,
		strings.Join(recordData[0].Fields, ","),
		valStmnt,
	)
	res, err := t._driver.Exec(stmnt)
	if err != nil {
		return err
	}
	res.RowsAffected()
	return nil
}

func (t *Table) ReadRecordByField(field string, val any, out *any) error {
	data, err := getStructData(out)
	if err != nil {
		return err
	}
	stmnt := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = ?",
		strings.Join(data.Fields, ","),
		t._name,
		field,
	)
	row := t._driver.QueryRow(stmnt, val)
	return row.Scan(data.Refs)
}
