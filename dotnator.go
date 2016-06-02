package dotnator

import (
	"fmt"
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

func Dotnate(email, service, salt string) string {
	address := strings.Split(email, "@")
	username, host := address[0], address[1]
	plus := ""

	if i := strings.Index(username, "+"); i > -1 {
		plus = username[i:]
		username = username[:i]
	}

	key := append([]byte(salt), []byte(service)...)
	crc := crc64.Checksum(key, crc64.MakeTable(crc64.ECMA))
	mask := fmt.Sprintf("%063s", strconv.FormatInt(int64(crc), 2))
	name := make([]byte, 0)
	size := len(username)
	index := 0

	if modulus := len(mask) - size; modulus > 0 {
		index = int(crc>>32) % modulus

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

	return fmt.Sprintf("%s%s@%s", name, plus, host)
}
