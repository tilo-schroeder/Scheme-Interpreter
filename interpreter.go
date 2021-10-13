//https://gist.github.com/pkelchte/c2bd76b9f8f9cd603b3c
//TODO: Add possibility to initialize lists explicitly via '(1 2 3 4)

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

//type Number float64
type Symbol string
type Variable string
type expr interface{}
type Bool bool

//A variable is just an alias for an expression
type vars map[Variable]expr

//A procedure is executed within an environment and has parameters as well as a body expression
type proc struct {
	params expr
	body   expr
	en     *env
}

type env struct {
	vars
	outer *env
}

func tokenize(chars string) []string {
	//tokens := strings.Split(strings.Replace(strings.Replace(chars, "(", " ( ", -1), ")", " ) ", -1), " ")
	tokens := strings.Fields(strings.Replace(strings.Replace(chars, "(", " ( ", -1), ")", " ) ", -1))
	return tokens //Exclude leading and last whitespace
}

func (e *env) Find(s Variable) *env {
	if _, ok := e.vars[s]; ok {
		return e
	} else {
		//If the symbol is not in the env, look at the outer env
		return e.outer.Find(s)
	}
}

func readFromTokens(tokens *[]string) (expression expr) {
	token := (*tokens)[0]
	*tokens = (*tokens)[1:]
	switch token {
	case "(":
		L := make([]expr, 0)
		for (*tokens)[0] != ")" {
			if i := readFromTokens(tokens); i != Variable("") {
				L = append(L, i)
			}
		}
		*tokens = (*tokens)[1:]
		return L
	default:
		f, errF := strconv.ParseFloat(token, 64)
		if errF == nil {
			return f
		} else if token == "#true" {
			return Bool(true)
		} else if token == "#false" {
			return Bool(false)
		} else if token[0] == 39 { //39 is ' in ASCII
			return Symbol(token[1:])
		} else {
			return Variable(token)
		}
	}

}

func interpret(line string) (expression expr) {
	tokens := tokenize(line)
	return readFromTokens(&tokens)
}

func findParentheses(tokens []string, start int) []int {
	stack := []int{}
	i := start
	partBeginning := 0
	partEnding := 0
	for ok := true; ok; ok = len(stack) != 0 {
		if tokens[i] == "(" {
			stack = append(stack, i)
		} else if tokens[i] == ")" {
			openingPos := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			partBeginning = openingPos
			partEnding = i
		}
		i++
	}
	return []int{partBeginning, partEnding}
}

func findParts(tokens []string) [][]int {
	parts := [][]int{}
	currentPart := []int{}
	pos := 0
	for pos < len(tokens) {
		currentPart = findParentheses(tokens, pos)
		parts = append(parts, currentPart)
		pos = currentPart[1] + 1
	}
	return parts
}

func apply(procedure expr, args []expr) (value expr) {
	switch p := procedure.(type) {
	case func(...expr) expr:
		value = p(args...)
	case proc:
		en := &env{make(vars), p.en}
		switch params := p.params.(type) {
		case []expr:
			for i, param := range params {
				en.vars[param.(Variable)] = args[i]
			}
		default:
			en.vars[params.(Variable)] = args
		}
		value = eval(p.body, en)
	default:
		log.Println("Failed application - unknown procedure", p)
	}
	return
}

func eval(expression expr, en *env) (value expr) {
	switch e := expression.(type) {
	case int64:
		value = e
	case float64:
		value = e
	case Symbol:
		value = e
	case Variable:
		value = en.Find(e).vars[e]
	case Bool:
		value = e
	case []expr:
		switch car, _ := e[0].(Variable); car {
		case "quote":
			value = e[1]
		case "'":
			fmt.Println(e[1])
			value = e[1]
		case "if":
			if eval(e[1], en).(Bool) {
				value = eval(e[2], en)
			} else {
				value = eval(e[3], en)
			}
		case "cond":
			for i := 1; i < len(e)-1; i++ {
				if eval(e[i].([]expr)[0], en).(Bool) {
					value = eval(e[i].([]expr)[1], en)
					return
				}
			}
			//'else' in the cond expression currently is typed as a Variable
			var elseExpression expr = Variable("else")
			if e[len(e)-1].([]expr)[0] == elseExpression {
				value = eval(e[len(e)-1].([]expr)[1], en)
			}
		case "set!":
			v := e[1].(Variable)
			en.Find(v).vars[v] = eval(e[2], en)
			value = "ok"
		case "letrec":
			var localEnv env
			localEnv = env{vars{}, en} //create new local env and set the outer env to the current one
			definitions := e[1].([]expr)
			for i := 0; i < len(definitions); i++ {
				curr := definitions[i].([]expr)
				v := curr[0].(Variable)
				localEnv.vars[v] = eval(curr[1], &localEnv)
			}
			//Maybe extend for more than one expression
			body := e[2].(expr)
			value = eval(body, &localEnv)

		case "define":
			v := e[1].(Variable)
			en.vars[v] = eval(e[2], en)
			value = "ok"
		case "lambda":
			value = proc{e[1], e[2], en}
		case "λ":
			value = proc{e[1], e[2], en}
		case "begin":
			for _, i := range e[1:] {
				value = eval(i, en)
			}
		default:
			operands := e[1:]
			values := make([]expr, len(operands))
			for i, x := range operands {
				values[i] = eval(x, en)
			}
			value = apply(eval(e[0], en), values)
		}
	default:
		log.Println("Failed evaluation - unknown expression", e)
	}
	return
}

func REPL() {
	var quit bool = false
	reader := bufio.NewReader(os.Stdin)
	for !quit {
		fmt.Print(">> ")
		inputLine, _ := reader.ReadString('\n')

		if runtime.GOOS == "windows" {
			inputLine = strings.TrimRight(inputLine, "\r\n")
		} else {
			inputLine = strings.TrimRight(inputLine, "\n")
		}

		if strings.Compare("quit", inputLine) == 0 {
			quit = true
		} else {
			fmt.Println(eval(interpret(inputLine), &globalEnv))
		}
	}
}

func evalFile(path string) {
	data := tokenize(readFromFile(path))
	parts := findParts(data)
	for i := range parts {
		currentPart := data[parts[i][0] : parts[i][1]+1]
		r := readFromTokens(&currentPart)
		fmt.Println(eval(r, &globalEnv))
	}
}

func main() {
	//REPL()
	evalFile("./program_files/letrec.rscm")
}
