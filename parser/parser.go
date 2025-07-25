package parser

import (
	"bufio"
	"os"
	"strings"
)

type Instruction struct {
	Cmd  string
	Args []string
}

func ParseUnderdogfile(path string) ([][]Instruction, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stages [][]Instruction
	var currentStage []Instruction

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		cmd := strings.ToUpper(parts[0])
		args := parts[1:]

		if cmd == "FROM" && len(currentStage) > 0 {
			stages = append(stages, currentStage)
			currentStage = nil
		}
		currentStage = append(currentStage, Instruction{Cmd: cmd, Args: args})
	}
	stages = append(stages, currentStage)

	return stages, nil
}
