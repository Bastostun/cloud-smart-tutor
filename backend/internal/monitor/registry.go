package monitor

import (
	"log"
	"sync"
	"time"

	"cloud-smart-tutor-backend/internal/domain"
	"cloud-smart-tutor-backend/internal/pve"
	"cloud-smart-tutor-backend/internal/websocket"
)

type NodeRegistry struct {
	Nodes       map[string]*domain.EdgeNode
	Classrooms  map[string][]string
	Hub         *websocket.Hub
	PVEClient   *pve.PVEClient
	mu          sync.RWMutex
}

func NewNodeRegistry(hub *websocket.Hub, pveClient *pve.PVEClient) *NodeRegistry {
	return &NodeRegistry{
		Nodes:      make(map[string]*domain.EdgeNode),
		Classrooms: make(map[string][]string),
		Hub:        hub,
		PVEClient:  pveClient,
	}
}

func (nr *NodeRegistry) RegisterNode(node domain.EdgeNode) {
	nr.mu.Lock()
	defer nr.mu.Unlock()

	node.CreatedAt = time.Now()
	node.UpdatedAt = time.Now()
	nr.Nodes[node.ID] = &node

	if node.ClassroomID != "" {
		found := false
		for _, id := range nr.Classrooms[node.ClassroomID] {
			if id == node.ID {
				found = true
				break
			}
		}
		if !found {
			nr.Classrooms[node.ClassroomID] = append(nr.Classrooms[node.ClassroomID], node.ID)
		}
	}

	log.Printf("Node registered: id=%s, name=%s, classroom=%s", node.ID, node.Name, node.ClassroomID)
}

func (nr *NodeRegistry) UpdateHeartbeat(heartbeat domain.NodeHeartbeat) {
	nr.mu.Lock()
	defer nr.mu.Unlock()

	node, exists := nr.Nodes[heartbeat.NodeID]
	if !exists {
		log.Printf("Received heartbeat for unknown node: %s", heartbeat.NodeID)
		return
	}

	node.LastHeartbeat = heartbeat.Timestamp
	node.CPUUsage = heartbeat.CPUUsage
	node.MemoryUsed = heartbeat.MemoryUsed
	node.MemoryBalloon = heartbeat.MemoryBalloon
	node.DiskReadIOPS = heartbeat.DiskReadIOPS
	node.DiskWriteIOPS = heartbeat.DiskWriteIOPS
	node.DiskLatency = heartbeat.DiskLatency
	node.KSMMergedPages = heartbeat.KSMMergedPages
	node.UpdatedAt = time.Now()

	if heartbeat.CPUUsage > 85 {
		node.Status = domain.NodeWarning
	} else {
		node.Status = domain.NodeOnline
	}
}

func (nr *NodeRegistry) GetNode(nodeID string) (*domain.EdgeNode, bool) {
	nr.mu.RLock()
	defer nr.mu.RUnlock()
	node, exists := nr.Nodes[nodeID]
	return node, exists
}

func (nr *NodeRegistry) GetAllNodes() []*domain.EdgeNode {
	nr.mu.RLock()
	defer nr.mu.RUnlock()

	nodes := make([]*domain.EdgeNode, 0, len(nr.Nodes))
	for _, node := range nr.Nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (nr *NodeRegistry) GetClassroomNodes(classroomID string) []*domain.EdgeNode {
	nr.mu.RLock()
	defer nr.mu.RUnlock()

	var nodes []*domain.EdgeNode
	for _, nodeID := range nr.Classrooms[classroomID] {
		if node, exists := nr.Nodes[nodeID]; exists {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (nr *NodeRegistry) StartPolling(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			nr.pollAllNodes()
			nr.broadcastMetrics()
			nr.checkStaleNodes()
		}
	}()

	log.Printf("Node polling started with interval %v", interval)
}

func (nr *NodeRegistry) pollAllNodes() {
	if nr.PVEClient == nil {
		return
	}

	status, err := nr.PVEClient.GetClusterStatus()
	if err != nil {
		log.Printf("Failed to poll PVE cluster: %v", err)
		return
	}

	nr.mu.Lock()
	for _, pveNode := range status.Nodes {
		node, exists := nr.Nodes[pveNode.Node]
		if !exists {
			continue
		}

		node.CPUUsage = pveNode.CPU * 100
		node.MemoryUsed = pveNode.Mem
		node.MemoryTotal = pveNode.MaxMem
		node.LastHeartbeat = time.Now()
		node.Status = domain.NodeStatus(pveNode.Status)
	}
	nr.mu.Unlock()
}

func (nr *NodeRegistry) broadcastMetrics() {
	nr.mu.RLock()
	defer nr.mu.RUnlock()

	metrics := make([]map[string]interface{}, 0, len(nr.Nodes))
	for _, node := range nr.Nodes {
		memUsage := 0.0
		if node.MemoryTotal > 0 {
			memUsage = float64(node.MemoryUsed) / float64(node.MemoryTotal) * 100
		}

		metrics = append(metrics, map[string]interface{}{
			"id":                  node.ID,
			"name":                node.Name,
			"type":                node.Type,
			"status":              node.Status,
			"classroom_id":        node.ClassroomID,
			"cpu_usage":           node.CPUUsage,
			"memory_usage":        memUsage,
			"memory_balloon":      node.MemoryBalloon,
			"disk_read_iops":      node.DiskReadIOPS,
			"disk_write_iops":     node.DiskWriteIOPS,
			"disk_latency":        node.DiskLatency,
			"ksm_merged_pages":    node.KSMMergedPages,
			"last_heartbeat":      node.LastHeartbeat,
		})
	}

	nr.Hub.SendToRoom("idc-dashboard", &websocket.WSMessage{
		Type:      websocket.WSNodeMetrics,
		Timestamp: time.Now(),
		Payload:   metrics,
	})
}

func (nr *NodeRegistry) checkStaleNodes() {
	nr.mu.Lock()
	defer nr.mu.Unlock()

	staleThreshold := 60 * time.Second
	now := time.Now()

	for _, node := range nr.Nodes {
		if node.Status == domain.NodeOffline {
			continue
		}

		if now.Sub(node.LastHeartbeat) > staleThreshold {
			node.Status = domain.NodeOffline
			log.Printf("Node marked as offline: id=%s, name=%s", node.ID, node.Name)
		}
	}
}

func (nr *NodeRegistry) GetClassroomStats(classroomID string) map[string]interface{} {
	nodes := nr.GetClassroomNodes(classroomID)

	total := len(nodes)
	online := 0
	warning := 0
	offline := 0

	var totalCPU float64
	var totalMemory float64

	for _, node := range nodes {
		switch node.Status {
		case domain.NodeOnline:
			online++
		case domain.NodeWarning:
			online++
			warning++
		case domain.NodeOffline:
			offline++
		}
		totalCPU += node.CPUUsage
		if node.MemoryTotal > 0 {
			totalMemory += float64(node.MemoryUsed) / float64(node.MemoryTotal) * 100
		}
	}

	avgCPU := 0.0
	avgMemory := 0.0
	if total > 0 {
		avgCPU = totalCPU / float64(total)
		avgMemory = totalMemory / float64(total)
	}

	return map[string]interface{}{
		"classroom_id": classroomID,
		"total":        total,
		"online":       online,
		"warning":      warning,
		"offline":      offline,
		"avg_cpu":      avgCPU,
		"avg_memory":   avgMemory,
	}
}
