/**
 * @Author: guxline zjguoxin@163.com
 * @Date: 2025/7/3 07:16:31
 * @LastEditors: guxline zjguoxin@163.com
 * @LastEditTime: 2025/7/3 07:16:31
 * Description: 核心功能实现
 * Copyright: Copyright (©) 2025 中易综服. All rights reserved.
 */
package goslogviewer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LogEntry struct {
	Level string `json:"level"`
	Time  string `json:"time"`
	Msg   string `json:"msg"`
}

// GetLogFiles 获取日志文件列表
func (lv *LogViewer) GetLogFiles() ([]string, error) {
	var files []string
	entries, err := os.ReadDir(lv.config.LogDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// GetLogContent 获取日志内容
func (lv *LogViewer) GetLogContent(filename string) ([]LogEntry, error) {
	path := filepath.Join(lv.config.LogDir, filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logs []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var log LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &log); err == nil {
			logs = append(logs, log)
		}
	}
	return logs, scanner.Err()
}

// DeleteAllLogs 删除所有日志文件
func (lv *LogViewer) DeleteAllLogs() error {
	if !lv.config.DevMode || !lv.config.EnableDelete {
		return fmt.Errorf("delete operation is disabled")
	}

	entries, err := os.ReadDir(lv.config.LogDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(lv.config.LogDir, entry.Name())
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

// ClearFileContent 清空文件内容
func (lv *LogViewer) ClearFileContent(filename string) error {
	if !lv.config.DevMode || !lv.config.EnableClear {
		return fmt.Errorf("clear operation is disabled")
	}
	path := filepath.Join(lv.config.LogDir, filename)
	return os.WriteFile(path, []byte{}, 0644)
}

func (lv *LogViewer) ExportFile(filename string) error {
	path := filepath.Join(lv.config.LogDir, filename)
	return os.Remove(path)
}
