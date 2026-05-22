package cluster

import (
	"hash/fnv"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"cloud-smart-tutor-backend/internal/domain"
	"cloud-smart-tutor-backend/internal/websocket"
)

type ErrorClusterKey struct {
	Hash        uint32
	ClassroomID string
}

type ErrorCluster struct {
	Key         ErrorClusterKey
	Signature   string
	Count       int
	Students    map[string]bool
	FirstSeen   time.Time
	LastSeen    time.Time
	Severity    string
	Suggestion  *domain.AITeachingSuggestion
}

type ErrorClusteringEngine struct {
	clusters     map[ErrorClusterKey]*ErrorCluster
	threshold    float64
	hub          *websocket.Hub
	mu           sync.RWMutex
	alertCooldown map[ErrorClusterKey]time.Time
	alertMu      sync.Mutex
}

func NewErrorClusteringEngine(hub *websocket.Hub, threshold float64) *ErrorClusteringEngine {
	return &ErrorClusteringEngine{
		clusters:      make(map[ErrorClusterKey]*ErrorCluster),
		threshold:     threshold,
		hub:           hub,
		alertCooldown: make(map[ErrorClusterKey]time.Time),
	}
}

func (e *ErrorClusteringEngine) SignatureHash(errorMessage string) uint32 {
	normalized := strings.ToLower(strings.TrimSpace(errorMessage))
	normalized = strings.ReplaceAll(normalized, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\r", " ")
	normalized = strings.Join(strings.Fields(normalized), " ")

	h := fnv.New32a()
	h.Write([]byte(normalized))
	return h.Sum32()
}

func (e *ErrorClusteringEngine) SignatureCategory(hash uint32, message string) string {
	msg := strings.ToLower(message)
	if strings.Contains(msg, "nginx") && strings.Contains(msg, "syntax") {
		return "Nginx 配置语法错误"
	}
	if strings.Contains(msg, "port") && strings.Contains(msg, "already in use") {
		return "端口占用冲突"
	}
	if strings.Contains(msg, "permission denied") {
		return "权限不足"
	}
	if strings.Contains(msg, "not found") || strings.Contains(msg, "no such file") {
		return "文件或模块未找到"
	}
	if strings.Contains(msg, "connection refused") {
		return "连接被拒绝"
	}
	if strings.Contains(msg, "timeout") {
		return "连接超时"
	}
	return "其他错误"
}

func (e *ErrorClusteringEngine) IngestError(event domain.ProbeEvent, totalStudents int) {
	key := ErrorClusterKey{
		Hash:        e.SignatureHash(event.ErrorMessage),
		ClassroomID: event.ClassroomID,
	}

	e.mu.Lock()
	cluster, exists := e.clusters[key]
	if !exists {
		category := e.SignatureCategory(key.Hash, event.ErrorMessage)
		cluster = &ErrorCluster{
			Key:       key,
			Signature: category,
			Count:     0,
			Students:  make(map[string]bool),
			FirstSeen: time.Now(),
			LastSeen:  time.Now(),
		}
		e.clusters[key] = cluster
	}

	if !cluster.Students[event.StudentID] {
		cluster.Students[event.StudentID] = true
		cluster.Count++
	}
	cluster.LastSeen = time.Now()
	e.mu.Unlock()

	percentage := float64(cluster.Count) / float64(totalStudents) * 100

	if percentage >= e.threshold {
		e.maybeSendAlert(key, cluster, percentage, event.ClassroomID)
	}

	if percentage >= 5 {
		cluster.Severity = "high"
	} else if percentage >= 2 {
		cluster.Severity = "medium"
	} else {
		cluster.Severity = "low"
	}

	e.hub.SendToRoom("teacher:"+event.ClassroomID, &websocket.WSMessage{
		Type:      websocket.WSErrorCluster,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"signature":  cluster.Signature,
			"count":      cluster.Count,
			"percentage": percentage,
			"severity":   cluster.Severity,
			"students":   len(cluster.Students),
			"last_seen":  cluster.LastSeen,
		},
		Room: "teacher:" + event.ClassroomID,
	})
}

