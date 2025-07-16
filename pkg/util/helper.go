package util

import (
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/enums"
	"math"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
)

type Identifiable interface {
	GetId() string
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Paginate(r *http.Request) (sk int64, lim int64, er error) {
	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("perPage")

	var skip, limit int64
	var convErr error
	if perPage == "" {
		limit = 40
	} else {
		limit, convErr = strconv.ParseInt(perPage, 10, 64)
		if convErr != nil {
			return 0, 0, convErr
		}
	}
	if page == "" {
		skip = 0
	} else {
		skip, convErr = strconv.ParseInt(page, 10, 64)
		if convErr != nil {
			return 0, 0, convErr
		}
		skip = (skip - 1) * limit
	}

	return skip, limit, nil
}

func PaginatePage(r *http.Request) (int, int) {
	perPage := r.URL.Query().Get("perPage")
	if perPage == "" {
		perPage = "10"
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		return 0, 0
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return 0, 0
	}
	return perPageInt, pageInt
}

func CacheKeyForMulti(key enums.CacheKey, db, path, rawQuery string) string {
	pattern := "%s:p:c:%s:%s?%s"
	return fmt.Sprintf(pattern, key, db, path, rawQuery)
}

func CacheKeyForOne(key enums.CacheKey, db, id string) string {
	pattern := "%s:g:c:%s:%s"
	return fmt.Sprintf(pattern, key, db, id)
}

func CacheKeyDeleteMul(key enums.CacheKey, db string) string {
	pattern := "%s:p:c:%s:*"
	return fmt.Sprintf(pattern, key, db)
}
func CacheKeyDeleteOne(key enums.CacheKey, db, id string) string {
	pattern := "%s:g:c:%s:%s"
	return fmt.Sprintf(pattern, key, db, id)
}

func InSlice[T comparable](ele T, s []T) bool {
	for _, e := range s {
		if e == ele {
			return true
		}
	}

	return false

}

func GetBsonFromString(query string) (bson.M, error) {
	var bs bson.M = bson.M{}
	if query != "" {
		err := bson.UnmarshalExtJSON([]byte(query), true, &bs)
		if err != nil {
			return nil, err
		}
	}

	return bs, nil
}

func GeneratePassword(passWordLength, minSpecialChar, minUpper, minNum int) string {
	var (
		lowerCharSet   = "abcdedfghijklmnopqrst"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialCharSet = "!@#$%&*"
		numberSet      = "0123456789"
		allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	)
	var password strings.Builder
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < minSpecialChar; i++ {
		index := r.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[index]))
	}

	for i := 0; i < minUpper; i++ {
		index := r.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[index]))
	}

	for i := 0; i < minNum; i++ {
		index := r.Intn(len(numberSet))
		password.WriteString(string(numberSet[index]))
	}

	remainLength := passWordLength - minSpecialChar - minUpper - minNum

	for i := 0; i < remainLength; i++ {
		index := r.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[index]))
	}

	return password.String()

}

// as is date or time
func GetDateAs(date time.Time, as string) time.Time {
	if as == "date" {
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	} else if as == "time" {
		return time.Date(0, 0, 0, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), time.UTC)
	}

	return date

}

type ReqArgs struct {
	Sort    bson.M
	GroupBy bson.M
	//any other args needed
}

func (a *ReqArgs) SetSort(r *http.Request) error {
	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")

	if sortBy != "" && sortOrder != "" {
		sort := bson.M{}
		order, convErr := strconv.ParseInt(sortOrder, 10, 64)
		if convErr != nil {
			return convErr
		}
		sort[sortBy] = order
		a.Sort = sort
	}

	return nil
}

func (a *ReqArgs) SetGroupBy(r *http.Request) error {
	groupBy := r.URL.Query().Get("groupBy")

	if groupBy != "" {
		group := bson.M{
			"_id":     "$" + groupBy,
			"records": bson.M{"$push": "$$ROOT"},
		}

		a.GroupBy = group
	}

	return nil

}

func GetWeekdayNames(startDate, endDate time.Time) []string {
	var weekdays []string

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		weekdays = append(weekdays, d.Weekday().String())
	}

	return weekdays
}

func BsonAToSlice[T any](bsonArray bson.A) ([]T, error) {
	var result []T
	for _, ele := range bsonArray {
		t, ok := ele.(T)
		if !ok {
			return nil, errors.New("error type")
		}

		result = append(result, t)
	}
	return result, nil
}

func Contains[T comparable](arr []T, key T) bool {
	for _, value := range arr {
		if value == key {
			return true
		}
	}
	return false
}

func SliceFilter[T any](arr []T, callBack func(el T) bool) []T {
	var result []T
	for _, e := range arr {
		ok := callBack(e)
		if ok {
			result = append(result, e)
		}
	}

	return result
}

func SliceFind[T any](arr []T, callback func(el T) bool) (T, bool) {
	var result T
	var ok bool
	for _, e := range arr {
		ok = callback(e)
		if ok {
			result = e
			break
		}
	}

	return result, ok
}

func RoundToDecimalPlaces(value float64, decimalPlaces int) float64 {
	factor := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*factor) / factor
}

func IsValidEmail(email string) bool {
	// Define the regex pattern for a valid email address
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	re := regexp.MustCompile(emailPattern)

	// Match the email address against the pattern
	return re.MatchString(email)
}

func GenerateUniqueString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func HasDuplicateIDs[T Identifiable](slice []T) bool {
	seen := make(map[string]bool)
	for _, item := range slice {
		if seen[item.GetId()] {
			return true
		}
		seen[item.GetId()] = true
	}
	return false
}

func SliceReorder[T Identifiable](s []T, oldIndex, newIndex int) ([]T, error) {
	if oldIndex >= 0 && oldIndex < len(s) && newIndex >= 0 && newIndex < len(s) {
		// Remove the element from the old position
		ele := s[oldIndex]
		s = append(s[:oldIndex], s[oldIndex+1:]...)
		// Insert the element at the new position
		s = append(s[:newIndex], append([]T{ele}, s[newIndex:]...)...)
		return s, nil
	} else {
		return s, errors.New("invalid indices")
	}
}

func InsertAt[T any](slice []T, index int, element T) ([]T, error) {
	// Validate the index
	if index < 0 || index > len(slice) {
		return nil, fmt.Errorf("index out of range: %d", index)
	}

	// Insert the element at the specified index
	newSlice := append(slice[:index], append([]T{element}, slice[index:]...)...)

	return newSlice, nil
}

// StructToMap converts a struct to a map[string]any
func StructToMap(obj any) (map[string]any, error) {
	result := make(map[string]any)

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct, got %v", v.Kind())
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get the JSON tag name, fallback to field name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue // Skip fields without JSON tags
		}

		// Handle comma-separated tags (e.g., "name,omitempty")
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonTag = jsonTag[:idx]
		}

		if jsonTag == "" {
			jsonTag = fieldType.Name
		}

		// Convert field value
		value := field.Interface()
		result[jsonTag] = value
	}

	return result, nil
}

// input: "syria"
// output: "Syria"
func CapitalizeFirstLowerRest(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
