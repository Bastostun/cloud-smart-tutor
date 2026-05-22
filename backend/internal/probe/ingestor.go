package probe

import (
	"fmt"
	"log"
	"time"

	"cloud-smart-tutor-backend/internal/cluster"
	"cloud-smart-tutor-backend/internal/domain"
	"cloud-smart-tutor-backend/internal/sandbox"
	"cloud-smart-tutor-backend/internal/websocket"
)

type ProbeIngestor struct {
	Hub              *websocket.Hub
	ErrorEngine      *cluster.ErrorClusteringEngine
	AuditStream      *sandbox.AuditStream
	TotalStudents    int
	AIResponseEngine *AIResponseEngine
}

func NewProbeIngestor(
	hub *websocket.Hub,
	errorEngine *cluster.ErrorClusteringEngine,
	auditStream *sandbox.AuditStream,
	totalStudents int,
) *ProbeIngestor {
	return &ProbeIngestor{
		Hub:              hub,
		ErrorEngine:      errorEngine,
		AuditStream:      auditStream,
		TotalStudents:    totalStudents,
		AIResponseEngine: NewAIResponseEngine(hub),
	}
}

func (pi *ProbeIngestor) Ingest(event domain.ProbeEvent) {
	switch event.EventType {
	case domain.ProbeStderrCapture, domain.ProbeExitCode:
		pi.handleStudentError(event)

	case domain.ProbeClipboardCopy, domain.ProbeDLPViolation, domain.ProbeNetworkAttempt:
		pi.handleSandboxEvent(event)

	case domain.ProbeFileAccess:
		pi.handleFileAccess(event)

	case domain.ProbeAnomalyDetected:
		pi.handleAnomaly(event)

	case domain.ProbeHeartbeat:
		// Heartbeat handled separately by node monitor
	}
}

func (pi *ProbeIngestor) handleStudentError(event domain.ProbeEvent) {
	pi.ErrorEngine.IngestError(event, pi.TotalStudents)

	go pi.AIResponseEngine.GenerateResponse(event)

	pi.Hub.SendToRoom("student:"+event.StudentID, &websocket.WSMessage{
		Type:      websocket.WSProbeEvent,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"event_type":    event.EventType,
			"error_message": event.ErrorMessage,
			"command":       event.Command,
			"exit_code":     event.ExitCode,
			"student_id":    event.StudentID,
		},
		Room: "student:" + event.StudentID,
	})
}

func (pi *ProbeIngestor) handleSandboxEvent(event domain.ProbeEvent) {
	event.RiskLevel = "high"

	pi.AuditStream.ProcessEvent(event)

	log.Printf("SANDBOX EVENT: VM=%s, type=%s, risk=%s", event.VMID, event.EventType, event.RiskLevel)
}

func (pi *ProbeIngestor) handleFileAccess(event domain.ProbeEvent) {
	if event.RiskLevel == "critical" {
		pi.AuditStream.ProcessEvent(event)
	}
}

func (pi *ProbeIngestor) handleAnomaly(event domain.ProbeEvent) {
	event.RiskLevel = "critical"
	pi.AuditStream.ProcessEvent(event)

	pi.Hub.SendToRoom("admin", &websocket.WSMessage{
		Type:      websocket.WSProbeEvent,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"event_type":   event.EventType,
			"vm_id":        event.VMID,
			"error_message": event.ErrorMessage,
			"risk_level":   "critical",
		},
	})
}

func (pi *ProbeIngestor) SetTotalStudents(count int) {
	pi.TotalStudents = count
}

type AIResponseEngine struct {
	Hub *websocket.Hub
}

func NewAIResponseEngine(hub *websocket.Hub) *AIResponseEngine {
	return &AIResponseEngine{Hub: hub}
}

