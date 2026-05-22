package sandbox

import (
	"fmt"
	"log"
	"sync"
	"time"

	"cloud-smart-tutor-backend/internal/domain"
	"cloud-smart-tutor-backend/internal/websocket"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type InterventionAction string

const (
	InterventionNone    InterventionAction = "none"
	InterventionWarn    InterventionAction = "warn"
	InterventionBlock   InterventionAction = "block"
	InterventionSuspend InterventionAction = "suspend_vm"
)

type SecurityRule struct {
	RuleID          string
	Name            string
	Description     string
	TriggerPatterns []string
	RiskLevel       RiskLevel
	Intervention    InterventionAction
	Enabled         bool
}

type AuditStream struct {
	Rules         map[string]*SecurityRule
	VMStates      map[string]*VMSecurityState
	PVEController PVEController
	Hub           *websocket.Hub
	LogBuffer     chan domain.AuditLogEntry
	mu            sync.RWMutex
}

type VMSecurityState struct {
	VMID            string
	EmployeeID      string
	ClassroomID     string
	RiskScore       float64
	BlockCount      int
	ScanCount       int
	ClipboardCopies int
	NetworkAttempts int
	LastActivity    time.Time
	Suspended       bool
}

type PVEController interface {
	SuspendVM(node string, vmID int) error
}

func NewAuditStream(hub *websocket.Hub, pveCtrl PVEController) *AuditStream {
	as := &AuditStream{
		Rules:         make(map[string]*SecurityRule),
		VMStates:      make(map[string]*VMSecurityState),
		PVEController: pveCtrl,
		Hub:           hub,
		LogBuffer:     make(chan domain.AuditLogEntry, 500),
	}

	as.loadDefaultRules()

	go as.logProcessor()

	return as
}

func (as *AuditStream) loadDefaultRules() {
	rules := []*SecurityRule{
		{
			RuleID:          "rule-clipboard",
			Name:            "剪贴板拦截",
			Description:     "拦截沙箱内剪贴板复制操作",
			TriggerPatterns: []string{"clipboard_copy", "clipboard_paste", "ctrl+c", "ctrl+v"},
			RiskLevel:       RiskHigh,
			Intervention:    InterventionBlock,
			Enabled:         true,
		},
		{
			RuleID:          "rule-network-egress",
			Name:            "网络外发阻断",
			Description:     "阻断从沙箱向外网的任何数据传输",
			TriggerPatterns: []string{"network_attempt", "http_post", "curl", "wget", "scp"},
			RiskLevel:       RiskCritical,
			Intervention:    InterventionBlock,
			Enabled:         true,
		},
		{
			RuleID:          "rule-keyboard-audit",
			Name:            "键盘行为审计",
			Description:     "记录所有键盘输入行为用于审计",
			TriggerPatterns: []string{"keystroke", "keyboard"},
			RiskLevel:       RiskMedium,
			Intervention:    InterventionNone,
			Enabled:         true,
		},
		{
			RuleID:          "rule-dlp",
			Name:            "防泄漏保护",
			Description:     "检测并阻止敏感数据外发",
			TriggerPatterns: []string{"dlp_violation", "data_exfiltration", "sensitive_data"},
			RiskLevel:       RiskCritical,
			Intervention:    InterventionSuspend,
			Enabled:         true,
		},
		{
			RuleID:          "rule-screenshot",
			Name:            "截屏防护",
			Description:     "检测并模糊化截屏行为",
			TriggerPatterns: []string{"screenshot", "print_screen"},
			RiskLevel:       RiskHigh,
			Intervention:    InterventionWarn,
			Enabled:         true,
		},
		{
			RuleID:          "rule-usb-storage",
			Name:            "USB 存储禁用",
			Description:     "禁用 USB 存储设备访问",
			TriggerPatterns: []string{"usb_mount", "usb_storage"},
			RiskLevel:       RiskHigh,
			Intervention:    InterventionBlock,
			Enabled:         true,
		},
	}

	for _, rule := range rules {
		as.Rules[rule.RuleID] = rule
	}
}

func (as *AuditStream) ProcessEvent(event domain.ProbeEvent) {
	as.mu.Lock()
	vmState, exists := as.VMStates[event.VMID]
	if !exists {
		vmState = &VMSecurityState{
			VMID:         event.VMID,
			EmployeeID:   event.StudentID,
			ClassroomID:  event.ClassroomID,
			LastActivity: time.Now(),
		}
		as.VMStates[event.VMID] = vmState
	}
	vmState.LastActivity = time.Now()
	as.mu.Unlock()

	matchedRules := as.matchRules(event)

	for _, rule := range matchedRules {
		if !rule.Enabled {
			continue
		}

		logEntry := as.createAuditLog(event, rule)

		switch rule.Intervention {
		case InterventionNone:
			logEntry.Result = "allowed"
			vmState.ScanCount++

		case InterventionWarn:
			logEntry.Result = "warned"
			logEntry.Intervention = "warn"

		case InterventionBlock:
			logEntry.Result = "blocked"
			logEntry.Intervention = "block"
			vmState.BlockCount++

		case InterventionSuspend:
			logEntry.Result = "blocked"
			logEntry.Intervention = "suspend_vm"
			vmState.BlockCount++
			vmState.Suspended = true

			go func() {
				log.Printf("SUSPENDING VM: %s due to critical security violation", event.VMID)
				as.Hub.SendToRoom("enterprise:"+event.ClassroomID, &websocket.WSMessage{
					Type:      websocket.WSAuditLog,
					Timestamp: time.Now(),
					Payload: map[string]interface{}{
						"type":       "vm_suspended",
						"vm_id":      event.VMID,
						"reason":     rule.Name,
						"risk_level": rule.RiskLevel,
					},
					Room: "enterprise:" + event.ClassroomID,
				})
			}()
		}

		vmState.RiskScore = as.calculateRiskScore(vmState)
		logEntry.RiskScore = vmState.RiskScore

		select {
		case as.LogBuffer <- logEntry:
		default:
			log.Printf("Audit log buffer full, dropping entry")
		}

		as.Hub.SendToRoom("enterprise:"+event.ClassroomID, &websocket.WSMessage{
			Type:      websocket.WSAuditLog,
			Timestamp: time.Now(),
			Payload:   logEntry,
			Room:      "enterprise:" + event.ClassroomID,
		})
	}
}

func (as *AuditStream) matchRules(event domain.ProbeEvent) []*SecurityRule {
	var matched []*SecurityRule

	for _, rule := range as.Rules {
		if !rule.Enabled {
			continue
		}

		eventTypeMatch := false
		for _, pattern := range rule.TriggerPatterns {
			if string(event.EventType) == pattern {
				eventTypeMatch = true
				break
			}
		}

		messageMatch := false
		if !eventTypeMatch {
			for _, pattern := range rule.TriggerPatterns {
				if containsIgnoreCase(event.ErrorMessage, pattern) {
					messageMatch = true
					break
				}
			}
		}

		if eventTypeMatch || messageMatch {
			matched = append(matched, rule)
		}
	}

	return matched
}

func (as *AuditStream) createAuditLog(event domain.ProbeEvent, rule *SecurityRule) domain.AuditLogEntry {
	return domain.AuditLogEntry{
		LogID:      generateLogID(),
		VMID:       event.VMID,
		EmployeeID: event.StudentID,
		ClassroomID: event.ClassroomID,
		Timestamp:  time.Now(),
		Action:     rule.Name,
		Result:     "pending",
		Details:    event.ErrorMessage,
		RiskScore:  0,
	}
}

func (as *AuditStream) calculateRiskScore(state *VMSecurityState) float64 {
	score := 0.0

	score += float64(state.BlockCount) * 15.0
	score += float64(state.ClipboardCopies) * 5.0
	score += float64(state.NetworkAttempts) * 10.0

	if time.Since(state.LastActivity) > 5*time.Minute {
		score += 10.0
	}

	if score > 100 {
		score = 100
	}

	return score
}

func (as *AuditStream) logProcessor() {
	for logEntry := range as.LogBuffer {
		log.Printf("[AUDIT] VM=%s, Action=%s, Result=%s, Risk=%.1f",
			logEntry.VMID, logEntry.Action, logEntry.Result, logEntry.RiskScore)
	}
}

func (as *AuditStream) GetVMState(vmID string) *VMSecurityState {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.VMStates[vmID]
}

func (as *AuditStream) GetPolicies() []map[string]interface{} {
	as.mu.RLock()
	defer as.mu.RUnlock()

	var policies []map[string]interface{}
	for _, rule := range as.Rules {
		policies = append(policies, map[string]interface{}{
			"name":        rule.Name,
			"enabled":     rule.Enabled,
			"risk_level":  rule.RiskLevel,
			"intervention": rule.Intervention,
		})
	}

	return policies
}

func containsIgnoreCase(s, substr string) bool {
	if len(s) == 0 {
		return false
	}
	sLower := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			sLower += string(c + 32)
		} else {
			sLower += string(c)
		}
	}
	substrLower := ""
	for _, c := range substr {
		if c >= 'A' && c <= 'Z' {
			substrLower += string(c + 32)
		} else {
			substrLower += string(c)
		}
	}
	return containsString(sLower, substrLower)
}

func containsString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

var logCounter int64
var logMu sync.Mutex

func generateLogID() string {
	logMu.Lock()
	defer logMu.Unlock()
	logCounter++
	return time.Now().Format("20060102150405") + "-" + fmt.Sprintf("%06d", logCounter)
}
