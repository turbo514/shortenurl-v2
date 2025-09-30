package base62

import "errors"

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeBase62(u uint64) string {
	if u == 0 {
		return "0"
	}
	var buf [11]byte // uint64 最大 2^64-1 → base62 最多 11 位
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = base62Chars[u%62]
		u /= 62
	}
	return string(buf[i:])
}

func DecodeBase62(s string) (uint64, error) {
	var u uint64
	for _, c := range s {
		var v byte
		switch {
		case c >= '0' && c <= '9':
			v = byte(c - '0')
		case c >= 'A' && c <= 'Z':
			v = byte(c-'A') + 10
		case c >= 'a' && c <= 'z':
			v = byte(c-'a') + 36
		default:
			return 0, errors.New("invalid base62 character")
		}
		u = u*62 + uint64(v)
	}
	return u, nil
}
