package utils

import "math/rand"
import crd "crypto/rand"

func RandomString(length int) string {
	return RandomStringFrom(length, "123456789123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ@#$%&~+=@#$%&~+=")
}

func RandomStringFrom(length int, characters string) string {
	originRunes, genRunes := []rune(characters), make([]rune, length)
	charLen := len(originRunes)
	for i := 0; i < length; i++ {
		index := rand.Intn(charLen)
		genRunes[i] = originRunes[index]
	}
	return string(genRunes)
}

func RandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, _ = crd.Read(bytes)
	return bytes
}
