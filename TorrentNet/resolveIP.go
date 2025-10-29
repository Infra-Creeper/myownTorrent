package TorrentNet

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"fmt"
	"io"
)

func GetLocalIP() (string, error) {
	// 1) Try interfaces
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			// skip down or loopback interfaces
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip == nil {
					continue
				}
				// prefer IPv4
				ip4 := ip.To4()
				if ip4 == nil {
					continue
				}
				return ip4.String(), nil
			}
		}
	}

	// 2) Fallback: use UDP dial to determine outbound IP (no packets are sent).
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", errors.New("could not determine local IP address")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	if localAddr.IP == nil {
		return "", errors.New("could not determine local IP address")
	}
	return localAddr.IP.String(), nil
}

// Clients requests the File from server
func RequestFile(serverAddr, fileName, saveAs string) error {
	// Connect to server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("error connecting to server: %v", err)
	}
	defer conn.Close()

	// Send filename length (8 bytes, padded)
	fileNameLenStr := fmt.Sprintf("%08d", len(fileName))
	conn.Write([]byte(fileNameLenStr))

	// Send filename
	conn.Write([]byte(fileName))

	fmt.Printf("Requesting file: %s\n", fileName)

	// Read response header (8 bytes: SUCCESS_ or ERROR___)
	headerBuf := make([]byte, 8)
	_, err = io.ReadFull(conn, headerBuf)
	if err != nil {
		return fmt.Errorf("ERROR reading response header: %v", err)
	}

	header := string(headerBuf)

	if header == "ERROR___" {
		// Read error message length
		errLenBuf := make([]byte, 8)
		_, err = io.ReadFull(conn, errLenBuf)
		if err != nil {
			return fmt.Errorf("ERROR reading error message length: %v", err)
		}
		errLen, _ := strconv.ParseInt(strings.TrimSpace(string(errLenBuf)), 10, 64)

		// Read error message
		errMsgBuf := make([]byte, errLen)
		_, err = io.ReadFull(conn, errMsgBuf)
		if err != nil {
			return fmt.Errorf("ERROR reading error message: %v", err)
		}

		return fmt.Errorf("server error: %s", string(errMsgBuf))
	}

	// Read file size
	fileSizeBuf := make([]byte, 16)
	_, err = io.ReadFull(conn, fileSizeBuf)
	if err != nil {
		return fmt.Errorf("error reading file size: %v", err)
	}
	fileSize, _ := strconv.ParseInt(strings.TrimSpace(string(fileSizeBuf)), 10, 64)

	// Create output file
	if saveAs == "" {
		saveAs = fileName
	}
	outFile, err := os.Create(saveAs)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	// Receive file data
	received := int64(0)
	buf := make([]byte, 4096)
	for received < fileSize {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("error reading data: %v", err)
			}
			break
		}
		outFile.Write(buf[:n])
		received += int64(n)
	}

	fmt.Printf("Downloaded file: %s (%d bytes)\n", saveAs, received)
	return nil
}

// Server serves files from a directory
func StartServer(ipAddr, port, shareDir string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	listener, err := net.Listen("tcp", ipAddr+":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %s, sharing: %s\n", port, shareDir)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// Each request handled in its own goroutine
		go handleRequest(conn, shareDir)
	}
}

func handleRequest(conn net.Conn, shareDir string) {
	defer conn.Close()

	// Read requested filename length
	fileNameLenBuf := make([]byte, 8)
	_, err := io.ReadFull(conn, fileNameLenBuf)
	if err != nil {
		fmt.Println("Error reading filename length:", err)
		return
	}
	fileNameLen, _ := strconv.ParseInt(strings.TrimSpace(string(fileNameLenBuf)), 10, 64)

	// Read requested filename
	fileNameBuf := make([]byte, fileNameLen)
	_, err = io.ReadFull(conn, fileNameBuf)
	if err != nil {
		fmt.Println("Error reading filename:", err)
		return
	}
	requestedFile := string(fileNameBuf)

	fmt.Printf("[SERVER] Request for file: %s\n", requestedFile)

	// Construct full path and clean it to prevent directory traversal
	// filepath.Clean removes .. and . elements
	cleanPath := filepath.Clean(requestedFile)

	// Prevent absolute paths and paths that go outside share directory
	if filepath.IsAbs(cleanPath) || strings.HasPrefix(cleanPath, "..") {
		conn.Write([]byte("ERROR___"))
		errorMsg := "Invalid file path"
		conn.Write([]byte(fmt.Sprintf("%08d", len(errorMsg))))
		conn.Write([]byte(errorMsg))
		fmt.Printf("[SERVER] Invalid path rejected: %s\n", requestedFile)
		return
	}

	filePath := filepath.Join(shareDir, cleanPath)

	// Check if file exists
	file, err := os.Open(filePath)
	if err != nil {
		// Send error response
		conn.Write([]byte("ERROR___"))
		errorMsg := "File not found"
		conn.Write([]byte(fmt.Sprintf("%08d", len(errorMsg))))
		conn.Write([]byte(errorMsg))
		fmt.Printf("[SERVER] File not found: %s\n", requestedFile)
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		conn.Write([]byte("ERROR___"))
		errorMsg := "Cannot read file info"
		conn.Write([]byte(fmt.Sprintf("%08d", len(errorMsg))))
		conn.Write([]byte(errorMsg))
		return
	}

	fileSize := fileInfo.Size()

	// Send success header
	conn.Write([]byte("SUCCESS_"))

	// Send file size (16 bytes, padded)
	fileSizeStr := fmt.Sprintf("%016d", fileSize)
	conn.Write([]byte(fileSizeStr))

	// Send file data
	sent := int64(0)
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("[SERVER] Error reading file:", err)
			return
		}
		conn.Write(buf[:n])
		sent += int64(n)
	}

	fmt.Printf("[SERVER] Sent file: %s (%d bytes)\n", requestedFile, sent)
}
