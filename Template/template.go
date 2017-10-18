package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var genHTML []string // code generation genHTML function body
type constructSort []string

var construct constructSort // code generation construct function body

var varReady map[string]chan string     // used in VarFill
var vars []string                       // list of variables
var varInferring map[string]chan string // middleware for variable type inference
var varType map[string]string           // storage for variable type inference

func parseVariable(variable string) string {
	if strings.Contains(variable, "__") {
		fmt.Fprintln(os.Stderr, "variable name invalid: contains '__'")
		panic(variable)
	}
	for _, i := range []rune(variable) {
		if !strings.ContainsRune(" abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY_.", i) {
			fmt.Fprintln(os.Stderr, "variable name invalid: illegal literal '"+string(i)+"'")
			panic(variable)
		}
	}
	return strings.Replace(variable, ".", "__", -1)
}

// AddVariable identifies all the variables within the {{ }} like {{ a + b * 2 }}, where every operator and variable
// should be separated by a space
func addVariable(variable string) {
	for _, data := range strings.Split(variable, " ") {
		if len(data) > 0 {
			if (data[0] <= 'z' && data[0] >= 'a') || (data[0] <= 'Z' && data[0] >= 'A') {
				vars = append(vars, data)
			}
		}
	}
	genHTML = append(genHTML, "_gen += "+parseVariable(variable))
}

func addOperation(content string) {
	if strings.Contains(content, " for ") {
		vars = append(vars, strings.Split(strings.TrimPrefix(content, " for "), " in ")[0]+"|"+strings.TrimSuffix(strings.Split(content, " in ")[1], " "))
		genHTML = append(genHTML, strings.Replace(strings.Replace(content, " in ", " := range ", 1)[1:]+"{", "for ", "for _, ", -1))
		return
	}
	if strings.Contains(content, " if ") {
		inc := "if "
		for _, data := range strings.Split(strings.Replace(content, " if ", "", 1), " ") {
			if len(data) > 0 {
				if (data[0] <= 'z' && data[0] >= 'a') || (data[0] <= 'Z' && data[0] >= 'A') {
					vars = append(vars, data)
					inc += parseVariable(data) + " "
				} else {
					inc += data + " "
				}
			}
		}
		genHTML = append(genHTML, inc+"{")
		return
	}
	if strings.Contains(content, " endfor ") || strings.Contains(content, " endif ") {
		genHTML = append(genHTML, "}")
		return
	}
	if strings.Contains(content, " else ") {
		genHTML = append(genHTML, "} else {")
		return
	}
}

func walkThrough(content []byte) {
	preIter := 0
	for iter := 0; iter < len(content); iter++ {
		if content[iter] == '{' && iter < len(content)-1 {
			switch content[iter+1] {
			case '{':
				if preIter < iter {
					genHTML = append(genHTML, "_gen += `"+strings.Replace(string(content[preIter:iter]), "\n", `\n`, -1)+"`")
				}
				found := false
				for nextIter := iter; nextIter < len(content)-1; nextIter++ {
					if (content[nextIter] == '}') && (content[nextIter+1] == '}') {
						addVariable(string(content[iter+2 : nextIter]))
						found = true
						preIter = nextIter + 2
						iter = nextIter + 2
						break
					}
				}
				if !found {
					fmt.Fprintln(os.Stderr, `expected }} found EOF`)
					panic(string(content[preIter : iter-preIter+1]))
				}

			case '%':
				if preIter < iter {
					genHTML = append(genHTML, "_gen += `"+strings.Replace(string(content[preIter:iter]), "\n", `\n`, -1)+"`")
				}
				found := false
				for nextIter := iter; nextIter < len(content)-1; nextIter++ {
					if (content[nextIter] == '%') && (content[nextIter+1] == '}') {
						addOperation(string(content[iter+2 : nextIter]))
						found = true
						preIter = nextIter + 2
						iter = nextIter + 2
						break
					}
				}
				if !found {
					fmt.Fprintln(os.Stderr, `expected %} found EOF`)
					panic(string(content[preIter : iter-preIter+1]))
				}

			}
		}
	}
	if preIter < len(content) {
		genHTML = append(genHTML, "_gen += `"+strings.Replace(string(content[preIter:]), "\n", `\n`, -1)+"`")
	}
}

func write(c chan string, i string) {
	c <- i
}

func inferTypes(variable string, r chan string) {
	if strings.ContainsRune(variable, '|') { // for loop
		vars := strings.Split(variable, "|")
		middleWare := <-varInferring[vars[0]]
		go write(varInferring[vars[0]], middleWare)
		if strings.HasPrefix(middleWare, "type") {
			r <- "[]" + parseVariable(vars[0]) + "___class"
		} else {
			r <- "[]" + middleWare
		}
		return
	}
	var instances []string
	for _, allVariable := range vars {
		if strings.Contains(allVariable, variable+".") && !strings.ContainsRune(strings.TrimPrefix(allVariable, variable+"."), '.') {
			instances = append(instances, allVariable)
		}
	}
	if len(instances) == 0 {
		r <- "string"
	} else {
		r <- parseVariable(variable) + "___class"
	}
}

