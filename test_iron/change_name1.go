package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func randomExcellentBandwidthchange_name1() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func generateRandomPrivacyContextchange_name1() string {
	return strconv.FormatInt(rand.Int63n(899999999999999)+100000000000000, 10)
}

func Runchange_name1(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	rand.Seed(time.Now().UnixNano())

	host := "graph.facebook.com"

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "fb_comment_system.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
		return
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/fb_comment_system.db")

	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	privacyContext := generateRandomPrivacyContextchange_name1()

	contextData := map[string]interface{}{
		"using_white_navbar": true,
		"pixel_ratio":        3,
		"is_push_on":         true,
		"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
		"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
	}

	variables := map[string]interface{}{
		"context":    contextData,
		"nt_context": contextData,
		"scale":      "3",
	}
	variablesJson, _ := json.Marshal(variables)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "PrivacySettingsNTActionQuery")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "101861545817786939371120918504")
	form.Set("variables", string(variablesJson))
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(form.Encode()))
	gz.Close()

	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", &buf)
	req.Header = http.Header{
		"Authorization":               {"OAuth " + accessToken},
		"Content-Encoding":            {"gzip"},
		"Content-Type":                {"application/x-www-form-urlencoded"},
		"User-Agent":                  {userAgent},
		"X-FB-Friendly-Name":          {"PrivacySettingsNTActionQuery"},
		"X-FB-Connection-Type":        {"MOBILE.HSDPA"},
		"X-FB-HTTP-Engine":            {"Liger"},
		"X-FB-Background-State":       {"1"},
		"x-fb-client-ip":              {"True"},
		"x-fb-device-group":           {devicegroup},
		"x-fb-privacy-context":        {privacyContext},
		"X-FB-Request-Analytics-Tags": {`{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`},
		"x-fb-server-cluster":         {"True"},
		"x-graphql-client-library":    {"graphservice"},
		"x-tigon-is-retry":            {"False"},
	}

	// ---------- SEND ----------
	// ---------- SEND ----------
	bw := tlsConns.RWGraph.Writer
	br := tlsConns.RWGraph.Reader

	err = req.Write(bw)
	if err != nil {
		fmt.Println("❌ Write fail: " + err.Error())
		return

	}
	bw.Flush() // ✅ ต้อง flush เพื่อให้ข้อมูลถูกส่งออกจริง ๆ

	// ✅ ใช้ reader ตัวเดียวกับที่รับมาจาก utls
	resp, err := http.ReadResponse(br, req)
	if err != nil {
		fmt.Println("❌ Read fail: " + err.Error())
		return

	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("❌ GZIP decompress fail: " + err.Error())
			return

		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}

	bodyResp, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("❌ Body read fail: " + err.Error())
		return

	}

	fmt.Println("✅ Status:", resp.Status)
	fmt.Println("📦 Response:", string(bodyResp))

}
