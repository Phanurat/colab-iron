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

// ------------- ตั้งค่าตัวแปรแก้ไขได้ตรงนี้ -----------------

var (
	//deviceGroup = "6301"

	host         = "graph.facebook.com"
	friendlyName = "FbStoriesLightBucketsQuery"
	clientDocID  = "134400997213403767931215154379"
	privacyCtx   = "908202689602531"
)

func randomExcellentBandwidth() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000 // 20 Mbps
	max := 35000000 // 35 Mbps
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// ------------- MAIN --------------

func Runopenapp1(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	rand.Seed(time.Now().UnixNano())
	connToken := generateHex32()
	sessionID := fmt.Sprintf("nid=TSRsHSL+wunc;tid=%d;nc=0;fc=0;bc=0;cid=%s", rand.Intn(300)+100, connToken)
	traceID := uuid.New().String()
	clientReqID := uuid.New().String()
	taLoggingID := "graphql:" + traceID

	momentBucket := randomFbID() // *** เจนใหม่ทุกครั้ง ***

	variables := fmt.Sprintf(
		`{"nt_context":{"using_white_navbar":true,"pixel_ratio":3,"is_push_on":true,"styles_id":"196702b4d5dfb9dbf1ded6d58ee42767","bloks_version":"c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722"},"video_thumbnail_width":672,"use_server_thumbnail":true,"height":1920,"video_thumbnail_height":672,"restore_all_removed_fields_wave1":true,"profile_image_size":105,"should_include_first_media":true,"query_trigger":"COLD_START","width":1080,"moment_bucket_id":["%s"],"fbstory_tray_sizing_type":"cover-fill-cropped","page_profile_image_size_experiment":105,"fbstory_tray_preview_width":318,"scale":"3","fbstory_tray_preview_height":565,"enable_cix_screen_rollout":true,"client_request_id":"%s","unified_stories_buckets_paginated_for_light_query_first":6,"should_include_expiration_time":true,"min_story_page_size_from_client":4}`,
		momentBucket, clientReqID,
	)

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
	form.Set("fb_api_req_friendly_name", friendlyName)
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", clientDocID)
	form.Set("variables", variables)
	form.Set("fb_api_analytics_tags", `["type=head_load","GraphServices","At_Connection","surface=story_tray","trigger=cold_start"]`)
	form.Set("client_trace_id", traceID)

	// address := host + ":443"
	// var conn net.Conn

	// proxy := os.Getenv("USE_PROXY")
	// auth := os.Getenv("USE_PROXY_AUTH")

	// conn, err = net.DialTimeout("tcp", proxy, 10*time.Second)
	// if err != nil {
	// 	panic("❌ Proxy fail: " + err.Error())
	// }

	// reqLine := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n", address, host)
	// if auth != "" {
	// 	reqLine += "Proxy-Authorization: Basic " + auth + "\r\n"
	// }
	// reqLine += "\r\n"
	// fmt.Fprintf(conn, reqLine)

	// br := bufio.NewReader(conn)
	// respLine, _ := br.ReadString('\n')
	// if !strings.Contains(respLine, "200") {
	// 	panic("❌ CONNECT fail: " + respLine)
	// }
	// for {
	// 	line, _ := br.ReadString('\n')
	// 	if line == "\r\n" || line == "" {
	// 		break
	// 	}
	// }

	// utlsConn := utls.UClient(conn, &utls.Config{ServerName: host}, utls.HelloAndroid_11_OkHttp)
	// if err := utlsConn.Handshake(); err != nil {
	// 	panic("❌ TLS handshake fail: " + err.Error())
	// }

	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	// req.Header.Set("X-FB-Connection-Type", "WIFI")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-connection-token", connToken)
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", friendlyName)
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-session-id", sessionID)
	req.Header.Set("x-fb-ta-logging-ids", taLoggingID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-privacy-context", privacyCtx)

	req.Header.Set("x-fb-net-hni", netHni)                                  // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                  // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidth()) //เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

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

func generateHex32() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// ******* เจน random FB 15-17 หลัก ********
func randomFbID() string {
	digits := 15 + rand.Intn(3)
	out := make([]byte, digits)
	for i := range out {
		out[i] = byte('0' + rand.Intn(10))
	}
	return string(out)
}
