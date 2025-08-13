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

func extractCommentIDcomment_comment(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	commentID := q.Get("comment_id")
	if commentID != "" {
		return commentID, nil
	}
	return "", fmt.Errorf("❌ ไม่พบ comment_id")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func generateFeedbackIDcomment_comment(postID string) string {
	feedbackID := "feedback:" + postID
	return base64.StdEncoding.EncodeToString([]byte(feedbackID))
}

func randomExcellentBandwidthcomment_comment() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(15000000) + 20000000)
}

func extractPostIDcomment_comment(rawurl string) (string, error) {
	var postID string

	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	reStory := regexp.MustCompile(`story_fbid=(\d+)`)
	rePath := regexp.MustCompile(`facebook\.com/(\d+)/(?:videos|posts)/(\d+)`)

	if match := reStory.FindStringSubmatch(rawurl); len(match) > 1 {
		postID = match[1]
	}
	if match := rePath.FindStringSubmatch(rawurl); len(match) > 2 {
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

func isNumericcomment_comment(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func getFBIDFromUsernamecomment_comment(username string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", "https://mbasic.facebook.com/"+username, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if strings.HasPrefix(location, "intent://profile/") {
		re := regexp.MustCompile(`intent://profile/(\d+)`)
		match := re.FindStringSubmatch(location)
		if len(match) > 1 {
			return match[1], nil
		}
	}

	resp, err = http.Get("https://mbasic.facebook.com/" + username)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	text := string(body)

	re := regexp.MustCompile(`owner_id=(\d+)`)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1], nil
	}
	re = regexp.MustCompile(`profile\.php\?id=(\d+)`)
	match = re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("❌ ไม่พบ owner_id จาก username")
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func Runcomment_comment(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var accessToken, actorID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &actorID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	var commentText, postLink string
	err = db.QueryRow("SELECT comment_text, link FROM like_comment_and_reply_comment_table LIMIT 1").Scan(&commentText, &postLink)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล like_comment_and_reply_comment_table ไม่สำเร็จ: " + err.Error())
		return
	}

	commentID, _ := extractCommentIDcomment_comment(postLink)
	var feedbackID string

	if commentID != "" {
		// ถ้ามี comment_id → แปลว่า reply
		feedbackID = generateFeedbackIDcomment_comment(commentID)
	} else {
		// ไม่มี → ใช้ postID แบบเดิม
		postID, err := extractPostIDcomment_comment(postLink)
		if err != nil {
			fmt.Println("❌ ขุด post_id ไม่สำเร็จ: " + err.Error())
			return
		}
		feedbackID = generateFeedbackIDcomment_comment(postID)
	}

	host := "graph.facebook.com"

	idempotenceToken := uuid.New().String()
	clientMutationID := uuid.New().String()
	attributionID := generateAttributionIDV2comment_comment()

	input := map[string]interface{}{
		"input": map[string]interface{}{
			"actor_id":           actorID,
			"feedback_id":        feedbackID,
			"message":            map[string]interface{}{"text": commentText},
			"idempotence_token":  idempotenceToken,
			"client_mutation_id": clientMutationID,
			"attribution_id_v2":  attributionID,
			"feedback_source":    "feedback_comments",
			"nectar_module":      "feed_inline_comment_composer",
			"entry_point":        "TAP_FEED_INLINE_COMMENT_COMPOSER",
		},
	}

	payload := map[string]string{
		"access_token":             accessToken,
		"fb_api_caller_class":      "graphservice",
		"fb_api_req_friendly_name": "CommentCreateMutation",
		"client_doc_id":            "847448985557369682546426351",
		"server_timestamps":        "true",
		"locale":                   "en_US",
		"variables":                encodeJSONcomment_comment(input),
	}

	formBody := encodeFormcomment_comment(payload)
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	_, _ = gz.Write([]byte(formBody))
	gz.Close()

	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", &compressed)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Friendly-Name", "CommentCreateMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-graphql-request-purpose", "fetch")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthcomment_comment())
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

	//	_, err = db.Exec("DELETE FROM like_comment_and_reply_comment_table WHERE comment_text = ?", commentText) // commentText, postLink
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ like_comment_and_reply_comment_table ออกจากฐานข้อมูลแล้ว:", commentText)
	//	}

	_, err = db.Exec("DELETE FROM like_comment_and_reply_comment_table WHERE link = ?", postLink) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ like_comment_and_reply_comment_table ออกจากฐานข้อมูลแล้ว:", postLink)
	}

	// 🔽 บันทึก response ลงตาราง respond_for_comment_table
	_, err = db.Exec("INSERT INTO respond_for_comment_comment_table (respond_txt) VALUES (?)", string(bodyResp))
	if err != nil {
		fmt.Println("❌ บันทึก response ลงตาราง respond_for_comment_comment_table ไม่สำเร็จ:", err)
	} else {
		fmt.Println("💾 บันทึก response แล้วลงตาราง respond_for_comment_comment_table")
	}

}

func encodeJSONcomment_comment(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func encodeFormcomment_comment(data map[string]string) string {
	var buf bytes.Buffer
	for k, v := range data {
		buf.WriteString(fmt.Sprintf("%s=%s&", k, urlEncodecomment_comment(v)))
	}
	return buf.String()[:buf.Len()-1]
}

func urlEncodecomment_comment(s string) string {
	return (&url.URL{Path: s}).EscapedPath()
}

func generateAttributionIDV2comment_comment() string {
	t1 := fmt.Sprintf("%.2f", float64(time.Now().Unix())+rand.Float64())
	t2 := fmt.Sprintf("%.3f", float64(time.Now().Unix())+rand.Float64())
	s1 := rand.Intn(89999999) + 10000000
	s2 := rand.Intn(899999999) + 100000000

	part1 := fmt.Sprintf("tap_feed_inline_comment_composer,%s,%d,,,", t1, s1)
	part2 := fmt.Sprintf("NewsFeedFragment,native_newsfeed,tap_top_jewel_bar,%s,%d,%d,36#301,%d", t2, s2)

	return part1 + ";" + part2
}
