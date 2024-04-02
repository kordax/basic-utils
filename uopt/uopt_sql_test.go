////go:build integration_test

/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uopt_test

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kordax/basic-utils/uopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestOptTable(db *sql.DB) error {
	// Define the SQL statement to create the table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS test_opt (
		id INT AUTO_INCREMENT PRIMARY KEY,
		int_col BIGINT,
		float_col DOUBLE,
		string_col VARCHAR(255),
		bool_col BOOLEAN,
		date_col DATETIME
	);`

	// Execute the SQL statement to create the table
	_, err := db.Exec(createTableSQL)
	return err
}

func clearTestOptTable(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM test_opt")
	return err
}

func TestOpt_InsertValues_Integration(t *testing.T) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db?parseTime=true")
	require.NoError(t, err)
	defer db.Close()

	err = setupTestOptTable(db)
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err = clearTestOptTable(db)
		require.NoError(t, err)
	}(db)

	stmt, err := db.Prepare("INSERT INTO test_opt (int_col, float_col, string_col, bool_col, date_col) VALUES (?, ?, ?, ?, ?)")
	require.NoError(t, err)

	intVal := 42
	floatVal := 3.14
	stringVal := "hello"
	boolVal := true
	dateVal := time.Now()

	_, err = stmt.Exec(
		uopt.Of(intVal),
		uopt.Of(floatVal),
		uopt.Of(stringVal),
		uopt.Of(boolVal),
		uopt.Of(dateVal),
	)
	require.NoError(t, err)

	var readInt int
	err = db.QueryRow("SELECT int_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readInt)
	require.NoError(t, err)
	assert.Equal(t, intVal, readInt)

	var readFloat float64
	err = db.QueryRow("SELECT float_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readFloat)
	require.NoError(t, err)
	assert.Equal(t, floatVal, readFloat)

	var readStr string
	err = db.QueryRow("SELECT string_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readStr)
	require.NoError(t, err)
	assert.Equal(t, stringVal, readStr)

	var readBool bool
	err = db.QueryRow("SELECT bool_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readBool)
	require.NoError(t, err)
	assert.Equal(t, boolVal, readBool)

	var readDate time.Time
	err = db.QueryRow("SELECT date_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readDate)
	require.NoError(t, err)
	assert.WithinDuration(t, dateVal, readDate, time.Second)
}

func TestOpt_ReadValues_Integration(t *testing.T) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db?parseTime=true")
	require.NoError(t, err)
	defer db.Close()

	err = setupTestOptTable(db)
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err = clearTestOptTable(db)
		require.NoError(t, err)
	}(db)

	stmt, err := db.Prepare("INSERT INTO test_opt (int_col, float_col, string_col, bool_col, date_col) VALUES (?, ?, ?, ?, ?)")
	require.NoError(t, err)

	intVal := 42
	floatVal := 3.14
	stringVal := "hello"
	boolVal := true
	dateVal := time.Now()

	_, err = stmt.Exec(
		intVal,
		floatVal,
		stringVal,
		boolVal,
		dateVal,
	)
	require.NoError(t, err)

	var readInt uopt.Opt[int]
	err = db.QueryRow("SELECT int_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readInt)
	require.NoError(t, err)
	assert.Equal(t, intVal, *readInt.Get())

	var readFloat uopt.Opt[float64]
	err = db.QueryRow("SELECT float_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readFloat)
	require.NoError(t, err)
	assert.Equal(t, floatVal, *readFloat.Get())

	var readStr uopt.Opt[string]
	err = db.QueryRow("SELECT string_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readStr)
	require.NoError(t, err)
	assert.Equal(t, stringVal, *readStr.Get())

	var readBool uopt.Opt[bool]
	err = db.QueryRow("SELECT bool_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readBool)
	require.NoError(t, err)
	assert.Equal(t, boolVal, *readBool.Get())

	var readDate uopt.Opt[time.Time]
	err = db.QueryRow("SELECT date_col FROM test_opt ORDER BY id DESC LIMIT 1").Scan(&readDate)
	require.NoError(t, err)
	assert.True(t, readDate.Present())
	assert.WithinDuration(t, dateVal, *readDate.Get(), time.Second)
}