func (ae *AIResponseEngine) GenerateResponse(event domain.ProbeEvent) {
	time.Sleep(500 * time.Millisecond)

	analysis := ae.analyzeError(event.ErrorMessage)

	response := domain.AICompanionResponse{
		ResponseID:   generateResponseID(),
		StudentID:    event.StudentID,
		VMID:         event.VMID,
		Timestamp:    time.Now(),
		OriginalError: event.ErrorMessage,
		Explanation:  analysis.Explanation,
		Solution:     analysis.Solution,
		Confidence:   analysis.Confidence,
		Proactive:    true,
	}

	ae.Hub.SendToRoom("student:"+event.StudentID, &websocket.WSMessage{
		Type:      websocket.WSAICompanion,
		Timestamp: time.Now(),
		Payload:   response,
		Room:      "student:" + event.StudentID,
	})

	log.Printf("AI Response generated for student %s: %s", event.StudentID, analysis.Solution[:50])
}

type AIAnalysis struct {
	Explanation string
	Solution    string
	Confidence  float64
}

func (ae *AIResponseEngine) analyzeError(errorMessage string) AIAnalysis {
	msg := errorMessage

	if contains(msg, "nginx") && contains(msg, "syntax") {
		return AIAnalysis{
			Explanation: "Nginx 配置文件存在语法错误。常见原因是缺少分号 (;) 或括号不匹配。",
			Solution: `检查配置文件语法：
1. sudo nginx -t 测试配置
2. 检查 /etc/nginx/nginx.conf 和 sites-enabled/ 下的文件
3. 确保每个指令以分号结尾
4. 确保 server { } 块正确闭合`,
			Confidence: 0.95,
		}
	}

	if contains(msg, "port") && (contains(msg, "already in use") || contains(msg, "bind")) {
		return AIAnalysis{
			Explanation: "目标端口已被其他进程占用。需要先找出占用进程并处理。",
			Solution: `排查端口占用：
1. sudo lsof -i :80 查看占用 80 端口的进程
2. sudo netstat -tlnp | grep :80 确认
3. sudo systemctl stop apache2 (如果是 apache 占用)
4. 或修改 Nginx 配置使用其他端口`,
			Confidence: 0.92,
		}
	}

	if contains(msg, "permission denied") {
		return AIAnalysis{
			Explanation: "当前用户没有执行该操作的权限。需要使用 sudo 或检查文件权限。",
			Solution: `解决权限问题：
1. 使用 sudo 执行特权操作：sudo systemctl restart nginx
2. 检查文件权限：ls -la /path/to/file
3. 修改文件所有者：sudo chown user:group /path
4. 添加执行权限：chmod +x script.sh`,
			Confidence: 0.88,
		}
	}

	if contains(msg, "not found") || contains(msg, "no such file") {
		return AIAnalysis{
			Explanation: "指定的文件、目录或软件包不存在。",
			Solution: `排查文件未找到问题：
1. 确认路径是否正确：ls -la /path/to/check
2. 检查软件包是否安装：dpkg -l | grep nginx
3. 安装缺失的软件包：sudo apt install nginx
4. 查找文件位置：find / -name "filename" 2>/dev/null`,
			Confidence: 0.85,
		}
	}

	if contains(msg, "connection refused") {
		return AIAnalysis{
			Explanation: "目标服务没有在该端口监听，或连接被防火墙拦截。",
			Solution: `解决连接被拒绝：
1. 确认服务已启动：sudo systemctl status nginx
2. 检查端口监听：sudo netstat -tlnp
3. 检查防火墙：sudo ufw status
4. 允许端口：sudo ufw allow 80/tcp`,
			Confidence: 0.82,
		}
	}

	if contains(msg, "timeout") {
		return AIAnalysis{
			Explanation: "操作超时，可能是网络问题或服务响应过慢。",
			Solution: `排查超时问题：
1. 检查网络连接：ping target-server
2. 检查 DNS 解析：nslookup target-server
3. 增加超时时间配置
4. 检查目标服务是否正常运行`,
			Confidence: 0.75,
		}
	}

	return AIAnalysis{
		Explanation: "检测到系统错误，正在分析中...",
		Solution: `通用排查步骤：
1. 查看错误详情：journalctl -xe
2. 检查服务状态：systemctl status <service>
3. 查看系统日志：tail -f /var/log/syslog
4. 尝试重启服务：sudo systemctl restart <service>`,
		Confidence: 0.60,
	}
}

func contains(s, substr string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

var responseCounter int64

func generateResponseID() string {
	responseCounter++
	return time.Now().Format("20060102150405") + "-ai-" + fmt.Sprintf("%04d", responseCounter%10000)
}
