package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io/ioutil"

	//"compiler/cmd"
)
const (
	//errors
	CMD_ARG_ERR string= "\u001B[40m\u001B[91mINVALID COMMAND\u001B[0m\n"
	METHOD_DECLARATION_ERR string = "The '{' token should not have any following token except new line feed or space"
	INSUFFICIENT_VARS_ERR string = "the number of vars on the lhs dont match with the values on the rhs"
	ALREADY_DECLARED_ERR string = "the variables have already been declared please remove the \"var\" specifier"
	BUILD_FAIL_ERR string = "FATAL ERROR\nBUILD FAILED"

	//syntax
	FUNCTION_PARAM string = "<-"
	NORMAL_ASSIGNMENT string = "="
	SYNTACTIC_ASSIGNMENT string = "<-"

	COMMENT_START string = "#" 

	METHOD_DECLARATION string = "method "
	VAR_DECALRATION string = "var "

	//TYPES
	//varTypes
)

//structs
type actionFunc func(string) (bool,string)

type CmdArgs struct {
	action actionFunc
}

type Node struct {
	value string
	childAction string
}

type AssignmentLinkedList struct {
	child *Node
}

type Data struct {
	value string
}

type MethodData struct {
	parameters string
	data Block
}

type VarSyntaxData struct {
	data Block
}

type Block struct {
	data []string
}

//global variables
var headNode *Node

var varSave map[string]Data
var varSyntaxSave map[string]VarSyntaxData
var varSyntaxNames []string
var methodSave map[string]MethodData

var buldFailure bool

func main() {
	fmt.Println("Welcome to VIPER Lang")

	for true {
		methodSave = make(map[string]MethodData)
		varSave = make(map[string]Data)
		varSyntaxSave = make(map[string]VarSyntaxData)
		userInput := ""
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("\u001B[92m> \u001B[0m")

		userInput ,_ = reader.ReadString('\n')
		userInput = strings.Replace(userInput, "\n", "", -1)

		argHandler := InititializeCMD()
		program := Interpret(userInput,argHandler)

		startExec(program)
	}
}

//initializes the commands
func InititializeCMD() [3]CmdArgs{
	var argHandler[3] CmdArgs

	defFile := "lethalityTest.vpr"

	argHandler[0].action = func(args string) (bool,string){
		if (strings.TrimSpace(args) == "run -d") {
			str := CommenceReading(defFile)
			return true,str
		}
		return false,""
	}

	argHandler[1].action = func(args string) (bool,string){

		if strings.HasPrefix(args,"run ") {
			parts := []rune(args)
			fileName := string(parts[4:])

			str := CommenceReading(fileName)
			return true,str
		}

		return false,""

	}

	argHandler[2].action = func(args string) (bool,string){
		if args == "quit"||args == "exit" {
			os.Exit(0)
		}
		return false,""
	}

	return argHandler
}

//reads the file
func CommenceReading(fileName string) string {
	data,ERR := ioutil.ReadFile(fileName)

	if ERR != nil {
		fmt.Print("\u001B[40m\u001B[91m",ERR,"\u001B[0m\n")
		return ""
	}
	program := string(data)

	return program
}

//interprets whatever commans is given
func Interpret(input string,argHandler[3] CmdArgs) string{
	for i := 0; i < len(argHandler); i++ {
		
		success,data := argHandler[i].action(input)


		if success {
			return data
		}
	}
	return CMD_ARG_ERR
}

//
func startExec(args string) {
	if args == CMD_ARG_ERR {
		fmt.Print(CMD_ARG_ERR)
	} else {
		AssignmentRun(args)
	}
}

func AssignmentRun(program string) {
	//the splitting on new line works
	splitCode := strings.Split(program,"\n")

	StaticallyInitialize(splitCode)
}

