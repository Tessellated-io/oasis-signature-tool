package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	args := os.Args

	// Validate and decode input
	if len(args) != 2 {
		panic(fmt.Errorf("usage: oasis-signature-tool <bytes base64>"))
	}
	inputBase64 := os.Args[1]

	bytes, err := base64.StdEncoding.DecodeString(inputBase64)
	if err != nil {
		panic(fmt.Errorf("unable to decode input. Are you sure it is base64? Found: \"%s\"", inputBase64))
	}

	// Oasis uses a domain separation context on sigantures
	// context for the consensus network is:
	// ASCII: oasis-core/consensus: tx for chain bb3d748def55bdfb797a2ac53ee6ee141e54cd2ab2dc2375f4a0703a178e6e55
	// base64: b2FzaXMtY29yZS9jb25zZW5zdXM6IHR4IGZvciBjaGFpbiBiYjNkNzQ4ZGVmNTViZGZiNzk3YTJhYzUzZWU2ZWUxNDFlNTRjZDJhYjJkYzIzNzVmNGEwNzAzYTE3OGU2ZTU1"
	domain := []byte("oasis-core/consensus: tx for chain bb3d748def55bdfb797a2ac53ee6ee141e54cd2ab2dc2375f4a0703a178e6e55")

	h := sha512.New512_256()
	_, _ = h.Write(bytes)
	_, _ = h.Write(domain)
	sum := h.Sum(nil)

	fmt.Println("Bytes to sign: ")
	fmt.Printf("Hex: %s\n", hex.EncodeToString(sum))
	fmt.Printf("B64: %s\n", base64.StdEncoding.EncodeToString(sum))
}
