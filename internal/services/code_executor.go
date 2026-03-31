package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gin/internal/models"
	"gin/internal/utils"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// CodeExecutor 代码执行器接口
type CodeExecutor interface {
	Execute(ctx context.Context, language, code, stdin string) (*models.ExecuteCodeResponse, error)
	GetSupportedLanguages() []models.LanguageInfo
}

// Judge0CodeExecutor Judge0 API 代码执行器实现
type Judge0CodeExecutor struct {
	apiURL      string
	timeout     time.Duration
	maxMemoryKB int
	client      *http.Client
}

// 支持的语言配置
var supportedLanguages = map[string]models.LanguageInfo{
	"python": {
		ID:         "python",
		Name:       "Python",
		Version:    "3.13.2",
		PistonName: "python",
		Judge0ID:   109,
		DefaultCode: `# Python 示例代码
print("Hello, World!")

# 获取用户输入
# name = input("请输入你的名字: ")
# print(f"你好, {name}!")`,
	},
	"javascript": {
		ID:         "javascript",
		Name:       "JavaScript (Node.js)",
		Version:    "22.08.0",
		PistonName: "javascript",
		Judge0ID:   102,
		DefaultCode: `// JavaScript 示例代码
console.log("Hello, World!");

// 获取用户输入
// const readline = require('readline');
// const rl = readline.createInterface({
//   input: process.stdin,
//   output: process.stdout
// });`,
	},
	"java": {
		ID:         "java",
		Name:       "Java",
		Version:    "17.0.6",
		PistonName: "java",
		Judge0ID:   91,
		DefaultCode: `// Java 示例代码
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
        
        // 获取用户输入
        // Scanner scanner = new Scanner(System.in);
        // String name = scanner.nextLine();
        // System.out.println("你好, " + name + "!");
    }
}`,
	},
	"cpp": {
		ID:         "cpp",
		Name:       "C++",
		Version:    "14.1.0",
		PistonName: "cpp",
		Judge0ID:   105,
		DefaultCode: `// C++ 示例代码
#include <iostream>
using namespace std;

int main() {
    cout << "Hello, World!" << endl;
    
    // 获取用户输入
    // string name;
    // cout << "请输入你的名字: ";
    // cin >> name;
    // cout << "你好, " << name << "!" << endl;
    
    return 0;
}`,
	},
	"c": {
		ID:         "c",
		Name:       "C",
		Version:    "14.1.0",
		PistonName: "c",
		Judge0ID:   103,
		DefaultCode: `// C 示例代码
#include <stdio.h>

int main() {
    printf("Hello, World!\n");
    
    // 获取用户输入
    // char name[100];
    // printf("请输入你的名字: ");
    // scanf("%s", name);
    // printf("你好, %s!\n", name);
    
    return 0;
}`,
	},
	"go": {
		ID:         "go",
		Name:       "Go",
		Version:    "1.23.5",
		PistonName: "go",
		Judge0ID:   107,
		DefaultCode: `// Go 示例代码
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    
    // 获取用户输入
    // var name string
    // fmt.Print("请输入你的名字: ")
    // fmt.Scanln(&name)
    // fmt.Printf("你好, %s!\n", name)
}`,
	},
	"rust": {
		ID:         "rust",
		Name:       "Rust",
		Version:    "1.85.0",
		PistonName: "rust",
		Judge0ID:   108,
		DefaultCode: `// Rust 示例代码
fn main() {
    println!("Hello, World!");
    
    // 获取用户输入
    // use std::io;
    // let mut name = String::new();
    // io::stdin().read_line(&mut name).expect("Failed to read line");
    // println!("你好, {}!", name.trim());
}`,
	},
	"php": {
		ID:         "php",
		Name:       "PHP",
		Version:    "8.3.11",
		PistonName: "php",
		Judge0ID:   98,
		DefaultCode: `<?php
// PHP 示例代码
echo "Hello, World!\n";
?>`,
	},
	"ruby": {
		ID:         "ruby",
		Name:       "Ruby",
		Version:    "2.7.0",
		PistonName: "ruby",
		Judge0ID:   72,
		DefaultCode: `# Ruby 示例代码
puts "Hello, World!"`,
	},
	"swift": {
		ID:         "swift",
		Name:       "Swift",
		Version:    "5.2.3",
		PistonName: "swift",
		Judge0ID:   83,
		DefaultCode: `// Swift 示例代码
print("Hello, World!")`,
	},
	"bash": {
		ID:         "bash",
		Name:       "Bash",
		Version:    "5.0.0",
		PistonName: "bash",
		Judge0ID:   46,
		DefaultCode: `#!/bin/bash
# Bash 示例代码
echo "Hello, World!"`,
	},
	"lua": {
		ID:         "lua",
		Name:       "Lua",
		Version:    "5.3.5",
		PistonName: "lua",
		Judge0ID:   64,
		DefaultCode: `-- Lua 示例代码
print("Hello, World!")`,
	},
	"scala": {
		ID:         "scala",
		Name:       "Scala",
		Version:    "3.4.2",
		PistonName: "scala",
		Judge0ID:   112,
		DefaultCode: `// Scala 示例代码
object Main extends App {
  println("Hello, World!")
}`,
	},
	"haskell": {
		ID:         "haskell",
		Name:       "Haskell",
		Version:    "8.8.1",
		PistonName: "haskell",
		Judge0ID:   61,
		DefaultCode: `-- Haskell 示例代码
main :: IO ()
main = putStrLn "Hello, World!"`,
	},
	"perl": {
		ID:         "perl",
		Name:       "Perl",
		Version:    "5.28.1",
		PistonName: "perl",
		Judge0ID:   85,
		DefaultCode: `# Perl 示例代码
print "Hello, World!\n";`,
	},
}

