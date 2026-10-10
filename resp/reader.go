package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

func ReadRESP(reader *bufio.Reader) (RESPValue, error) {
	typeByte, err := reader.ReadByte()
	if err != nil {
		return RESPValue{}, err
	}

	switch RESPType(typeByte) {
	case SimpleString:
		return readSimpleString(reader)
	case Error:
		return readError(reader)
	case Integer:
		return readInteger(reader)
	case BulkString:
		return readBulkString(reader)
	case Array:
		return readArray(reader)
	default:
		return RESPValue{}, fmt.Errorf("unknown RESP type: %c", typeByte)
	}
}

func readSimpleString(reader *bufio.Reader) (RESPValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return RESPValue{}, err
	}

	return RESPValue{
		Type: SimpleString,
		Str:  line,
	}, nil
}

func readError(reader *bufio.Reader) (RESPValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return RESPValue{}, err
	}

	return RESPValue{
		Type: Error,
		Str:  line,
	}, nil
}

func readInteger(reader *bufio.Reader) (RESPValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return RESPValue{}, err
	}

	num, err := strconv.Atoi(line)
	if err != nil {
		return RESPValue{}, fmt.Errorf("invalid integer: %s", line)
	}

	return RESPValue{
		Type: Integer,
		Num:  num,
	}, nil
}

func readBulkString(reader *bufio.Reader) (RESPValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return RESPValue{}, err
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return RESPValue{}, fmt.Errorf("invalid bulk string length: %s", line)
	}

	if length == -1 {
		return RESPValue{
			Type: BulkString,
			Null: true,
		}, nil
	}

	bulk := make([]byte, length)
	_, err = io.ReadFull(reader, bulk)
	if err != nil {
		return RESPValue{}, err
	}

	reader.ReadByte()
	reader.ReadByte()

	return RESPValue{
		Type: BulkString,
		Bulk: string(bulk),
	}, nil
}

func readArray(reader *bufio.Reader) (RESPValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return RESPValue{}, err
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return RESPValue{}, fmt.Errorf("invalid array length: %s", line)
	}

	if length == -1 {
		return RESPValue{
			Type: Array,
			Null: true,
		}, nil
	}

	array := make([]RESPValue, length)
	for i := 0; i < length; i++ {
		value, err := ReadRESP(reader)
		if err != nil {
			return RESPValue{}, err
		}
		array[i] = value
	}

	return RESPValue{
		Type:  Array,
		Array: array,
	}, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if len(line) < 2 || line[len(line)-2] != '\r' {
		return "", errors.New("invalid RESP format: missing \\r\\n")
	}

	return line[:len(line)-2], nil
}
