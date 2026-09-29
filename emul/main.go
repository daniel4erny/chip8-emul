package main

import (
	"fmt"
	"log"
	"math/rand/v2"
)

type CPU struct {
	ram [4096]byte
	regs [16]byte
	i_reg uint16
	pc_reg uint16
	sp_reg uint8
	stack [16]uint16
	screen [64 * 32]byte //I will hate myself for this

	shift_quirk bool
}

func newCPU() *CPU {
	fonts := []byte{
	    0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	    0x20, 0x60, 0x20, 0x20, 0x70, // 1
	    0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	    0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	    0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	    0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	    0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	    0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	    0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	    0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	    0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	    0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	    0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	    0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	    0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	    0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	var ram [4096]byte
	copy(ram[0:], fonts)
	return &CPU{
		pc_reg: 0x200,
		ram: ram,
	}
}

func (c *CPU) set_i(val uint16) {
	c.i_reg = val % 4096
}

func (c *CPU) draw(op [2]byte) {
	// to veme info a zacne na vr a ry, pak to jde dal a je to ve for loopu s poctem iteraci S,
	// ten vzdycky vezme pointer na byte a pro kazdy bit vyxoruje pixel na kterem stojime
	// vf je 1 pokud jsme xorem vypli nejaky bit (1^1) 
	// drys		 sprite rx,ry,s 	 Draw sprite at screen location rx,ry height s WIDTH 8
	reg_y := op[0] & 0x0F
	reg_x := op[1] >> 4

	pos_y := c.regs[reg_y]
	pos_x := c.regs[reg_x]
	height := int(op[1] & 0x0F)
	sprite_ptr := c.i_reg
	masky := []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01}

	for y := 0; y < height; y++ {
		for x, mask := range masky {
			
		}
	}
}

func (c *CPU) exec() {
	op := [2]byte{c.ram[c.pc_reg], c.ram[c.pc_reg+1]}
	whole_op := uint16(op[0])<<8 | uint16(op[1])
	first_hex := op[0] >> 4
	const vf = 15

	switch first_hex {
	case 0x1: // jmp to last 12 bits
		high_addr := op[0] << 4
		low_addr := op[1]
		addr := (uint16(high_addr) << 8) | uint16(low_addr)
		c.pc_reg = addr
	case 0x2:
		high_addr := op[0] << 4
		low_addr := op[1]
		addr := (uint16(high_addr) << 8) | uint16(low_addr)
		if c.sp_reg >= 15 {
			log.Fatal("SUBROUTINES EXCEEDED MAXIMUM DEPTH (16)")
		} 
		c.stack[c.sp_reg] = addr
		c.sp_reg++
	case 0x3:
		reg := int(op[0] & 0x0F)
		val := op[1]
		if c.regs[reg] == val {
			c.pc_reg += 4
		} 
		c.pc_reg += 4
	case 0x4:
		reg := int(op[0] & 0x0F)
		val := op[1]
		if c.regs[reg] < val || c.regs[reg] > val {
			c.pc_reg += 4
		} 
		c.pc_reg += 4
	case 0x6: // mov reg, const
		reg := int(op[0] & 0x0F)
		val := op[1]
		c.regs[reg] = val
		c.pc_reg += 2
	case 0x7: // add reg, const
		reg := int(op[0] & 0x0F)
		val := op[1]
		res := val + c.regs[reg]
		c.regs[reg] = res
		c.pc_reg += 2
	case 0x8: // reg, reg operations
		last_hex := op[1] & 0x0F
		reg_y := op[1] >> 4
		reg_x := op[0] & 0x0F
		switch last_hex {
		case 0x0: // mov reg_x, reg_y
			c.regs[reg_x] = c.regs[reg_y]
			c.pc_reg += 2
		case 0x1: // OR reg_x, reg_y => reg_x
			c.regs[reg_x] = c.regs[reg_x] | c.regs[reg_y]
			c.pc_reg += 2
		case 0x2: // AND reg_x, reg_y => reg_x
			c.regs[reg_x] = c.regs[reg_x] & c.regs[reg_y]
			c.pc_reg += 2
		case 0x3: // XOR reg_x, reg_y => reg_x
			c.regs[reg_x] = c.regs[reg_x] ^ c.regs[reg_y]
			c.pc_reg += 2
		case 0x4: // reg_x + reg_y => reg_x (Carry)
			res := uint16(c.regs[reg_x]) + uint16(c.regs[reg_y])
			if res > 255 {
				c.regs[vf] = 1
			} else {
				c.regs[vf] = 0
			}
			c.regs[reg_x] = byte(res)
			c.pc_reg += 2
		case 0x5: // reg_x - reg_y => reg_x (Borrow)
			if c.regs[reg_x] >= c.regs[reg_y] {
				c.regs[vf] = 1
			} else {
				c.regs[vf] = 0
			}
			c.regs[reg_x] = c.regs[reg_x] - c.regs[reg_y]
			c.pc_reg += 2
		case 0x6: // note pro mě, je to dělení dvěma a pak je v vf jestli to bylo lichý (1 jestli jo a 0 jestli ne)
			if c.shift_quirk {
			    val := c.regs[reg_y]
			    c.regs[vf] = val & 1
			    c.regs[reg_x] = val >> 1
			} else {
			    c.regs[vf] = c.regs[reg_x] & 1
			    c.regs[reg_x] = c.regs[reg_x] >> 1
			}					
			c.pc_reg += 1
		case 0x7:
			if c.regs[reg_x] >= c.regs[reg_y] {
				c.regs[vf] = 1
			} else {
				c.regs[vf] = 0
			}
			c.regs[reg_x] = c.regs[reg_y] - c.regs[reg_x]
			c.pc_reg += 2
		case 0xe:
			if c.shift_quirk {
				val := c.regs[reg_y]
				c.regs[vf] = (val >> 7) & 1
				c.regs[reg_x] = val << 1
			} else {
				c.regs[vf] = (c.regs[reg_x] >> 7) & 1
				c.regs[reg_x] = c.regs[reg_x] << 1
			}
			c.pc_reg += 2
		}
	case 0x9: // another skip func
		reg_y := op[1] >> 4
		reg_x := op[0] & 0x0F
		if c.regs[reg_y] > c.regs[reg_x] || c.regs[reg_y] < c.regs[reg_x] {
			c.pc_reg += 2
		} 
		c.pc_reg += 2
	case 0xa:
		high_val := op[0] << 4
		low_val := op[1]
		val := (uint16(high_val) << 8) | uint16(low_val)
		c.set_i(val)
		c.pc_reg += 2
	case 0xb:
		high_val := op[0] << 4
		low_val := op[1]
		val := (uint16(high_val) << 8) | uint16(low_val)
		val += uint16(c.regs[0])
		c.set_i(val)
		c.pc_reg += 2
	case 0xc: //the number is actually pseudorandom (hate chip-8)
		reg := op[0] & 0x0F
		val := op[1]
		random_byte := byte(rand.IntN(256))
		c.regs[reg] = random_byte & val
		c.pc_reg += 2
	case 0xd:
		c.draw(op) // gotta do seperate function for ts
		c.pc_reg += 2
	case 0xf:
		switch op[1] {
		case 0x29: // put the pointer to char in ram to vi
			reg := op[0] & 0x0F
			c.i_reg = uint16(c.regs[reg]) * 5 // one char is 5 byte, therefore we multiply the value by five
			c.pc_reg += 2
		}
	default:
		if whole_op == 0x00ee {  // return from subroutine call 
			if c.sp_reg == 0 {
				log.Fatal("SUBROUTINES EXCEEDED MINIMUM DEPTH (0)")
			} 
			c.pc_reg = c.stack[c.sp_reg -1]
			c.sp_reg--
		}
	} 
}

