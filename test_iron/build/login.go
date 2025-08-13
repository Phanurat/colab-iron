package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"strings"
	"time"

	"github.com/google/uuid"
)

func main() {
	email := "ottzenxhumra@hotmail.com"
	password := "PhanTuyet8m4mKC"
	token, err := getToken(email, password)
	if err != nil {
		fmt.Println("❌ ERROR:", err)
		return
	}
	fmt.Println("✅ RESULT:", token)
}

func getToken(email, password string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	sim := randBetween(20000, 40000)
	deviceID := uuid.New().String()
	adID := uuid.New().String()

	form := url.Values{
		"adid":                       {adID},
		"format":                     {"json"},
		"device_id":                  {deviceID},
		"email":                      {email},
		"password":                   {password},
		"cpl":                        {"true"},
		"family_device_id":           {deviceID},
		"credentials_type":           {"device_based_login_password"},
		"generate_session_cookies":   {"1"},
		"error_detail_type":          {"button_with_disabled"},
		"source":                     {"device_based_login"},
		"machine_id":                 {randString(24)},
		"meta_inf_fbmeta":            {""},
		"advertiser_id":              {adID},
		"currently_logged_in_userid": {"0"},
		"locale":                     {"en_US"},
		"client_country_code":        {"US"},
		"method":                     {"auth.login"},
		"fb_api_req_friendly_name":   {"authenticate"},
		"fb_api_caller_class":        {"com.facebook.account.login.protocol.Fb4aAuthHandler"},
		"api_key":                    {"882a8490361da98702bf97a021ddc14d"},
	}

	// สร้าง sig แล้วใส่เข้าไป
	form["sig"] = []string{getSig(form)}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://b-api.facebook.com/method/auth.login", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("x-fb-connection-bandwidth", fmt.Sprintf("%d", randBetween(20000000, 30000000)))
	req.Header.Set("x-fb-sim-hni", fmt.Sprintf("%d", sim))
	req.Header.Set("x-fb-net-hni", fmt.Sprintf("%d", sim))
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")
	req.Header.Set("x-fb-connection-type", "cell.CTRadioAccessTechnologyHSDPA")
	req.Header.Set("user-agent", "Dalvik/1.6.0 (Linux; U; Android 4.4.2; NX55 Build/KOT5506) [FBAN/FB4A;FBAV/106.0.0.26.68;FBBV/45904160;FBDM/{density=3.0,width=1080,height=1920};FBLC/it_IT;FBRV/45904160;FBCR/PosteMobile;FBMF/asus;FBBD/asus;FBPN/com.facebook.katana;FBDV/ASUS_Z00AD;FBSV/5.0;FBOP/1;FBCA/x86:armeabi-v7a;]")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("x-fb-http-engine", "Liger")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func getSig(values url.Values) string {
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}
	// Sort keys alphabetically
	sortStrings(keys)

	var sigBuilder strings.Builder
	for _, key := range keys {
		sigBuilder.WriteString(fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	sigBuilder.WriteString("62f8ce9f74b12f84c123cc23437a4a32") // secret key

	sum := md5.Sum([]byte(sigBuilder.String()))
	return hex.EncodeToString(sum[:])
}

func randBetween(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func sortStrings(arr []string) {
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}
