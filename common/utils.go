package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/datatypes"
)

const (
	ERR_PARSE_VALUE_ENV = "cannot parse value of %v env"
)

var (
	FormatErr = func(prefix string, err error) error {
		return ErrorWrapper(fmt.Sprintf(ERR_PARSE_VALUE_ENV, prefix), err)
	}
)

func ConvertStruct2Map(ctx context.Context, obj interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	if obj != nil {
		values := reflect.ValueOf(obj).Elem()
		types := values.Type()

		for i := 0; i < values.NumField(); i++ {
			key := types.Field(i).Tag.Get("filter")
			fmt.Sprint(key)
			value := values.Field(i)

			if key != "" && !value.IsNil() {
				m[key] = value.Interface()
			}
		}
	}

	return m
}

var ConvertMap2StringSQL = func(cond map[string]interface{}) ([]string, []interface{}) {
	sqls := []string{}
	values := []interface{}{}

	for k, v := range cond {
		operator := "="
		if k != "" && v != nil {
			typeValue := fmt.Sprintf("%T", v)
			if strings.Contains(typeValue, "[]") {
				operator = "IN"
			}
			sqls = append(sqls, fmt.Sprintf("%s %s ?", k, operator))

			values = append(values, v)
		}
	}

	return sqls, values
}

type osENV struct {
	name  string
	value string
}

func (o *osENV) ParseInt() (value int64, err error) {
	v, err := strconv.ParseInt(o.value, 10, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseUInt() (value uint64, err error) {
	v, err := strconv.ParseUint(o.value, 10, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseString() (value string, err error) {
	return o.value, nil
}

func (o *osENV) ParseBool() (value bool, err error) {
	v, err := strconv.ParseBool(o.value)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseFloat() (value float64, err error) {
	v, err := strconv.ParseFloat(o.value, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func GetOSEnv(envName string) *osENV {
	value := os.Getenv(envName)
	return &osENV{name: envName, value: value}
}

func GetOffset(page int, pageSize int) int {
	return pageSize * (page - 1)
}

func UnmarshalJSON(input string) (datatypes.JSON, error) {
	var result datatypes.JSON
	if input != "" {
		err := json.Unmarshal([]byte(input), &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func Contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func SendOTPEmail(fromEmail, toEmail, otp string) error {
	apiKey := os.Getenv("MAIL_API_KEY")
	from := mail.NewEmail("MePass", fromEmail)
	to := mail.NewEmail("Recipient Name", toEmail)
	subject := "Your OTP Code to reset password"
	plainTextContent := fmt.Sprintf("Your OTP to reset password is: %s", otp)
	htmlContent := fmt.Sprintf("<strong>Your OTP to reset password is: %s</strong>", otp)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	log.Printf("Email sent! Status Code: %d, Body: %s, Headers: %v", response.StatusCode, response.Body, response.Headers)
	return nil
}

func GenerateRandomOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func NormalizeToBangkokTimezone(t time.Time) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Time{}, errors.New("failed to load timezone")
	}
	return t.In(location), nil
}

const (
	defaultLimit    = 20
	defaultPage     = 1
	defaultPageSize = 10
	maxLimit        = 200
)

func GetPageAndPageSize(page, pageSize int) (int, int) {
	if page == 0 {
		page = defaultPage
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxLimit {
		pageSize = maxLimit
	}
	return page, pageSize
}
