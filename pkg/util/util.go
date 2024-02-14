package util

import (
	"crypto/rand"
	"fmt"
	"os"
)

const (
	checkMark    = "\u2714"
	wrongMark    = "\u2718"
	wrongSpec    = "(Wrong fields)"
	correctColor = "\033[1;34m%s\033[0m"
	errColor     = "\033[1;31m%s\033[0m"
	characters   = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func Print(resource string, check, diff bool) {
	if check {
		if diff {
			fmt.Fprintf(os.Stdout, "\033[1;31m- [%s] %s %s\033[0m\n", resource, wrongMark, wrongSpec)
		} else {
			fmt.Fprintf(os.Stdout, "\033[1;32m- [%s] %s\033[0m\n", resource, checkMark)
		}
	} else {
		fmt.Fprintf(os.Stdout, "\033[1;31m- [%s] %s\033[0m\n", resource, wrongMark)
	}
}

func GenerateName(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("failed to generate random string")
		os.Exit(1)
	}

	out := make([]byte, length)
	for i := range out {
		index := uint8(bytes[i]) % uint8(len(characters))
		out[i] = characters[index]
	}

	return string(out)
}
