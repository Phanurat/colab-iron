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

func randomExcellentBandwidthRunup_pic_caption() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func findSmallestImagePathRunup_pic_caption(folder string) string {
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

func Runup_pic_caption(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}
	profileFolder := filepath.Join(folder, "caption_photo")
	imagePath := findSmallestImagePathRunup_pic_caption(profileFolder)

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

	imgFile, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("❌ Open image fail: " + err.Error())
		return
	}

	// Decode → Encode
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("❌ Decode image fail: " + err.Error())
		return
	}

	// ✅ Close ทันทีหลังใช้เสร็จ
	imgFile.Close()

	var imgBuf bytes.Buffer
	if err := jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: 90}); err != nil {
		fmt.Println("❌ Encode JPEG fail: " + err.Error())
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	boundary := writer.Boundary()
	filename := uuid.New().String() + ".tmp"

	writer.WriteField("published", "false")
	writer.WriteField("audience_exp", "true")
	writer.WriteField("qn", uuid.New().String())
	writer.WriteField("composer_session_id", uuid.New().String())
	writer.WriteField("idempotence_token", uuid.New().String()+"_"+fmt.Sprint(time.Now().Unix()))
	writer.WriteField("composer_entry_point", "camera_roll")
	writer.WriteField("locale", "en_US")
	writer.WriteField("client_country_code", "TH")
	writer.WriteField("fb_api_req_friendly_name", "upload-photo")
	writer.WriteField("fb_api_caller_class", "MultiPhotoUploader")

	part, _ := writer.CreateFormFile("source", filename)
	part.Write(imgBuf.Bytes())
	writer.Close()

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
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthRunup_pic_caption())
	req.Header.Set("X-FB-Connection-Quality", "EXCELLENT")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "upload-photo")
	req.Header.Set("X-FB-Photo-Source", "photo_picker")
	req.Header.Set("X-FB-Upload-Phase", "transfer")
	req.Header.Set("X-FB-Photo-Waterfall-ID", uuid.New().String())
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"PROFILE_PIC"}`)
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)

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

	var jsonResp map[string]interface{}
	if err := json.Unmarshal(bodyResp, &jsonResp); err != nil {
		fmt.Println("❌ แปลง JSON ไม่ได้: " + err.Error())
		return
	}

	photoID, ok := jsonResp["id"].(string)
	if !ok {
		fmt.Println("❌ ไม่มีฟิลด์ id ใน response หรือไม่ได้เป็น string")
		return
	}

	fmt.Println("📸 Uploaded photo ID:", photoID)

	_, err = db.Exec(`DELETE FROM pic_caption_table`)
	if err != nil {
		fmt.Println("❌ ลบแถวเดิมไม่สำเร็จ: " + err.Error())
		return
	}

	_, err = db.Exec(`INSERT INTO pic_caption_table (media_id) VALUES (?)`, photoID)
	if err != nil {
		fmt.Println("❌ INSERT แถวใหม่ไม่สำเร็จ: " + err.Error())
		return
	}

	fmt.Println("💾 บันทึก media_id ลง DB แล้ว:", photoID)

	if err := os.Remove(imagePath); err != nil {
		fmt.Println("⚠️ ลบไฟล์ไม่สำเร็จ:", err)
	} else {
		fmt.Println("🗑️ ลบไฟล์ที่ใช้แล้ว:", imagePath)
	}
}
