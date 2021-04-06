package base

const smalls = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

func ItoA(dst *[]byte, i int, w int) {
	var b = [20]byte{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}
	var pos = 20
	var is, ii int
	var neg bool
	if i < 0 {
		neg = true
		i = -i
	}
	for i > 99 {
		ii = i / 100
		is = (i - ii*100) * 2
		pos -= 2
		b[pos+1] = smalls[is+1]
		b[pos] = smalls[is]
		i = ii
	}
	is = i * 2
	pos--
	b[pos] = smalls[is+1]
	if i > 9 {
		pos--
		b[pos] = smalls[is]
	}
	var start = 20 - w
	if start > pos {
		start = pos
	}
	if neg {
		start--
		b[start] = '-'
	}
	*dst = append(*dst, b[start:]...)
}

func ItoAW(dst *[]byte, i int, w int) {
	var b = [20]byte{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}
	var pos = 20
	var is, ii int
	var start = 20 - w
	var neg bool
	if i < 0 {
		neg = true
		i = -i
	}
	for i > 99 {
		ii = i / 100
		is = (i - ii*100) * 2
		pos -= 2
		b[pos+1] = smalls[is+1]
		b[pos] = smalls[is]
		i = ii

		if pos <= start {
			*dst = append(*dst, b[start:]...)
			return
		}
	}
	is = i * 2
	pos--
	b[pos] = smalls[is+1]
	if i > 9 {
		pos--
		b[pos] = smalls[is]
	}
	if neg {
		start--
		b[start] = '-'
	}
	*dst = append(*dst, b[start:]...)
}