func (e *ErrorClusteringEngine) maybeSendAlert(key ErrorClusterKey, cluster *ErrorCluster, percentage float64, classroomID string) {
	e.alertMu.Lock()
	defer e.alertMu.Unlock()

	if lastAlert, exists := e.alertCooldown[key]; exists {
		if time.Since(lastAlert) < 5*time.Minute {
			return
		}
	}

	e.alertCooldown[key] = time.Now()

	log.Printf("TEACHING ALERT: %.1f%% students hit error '%s' in classroom %s", percentage, cluster.Signature, classroomID)

	suggestion := e.generateTeachingSuggestion(key, cluster, classroomID)

	e.hub.SendToRoom("teacher:"+classroomID, &websocket.WSMessage{
		Type:      websocket.WSTeachingAlert,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"type":        "TEACHING_ALERT",
			"signature":   cluster.Signature,
			"count":       cluster.Count,
			"percentage":  percentage,
			"severity":    "critical",
			"students":    len(cluster.Students),
			"suggestion":  suggestion,
			"action":      "consider_group_explanation",
		},
		Room: "teacher:" + classroomID,
	})
}

func (e *ErrorClusteringEngine) generateTeachingSuggestion(key ErrorClusterKey, cluster *ErrorCluster, classroomID string) *domain.AITeachingSuggestion {
	var actionType, title, content string

	switch cluster.Signature {
	case "Nginx 配置语法错误":
		actionType = "demo"
		title = "建议集中讲解 Nginx 配置文件语法"
		content = "检测到多名学生在 Nginx 配置文件中出现同类语法错误（缺少分号或括号不匹配）。建议在接下来 10 分钟内切回大屏，集中演示 Nginx 配置文件的正确写法。"
	case "端口占用冲突":
		actionType = "handout"
		title = "建议增加端口排查指南"
		content = "多名学生遇到端口占用问题。建议在实验手册中增加端口排查步骤，或演示如何使用 lsof/netstat 定位占用端口的进程。"
	case "权限不足":
		actionType = "explain"
		title = "建议讲解 Linux 权限模型"
		content = "多名学生遇到权限问题。建议集中讲解 sudo 使用场景和 Linux 文件权限体系。"
	default:
		actionType = "explain"
		title = "检测到高频错误聚类"
		content = cluster.Signature
	}

	return &domain.AITeachingSuggestion{
		SuggestionID:  cluster.Signature,
		ErrorClusterID: cluster.Signature,
		ClassroomID:   classroomID,
		Timestamp:     time.Now(),
		Title:         title,
		Content:       content,
		Confidence:    0.85,
		ActionType:    actionType,
		Priority:      "high",
	}
}

func (e *ErrorClusteringEngine) GetTopErrors(classroomID string, limit int) []map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var errors []map[string]interface{}
	for _, cluster := range e.clusters {
		if cluster.ClassroomID == classroomID {
			errors = append(errors, map[string]interface{}{
				"signature":  cluster.Signature,
				"count":      cluster.Count,
				"students":   len(cluster.Students),
				"severity":   cluster.Severity,
				"first_seen": cluster.FirstSeen,
				"last_seen":  cluster.LastSeen,
			})
		}
	}

	sort.Slice(errors, func(i, j int) bool {
		return errors[i]["count"].(int) > errors[j]["count"].(int)
	})

	if limit > 0 && len(errors) > limit {
		errors = errors[:limit]
	}

	return errors
}

func (e *ErrorClusteringEngine) CleanupExpired(maxAge time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	for key, cluster := range e.clusters {
		if now.Sub(cluster.LastSeen) > maxAge {
			delete(e.clusters, key)
		}
	}
}
