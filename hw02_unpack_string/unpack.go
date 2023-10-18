package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const unicodeValueBackslash = 92

func Unpack(s string) (string, error) {
	var b strings.Builder
	var sbuf int // количество слэшей до текущего символа
	var cur, next rune

	r := []rune(s)

	if len(r) == 0 {
		return "", nil
	}

	// для обхода слайса рун используется окно на два элемента
	for i, j := 0, 1; i < len(r); i, j = i+1, j+1 {
		cur = r[i]
		// если текущая цифра неэкранирована - ошибка
		if isUnescapedDigit(cur, sbuf) {
			return "", ErrInvalidString
		}
		// если i еще попадает в слайс рун, а j - уже нет
		if j == len(r) {
			if canWriteLastCharacter(cur, sbuf) {
				b.WriteString(string(cur))
			}
			break
		}

		next = r[j]

		// из условия задачи - экранировать можно только цифру или слэш
		if isInvalidEscapeSequence(cur, next) {
			return "", ErrInvalidString
		}

		if canIncreaseSlashBuffer(cur, next, sbuf) {
			sbuf++
			continue
		}

		// исходя из логики дефолтных тест-кейсов:
		// если слэшей в буфере больше одного, тогда пишем текущий слэш и следующую цифру как одиночный символ
		// если слэш в буфере один - значит текущий слэш можно обработать с множителем, и мы идем дальше
		if isMultiSlashSequence(cur, next, sbuf) {
			b.WriteString(string(cur))
			b.WriteString(string(next))
			// после записей, переход на позицию после записанной одиночной цифры
			i++
			j++
			sbuf = 0
			continue
		}

		// запись повторяющихся символов
		if unicode.IsDigit(next) {
			d, err := strconv.Atoi(string(next))
			if err != nil {
				return "", ErrInvalidString
			}
			b.WriteString(strings.Repeat(string(cur), d))

			// после репита, переход на позицию после числа-множителя
			i++
			j++
			sbuf = 0
			continue
		}
		// запись символов без множителя
		b.WriteString(string(cur))
		sbuf = 0
	}

	return b.String(), nil
}

func canWriteLastCharacter(cur rune, sbuf int) bool {
	return !isSlash(cur) || sbuf > 0
}

func canIncreaseSlashBuffer(cur, next rune, sbuf int) bool {
	return isSlash(cur) && (isSlash(next) || unicode.IsDigit(next) && sbuf == 0)
}

func isMultiSlashSequence(cur, next rune, sbuf int) bool {
	return isSlash(cur) && unicode.IsDigit(next) && sbuf > 1
}

func isUnescapedDigit(cur rune, sbuf int) bool {
	return unicode.IsDigit(cur) && sbuf == 0
}

func isInvalidEscapeSequence(cur, next rune) bool {
	return isSlash(cur) && !isSlash(next) && !unicode.IsDigit(next)
}

func isSlash(r rune) bool {
	return r == unicodeValueBackslash
}
