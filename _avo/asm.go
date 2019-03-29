// +build ignore

package main

import (
	"strconv"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

// siphash 1-3
const (
	cROUND = 1
	dROUND = 3
)

func sipround(v0, v1, v2, v3 Register) {
	ADDQ(v1, v0)
	ADDQ(v3, v2)
	ROLQ(Imm(13), v1)
	ROLQ(Imm(16), v3)
	XORQ(v0, v1)
	XORQ(v2, v3)

	ROLQ(Imm(32), v0)

	ADDQ(v1, v2)
	ADDQ(v3, v0)
	ROLQ(Imm(17), v1)
	ROLQ(Imm(21), v3)
	XORQ(v2, v1)
	XORQ(v0, v3)

	ROLQ(Imm(32), v2)

}

func main() {
	Package("github.com/dgryski/go-sip13")

	makeSip("Sum64", "func(k0, k1 uint64, p []byte) uint64")
	makeSip("Sum64Str", "func(k0, k1 uint64, p string) uint64")

	Generate()

}

func makeSip(fname, fproto string) {
	TEXT(fname, NOSPLIT, fproto)

	reg_v0 := GP64()
	reg_v1 := GP64()
	reg_v2 := GP64()
	reg_v3 := GP64()

	Load(Param("k0"), reg_v0)
	MOVQ(reg_v0, reg_v2)
	Load(Param("k1"), reg_v1)
	MOVQ(reg_v1, reg_v3)

	reg_magic := GP64()
	MOVQ(Imm(0x736f6d6570736575), reg_magic)
	XORQ(reg_magic, reg_v0)
	MOVQ(Imm(0x646f72616e646f6d), reg_magic)
	XORQ(reg_magic, reg_v1)
	MOVQ(Imm(0x6c7967656e657261), reg_magic)
	XORQ(reg_magic, reg_v2)
	MOVQ(Imm(0x7465646279746573), reg_magic)
	XORQ(reg_magic, reg_v3)

	reg_p := Load(Param("p").Base(), GP64())
	reg_p_len := Load(Param("p").Len(), GP64())

	reg_b := GP64()
	MOVQ(reg_p_len, reg_b)
	SHLQ(Imm(56), reg_b)

	reg_m := GP64()

	loop_end := "loop_end"
	loop_begin := "loop_begin"
	CMPQ(reg_p_len, Imm(8))
	JL(LabelRef(loop_end))

	Label(loop_begin)
	MOVQ(Mem{Base: reg_p}, reg_m)
	XORQ(reg_m, reg_v3)
	for i := 0; i < cROUND; i++ {
		sipround(reg_v0, reg_v1, reg_v2, reg_v3)
	}
	XORQ(reg_m, reg_v0)

	ADDQ(Imm(8), reg_p)
	SUBQ(Imm(8), reg_p_len)
	CMPQ(reg_p_len, Imm(8))
	JGE(LabelRef(loop_begin))
	Label(loop_end)

	var labels []string
	for i := 0; i < 8; i++ {
		labels = append(labels, "sw"+strconv.Itoa(i))
	}

	for i := 0; i < 7; i++ {
		CMPQ(reg_p_len, Imm(uint64(i)))
		JE(LabelRef(labels[i]))
	}

	char := GP64()
	for i := 7; i > 0; i-- {
		Label(labels[i])
		MOVBQZX(Mem{Base: reg_p, Disp: i - 1}, char)
		SHLQ(Imm(uint64(i-1)*8), char)
		ORQ(char, reg_b)

	}

	Label(labels[0])

	XORQ(reg_b, reg_v3)
	for i := 0; i < cROUND; i++ {
		sipround(reg_v0, reg_v1, reg_v2, reg_v3)
	}
	XORQ(reg_b, reg_v0)

	XORQ(Imm(0xff), reg_v2)

	for i := 0; i < dROUND; i++ {
		sipround(reg_v0, reg_v1, reg_v2, reg_v3)
	}

	XORQ(reg_v1, reg_v0)
	XORQ(reg_v3, reg_v2)
	XORQ(reg_v2, reg_v0)

	Store(reg_v0, ReturnIndex(0))
	RET()
}
