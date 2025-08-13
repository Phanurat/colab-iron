// befor_reactiontype_like_DelightsMLEAnimationQuery.go (FIXED – compiles)

package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// ---------- HELPERS ----------

func isNumericlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func extractPostIDlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(link string) (string, error) {
	var postID string
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	reStory := regexp.MustCompile(`story_fbid=(\d+)`)
	rePath := regexp.MustCompile(`facebook\.com/(\d+)/(?:videos|posts)/(\d+)`)

	if match := reStory.FindStringSubmatch(link); len(match) > 1 {
		postID = match[1]
	}
	if match := rePath.FindStringSubmatch(link); len(match) > 2 {
		postID = match[2]
	}
	if postID == "" {
		re := regexp.MustCompile(`/posts/(\d+)|/videos/(\d+)`)
		match := re.FindStringSubmatch(u.Path)
		if len(match) > 1 {
			if match[1] != "" {
				postID = match[1]
			} else {
				postID = match[2]
			}
		}
	}
	return postID, nil
}

func genFeedbackIDlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(postID string) string {
	return base64.StdEncoding.EncodeToString([]byte("feedback:" + postID))
}

// ---------- GENERATORS ----------

func generateTraceIDlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery() string {
	return uuid.New().String()
}

func generateLoggingIDslike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(t string) string {
	return "graphql:" + t
}

// ---------- INIT ----------

func initlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery() {
	rand.Seed(time.Now().UnixNano())
}

// ---------- DB ----------

func loadAppProfilelike_only_befor_reactiontype_like_DelightsMLEAnimationQuery() (token, userAgent, netHni, simHni, deviceGroup string) {
	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "fb_comment_system.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/fb_comment_system.db")

	err = db.QueryRow(`
		SELECT access_token, user_agent, net_hni, sim_hni, device_group
		FROM app_profiles LIMIT 1`,
	).Scan(&token, &userAgent, &netHni, &simHni, &deviceGroup)
	if err != nil {
		panic("❌ ดึง app_profiles ไม่สำเร็จ: " + err.Error())
	}
	return
}

func randomExcellentBandwidthlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000 // 20 Mbps
	max := 35000000 // 35 Mbps
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// ---------- MAIN ----------

func Runlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	// --- OPEN DB ---
	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "fb_comment_system.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/fb_comment_system.db")

	var accessToken, userID, userAgent, netHni, simHni string
	err = db.QueryRow(`
		SELECT access_token, actor_id, user_agent, net_hni, sim_hni
		FROM app_profiles LIMIT 1`,
	).Scan(&accessToken, &userID, &userAgent, &netHni, &simHni)
	if err != nil {
		panic("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
	}

	// --- GET LINK ---
	var link string
	err = db.QueryRow(`SELECT link FROM like_only_table LIMIT 1`).Scan(&link)
	if err != nil {
		panic("❌ ดึงลิงก์จาก like_only_table ไม่สำเร็จ: " + err.Error())
	}

	// --- EXTRACT POST ID ---
	postID, err := extractPostIDlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(link)
	if err != nil || postID == "" {
		panic("❌ ดึง postID จากลิงก์ไม่สำเร็จ: " + link)
	}
	feedbackID := genFeedbackIDlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(postID)

	// --- LOAD PROFILE VARIABLES ---
	accessToken, userAgent, netHni, simHni, deviceGroup := loadAppProfilelike_only_befor_reactiontype_like_DelightsMLEAnimationQuery()

	// --- CONSTANTS ---
	host := "graph.facebook.com"
	//address := host + ":443"
	clientDocID := "23228860342874045091256064671"

	variables := map[string]any{
		"id":          feedbackID,
		"user_action": "REACTION_LIKE_SENT",
	}

	// --- DYNAMIC VALUES ---
	traceID := generateTraceIDlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery()
	loggingIDs := generateLoggingIDslike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(traceID)
	bandwidth := randomExcellentBandwidthlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery()

	// --- BUILD BODY ---
	variablesJSON, _ := json.Marshal(variables)
	form := url.Values{
		"method":                   {"post"},
		"pretty":                   {"false"},
		"format":                   {"json"},
		"server_timestamps":        {"true"},
		"locale":                   {"en_US"},
		"fb_api_req_friendly_name": {"DelightsMLEAnimationQuery"},
		"fb_api_caller_class":      {"graphservice"},
		"client_doc_id":            {clientDocID},
		"variables":                {string(variablesJSON)},
		"fb_api_analytics_tags":    {`["GraphServices"]`},
		"client_trace_id":          {loggingIDs},
	}

	bodyBuf := new(bytes.Buffer)
	gz := gzip.NewWriter(bodyBuf)
	_, _ = gz.Write([]byte(form.Encode()))
	gz.Close()

	// // --- PROXY / DIRECT CONNECT ---
	// proxyAddr := os.Getenv("USE_PROXY")
	// proxyAuth := os.Getenv("USE_PROXY_AUTH")

	// var conn net.Conn
	// if proxyAddr != "" {
	// 	conn, err = net.DialTimeout("tcp", proxyAddr, 10*time.Second)
	// 	if err != nil {
	// 		panic("❌ Proxy fail: " + err.Error())
	// 	}
	// 	connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n", address, host)
	// 	if proxyAuth != "" {
	// 		connectReq += "Proxy-Authorization: Basic " + proxyAuth + "\r\n"
	// 	}
	// 	connectReq += "\r\n"
	// 	fmt.Fprint(conn, connectReq)

	// 	br := bufio.NewReader(conn)
	// 	respLine, _ := br.ReadString('\n')
	// 	if !strings.Contains(respLine, "200") {
	// 		panic("❌ CONNECT fail: " + respLine)
	// 	}
	// 	for {
	// 		l, _ := br.ReadString('\n')
	// 		if l == "\r\n" || l == "" {
	// 			break
	// 		}
	// 	}
	// } else {
	// 	conn, err = net.DialTimeout("tcp", address, 10*time.Second)
	// 	if err != nil {
	// 		panic("❌ Dial fail: " + err.Error())
	// 	}
	// }

	// // --- TLS HANDSHAKE ---
	// utlsConn := utls.UClient(conn, &utls.Config{ServerName: host}, utls.HelloAndroid_11_OkHttp)
	// if err = utlsConn.Handshake(); err != nil {
	// 	panic("❌ TLS handshake fail: " + err.Error())
	// }

	// --- BUILD REQUEST ---
	req, _ := http.NewRequest(
		"POST",
		"https://"+host+"/graphql",
		bytes.NewReader(bodyBuf.Bytes()),
	)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", deviceGroup)
	req.Header.Set("X-FB-Friendly-Name", "DelightsMLEAnimationQuery")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "697694927744066")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", bandwidth)
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")
	req.Header.Set("x-fb-ta-logging-ids", loggingIDs)

	// ถ้าต้อง spoof header เพิ่มเติมก็เติมตรงนี้
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
