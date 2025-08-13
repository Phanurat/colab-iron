package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

// ✅ กำหนด binary ตาม OS (Linux ไม่มี .exe)
func targetBinaryName() string {
	name := "main333111"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

// ✅ อ่านค่า MAX_PROCS จาก ENV (ค่าเริ่มต้น 32)
func maxParallel() int {
	if v := os.Getenv("MAX_PROCS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 32 // ปรับตามกำลังเครื่อง/Colab ได้
}

func main() {
	total := 400
	bin := targetBinaryName()

	// ✅ หา path โฟลเดอร์ของ runner แล้วชี้ไปยัง binary ข้างเคียง
	self, err := os.Executable()
	if err != nil {
		fmt.Println("❌ cannot resolve runner executable path:", err)
		return
	}
	baseDir := filepath.Dir(self)
	exePath := filepath.Join(baseDir, bin)

	// ✅ เช็คว่ามีไฟล์เป้าหมายจริงไหม
	if st, err := os.Stat(exePath); err != nil || st.IsDir() {
		fmt.Printf("❌ target binary not found: %s (build สำหรับ Linux แล้วอัปโหลดให้เรียบร้อยก่อน)\n", exePath)
		return
	}

	fmt.Printf("⚡ STARTING %d PROCESSES (parallel=%d)\n", total, maxParallel())

	// ✅ semaphore จำกัดจำนวนโปรเซสพร้อมกัน
	sem := make(chan struct{}, maxParallel())
	var wg sync.WaitGroup

	for i := 1; i <= total; i++ {
		dbFolder := fmt.Sprintf("acc%03d", i)
		wg.Add(1)

		sem <- struct{}{} // เข้าคิว
		go func(db string) {
			defer wg.Done()
			defer func() { <-sem }() // ออกคิว

			runProcess(exePath, db)
		}(dbFolder)
	}

	wg.Wait()
	fmt.Println("✅ All processes finished.")
}

func runProcess(exePath, dbFolder string) {
	fmt.Printf("🟢 Spawning process for %-8s → %s\n", dbFolder, filepath.Base(exePath))

	cmd := exec.Command(exePath)
	// ✅ ส่งต่อ ENV เดิม + ตั้ง DBFOLDER ให้แต่ละโปรเซส
	cmd.Env = append(os.Environ(), "DBFOLDER="+dbFolder)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ [%-8s] failed to start: %v\n", dbFolder, err)
		return
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("⚠️  [%-8s] exited with error: %v\n", dbFolder, err)
	} else {
		fmt.Printf("✅ [%-8s] completed\n", dbFolder)
	}
}
