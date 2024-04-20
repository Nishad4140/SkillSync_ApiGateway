package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"google.golang.org/grpc"
)

func DialGrpc(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}

func CheckString(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, s := range str {
		if unicode.IsNumber(s) {
			return false
		}
	}
	return true
}

func ConvertStringToDate(data string) (time.Time, error) {
	layOut := "02-01-2006"
	date, err := time.Parse(layOut, data)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while converting to time")
	}
	return date, nil
}

func CheckStringNumber(str string) bool {
	_, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	return len(str) == 10
}

func CheckNegativeStringNumber(s string) bool {
	return strings.HasPrefix(s, "-")
}

func CheckNumberInString(s string) bool {
	for _, sr := range s {
		if unicode.IsNumber(sr) {
			return true
		}
	}
	return false
}

func CheckDate(date string) bool {
	layOut := "02-01-2006"
	_, err := time.Parse(layOut, date)
	if err != nil {
		return false
	}
	return true
}

func CheckYear(s string) bool {
	return strings.HasSuffix(s, "years")
}

func CheckNegative(num int32) bool {
	return num < 0
}

func ValidMail(mail string) bool {
	return strings.Contains(mail, "@")
}

func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasLower && hasUpper && hasNumber && hasSpecial
}
