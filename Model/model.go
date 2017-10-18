package Model

import "database/sql"


var tables map[string]*TableType

func init(){
	database.connect(`main.sqlite`)
	tables = make(map[string]*TableType)

	tables[`user`] = &TableType{}
	tables[`user`].name = `user`
	tables[`user`].v = []CellStruct{{`id`, INT},{`username`, TEXT}}
	tables[`user`].query = func(rows *sql.Rows) (t Table) {
		t.m = *tables[`user`]
		var tabledata [][]Cell
		defer rows.Close()
		for rows.Next() {
			tabulardata :=[]Cell{{`id`, INT, nil},{`username`, TEXT, nil}}
			err := rows.Scan(&tabulardata[0].data,&tabulardata[1].data)
			checkErrorDB(err)
			tabledata = append(tabledata, tabulardata)
		}
		t.v = tabledata
		return t
	}
	tables[`test`] = &TableType{}
	tables[`test`].name = `test`
	tables[`test`].v = []CellStruct{{`id`, INT},{`name`, INT},{`time`, DATETIME}}
	tables[`test`].query = func(rows *sql.Rows) (t Table) {
		t.m = *tables[`test`]
		var tabledata [][]Cell
		defer rows.Close()
		for rows.Next() {
			tabulardata :=[]Cell{{`id`, INT, nil},{`name`, INT, nil},{`time`, DATETIME, nil}}
			err := rows.Scan(&tabulardata[0].data,&tabulardata[1].data,&tabulardata[2].data)
			checkErrorDB(err)
			tabledata = append(tabledata, tabulardata)
		}
		t.v = tabledata
		return t
	}

}