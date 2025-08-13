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

// ---------- CONFIG & INIT ----------

func initRunlike_comment_only() {
	rand.Seed(time.Now().UnixNano())
}

// ---------- DB HELPERS ----------

//func loadAppProfile() (accessToken, actorID, userAgent, netHni, simHni, deviceGroup string) {
//	db, err := sql.Open("sqlite3", "./fb_comment_system.db")
//	if err != nil {
//		panic("❌ เปิดฐานข้อมูลไม่ได้: " + err.Error())
//	}
//	defer db.Close()
//
//	err = db.QueryRow(`
//		SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group
//		FROM app_profiles LIMIT 1`,
//	).Scan(&accessToken, &actorID, &userAgent, &netHni, &simHni, &deviceGroup)
//	if err != nil {
//		panic("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
//	}
//	return
//}

//func loadReactionJob() (reactionType, link string) {
//	db, err := sql.Open("sqlite3", "./fb_comment_system.db")
//	if err != nil {
//		panic("❌ เปิดฐานข้อมูลไม่ได้: " + err.Error())
//	}
//	defer db.Close()
//
//	err = db.QueryRow(`SELECT reaction_type, like FROM like_comment_only LIMIT 1`).Scan(&reactionType, &link)
//	if err != nil {
//		panic("❌ ดึงข้อมูล like_comment_only ไม่สำเร็จ: " + err.Error())
//	}
//	return
//}

// ---------- UTILITIES ----------