func varInfer() {
	varInferring = make(map[string]chan string)
	varType = make(map[string]string)

	for _, v := range vars {
		varInferring[v] = make(chan string)
		reducedVariable := v
		for strings.ContainsRune(reducedVariable, '.') {
			reducedVariable = v[:strings.LastIndex(reducedVariable, ".")]
			if _, ok := varInferring[reducedVariable]; !ok {
				varInferring[reducedVariable] = make(chan string)
				vars = append(vars, reducedVariable)
			}
		}
	}
	for _, v := range vars {
		go inferTypes(v, varInferring[v])
	}

	for _, v := range vars {
		if strings.ContainsRune(v, '|') {
			varType[strings.Split(v, "|")[1]] = <-varInferring[v]
		} else {
			varType[v] = <-varInferring[v]
		}
		// close(varInferring[v])
	}
}

func singleVarFill(v string, c chan string) {
	switch varType[v] {
	case "string":
		if strings.ContainsRune(v, '.') {
			v_ := parseVariable(v)
			v_parent := v_[:strings.LastIndex(v_, "__")]
			v_name := strings.Split(v, ".")
			c <- v_ + " := " + v_parent + "[`" + v_name[len(v_name)-1] + "`].(string)"
		} else {
			c <- v + " := m[`" + v + "`].(string)"
		}
	case "[]string":
		if strings.ContainsRune(v, '.') {
			v_ := parseVariable(v)
			v_parent := v_[:strings.LastIndex(v_, "__")]
			v_name := strings.Split(v, ".")
			c <- v_ + " := " + v_parent + "[`" + v_name[len(v_name)-1] + "`].([]string)"
		} else {
			c <- v + " := m[`" + v + "`].([]string)"
		}
	default:
		if strings.HasPrefix(varType[v], "[]") {
			if strings.ContainsRune(v, '.') {
				v_ := parseVariable(v)
				v_parent := v_[:strings.LastIndex(v_, "__")]
				v_name := strings.Split(v, ".")
				c <- v_ + " := " + v_parent + "[`" + v_name[len(v_name)-1] + "`].([]map[string]interface{})"
			} else {
				c <- v + " := m[`" + v + "`].([]map[string]interface{})"
			}
		} else {
			if strings.ContainsRune(v, '.') {
				v_ := parseVariable(v)
				v_parent := v_[:strings.LastIndex(v_, "__")]
				v_name := strings.Split(v, ".")
				c <- v_ + " := " + v_parent + "[`" + v_name[len(v_name)-1] + "`].(map[string]interface{})"
			} else {
				c <- v + " := m[`" + v + "`].(map[string]interface{})"
			}
		}
	}
}

func (s constructSort) Len() int {
	return len(s)
}

func (s constructSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s constructSort) Less(i, j int) bool {
	return strings.Count(s[i], "__") < strings.Count(s[j], "__")
}

func varFill() {
	varReady = make(map[string]chan string)
	for v := range varType {
		varReady[v] = make(chan string)
	}
	for v := range varType {
		go singleVarFill(v, varReady[v])
	}

	for v := range varType {
		construct = append(construct, <-varReady[v])
	}

	sort.Sort(construct)
}

func inFor(v string) bool {
	for _, i := range vars {
		if strings.HasPrefix(i, v+"|") {
			return true
		}
	}
	return false
}

func Gen(filename string) {
	// Initialize
	genHTML = []string{}
	construct = []string{}
	vars = []string{}

	content, _ := ioutil.ReadFile(filename + ".html")
	walkThrough(content)
	varInfer()
	varFill()
	gen := ""
	gen += "package template\n"
	gen += "\nfunc genHTML_" + filename[strings.LastIndex(filename, "/")+1:] + "("

	vCall := ""

	for k, v := range varType {
		if !strings.ContainsRune(k, '.') && !inFor(k) {
			if strings.HasPrefix(v, "[]") && v != "[]string" {
				gen += k + " []map[string]interface{},"
			} else {
				if v != "string" {
					gen += k + " map[string]interface{}"
				} else {
					gen += k + " " + v + ","
				}
			}
			vCall += k + ","
		}
	}
	gen = gen[:len(gen)-1] + ")(_gen string){\n_gen = ``\n"
	for _, v := range genHTML {
		gen += v + "\n"
		if strings.HasPrefix(v, "for ") {
			variable := strings.Split(strings.TrimPrefix(v, "for _, "), " :=")[0]
			for _, line := range construct {
				if strings.HasPrefix(line, variable+"__") {
					gen += line + "\n"
				}
			}
		}
	}
	gen += "\nreturn\n}\n\nfunc Construct_" + filename[strings.LastIndex(filename, "/")+1:] + "(m map[string]interface{})(string) {\n"
	for _, v := range construct {
		curVar := strings.Split(v, " :=")[0]
		if !strings.Contains(curVar, "__") && !inFor(curVar) {
			gen += v + "\n"
		}
	}
	gen += "\nreturn genHTML_" + filename[strings.LastIndex(filename, "/")+1:] + "(" + vCall[:len(vCall)-1] + ")\n}"
	ioutil.WriteFile(filename+".go", []byte(gen), os.ModePerm)
	exec.Command("go", "fmt", filename+".go").Output()
}

func main() {
	flag.Parse()

	dir, err := ioutil.ReadDir(flag.Arg(0))
	if err != nil {
		fmt.Println("permission denied", flag.Arg(0))
		os.Exit(-1)
	}
	for _, f := range dir {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".html") {
			Gen(flag.Arg(0) + "/" + strings.TrimSuffix(f.Name(), ".html"))
		}
	}
}
