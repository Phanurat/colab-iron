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

var reactionIDslike_before_comment = map[string]string{
	"like":  "1635855486666999",
	"love":  "1678524932434102",
	"haha":  "115940658764963",
	"wow":   "478547315650144",
	"sad":   "444813342392137",
	"angry": "604753422931501",
	"care":  "613557422527858",
}

func generateLikeMetalike_before_comment(postID string) (string, string, string, string) {
	feedbackID := "feedback:" + postID
	feedbackIDB64 := base64.StdEncoding.EncodeToString([]byte(feedbackID))
	clientMutationID := uuid.New().String()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	qid := fmt.Sprintf("%d", rand.Int63n(9e18)*-1)
	return feedbackIDB64, clientMutationID, timestamp, qid
}

func randomExcellentBandwidthlike_before_comment() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func extractFacebookIDslike_before_comment(rawurl string) (string, string, error) {
	var postID, ownerID string

	u, err := url.Parse(rawurl)
	if err != nil {
		return "", "", err
	}
	query := u.Query()

	reStory := regexp.MustCompile(`story_fbid=(\d+)`)
	rePath := regexp.MustCompile(`facebook\.com/(\d+)/(?:videos|posts)/(\d+)`)
	reUsername := regexp.MustCompile(`facebook\.com/([^/?&]+)`)

	if match := reStory.FindStringSubmatch(rawurl); len(match) > 1 {
		postID = match[1]
	}
	if match := rePath.FindStringSubmatch(rawurl); len(match) > 2 {
		ownerID = match[1]
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
	if id := query.Get("id"); id != "" {
		ownerID = id
	}
	if ownerID == "" {
		if match := reUsername.FindStringSubmatch(rawurl); len(match) > 1 {
			username := match[1]
			if isNumericlike_before_comment(username) {
				ownerID = username
			} else {
				fbid, err := getFBIDFromUsernamelike_before_comment(username)
				if err != nil {
					return "", "", err
				}
				ownerID = fbid
			}
		}
	}
	return ownerID, postID, nil
}

func isNumericlike_before_comment(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func getFBIDFromUsernamelike_before_comment(username string) (string, error) {
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

// Runlike_before_comment
func Runlike_before_comment(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	rand.Seed(time.Now().UnixNano())

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

	// db, err := sql.Open("sqlite3", "./fb_comment_system.db")
	// if err != nil {
	// 	fmt.Println("❌ เปิดฐานข้อมูลไม่ได้: " + err.Error())
	// 	return
	// }
	// defer db.Close()

	var accessToken, userID, userAgent, netHni, simHni string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni FROM app_profiles LIMIT 1").Scan(&accessToken, &userID, &userAgent, &netHni, &simHni)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	var reactionType, link string
	err = db.QueryRow("SELECT reaction_type, link FROM like_and_comment_table LIMIT 1").Scan(&reactionType, &link)
	if err != nil {
		fmt.Println("❌ ดึงลิงก์จาก like_and_comment_table ไม่สำเร็จ: " + err.Error())
		return
	}

	ownerID, postID, err := extractFacebookIDslike_before_comment(link)
	if err != nil {
		fmt.Println("❌ ขุด owner_id/post_id ไม่สำเร็จ: " + err.Error())
		return
	}

	reactionID := reactionIDslike_before_comment[reactionType]
	feedbackIDB64, clientMutationID, timestamp, qid := generateLikeMetalike_before_comment(postID)

	variables := `{"input":{"tracking":["{\"qid\":\"` + qid + `\",\"mf_story_key\":\"` + postID + `\",\"top_level_post_id\":\"` + postID + `\",\"tl_objid\":\"` + postID + `\",\"content_owner_id_new\":\"` + ownerID + `\",\"throwback_story_fbid\":\"` + postID + `\",\"page_id\":\"1263483477072753\",\"story_location\":4,\"sty\":22,\"page_insights\":{\"` + ownerID + `\":{\"page_id\":\"` + ownerID + `\",\"page_id_type\":\"page\",\"actor_id\":\"` + ownerID + `\",\"dm\":{\"isShare\":0,\"originalPostOwnerID\":0,\"sharedMediaID\":0,\"sharedMediaOwnerID\":0},\"psn\":\"EntStatusCreationStory\",\"post_context\":{\"object_fbtype\":266,\"publish_time\":1747916356,\"story_name\":\"EntStatusCreationStory\",\"story_fbid\":[\"` + postID + `\"]},\"role\":1,\"sl\":4,\"targets\":[{\"actor_id\":\"` + ownerID + `\",\"page_id\":\"` + ownerID + `\",\"post_id\":\"` + postID + `\",\"role\":1,\"share_id\":0}]}},\"profile_id\":\"` + ownerID + `\",\"profile_relationship_type\":3,\"actrs\":\"` + ownerID + `\",\"tds_flgs\":3}","{\"image_loading_state\":0,\"radio_type\":\"wifi-none\",\"client_viewstate_position\":-3}"],"nectar_module":"timeline_ufi","feedback_source":"native_timeline","feedback_referrer":"native_timeline","feedback_id":"` + feedbackIDB64 + `","client_mutation_id":"` + clientMutationID + `","attribution_id_v2":"ProfileFragment,...","actor_id":"` + userID + `","feedback_reaction_id":"` + reactionID + `","action_timestamp":` + timestamp + `}}`

	data := url.Values{}
	data.Set("method", "post")
	data.Set("pretty", "false")
	data.Set("format", "json")
	data.Set("server_timestamps", "true")
	data.Set("locale", "en_US")
	data.Set("fb_api_req_friendly_name", "ViewerReactionsMutation")
	data.Set("fb_api_caller_class", "graphservice")
	data.Set("client_doc_id", "285778409315553568300335455481")
	data.Set("variables", variables)

	encodedBody := data.Encode()
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

	req, err := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBufferString(encodedBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Host = host
	req.ContentLength = int64(len(encodedBody))

	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Friendly-Name", "ViewerReactionsMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthlike_before_comment())
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("X-FB-Background-State", "1")

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

	// 🔽 บันทึก response ลงตาราง respond_for_comment_table
	_, err = db.Exec("INSERT INTO respond_for_like_before_comment_table (respond_txt) VALUES (?)", string(bodyResp))
	if err != nil {
		fmt.Println("❌ บันทึก response ลงตาราง respond_for_like_before_comment_table ไม่สำเร็จ:", err)
	} else {
		fmt.Println("💾 บันทึก response แล้วลงตาราง respond_for_like_before_comment_table")
	}

	//	_, err = db.Exec("DELETE FROM like_before_comment_table WHERE reaction_type = ?", reactionType) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid) reactionType, link
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ reaction_type ออกจากฐานข้อมูลแล้ว:", reactionType)
	//	}

	//	_, err = db.Exec("DELETE FROM like_before_comment_table WHERE link = ?", link) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid)
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ link ออกจากฐานข้อมูลแล้ว:", link)
	//	}

}
