package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
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

func randomExcellentBandwidthlock_profile1() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// ---------- ปรับค่านี้ได้เอง (ตามที่ดักมา) ----------
var (
	stylesIDlock_profile1     = "196702b4d5dfb9dbf1ded6d58ee42767"
	bloksVersionlock_profile1 = "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722"
	pathlock_profile1         = "/private_sharing/new_home_screen/"
	clientDocIDlock_profile1  = "22108083522995106186536005950"
	fbFriendlylock_profile1   = "NativeTemplateScreenQuery"
	fbPurposelock_profile1    = "fetch"
)

// ------------------------------------------------------

func Runlock_profile1(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	host := "graph.facebook.com"

	connToken := generateConnectionTokenlock_profile1()
	sessionID := generateSessionIDlock_profile1(connToken)
	traceID := uuid.New().String()
	qplActive := `{"schema_version":"v2","inprogress_qpls":[{"marker_id":25952257,"annotations":{"current_endpoint":"ProfileDynamicActionBarOverflowActivity"}}],"snapshot_attributes":{}}`

	variables := fmt.Sprintf(`{"params":{"nt_context":{"using_white_navbar":true,"pixel_ratio":3,"is_push_on":true,"styles_id":"%s","bloks_version":"%s"},"path":"%s","params":"{\"entry_point\":\"profile_section\"}","extra_client_data":{}},"scale":"3","nt_context":{"using_white_navbar":true,"pixel_ratio":3,"is_push_on":true,"styles_id":"%s","bloks_version":"%s"}}`, stylesIDlock_profile1, bloksVersionlock_profile1, pathlock_profile1, stylesIDlock_profile1, bloksVersionlock_profile1)

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

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("purpose", fbPurposelock_profile1)
	form.Set("fb_api_req_friendly_name", fbFriendlylock_profile1)
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", clientDocIDlock_profile1)
	form.Set("variables", variables)
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)
	form.Set("client_trace_id", traceID)

	// GZIP compress
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	_, _ = gz.Write([]byte(form.Encode()))
	_ = gz.Close()

	req, _ := http.NewRequest("POST", "https://graph.facebook.com/graphql", &compressed)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "zstd, gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-token", connToken)
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", fbFriendlylock_profile1)
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-qpl-active-flows-json", qplActive)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"fetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-session-id", sessionID)
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+traceID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-graphql-request-purpose", fbPurposelock_profile1)
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)                                               // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                               // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthlock_profile1()) //เพิ่มเข้าไป
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

}

// --- สำหรับเจนค่าที่ต้องเปลี่ยนใหม่ทุกรอบ ---
func generateConnectionTokenlock_profile1() string {
	// 32 char, hex lower
	const chars = "abcdef0123456789"
	rnd := make([]byte, 32)
	for i := range rnd {
		rnd[i] = chars[rand.Intn(len(chars))]
	}
	return string(rnd)
}

func generateSessionIDlock_profile1(cid string) string {
	return fmt.Sprintf("nid=TSRsHSL+wunc;tid=204;nc=0;fc=0;bc=0;cid=%s", cid)
}
