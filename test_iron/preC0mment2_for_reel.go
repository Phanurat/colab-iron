package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
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

func extractPostIDpreC0mment2_for_reel(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	reStory := regexp.MustCompile(`story_fbid=(\d+)`)
	rePath := regexp.MustCompile(`facebook\.com/(\d+)/(?:videos|posts)/(\d+)`)
	reReel := regexp.MustCompile(`/reel/(\d+)`)
	reGeneric := regexp.MustCompile(`/posts/(\d+)|/videos/(\d+)`)

	if match := reStory.FindStringSubmatch(rawurl); len(match) > 1 {
		return match[1], nil
	}
	if match := rePath.FindStringSubmatch(rawurl); len(match) > 2 {
		return match[2], nil
	}
	if match := reGeneric.FindStringSubmatch(u.Path); len(match) > 1 {
		if match[1] != "" {
			return match[1], nil
		}
		return match[2], nil
	}
	if match := reReel.FindStringSubmatch(u.Path); len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("ไม่พบ postID ในลิงก์")
}

func generateFeedbackIDpreC0mment2_for_reel(postID string) string {
	return base64.StdEncoding.EncodeToString([]byte("feedback:" + postID))
}

func randomExcellentBandwidthpreC0mment2_for_reel() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func generateUUIDpreC0mment2_for_reel() string {
	return uuid.New().String()
}

func generateHex32preC0mment2_for_reel() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// connToken := generateHex32()

func RunpreC0mment2_for_reel(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	host := "graph.facebook.com"
	//	address := host + ":443"
	clientDocID := "33213602967574118353564342259"
	friendlyName := "FetchPredictiveTextSuggestions"
	connToken := generateHex32preC0mment2_for_reel()

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

	var token, actorID, userAgent, netHni, simHni, deviceGroup string
	err = db.QueryRow(`SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1`).Scan(&token, &actorID, &userAgent, &netHni, &simHni, &deviceGroup)
	if err != nil {
		fmt.Println("❌ ดึง app_profiles ไม่ได้: " + err.Error())
		return
	}

	var link string
	err = db.QueryRow(`SELECT link FROM like_reel_and_comment_reel_table LIMIT 1`).Scan(&link)
	if err != nil {
		fmt.Println("❌ ดึงลิงก์ไม่สำเร็จ: " + err.Error())
		return
	}

	postID, err := extractPostIDpreC0mment2_for_reel(link)
	if err != nil {
		fmt.Println("❌ ขุด postID ล้มเหลว: " + err.Error())
		return
	}
	feedbackID := generateFeedbackIDpreC0mment2_for_reel(postID)
	traceID := generateUUIDpreC0mment2_for_reel()

	variables := fmt.Sprintf(`{"queryParams":{"text":{"sensitive_string_value":""},"surface":"comment","is_mobile":true,"feedback_id":"%s"}}`, feedbackID)

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
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)
	form.Set("client_trace_id", traceID)

	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	_, _ = gz.Write([]byte(form.Encode()))
	_ = gz.Close()

	// proxy := os.Getenv("USE_PROXY")
	// auth := os.Getenv("USE_PROXY_AUTH")

	// conn, err := net.DialTimeout("tcp", proxy, 10*time.Second)
	// if err != nil {
	// 	panic("❌ Proxy failed: " + err.Error())
	// }

	// reqLine := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n", address, host)
	// if auth != "" {
	// 	reqLine += "Proxy-Authorization: Basic " + auth + "\r\n"
	// }
	// reqLine += "\r\n"
	// fmt.Fprintf(conn, reqLine)

	// br := bufio.NewReader(conn)
	// status, _ := br.ReadString('\n')
	// if !strings.Contains(status, "200") {
	// 	panic("❌ CONNECT failed: " + status)
	// }
	// for {
	// 	line, _ := br.ReadString('\n')
	// 	if line == "\r\n" || line == "" {
	// 		break
	// 	}
	// }

	// tlsConn := utls.UClient(conn, &utls.Config{ServerName: host}, utls.HelloAndroid_11_OkHttp)
	// if err := tlsConn.Handshake(); err != nil {
	// 	panic("❌ TLS handshake fail: " + err.Error())
	// }

	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", &compressed)
	req.Header.Set("Authorization", "OAuth "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-token", connToken)
	req.Header.Set("X-FB-Connection-Type", "WIFI")
	req.Header.Set("x-fb-device-group", deviceGroup)
	req.Header.Set("X-FB-Friendly-Name", friendlyName)
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "204001154719270")
	req.Header.Set("x-fb-qpl-active-flows-json", `{"schema_version":"v2","inprogress_qpls":[],"snapshot_attributes":{}}`)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-rmd", "state=URL_ELIGIBLE")
	req.Header.Set("x-fb-server-cluster", "True")
	//	req.Header.Set("x-fb-session-id", "nid=yTInr91goUUA;tid=4907;nc=0;fc=2;bc=2;cid=70f2c354d4a435885857b75a7b5fbd38")
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+traceID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthpreC0mment2_for_reel())
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
