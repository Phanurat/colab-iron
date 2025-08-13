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

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func generateFeedbackIDcomment_reel(postID string) string {
	feedbackID := "feedback:" + postID
	return base64.StdEncoding.EncodeToString([]byte(feedbackID))
}

func randomExcellentBandwidthcomment_reel() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func buildEncodedPayloadcomment_reel(actorID, feedbackID, comment string) string {
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

func fetchOwnerIDcomment_reel(objectID, token string) (string, error) {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?fields=from&access_token=%s", objectID, token)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		From struct {
			ID string `json:"id"`
		} `json:"from"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.From.ID == "" {
		return "", fmt.Errorf("ไม่พบ from.id ใน response")
	}
	return result.From.ID, nil
}

func extractIDFromLinkcomment_reel(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	parts := strings.Split(u.Path, "/")
	for _, p := range parts {
		if len(p) > 10 && isNumericcomment_reel(p) {
			return p
		}
	}
	return ""
}

func isNumericcomment_reel(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func Runcomment_reel(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var accessToken, userID, userAgent, netHni, simHni string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	var link, commentText string
	err = db.QueryRow("SELECT link, comment_text FROM like_reel_and_comment_reel_table LIMIT 1").Scan(
		&link, &commentText)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล like_reel_and_comment_reel_table ไม่สำเร็จ: " + err.Error())
		return
	}

	postID := extractIDFromLinkcomment_reel(link)
	if postID == "" {
		fmt.Println("❌ ดึง postID จากลิงก์ไม่สำเร็จ: " + link)
		return
	}

	feedbackID := generateFeedbackIDcomment_reel(postID)
	payload := buildEncodedPayloadcomment_reel(userID, feedbackID, commentText)
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
		"X-FB-Friendly-Name":        {"CommentCreateMutation"},
		"X-FB-Connection-Type":      {"MOBILE.HSDPA"},
		"X-FB-HTTP-Engine":          {"Liger"},
		"x-fb-client-ip":            {"True"},
		"x-fb-server-cluster":       {"True"},
		"x-fb-connection-bandwidth": {randomExcellentBandwidthcomment_reel()},
		"x-fb-connection-quality":   {"EXCELLENT"},
		"x-fb-net-hni":              {netHni},
		"x-fb-sim-hni":              {simHni},
		"x-graphql-client-library":  {"graphservice"},
		"x-tigon-is-retry":          {"False"},
		"x-fb-ta-logging-ids":       {fmt.Sprintf("graphql:%s", uuid.New().String())},
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

	//	_, err = db.Exec("DELETE FROM like_reel_and_comment_reel_table WHERE comment_text = ?", commentText) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid) reactionType, link
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ reaction_type ออกจากฐานข้อมูลแล้ว:", commentText)
	//	}

	_, err = db.Exec("DELETE FROM like_reel_and_comment_reel_table WHERE link = ?", link) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid)
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ link ออกจากฐานข้อมูลแล้ว:", link)
	}

	// 🔽 บันทึก response ลงตาราง respond_for_comment_table
	_, err = db.Exec("INSERT INTO respond_for_comment_reel_table (respond_txt) VALUES (?)", string(bodyResp))
	if err != nil {
		fmt.Println("❌ บันทึก response ลงตาราง respond_for_comment_reel_table ไม่สำเร็จ:", err)
	} else {
		fmt.Println("💾 บันทึก response แล้วลงตาราง respond_for_comment_reel_table")
	}

}
