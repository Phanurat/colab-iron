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

func buildEncodedPayloadRunlike_page2_FbBloksActionRootQuery() string {
	serverParams := map[string]interface{}{
		"INTERNAL__latency_qpl_marker_id":   36707139,
		"INTERNAL__latency_qpl_instance_id": generateQPLIDRunlike_page2_FbBloksActionRootQuery(),
	}
	params := map[string]interface{}{
		"params": map[string]interface{}{
			"params": map[string]interface{}{
				"client_input_params": map[string]interface{}{},
				"server_params":       serverParams,
			},
			"bloks_versioning_id": "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
			"app_id":              "com.bloks.www.fb.profile.recommendation.down.caret.action.tap.expose.chaining.backtest",
		},
		"scale": "3",
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
	}
	jsonVars, _ := json.Marshal(params)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("purpose", "fetch")
	form.Set("fb_api_req_friendly_name", "FbBloksActionRootQuery-com.bloks.www.fb.profile.recommendation.down.caret.action.tap.expose.chaining.backtest")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "11994080424240083948543644217")
	form.Set("variables", string(jsonVars))
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)
	form.Set("client_trace_id", uuid.New().String())
	return form.Encode()
}

func generateQPLIDRunlike_page2_FbBloksActionRootQuery() float64 {
	rand.Seed(time.Now().UnixNano())
	return float64(rand.Int63n(899999999999999) + 100000000000000)
}

func generateSessionIDRunlike_page2_FbBloksActionRootQuery() string {
	u := uuid.New().String()
	tid := strconv.Itoa(rand.Intn(900) + 100)
	return fmt.Sprintf("nid=%s;tid=%s;nc=1;fc=2;bc=1;cid=%s", u, tid, uuid.New().String())
}

func gzipCompressRunlike_page2_FbBloksActionRootQuery(data []byte) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(data)
	writer.Close()
	return buf.Bytes()
}

func Runlike_page2_FbBloksActionRootQuery(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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
	var token, userAgent, netHni, simHni, deviceGroup string
	err = db.QueryRow(`SELECT access_token, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1`).Scan(
		&token, &userAgent, &netHni, &simHni, &deviceGroup)
	if err != nil {
		fmt.Println("❌ ดึง app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	payload := buildEncodedPayloadRunlike_page2_FbBloksActionRootQuery()
	host := "graph.facebook.com"

	req, err := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBuffer(gzipCompressRunlike_page2_FbBloksActionRootQuery([]byte(payload))))
	if err != nil {
		fmt.Println("❌ NewRequest fail: " + err.Error())
		return
	}

	traceID := uuid.New().String()
	sessionID := generateSessionIDRunlike_page2_FbBloksActionRootQuery()
	connectionToken := uuid.New().String()

	req.Header = map[string][]string{
		"Authorization":               {"OAuth " + token},
		"Accept-Encoding":             {"zstd, gzip, deflate"},
		"Connection":                  {"keep-alive"},
		"Content-Encoding":            {"gzip"},
		"Content-Type":                {"application/x-www-form-urlencoded"},
		"Host":                        {host},
		"User-Agent":                  {userAgent},
		"X-FB-Friendly-Name":          {"FbBloksActionRootQuery-com.bloks.www.fb.profile.recommendation.down.caret.action.tap.expose.chaining.backtest"},
		"x-fb-client-ip":              {"True"},
		"x-fb-connection-token":       {connectionToken},
		"X-FB-Connection-Type":        {"MOBILE.HSDPA"},
		"x-fb-device-group":           {deviceGroup},
		"X-FB-HTTP-Engine":            {"Liger"},
		"x-fb-privacy-context":        {"3643298472347298"},
		"x-fb-qpl-active-flows-json":  {`{"schema_version":"v2","inprogress_qpls":[],"snapshot_attributes":{}}`},
		"X-FB-Request-Analytics-Tags": {`{"network_tags":{"product":"350685531728","purpose":"fetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`},
		"x-fb-rmd":                    {"state=URL_ELIGIBLE"},
		"x-fb-server-cluster":         {"True"},
		"x-fb-session-id":             {sessionID},
		"x-fb-ta-logging-ids":         {"graphql:" + traceID},
		"x-graphql-client-library":    {"graphservice"},
		"x-graphql-request-purpose":   {"fetch"},
		"x-tigon-is-retry":            {"False"},
	}

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
