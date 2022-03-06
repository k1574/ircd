package input

import(
	"bufio"
	"errors"
)

/* Returns truncated bytes by size or delimiter. */
func
ReadTrunc(rd *bufio.Reader, siz int, delim []byte) ([]byte, error) {
	dlen := len(delim)	
	if dlen <= 0 {
		return nil, errors.New("delimiter length cannot be 0 or less")
	}
	var(
		ret []byte
		peakLen = siz - dlen 
		b byte
		buf []byte
	)

	buf = make([]byte, 1)
	
	i := 0
	j := 0
	for ;  i < peakLen ; i++ {
		n, err := rd.Read(buf)

		if n == 0 {
			return nil, CIC
		} else if err == io.EOF {
			break;
		} else if err != nil {
			return nil, err
		}

		b = buf[0]
		ret = append(ret, b)

		if b == delim[j] {
			if j == dlen - 1 { break }
			j++
		} else {
			j = 0
		}
	}

	if i == peakLen {
		ret = append(ret, delim...)
	}
	
	
	return ret, nil
}