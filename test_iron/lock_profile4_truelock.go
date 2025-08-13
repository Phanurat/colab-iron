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

func randomExcellentBandwidthlock_profile4_truelock() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// ------------------------------
var (
	hostlock_profile4_truelock                 = "graph.facebook.com"
	friendlyNamelock_profile4_truelock         = "NativeTemplateAsyncQuery"
	clientDocIDlock_profile4_truelock          = "30749539925783572052554034822"
	privacyMutationTokenlock_profile4_truelock = "eyJ0eXBlIjo1LCJjcmVhdGlvbl90aW1lIjoxNzQ5MjI2OTE5fQ%253D%253D" // <<< แก้เองตามที่ดักมา
)

// ------------------------------

func Runlock_profile4_truelock(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	connToken := generateHex32lock_profile4_truelock()
	sessionUUID := uuid.New().String()
	sessionID := fmt.Sprintf("nid=TSRsHSL+wunc;tid=207;nc=0;fc=0;bc=0;cid=%s", connToken)
	traceID := uuid.New().String()
	taLoggingID := "graphql:" + traceID

	// GEN ใหม่ทุกครั้ง
	payloadSessionID := sessionUUID // ใช้ uuid.New().String() ทุกรอบ

	variables := fmt.Sprintf(
		`{"scale":"3","params":{"payload":"/nt/private_sharing/private_sharing_lock_toggle/?should_revert_past_posts=1&is_eligible_vertical_revert_past_posts=0&new_private_sharing_state=1&entry_point=profile_section&current_page=home&session_id=%s&privacy_mutation_token=%s","nt_context":{"using_white_navbar":true,"pixel_ratio":3,"is_push_on":true,"styles_id":"196702b4d5dfb9dbf1ded6d58ee42767","bloks_version":"c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722"}},"nt_context":{"using_white_navbar":true,"pixel_ratio":3,"is_push_on":true,"styles_id":"196702b4d5dfb9dbf1ded6d58ee42767","bloks_version":"c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722"},"profile_image_size":246,"include_image_ranges":true}`,
		payloadSessionID, privacyMutationTokenlock_profile4_truelock)

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
	form.Set("fb_api_req_friendly_name", friendlyNamelock_profile4_truelock)
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", clientDocIDlock_profile4_truelock)
	form.Set("variables", variables)
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)
	form.Set("client_trace_id", traceID)

	req, _ := http.NewRequest("POST", "https://"+hostlock_profile4_truelock+"/graphql", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "zstd, gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostlock_profile4_truelock)
	req.Header.Set("Priority", "u=3, i")
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-token", connToken)
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", friendlyNamelock_profile4_truelock)
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-qpl-active-flows-json", `{"schema_version":"v2","inprogress_qpls":[{"marker_id":25952257,"annotations":{"current_endpoint":"FbScreenFragment:private_sharing"}}],"snapshot_attributes":{}}`)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-session-id", sessionID)
	req.Header.Set("x-fb-ta-logging-ids", taLoggingID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)                                                        // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                                        // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthlock_profile4_truelock()) //เพิ่มเข้าไป
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

func generateHex32lock_profile4_truelock() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
