package main

import (
	"fmt"
	"hash/crc64"
	"os"
	"strconv"
	"strings"
)

/*

From https://support.google.com/mail/answer/10313
	"Gmail doesn't recognize dots as characters within usernames,
	you can add or remove the dots from a Gmail address without
	changing the actual destination address; they'll all go to
	your inbox, and only yours."

	Usage: dotnator <email address> <service name or address> [salt]

*/

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:\n\tdotinator <email address> <service name or address> [salt]")
		os.Exit(1)
	}

	email := strings.Split(os.Args[1], "@")
	username, server := email[0], email[1]
	service := os.Args[2]
	salt := ""

	if len(os.Args) == 4 {
		salt = os.Args[3]
	}

	plus := ""

	if i := strings.Index(username, "+"); i > -1 {
		plus = username[i+1:]
		username = username[0:i]
	}

	key := append([]byte(salt), []byte(service)...)
	crc := crc64.Checksum(key, crc64.MakeTable(crc64.ECMA))
	crcp := fmt.Sprintf("%063s", strconv.FormatInt(int64(crc), 2))
	index := int(service[0]) % (63 - len(username))
	part := crcp[index : index+len(username)-1]
	dots := make(map[int]bool)

	for i := range part {
		if part[i] == '1' {
			dots[i] = true
		}
	}

	name := ""

	for i := 0; i < len(username); i++ {
		name += string(username[i])

		if _, ok := dots[i]; ok {
			name += "."
		}
	}

	if len(plus) > 0 {
		plus = "+" + plus
	}

	fmt.Printf("%s%s@%s\n", name, plus, server)
}
