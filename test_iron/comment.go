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

func buildEncodedPayloadcomment(actorID, feedbackID, comment string) string {
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"actor_id":    actorID,
			"message":     map[string]string{"text": comment},
			"feedback_id": feedbackID,
		},
	}
	jsonVars, _ := json.Marshal(variables)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "CommentCreateMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "847448985557369682546426351")
	form.Set("variables", string(jsonVars))

	return form.Encode()
}

func generateFeedbackIDcomment(postID string) string {
	feedbackID := "feedback:" + postID
	return base64.StdEncoding.EncodeToString([]byte(feedbackID))
}

func randomExcellentBandwidthcomment() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(15000000) + 20000000)
}

func extractPostIDcomment(rawurl string) (string, error) {
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

func isNumericcomment(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func getFBIDFromUsernamecomment(username string) (string, error) {
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

func Runcomment(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	var commentText, postLink string
	err = db.QueryRow("SELECT comment_text, link FROM like_and_comment_table LIMIT 1").Scan(&commentText, &postLink)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล like_and_comment_table ไม่สำเร็จ: " + err.Error())
		return
	}

	postID, err := extractPostIDcomment(postLink)
	if err != nil {
		fmt.Println("❌ ขุด post_id ไม่สำเร็จ: " + err.Error())
		return
	}

	feedbackID := generateFeedbackIDcomment(postID)
	payload := buildEncodedPayloadcomment(userID, feedbackID, commentText)

	host := "graph.facebook.com"
	//	address := host + ":443"

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

	req, err := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBufferString(payload))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header = map[string][]string{
		"Authorization":             {"OAuth " + accessToken},
		"Accept-Encoding":           {"gzip, deflate"},
		"Connection":                {"keep-alive"},
		"Host":                      {host},
		"Content-Type":              {"application/x-www-form-urlencoded"},
		"User-Agent":                {userAgent},
		"x-fb-device-group":         {devicegroup},
		"X-FB-Friendly-Name":        {"CommentCreateMutation"},
		"X-FB-Connection-Type":      {"MOBILE.HSDPA"},
		"X-FB-HTTP-Engine":          {"Liger"},
		"x-fb-client-ip":            {"True"},
		"x-fb-server-cluster":       {"True"},
		"x-fb-connection-bandwidth": {randomExcellentBandwidthcomment()},
		"x-fb-connection-quality":   {"EXCELLENT"},
		"x-fb-net-hni":              {netHni},
		"x-fb-sim-hni":              {simHni},
		"x-graphql-client-library":  {"graphservice"},
		"x-tigon-is-retry":          {"False"},
		"x-fb-ta-logging-ids":       {fmt.Sprintf("graphql:%s", uuid.New().String())},
	}
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

	//	_, err = db.Exec("DELETE FROM like_and_comment_table WHERE comment_text = ?", commentText) // commentText, postLink
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ like_and_comment_table ออกจากฐานข้อมูลแล้ว:", commentText)
	//	}

	_, err = db.Exec("DELETE FROM like_and_comment_table WHERE link = ?", postLink) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ like_and_comment_table ออกจากฐานข้อมูลแล้ว:", postLink)
	}

	// 🔽 บันทึก response ลงตาราง respond_for_comment_table
	_, err = db.Exec("INSERT INTO respond_for_comment_table (respond_txt) VALUES (?)", string(bodyResp))
	if err != nil {
		fmt.Println("❌ บันทึก response ลงตาราง respond_for_comment_table ไม่สำเร็จ:", err)
	} else {
		fmt.Println("💾 บันทึก response แล้วลงตาราง respond_for_comment_table")
	}

}
