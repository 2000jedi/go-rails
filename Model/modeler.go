package Model

// Example: table["user"].all(); table["user"].count()

import (
	"database/sql"
	"strconv"
	"time"
)

type CellType int

type CellStruct struct {
	name string
	T    CellType
}

type Cell struct {
	name string
	T    CellType
	data string
}

type TableType struct {
	name  string
	v     []CellStruct
	query func(rows *sql.Rows) Table // perform rows.Scan() action
}

type Table [][]Cell

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

func (m TableType) all() Table {
	queryString := `select * from '` + m.name + `'`
	return database.query(queryString, m)
}

func (m TableType) count() int {
	return database.count(m.name)
}

const (
	CHAR       CellType = iota // varchar
	TEXT                       // text
	INT                        // int
	FLOAT                      // float
	FILE                       // varchar(255)
	DATETIME                   // datetime
	FOREIGNKEY                 // int (id)
)

var database DB
