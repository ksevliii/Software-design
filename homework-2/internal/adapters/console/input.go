package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ErrInputRead = errors.New("input read failed")

type Input struct {
	scanner *bufio.Scanner
	output  io.Writer
}

func NewInput(reader io.Reader, output io.Writer) *Input {
	return &Input{scanner: bufio.NewScanner(reader), output: output}
}

func (in *Input) ReadLine(message string) (string, error) {
	if _, err := fmt.Fprint(in.output, message); err != nil {
		return "", fmt.Errorf("write message: %w", err)
	}
	if !in.scanner.Scan() {
		if err := in.scanner.Err(); err != nil {
			return "", fmt.Errorf("%w: %v", ErrInputRead, err)
		}
		return "", io.EOF
	}
	return strings.TrimSpace(in.scanner.Text()), nil
}

func (in *Input) ReadInt(message string) (int, error) {
	for {
		line, err := in.ReadLine(message)
		if err != nil {
			return 0, err
		}
		value, err := strconv.Atoi(line)
		if err == nil {
			return value, nil
		}
		if _, err := fmt.Fprintln(in.output, "Введите целое число в допустимом диапазоне."); err != nil {
			return 0, fmt.Errorf("write input error message: %w", err)
		}
	}
}
