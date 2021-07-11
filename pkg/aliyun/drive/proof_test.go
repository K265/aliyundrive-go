package drive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestI(t *testing.T) {
	assert.True(t, 1924017899 == I(-1483587690, -887361707))
}

func TestO(t *testing.T) {
	assert.True(t, 1924017899 == O(-1967175164, 1134877475, -887361707, 0, 9, -51403784))
}

func TestL(t *testing.T) {
	assert.True(t, 984724524 == L(271733878, -690939994, -271733879, -1732584194, 8414821, 12, -389564586))
}

func TestA(t *testing.T) {
	assert.True(t, 1924017899 == A(1134877475, -887361707, -1363207668, 643323946, 0, 9, -51403784))
}

func TestS(t *testing.T) {
	assert.True(t, 1154 == S(1, 2, 3, 4, 5, 6, 7))
}

func TestC(t *testing.T) {
	assert.True(t, 322 == C(1, 2, 3, 4, 5, 6, 7))
}

func TestD(t *testing.T) {
	assert.Equal(t, []int32{1684234849, 26213}, D("abcdef"))
}

func TestP(t *testing.T) {
	assert.Equal(t, []int32{391121896, -61830903, 1020832344, -1902701428}, P([]int32{1684234849, 26213}, 48))
}

func TestU(t *testing.T) {
	assert.Equal(t, "è\vP\u0017\t\u0089PüXªØ<\u008C\u0014\u0097\u008e", U([]int32{391121896, -61830903, 1020832344, -1902701428}))
}

func TestH(t *testing.T) {
	assert.Equal(t, "è\vP\u0017\t\u0089PüXªØ<\u008C\u0014\u0097\u008e", H("abcdef"))
}

func TestF(t *testing.T) {
	assert.Equal(t, "e80b5017098950fc58aad83c8c14978e", F("è\vP\u0017\t\u0089PüXªØ<\u008C\u0014\u0097\u008E"))
}

func TestG(t *testing.T) {
	assert.Equal(t, "e80b5017098950fc58aad83c8c14978e", G("abcdef"))
}

func TestCalcProofOffset(t *testing.T) {
	accessToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI0MjQyNDI0MjQyNDI0MjQyNDI0MjQyNDI0MjQyNDI0MiIsImN1c3RvbUpzb24iOiJ7XCJjbGllbnRJZFwiOlwiMjVkelgzdmJZcWt0Vnh5WFwiLFwiZG9tYWluSWRcIjpcImJqMjlcIixcInNjb3BlXCI6W1wiRFJJVkUuQUxMXCIsXCJTSEFSRS5BTExcIixcIkZJTEUuQUxMXCIsXCJVU0VSLkFMTFwiLFwiU1RPUkFHRS5BTExcIixcIlNUT1JBR0VGSUxFLkxJU1RcIixcIkJBVENIXCIsXCJPQVVUSC5BTExcIixcIklNQUdFLkFMTFwiLFwiSU5WSVRFLkFMTFwiLFwiQUNDT1VOVC5BTExcIl0sXCJyb2xlXCI6XCJ1c2VyXCIsXCJyZWZcIjpcImh0dHBzOi8vd3d3LmFsaXl1bmRyaXZlLmNvbS9cIn0iLCJleHAiOjE2MjAxMjAyNzUsImlhdCI6MTYyMDExMzAxNX0K.AlMk63QdYz8pDxbyg7CoxzJ_pMh9t1RO9e6Ri4yTLQyDghbpCuAxY4W43RHkShyYd4MmE7xNywT4IDxARKkCGkKqwCD1AOdTRLta_JVos242QPJbaEC1XEs2CdwjmID8uv4IiB1FmVb37_LVaVg-oUvgm94-zfdAocoxkUBavHE"
	fileSize := int64(1477419708)
	offset := CalcProofOffset(accessToken, fileSize)
	assert.True(t, 353042511 == offset)
}
