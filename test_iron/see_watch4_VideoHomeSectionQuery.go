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

func Runsee_watch4_VideoHomeSectionQuery(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	host := "graph.facebook.com"
	//	address := host + ":443"

	sectionID := "dmg6MTY2MzY5NzM0MDAyOTY0"
	deviceID := uuid.New().String()

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

	vars := map[string]interface{}{
		"section_id":                    sectionID,
		"at_stream_label":               "watch_feed",
		"enable_single_image_ads":       true,
		"max_immersive_image_height":    2048,
		"num_reels_for_watch_ifu":       3,
		"bloks_version":                 "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		"ads_context_data":              map[string]interface{}{},
		"should_fetch_fallback_actions": true,
		"should_include_friend_actions": true,
		"enable_watch_feed_stream":      true,
		"caller":                        "WATCH_FEED_TAB",
		"action_location":               "VIDEO_HOME",
		"video_channel_context_data": map[string]interface{}{
			"player_origin":   "video_home::feed",
			"player_behavior": "WATCH_FEED",
		},
		"default_image_scale":                "3",
		"feed_story_render_location":         "video_home",
		"profile_entry_point":                "VIDEO_HOME",
		"device_id":                          deviceID,
		"enable_watch_feed_pivots":           true,
		"fb_shorts_location":                 "video_home",
		"image_large_aspect_width":           1080,
		"fetch_aggregations":                 true,
		"friend_facepile_profile_image_size": 84,
		"fetch_request_id":                   true,
		"fetch_unit_metadata_social_context": true,
		"enable_watch_feed_edge_header":      true,
		"image_high_width":                   1080,
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
		"scale":                     "3",
		"profile_pic_media_type":    "image/x-auto",
		"image_large_aspect_height": 565,
		"should_fetch_adaptive_ufi": true,
		"image_low_width":           617,
		"watch_nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"pixel_ratio":        3,
			"is_push_on":         true,
			"extra_data":         "{\"is_fullbleed\":true}",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
		"video_channel_entry_point": "VIDEO_HOME",
		"image_medium_width":        814,
		"media_type":                "image/x-auto",
		"profile_image_size":        945,
		"tail_load_type":            map[string]interface{}{"tail_load_type": "NORMAL"},
		"section_after_cursor":      "",
	}

	payload := map[string]string{
		"method":                   "post",
		"pretty":                   "false",
		"format":                   "json",
		"server_timestamps":        "true",
		"locale":                   "en_US",
		"fb_api_req_friendly_name": "VideoHomeSectionQuery",
		"fb_api_caller_class":      "graphservice",
		"client_doc_id":            "200821220517377489446718798218",
		"variables":                encodeJSONsee_watch4_VideoHomeSectionQuery(vars),
		"fb_api_analytics_tags":    "[\"GraphServices\"]",
		"access_token":             accessToken,
	}

	body := encodeFormsee_watch4_VideoHomeSectionQuery(payload)
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	gz.Write([]byte(body))
	gz.Close()

	req, _ := http.NewRequest("POST", "https://"+host+"graphql", &compressed)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "VideoHomeSectionQuery")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "2444622522461689")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthsee_watch4_VideoHomeSectionQuery())
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

func encodeJSONsee_watch4_VideoHomeSectionQuery(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func encodeFormsee_watch4_VideoHomeSectionQuery(data map[string]string) string {
	var buf bytes.Buffer
	for k, v := range data {
		buf.WriteString(fmt.Sprintf("%s=%s&", k, urlEncodesee_watch4_VideoHomeSectionQuery(v)))
	}
	return buf.String()[:buf.Len()-1]
}

func urlEncodesee_watch4_VideoHomeSectionQuery(s string) string {
	return (&url.URL{Path: s}).EscapedPath()
}

func randomExcellentBandwidthsee_watch4_VideoHomeSectionQuery() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}
