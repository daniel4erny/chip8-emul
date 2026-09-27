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

func (c *CPU) exec(op [2]byte) {
	first_hex := op[0] >> 4
	const vf = 15

	switch first_hex {
	case 0x6: // mov reg, const
		reg := int(op[0] & 0x0F)
		val := op[1]
		c.regs[reg] = val
	case 0x7: // add reg, const
		reg := int(op[0] & 0x0F)
		val := op[1]
		res := val + c.regs[reg]
		c.regs[reg] = res
	case 0x8: // reg, reg operations
		last_hex := op[1] & 0x0F
		reg_y := op[1] >> 4
		reg_x := op[0] & 0x0F
		switch last_hex {
		case 0x0: // mov reg_x, reg_y
			c.regs[reg_x] = c.regs[reg_y]
		case 0x1: // OR reg_x, reg_y => reg_x
			c.regs[reg_x] = c.regs[reg_x] | c.regs[reg_y]
		case 0x2: // AND reg_x, reg_y => reg_x
			c.regs[reg_x] = c.regs[reg_x] & c.regs[reg_y]
		case 0x3: // XOR reg_x, reg_y => reg_x
			c.regs[reg_x] = c.regs[reg_x] ^ c.regs[reg_y]
		case 0x4: // reg_x + reg_y => reg_x (Carry)
			res := uint16(c.regs[reg_x]) + uint16(c.regs[reg_y])
			if res > 255 {
				c.regs[vf] = 1
			} else {
				c.regs[vf] = 0
			}
			c.regs[reg_x] = byte(res)
		case 0x5: // reg_x - reg_y => reg_x (Borrow)
			if c.regs[reg_x] >= c.regs[reg_y] {
				c.regs[vf] = 1
			} else {
				c.regs[vf] = 0
			}
			c.regs[reg_x] = c.regs[reg_x] - c.regs[reg_y]
		}
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
// 8ry4 	add vr,vy 	add register vy to vr,carry in vf 	(moje note: vf je 0 pokud borrow a 1 pokud ne)
// 8ry5 	sub vr,vy 	subtract register vy from vr,borrow in vf  (moje note: vf je 0 pokud borrow a 1 pokud ne)