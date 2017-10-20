package main

import (
	"database/sql"
	"go-rails/Model"
)

var Tables map[string]*Model.TableType
var database Model.DB

func init() {
	database.Connect(`main.sqlite`)
	Tables = make(map[string]*Model.TableType)

	Tables[`user`] = &Model.TableType{}
	Tables[`user`].Name = `user`
	Tables[`user`].Database = database
	Tables[`user`].V = []Model.CellStruct{{`id`, Model.INT}, {`username`, Model.TEXT}}
	Tables[`user`].Query = func(rows *sql.Rows) (t Model.Table) {
		var tabledata [][]Model.Cell
		defer rows.Close()
		for rows.Next() {
			tabulardata := []Model.Cell{{`id`, Model.INT, ``}, {`username`, Model.TEXT, ``}}
			var scan [2]interface{}
			err := rows.Scan(&scan[0], &scan[1])
			for i := 0; i <= 1; i++ {
				tabulardata[i].Data = Model.Cast(tabulardata[i], scan[i])
			}
			Model.CheckErrorDB(err)
			tabledata = append(tabledata, tabulardata)
		}
		t = tabledata
		return t
	}
	Tables[`test`] = &Model.TableType{}
	Tables[`test`].Name = `test`
	Tables[`test`].Database = database
	Tables[`test`].V = []Model.CellStruct{{`id`, Model.INT}, {`time`, Model.DATETIME}, {`name`, Model.INT}}
	Tables[`test`].Query = func(rows *sql.Rows) (t Model.Table) {
		var tabledata [][]Model.Cell
		defer rows.Close()
		for rows.Next() {
			tabulardata := []Model.Cell{{`id`, Model.INT, ``}, {`time`, Model.DATETIME, ``}, {`name`, Model.INT, ``}}
			var scan [3]interface{}
			err := rows.Scan(&scan[0], &scan[1], &scan[2])
			for i := 0; i <= 2; i++ {
				tabulardata[i].Data = Model.Cast(tabulardata[i], scan[i])
			}
			Model.CheckErrorDB(err)
			tabledata = append(tabledata, tabulardata)
		}
		t = tabledata
		return t
	}

}
