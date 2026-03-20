package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const NotifyServiceURL = "http://localhost:3008/notify"

// --- Structs cập nhật chính xác theo logic data.total ---

type Step1Response struct {
	Status int `json:"status"`
	Data   struct {
		CustomToken string `json:"customToken"`
	} `json:"data"`
}

type Step2Response struct {
	IdToken string `json:"idToken"`
}

type TripResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Total int `json:"total"` // Thông số quan trọng bạn yêu cầu
		List  []struct {
			Id            string `json:"id"`
			DepartureTime string `json:"departureTime"`
			Price         int    `json:"price"`
			AvailableSeat int    `json:"availableSeat"`
		} `json:"data"` // Một số API để trong field 'data' hoặc 'list' tùy version
	} `json:"data"`
}

func main() {
	client := &http.Client{Timeout: 30 * time.Second}
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	fmt.Println("🚀 Worker đang chạy... Kiểm tra data.total mỗi 5 phút.")
	executeTask(client)

	for range ticker.C {
		executeTask(client)
	}
}

func executeTask(client *http.Client) {
	fmt.Printf("\n==================== KIỂM TRA LÚC %s ====================\n", time.Now().Format("15:04:05"))

	// BƯỚC 1: LOGIN
	loginReq := map[string]string{
		"grant_type": "PASSWORD",
		"provider":   "PHONE",
		"username":   "0333771800",
		"password":   "Th@nhH@i3303703",
	}
	b1, err := callAPI(client, "POST", "https://api.vato.vn/api/authenticate/login", nil, loginReq)
	if err != nil {
		return
	}
	var res1 Step1Response
	json.Unmarshal(b1, &res1)

	// BƯỚC 2: VERIFY
	verifyReq := map[string]interface{}{"token": res1.Data.CustomToken, "returnSecureToken": true}
	b2, err := callAPI(client, "POST", "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=AIzaSyCmNFcFRBwaOa2cVPTjwzh9mHAN7sKymd4", nil, verifyReq)
	if err != nil {
		return
	}
	var res2 Step2Response
	json.Unmarshal(b2, &res2)

	// BƯỚC 3: CHECK TRIPS (KIỂM TRA TOTAL)
	tripReq := map[string]interface{}{
		"channel":    "web_client",
		"routeIds":   []string{"11f09e78-0e39-1316-a90c-42964ca74856", "d782c138-75b8-40b5-985b-b4bc3a5b9ed1"},
		"fromDate":   "2026-04-27T17:00:00.000Z",
		"toDate":     "2026-04-28T16:59:59.999Z",
		"minNumSeat": 1,
	}
	headers := map[string]string{
		"x-access-token": res2.IdToken,
		"x-app-id":       "client",
		"origin":         "https://futabus.vn",
	}
	b3, err := callAPI(client, "POST", "https://api-online.futabus.vn/vato/v1/search/trip-by-route", headers, tripReq)
	if err != nil {
		return
	}

	var res3 TripResponse
	json.Unmarshal(b3, &res3)

	// KIỂM TRA THÔNG SỐ TOTAL
	totalTrips := res3.Data.Total
	fmt.Printf("[BƯỚC 3] Kết quả: data.total = %d\n", totalTrips)

	if totalTrips > 0 {
		fmt.Printf("🎯 THÀNH CÔNG: Tìm thấy %d chuyến xe!\n", totalTrips)
		msg := fmt.Sprintf("🎫 CÓ VÉ! Hệ thống Futa báo có %d chuyến xe mới.", totalTrips)
		callAPI(client, "POST", NotifyServiceURL, nil, map[string]string{"message": msg})
	} else {
		fmt.Println("ℹ️ Hiện tại data.total = 0 (Chưa có vé).")
	}
}

// Hàm callAPI tích hợp log response đầy đủ khi cần debug
func callAPI(client *http.Client, method, url string, headers map[string]string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		js, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(js)
	}

	req, _ := http.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Lỗi mạng: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("[%d] %s\n", resp.StatusCode, url)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("⚠️ Chi tiết lỗi Body: %s\n", string(respBody))
		return respBody, fmt.Errorf("Status %d", resp.StatusCode)
	}

	// Bạn có thể uncomment dòng dưới đây nếu muốn soi toàn bộ JSON trả về khi chạy tốt
	// fmt.Printf("DEBUG JSON: %s\n", string(respBody))

	return respBody, nil
}
