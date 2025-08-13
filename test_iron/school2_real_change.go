package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// func randomExcellentBandwidthschool2_real_change() string {
// 	rand.Seed(time.Now().UnixNano())
// 	min := 20000000
// 	max := 35000000
// 	return strconv.Itoa(rand.Intn(max-min+1) + min)
// }

func Runschool2_real_change(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	//	host := "graph.facebook.com"

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

	var schoolName string
	err = db.QueryRow("SELECT school_name FROM change_school_table LIMIT 1").Scan(
		&schoolName)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล change_school_table ไม่สำเร็จ: " + err.Error())
		return
	}

	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"client_mutation_id":      "a1",
			"actor_id":                userID,
			"concentration_id":        nil,
			"concentration_name":      nil,
			"experience_id":           nil,
			"has_graduated":           false,
			"life_event_publish_type": "SUPPRESS_ALL",
			"privacy": map[string]interface{}{
				"allow":               []string{},
				"base_state":          "EVERYONE",
				"deny":                []string{},
				"tag_expansion_state": "UNSPECIFIED",
			},
			"school_id":        nil, // <--- จุดสำคัญ คือปล่อยให้ null
			"school_name":      schoolName,
			"school_type":      "hs",
			"start":            map[string]interface{}{},
			"end":              map[string]interface{}{},
			"ref":              "react_native_form",
			"mutation_surface": "PROFILE",
			"session_id":       uuid.New().String(),
		},
	}

	payload := map[string]string{
		"access_token":             accessToken,
		"fb_api_caller_class":      "RelayModern",
		"fb_api_req_friendly_name": "ProfileEditEducationExperienceSaveMutation",
		"variables":                encodeJSONschool2_real_change(vars),
		"server_timestamps":        "true",
		"doc_id":                   "2228867157143096", // doc ของ Education Mutation
	}

	body := encodeFormschool2_real_change(payload)
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	_, _ = gz.Write([]byte(body))
	gz.Close()

	req, _ := http.NewRequest("POST", "https://graph.facebook.com/graphql?locale=en_US", &compressed)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-fb-friendly-name", "ProfileEditEducationExperienceSaveMutation")

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

	_, err = db.Exec("DELETE FROM change_school_table WHERE school_name = ?", schoolName) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ change_school_table ออกจากฐานข้อมูลแล้ว:", schoolName)
	}

}

func encodeJSONschool2_real_change(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func encodeFormschool2_real_change(data map[string]string) string {
	var buf bytes.Buffer
	for k, v := range data {
		buf.WriteString(fmt.Sprintf("%s=%s&", k, urlEncodeschool2_real_change(v)))
	}
	return buf.String()[:buf.Len()-1]
}

func urlEncodeschool2_real_change(s string) string {
	return (&url.URL{Path: s}).EscapedPath()
}
