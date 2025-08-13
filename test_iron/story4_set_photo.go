// composer_story.go (FULL)
// ยิง ComposerStoryCreateMutation พร้อม spoof headers, gzip form, uTLS + proxy + ฟังก์ชันเจน UUID/Session แบบไม่ซ้ำ
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

	_ "github.com/mattn/go-sqlite3"
)

func getLatestPhotoIDRunstory4_set_photo() string {
	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "fb_comment_system.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())

	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/fb_comment_system.db")

	row := db.QueryRow("SELECT pic_id FROM story_photo_id_table ORDER BY id DESC LIMIT 1")
	var photoID string
	err = row.Scan(&photoID)
	if err != nil {
		fmt.Println("❌ ดึง pic_id ไม่สำเร็จ: " + err.Error())
	}
	return photoID
}

var (
	hostRunstory4_set_photo     = "graph.facebook.com"
	endpointRunstory4_set_photo = "https://graph.facebook.com/graphql"
)

func genUUIDRunstory4_set_photo() string {
	u := make([]byte, 16)
	rand.Read(u)
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func genVisitationIDRunstory4_set_photo() string {
	t := time.Now()
	epoch := t.Unix()
	ms := t.UnixNano() / 1e6 % 100000
	return fmt.Sprintf("4748854339:dcb73:0:%d.%d", epoch, ms)
}

func buildFormDataRunstory4_set_photo() []byte {

	photoID := getLatestPhotoIDRunstory4_set_photo()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "fb_comment_system.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/fb_comment_system.db")

	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	sessionID := genUUIDRunstory4_set_photo()
	mutationID := genUUIDRunstory4_set_photo()
	timestamp := time.Now().Unix()
	attribution := fmt.Sprintf("NewsFeedFragment,native_newsfeed,,%d,42968185,4748854339,,;InspirationComposerActivity,stories_composer,,%d,81049347,,,", timestamp, timestamp-12)

	vars := map[string]interface{}{
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
			"tag_expansion_metadata":   map[string]interface{}{"tag_expansion_ids": []string{}},
			"place_attachment_setting": "SHOW_ATTACHMENT",
			"past_time":                map[string]int{"time_since_original_post": 2},
			"logging":                  map[string]string{"composer_session_id": sessionID},
			"is_throwback_post":        "NOT_THROWBACK_POST",
			"inspiration_prompts": []map[string]string{
				{"prompt_type": "MANUAL", "prompt_tracking_string": "0", "prompt_id": "1752514608329267"},
			},
			"navigation_data":       map[string]string{"attribution_id_v2": attribution},
			"reshare_original_post": "SHARE_LINK_ONLY",
			"idempotence_token":     "STORIES_" + mutationID,
			"camera_post_context": map[string]string{
				"source": "COMPOSER", "platform": "FACEBOOK", "deduplication_id": mutationID,
			},
			"connection_class":            "EXCELLENT",
			"composer_type":               "story",
			"composer_source_surface":     "newsfeed",
			"message":                     map[string]string{"text": ""},
			"implicit_with_tags_ids":      []string{},
			"composer_entry_point":        "add_to_story_first_pog",
			"composer_entry_picker":       "NULL",
			"client_mutation_id":          mutationID,
			"producer_supported_features": []string{"LIGHTWEIGHT_REPLY"},
			"audiences": []map[string]interface{}{
				{"stories": map[string]interface{}{"self": map[string]string{"target_id": userID}}},
			},
			"source":                "MOBILE",
			"actor_id":              userID,
			"audiences_is_complete": true,
			"attachments": []map[string]interface{}{
				{"photo": map[string]interface{}{
					"ml_media_tracking_data": map[string]string{"media_tracking_id": "281202717"},
					"id":                     photoID, "filter_state": "FILTERED",
					"unified_stories_media_source": "COMPOSER_GALLERY",
					"filter_name":                  "PassThrough",
					"story_media_audio_data":       map[string]string{"raw_media_type": "PHOTO"},
					"capture_mode":                 "NORMAL",
				}},
			},
			"action_timestamp": timestamp,
			"composer_session_events_log": map[string]int{
				"number_of_keystrokes": 0, "number_of_copy_pastes": 0, "composition_duration": 0,
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

	varsJSON, _ := json.Marshal(vars)

	visitation := genVisitationIDRunstory4_set_photo()
	session := "UFS-" + strings.ReplaceAll(genUUIDRunstory4_set_photo(), "-", "")
	analytics := fmt.Sprintf(`["surface_hierarchy=NewsFeedFragment,native_newsfeed,null;FbChromeFragment,null,tap_back_button;FbMainTabActivity,unknown,null","visitation_id=%s","session_id=%s","GraphServices"]`, visitation, session)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "ComposerStoryCreateMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "91093790612716748765152950249")
	form.Set("variables", string(varsJSON))
	form.Set("fb_api_analytics_tags", analytics)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(form.Encode()))
	gz.Close()

	return buf.Bytes()
}

func randomExcellentBandwidthRunstory4_set_photo() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runstory4_set_photo(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	photoID := getLatestPhotoIDRunstory4_set_photo()
	body := buildFormDataRunstory4_set_photo()
	//address := host + ":443"

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

	req, _ := http.NewRequest("POST", endpointRunstory4_set_photo, bytes.NewReader(body))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostRunstory4_set_photo)
	req.Header.Set("User-Agent", userAgent)
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
	req.Header.Set("x-fb-net-hni", netHni)                                                     // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                                     // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthRunstory4_set_photo()) //เพิ่มเข้าไป
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

	_, err = db.Exec(`DELETE FROM story_photo_id_table WHERE pic_id = ?`, photoID)
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ friend_id ออกจากฐานข้อมูลแล้ว:", photoID)
	}

}
