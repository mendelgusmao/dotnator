package dotnator

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/crc64"
	"strconv"
	"strings"
)

/*

From https://support.google.com/mail/answer/10313
	"Gmail doesn't recognize dots as characters within usernames,
	you can add or remove the dots from a Gmail address without
	changing the actual destination address; they'll all go to
	your inbox, and only yours."

	This piece of code provides a deterministic way of generating different e-mail addresses
	from an e-mail address by putting dots between its characters. So, after receiving an spam
	sent to one of your "new" addresses, you can track and have the information of what service
	forwarded your address.

	It's just a proof of concept and should not be used seriously. Some smart spammers might just
	remove the dots (and/or the plus sign and what comes after) from the addresses they receive.
	Also, the services you use can simply consider your address as not valid.

*/

func Dotnate(email, service, salt, algorithm string) (string, error) {
	address := strings.Split(email, "@")
	username, host := address[0], address[1]
	plus := ""

	if i := strings.Index(username, "+"); i > -1 {
		plus = username[i:]
		username = username[:i]
	}

	key := append([]byte(salt), []byte(service)...)
	checksum, err := checksum(key, algorithm)

	if err != nil {
		return "", err
	}

	mask := fmt.Sprintf("%063s", strconv.FormatInt(checksum, 2))
	name := make([]byte, 0)
	size := len(username)
	index := 0

	if modulus := len(mask) - size; modulus > 0 {
		index = int(checksum>>32) % modulus

		if index < 0 {
			index *= -1
		}
	}

	for i := 0; i < size; i++ {
		if i+index < len(mask) && mask[i+index] == '1' && i < size-1 {
			name = append(name, username[i], '.')
		} else {
			name = append(name, username[i])
		}
	}

	return fmt.Sprintf("%s%s@%s", name, plus, host), nil
}

func checksum(key []byte, algorithm string) (int64, error) {
	hash := func(h hash.Hash, in []byte) int64 {
		h.Write(in)
		result := h.Sum(nil)
		size := len(result)

		for i := 0; i < 64/greatestCommonDivisor(64, size)-1 && len(result) <= 64; i++ {
			result = append(result, result[i*size:]...)
		}

		return int64(binary.BigEndian.Uint64(result[0:64]))
	}

	switch algorithm {
	case "md5":
		return hash(md5.New(), key), nil
	case "sha1":
		return hash(sha1.New(), key), nil
	case "sha256":
		return hash(sha256.New(), key), nil
	case "sha512":
		return hash(sha512.New(), key), nil
	case "crc64:iso":
		return int64(crc64.Checksum(key, crc64.MakeTable(crc64.ISO))), nil
	case "crc64:ecma", "crc64":
		return int64(crc64.Checksum(key, crc64.MakeTable(crc64.ECMA))), nil
	default:
		return 0, fmt.Errorf("algorithm '%s' not implemented", algorithm)
	}
}

// https://en.wikipedia.org/wiki/Binary_GCD_algorithm
func greatestCommonDivisor(u, v int) int {
	if u == v || v == 0 {
		return u
	}

	if u == 0 {
		return v
	}

	if ^u&1 == 1 {
		if v&1 == 1 {
			return greatestCommonDivisor(u>>1, v)
		}

		return greatestCommonDivisor(u>>1, v>>1) << 1
	}

	if ^v&1 == 1 {
		return greatestCommonDivisor(u, v>>1)
	}

	if u > v {
		return greatestCommonDivisor((u-v)>>1, v)
	}

	return greatestCommonDivisor((v-u)>>1, u)
}