func randomExcellentBandwidthRunlike_comment_only() string {
	min, max := 20_000_000, 35_000_000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func isNumericRunlike_comment_only(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func randIntRunlike_comment_only() int {
	return int(time.Now().UnixNano() % 100_000_000)
}

// ---------- FACEBOOK ID EXTRACTION ----------

func extractCommentIDRunlike_comment_only(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if id := u.Query().Get("comment_id"); id != "" {
		return id, nil
	}
	return "", fmt.Errorf("❌ ไม่พบ comment_id")
}

func extractFacebookIDsRunlike_comment_only(rawURL string) (ownerID, postID string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}
	query := u.Query()

	reStory := regexp.MustCompile(`story_fbid=(\d+)`)
	rePath := regexp.MustCompile(`facebook\.com/(\d+)/(?:videos|posts)/(\d+)`)
	reUsername := regexp.MustCompile(`facebook\.com/([^/?&]+)`)

	if m := reStory.FindStringSubmatch(rawURL); len(m) > 1 {
		postID = m[1]
	}
	if m := rePath.FindStringSubmatch(rawURL); len(m) > 2 {
		ownerID = m[1]
		postID = m[2]
	}
	if postID == "" {
		re := regexp.MustCompile(`/posts/(\d+)|/videos/(\d+)`)
		if m := re.FindStringSubmatch(u.Path); len(m) > 1 {
			if m[1] != "" {
				postID = m[1]
			} else {
				postID = m[2]
			}
		}
	}
	if id := query.Get("id"); id != "" {
		ownerID = id
	}
	if ownerID == "" {
		if m := reUsername.FindStringSubmatch(rawURL); len(m) > 1 {
			username := m[1]
			if isNumericRunlike_comment_only(username) {
				ownerID = username
			} else {
				ownerID, err = getFBIDFromUsernameRunlike_comment_only(username)
				if err != nil {
					return "", "", err
				}
			}
		}
	}
	return ownerID, postID, nil
}

func getFBIDFromUsernameRunlike_comment_only(username string) (string, error) {
	client := &http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	req, _ := http.NewRequest("GET", "https://mbasic.facebook.com/"+username, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if loc := resp.Header.Get("Location"); strings.HasPrefix(loc, "intent://profile/") {
		re := regexp.MustCompile(`intent://profile/(\d+)`)
		if m := re.FindStringSubmatch(loc); len(m) > 1 {
			return m[1], nil
		}
	}

	resp2, err := http.Get("https://mbasic.facebook.com/" + username)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()
	body, _ := io.ReadAll(resp2.Body)
	text := string(body)

	re := regexp.MustCompile(`owner_id=(\d+)`)
	if m := re.FindStringSubmatch(text); len(m) > 1 {
		return m[1], nil
	}
	re = regexp.MustCompile(`profile\.php\?id=(\d+)`)
	if m := re.FindStringSubmatch(text); len(m) > 1 {
		return m[1], nil
	}
	return "", fmt.Errorf("❌ ไม่พบ owner_id จาก username")
}

// ---------- META FETCH (Graph API) ----------

type MetaData struct {
	PageID   string
	ownerID  string
	Tracking []interface{}
}

func fetchMetaRunlike_comment_only(postID, commentID, token string) *MetaData {
	var endpoint string
	if commentID != "" {
		endpoint = fmt.Sprintf("https://graph.facebook.com/%s?fields=from{ id }&access_token=%s", commentID, url.QueryEscape(token))
	} else {
		endpoint = fmt.Sprintf("https://graph.facebook.com/%s?fields=from{ id }&access_token=%s", postID, url.QueryEscape(token))
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("⚠️  ดึง meta ไม่สำเร็จ:", err)
		return &MetaData{Tracking: []interface{}{}}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	type fbResp struct {
		From struct {
			ID string `json:"id"`
		} `json:"from"`
	}
	var data fbResp
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("⚠️  แปลง meta ไม่สำเร็จ:", err)
		return &MetaData{Tracking: []interface{}{}}
	}

	return &MetaData{
		PageID:   data.From.ID,
		ownerID:  data.From.ID,
		Tracking: []interface{}{},
	}
}

// ---------- FEEDBACK / ATTRIBUTION ----------

func generateFeedbackIDRunlike_comment_only(id string) string {
	return base64.StdEncoding.EncodeToString([]byte("feedback:" + id))
}

// generateAttributionID สร้างค่า attribution_id_v2 แบบไม่ซ้ำ
func generateAttributionIDRunlike_comment_only() string {
	now := float64(time.Now().UnixNano()) / 1e9 // epoch วินาที (มีจุดทศนิยม)
	r1 := rand.Uint64() % 1_000_000_000         // สุ่ม 9 หลัก
	r2 := rand.Uint64() % 1_000_000_000         // สุ่ม 9 หลัก
	return fmt.Sprintf(
		"SimpleUFIPopoverFragment,story_feedback_flyout,tap_bling_bar_comment,%.6f,%d,,,,;NewsFeedFragment,native_newsfeed,cold_start,%.6f,%d,4748854339,36#301,1330559721764297",
		now, r1, now-12, r2,
	)
}

// ---------- REACTION ID MAP ----------

var reactionIDsRunlike_comment_only = map[string]string{
	"like":  "1635855486666999",
	"love":  "1678524932434102",
	"haha":  "115940658764963",
	"wow":   "478547315650144",
	"sad":   "444813342392137",
	"angry": "604753422931501",
	"care":  "613557422527858",
}

// ---------- MAIN ----------

func Runlike_comment_only(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var accessToken, userID, userAgent, netHni, simHni, deviceGroup string
	err = db.QueryRow(`
		SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group
		FROM app_profiles LIMIT 1`,
	).Scan(&accessToken, &userID, &userAgent, &netHni, &simHni, &deviceGroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}
	//return

	var reactionType, link string
	err = db.QueryRow("SELECT reaction_type, link FROM like_comment_only_table LIMIT 1").Scan(&reactionType, &link)
	if err != nil {
		fmt.Println("❌ ดึงลิงก์จาก like_comment_only_table ไม่สำเร็จ: " + err.Error())
		return
	}

	// 1. โหลดโปรไฟล์แอป + งานที่จะทำ
	//	accessToken, actorID, userAgent, netHni, simHni, deviceGroup := loadAppProfile()
	//	reactionType, postLink := loadReactionJob()

	// 2. แยก owner / post / comment
	_, postID, err := extractFacebookIDsRunlike_comment_only(link)
	if err != nil {
		fmt.Println("❌ ขุด owner_id/post_id ไม่สำเร็จ: " + err.Error())
		return
	}
	commentID, _ := extractCommentIDRunlike_comment_only(link)

	// 3. สร้าง feedback_id
	targetID := postID
	if commentID != "" {
		targetID = commentID
	}
	feedbackID := generateFeedbackIDRunlike_comment_only(targetID)

	// 4. ดึง meta
	meta := fetchMetaRunlike_comment_only(postID, commentID, accessToken)

	// 5. เตรียม variables
	mutationID := uuid.New().String()
	traceID := uuid.New().String()
	attribution := generateAttributionIDRunlike_comment_only() ////////////////////////////////////////
	actionTimestamp := time.Now().Unix()
	reactionID := reactionIDsRunlike_comment_only[reactionType]

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"tracking":             meta.Tracking,
			"nectar_module":        "newsfeed_blingbar",
			"feedback_source":      "feedback_comments",
			"feedback_referrer":    "native_newsfeed",
			"feedback_id":          feedbackID,
			"client_mutation_id":   mutationID,
			"attribution_id_v2":    attribution,
			"actor_id":             userID,
			"feedback_reaction_id": reactionID,
			"action_timestamp":     actionTimestamp,
			"page_id":              meta.PageID,
			"content_owner_id_new": meta.ownerID,
		},
	}
	variablesJSON, _ := json.Marshal(variables)

	// 6. ฟอร์มข้อมูล
	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "ViewerReactionsMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "285778409315553568300335455481")
	form.Set("variables", string(variablesJSON))
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)
	form.Set("client_trace_id", traceID)

	// 7. เตรียมการเชื่อมต่อ
	host := "graph.facebook.com"

	// 9. สร้าง HTTP Request
	req, _ := http.NewRequest("POST", "https://graph.facebook.com/graphql", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", deviceGroup)
	req.Header.Set("X-FB-Friendly-Name", "ViewerReactionsMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", meta.PageID)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+traceID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthRunlike_comment_only())
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	// 10. ส่ง Request
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

	// 🔽 บันทึก response ลงตาราง respond_for_comment_table
	_, err = db.Exec("INSERT INTO respond_for_like_comment_only_table (respond_txt) VALUES (?)", string(bodyResp))
	if err != nil {
		fmt.Println("❌ บันทึก response ลงตาราง respond_for_like_comment_only_table ไม่สำเร็จ:", err)
	} else {
		fmt.Println("💾 บันทึก response แล้วลงตาราง respond_for_like_comment_only_table")
	}

	//	_, err = db.Exec("DELETE FROM like_comment_only_table WHERE reaction_type = ?", reactionType) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid) reactionType, link
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ reaction_type ออกจากฐานข้อมูลแล้ว:", reactionType)
	//	}

	_, err = db.Exec("DELETE FROM like_comment_only_table WHERE link = ?", link) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid)
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ link ออกจากฐานข้อมูลแล้ว:", link)
	}

}
