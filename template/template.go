package template

import (
	"io/ioutil"
	"fmt"
	"os"
	"strings"
)

var format []string
var vars []string
var varInferring map[string]chan string
var varInferred map[string]string
var typedefs []string

func appendString(content string) {
	format = append(format, "_gen += \""+strings.Replace(content, "\n", `\n`, -1)+"\"")
}

func addVariable(variable string) {
	for _, data := range strings.Split(variable, " ") {
		if len(data) > 0 {
			if (data[0] <= 'z' && data[0] >= 'a') || (data[0] <= 'Z' && data[0] >= 'A') {
				vars = append(vars, data)
			}
		}
	}
	format = append(format, "_gen += "+variable)
}

func addOperation(content string) {
	if strings.Contains(content, " for ") {
		vars = append(vars, strings.Split(strings.TrimPrefix(content, " for "), " in ")[0]+"|"+strings.TrimSuffix(strings.Split(content, " in ")[1], " "))
		format = append(format, strings.Replace(content, " in ", " := ", 1)[1:]+"{")
		return
	}
	if strings.Contains(content, " if ") {
		for _, data := range strings.Split(strings.Replace(content, " if ", "", 1), " ") {
			if len(data) > 0 {
				if (data[0] <= 'z' && data[0] >= 'a') || (data[0] <= 'Z' && data[0] >= 'A') {
					vars = append(vars, data)
				}
			}
		}
		format = append(format, content[1:]+"{")
		return
	}
	if strings.Contains(content, " endfor ") || strings.Contains(content, " endif ") {
		format = append(format, "}")
		return
	}
}

func walkThrough(content []byte) {
	preIter := 0
	for iter := 0; iter < len(content); iter++ {
		if content[iter] == '{' && iter < len(content)-1 {
			switch content[iter+1] {
			case '{':
				//fmt.Println(preIter, iter)
				if preIter < iter {
					appendString(string(content[preIter:iter]))
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
					fmt.Printf(`expected }} found EOF`)
					os.Exit(-1)
				}

			case '%':
				//fmt.Println(preIter, iter)
				if preIter < iter {
					appendString(string(content[preIter:iter]))
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
					fmt.Printf(`expected %} found EOF`)
					os.Exit(-1)
				}

			}
		}
	}
	if preIter < len(content) {
		appendString(string(content[preIter:]))
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
		if strings.HasPrefix(middleWare, "type"){
			r <- "[]" + strings.Replace(vars[0], ".", "_", -1) + "__class"
		} else {
			r <- "[]" + middleWare
		}
		return
	}
	var instances []string
	for _, allVariable := range vars {
		if strings.Contains(allVariable, variable + ".") && !strings.ContainsRune(strings.TrimPrefix(allVariable, variable + "."), '.'){
			instances = append(instances, allVariable)
		}
	}
	if len(instances) == 0 {
		r <- "string"
	} else {
		typedef := "type " + strings.Replace(variable, ".", "_", -1) + "__class struct { \n"
		for _, subVar := range instances {
			middleWare := <-varInferring[subVar]
			go write(varInferring[subVar], middleWare)
			typedef += "    " + strings.TrimPrefix(subVar, variable + ".") + " " + middleWare + "\n"
		}
		typedefs = append(typedefs, typedef + "}")
		r <- strings.Replace(variable, ".", "_", -1) + "__class"
	}
}

func Gen(filename string) {
	content, _ := ioutil.ReadFile(filename)
	walkThrough(content)

	//for _, v := range format {
	//	fmt.Println(v)
	//}

	varInferring = make(map[string]chan string)
	varInferred = make(map[string]string)

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
		if !strings.ContainsRune(v, '.') {
			varInferred[v] = <-varInferring[v]
		}
	}

	for k, v := range varInferred {
		fmt.Println(k, v)
	}
	for _, v := range typedefs {
		fmt.Println(v)
	}
}
