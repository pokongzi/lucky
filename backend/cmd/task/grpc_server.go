package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// handleGRPCRequest 处理 gRPC HTTP/2 请求
func handleGRPCRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("[gRPC] 收到请求: Method=%s, Path=%s, Proto=%s, ContentType=%s",
		r.Method, r.URL.Path, r.Proto, r.Header.Get("Content-Type"))

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[gRPC] 读取请求失败: %v", err)
		sendGRPCError(w, 13, "读取请求失败")
		return
	}

	log.Printf("[gRPC] 请求体长度: %d bytes", len(body))

	// 根据路径分发请求
	switch r.URL.Path {
	case "/rpc.Task/Run":
		handleTaskRun(w, r, body)
	case "/rpc.Task/Check":
		handleTaskCheck(w, r)
	default:
		log.Printf("[gRPC] 未知路径: %s", r.URL.Path)
		sendGRPCError(w, 12, "未实现的方法")
	}
}

// handleTaskRun 处理任务执行请求
func handleTaskRun(w http.ResponseWriter, r *http.Request, body []byte) {
	// 尝试多种方式解析请求
	var command string
	var req map[string]interface{}

	// 方式1: 尝试 JSON 解析
	if len(body) > 5 {
		// 跳过 gRPC 的 5 字节前缀（1字节压缩标志 + 4字节消息长度）
		actualBody := body
		if body[0] == 0 && len(body) > 5 {
			actualBody = body[5:]
		}

		if err := json.Unmarshal(actualBody, &req); err == nil {
			if cmd, ok := req["command"].(string); ok {
				command = cmd
			} else if name, ok := req["name"].(string); ok {
				command = name
			}
		}
	}

	// 方式2: 从 URL 参数获取
	if command == "" {
		command = r.URL.Query().Get("command")
	}

	// 方式3: 直接使用原始字符串
	if command == "" && len(body) > 0 {
		command = strings.Trim(string(body), "\x00\n\r\t \"")
	}

	log.Printf("[gRPC] Run 接收到任务请求: command=%s, body_len=%d", command, len(body))

	// gRPC 服务仅用于 gocron 节点检测和任务调度
	// 不执行实际抓取逻辑，直接返回成功
	log.Printf("[gRPC] 任务已接收: %s（不执行实际抓取）", command)

	// gocron 期望的格式：output + error（error 为空表示成功）
	sendGRPCResponse(w, map[string]interface{}{
		"output": "task received successfully",
		"error":  "", // 空字符串表示成功
	})
}

// handleTaskCheck 处理健康检查请求
func handleTaskCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("[gRPC] Check 健康检查请求")
	sendGRPCResponse(w, map[string]interface{}{
		"status":  "ok",
		"message": "服务运行正常",
	})
}

// sendGRPCResponse 发送 gRPC 响应（Protobuf 格式）
func sendGRPCResponse(w http.ResponseWriter, data interface{}) {
	// 构造简单的 Protobuf 消息
	protobufMsg := encodeSimpleProtobuf(data)

	// gRPC 消息格式：1字节压缩标志 + 4字节消息长度 + 消息体
	msgLen := len(protobufMsg)
	grpcMsg := make([]byte, 5+msgLen)

	// 第1字节：压缩标志（0 = 不压缩）
	grpcMsg[0] = 0

	// 第2-5字节：消息长度（大端序）
	grpcMsg[1] = byte(msgLen >> 24)
	grpcMsg[2] = byte(msgLen >> 16)
	grpcMsg[3] = byte(msgLen >> 8)
	grpcMsg[4] = byte(msgLen)

	// 第6字节开始：Protobuf 消息体
	copy(grpcMsg[5:], protobufMsg)

	// 必须先声明要发送的 trailers
	w.Header().Set("Trailer", "Grpc-Status, Grpc-Message")

	// 设置响应头
	w.Header().Set("Content-Type", "application/grpc")

	// 发送响应头和消息体
	w.WriteHeader(http.StatusOK)
	n, _ := w.Write(grpcMsg)

	// 强制 flush，确保消息体发送
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// 在消息体发送后设置 trailers（注意大小写：Grpc-Status）
	w.Header().Set("Grpc-Status", "0")
	w.Header().Set("Grpc-Message", "")

	log.Printf("[gRPC] Protobuf 响应已发送: %d bytes (消息体: %d bytes), trailers 已设置", n, msgLen)
}

// encodeSimpleProtobuf 编码简单的 Protobuf 消息
// gocron 期望的格式：Result { output, error }
func encodeSimpleProtobuf(data interface{}) []byte {
	// 将数据转为 map
	var dataMap map[string]interface{}
	jsonBytes, _ := json.Marshal(data)
	json.Unmarshal(jsonBytes, &dataMap)

	buf := []byte{}

	// Field 1: output (string)
	if output, ok := dataMap["output"].(string); ok {
		buf = append(buf, 0x0a) // tag: field=1, wire_type=2 (length-delimited)
		buf = appendString(buf, output)
	}

	// Field 2: error (string) - 空表示成功
	if errMsg, ok := dataMap["error"].(string); ok {
		buf = append(buf, 0x12) // tag: field=2, wire_type=2 (length-delimited)
		buf = appendString(buf, errMsg)
	}

	// Field 1: status (string) - for Check response
	if status, ok := dataMap["status"].(string); ok {
		buf = append(buf, 0x0a) // tag: field=1, wire_type=2 (length-delimited)
		buf = appendString(buf, status)
	}

	return buf
}

// appendVarint 添加 varint 编码的整数
func appendVarint(buf []byte, value int64) []byte {
	for value >= 0x80 {
		buf = append(buf, byte(value)|0x80)
		value >>= 7
	}
	buf = append(buf, byte(value))
	return buf
}

// appendString 添加 length-delimited 字符串
func appendString(buf []byte, s string) []byte {
	// 添加长度
	buf = appendVarint(buf, int64(len(s)))
	// 添加字符串内容
	buf = append(buf, []byte(s)...)
	return buf
}

// sendGRPCError 发送 gRPC 错误
func sendGRPCError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/grpc+json")
	w.Header().Set("grpc-status", fmt.Sprintf("%d", code))
	w.Header().Set("grpc-message", message)
	w.WriteHeader(http.StatusOK)
}
