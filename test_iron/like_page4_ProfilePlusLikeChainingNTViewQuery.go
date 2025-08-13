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

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func buildEncodedPayloadRunlike_page4_ProfilePlusLikeChainingNTViewQuery(profileID string) string {
	variables := map[string]interface{}{
		"profile_id": profileID,
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
		"scale": "3",
	}
	jsonVars, _ := json.Marshal(variables)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "ProfilePlusLikeChainingNTViewQuery")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "27407771627356531225900593843")
	form.Set("variables", string(jsonVars))
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)
	form.Set("client_trace_id", uuid.New().String())
	return form.Encode()
}

func generateSessionIDRunlike_page4_ProfilePlusLikeChainingNTViewQuery() string {
	u := uuid.New().String()
	tid := strconv.Itoa(rand.Intn(900) + 100)
	return fmt.Sprintf("nid=%s;tid=%s;nc=1;fc=2;bc=1;cid=%s", u, tid, uuid.New().String())
}

func gzipCompressRunlike_page4_ProfilePlusLikeChainingNTViewQuery(data []byte) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(data)
	writer.Close()
	return buf.Bytes()
}

func Runlike_page4_ProfilePlusLikeChainingNTViewQuery(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var profileID string
	err = db.QueryRow(`SELECT page_id FROM link_page_for_id_page_table LIMIT 1`).Scan(&profileID)
	if err != nil {
		fmt.Println("❌ ดึง profile_id ไม่สำเร็จ: " + err.Error())
		return
	}

	payload := buildEncodedPayloadRunlike_page4_ProfilePlusLikeChainingNTViewQuery(profileID)
	host := "graph.facebook.com"

	req, err := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBuffer(gzipCompressRunlike_page4_ProfilePlusLikeChainingNTViewQuery([]byte(payload))))
	if err != nil {
		fmt.Println("❌ NewRequest fail: " + err.Error())
		return
	}

	traceID := uuid.New().String()
	sessionID := generateSessionIDRunlike_page4_ProfilePlusLikeChainingNTViewQuery()
	connectionToken := uuid.New().String()

	req.Header = map[string][]string{
		"Authorization":               {"OAuth " + token},
		"Accept-Encoding":             {"zstd, gzip, deflate"},
		"Connection":                  {"keep-alive"},
		"Content-Encoding":            {"gzip"},
		"Content-Type":                {"application/x-www-form-urlencoded"},
		"Host":                        {host},
		"User-Agent":                  {userAgent},
		"X-FB-Friendly-Name":          {"ProfilePlusLikeChainingNTViewQuery"},
		"x-fb-client-ip":              {"True"},
		"x-fb-connection-token":       {connectionToken},
		"X-FB-Connection-Type":        {"MOBILE.HSDPA"},
		"x-fb-device-group":           {deviceGroup},
		"X-FB-HTTP-Engine":            {"Liger"},
		"x-fb-privacy-context":        {"595353684716822"},
		"x-fb-qpl-active-flows-json":  {`{"schema_version":"v2","inprogress_qpls":[],"snapshot_attributes":{}}`},
		"X-FB-Request-Analytics-Tags": {`{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`},
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

	_, err = db.Exec("DELETE FROM link_page_for_id_page_table WHERE page_id = ?", profileID)
	if err != nil {
		fmt.Println("❌ ลบ profile_id ไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ profile_id สำเร็จ:", profileID)
	}
}
