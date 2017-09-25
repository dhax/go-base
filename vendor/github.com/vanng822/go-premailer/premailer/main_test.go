package premailer

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Test starting")
	args := os.Args[:]
	retCode := m.Run()
	os.Args = args
	fmt.Println("Test ending")
	os.Exit(retCode)
}
