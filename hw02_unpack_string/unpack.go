package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")
var escapeRune rune = []rune(`\`)[0]

func getSequenceLen(arr []rune, index int) int {
	elem := arr[index]
	result := 1
	newIndex := index - 1
	for newIndex > 0 {
		if arr[newIndex] == elem {
			result++
		} else {
			return result
		}
		newIndex--
	}
	return result
}

func isEscaped(arr []rune, index int) bool {
	if index <= 0 || index >= len(arr) {
		return false
	}
	return (arr[index-1] == escapeRune && getSequenceLen(arr, index-1)%2 != 0)
}

func isCount(arr []rune, index int) bool {
	if index < 0 || index >= len(arr) {
		return false
	}
	return unicode.IsDigit(arr[index]) && !isEscaped(arr, index)
}

func Unpack(input string) (string, error) {
	var result strings.Builder

	inputAsRuneArr := []rune(input)

	for i, v := range inputAsRuneArr {
		if isEscaped(inputAsRuneArr, i) && !isCount(inputAsRuneArr, i+1) {
			result.WriteRune(v)
		} else if isCount(inputAsRuneArr, i) {
			if (i-1 < 0) || (isCount(inputAsRuneArr, i-1)) {
				return "", ErrInvalidString
			}
			literalStr := string(inputAsRuneArr[i-1])
			digit, _ := strconv.Atoi(string(v))
			result.WriteString(strings.Repeat(literalStr, digit))

		} else if v != escapeRune && !isCount(inputAsRuneArr, i+1) {
			result.WriteRune(v)
		}
	}

	return result.String(), nil
}
