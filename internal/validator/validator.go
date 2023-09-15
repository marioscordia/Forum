package validator

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors map[string]string
} 

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
} 

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	} 
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
} 

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
} 

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
} 

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
} 

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func ConfirmPassword(value string, value2 string) bool{
	return value == value2
}

func CheckCategory(category []string) bool{
	return len(category) > 0
}

func CheckFileName(image_name string) bool {
	ext := filepath.Ext(image_name)
	if ext != ".jpeg" && ext != ".jpg" && ext != ".gif" && ext != ".png" {
		return false
	}
	return true
}

func CheckFileSize(image_size int) bool {
	return image_size <= 20*1024*1024
}


