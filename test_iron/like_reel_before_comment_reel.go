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

var reactionIDslike_reel_before_comment_reel = map[string]string{
	"like":  "1635855486666999",
	"love":  "1678524932434102",
	"haha":  "115940658764963",
	"wow":   "478547315650144",
	"sad":   "444813342392137",
	"angry": "604753422931501",
	"care":  "613557422527858",
}

func generateLikeMetalike_reel_before_comment_reel(postID string) (string, string, string, string) {
	feedbackID := "feedback:" + postID
	feedbackIDB64 := base64.StdEncoding.EncodeToString([]byte(feedbackID))
	clientMutationID := uuid.New().String()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	qid := fmt.Sprintf("%d", rand.Int63n(9e18)*-1)
	return feedbackIDB64, clientMutationID, timestamp, qid
}

func randomExcellentBandwidthlike_reel_before_comment_reel() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000 // 20 Mbps
	max := 35000000 // 35 Mbps
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func fetchOwnerIDlike_reel_before_comment_reel(objectID, token string) (string, error) {
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

func extractIDFromLinklike_reel_before_comment_reel(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	parts := strings.Split(u.Path, "/")
	for _, p := range parts {
		if len(p) > 10 && isNumericlike_reel_before_comment_reel(p) {
			return p
		}
	}
	return ""
}

func isNumericlike_reel_before_comment_reel(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func Runlike_reel_before_comment_reel(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var accessToken, userID, userAgent, netHni, simHni string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return

	}

	var reactionType, link string
	err = db.QueryRow("SELECT reaction_type, link FROM like_reel_and_comment_reel_table LIMIT 1").Scan(
		&reactionType, &link)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล like_reel_and_comment_reel_table ไม่สำเร็จ: " + err.Error())
		return

	}

	postID := extractIDFromLinklike_reel_before_comment_reel(link)
	if postID == "" {
		fmt.Println("❌ ดึง postID จากลิงก์ไม่สำเร็จ: " + link)
		return

	}

	ownerID, err := fetchOwnerIDlike_reel_before_comment_reel(postID, accessToken)
	if err != nil {
		fmt.Println("❌ ดึง owner_id ไม่สำเร็จ: " + err.Error())
		return

	}

	reactionID := reactionIDslike_reel_before_comment_reel[reactionType]
	feedbackIDB64, clientMutationID, timestamp, qid := generateLikeMetalike_reel_before_comment_reel(postID)

	variables := `{"input":{"tracking":["{\"qid\":\"` + qid + `\",\"mf_story_key\":\"` + postID + `\",\"top_level_post_id\":\"` + postID + `\",\"tl_objid\":\"` + postID + `\",\"content_owner_id_new\":\"` + ownerID + `\",\"throwback_story_fbid\":\"` + postID + `\",\"page_id\":\"1263483477072753\",\"story_location\":4,\"sty\":22,\"page_insights\":{\"` + ownerID + `\":{\"page_id\":\"` + ownerID + `\",\"page_id_type\":\"page\",\"actor_id\":\"` + ownerID + `\",\"dm\":{\"isShare\":0,\"originalPostOwnerID\":0,\"sharedMediaID\":0,\"sharedMediaOwnerID\":0},\"psn\":\"EntStatusCreationStory\",\"post_context\":{\"object_fbtype\":266,\"publish_time\":1747916356,\"story_name\":\"EntStatusCreationStory\",\"story_fbid\":[\"` + postID + `\"]},\"role\":1,\"sl\":4,\"targets\":[{\"actor_id\":\"` + ownerID + `\",\"page_id\":\"` + ownerID + `\",\"post_id\":\"` + postID + `\",\"role\":1,\"share_id\":0}]}},\"profile_id\":\"` + ownerID + `\",\"profile_relationship_type\":3,\"actrs\":\"` + ownerID + `\",\"tds_flgs\":3}","{\"image_loading_state\":0,\"radio_type\":\"wifi-none\",\"client_viewstate_position\":-3,\"feed_unit_trace_info\":\"{\\\"groups_tab_feed_unit_type_name\\\":\\\"Story\\\",\\\"groups_tab_story_is_crf\\\":2}\",\"feed_unit_trace_info\":\"{}\"}"],"nectar_module":"timeline_ufi","feedback_source":"native_timeline","feedback_referrer":"native_timeline","feedback_id":"` + feedbackIDB64 + `","client_mutation_id":"` + clientMutationID + `","attribution_id_v2":"ProfileFragment,...","actor_id":"` + userID + `","feedback_reaction_id":"` + reactionID + `","action_timestamp":` + timestamp + `}}`

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
	// address := host + ":443"

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
		panic(err)
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
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthlike_reel_before_comment_reel())
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

	//	_, err = db.Exec("DELETE FROM like_reel_before_comment_reel_table WHERE reaction_type = ?", reactionType) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid) reactionType, link
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ reaction_type ออกจากฐานข้อมูลแล้ว:", reactionType)
	//	}

	//	_, err = db.Exec("DELETE FROM like_reel_before_comment_reel_table WHERE link = ?", link) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid)
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ link ออกจากฐานข้อมูลแล้ว:", link)
	//	}

}
