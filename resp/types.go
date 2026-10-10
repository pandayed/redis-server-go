package resp

type RESPType byte

const (
	SimpleString RESPType = '+'
	Error        RESPType = '-'
	Integer      RESPType = ':'
	BulkString   RESPType = '$'
	Array        RESPType = '*'
)

type RESPValue struct {
	Type  RESPType
	Str   string
	Num   int
	Bulk  string
	Array []RESPValue
	Null  bool
}
