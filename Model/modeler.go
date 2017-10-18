package Model

import (
	"database/sql"
	"strconv"
	"time"
	"fmt"
)

type CellType int

type CellStruct struct {
	name string
	T    CellType
}

type Cell struct {
	name string
	T    CellType
	data interface{}
}

type TableType struct {
	name  string
	v     []CellStruct
	query func(rows *sql.Rows) Table // perform rows.Scan() action
}

type Table struct {
	m TableType
	v [][]Cell
}

func Cast(cell Cell) string {
	switch cell.T {
	case INT:
		return strconv.FormatInt(cell.data.(int64), 10)
	case FLOAT:
		return strconv.FormatFloat(cell.data.(float64), 'f', 6, 64)
	case TEXT:
		return cell.data.(string)
	case FILE:
		return cell.data.(string)
	case CHAR:
		return string([]byte{byte(cell.data.(rune))})
	case DATETIME:
		return cell.data.(time.Time).String()
	case FOREIGNKEY:
		return strconv.FormatInt(cell.data.(int64), 10)
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

func Table_test() {
	fmt.Println(Cast(tables["user"].all().v[1][1]))
}
