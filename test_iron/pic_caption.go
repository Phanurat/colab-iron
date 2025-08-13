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
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func randomExcellentBandwidthpic_caption() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func randomUUIDpic_caption() string {
	b := make([]byte, 16)
	rand.Read(b)
	// UUID v4 format
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func generateIdempotenceTokenpic_caption(sessionID string) string {
	return "FEED_" + sessionID
}

func generateEncodedRequestIDpic_caption() string {
	randomUID := randomUUIDpic_caption()
	payload := fmt.Sprintf("GET_WHATSAPP_MESSAGES_UNDIRECTED_FEED_COMPOSER:%s", randomUID)
	return base64.StdEncoding.EncodeToString([]byte(payload))
}

func encodeGzippic_caption(data []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(data)
	w.Close()
	return buf.Bytes()
}

func Runpic_caption(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	sessionID := randomUUIDpic_caption()
	clientMutationID := randomUUIDpic_caption()
	traceID := randomUUIDpic_caption()
	requestID := generateEncodedRequestIDpic_caption()
	timestamp := time.Now().Unix()
	idempotenceToken := generateIdempotenceTokenpic_caption(sessionID)

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

	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		panic("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
	}

	var mediaID string
	err = db.QueryRow("SELECT media_id FROM pic_caption_table LIMIT 1").Scan(
		&mediaID)
	if err != nil {
		panic("❌ ดึงข้อมูล pic_caption_table ไม่สำเร็จ: " + err.Error())
	}

	var captiontext string
	err = db.QueryRow("SELECT caption_text FROM pic_caption_text_table LIMIT 1").Scan(
		&captiontext)
	if err != nil {
		panic("❌ ดึงข้อมูล pic_caption_text_table ไม่สำเร็จ: " + err.Error())
	}

	variables := map[string]interface{}{
		"image_low_height":                          2048,
		"image_medium_width":                        540,
		"automatic_photo_captioning_enabled":        "false",
		"angora_attachment_profile_image_size":      105,
		"poll_facepile_size":                        105,
		"default_image_scale":                       "3",
		"image_high_height":                         2048,
		"should_fetch_adaptive_ufi":                 true,
		"image_large_aspect_height":                 565,
		"image_low_width":                           360,
		"image_medium_height":                       2048,
		"include_mentions_messenger_sharing_params": true,
		"media_type":                                "image/jpeg",
		"size_style":                                "contain-fit",
		"image_high_width":                          1080,
		"input": map[string]interface{}{
			"producer_supported_features": []string{"LIGHTWEIGHT_REPLY"},
			"tag_expansion_metadata": map[string]interface{}{
				"tag_expansion_ids": []string{},
			},
			"place_attachment_setting": "SHOW_ATTACHMENT",
			"past_time": map[string]interface{}{
				"time_since_original_post": 2,
			},
			"logging": map[string]interface{}{
				"composer_session_id": sessionID,
			},
			"is_throwback_post": "NOT_THROWBACK_POST",
			"navigation_data": map[string]interface{}{
				"attribution_id_v2": "NewsFeedFragment,native_newsfeed,,1749801007.610,267801139,4748854339,,",
			},
			"reshare_original_post": "SHARE_LINK_ONLY",
			"idempotence_token":     idempotenceToken,
			"camera_post_context": map[string]string{
				"source":           "COMPOSER",
				"platform":         "FACEBOOK",
				"deduplication_id": sessionID,
			},
			"connection_class":        "EVERYONE",
			"composer_type":           "status",
			"composer_source_surface": "timeline",
			"message": map[string]string{
				"text": captiontext,
			},

			"implicit_with_tags_ids": []string{},
			"composer_entry_point":   "inline_composer",
			"nectar_module":          "timeline_composer",
			"extensible_sprouts_ranker_request": map[string]string{
				"RequestID": requestID,
			},
			"composer_entry_picker": "NULL",
			"client_mutation_id":    clientMutationID,
			"audiences": []map[string]interface{}{
				{
					"undirected": map[string]interface{}{
						"privacy": map[string]interface{}{
							"tag_expansion_state": "UNSPECIFIED",
							"deny":                []string{},
							"base_state":          "EVERYONE",
							"allow":               []string{},
						},
					},
				},
			},
			"source":                "MOBILE",
			"actor_id":              userID,
			"audiences_is_complete": true,
			"attachments": []map[string]interface{}{
				{
					"photo": map[string]interface{}{
						"unified_stories_media_source": "CAMERA_ROLL",
						"story_media_audio_data": map[string]string{
							"raw_media_type": "PHOTO",
						},
						"id": mediaID,
					},
				},
			},
			"action_timestamp": timestamp,
			"composer_session_events_log": map[string]interface{}{
				"number_of_keystrokes":  100,
				"number_of_copy_pastes": 0,
				"composition_duration":  5,
			},
		},
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
		"poll_voters_count":                       5,
		"action_location":                         "feed",
		"reading_attachment_profile_image_height": 354,
		"include_image_ranges":                    true,
		"profile_image_size":                      105,
		"bloks_version":                           "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		"profile_pic_media_type":                  "image/x-auto",
		"angora_attachment_cover_image_size":      1260,
		"question_poll_count":                     100,
		"image_large_aspect_width":                1080,
		"reading_attachment_profile_image_width":  236,
		"fetch_fbc_header":                        true,
		"should_fetch_fallback_actions":           true,
	}

	varJSON, _ := json.Marshal(variables)
	form := fmt.Sprintf("method=post&pretty=false&format=json&server_timestamps=true&locale=en_US&fb_api_req_friendly_name=ComposerStoryCreateMutation&fb_api_caller_class=graphservice&client_doc_id=91093790612716748765152950249&variables=%s", url.QueryEscape(string(varJSON)))
	gzipBody := encodeGzippic_caption([]byte(form))

	req, _ := http.NewRequest("POST", "https://graph.facebook.com/graphql", bytes.NewReader(gzipBody))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "ComposerStoryCreateMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "496463117678580")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+traceID)
	req.Header.Set("client_trace_id", traceID)
	req.Header.Set("x-fb-net-hni", netHni)                                             // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                             // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthpic_caption()) //เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	req.ContentLength = int64(len(gzipBody))

	// // ---------- SEND ----------
	// err = req.Write(utlsConn)
	// if err != nil {
	// 	panic("❌ Write fail: " + err.Error())
	// }

	// resp, err := http.ReadResponse(bufio.NewReader(utlsConn), req)
	// if err != nil {
	// 	panic("❌ Read fail: " + err.Error())
	// }
	// defer resp.Body.Close()

	// var reader io.ReadCloser
	// switch resp.Header.Get("Content-Encoding") {
	// case "gzip":
	// 	reader, err = gzip.NewReader(resp.Body)
	// 	if err != nil {
	// 		panic("❌ GZIP decompress fail: " + err.Error())
	// 	}
	// 	defer reader.Close()
	// default:
	// 	reader = resp.Body
	// }

	// bodyResp, err := io.ReadAll(reader)
	// if err != nil {
	// 	panic("❌ Body read fail: " + err.Error())
	// }

	// fmt.Println("✅ Status:", resp.Status)
	// fmt.Println("📦 Response:", string(bodyResp))
	// fmt.Println("🧩 Proxy:", proxy)

	bw := tlsConns.RWGraph.Writer
	br := tlsConns.RWGraph.Reader

	err = req.Write(bw)
	if err != nil {
		fmt.Println("❌ Write fail: " + err.Error())
		return
	}
	bw.Flush()

	resp, err := http.ReadResponse(br, req)
	if err != nil {
		fmt.Println("❌ Read fail: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// Debug response headers
	fmt.Println("📥 Response Headers:")
	for name, values := range resp.Header {
		fmt.Printf("  %s: %s\n", name, values[0])
	}

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
	fmt.Printf("📦 Response length: %d bytes\n", len(bodyResp))
	fmt.Printf("📦 Response content: '%s'\n", string(bodyResp))

	//mediaID, statusText

	_, err = db.Exec(`DELETE FROM pic_caption_text_table WHERE caption_text = ?`, captiontext)
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ caption_text ออกจากฐานข้อมูลแล้ว:", captiontext)
	}

}
