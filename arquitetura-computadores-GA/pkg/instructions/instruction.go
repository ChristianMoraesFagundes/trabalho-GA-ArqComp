package instructions

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	SUPPORTED_INSTRUCTIONS                = []string{"add", "sum", "beq", "lw", "sw", "noop", "halt"}
	labels                 map[string]int = make(map[string]int)
)

type MemInstruction struct {
	OpCode   string
	Operand1 int
	Operand2 int
	Operand3 int
}

func ParseInstructions(buffer []byte) (instructions []MemInstruction, err error) {
	str := string(buffer)
	str = strings.TrimRight(str, "\n")
	lines := strings.Split(str, "\n")

	for i, line := range lines {
		lines[i] = strings.ReplaceAll(line, "\r", "")
	}

	for i, line := range lines {
		lineParsed := strings.Split(line, " ")
		if !slices.Contains(SUPPORTED_INSTRUCTIONS, lineParsed[0]) {
			if _, ok := labels[lineParsed[0]]; !ok {
				labels[lineParsed[0]] = i
				lines[i] = strings.ReplaceAll(lines[i], lineParsed[0], "")
			}
		}
	}

	for _, line := range lines {
		line = strings.TrimLeft(line, " ")
		lineParsed := strings.Split(line, " ")
		var instruction MemInstruction

		instruction, err = newInstruction(lineParsed)
		if err != nil {
			return
		}
		instructions = append(instructions, instruction)
	}
	return
}

func newInstruction(line []string) (MemInstruction, error) {
	operands := make([]int, 3)
	var err error

	if !slices.Contains(SUPPORTED_INSTRUCTIONS[4:], line[0]) {
		for i := 1; i < 4; i++ {
			operands[i-1], err = strconv.Atoi(line[i])
			if err != nil {
				if v, ok := labels[line[i]]; ok {
					operands[i-1] = v
				} else {
					return MemInstruction{}, err
				}
			}
		}
	}

	instruction := MemInstruction{
		OpCode:   line[0],
		Operand1: operands[0],
		Operand2: operands[1],
		Operand3: operands[2],
	}

	return instruction, nil
}

func (m *MemInstruction) String() string {
	return fmt.Sprintf("{\n\t\"opCode\": %s,\n\t\"operand1\": %d,\n\t\"operand2\": %d,\n\t\"operand3\": %d,\n}", m.OpCode, m.Operand1, m.Operand2, m.Operand3)
}
