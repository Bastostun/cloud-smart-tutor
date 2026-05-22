
package domain

import "time"

type NodeStatus string

const (
	NodeOnline    NodeStatus = "online"
	NodeWarning   NodeStatus = "warning"
	NodeOffline   NodeStatus = "offline"
	NodeScheduled NodeStatus = "scheduled"
)

type NodeType string

const (
	NodeTeacherServer NodeType = "teacher_server"
	NodeStudentVM     NodeType = "student_vm"
	NodeEnterpriseSandbox NodeType = "enterprise_sandbox"
	NodeStoragePool   NodeType = "storage_pool"
)

type EdgeNode struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Type            NodeType   `json:"type"`
	Status          NodeStatus `json:"status"`
	IPAddress       string     `json:"ip_address"`
	ClassroomID     string     `json:"classroom_id"`
	
	// Resource metrics
	CPUUsage        float64    `json:"cpu_usage"`
	MemoryTotal     uint64     `json:"memory_total"`
	MemoryUsed      uint64     `json:"memory_used"`
	MemoryBalloon   uint64     `json:"memory_balloon"`
	DiskReadIOPS    float64    `json:"disk_read_iops"`
	DiskWriteIOPS   float64    `json:"disk_write_iops"`
	DiskReadThroughput  float64 `json:"disk_read_throughput"`
	DiskWriteThroughput float64 `json:"disk_write_throughput"`
	DiskLatency     float64    `json:"disk_latency"`
	
	// KSM metrics
	KSMMergedPages  uint64     `json:"ksm_merged_pages"`
	KSMMemorySaved  uint64     `json:"ksm_memory_saved"`
	
	// ZFS metrics
	ZFSPoolName     string     `json:"zfs_pool_name"`
	ZFSUsedSpace    uint64     `json:"zfs_used_space"`
	ZFSFreeSpace    uint64     `json:"zfs_free_space"`
	ZFSCompressionRatio float64 `json:"zfs_compression_ratio"`
	
	// Heartbeat
	LastHeartbeat   time.Time  `json:"last_heartbeat"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type NodeHeartbeat struct {
	NodeID          string    `json:"node_id"`
	Timestamp       time.Time `json:"timestamp"`
	CPUUsage        float64   `json:"cpu_usage"`
	MemoryUsed      uint64    `json:"memory_used"`
	MemoryBalloon   uint64    `json:"memory_balloon"`
	DiskReadIOPS    float64   `json:"disk_read_iops"`
	DiskWriteIOPS   float64   `json:"disk_write_iops"`
	DiskLatency     float64   `json:"disk_latency"`
	KSMMergedPages  uint64    `json:"ksm_merged_pages"`
}

type VMCloneRequest struct {
	TemplateVMID    string   `json:"template_vm_id"`
	TargetCount     int      `json:"target_count"`
	TargetNodeIDs   []string `json:"target_node_ids"`
	ClassroomID     string   `json:"classroom_id"`
	ResourcePool    string   `json:"resource_pool"`
	CloneType       string   `json:"clone_type"` // "linked" or "full"
	NetworkConfig   NetworkConfig `json:"network_config"`
}

type NetworkConfig struct {
	Bridge          string   `json:"bridge"`
	VLANID          int      `json:"vlan_id"`
	IPPrefix        string   `json:"ip_prefix"`
	DNSServers      []string `json:"dns_servers"`
	Gateway         string   `json:"gateway"`
}

type CloneResult struct {
	TaskID          string    `json:"task_id"`
	TargetVMID      string    `json:"target_vm_id"`
	TargetNodeID    string    `json:"target_node_id"`
	Status          string    `json:"status"` // "success", "failed", "pending"
	Error           string    `json:"error,omitempty"`
	CompletedAt     time.Time `json:"completed_at"`
}

type CloneBatchStatus struct {
	BatchID         string        `json:"batch_id"`
	TotalCount      int           `json:"total_count"`
	SuccessCount    int           `json:"success_count"`
	FailedCount     int           `json:"failed_count"`
	PendingCount    int           `json:"pending_count"`
	StartedAt       time.Time     `json:"started_at"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
	Results         []CloneResult `json:"results"`
}

type PVEClusterConfig struct {
	APIEndpoint     string `json:"api_endpoint"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Realm           string `json:"realm"`
	TLSSkipVerify   bool   `json:"tls_skip_verify"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
}
