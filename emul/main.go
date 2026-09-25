package main

import (
	"fmt"
)

type CPU struct {
	ram [4096]byte
	regs [16]byte
	i_reg uint16
	pc_reg uint16
	sp_reg [1]byte
	Stack [16]uint16
}

func (CPU) new() *CPU {
	return &CPU{
		pc_reg: 0x200,		
	}
}

func main() {
	fmt.Println("hello world")
}

// 6rxx 	mov vr,xx 	move constant to register r 	
// 7rxx 	add vr,vx 	add constant to register r 	No carry generated
// 8ry0 	mov vr,vy 	move register vy into vr 	
// 8ry1 	or rx,ry 	or register vy into register vx 	
// 8ry2 	and rx,ry 	and register vy into register vx 	
// 8ry3 	xor rx,ry 	exclusive or register ry into register rx 	
// 8ry4 	add vr,vy 	add register vy to vr,carry in vf 	
// 8ry5 	sub vr,vy 	subtract register vy from vr,borrow in vf 	vf set to 1 if borroesws (moje note: vf je 0 pokud borrow a 1 pokud ne)