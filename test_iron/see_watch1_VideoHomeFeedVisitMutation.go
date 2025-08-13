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
	hostsee_watch1_VideoHomeFeedVisitMutation     = "graph.facebook.com"
	endpointsee_watch1_VideoHomeFeedVisitMutation = "https://graph.facebook.com/graphql"
)

func Runsee_watch1_VideoHomeFeedVisitMutation(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	//address := hostsee_watch1_VideoHomeFeedVisitMutation + ":443"
	clientTraceID := genUUIDsee_watch1_VideoHomeFeedVisitMutation()
	clientMutationID := genUUIDsee_watch1_VideoHomeFeedVisitMutation()

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

	var token, actorID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&token, &actorID, &userAgent, &netHni, &simHni, &devicegroup)
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

	// utlsConn := utls.UClient(conn, &utls.Config{ServerName: hostsee_watch1_VideoHomeFeedVisitMutation}, utls.HelloAndroid_11_OkHttp)
	// if err := utlsConn.Handshake(); err != nil {
	// 	panic("❌ TLS handshake fail: " + err.Error())
	// }

	variables := map[string]interface{}{
		"input": map[string]string{
			"client_mutation_id": clientMutationID,
			"actor_id":           actorID,
		},
	}
	varBuf, _ := json.Marshal(variables)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "VideoHomeFeedVisitMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "109902923817258724998244020358")
	form.Set("variables", string(varBuf))
	form.Set("fb_api_analytics_tags", `["visitation_id=null","GraphServices"]`)
	form.Set("client_trace_id", clientTraceID)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(form.Encode()))
	gz.Close()

	req, _ := http.NewRequest("POST", endpointsee_watch1_VideoHomeFeedVisitMutation, &buf)
	req.Header.Set("Authorization", "OAuth "+token)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Encoding", "zstd, gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostsee_watch1_VideoHomeFeedVisitMutation)
	req.Header.Set("Priority", "u=3, i")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-token", "1096c31def9028bbc5c6f4f50d7dabe9")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "VideoHomeFeedVisitMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "2444622522461689")
	req.Header.Set("x-fb-qpl-active-flows-json", `{"schema_version":"v2","inprogress_qpls":[{"marker_id":25952257,"annotations":{"current_endpoint":"NewsFeedFragment:native_newsfeed"}}],"snapshot_attributes":{}}`)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-rmd", "state=URL_ELIGIBLE")
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-session-id", "nid=eaVy+PHPWmXj;tid=371;nc=1;fc=2;bc=1;cid=1096c31def9028bbc5c6f4f50d7dabe9")
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+clientTraceID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthsee_watch1_VideoHomeFeedVisitMutation())
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

func genUUIDsee_watch1_VideoHomeFeedVisitMutation() string {
	b := make([]byte, 16)
	crand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func randomExcellentBandwidthsee_watch1_VideoHomeFeedVisitMutation() string {
	mrand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", mrand.Intn(15000000)+20000000)
}
