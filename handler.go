/**
 * @Author: guxline zjguoxin@163.com
 * @Date: 2025/7/3 07:16:52
 * @LastEditors: guxline zjguoxin@163.com
 * @LastEditTime: 2025/7/3 07:16:52
 * Description: HTTP处理器
 * Copyright: Copyright (©) 2025 中易综服. All rights reserved.
 */
package goslogviewer

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type LogViewer struct {
	config *Config
}

func (lv *LogViewer) GetConfig() *Config {
	return lv.config
}

func New(config *Config) *LogViewer {
	if config == nil {
		config = DefaultConfig()
	}

	return &LogViewer{config: config}
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetFilesHandler 获取文件列表处理器
func (lv *LogViewer) GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := lv.GetLogFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, map[string]interface{}{
		"code":  200,
		"files": files,
		"msg":   "success",
	})
}

// GetContentHandler 获取日志内容处理器
func (lv *LogViewer) GetContentHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("name")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	logs, err := lv.GetLogContent(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]interface{}{
		"code": 200,
		"data": logs,
		"msg":  "success",
	})
}

// ClearFileContentHandler 清空文件内容
func (lv *LogViewer) ClearFileContentHandler(w http.ResponseWriter, r *http.Request) {
	if !lv.config.DevMode || !lv.config.EnableClear {
		respondJSON(w, map[string]interface{}{
			"code":  3001,
			"files": nil,
			"msg":   "clear operation is disabled",
		})
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	filename := r.PostForm.Get("name")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}
	err := lv.ClearFileContent(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	respondJSON(w, map[string]interface{}{
		"code": 200,
		"data": nil,
		"msg":  "success",
	})
}

func (lv *LogViewer) DeleteAllFilesHandler(w http.ResponseWriter, r *http.Request) {
	if !lv.config.DevMode || !lv.config.EnableDelete {
		respondJSON(w, map[string]interface{}{
			"code":  3001,
			"files": nil,
			"msg":   "delete operation is disabled",
		})
		return
	}
	err := lv.DeleteAllLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	respondJSON(w, map[string]interface{}{
		"code": 200,
		"data": nil,
		"msg":  "success",
	})
}

// 导出指定文件
func (lv *LogViewer) ExportFileHandler(w http.ResponseWriter, r *http.Request) {
	if !lv.config.EnableExport {
		respondJSON(w, map[string]interface{}{
			"code":  3001,
			"files": nil,
			"msg":   "export operation is disabled",
		})
		return
	}

	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		respondJSON(w, map[string]interface{}{
			"code":  3002,
			"files": nil,
			"msg":   "Missing 'name' parameter",
		})
		return
	}

	// 文件名安全校验
	if strings.Contains(fileName, "..") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		respondJSON(w, map[string]interface{}{
			"code":  3003,
			"files": nil,
			"msg":   "Invalid filename",
		})
		return
	}

	filePath := filepath.Join(lv.config.LogDir, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			respondJSON(w, map[string]interface{}{
				"code":  3004,
				"files": nil,
				"msg":   "File not found",
			})
		} else {
			respondJSON(w, map[string]interface{}{
				"code":  3005,
				"files": nil,
				"msg":   "Failed to read file",
			})
		}
		return
	}
	respondJSON(w, map[string]interface{}{
		"code": 200,
		"data": string(content),
		"msg":  "success",
	})
}