// NewJudge0CodeExecutor 创建新的 Judge0 代码执行器
func NewJudge0CodeExecutor(apiURL string, timeout time.Duration, maxMemoryMB int, maxIdleConns int, maxIdleConnsPerHost int, idleConnTimeout int) CodeExecutor {
	if maxMemoryMB <= 0 {
		maxMemoryMB = 128
	}

	return &Judge0CodeExecutor{
		apiURL:      strings.TrimRight(apiURL, "/"),
		timeout:     timeout,
		maxMemoryKB: maxMemoryMB * 1024,
		client: &http.Client{
			Timeout: timeout + 2*time.Second,
			Transport: &http.Transport{
				// 当前环境访问 ce.judge0.com 的 IPv6 连接会被远端重置，固定走 IPv4。
				ForceAttemptHTTP2:   false,
				MaxIdleConns:        maxIdleConns,
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
				IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
				DisableCompression:  false,
				DisableKeepAlives:   false,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					dialer := &net.Dialer{Timeout: timeout}
					return dialer.DialContext(ctx, "tcp4", addr)
				},
			},
		},
	}
}

// Execute 执行代码
func (e *Judge0CodeExecutor) Execute(ctx context.Context, language, code, stdin string) (*models.ExecuteCodeResponse, error) {
	logger := utils.GetLogger()

	langInfo, ok := supportedLanguages[language]
	if !ok {
		return nil, fmt.Errorf("不支持的语言: %s", language)
	}

	reqPayload := models.Judge0SubmissionRequest{
		SourceCode:    base64.StdEncoding.EncodeToString([]byte(code)),
		LanguageID:    langInfo.Judge0ID,
		Stdin:         base64.StdEncoding.EncodeToString([]byte(stdin)),
		CPUTimeLimit:  e.timeout.Seconds(),
		WallTimeLimit: e.timeout.Seconds() + 1,
		MemoryLimit:   e.maxMemoryKB,
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	startTime := time.Now()

	buf := utils.GetBuffer()
	defer utils.PutBuffer(buf)
	buf.Write(reqBody)

	requestURL := e.apiURL + "/submissions?base64_encoded=true&wait=true"
	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		logger.Error("Judge0 API 请求失败", "error", err)
		return &models.ExecuteCodeResponse{
			Output:        "",
			Error:         "代码执行超时或服务不可用",
			ExecutionTime: int(time.Since(startTime).Milliseconds()),
			Status:        "timeout",
		}, nil
	}
	defer resp.Body.Close()

	executionTime := int(time.Since(startTime).Milliseconds())

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var judge0Err models.Judge0ErrorResponse
		if err := json.Unmarshal(respBody, &judge0Err); err == nil {
			message := firstNonEmpty(judge0Err.Message, judge0Err.Error)
			if message != "" {
				logger.Warn("Judge0 API 返回错误",
					"status_code", resp.StatusCode,
					"message", message)
				return &models.ExecuteCodeResponse{
					Output:        "",
					Error:         message,
					ExecutionTime: executionTime,
					MemoryUsage:   0,
					Status:        "error",
				}, nil
			}
		}

		logger.Warn("Judge0 API 返回非成功状态",
			"status_code", resp.StatusCode,
			"response_body", utils.TruncateString(string(respBody), 512))
		return &models.ExecuteCodeResponse{
			Output:        "",
			Error:         fmt.Sprintf("代码执行服务暂时不可用（HTTP %d）", resp.StatusCode),
			ExecutionTime: executionTime,
			MemoryUsage:   0,
			Status:        "error",
		}, nil
	}

	var judge0Resp models.Judge0SubmissionResponse
	if err := json.Unmarshal(respBody, &judge0Resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if apiError := nullableString(judge0Resp.Error); apiError != "" {
		logger.Warn("Judge0 提交响应返回错误", "error", apiError)
		return &models.ExecuteCodeResponse{
			Output:        "",
			Error:         apiError,
			ExecutionTime: executionTime,
			MemoryUsage:   0,
			Status:        "error",
		}, nil
	}

	stdout := decodeMaybeBase64(nullableString(judge0Resp.Stdout))
	stderr := decodeMaybeBase64(nullableString(judge0Resp.Stderr))
	compileOutput := decodeMaybeBase64(nullableString(judge0Resp.CompileOutput))
	message := decodeMaybeBase64(nullableString(judge0Resp.Message))

	result := &models.ExecuteCodeResponse{
		Output:        stdout,
		ExecutionTime: executionTime,
		MemoryUsage:   judge0Resp.Memory * 1024,
	}

	switch judge0Resp.Status.ID {
	case 3:
		result.Status = "success"
	case 5:
		result.Status = "timeout"
		result.Error = firstNonEmpty(message, stderr, compileOutput, judge0Resp.Status.Description, "代码执行超时")
	case 1, 2:
		result.Status = "error"
		result.Error = "代码仍在处理中，请稍后重试"
	default:
		result.Status = "error"
		result.Error = firstNonEmpty(compileOutput, stderr, message, judge0Resp.Status.Description, "程序执行失败，但未返回错误详情")
	}

	logger.Info("代码执行完成",
		"language", language,
		"judge0_language_id", langInfo.Judge0ID,
		"judge0_status_id", judge0Resp.Status.ID,
		"judge0_status", judge0Resp.Status.Description,
		"status", result.Status,
		"execution_time", executionTime,
		"code_length", len(code))

	return result, nil
}

// GetSupportedLanguages 获取支持的语言列表
func (e *Judge0CodeExecutor) GetSupportedLanguages() []models.LanguageInfo {
	languages := make([]models.LanguageInfo, 0, len(supportedLanguages))
	for _, lang := range supportedLanguages {
		languages = append(languages, lang)
	}
	return languages
}

func nullableString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func decodeMaybeBase64(value string) string {
	if value == "" {
		return ""
	}

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return value
	}

	return string(decoded)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
