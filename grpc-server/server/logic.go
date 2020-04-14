package main

import (
	"crypto/md5"
	fmt "fmt"
)

func gravatarHash(email string) [16]byte {
	return md5.Sum([]byte(email))
}

func gravatarURL(hash [16]byte, size int32) string {
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d", hash, size)
}

func gravatar(email string, size int32) string {
	hash := gravatarHash(email)

	return gravatarURL(hash, size)
}
