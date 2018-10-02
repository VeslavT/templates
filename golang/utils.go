package golang

import (
	"net/http"
	"io/ioutil"
	"io"
	"strings"
	"fmt"
	"reflect"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// Convert to map[string]interface{}
func StructToMap(st interface{}) (map[string]interface{}) {
	dataMap := make(map[string]interface{})
	iVal := reflect.ValueOf(st).Elem()
	iType := iVal.Type()

	for i := 0; i< iVal.NumField(); i++ {
		f := iVal.Field(i)
		dataMap[iType.Field(i).Name] = f
	}

	return dataMap
}

func ReadResponseData(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

func ReadBodyData(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}

func ReverseString(s string) string {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	return string(runes)
}

func StringInSlice(a string, list []string) bool {
	for i := range list {
		if a == list[i] {
			return true
		}
	}

	return false
}

// StringSliceDiff compares a against b and returns the values in a that are not present in b
func StringSliceDiff(a, b []string) []string {
	mb := map[string]bool{}
	for i := range b {
		mb[b[i]] = true
	}

	diff := make([]string, 0, len(a))
	for i := range a {
		if _, ok := mb[a[i]]; !ok {
			diff = append(diff, a[i])
		}
	}
	return diff
}

func ExtendMap(res map[string]bool, slc []string) {
	for _, val := range slc {
		res[val] = true
	}
}

func ExtendMapList(res map[string]map[string]bool, slc map[string][]string) {
	for key, vals := range slc {
		if len(vals) == 0 {
			continue
		}
		for _, val := range vals {
			if _, ok := res[val]; !ok {
				res[val] = make(map[string]bool)
			}
			res[val][key] = true
		}
	}
}

func InStringSliceEqual(str string, haystack []string) bool {
	for _, row := range haystack {
		if str == row {
			return true
		}
	}
	return false
}

func RetrieveCommandFromBody(str string) string {
	trimmed := strings.TrimSpace(str)
	index := strings.Index(trimmed, " ")
	if index == -1 {
		return trimmed
	}
	return trimmed[:index]
}

func EqualSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func ExcludeEmptyStrings(lst []string) []string {
	slc := make([]string, 0, len(lst))
	for _, val := range lst {
		if val != "" {
			slc = append(slc, val)
		}
	}
	return slc
}

func ConcatList(values []string) string {
	if len(values) == 0 {
		return "[]"
	} else {
		return fmt.Sprintf(`["%s"]`, strings.Join(values, `","`))
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}