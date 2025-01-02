package services

import (
	"context"
	"os"

	"errors"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
)

type object map[string]interface{}

var (
	app_id int
	key1   string
)

func createOrderZalopay(amount string, description string) (*string, *string, error) {
	rand.Seed(time.Now().UnixNano())
	transID := rand.Intn(1000000) // Generate random trans id
	embedData, _ := json.Marshal(object{
		"redirecturl": os.Getenv("ZALOPAY_REDIRECT_URL"),
	})
	items, _ := json.Marshal([]object{})
	// request data
	params := make(url.Values)
	params.Add("app_id", strconv.Itoa(app_id))
	params.Add("expire_duration_seconds", "300") // 5 minutes
	params.Add("amount", amount)
	params.Add("app_user", "user123")
	params.Add("embed_data", string(embedData))
	params.Add("item", string(items))
	params.Add("description", description)
	params.Add("bank_code", "")

	now := time.Now()
	params.Add("app_time", strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)) // miliseconds

	params.Add("app_trans_id", fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID)) // translation missing: vi.docs.shared.sample_code.comments.app_trans_id

	// appid|app_trans_id|appuser|amount|apptime|embeddata|item
	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

	// Content-Type: application/x-www-form-urlencoded
	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)

	// parse response
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, err
	}

	returnCode, ok := result["return_code"].(float64)
	if !ok || returnCode != 1 {
		return nil, nil, errors.New("Create order failed")
	}

	orderURL := result["order_url"].(string)
	appTransID := params.Get("app_trans_id")

	return &orderURL, &appTransID, nil
}

func checkZalopayOrderStatus(appTransID string) (*int, error) {
	data := fmt.Sprintf("%v|%s|%s", app_id, appTransID, key1) // appid|apptransid|key1
	params := map[string]interface{}{
		"app_id":       app_id,
		"app_trans_id": appTransID,
		"mac":          hmacutil.HexStringEncode(hmacutil.SHA256, key1, data),
	}

	jsonStr, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return_code := int(result["return_code"].(float64))
	return &return_code, nil
}

func (s *Service) UpdateVocabCount(ctx context.Context, userId string) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.VocabUsageCount += 10
	_, err = s.userRepo.UpdateColumns(ctx, user.ID, map[string]interface{}{
		"vocab_usage_count": user.VocabUsageCount,
	})
	if err != nil {
		return err
	}

	return nil
}

// Polling function with a maximum duration of 5 minutes (no max attempts)
func (s *Service) pollOrderStatusForBuyingAiVocabTurn(appTransId string, userID string) {
	interval := 10 * time.Second   // Check every 10 seconds
	maxDuration := 5 * time.Minute // Stop after 5 minutes

	startTime := time.Now()

	for {
		// Check if the polling has exceeded the max duration
		if time.Since(startTime) > maxDuration {
			fmt.Printf("Polling stopped: Max duration of 5 minutes reached for OrderID: %s\n", appTransId)
			return
		}

		// Check order status
		status, _ := checkZalopayOrderStatus(appTransId)
		fmt.Printf("Polling Status: %d\n", *status)

		if *status == 1 {
			_ = s.UpdateVocabCount(context.Background(), userID)
			fmt.Printf("Payment SUCCESS for OrderID: %s, UserID: %s. Updating vocab_count...\n", appTransId, userID)
			return
		} else if *status == 2 {
			fmt.Printf("Payment FAILED for OrderID: %s. Stopping polling.\n", appTransId)
			return
		}

		time.Sleep(interval)
	}
}

func (s *Service) BuyMoreAiVocabTurn(ctx context.Context, userId string) (*string, error) {
	app_id, _ = strconv.Atoi(os.Getenv("ZALOPAY_APP_ID"))
	key1 = os.Getenv("ZALOPAY_KEY1")

	amount := "20000"
	description := "Mua thêm lượt tra từ vựng với AI"
	orderURL, appTransID, err := createOrderZalopay(amount, description)
	if err != nil {
		return nil, err
	}

	go s.pollOrderStatusForBuyingAiVocabTurn(*appTransID, userId)

	return orderURL, nil
}
