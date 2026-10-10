package resp

import (
	"errors"
	"fmt"
)

func (v RESPValue) ToCommand() ([]string, error) {
	if v.Type != Array {
		return nil, errors.New("command must be an array")
	}

	if v.Null {
		return nil, errors.New("command array is null")
	}

	command := make([]string, len(v.Array))
	for i, elem := range v.Array {
		if elem.Type == BulkString {
			if elem.Null {
				return nil, errors.New("command contains null bulk string")
			}
			command[i] = elem.Bulk
		} else if elem.Type == SimpleString {
			command[i] = elem.Str
		} else {
			return nil, fmt.Errorf("invalid command element type: %c", elem.Type)
		}
	}

	return command, nil
}
