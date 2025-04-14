package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*()-_=+,.?/:;{}[]~"
)

// Генерация случайного email
func GenerateRandomEmail() string {
	username := randomString(8, lowercaseLetters+digits) // 8 символов для имени пользователя
	domain := randomString(5, lowercaseLetters)          // 5 символов для домена
	tld := randomString(3, lowercaseLetters)             // 3 символа для TLD

	return fmt.Sprintf("%s@%s.%s", username, domain, tld)
}

// Генерация случайного пароля
func GenerateRandomPassword(length int) string {
	// Все возможные символы для пароля
	allChars := lowercaseLetters + uppercaseLetters + digits + specialChars

	// Гарантируем, что пароль содержит хотя бы по одному символу из каждой группы
	var password strings.Builder

	// Добавляем по одному символу из каждой группы
	password.WriteByte(lowercaseLetters[randomInt(len(lowercaseLetters))])
	password.WriteByte(uppercaseLetters[randomInt(len(uppercaseLetters))])
	password.WriteByte(digits[randomInt(len(digits))])
	password.WriteByte(specialChars[randomInt(len(specialChars))])

	// Добавляем оставшиеся символы
	for i := 4; i < length; i++ {
		password.WriteByte(allChars[randomInt(len(allChars))])
	}

	// Перемешиваем символы в пароле
	pwd := []byte(password.String())
	for i := range pwd {
		j := randomInt(len(pwd))
		pwd[i], pwd[j] = pwd[j], pwd[i]
	}

	return string(pwd)
}

// Генерация случайной строки заданной длины из заданных символов
func randomString(length int, charset string) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[randomInt(len(charset))])
	}
	return sb.String()
}

// Генерация случайного числа в диапазоне [0, max)
func randomInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}
