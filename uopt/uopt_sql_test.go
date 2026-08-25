//go:build integration_test

/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uopt_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v4/uopt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type integrationJSONPayload struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

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

func setupTestOptJSONTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS test_opt_json (
		id INT AUTO_INCREMENT PRIMARY KEY,
		json_col JSON,
		json_struct_col JSON,
		decimal_col DECIMAL(10, 2),
		binary_col VARBINARY(255)
	);`

	_, err := db.Exec(createTableSQL)
	return err
}

func clearTestOptJSONTable(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM test_opt_json")
	return err
}

func setupMySQLDB(t *testing.T) *sql.DB {
	t.Helper()

	ctx := context.Background()
	container, err := tcmysql.Run(
		ctx,
		"mysql:8.4",
		tcmysql.WithDatabase("db"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("password"),
	)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(ctx, "parseTime=true")
	require.NoError(t, err)

	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	require.NoError(t, db.PingContext(ctx))

	return db
}

func setupPostgresDB(t *testing.T) *sql.DB {
	t.Helper()

	ctx := context.Background()
	container, err := tcpostgres.Run(
		ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("db"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("password"),
		tcpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	require.NoError(t, db.PingContext(ctx))

	return db
}

func TestOpt_InsertValues_Integration(t *testing.T) {
	db := setupMySQLDB(t)

	err := setupTestOptTable(db)
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
	db := setupMySQLDB(t)

	err := setupTestOptTable(db)
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

func TestOpt_JSONAndSQLTypes_Integration(t *testing.T) {
	db := setupMySQLDB(t)

	err := setupTestOptJSONTable(db)
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err = clearTestOptJSONTable(db)
		require.NoError(t, err)
	}(db)

	payload := integrationJSONPayload{Name: "test", Items: []string{"a", "b"}}
	raw := json.RawMessage(`{"name":"raw","items":["x"]}`)
	binaryValue := []byte{1, 2, 3}

	_, err = db.Exec(
		"INSERT INTO test_opt_json (json_col, json_struct_col, decimal_col, binary_col) VALUES (?, ?, ?, ?)",
		uopt.Of(map[string]any{"name": "map", "count": 2}),
		uopt.Of(payload),
		uopt.Of(123.45),
		uopt.Of(binaryValue),
	)
	require.NoError(t, err)

	var readMap uopt.Opt[map[string]any]
	err = db.QueryRow("SELECT json_col FROM test_opt_json ORDER BY id DESC LIMIT 1").Scan(&readMap)
	require.NoError(t, err)
	require.True(t, readMap.Present())
	assert.Equal(t, "map", (*readMap.Get())["name"])
	assert.EqualValues(t, 2, (*readMap.Get())["count"])

	var readStruct uopt.Opt[integrationJSONPayload]
	err = db.QueryRow("SELECT json_struct_col FROM test_opt_json ORDER BY id DESC LIMIT 1").Scan(&readStruct)
	require.NoError(t, err)
	require.True(t, readStruct.Present())
	assert.Equal(t, payload, *readStruct.Get())

	var readDecimal uopt.Opt[float64]
	err = db.QueryRow("SELECT decimal_col FROM test_opt_json ORDER BY id DESC LIMIT 1").Scan(&readDecimal)
	require.NoError(t, err)
	require.True(t, readDecimal.Present())
	assert.InDelta(t, 123.45, *readDecimal.Get(), 0.001)

	var readBinary uopt.Opt[[]byte]
	err = db.QueryRow("SELECT binary_col FROM test_opt_json ORDER BY id DESC LIMIT 1").Scan(&readBinary)
	require.NoError(t, err)
	require.True(t, readBinary.Present())
	assert.Equal(t, binaryValue, *readBinary.Get())

	_, err = db.Exec(
		"INSERT INTO test_opt_json (json_col, json_struct_col, decimal_col, binary_col) VALUES (?, ?, ?, ?)",
		uopt.Of(raw),
		uopt.Null[integrationJSONPayload](),
		uopt.Null[float64](),
		uopt.Null[[]byte](),
	)
	require.NoError(t, err)

	var readRaw uopt.Opt[json.RawMessage]
	err = db.QueryRow("SELECT json_col FROM test_opt_json ORDER BY id DESC LIMIT 1").Scan(&readRaw)
	require.NoError(t, err)
	require.True(t, readRaw.Present())
	assert.JSONEq(t, string(raw), string(*readRaw.Get()))

	var nullStruct uopt.Opt[integrationJSONPayload]
	err = db.QueryRow("SELECT json_struct_col FROM test_opt_json ORDER BY id DESC LIMIT 1").Scan(&nullStruct)
	require.NoError(t, err)
	assert.False(t, nullStruct.Present())
}

func TestOpt_PostgreSQLTypes_Integration(t *testing.T) {
	db := setupPostgresDB(t)

	_, err := db.Exec(`
		CREATE TABLE test_opt_postgres (
			id SERIAL PRIMARY KEY,
			text_array_col TEXT[],
			int_array_col INTEGER[],
			text_array_text_col TEXT,
			jsonb_col JSONB,
			jsonb_struct_col JSONB,
			numeric_col NUMERIC(10, 2),
			bytea_col BYTEA,
			bool_col BOOLEAN,
			timestamptz_col TIMESTAMPTZ
		);
	`)
	require.NoError(t, err)

	payload := integrationJSONPayload{Name: "postgres", Items: []string{"jsonb", "array"}}
	binaryValue := []byte{4, 5, 6}
	now := time.Now().UTC().Truncate(time.Microsecond)

	_, err = db.Exec(
		`INSERT INTO test_opt_postgres (
			text_array_col,
			int_array_col,
			text_array_text_col,
			jsonb_col,
			jsonb_struct_col,
			numeric_col,
			bytea_col,
			bool_col,
			timestamptz_col
		) VALUES (
			ARRAY['alpha', 'beta,quoted', 'slash\value']::TEXT[],
			ARRAY[1, 2, 3]::INTEGER[],
			ARRAY['plain', 'text']::TEXT,
			$1::JSONB,
			$2::JSONB,
			$3::NUMERIC,
			$4::BYTEA,
			$5::BOOLEAN,
			$6::TIMESTAMPTZ
		)`,
		uopt.Of(map[string]any{"name": "map", "count": 3}),
		uopt.Of(payload),
		uopt.Of(456.78),
		uopt.Of(binaryValue),
		uopt.Of(true),
		uopt.Of(now),
	)
	require.NoError(t, err)

	var textArray uopt.Opt[[]string]
	err = db.QueryRow("SELECT text_array_col::TEXT FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&textArray)
	require.NoError(t, err)
	require.True(t, textArray.Present())
	assert.Equal(t, []string{"alpha", "beta,quoted", `slash\value`}, *textArray.Get())

	var intArray uopt.Opt[[]int]
	err = db.QueryRow("SELECT int_array_col::TEXT FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&intArray)
	require.NoError(t, err)
	require.True(t, intArray.Present())
	assert.Equal(t, []int{1, 2, 3}, *intArray.Get())

	var arrayText uopt.Opt[[]string]
	err = db.QueryRow("SELECT text_array_text_col FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&arrayText)
	require.NoError(t, err)
	require.True(t, arrayText.Present())
	assert.Equal(t, []string{"plain", "text"}, *arrayText.Get())

	var readMap uopt.Opt[map[string]any]
	err = db.QueryRow("SELECT jsonb_col FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&readMap)
	require.NoError(t, err)
	require.True(t, readMap.Present())
	assert.Equal(t, "map", (*readMap.Get())["name"])
	assert.EqualValues(t, 3, (*readMap.Get())["count"])

	var readStruct uopt.Opt[integrationJSONPayload]
	err = db.QueryRow("SELECT jsonb_struct_col FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&readStruct)
	require.NoError(t, err)
	require.True(t, readStruct.Present())
	assert.Equal(t, payload, *readStruct.Get())

	var rawJSON uopt.Opt[json.RawMessage]
	err = db.QueryRow("SELECT jsonb_struct_col::TEXT FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&rawJSON)
	require.NoError(t, err)
	require.True(t, rawJSON.Present())
	assert.JSONEq(t, `{"name":"postgres","items":["jsonb","array"]}`, string(*rawJSON.Get()))

	var numeric uopt.Opt[float64]
	err = db.QueryRow("SELECT numeric_col FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&numeric)
	require.NoError(t, err)
	require.True(t, numeric.Present())
	assert.InDelta(t, 456.78, *numeric.Get(), 0.001)

	var bytes uopt.Opt[[]byte]
	err = db.QueryRow("SELECT bytea_col FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&bytes)
	require.NoError(t, err)
	require.True(t, bytes.Present())
	assert.Equal(t, binaryValue, *bytes.Get())

	var boolean uopt.Opt[bool]
	err = db.QueryRow("SELECT bool_col FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&boolean)
	require.NoError(t, err)
	require.True(t, boolean.Present())
	assert.True(t, *boolean.Get())

	var timestamp uopt.Opt[time.Time]
	err = db.QueryRow("SELECT timestamptz_col FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&timestamp)
	require.NoError(t, err)
	require.True(t, timestamp.Present())
	assert.WithinDuration(t, now, *timestamp.Get(), time.Second)

	_, err = db.Exec(`
		INSERT INTO test_opt_postgres (text_array_col, jsonb_col, numeric_col, bytea_col)
		VALUES (NULL, NULL, NULL, NULL)
	`)
	require.NoError(t, err)

	var nullArray uopt.Opt[[]string]
	err = db.QueryRow("SELECT text_array_col::TEXT FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&nullArray)
	require.NoError(t, err)
	assert.False(t, nullArray.Present())

	var nullJSON uopt.Opt[map[string]any]
	err = db.QueryRow("SELECT jsonb_col FROM test_opt_postgres ORDER BY id DESC LIMIT 1").Scan(&nullJSON)
	require.NoError(t, err)
	assert.False(t, nullJSON.Present())
}
