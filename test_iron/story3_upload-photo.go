// cover_pic1up.go (FULL)
// ยิงอัปโหลดรูปภาพไปยัง Facebook ผ่าน /me/photos + uTLS + Proxy + Header spoof + ตรวจ gzip response + ฟิลด์เจนใหม่ครบทุกตัว

package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/fs"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func randomExcellentBandwidthRunstory3_upload_photo() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func findSmallestImagePathRunstory3_upload_photo(folder string) string {
	var smallest string
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".jpg") || strings.HasSuffix(strings.ToLower(d.Name()), ".jpeg") {
			if smallest == "" || d.Name() < smallest {
				smallest = d.Name()
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("❌ WalkDir fail: " + err.Error())
	}
	if smallest == "" {
		fmt.Println("❌ ไม่มีไฟล์ .jpg ในโฟลเดอร์เลย")
	}
	return filepath.Join(folder, smallest)
}

func generateIDRunstory3_upload_photo() string {
	return uuid.New().String()
}

func generateIdempotenceTokenRunstory3_upload_photo() string {
	return fmt.Sprintf("%s_%d", uuid.New().String(), time.Now().UnixNano())
}

func Runstory3_upload_photo(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr)

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}
	functionFolder := filepath.Join(folder, "story_photo")
	imagePath := findSmallestImagePathRunstory3_upload_photo(functionFolder)

	// เพิ่ม debug
	fmt.Println("📁 Story folder:", functionFolder)
	fmt.Println("🖼️ Image path:", imagePath)

	imgFile, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("❌ Open image fail: " + err.Error())
		return
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("❌ Decode image fail: " + err.Error())
		return
	}
	imgFile.Close()

	var imgBuf bytes.Buffer
	if err := jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: 90}); err != nil {
		fmt.Println("❌ Encode JPEG fail: " + err.Error())
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	boundary := writer.Boundary()
	filename := generateIDRunstory3_upload_photo() + ".tmp"

	// ✅ แก้ไข application_tags ให้เหมาะกับ story
	writer.WriteField("published", "false")
	writer.WriteField("audience_exp", "true")
	writer.WriteField("qn", generateIDRunstory3_upload_photo())
	writer.WriteField("composer_session_id", generateIDRunstory3_upload_photo())
	writer.WriteField("idempotence_token", generateIdempotenceTokenRunstory3_upload_photo())
	writer.WriteField("composer_entry_point", "camera_roll")
	writer.WriteField("locale", "en_US")
	writer.WriteField("client_country_code", "TH")
	writer.WriteField("fb_api_req_friendly_name", "upload-photo")
	writer.WriteField("fb_api_caller_class", "MultiPhotoUploader")

	part, _ := writer.CreateFormFile("source", filename)
	part.Write(imgBuf.Bytes())
	writer.Close()

	dbPath := filepath.Join(folder, "fb_comment_system.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
		return
	}
	defer db.Close()

	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	// Debug token
	fmt.Printf("🔑 Token: %s...%s (length: %d)\n",
		accessToken[:10], accessToken[len(accessToken)-10:], len(accessToken))

	host := "graph.facebook.com"

	req, err := http.NewRequest("POST", "https://"+host+"/me/photos", body)
	if err != nil {
		fmt.Println("❌ Request build fail: " + err.Error())
		return
	}

	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthRunstory3_upload_photo())
	req.Header.Set("X-FB-Connection-Quality", "EXCELLENT")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "upload-photo")
	req.Header.Set("X-FB-Photo-Source", "photo_picker")
	req.Header.Set("X-FB-Upload-Phase", "transfer")
	req.Header.Set("X-FB-Photo-Waterfall-ID", generateIDRunstory3_upload_photo())
	req.Header.Set("X-FB-HTTP-Engine", "Liger")

	// ✅ แก้ application_tags ให้เหมาะกับ story
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"STORY_UPLOAD"}`)
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)

	// Debug headers
	fmt.Println("📋 Request Headers:")
	for name, values := range req.Header {
		fmt.Printf("  %s: %s\n", name, values[0])
	}

	fmt.Printf("📊 Body size: %d bytes\n", req.ContentLength)

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

	// ตรวจสอบว่าเป็น HTTP error response ไหม
	if resp.StatusCode != 200 {
		fmt.Printf("❌ HTTP Error %d: %s\n", resp.StatusCode, string(bodyResp))
		return
	}

	if len(bodyResp) == 0 {
		fmt.Println("❌ Empty response - Facebook may have rejected the request")
		fmt.Println("💡 Suggestions:")
		fmt.Println("  - Check if access token is valid")
		fmt.Println("  - Check if account has permission to upload")
		fmt.Println("  - Try uploading to /me/photos with different parameters")
		return
	}

	var jsonResp map[string]interface{}
	if err := json.Unmarshal(bodyResp, &jsonResp); err != nil {
		fmt.Println("❌ แปลง JSON ไม่ได้: " + err.Error())
		fmt.Printf("📄 Raw response: %q\n", string(bodyResp))
		return
	}

	photoID, ok := jsonResp["id"].(string)
	if !ok {
		fmt.Println("⚠️ JSON Response:", jsonResp)
		fmt.Println("❌ ไม่มีฟิลด์ id ใน response หรือไม่ได้เป็น string")
		return
	}

	fmt.Println("📸 Uploaded photo ID:", photoID)

	// ✅ แก้ให้ตรงกับชื่อคอลัมน์ในตาราง
	_, err = db.Exec(`INSERT INTO story_photo_id_table (pic_id) VALUES (?)`, photoID)
	if err != nil {
		fmt.Println("❌ INSERT pic_id ไม่สำเร็จ: " + err.Error())
		return
	}
	fmt.Println("💾 บันทึก story_photo_id_table ลง DB แล้ว:", photoID)

	if err := os.Remove(imagePath); err != nil {
		fmt.Println("⚠️ ลบไฟล์ไม่สำเร็จ:", err)
	} else {
		fmt.Println("🗑️ ลบไฟล์ที่ใช้แล้ว:", imagePath)
	}
}
