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

// ✅ เพิ่ม init เพื่อ seed rand ให้ไม่สุ่มซ้ำ
func init() {
	rand.Seed(time.Now().UnixNano())
}

// ✅ ส่วน helper และ util คงเดิม ไม่มีการเปลี่ยนแปลง

func isNumericbefor_reactiontype_like_DelightsMLEAnimationQuery(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func extractPostIDbefor_reactiontype_like_DelightsMLEAnimationQuery(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	reStory := regexp.MustCompile(`story_fbid=(\d+)`)
	rePath := regexp.MustCompile(`facebook\\.com/(\d+)/(?:videos|posts)/(\d+)`)

	if match := reStory.FindStringSubmatch(link); len(match) > 1 {
		return match[1], nil
	}
	if match := rePath.FindStringSubmatch(link); len(match) > 2 {
		return match[2], nil
	}

	re := regexp.MustCompile(`/posts/(\d+)|/videos/(\d+)`)
	match := re.FindStringSubmatch(u.Path)
	if len(match) > 1 {
		if match[1] != "" {
			return match[1], nil
		} else {
			return match[2], nil
		}
	}
	return "", nil
}

func genFeedbackIDbefor_reactiontype_like_DelightsMLEAnimationQuery(postID string) string {
	return base64.StdEncoding.EncodeToString([]byte("feedback:" + postID))
}

func generateTraceIDbefor_reactiontype_like_DelightsMLEAnimationQuery() string {
	return uuid.New().String()
}
func generateLoggingIDsbefor_reactiontype_like_DelightsMLEAnimationQuery(t string) string {
	return "graphql:" + t
}

func randomExcellentBandwidthbefor_reactiontype_like_DelightsMLEAnimationQuery() string {
	min := 20000000 // 20 Mbps
	max := 35000000 // 35 Mbps
	return strconv.Itoa(min + rand.Intn(max-min+1))
}

func loadAppProfilebefor_reactiontype_like_DelightsMLEAnimationQuery() (token, userAgent, netHni, simHni, deviceGroup string) {
	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}
	dbPath := filepath.Join(folder, "fb_comment_system.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
		return
	}
	defer db.Close()

	err = db.QueryRow(`SELECT access_token, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1`).Scan(&token, &userAgent, &netHni, &simHni, &deviceGroup)
	if err != nil {
		fmt.Println("❌ ดึง app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}
	return
}

// ✅ เปลี่ยนจาก main เป็น RunDelightsMLEAnimationQuery พร้อมรับ TLSConnection
func Runbefor_reactiontype_like_DelightsMLEAnimationQuery(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ แสดง Proxy ที่ใช้อยู่

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}
	dbPath := filepath.Join(folder, "fb_comment_system.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
		return
	}
	defer db.Close()

	// --- GET LINK ---
	var link string
	err = db.QueryRow(`SELECT link FROM like_and_comment_table LIMIT 1`).Scan(&link)
	if err != nil {
		fmt.Println("❌ ดึงลิงก์จาก like_and_comment_table ไม่สำเร็จ: " + err.Error())
		return
	}

	postID, err := extractPostIDbefor_reactiontype_like_DelightsMLEAnimationQuery(link)
	if err != nil || postID == "" {
		fmt.Println("❌ ดึง postID จากลิงก์ไม่สำเร็จ: " + link)
		return
	}
	feedbackID := genFeedbackIDbefor_reactiontype_like_DelightsMLEAnimationQuery(postID)

	accessToken, userAgent, netHni, simHni, deviceGroup := loadAppProfilebefor_reactiontype_like_DelightsMLEAnimationQuery()

	host := "graph.facebook.com"
	clientDocID := "23228860342874045091256064671"
	traceID := generateTraceIDbefor_reactiontype_like_DelightsMLEAnimationQuery()
	loggingIDs := generateLoggingIDsbefor_reactiontype_like_DelightsMLEAnimationQuery(traceID)
	bandwidth := randomExcellentBandwidthbefor_reactiontype_like_DelightsMLEAnimationQuery()

	variables := map[string]any{
		"id":          feedbackID,
		"user_action": "REACTION_LIKE_SENT",
	}
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

	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewReader(bodyBuf.Bytes()))
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
