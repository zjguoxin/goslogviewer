/**
 * @Author: guxline zjguoxin@163.com
 * @Date: 2025/7/4 20:01:07
 * @LastEditors: guxline zjguoxin@163.com
 * @LastEditTime: 2025/7/4 20:01:07
 * Description: 处理器测试
 * Copyright: Copyright (©) 2025 中易综服. All rights reserved.
 */
package goslogviewer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHttpGetFilesHandler(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{LogDir: tempDir}
	lv := New(config)

	// 创建测试文件
	testFiles := []string{"test1.log", "test2.log"}
	for _, filename := range testFiles {
		file, err := os.Create(filepath.Join(tempDir, filename))
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()
	}

	req := httptest.NewRequest("GET", "/log/getLogFilesList", nil)
	w := httptest.NewRecorder()

	lv.GetFilesHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response struct {
		Code  int      `json:"code"`
		Files []string `json:"files"`
		Msg   string   `json:"msg"`
	}
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Files) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(response.Files))
	}
}

func TestHttpGetContentHandler(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{LogDir: tempDir}
	lv := New(config)

	// 创建测试日志文件
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

	req := httptest.NewRequest("GET", "/log/getFileContent?name="+testFile, nil)
	w := httptest.NewRecorder()

	lv.GetContentHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response struct {
		Code int        `json:"code"`
		Data []LogEntry `json:"data"`
		Msg  string     `json:"msg"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Data) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(response.Data))
	}

	if response.Data[0] != testEntry {
		t.Errorf("Entry mismatch: expected %v, got %v", testEntry, response.Data[0])
	}
}

func TestHttpClearFileContentHandler(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{LogDir: tempDir}
	lv := New(config)

	// 创建测试文件
	testFile := "test.log"
	filePath := filepath.Join(tempDir, testFile)
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.WriteString("test content")
	file.Close()

	req := httptest.NewRequest("POST", "/log/clearFileContent", strings.NewReader("name="+testFile))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	lv.ClearFileContentHandler(w, req)

	resp := w.Result()
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
}

func TestHttpExportFileHandler(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{LogDir: tempDir}
	lv := New(config)

	// 创建测试日志文件
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

	// 测试JSON导出
	req := httptest.NewRequest("GET", "/log/exportFile?name="+testFile+"&format=json", nil)
	w := httptest.NewRecorder()

	lv.ExportFileHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", contentType)
	}

	// 测试文本导出
	req = httptest.NewRequest("GET", "/log/exportFile?name="+testFile, nil)
	w = httptest.NewRecorder()

	lv.ExportFileHandler(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "text/plain" {
		t.Errorf("Expected content type text/plain, got %s", contentType)
	}
}

func TestGetFilesHandler_Error(t *testing.T) {
	// 使用不存在的目录强制出错
	config := &Config{LogDir: "/nonexistent/directory"}
	lv := New(config)

	req := httptest.NewRequest("GET", "/log/getLogFilesList", nil)
	w := httptest.NewRecorder()

	lv.GetFilesHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status InternalServerError, got %v", resp.Status)
	}
}

func TestGetContentHandler_InvalidFile(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{LogDir: tempDir}
	lv := New(config)

	// 测试不存在的文件
	req := httptest.NewRequest("GET", "/log/getFileContent?name=nonexistent.log", nil)
	w := httptest.NewRecorder()

	lv.GetContentHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status InternalServerError, got %v", resp.Status)
	}
}

func TestGetContentHandler_MissingParam(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{LogDir: tempDir}
	lv := New(config)

	// 测试缺少文件名参数
	req := httptest.NewRequest("GET", "/log/getFileContent", nil)
	w := httptest.NewRecorder()

	lv.GetContentHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", resp.Status)
	}
}