func main() {
	fmt.Println("hello world")
}


//FIRST TO-DO
// 1 6rxx 	mov vr,xx 	move constant to register r 	
// 1 7rxx 	add vr,vx 	add constant to register r 	No carry generated
// 1 8ry0 	mov vr,vy 	move register vy into vr 	
// 1 8ry1 	or rx,ry 	or register vy into register vx 	
// 1 8ry2 	and rx,ry 	and register vy into register vx 	
// 1 8ry3 	xor rx,ry 	exclusive or register ry into register rx 	
// 1 8ry4 	add vr,vy 	add register vy to vr,carry in vf 	(moje note: vf je 0 pokud borrow a 1 pokud ne)
// 1 8ry5 	sub vr,vy 	subtract register vy from vr,borrow in vf  (moje note: vf je 0 pokud borrow a 1 pokud ne)

// SECOND TO-DO
// 1 1xxx 	jmp xxx 	jump to address xxx 	
// 1 2xxx 	jsr xxx 	jump to subroutine at address xxx 	16 levels maximum
// 1 3rxx 	skeq vr,xx 	skip if register r = constant 	
// 1 4rxx 	skne vr,xx 	skip if register r <> constant 	
// 1 5ry0 	skeq vr,vy 	skip f register r = register y 
// 1 00EE 	rts 	return from subroutine call 

// THIRD TO-DO
// 1 8r06 	shr vr 	shift register vy right, bit 0 goes into register vf 	
// 1 8ry7 	rsb vr,vy 	subtract register vr from register vy, result in vr 	vf set to 1 if borrows
// 1 8r0e 	shl vr 	shift register vr left,bit 7 goes into register vf 	
// 1 9ry0 	skne rx,ry 	skip if register rx <> register ry 	
// 1 axxx 	mvi xxx 	Load index register with constant xxx 	
// 1 bxxx 	jmi xxx 	Jump to address xxx+register v0 	
// 1 crxx 	rand vr,xxx    	vr = random number less than or equal to xxx 	

// FOURTH-TODO
// 1 fr29 	font vr 	point I to the sprite for hexadecimal character in vr 	Sprite is 5 bytes high
// 1 fr1e 	adi vr 	add register vr to the index register
// drys 	sprite rx,ry,s 	Draw sprite at screen location rx,ry height s 	Sprites stored in memory at location in index register, maximum 8 bits wide. Wraps around the screen. If when drawn, clears a pixel, vf is set to 1 otherwise it is zero. All drawing is xor drawing (e.g. it toggles the screen pixels
