package main

import (
	"bytes"
	"compress/gzip"
	crand "crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	hostRunsee_watch2_SurveyIntegrationPointQuery     = "graph.facebook.com"
	endpointRunsee_watch2_SurveyIntegrationPointQuery = "https://graph.facebook.com/graphql"
)

func Runsee_watch2_SurveyIntegrationPointQuery(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	//	address := host + ":443"
	clientTraceID := genUUIDRunsee_watch2_SurveyIntegrationPointQuery()

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

	var token, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&token, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ Load profile fail: " + err.Error())
		return
	}

	// proxy := os.Getenv("USE_PROXY")
	// auth := os.Getenv("USE_PROXY_AUTH")

	// conn, err := net.DialTimeout("tcp", proxy, 10*time.Second)
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

	variables := map[string]interface{}{
		"integration_point_id": "420270581758177",
		"survey_context_data": []map[string]string{
			{"context_key": "is_indicator_shown", "context_value": "false"},
			{"context_key": "is_ad_complete", "context_value": "false"},
			{"context_key": "ad_type", "context_value": "NONLIVE"},
			{"context_key": "player_sub_origin", "context_value": "feed"},
			{"context_key": "player_format", "context_value": "inline"},
			{"context_key": "instream_is_skippable", "context_value": "false"},
			{"context_key": "is_ad_shown", "context_value": "false"},
			{"context_key": "is_warion_entry_video", "context_value": "false"},
			{"context_key": "host_video_watch_time_ms", "context_value": "17531"},
			{"context_key": "player_origin", "context_value": "video_home"},
			{"context_key": "video_id", "context_value": "1199647391433909"},
			{"context_key": "integration_point_name", "context_value": "VIDEO_AD_BREAK_STOP_WATCHING"},
			{"context_key": "ad_break_index", "context_value": "-1"},
		},
		"version_number": "2_0_0",
	}
	varBuf, _ := json.Marshal(variables)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "SurveyIntegrationPointQuery")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "48277790717200379631020126574")
	form.Set("variables", string(varBuf))
	form.Set("fb_api_analytics_tags", `["420270581758177","GraphServices"]`)
	form.Set("client_trace_id", clientTraceID)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(form.Encode()))
	gz.Close()

	req, _ := http.NewRequest("POST", endpointRunsee_watch2_SurveyIntegrationPointQuery, &buf)
	req.Header.Set("Authorization", "OAuth "+token)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Encoding", "zstd, gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostRunsee_watch2_SurveyIntegrationPointQuery)
	req.Header.Set("Priority", "u=3, i")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-token", "1096c31def9028bbc5c6f4f50d7dabe9")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "SurveyIntegrationPointQuery")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "c000000000000000")
	req.Header.Set("x-fb-qpl-active-flows-json", `{"schema_version":"v2","inprogress_qpls":[{"marker_id":25952257,"annotations":{"current_endpoint":"WatchFeedOrWarionFragment:video_home_root"}}],"snapshot_attributes":{}}`)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-rmd", "state=URL_ELIGIBLE")
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-session-id", "nid=eaVy+PHPWmXj;tid=378;nc=1;fc=2;bc=1;cid=1096c31def9028bbc5c6f4f50d7dabe9")
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+clientTraceID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthRunsee_watch2_SurveyIntegrationPointQuery())
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

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

func genUUIDRunsee_watch2_SurveyIntegrationPointQuery() string {
	b := make([]byte, 16)
	crand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func randomExcellentBandwidthRunsee_watch2_SurveyIntegrationPointQuery() string {
	mrand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", mrand.Intn(15000000)+20000000)
}
