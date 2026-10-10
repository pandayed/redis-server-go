package resp

import "fmt"

func SerializeSimpleString(s string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", s))
}

func SerializeError(msg string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", msg))
}

func SerializeInteger(n int) []byte {
	return []byte(fmt.Sprintf(":%d\r\n", n))
}

func SerializeBulkString(s string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
}

func SerializeNullBulkString() []byte {
	return []byte("$-1\r\n")
}

func SerializeArray(elements [][]byte) []byte {
	result := []byte(fmt.Sprintf("*%d\r\n", len(elements)))
	for _, elem := range elements {
		result = append(result, elem...)
	}
	return result
}

func SerializeNullArray() []byte {
	return []byte("*-1\r\n")
}
