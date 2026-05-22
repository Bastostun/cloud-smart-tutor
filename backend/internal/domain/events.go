package domain

import "time"

type ProbeEventType string

const (
	ProbeStderrCapture    ProbeEventType = "stderr_capture"
	ProbeExitCode         ProbeEventType = "exit_code"
	ProbeFileAccess       ProbeEventType = "file_access"
	ProbeNetworkAttempt   ProbeEventType = "network_attempt"
	ProbeClipboardCopy    ProbeEventType = "clipboard_copy"
	ProbeDLPViolation     ProbeEventType = "dlp_violation"
	ProbeAnomalyDetected  ProbeEventType = "anomaly_detected"
	ProbeHeartbeat        ProbeEventType = "heartbeat"
)

type ProbeEvent struct {
	EventID         string        `json:"event_id"`
	VMID            string        `json:"vm_id"`
	NodeID          string        `json:"node_id"`
	StudentID       string        `json:"student_id"`
	ClassroomID     string        `json:"classroom_id"`
	EventType       ProbeEventType `json:"event_type"`
	Timestamp       time.Time     `json:"timestamp"`
	
	// Error/stderr capture
	ErrorMessage    string        `json:"error_message,omitempty"`
	ErrorStack      string        `json:"error_stack,omitempty"`
	Command         string        `json:"command,omitempty"`
	ExitCode        int           `json:"exit_code,omitempty"`
	
	// Security events (sandbox)
	TargetPath      string        `json:"target_path,omitempty"`
	TargetURL       string        `json:"target_url,omitempty"`
	ClipboardContent string       `json:"clipboard_content,omitempty"`
	RiskLevel       string        `json:"risk_level"` // "low", "medium", "high", "critical"
	
	// Metadata
	ProcessName     string        `json:"process_name,omitempty"`
	User            string        `json:"user,omitempty"`
	WorkingDir      string        `json:"working_dir,omitempty"`
}

type AITeachingSuggestion struct {
	SuggestionID    string        `json:"suggestion_id"`
	ErrorClusterID  string        `json:"error_cluster_id"`
	ClassroomID     string        `json:"classroom_id"`
	Timestamp       time.Time     `json:"timestamp"`
	
	Title           string        `json:"title"`
	Content         string        `json:"content"`
	Confidence      float64       `json:"confidence"`
	ActionType      string        `json:"action_type"` // "explain", "demo", "handout", "pause"
	Priority        string        `json:"priority"` // "low", "medium", "high", "urgent"
}

type AICompanionResponse struct {
	ResponseID      string        `json:"response_id"`
	StudentID       string        `json:"student_id"`
	VMID            string        `json:"vm_id"`
	Timestamp       time.Time     `json:"timestamp"`
	
	OriginalError   string        `json:"original_error"`
	Explanation     string        `json:"explanation"`
	Solution        string        `json:"solution"`
	Confidence      float64       `json:"confidence"`
	Proactive       bool          `json:"proactive"` // true if AI initiated without student request
}

type AuditLogEntry struct {
	LogID           string        `json:"log_id"`
	VMID            string        `json:"vm_id"`
	EmployeeID      string        `json:"employee_id"`
	ClassroomID     string        `json:"classroom_id"`
	Timestamp       time.Time     `json:"timestamp"`
	
	Action          string        `json:"action"`
	Result          string        `json:"result"` // "allowed", "blocked", "warned"
	Details         string        `json:"details"`
	RiskScore       float64       `json:"risk_score"`
	Intervention    string        `json:"intervention,omitempty"` // "none", "warn", "block", "suspend_vm"
}
