// A package to preprocess sql database types
package Model

import (
	"database/sql"
	"strconv"
	"time"
)

type cellType int

type CellStruct struct {
	Name string
	T    cellType
}

type Cell struct {
	Name string
	T    cellType
	Data string
}

type TableType struct {
	Name  string
	V     []CellStruct
	Query func(rows *sql.Rows) Table // perform rows.Scan() action
	Database DB
}

type Table [][]Cell

// Cast all cell to string
func Cast(cell Cell, data interface{}) string {
	switch cell.T {
	case INT:
		return strconv.FormatInt(data.(int64), 10)
	case FLOAT:
		return strconv.FormatFloat(data.(float64), 'f', 6, 64)
	case TEXT:
		return string(data.([]uint8))
	case FILE:
		return data.(string)
	case CHAR:
		return string([]byte{byte(data.(rune))})
	case DATETIME:
		return data.(time.Time).String()
	case FOREIGNKEY:
		return strconv.FormatInt(data.(int64), 10)
	default:
		return ""
	}
}

// Selecting data from database
func (m TableType) All() Table {
	queryString := `select * from '` + m.Name + `'`
	return m.Database.query(queryString, m)
}

func (m TableType) Count() int {
	return m.Database.count(m.Name)
}

// Types in database
const (
	CHAR       cellType = iota // varchar
	TEXT                       // text
	INT                        // int
	FLOAT                      // float
	FILE                       // varchar(255)
	DATETIME                   // datetime
	FOREIGNKEY                 // int (id)
)
