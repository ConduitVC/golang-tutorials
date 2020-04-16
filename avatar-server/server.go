package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ConduitVC/avatar/identicon"
)

func getPath(hash []byte) string {
	dir, _ := os.Getwd()
	name := string(hash[:5]) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".png"
	path := filepath.Join(dir, "generated", name)

	return path
}

func hashToFile(hash []byte) {
	identicon := identicon.FromHash(hash)
	path := getPath(hash)
	file, _ := os.Create(path)

	png.Encode(file, identicon)
}

func getInput(question string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	return reader.ReadString('\n')
}

func getByteArray(values []string) string {
	sha := sha256.New()

	for _, value := range values {
		sha.Write([]byte(value))
	}

	return hex.EncodeToString(sha.Sum(nil))
}

func main() {
	ip, _ := getInput("Enter IP:")
	email, _ := getInput("Enter Email:")
	values := []string{ip, email}
	fullByteArray := getByteArray(values)

	fmt.Printf("Identicon generated for %s\n", fullByteArray[:5])

	hashToFile([]byte(fullByteArray))
}
