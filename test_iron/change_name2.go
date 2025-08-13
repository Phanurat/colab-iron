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

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func randomExcellentBandwidthchange_name2() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runchange_name2(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

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

	var firstName, lastName string
	err = db.QueryRow("SELECT first_Name, last_Name FROM change_name_table LIMIT 1").Scan(
		&firstName, &lastName)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล change_name_table ไม่สำเร็จ: " + err.Error())
		return
	}

	//firstName := "Mind"
	//lastName := "Sukanya"
	fullName := firstName + " " + lastName
	//uid := "61576848086508"

	familyDeviceID := uuid.New().String()
	qplInstanceID := fmt.Sprintf("%.0f", float64(rand.Int63n(9e12)+1e12))
	privacyContext := randomNumericStringchange_name2(16)

	clientInput := map[string]interface{}{
		"client_input_params": map[string]interface{}{
			"first_name":       firstName,
			"last_name":        lastName,
			"full_name":        fullName,
			"middle_name":      "",
			"family_device_id": familyDeviceID,
		},
		"server_params": map[string]interface{}{
			"identity_identifiers":                  []string{},
			"INTERNAL__latency_qpl_marker_id":       36707139,
			"identity_ids_DEPRECATED":               userID,
			"should_dismiss_screen_after-save":      0,
			"is_meta_verified_profile_editing_flow": 0,
			"requested_screen_component_type":       2,
			"machine_id":                            nil,
			"INTERNAL__latency_qpl_instance_id":     qplInstanceID,
		},
	}
	l1, _ := json.Marshal(clientInput)
	l2, _ := json.Marshal(map[string]string{"params": string(l1)})
	l3, _ := json.Marshal(map[string]interface{}{
		"params":              string(l2),
		"bloks_versioning_id": "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		"app_id":              "com.bloks.www.fxim.settings.name.change.async",
	})
	finalVars := map[string]interface{}{
		"params": json.RawMessage(l3),
		"scale":  "3",
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
	}
	varBuf, _ := json.Marshal(finalVars)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("purpose", "fetch")
	form.Set("fb_api_req_friendly_name", "FbBloksActionRootQuery-com.bloks.www.fxim.settings.name.change.async")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "11994080424240083948543644217")
	form.Set("variables", string(varBuf))
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)

	var zipped bytes.Buffer
	gw := gzip.NewWriter(&zipped)
	gw.Write([]byte(form.Encode()))
	gw.Close()

	host := "graph.facebook.com"

	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", &zipped)
	req.Host = host
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Friendly-Name", "FbBloksActionRootQuery-com.bloks.www.fxim.settings.name.change.async")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("x-fb-privacy-context", privacyContext)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"fetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-graphql-request-purpose", "fetch")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthchange_name2())
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

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

	_, err = db.Exec("DELETE FROM change_name_table WHERE first_Name = ?", firstName) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ change_name_table ออกจากฐานข้อมูลแล้ว:", firstName)
	}

	_, err = db.Exec("DELETE FROM change_name_table WHERE last_Name = ?", lastName) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ change_name_table ออกจากฐานข้อมูลแล้ว:", lastName)
	}

}

func randomNumericStringchange_name2(length int) string {
	digits := "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}
