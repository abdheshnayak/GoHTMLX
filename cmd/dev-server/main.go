package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// DevServer provides a simple development server with auto-restart
type DevServer struct {
	dir     string
	cmd     []string
	process *exec.Cmd
	ctx     context.Context
	cancel  context.CancelFunc
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Println("GoHTMLX Development Server")
		fmt.Println("Usage: go run cmd/dev-server/main.go <directory> [command...]")
		fmt.Println("Example: go run cmd/dev-server/main.go ./example go run .")
		fmt.Println("")
		fmt.Println("This tool watches for Go file changes and automatically restarts your server.")
		os.Exit(0)
	}

	dir := os.Args[1]
	cmd := []string{"go", "run", "."}
	if len(os.Args) > 2 {
		cmd = os.Args[2:]
	}

	server := &DevServer{
		dir: dir,
		cmd: cmd,
	}

	server.ctx, server.cancel = context.WithCancel(context.Background())

	// Handle shutdown gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n🛑 Shutting down development server...")
		server.stop()
		os.Exit(0)
	}()

	fmt.Printf("🚀 Starting development server in %s\n", dir)
	fmt.Printf("📝 Command: %s\n", strings.Join(cmd, " "))
	fmt.Println("👀 Watching for Go file changes...")
	fmt.Println("Press Ctrl+C to stop")

	server.run()
}

func (s *DevServer) run() {
	// Start initial process
	s.start()

	// Watch for file changes
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastModTime := s.getLastModTime()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			currentModTime := s.getLastModTime()
			if currentModTime.After(lastModTime) {
				fmt.Println("📝 Go files changed, restarting server...")
				s.restart()
				lastModTime = currentModTime
			}
		}
	}
}

func (s *DevServer) start() {
	s.process = exec.CommandContext(s.ctx, s.cmd[0], s.cmd[1:]...)
	s.process.Dir = s.dir
	s.process.Stdout = os.Stdout
	s.process.Stderr = os.Stderr

	err := s.process.Start()
	if err != nil {
		log.Printf("❌ Failed to start process: %v", err)
		return
	}

	fmt.Printf("✅ Server started (PID: %d)\n", s.process.Process.Pid)

	// Wait for process in background
	go func() {
		err := s.process.Wait()
		if err != nil && s.ctx.Err() == nil {
			log.Printf("⚠️  Process exited: %v", err)
		}
	}()
}

func (s *DevServer) stop() {
	if s.process != nil && s.process.Process != nil {
		fmt.Printf("🛑 Stopping server (PID: %d)\n", s.process.Process.Pid)
		s.process.Process.Kill()
		s.process.Wait()
	}
}

func (s *DevServer) restart() {
	s.stop()
	time.Sleep(500 * time.Millisecond) // Brief pause before restart
	s.start()
}

func (s *DevServer) getLastModTime() time.Time {
	var lastMod time.Time

	err := filepath.Walk(s.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip generated files and vendor directories
		if strings.Contains(path, "/dist/") ||
			strings.Contains(path, "/vendor/") ||
			strings.Contains(path, "/.git/") {
			return nil
		}

		if info.ModTime().After(lastMod) {
			lastMod = info.ModTime()
		}

		return nil
	})

	if err != nil {
		log.Printf("Error scanning files: %v", err)
	}

	return lastMod
}
