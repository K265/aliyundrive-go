// go lang implementation for calculating proof_code
// see aliyun-drive-fe/aliyun-drive/2.0.12-web/web/bundle.js
package drive

import "strconv"

func I(n int32, t int32) int32 {
	var e = (65535 & n) + (65535 & t)
	return ((n>>16)+(t>>16)+(e>>16))<<16 | 65535&e
}

func O(n int32, t int32, e int32, a int32, r int32, o int32) int32 {
	l := I(I(t, n), I(a, o))
	A := r
	return I(l<<A|int32(uint32(l)>>(32-A)), e)
}

func L(n int32, t int32, e int32, a int32, r int32, i int32, l int32) int32 {
	return O(t&e|^t&a, n, t, r, i, l)
}

func A(n int32, t int32, e int32, a int32, r int32, i int32, l int32) int32 {
	return O(t&a | e & ^a, n, t, r, i, l)
}

func S(n int32, t int32, e int32, a int32, r int32, i int32, l int32) int32 {
	return O(t^e^a, n, t, r, i, l)
}

func C(n int32, t int32, e int32, a int32, r int32, i int32, l int32) int32 {
	return O(e^(t|^a), n, t, r, i, l)
}

func D(n string) []int32 {
	e := make([]int32, (len(n)-1)>>2+1)
	var a = 8 * len(n)
	t := 0
	for t < a {
		e[t>>5] |= int32((255 & int(n[t/8])) << (t % 32))
		t += 8
	}
	return e
}

func P(n []int32, t int32) []int32 {
	var a, r, o, p int32
	arr := make([]int32, 14+((t+64)>>9<<4)+2)
	copy(arr, n)
	n = arr
	n[t>>5] |= 128 << (t % 32)
	n[14+((t+64)>>9<<4)] = t
	u := int32(1732584193)
	d := int32(-271733879)
	f := int32(-1732584194)
	m := int32(271733878)
	for e := 0; e < len(n)-1; e += 16 {
		a = u
		r = d
		o = f
		p = m
		u = L(u, d, f, m, n[e], 7, -680876936)
		m = L(m, u, d, f, n[e+1], 12, -389564586)
		f = L(f, m, u, d, n[e+2], 17, 606105819)
		d = L(d, f, m, u, n[e+3], 22, -1044525330)
		u = L(u, d, f, m, n[e+4], 7, -176418897)
		m = L(m, u, d, f, n[e+5], 12, 1200080426)
		f = L(f, m, u, d, n[e+6], 17, -1473231341)
		d = L(d, f, m, u, n[e+7], 22, -45705983)
		u = L(u, d, f, m, n[e+8], 7, 1770035416)
		m = L(m, u, d, f, n[e+9], 12, -1958414417)
		f = L(f, m, u, d, n[e+10], 17, -42063)
		d = L(d, f, m, u, n[e+11], 22, -1990404162)
		u = L(u, d, f, m, n[e+12], 7, 1804603682)
		m = L(m, u, d, f, n[e+13], 12, -40341101)
		f = L(f, m, u, d, n[e+14], 17, -1502002290)
		d = L(d, f, m, u, n[e+15], 22, 1236535329)
		u = A(u, d, f, m, n[e+1], 5, -165796510)
		m = A(m, u, d, f, n[e+6], 9, -1069501632)
		f = A(f, m, u, d, n[e+11], 14, 643717713)
		d = A(d, f, m, u, n[e], 20, -373897302)
		u = A(u, d, f, m, n[e+5], 5, -701558691)
		m = A(m, u, d, f, n[e+10], 9, 38016083)
		f = A(f, m, u, d, n[e+15], 14, -660478335)
		d = A(d, f, m, u, n[e+4], 20, -405537848)
		u = A(u, d, f, m, n[e+9], 5, 568446438)
		m = A(m, u, d, f, n[e+14], 9, -1019803690)
		f = A(f, m, u, d, n[e+3], 14, -187363961)
		d = A(d, f, m, u, n[e+8], 20, 1163531501)
		u = A(u, d, f, m, n[e+13], 5, -1444681467)
		m = A(m, u, d, f, n[e+2], 9, -51403784)
		f = A(f, m, u, d, n[e+7], 14, 1735328473)
		d = A(d, f, m, u, n[e+12], 20, -1926607734)
		u = S(u, d, f, m, n[e+5], 4, -378558)
		m = S(m, u, d, f, n[e+8], 11, -2022574463)
		f = S(f, m, u, d, n[e+11], 16, 1839030562)
		d = S(d, f, m, u, n[e+14], 23, -35309556)
		u = S(u, d, f, m, n[e+1], 4, -1530992060)
		m = S(m, u, d, f, n[e+4], 11, 1272893353)
		f = S(f, m, u, d, n[e+7], 16, -155497632)
		d = S(d, f, m, u, n[e+10], 23, -1094730640)
		u = S(u, d, f, m, n[e+13], 4, 681279174)
		m = S(m, u, d, f, n[e], 11, -358537222)
		f = S(f, m, u, d, n[e+3], 16, -722521979)
		d = S(d, f, m, u, n[e+6], 23, 76029189)
		u = S(u, d, f, m, n[e+9], 4, -640364487)
		m = S(m, u, d, f, n[e+12], 11, -421815835)
		f = S(f, m, u, d, n[e+15], 16, 530742520)
		d = S(d, f, m, u, n[e+2], 23, -995338651)
		u = C(u, d, f, m, n[e], 6, -198630844)
		m = C(m, u, d, f, n[e+7], 10, 1126891415)
		f = C(f, m, u, d, n[e+14], 15, -1416354905)
		d = C(d, f, m, u, n[e+5], 21, -57434055)
		u = C(u, d, f, m, n[e+12], 6, 1700485571)
		m = C(m, u, d, f, n[e+3], 10, -1894986606)
		f = C(f, m, u, d, n[e+10], 15, -1051523)
		d = C(d, f, m, u, n[e+1], 21, -2054922799)
		u = C(u, d, f, m, n[e+8], 6, 1873313359)
		m = C(m, u, d, f, n[e+15], 10, -30611744)
		f = C(f, m, u, d, n[e+6], 15, -1560198380)
		d = C(d, f, m, u, n[e+13], 21, 1309151649)
		u = C(u, d, f, m, n[e+4], 6, -145523070)
		m = C(m, u, d, f, n[e+11], 10, -1120210379)
		f = C(f, m, u, d, n[e+2], 15, 718787259)
		d = C(d, f, m, u, n[e+9], 21, -343485551)
		u = I(u, a)
		d = I(d, r)
		f = I(f, o)
		m = I(m, p)
	}
	return []int32{u, d, f, m}
}

func U(n []int32) string {
	e := ""
	a := 32 * len(n)
	for t := 0; t < a; t += 8 {
		e += string(rune((uint32(n[t>>5]) >> (t % 32)) & 255))
	}
	return e
}

func H(n string) string {
	return U(P(D(n), int32(8*len(n))))
}

func F(n string) string {
	r := ""
	a := "0123456789abcdef"
	for _, i := range n {
		t := uint32(i)
		r += string(a[(t>>4)&15]) + string(a[15&t])
	}
	return r
}

func G(n string) string {
	return F(H(n))
}

func CalcProofOffset(accessToken string, fileSize int64) int64 {
	h := G(accessToken)[:16]
	a, _ := strconv.ParseUint(h, 16, 64)
	r := a % uint64(fileSize)
	return int64(r)
}
