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
		var tabledata [][]Cell
		defer rows.Close()
		for rows.Next() {
			tabulardata :=[]Cell{{`id`, INT, ``},{`username`, TEXT, ``}}
			var scan [2]interface{}
			err := rows.Scan(&scan[0],&scan[1])
			for i:=0;i<=1;i++ {
				tabulardata[i].data = Cast(tabulardata[i], scan[i])
			}
			checkErrorDB(err)
			tabledata = append(tabledata, tabulardata)
		}
		t = tabledata
		return t
	}
	tables[`test`] = &TableType{}
	tables[`test`].name = `test`
	tables[`test`].v = []CellStruct{{`id`, INT},{`time`, DATETIME},{`name`, INT}}
	tables[`test`].query = func(rows *sql.Rows) (t Table) {
		var tabledata [][]Cell
		defer rows.Close()
		for rows.Next() {
			tabulardata :=[]Cell{{`id`, INT, ``},{`time`, DATETIME, ``},{`name`, INT, ``}}
			var scan [3]interface{}
			err := rows.Scan(&scan[0],&scan[1],&scan[2])
			for i:=0;i<=2;i++ {
				tabulardata[i].data = Cast(tabulardata[i], scan[i])
			}
			checkErrorDB(err)
			tabledata = append(tabledata, tabulardata)
		}
		t = tabledata
		return t
	}

}