package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
	mu sync.RWMutex
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_timeout=5000&_sync=NORMAL&_busy_timeout=3000")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	d := &Database{DB: db}

	if err := d.initializeTables(); err != nil {
		return nil, fmt.Errorf("initialize tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return d, nil
}

func (d *Database) initializeTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL,
		classroom_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_login DATETIME,
		is_active BOOLEAN DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS edge_nodes (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		status TEXT NOT NULL,
		ip_address TEXT,
		classroom_id TEXT,
		cpu_usage REAL DEFAULT 0,
		memory_total INTEGER DEFAULT 0,
		memory_used INTEGER DEFAULT 0,
		memory_balloon INTEGER DEFAULT 0,
		disk_read_iops REAL DEFAULT 0,
		disk_write_iops REAL DEFAULT 0,
		disk_latency REAL DEFAULT 0,
		ksm_merged_pages INTEGER DEFAULT 0,
		ksm_memory_saved INTEGER DEFAULT 0,
		zfs_pool_name TEXT,
		zfs_used_space INTEGER DEFAULT 0,
		zfs_free_space INTEGER DEFAULT 0,
		last_heartbeat DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS vm_instances (
		id TEXT PRIMARY KEY,
		vmid INTEGER NOT NULL,
		name TEXT NOT NULL,
		node_id TEXT NOT NULL,
		classroom_id TEXT,
		student_id TEXT,
		status TEXT NOT NULL,
		ip_address TEXT,
		template_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS probe_events (
		id TEXT PRIMARY KEY,
		vm_id TEXT NOT NULL,
		node_id TEXT NOT NULL,
		student_id TEXT,
		classroom_id TEXT,
		event_type TEXT NOT NULL,
		error_message TEXT,
		error_stack TEXT,
		command TEXT,
		exit_code INTEGER,
		risk_level TEXT,
		timestamp DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS error_clusters (
		signature TEXT PRIMARY KEY,
		classroom_id TEXT NOT NULL,
		count INTEGER DEFAULT 0,
		students_count INTEGER DEFAULT 0,
		severity TEXT,
		first_seen DATETIME,
		last_seen DATETIME,
		alert_sent BOOLEAN DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id TEXT PRIMARY KEY,
		vm_id TEXT NOT NULL,
		employee_id TEXT,
		classroom_id TEXT,
		action TEXT NOT NULL,
		result TEXT NOT NULL,
		details TEXT,
		risk_score REAL DEFAULT 0,
		intervention TEXT,
		timestamp DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS clone_batches (
		batch_id TEXT PRIMARY KEY,
		template_vm_id TEXT NOT NULL,
		target_count INTEGER NOT NULL,
		success_count INTEGER DEFAULT 0,
		failed_count INTEGER DEFAULT 0,
		pending_count INTEGER DEFAULT 0,
		classroom_id TEXT,
		clone_type TEXT,
		started_at DATETIME,
		completed_at DATETIME,
		status TEXT DEFAULT 'running'
	);

	CREATE TABLE IF NOT EXISTS ai_responses (
		id TEXT PRIMARY KEY,
		student_id TEXT NOT NULL,
		vm_id TEXT NOT NULL,
		original_error TEXT,
		explanation TEXT,
		solution TEXT,
		confidence REAL,
		proactive BOOLEAN DEFAULT 0,
		timestamp DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_probe_events_vm_id ON probe_events(vm_id);
	CREATE INDEX IF NOT EXISTS idx_probe_events_classroom ON probe_events(classroom_id);
	CREATE INDEX IF NOT EXISTS idx_probe_events_timestamp ON probe_events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_error_clusters_classroom ON error_clusters(classroom_id);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_vm_id ON audit_logs(vm_id);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_vm_instances_classroom ON vm_instances(classroom_id);
	CREATE INDEX IF NOT EXISTS idx_vm_instances_student ON vm_instances(student_id);
	`

	_, err := d.DB.Exec(schema)
	return err
}

func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

func (d *Database) InsertUser(id, username, passwordHash string, role, classroomID string) error {
	_, err := d.DB.Exec(
		"INSERT OR REPLACE INTO users (id, username, password_hash, role, classroom_id) VALUES (?, ?, ?, ?, ?)",
		id, username, passwordHash, role, classroomID,
	)
	return err
}

func (d *Database) GetUserByUsername(username string) (*UserRecord, error) {
	row := d.DB.QueryRow(
		"SELECT id, username, password_hash, role, classroom_id, last_login, is_active FROM users WHERE username = ?",
		username,
	)

	var user UserRecord
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.ClassroomID, &user.LastLogin, &user.IsActive)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *Database) InsertNode(node *NodeRecord) error {
	_, err := d.DB.Exec(
		`INSERT OR REPLACE INTO edge_nodes (
			id, name, type, status, ip_address, classroom_id,
			cpu_usage, memory_total, memory_used, memory_balloon,
			disk_read_iops, disk_write_iops, disk_latency,
			ksm_merged_pages, ksm_memory_saved,
			zfs_pool_name, zfs_used_space, zfs_free_space,
			last_heartbeat
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		node.ID, node.Name, node.Type, node.Status, node.IPAddress, node.ClassroomID,
		node.CPUUsage, node.MemoryTotal, node.MemoryUsed, node.MemoryBalloon,
		node.DiskReadIOPS, node.DiskWriteIOPS, node.DiskLatency,
		node.KSMMergedPages, node.KSMMemorySaved,
		node.ZFSPoolName, node.ZFSUsedSpace, node.ZFSFreeSpace,
		node.LastHeartbeat,
	)
	return err
}

func (d *Database) UpdateNodeHeartbeat(nodeID string, heartbeat *NodeHeartbeatRecord) error {
	_, err := d.DB.Exec(
		`UPDATE edge_nodes SET
			cpu_usage = ?, memory_used = ?, memory_balloon = ?,
			disk_read_iops = ?, disk_write_iops = ?, disk_latency = ?,
			ksm_merged_pages = ?, last_heartbeat = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		heartbeat.CPUUsage, heartbeat.MemoryUsed, heartbeat.MemoryBalloon,
		heartbeat.DiskReadIOPS, heartbeat.DiskWriteIOPS, heartbeat.DiskLatency,
		heartbeat.KSMMergedPages, heartbeat.Timestamp, nodeID,
	)
	return err
}

func (d *Database) GetClassroomNodes(classroomID string) ([]NodeRecord, error) {
	rows, err := d.DB.Query(
		"SELECT id, name, type, status, ip_address, cpu_usage, memory_used, memory_total, last_heartbeat FROM edge_nodes WHERE classroom_id = ?",
		classroomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []NodeRecord
	for rows.Next() {
		var node NodeRecord
		err := rows.Scan(&node.ID, &node.Name, &node.Type, &node.Status, &node.IPAddress, &node.CPUUsage, &node.MemoryUsed, &node.MemoryTotal, &node.LastHeartbeat)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (d *Database) InsertProbeEvent(event *ProbeEventRecord) error {
	_, err := d.DB.Exec(
		`INSERT INTO probe_events (
			id, vm_id, node_id, student_id, classroom_id,
			event_type, error_message, error_stack, command, exit_code,
			risk_level, timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.VMID, event.NodeID, event.StudentID, event.ClassroomID,
		event.EventType, event.ErrorMessage, event.ErrorStack, event.Command, event.ExitCode,
		event.RiskLevel, event.Timestamp,
	)
	return err
}

func (d *Database) GetRecentProbeEvents(classroomID string, limit int) ([]ProbeEventRecord, error) {
	rows, err := d.DB.Query(
		"SELECT id, vm_id, student_id, event_type, error_message, command, exit_code, timestamp FROM probe_events WHERE classroom_id = ? ORDER BY timestamp DESC LIMIT ?",
		classroomID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []ProbeEventRecord
	for rows.Next() {
		var event ProbeEventRecord
		err := rows.Scan(&event.ID, &event.VMID, &event.StudentID, &event.EventType, &event.ErrorMessage, &event.Command, &event.ExitCode, &event.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (d *Database) InsertAuditLog(log *AuditLogRecord) error {
	_, err := d.DB.Exec(
		`INSERT INTO audit_logs (
			id, vm_id, employee_id, classroom_id, action, result,
			details, risk_score, intervention, timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		log.ID, log.VMID, log.EmployeeID, log.ClassroomID, log.Action, log.Result,
		log.Details, log.RiskScore, log.Intervention, log.Timestamp,
	)
	return err
}

func (d *Database) InsertAIResponse(response *AIResponseRecord) error {
	_, err := d.DB.Exec(
		`INSERT INTO ai_responses (
			id, student_id, vm_id, original_error, explanation,
			solution, confidence, proactive, timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		response.ID, response.StudentID, response.VMID, response.OriginalError,
		response.Explanation, response.Solution, response.Confidence,
		response.Proactive, response.Timestamp,
	)
	return err
}

func (d *Database) InsertCloneBatch(batch *CloneBatchRecord) error {
	_, err := d.DB.Exec(
		`INSERT INTO clone_batches (
			batch_id, template_vm_id, target_count, classroom_id,
			clone_type, started_at, status
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		batch.BatchID, batch.TemplateVMID, batch.TargetCount, batch.ClassroomID,
		batch.CloneType, batch.StartedAt, batch.Status,
	)
	return err
}

func (d *Database) UpdateCloneBatchStatus(batchID string, successCount, failedCount int, status string) error {
	_, err := d.DB.Exec(
		`UPDATE clone_batches SET
			success_count = ?, failed_count = ?,
			pending_count = target_count - success_count - failed_count,
			status = ?, completed_at = CURRENT_TIMESTAMP
		WHERE batch_id = ?`,
		successCount, failedCount, status, batchID,
	)
	return err
}

func (d *Database) GetClassroomStats(classroomID string) (*ClassroomStats, error) {
	var stats ClassroomStats

	row := d.DB.QueryRow(
		`SELECT
			COUNT(*) as total_nodes,
			COUNT(CASE WHEN status = 'online' THEN 1 END) as online_nodes,
			COUNT(CASE WHEN status = 'warning' THEN 1 END) as warning_nodes,
			COUNT(CASE WHEN status = 'offline' THEN 1 END) as offline_nodes,
			AVG(cpu_usage) as avg_cpu,
			AVG(CASE WHEN memory_total > 0 THEN CAST(memory_used AS REAL) / memory_total * 100 ELSE 0 END) as avg_memory
		FROM edge_nodes WHERE classroom_id = ?`,
		classroomID,
	)

	err := row.Scan(&stats.TotalNodes, &stats.OnlineNodes, &stats.WarningNodes, &stats.OfflineNodes, &stats.AvgCPU, &stats.AvgMemory)
	if err != nil {
		return nil, err
	}

	row2 := d.DB.QueryRow(
		"SELECT COUNT(*) FROM vm_instances WHERE classroom_id = ? AND status = 'running'",
		classroomID,
	)
	err = row2.Scan(&stats.RunningVMs)
	if err != nil {
		stats.RunningVMs = 0
	}

	row3 := d.DB.QueryRow(
		"SELECT COUNT(*) FROM probe_events WHERE classroom_id = ? AND timestamp > datetime('now', '-1 hour')",
		classroomID,
	)
	err = row3.Scan(&stats.RecentErrors)
	if err != nil {
		stats.RecentErrors = 0
	}

	return &stats, nil
}

type UserRecord struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
	ClassroomID  string
	LastLogin    *time.Time
	IsActive     bool
}

type NodeRecord struct {
	ID              string
	Name            string
	Type            string
	Status          string
	IPAddress       string
	ClassroomID     string
	CPUUsage        float64
	MemoryTotal     uint64
	MemoryUsed      uint64
	MemoryBalloon   uint64
	DiskReadIOPS    float64
	DiskWriteIOPS   float64
	DiskLatency     float64
	KSMMergedPages  uint64
	KSMMemorySaved  uint64
	ZFSPoolName     string
	ZFSUsedSpace    uint64
	ZFSFreeSpace    uint64
	LastHeartbeat   *time.Time
}

type NodeHeartbeatRecord struct {
	CPUUsage       float64
	MemoryUsed     uint64
	MemoryBalloon  uint64
	DiskReadIOPS   float64
	DiskWriteIOPS  float64
	DiskLatency    float64
	KSMMergedPages uint64
	Timestamp      time.Time
}

type ProbeEventRecord struct {
	ID           string
	VMID         string
	NodeID       string
	StudentID    string
	ClassroomID  string
	EventType    string
	ErrorMessage string
	ErrorStack   string
	Command      string
	ExitCode     int
	RiskLevel    string
	Timestamp    time.Time
}

type AuditLogRecord struct {
	ID          string
	VMID        string
	EmployeeID  string
	ClassroomID string
	Action      string
	Result      string
	Details     string
	RiskScore   float64
	Intervention string
	Timestamp   time.Time
}

type AIResponseRecord struct {
	ID            string
	StudentID     string
	VMID          string
	OriginalError string
	Explanation   string
	Solution      string
	Confidence    float64
	Proactive     bool
	Timestamp     time.Time
}

type CloneBatchRecord struct {
	BatchID      string
	TemplateVMID string
	TargetCount  int
	ClassroomID  string
	CloneType    string
	StartedAt    time.Time
	Status       string
}

type ClassroomStats struct {
	TotalNodes    int
	OnlineNodes   int
	WarningNodes  int
	OfflineNodes  int
	RunningVMs    int
	RecentErrors  int
	AvgCPU        float64
	AvgMemory     float64
}
