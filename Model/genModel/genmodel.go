package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"os/exec"
)

func jsonParse(data []byte) map[string]map[string]map[string]string {
	var temp_1 map[string]interface{}
	if err := json.Unmarshal(data, &temp_1); err != nil {
		panic(err)
	}

	temp_2 := make(map[string]map[string]interface{})
	for k, v := range temp_1 {
		temp_2[k] = v.(map[string]interface{})
	}
	dat := make(map[string]map[string]map[string]string)
	for k, v := range temp_2 {
		dat[k] = make(map[string]map[string]string)
		for k_, v_ := range v {
			temp_3 := v_.(map[string]interface{})
			temp_4 := make(map[string]string)
			for k__, v__ := range temp_3 {
				temp_4[k__] = v__.(string)
			}
			dat[k][k_] = temp_4
		}
	}
	return dat
}

var gensql []string

func genModel(data map[string]map[string]map[string]string) string {
	model := "package Model\n\nimport \"database/sql\"\n\n\nvar tables map[string]*TableType\n\nfunc init(){\ndatabase.connect(`main.sqlite`)\ntables = make(map[string]*TableType)\n\n"

	for table, v := range data {
		genTable := `CREATE TABLE ` + table + `(id INTEGER PRIMARY KEY AUTOINCREMENT,`

		model += "tables[`" + table + "`] = &TableType{}\n"
		model += "tables[`" + table + "`].name = `" + table + "`\n"
		model += "tables[`" + table + "`].v = []CellStruct{{`id`, INT},"

		cells := "[]Cell{{`id`, INT, ``},"
		scanint := 0
		for name, props := range v {
			if name == "id" {
				fmt.Fprintln(os.Stderr, `cell name should not contain 'id'`)
				panic(name)
			}
			if props["type"] == "FILE" {
				genTable += name + ` VARCHAR(255)`
			} else {
				genTable += name + ` ` + props["type"]
			}
			if val, ok := props["default"]; ok {
				genTable += ` DEFAULT ` + val
			}
			model += "{`" + name + "`, " + props["type"] + "},"
			cells += "{`" + name + "`, " + props["type"] + ", ``},"
			genTable += ","
			scanint ++
		}
		gensql = append(gensql, genTable[:len(genTable)-1]+`);`)
		model = model[:len(model)-1] + "}\ntables[`" + table + "`].query = func(rows *sql.Rows) (t Table) {\nvar tabledata [][]Cell\ndefer rows.Close()\nfor rows.Next() {\ntabulardata :="
		model += cells[:len(cells)-1] + "}\nvar scan [" + strconv.Itoa(scanint + 1) + "]interface{}\nerr := rows.Scan("
		for i := 0; i <= scanint;i++{
			model += "&scan[" + strconv.Itoa(i) + "],"
		}
		model = model[:len(model)-1] + ")\nfor i:=0;i<=" + strconv.Itoa(scanint) + ";i++ {\ntabulardata[i].data = Cast(tabulardata[i], scan[i])\n}\ncheckErrorDB(err)\ntabledata = append(tabledata, tabulardata)\n}\nt = tabledata\nreturn t\n}\n"
	}

	model += "\n}\n"

	fmt.Println(gensql)
	return model
}

func main() {
	flag.Parse()
	dest := flag.Arg(0)
	data, err := ioutil.ReadFile(dest)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading model file ", dest)
		panic(err)
	}

	ioutil.WriteFile("model.go", []byte(genModel(jsonParse(data))), os.ModePerm)
	exec.Command("go", "fmt", "model.go").Output()
}
