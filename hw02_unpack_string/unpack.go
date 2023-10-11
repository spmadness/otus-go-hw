package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

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
		if isDigit(cur) && sbuf == 0 {
			return "", ErrInvalidString
		}
		// если i еще попадает в слайс рун, а j - уже нет
		if j == len(r) {
			if !isSlash(cur) || sbuf > 0 {
				b.WriteString(string(cur))
			}
			break
		}

		next = r[j]

		// из условия задачи - экранировать можно только цифру или слэш
		if isSlash(cur) && !isSlash(next) && !isDigit(next) {
			return "", ErrInvalidString
		}

		if isSlash(cur) && (isSlash(next) || isDigit(next) && sbuf == 0) {
			sbuf++
			continue
		}

		// исходя из логики дефолтных тест-кейсов:
		// если слэшей в буфере больше одного, тогда пишем текущий слэш и следующую цифру как одиночный символ
		// если слэш в буфере один - значит текущий слэш можно обработать с множителем, и мы идем дальше
		if isSlash(cur) && isDigit(next) && sbuf > 1 {
			b.WriteString(string(cur))
			b.WriteString(string(next))
			// после записей, переход на позицию после записанной одиночной цифры
			i++
			j++
			sbuf = 0
			continue
		}

		// запись повторяющихся символов
		if isDigit(next) {
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

func isSlash(r rune) bool {
	return r == 92 // значение слэша в unicode
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}
