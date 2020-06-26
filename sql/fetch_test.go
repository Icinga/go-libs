package sql

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchRowsAsStructSlice(t *testing.T) {
	type row struct {
		I int8
		R float32
		T string
		B []byte
	}

	db, errOp := sql.Open("sqlite3", "file::memory:?cache=shared")
	if errOp != nil {
		t.Fatal(errOp)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, errEx := db.Exec("CREATE TABLE test (i INT, r REAL, t TEXT, b BLOB)"); errEx != nil {
		t.Fatal(errEx)
	}

	{
		_, errEx := db.Exec(
			"INSERT INTO test(i, r, t, b) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)",
			1, 2.5, "3", []byte{'4'},
			-5, -6.25, "-7", []byte("-8"),
			9, 10.125, "11", []byte("12"),
		)
		if errEx != nil {
			t.Fatal(errEx)
		}
	}

	testACase := func(limit int, out []row) {
		t.Helper()

		rows, errQr := db.Query("SELECT i, r, t, b FROM test ORDER BY i")
		if errQr != nil {
			t.Fatal(errQr)
		}

		defer rows.Close()

		res, errFR := FetchRowsAsStructSlice(rows, row{}, limit)

		assert.Nil(t, errFR)
		assert.Equal(t, out, res)
	}

	all := []row{
		{-5, -6.25, "-7", []byte("-8")},
		{1, 2.5, "3", []byte{'4'}},
		{9, 10.125, "11", []byte("12")},
	}

	testACase(0, []row{})
	testACase(1, all[:1])
	testACase(2, all[:2])
	testACase(3, all)
	testACase(4, all)
	testACase(-1, all)
}
