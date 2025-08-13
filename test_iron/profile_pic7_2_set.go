package main

import (
	"compress/gzip"
	"database/sql"
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

func genUUIDRunprofile_pic7_2_set() string {
	return uuid.NewString()
}

func generateVisitationIDRunprofile_pic7_2_set(userID string) string {
	shard := randomHexRunprofile_pic7_2_set(5)
	t1 := float64(time.Now().UnixNano()) / 1e9
	t2 := t1 - rand.Float64()*2
	hash := randomHexRunprofile_pic7_2_set(5)

	return fmt.Sprintf("User:%s:1:%.3f|%s:%s:1:%.3f", shard, t1, userID, hash, t2)
}

func generateSessionIDRunprofile_pic7_2_set() string {
	return fmt.Sprintf("UFS-%s-fg-%d", uuid.New().String(), rand.Intn(4)+1)
}

func randomHexRunprofile_pic7_2_set(n int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = hex[rand.Intn(len(hex))]
	}
	return string(b)
}

func randomExcellentBandwidthRunprofile_pic7_2_set() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// func getLatestPhotoIDRunprofile_pic7_2_set() string {
// 	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
// 	if folder == "" {
// 		folder = "."
// 	}

// 	dbPath := filepath.Join(folder, "photo_id.db")
// 	fmt.Println("📂 DB PATH:", dbPath)

// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
// 	}
// 	defer db.Close()

// 	fmt.Println("📂 DB PATH:", folder+"/photo_id.db")

// 	row := db.QueryRow("SELECT pic_id FROM photo_id_table ORDER BY id DESC LIMIT 1")
// 	var photoID string
// 	err = row.Scan(&photoID)
// 	if err != nil {
// 		fmt.Println("❌ ดึง pic_id ไม่สำเร็จ: " + err.Error())
// 	}
// 	return photoID
// }

func Runprofile_pic7_2_set(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	host := "graph.facebook.com"

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

	// Get latest photo ID
	//photoID := getLatestPhotoIDRunprofile_pic7_2_set()

	// Get profile data from database
	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	var photoID string
	err = db.QueryRow("SELECT pic_id FROM profile_photo_id_table LIMIT 1").Scan(
		&photoID)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล profile_photo_id_table ไม่สำเร็จ: " + err.Error())
		return
	}

	variables := fmt.Sprintf(`{
  "input":{
    "suppress_stories":false,
    "set_profile_photo_shield":"TURN_OFF",
    "scaled_crop_rect":{"y":0,"width":1,"x":0,"height":1},
    "composer_session_id":"%s",
    "profile_pic_source":"UNKNOWN",
    "client_mutation_id":"%s",
    "profile_id":"%s",
    "has_umg":false,
    "existing_photo_id":"%s",
    "frame_entrypoint":"camera_roll",
    "profile_pic_method":"unknown",
    "actor_id":"%s"
  }
}`,
		genUUIDRunprofile_pic7_2_set(), genUUIDRunprofile_pic7_2_set(), userID, photoID, userID)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "ProfilePictureSetMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "18097093386471407977628943042")
	form.Set("variables", variables)
	visitationID := generateVisitationIDRunprofile_pic7_2_set(userID)
	sessionID := generateSessionIDRunprofile_pic7_2_set()
	form.Set("fb_api_analytics_tags", fmt.Sprintf(`["visitation_id=%s","session_id=%s","GraphServices"]`, visitationID, sessionID))

	req, err := http.NewRequest("POST", "https://"+host+"/graphql", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println("❌ Request build fail: " + err.Error())
		return
	}

	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "ProfilePictureSetMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)                                                       // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                                       // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthRunprofile_pic7_2_set()) //เพิ่มเข้าไป
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

	// Clean up photo_id_table after successful operation
	_, err = db.Exec("DELETE FROM profile_photo_id_table WHERE pic_id = ?", photoID)
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ profile_photo_id_table ออกจากฐานข้อมูลแล้ว:", photoID)
	}

}
