package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	fmt.Println("⚡ STARTING ALL 400 PROCESSES SIMULTANEOUSLY")
	
	// สร้างโฟลเดอร์ logs ถ้ายังไม่มี
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 1; i <= 400; i++ {
		dbFolder := fmt.Sprintf("acc%03d", i)
		wg.Add(1)

		go func(db string) {
			defer wg.Done()
			fmt.Printf("🚀 Launching %s\n", db)
			runProcess(db)
		}(dbFolder)
	}

	fmt.Println("⏳ Waiting for all processes to complete...")
	wg.Wait()
	
	totalTime := time.Since(startTime)
	fmt.Printf("✅ ALL 400 PROCESSES FINISHED in %v\n", totalTime)
}

func runProcess(dbFolder string) {
	// หา path จริงของ exe
	exePath, err := filepath.Abs("main333111.exe")
	if err != nil {
		fmt.Printf("❌ [%-12s] failed to resolve exe path: %v\n", dbFolder, err)
		return
	}

	// สร้างไฟล์ log แยกสำหรับแต่ละ process
	logFile, err := os.Create(fmt.Sprintf("logs/%s.log", dbFolder))
	if err != nil {
		fmt.Printf("❌ [%-12s] failed to create log file: %v\n", dbFolder, err)
		return
	}
	defer logFile.Close()

	cmd := exec.Command(exePath)
	cmd.Env = append(os.Environ(), "DBFOLDER="+dbFolder)
	
	// แยก output ของแต่ละ process ไปยัง log file แยก
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	processStart := time.Now()
	
	err = cmd.Start()
	if err != nil {
		fmt.Printf("❌ [%-12s] failed to start: %v\n", dbFolder, err)
		return
	}

	fmt.Printf("✅ [%-12s] started with PID %d\n", dbFolder, cmd.Process.Pid)

	err = cmd.Wait()
	duration := time.Since(processStart)
	
	if err != nil {
		fmt.Printf("⚠️  [%-12s] exited with error after %v: %v\n", dbFolder, duration, err)
	} else {
		fmt.Printf("🎉 [%-12s] completed successfully in %v\n", dbFolder, duration)
	}
}