package main

import (
	"fmt"
	"strconv"
)

func validateMinArgs(command []string, min int, cmdName string) error {
	if len(command) < min {
		return fmt.Errorf("wrong number of arguments for '%s' command", cmdName)
	}
	return nil
}

func parseIntArgs(args []string) ([]int, error) {
	result := make([]int, len(args))
	for i, arg := range args {
		val, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("value is not an integer or out of range")
		}
		result[i] = val
	}
	return result, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

