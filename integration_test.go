/**
 * @Author: guxline zjguoxin@163.com
 * @Date: 2025/7/4 19:57:46
 * @LastEditors: guxline zjguoxin@163.com
 * @LastEditTime: 2025/7/4 19:57:46
 * Description: 集成测试
 * Copyright: Copyright (©) 2025 中易综服. All rights reserved.
 */
package goslogviewer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIntegration(t *testing.T) {
	// 设置测试环境
	tempDir := t.TempDir()
	config := &Config{
		LogDir:       tempDir,
		DevMode:      true,
		EnableDelete: true,
	}
	lv := New(config)

	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/log/getLogFilesList":
			lv.GetFilesHandler(w, r)
		case "/log/getFileContent":
			lv.GetContentHandler(w, r)
		case "/log/clearFileContent":
			lv.ClearFileContentHandler(w, r)
		case "/log/deleteAllFiles":
			lv.DeleteAllFilesHandler(w, r)
		case "/log/exportFile":
			lv.ExportFileHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	// 创建测试文件
	testFile := "test.log"
	filePath := filepath.Join(tempDir, testFile)
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	testEntry := LogEntry{Level: "INFO", Time: "2023-01-01T00:00:00Z", Msg: "Test message"}
	entryBytes, _ := json.Marshal(testEntry)
	file.Write(entryBytes)
	file.WriteString("\n")

	// 测试获取文件列表
	resp, err := http.Get(ts.URL + "/log/getLogFilesList")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var filesResponse struct {
		Code  int      `json:"code"`
		Files []string `json:"files"`
	}
	err = json.NewDecoder(resp.Body).Decode(&filesResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(filesResponse.Files) != 1 || filesResponse.Files[0] != testFile {
		t.Errorf("Expected file %s, got %v", testFile, filesResponse.Files)
	}

	// 测试获取文件内容
	resp, err = http.Get(ts.URL + "/log/getFileContent?name=" + testFile)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var contentResponse struct {
		Code int        `json:"code"`
		Data []LogEntry `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&contentResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(contentResponse.Data) != 1 || contentResponse.Data[0] != testEntry {
		t.Errorf("Entry mismatch: expected %v, got %v", testEntry, contentResponse.Data[0])
	}

	// 测试清空文件内容
	resp, err = http.Post(ts.URL+"/log/clearFileContent", "application/x-www-form-urlencoded", bytes.NewBufferString("name="+testFile))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// 验证文件是否为空
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty file, got %d bytes", len(content))
	}

	// 测试删除所有文件
	resp, err = http.Post(ts.URL+"/log/deleteAllFiles", "application/x-www-form-urlencoded", nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// 验证目录是否为空
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected empty directory, found %d files", len(entries))
	}
}
