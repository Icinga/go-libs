// go-libs | (c) 2020 Icinga GmbH | GPLv2+

package sql

import (
	"database/sql"
	"reflect"
)

// FetchRowsAsStructSlice fetches rows - at most limit (unless -1) - as slice of rowType's type.
// Consult the unit tests for usage examples.
func FetchRowsAsStructSlice(rows *sql.Rows, rowType interface{}, limit int) (interface{}, error) {
	blankRow := reflect.ValueOf(rowType)
	res := reflect.MakeSlice(reflect.SliceOf(blankRow.Type()), 0, 0)
	idx := -1
	scanDest := make([]interface{}, blankRow.NumField())

	for {
		if limit > -1 {
			if limit < 1 {
				break
			}

			limit--
		}

		if rows.Next() {
			res = reflect.Append(res, blankRow)
			idx++

			row := res.Index(idx)

			for i := range scanDest {
				scanDest[i] = row.Field(i).Addr().Interface()
			}

			if errSc := rows.Scan(scanDest...); errSc != nil {
				return nil, errSc
			}
		} else if errNx := rows.Err(); errNx == nil {
			break
		} else {
			return nil, errNx
		}
	}

	return res.Interface(), nil
}
