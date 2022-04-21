package report

import (
	"fmt"
	"os"
)

func Error(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func Info(info string) {
	fmt.Fprintln(os.Stdout, info)
}