func StaticallyInitialize(program []string) {
	for i := 0; i < len(program); i++ {
		program[i] = strings.TrimSpace(program[i])

		if strings.HasPrefix(program[i],COMMENT_START) {
			program[i] = ""
			continue
		}

		if strings.HasPrefix(program[i],METHOD_DECLARATION) {

			parts := []rune(program[i])
			name := string(parts[len(METHOD_DECLARATION):])
			temp := strings.Split(name,NORMAL_ASSIGNMENT)
			name = strings.TrimSpace(temp[0])
			parameters := strings.TrimSpace(temp[1])

			declareMethod(name,parameters,program,i)
		}

		if strings.HasPrefix(program[i],VAR_DECALRATION) {
			parts := []rune(program[i])
			name := string(parts[len(VAR_DECALRATION):])

			if strings.Contains(name,SYNTACTIC_ASSIGNMENT) {
				temp := strings.Split(name,SYNTACTIC_ASSIGNMENT)
				name = strings.TrimSpace(temp[0])
				if testVarDeclaration(name) {
					fmt.Println("\u001B[40m\u001B[91m",ALREADY_DECLARED_ERR,"\u001B[0m\nline:\t",i+1)
					buldFailure = true
				}
				declareVarSyntax(name,program,i)
				continue
			}

			if strings.Contains(name,NORMAL_ASSIGNMENT) {
				temp := strings.Split(name,NORMAL_ASSIGNMENT)
				allVarNames := strings.Split(temp[0],",")
				
				allValues := strings.Split(temp[1],",")

				if (len(allVarNames) != len(allValues)) && (len(allValues) != 1){

					fmt.Println("\u001B[40m\u001B[91m",INSUFFICIENT_VARS_ERR,"\u001B[0m\n","line:\t",i+1)
					buldFailure = true
				}

				if len(allValues) == 1 {
					assignToAll(allVarNames,allValues[0],i)
				}
			}
		}
	}

	if buldFailure == true {
		fmt.Println("\u001B[40m\u001B[91m",BUILD_FAIL_ERR,"\u001B[0m")
	}
}

func declareMethod (name string,parameters string,program []string,index int) {
	elems := []rune(strings.TrimSpace(parameters))

	var startIndex,endIndex int

	for i := 0; i < len(elems); i++ {
		if elems[i] == '(' {
			startIndex = i
			continue
		}

		if elems[i] == ')' {
			endIndex = i
			break
		}
	}

	param := string(elems[startIndex+1:endIndex])
	block := getBlock(program,index,'{','}')

	node := MethodData{param,block}

	methodSave[name] = node
}

func getBlock(program []string,startIndex int,start rune,end rune) (Block){
	num_of_nested_blocks := 0
	endIndex := startIndex

	for i := startIndex; i < len(program); i++ {
		program[i] = strings.TrimSpace(program[i])
		elems := []rune(program[i])

		for j := 0; j < len(elems); j++ {
			
			if elems[j] == '#' {
				break
			} 
			
			if elems[j] == start {
				num_of_nested_blocks++
				j = j+1
				continue
			}
			
			if elems[j] == end {
				num_of_nested_blocks--
			}
		}

		if num_of_nested_blocks == 0 {
			endIndex = i
			break
		}
	}
	if (startIndex == endIndex) {
		fmt.Println("\u001B[40m\u001B[91m",METHOD_DECLARATION_ERR,"\u001B[0m\n","line:\t",startIndex+1)
		buldFailure = true
		return Block{}
	}
	snippet := program[startIndex+1:endIndex]
	return Block{snippet}
}

func declareVarSyntax(name string,program []string,index int) {	
	block := getBlock(program,index,'{','}')
	varSyntaxSave[name] = VarSyntaxData{block}
	varSyntaxNames = append(varSyntaxNames,name)
}


func assignToAll(varNames []string,value string,index int) {
	data := compute(value)
	for i := 0; i < len(varNames); i++ {
		if testVarDeclaration(varNames[i]) {
			fmt.Println("\u001B[40m\u001B[91m",ALREADY_DECLARED_ERR,"\u001B[0m\nline:\t",index+1)
		}
		varSave[varNames[i]] = data
	}
}

func testVarDeclaration(name string) bool{
	t1 := varSyntaxSave[name].data.data
	t2 := string(varSave[name].value)

	if (len(t1) != 0 || t2 != ""){
		return true
	}

	return false
}


func compute(value string) (Data) {
	return Data{}
}
