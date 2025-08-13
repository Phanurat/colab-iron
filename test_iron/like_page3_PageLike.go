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

func buildEncodedPayloadRunlike_page3_PageLike(pageID, actorID string) string {
	clientMutationID := uuid.New().String()
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"source":             "page_profile",
			"page_id":            pageID,
			"actor_id":           actorID,
			"client_mutation_id": clientMutationID,
		},
	}
	jsonVars, _ := json.Marshal(variables)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "PageLike")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "92246462512975232024543564417")
	form.Set("variables", string(jsonVars))
	form.Set("fb_api_analytics_tags", `["visitation_id=null","GraphServices"]`)
	form.Set("client_trace_id", clientMutationID) // reuse for simplicity
	return form.Encode()
}

func generateSessionIDRunlike_page3_PageLike() string {
	u := uuid.New().String()
	tid := strconv.Itoa(rand.Intn(900) + 100)
	return fmt.Sprintf("nid=%s;tid=%s;nc=1;fc=2;bc=1;cid=%s", u, tid, uuid.New().String())
}

func gzipCompressRunlike_page3_PageLike(data []byte) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(data)
	writer.Close()
	return buf.Bytes()
}

func Runlike_page3_PageLike(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var token, userAgent, netHni, simHni, deviceGroup, actorID string
	err = db.QueryRow(`SELECT access_token, user_agent, net_hni, sim_hni, device_group, actor_id FROM app_profiles LIMIT 1`).Scan(
		&token, &userAgent, &netHni, &simHni, &deviceGroup, &actorID)
	if err != nil {
		fmt.Println("❌ ดึง app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	var pageID string
	err = db.QueryRow(`SELECT page_id FROM link_page_for_id_page_table LIMIT 1`).Scan(&pageID)
	if err != nil {
		fmt.Println("❌ ดึง page_id ไม่สำเร็จ: " + err.Error())
		return
	}

	payload := buildEncodedPayloadRunlike_page3_PageLike(pageID, actorID)
	host := "graph.facebook.com"

	req, err := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBuffer(gzipCompressRunlike_page3_PageLike([]byte(payload))))
	if err != nil {
		fmt.Println("❌ NewRequest fail: " + err.Error())
		return
	}

	traceID := uuid.New().String()
	sessionID := generateSessionIDRunlike_page3_PageLike()
	connectionToken := uuid.New().String()

	req.Header = map[string][]string{
		"Authorization":               {"OAuth " + token},
		"Accept-Encoding":             {"zstd, gzip, deflate"},
		"Connection":                  {"keep-alive"},
		"Content-Encoding":            {"gzip"},
		"Content-Type":                {"application/x-www-form-urlencoded"},
		"Host":                        {host},
		"User-Agent":                  {userAgent},
		"X-FB-Friendly-Name":          {"PageLike"},
		"x-fb-client-ip":              {"True"},
		"x-fb-connection-token":       {connectionToken},
		"X-FB-Connection-Type":        {"MOBILE.HSDPA"},
		"x-fb-device-group":           {deviceGroup},
		"X-FB-HTTP-Engine":            {"Liger"},
		"x-fb-privacy-context":        {"305228267119416"},
		"x-fb-qpl-active-flows-json":  {`{"schema_version":"v2","inprogress_qpls":[],"snapshot_attributes":{}}`},
		"X-FB-Request-Analytics-Tags": {`{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`},
		"x-fb-rmd":                    {"state=URL_ELIGIBLE"},
		"x-fb-server-cluster":         {"True"},
		"x-fb-session-id":             {sessionID},
		"x-fb-ta-logging-ids":         {"graphql:" + traceID},
		"x-graphql-client-library":    {"graphservice"},
		"x-graphql-request-purpose":   {"fetch"},
		"x-tigon-is-retry":            {"False"},
		"X-FB-Navigation-Chain":       {"ProfileFragment,profile_vnext_tab_posts,,1750787348.786,154347230,,,;ProfileFragment,timeline,tap_sponsored_link,1750787347.917,154347230,,,;NewsFeedFragment,native_newsfeed,tap_top_jewel_bar,1750787334.358,26876560,4748854339,312#10#230#132#230#231,609329695177089"},
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

	//	_, err = db.Exec("DELETE FROM like_page_table WHERE page_id = ?", pageID)
	//	if err != nil {
	//		fmt.Println("❌ ลบ page_id ไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ page_id สำเร็จ:", pageID)
	//	}
}
