package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func readFile(path string) (*os.File, *bytes.Buffer, error) {
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		return nil, nil, err
	}

	b := new(bytes.Buffer)
	if _, err := b.ReadFrom(f); err != nil {
		return nil, nil, err
	}

	return f, b, nil
}

func confirm(msg string) bool {
	fmt.Print(msg + "\nAre you sure you want to continue? [y/N]: ")

	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err.Error())
	}

	if strings.TrimSpace(strings.TrimSuffix(str, "\n")) == "y" {
		return true
	}

	return false
}
