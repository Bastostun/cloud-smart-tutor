package handlers

import (
	"log"
	"net/http"

	"cloud-smart-tutor-backend/internal/cluster"
	"cloud-smart-tutor-backend/internal/domain"
	"cloud-smart-tutor-backend/internal/monitor"
	"cloud-smart-tutor-backend/internal/probe"
	"cloud-smart-tutor-backend/internal/pve"
	"cloud-smart-tutor-backend/internal/sandbox"
	"cloud-smart-tutor-backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Hub           *websocket.Hub
	NodeRegistry  *monitor.NodeRegistry
	PVEScheduler  *pve.PVEScheduler
	ErrorEngine   *cluster.ErrorClusteringEngine
	AuditStream   *sandbox.AuditStream
	ProbeIngestor *probe.ProbeIngestor
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHandler(
	hub *websocket.Hub,
	nodeRegistry *monitor.NodeRegistry,
	pveScheduler *pve.PVEScheduler,
	errorEngine *cluster.ErrorClusteringEngine,
	auditStream *sandbox.AuditStream,
	probeIngestor *probe.ProbeIngestor,
) *Handler {
	return &Handler{
		Hub:           hub,
		NodeRegistry:  nodeRegistry,
		PVEScheduler:  pveScheduler,
		ErrorEngine:   errorEngine,
		AuditStream:   auditStream,
		ProbeIngestor: probeIngestor,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/probe/event", h.ReceiveProbeEvent)
		api.POST("/heartbeat", h.ReceiveHeartbeat)
		api.POST("/nodes/register", h.RegisterNode)

		api.GET("/nodes", h.GetAllNodes)
		api.GET("/nodes/:id", h.GetNode)
		api.GET("/classroom/:id/stats", h.GetClassroomStats)
		api.GET("/classroom/:id/errors", h.GetClassroomErrors)

		api.GET("/vm/list", h.GetVMList)
		api.GET("/vm/:id", h.GetVMDetail)
		api.POST("/vm/allocate", h.AllocateVMs)
		api.PATCH("/vm/:id", h.PatchVM)
		api.POST("/vm/:id/start", h.StartVM)
		api.POST("/vm/:id/stop", h.StopVM)
		api.POST("/vm/:id/reboot", h.RebootVM)
		api.DELETE("/vm/:id", h.DeleteVM)
		api.GET("/vm/:id/noVNC", h.GetVMNoVNCURL)
		api.GET("/vm/:id/status", h.GetVMStatus)
		api.POST("/vm/:id/reinstall", h.ReinstallVM)
		api.GET("/vm/:id/command", h.ExecuteVMCommand)
		api.GET("/vm/:id/network", h.GetVMNetworkInterfaces)
		api.GET("/vm/:id/ip", h.GetVMIPAddress)
		api.GET("/vm/:id/logs", h.GetVMLogs)
		api.GET("/vm/:id/firewall", h.GetVMFirewallRules)
		api.POST("/vm/:id/firewall", h.CreateVMFirewallRule)
		api.DELETE("/vm/:id/firewall/:ruleIndex", h.DeleteVMFirewallRule)
		api.GET("/vm/os-images", h.GetVMOSImages)
		api.GET("/vm/ai-models", h.GetVMAvailableModels)

		api.POST("/clone", h.CloneVMs)
		api.GET("/clone/:batchId/status", h.GetCloneStatus)

		api.GET("/sandbox/policies", h.GetSandboxPolicies)

		api.GET("/health", h.HealthCheck)
	}

	r.GET("/ws/idc", h.HandleIDCDashboardWS)
	r.GET("/ws/teacher/:classroomId", h.HandleTeacherWS)
	r.GET("/ws/student/:studentId", h.HandleStudentWS)
	r.GET("/ws/enterprise/:classroomId", h.HandleEnterpriseWS)
	r.GET("/ws/probe", h.HandleProbeWS)
}

func (h *Handler) HandleIDCDashboardWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed (idc): %v", err)
		return
	}

	client := &websocket.Client{
		Hub:      h.Hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		ClientID: "idc-dashboard-" + c.ClientIP(),
		Room:     "idc-dashboard",
		Role:     "admin",
	}

	h.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) HandleTeacherWS(c *gin.Context) {
	classroomID := c.Param("classroomId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed (teacher): %v", err)
		return
	}

	client := &websocket.Client{
		Hub:      h.Hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		ClientID: "teacher-" + classroomID + "-" + c.ClientIP(),
		Room:     "teacher:" + classroomID,
		Role:     "teacher",
	}

	h.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) HandleStudentWS(c *gin.Context) {
	studentID := c.Param("studentId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed (student): %v", err)
		return
	}

	client := &websocket.Client{
		Hub:      h.Hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		ClientID: "student-" + studentID,
		Room:     "student:" + studentID,
		Role:     "student",
	}

	h.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) HandleEnterpriseWS(c *gin.Context) {
	classroomID := c.Param("classroomId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed (enterprise): %v", err)
		return
	}

	client := &websocket.Client{
		Hub:      h.Hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		ClientID: "enterprise-" + classroomID + "-" + c.ClientIP(),
		Room:     "enterprise:" + classroomID,
		Role:     "enterprise",
	}

	h.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) HandleProbeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed (probe): %v", err)
		return
	}

	client := &websocket.Client{
		Hub:      h.Hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		ClientID: "probe-" + c.ClientIP(),
		Room:     "probe",
		Role:     "probe",
	}

	h.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) ReceiveProbeEvent(c *gin.Context) {
	var event map[string]interface{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received probe event: %v", event)

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *Handler) ReceiveHeartbeat(c *gin.Context) {
	var heartbeat struct {
		NodeID        string  `json:"node_id" binding:"required"`
		CPUUsage      float64 `json:"cpu_usage"`
		MemoryUsed    uint64  `json:"memory_used"`
		MemoryBalloon uint64  `json:"memory_balloon"`
		DiskReadIOPS  float64 `json:"disk_read_iops"`
		DiskWriteIOPS float64 `json:"disk_write_iops"`
		DiskLatency   float64 `json:"disk_latency"`
	}

	if err := c.ShouldBindJSON(&heartbeat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.NodeRegistry.UpdateHeartbeat(monitor.NodeHeartbeat{
		NodeID:        heartbeat.NodeID,
		CPUUsage:      heartbeat.CPUUsage,
		MemoryUsed:    heartbeat.MemoryUsed,
		MemoryBalloon: heartbeat.MemoryBalloon,
		DiskReadIOPS:  heartbeat.DiskReadIOPS,
		DiskWriteIOPS: heartbeat.DiskWriteIOPS,
		DiskLatency:   heartbeat.DiskLatency,
	})

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) RegisterNode(c *gin.Context) {
	var node struct {
		ID          string `json:"id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type"`
		IPAddress   string `json:"ip_address"`
		ClassroomID string `json:"classroom_id"`
	}

	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.NodeRegistry.RegisterNode(domain.EdgeNode{
		ID:          node.ID,
		Name:        node.Name,
		Type:        domain.NodeType(node.Type),
		IPAddress:   node.IPAddress,
		ClassroomID: node.ClassroomID,
	})

	c.JSON(http.StatusOK, gin.H{"status": "registered", "node_id": node.ID})
}

func (h *Handler) GetAllNodes(c *gin.Context) {
	nodes := h.NodeRegistry.GetAllNodes()
	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

func (h *Handler) GetNode(c *gin.Context) {
	nodeID := c.Param("id")
	node, exists := h.NodeRegistry.GetNode(nodeID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"node": node})
}

func (h *Handler) GetClassroomStats(c *gin.Context) {
	classroomID := c.Param("id")
	stats := h.NodeRegistry.GetClassroomStats(classroomID)
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *Handler) GetClassroomErrors(c *gin.Context) {
	classroomID := c.Param("id")
	errors := h.ErrorEngine.GetTopErrors(classroomID, 5)
	c.JSON(http.StatusOK, gin.H{"errors": errors})
}

func (h *Handler) CloneVMs(c *gin.Context) {
	var request struct {
		TemplateVMID  string   `json:"template_vm_id" binding:"required"`
		TargetCount   int      `json:"target_count" binding:"required,min=1,max=100"`
		TargetNodeIDs []string `json:"target_node_ids"`
		ClassroomID   string   `json:"classroom_id" binding:"required"`
		CloneType     string   `json:"clone_type"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.CloneType == "" {
		request.CloneType = "linked"
	}

	cloneReq := pve.VMCloneRequest{
		TemplateVMID:  request.TemplateVMID,
		TargetCount:   request.TargetCount,
		TargetNodeIDs: request.TargetNodeIDs,
		ClassroomID:   request.ClassroomID,
		CloneType:     request.CloneType,
	}

	batchID := h.PVEScheduler.CloneBatch(cloneReq)

	c.JSON(http.StatusOK, gin.H{
		"status":   "cloning_started",
		"batch_id": batchID,
		"count":    request.TargetCount,
	})
}

func (h *Handler) GetCloneStatus(c *gin.Context) {
	batchID := c.Param("batchId")
	status, exists := h.PVEScheduler.GetBatchStatus(batchID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (h *Handler) GetSandboxPolicies(c *gin.Context) {
	policies := h.AuditStream.GetPolicies()
	c.JSON(http.StatusOK, gin.H{"policies": policies})
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "cloud-smart-tutor-backend",
		"version": "1.0.0",
	})
}
