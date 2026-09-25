package main

import (
	"fmt"
)

type CPU struct {
	ram [4096]byte
	regs [16]byte
	i_reg [2]byte
	pc_reg [2]byte
	sp_reg [1]byte
	Stack [16]uint16
}

func main() {
	fmt.Println("hello world")
}