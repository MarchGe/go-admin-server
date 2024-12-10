package log

const EndFlag = "-[EOF]-"

func IsEnd(line string) bool {
	return line == EndFlag
}
