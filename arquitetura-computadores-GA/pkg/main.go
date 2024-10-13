package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/juniorrodes/arquitetura-computadores-GA/pkg/instructions"
)

type State struct {
	Pc                int
	InstructionMemory []instructions.MemInstruction
	Registers         [32]int
}

// InstructionMemory carrega e armazena instruções de um arquivo
type InstructionMemory struct {
	instructions []instructions.MemInstruction
}

// LoadInstructions carrega instruções de um arquivo
func (im *InstructionMemory) LoadInstructions(filepath string) error {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	instructions, err := instructions.ParseInstructions(fileContent)
	if err != nil {
		return err
	}
	im.instructions = instructions
	return nil
}

// RegisterFile representa o banco de registradores
type RegisterFile struct {
	registers [32]int // R0 a R31
}

// Read lê o valor de um registrador
func (rf *RegisterFile) Read(regNum int) int {
	return rf.registers[regNum]
}

// Write escreve um valor em um registrador (R0 é fixo em zero)
func (rf *RegisterFile) Write(regNum int, value int) {
	if regNum != 0 { // R0 não pode ser modificado
		rf.registers[regNum] = value
	}
}

// ControlUnit controla a política de desvio
type ControlUnit struct {
	branchTaken bool
}

// HandleBranch define a política de desvio
func (controlUnit *ControlUnit) HandleBranch(condition bool) {
	// Política fixa: Não tomado
	controlUnit.branchTaken = condition
}

// PipelineStage representa um estágio do pipeline
type PipelineStage struct {
	stageName   string
	instruction string
}

// Process processa a instrução no estágio atual
func (ps *PipelineStage) Process(instruction string) {
	ps.instruction = instruction
	//fmt.Printf("%s processing instruction: %s\n", ps.stageName, instruction)
}

// imprime os estágios do pipeline
func (ps *PipelineStage) PrintState() {
	fmt.Printf("%s: %s\n", ps.stageName, ps.instruction)
}

// gerencia o fluxo das instruções pelo pipeline
type Pipeline struct {
	instructions     []instructions.MemInstruction
	registers        *RegisterFile
	controlUnit      *ControlUnit
	stages           [5]*PipelineStage
	pc               int // Program Counter
	instructionCount int // Contador de instruções executadas
}

// NewPipeline cria um novo pipeline
func NewPipeline(state State) *Pipeline {
	p := &Pipeline{
		instructions: state.InstructionMemory,
		registers:    &RegisterFile{},
		controlUnit:  &ControlUnit{},
		pc:           state.Pc,
	}

	// Inicializa os estágios do pipeline
	p.stages[0] = &PipelineStage{"IF", ""}
	p.stages[1] = &PipelineStage{"ID", ""}
	p.stages[2] = &PipelineStage{"EX", ""}
	p.stages[3] = &PipelineStage{"MEM", ""}
	p.stages[4] = &PipelineStage{"WB", ""}

	return p
}

// Run executa o pipeline
func (p *Pipeline) Run() {
	for p.pc < len(p.instructions) {
		p.cycle()
		p.printStages() // Imprime o estado de todos os estágios após cada ciclo
	}
	
	fmt.Printf("Total de instruções executadas: %d\n", p.instructionCount)
}

// Cycle realiza um ciclo do pipeline
func (p *Pipeline) cycle() {
	// Move as instruções pelos estágios do pipeline
	for i := len(p.stages) - 1; i > 0; i-- {
		p.stages[i].instruction = p.stages[i-1].instruction
	}

	// Busca a instrução (IF)
	currentInstruction := ""
	if p.pc < len(p.instructions) {
		currentInstruction = p.instructions[p.pc].String() 
	}
	p.stages[0].Process(currentInstruction)

	// Decodifica (ID) e controla o desvio
	if currentInstruction != "" && strings.Contains(currentInstruction, "BRANCH") {
		p.controlUnit.HandleBranch(false) // Política: não tomado
		if p.controlUnit.branchTaken {
			p.pc = p.resolveBranch(currentInstruction)
		} else {
			p.pc++
		}
	} else {
		p.pc++
	}

	// Incrementa o contador de instruções executadas
	p.instructionCount++
}

// ResolveBranch processa o desvio caso a condição seja verdadeira
func (p *Pipeline) resolveBranch(instruction string) int {
	avançamos o PC em 1
	return p.pc + 1 
}

// printStages imprime o estado atual 
func (p *Pipeline) printStages() {
	fmt.Printf("Pipeline Stages (PC: %d):\n", p.pc)
	for _, stage := range p.stages {
		stage.PrintState()
	}
	fmt.Println() 
}

func main() {
	state := Init()

	/*for i, instruction := range state.InstructionMemory {
		fmt.Println("instruction ", i, ": ")
		fmt.Println(instruction)
	}*/

	// Inicializa e executa o pipeline
	pipeline := NewPipeline(state)
	pipeline.Run()

}

func Init() State {
	fileContent, err := os.ReadFile("test.asm")
	if err != nil {
		log.Fatal(err)
	}

	i, err := instructions.ParseInstructions(fileContent)
	if err != nil {
		log.Fatal(err)
	}

	return State{
		Pc:                0,
		InstructionMemory: i,
	}
}
