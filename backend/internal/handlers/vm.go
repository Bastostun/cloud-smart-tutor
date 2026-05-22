package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud-smart-tutor-backend/internal/domain"
	"cloud-smart-tutor-backend/internal/pve"

	"github.com/gin-gonic/gin"
)

type VMAllocateRequest struct {
	StudentIDs  []string `json:"student_ids" binding:"required"`
	ClassroomID string   `json:"classroom_id" binding:"required"`
	OSImage     string   `json:"os_image" binding:"required"`
	CPUCores    int      `json:"cpu_cores" binding:"min=1,max=16"`
	MemoryMB    int      `json:"memory_mb" binding:"min=512,max=32768"`
	DiskGB      int      `json:"disk_gb" binding:"min=10,max=500"`
	AIEnabled   bool     `json:"ai_enabled"`
	AIModel     string   `json:"ai_model"`
	AIAPIEndpoint string `json:"ai_api_endpoint"`
}

type VMConfig struct {
	CPUCores    int    `json:"cpu_cores"`
	MemoryMB    int    `json:"memory_mb"`
	DiskGB      int    `json:"disk_gb"`
	OSImage     string `json:"os_image"`
	AIEnabled   bool   `json:"ai_enabled"`
	AIModel     string `json:"ai_model"`
	AIAPIEndpoint string `json:"ai_api_endpoint"`
}

type VMPatchRequest struct {
	CPUCores    *int    `json:"cpu_cores"`
	MemoryMB    *int    `json:"memory_mb"`
	DiskGB      *int    `json:"disk_gb"`
	AIEnabled   *bool   `json:"ai_enabled"`
	AIModel     *string `json:"ai_model"`
	AIAPIEndpoint *string `json:"ai_api_endpoint"`
}

type VMListResponse struct {
	VMID      string         `json:"vm_id"`
	Name      string         `json:"name"`
	StudentID string         `json:"student_id"`
	ClassroomID string       `json:"classroom_id"`
	Status    string         `json:"status"`
	IPAddress string         `json:"ip_address"`
	Port      int            `json:"port"`
	Config    VMConfig       `json:"config"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `json:"expires_at"`
	NoVNCURL  string         `json:"no_vnc_url"`
}

type VMAllocateResponse struct {
	BatchID      string   `json:"batch_id"`
	TotalCount   int      `json:"total_count"`
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	VMs          []VMListResponse `json:"vms"`
}

func (h *Handler) AllocateVMs(c *gin.Context) {
	var req VMAllocateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	if req.CPUCores == 0 {
		req.CPUCores = 2
	}
	if req.MemoryMB == 0 {
		req.MemoryMB = 2048
	}
	if req.DiskGB == 0 {
		req.DiskGB = 30
	}
	if req.OSImage == "" {
		req.OSImage = "debian-12.0_x64"
	}

	batchID := fmt.Sprintf("vm-allocate-%d", time.Now().UnixNano())
	var allocatedVMs []VMListResponse
	successCount := 0
	failedCount := 0

	for _, studentID := range req.StudentIDs {
		vmID := fmt.Sprintf("vm-%s-%d", strings.ReplaceAll(studentID, "-", ""), time.Now().UnixNano()%100000)
		node := "pve-node-01"

		newVMID := 10000 + (successCount % 1000)
		config := map[string]interface{}{
			"cores":  req.CPUCores,
			"memory": req.MemoryMB,
			"balloon": req.MemoryMB,
			"ide0":   fmt.Sprintf("local:iso/%s.iso,media=cdrom", req.OSImage),
			"net0":   "virtio,bridge=vmbr0",
			"onboot": 1,
		}

		if err := h.PVEScheduler.PVEClient.SetVMConfig(node, newVMID, config); err != nil {
			log.Printf("Failed to configure VM %s for student %s: %v", vmID, studentID, err)
			failedCount++
			continue
		}

		if err := h.PVEScheduler.PVEClient.StartVM(node, newVMID); err != nil {
			log.Printf("Failed to start VM %s for student %s: %v", vmID, studentID, err)
			failedCount++
			continue
		}

		noVNCURL, _ := h.PVEScheduler.PVEClient.GetNoVNCURL(node, newVMID, 1280, 720)

		vm := VMListResponse{
			VMID:        vmID,
			Name:        fmt.Sprintf("student-%s", studentID),
			StudentID:   studentID,
			ClassroomID: req.ClassroomID,
			Status:      "running",
			IPAddress:   fmt.Sprintf("192.168.1.%d", 200+successCount),
			Port:        8006,
			Config: VMConfig{
				CPUCores:    req.CPUCores,
				MemoryMB:    req.MemoryMB,
				DiskGB:      req.DiskGB,
				OSImage:     req.OSImage,
				AIEnabled:   req.AIEnabled,
				AIModel:     req.AIModel,
				AIAPIEndpoint: req.AIAPIEndpoint,
			},
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
			NoVNCURL:  noVNCURL,
		}

		allocatedVMs = append(allocatedVMs, vm)
		successCount++
	}

	response := VMAllocateResponse{
		BatchID:      batchID,
		TotalCount:   len(req.StudentIDs),
		SuccessCount: successCount,
		FailedCount:  failedCount,
		VMs:          allocatedVMs,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetVMList(c *gin.Context) {
	classroomID := c.Query("classroom_id")
	studentID := c.Query("student_id")

	var vms []VMListResponse

	for _, node := range h.NodeRegistry.GetAllNodes() {
		if classroomID != "" && node.ClassroomID != classroomID {
			continue
		}

		vm := VMListResponse{
			VMID:        node.ID,
			Name:        node.Name,
			ClassroomID: node.ClassroomID,
			Status:      string(node.Status),
			Config: VMConfig{
				CPUCores: 2,
				MemoryMB: 2048,
				DiskGB:   30,
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}

		vms = append(vms, vm)
	}

	c.JSON(http.StatusOK, gin.H{"vms": vms, "total": len(vms)})
}

func (h *Handler) GetVMDetail(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	noVNCURL, _ := h.PVEScheduler.PVEClient.GetNoVNCURL(node.ID, 100, 1280, 720)

	vm := VMListResponse{
		VMID:        node.ID,
		Name:        node.Name,
		ClassroomID: node.ClassroomID,
		Status:      string(node.Status),
		IPAddress:   node.IPAddress,
		Config: VMConfig{
			CPUCores: 2,
			MemoryMB: 2048,
			DiskGB:   30,
		},
		CreatedAt: time.Now().Add(-24 * time.Hour),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		NoVNCURL:  noVNCURL,
	}

	c.JSON(http.StatusOK, vm)
}

func (h *Handler) PatchVM(c *gin.Context) {
	vmID := c.Param("id")

	var req VMPatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	config := make(map[string]interface{})
	if req.CPUCores != nil {
		config["cores"] = *req.CPUCores
	}
	if req.MemoryMB != nil {
		config["memory"] = *req.MemoryMB
		config["balloon"] = *req.MemoryMB
	}

	if len(config) > 0 {
		if err := h.PVEScheduler.PVEClient.SetVMConfig(node.ID, 100, config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update VM config: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated", "vm_id": vmID})
}

func (h *Handler) StartVM(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	if err := h.PVEScheduler.PVEClient.StartVM(node.ID, 100); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start VM: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "started", "vm_id": vmID})
}

func (h *Handler) StopVM(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	if err := h.PVEScheduler.PVEClient.StopVM(node.ID, 100); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stop VM: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "stopped", "vm_id": vmID})
}

func (h *Handler) RebootVM(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	if err := h.PVEScheduler.PVEClient.RebootVM(node.ID, 100); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reboot VM: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "rebooting", "vm_id": vmID})
}

func (h *Handler) DeleteVM(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	if err := h.PVEScheduler.PVEClient.DeleteVM(node.ID, 100); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete VM: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "vm_id": vmID})
}

func (h *Handler) GetVMNoVNCURL(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	noVNCURL, err := h.PVEScheduler.PVEClient.GetNoVNCURL(node.ID, 100, 1280, 720)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate noVNC URL: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"no_vnc_url": noVNCURL, "vm_id": vmID})
}

func (h *Handler) GetVMOSImages(c *gin.Context) {
	images := []map[string]interface{}{
		{
			"id":          "debian-12.0_x64",
			"name":        "Debian 12.0 x64",
			"description": "Linux 课程定制版",
			"version":     "12.0",
			"arch":        "x64",
			"type":        "linux",
			"size_gb":     2.1,
		},
		{
			"id":          "ubuntu-22.04_x64",
			"name":        "Ubuntu 22.04 LTS",
			"description": "标准服务器版",
			"version":     "22.04",
			"arch":        "x64",
			"type":        "linux",
			"size_gb":     1.8,
		},
		{
			"id":          "centos-9_x64",
			"name":        "CentOS Stream 9",
			"description": "企业级 Linux",
			"version":     "9",
			"arch":        "x64",
			"type":        "linux",
			"size_gb":     2.3,
		},
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (h *Handler) GetVMAvailableModels(c *gin.Context) {
	models := []map[string]interface{}{
		{
			"id":          "ollama-llama3",
			"name":        "Ollama Llama 3",
			"description": "本地部署，离线可用",
			"provider":    "ollama",
			"context_length": 8192,
		},
		{
			"id":          "ollama-qwen",
			"name":        "Ollama Qwen",
			"description": "中文优化，本地部署",
			"provider":    "ollama",
			"context_length": 4096,
		},
		{
			"id":          "custom-api",
			"name":        "自定义 API",
			"description": "对接第三方大模型",
			"provider":    "custom",
			"context_length": 0,
		},
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

func (h *Handler) GetVMFirewallRules(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	rules, err := h.PVEScheduler.PVEClient.GetFirewallRules(node.ID, 100)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"rules": []map[string]interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

func (h *Handler) CreateVMFirewallRule(c *gin.Context) {
	vmID := c.Param("id")

	var req struct {
		Direction string `json:"direction" binding:"required,oneof=in out"`
		Action    string `json:"action" binding:"required,oneof=accept drop reject"`
		Protocol  string `json:"protocol" binding:"required,oneof=tcp udp icmp"`
		DPort     string `json:"dport"`
		SPort     string `json:"sport"`
		Source    string `json:"source"`
		Dest      string `json:"dest"`
		Comment   string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	if err := h.PVEScheduler.PVEClient.CreateFirewallRule(node.ID, 100, req.Direction, req.Action, req.Protocol, req.DPort); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create firewall rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created", "vm_id": vmID})
}

func (h *Handler) DeleteVMFirewallRule(c *gin.Context) {
	vmID := c.Param("id")
	ruleIndex := 0

	if _, err := fmt.Sscanf(c.Param("ruleIndex"), "%d", &ruleIndex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule index"})
		return
	}

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	if err := h.PVEScheduler.PVEClient.DeleteFirewallRule(node.ID, 100, ruleIndex); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete firewall rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "vm_id": vmID, "rule_index": ruleIndex})
}

type VMCommandRequest struct {
	Command string   `json:"command" binding:"required"`
	Args    []string `json:"args"`
}

type VMCommandResponse struct {
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exit_code"`
}

func (h *Handler) ExecuteVMCommand(c *gin.Context) {
	vmID := c.Param("id")

	var req VMCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	command := []string{req.Command}
	command = append(command, req.Args...)

	output, err := h.PVEScheduler.PVEClient.ExecuteInVM(node.ID, 100, command)
	if err != nil {
		c.JSON(http.StatusOK, VMCommandResponse{
			Error:    err.Error(),
			ExitCode: 1,
		})
		return
	}

	c.JSON(http.StatusOK, VMCommandResponse{
		Output:   output,
		ExitCode: 0,
	})
}

func (h *Handler) GetVMNetworkInterfaces(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	networks, err := h.PVEScheduler.PVEClient.GetVMNetworks(node.ID, 100)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"interfaces": []string{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"interfaces": networks})
}

func (h *Handler) GetVMIPAddress(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	ips, err := h.PVEScheduler.PVEClient.GetVMIPAddresses(node.ID, 100)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ip_addresses": []string{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ip_addresses": ips})
}

func (h *Handler) GetVMLogs(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	logs, err := h.PVEScheduler.PVEClient.GetLogs(node.ID, 100, 0, 50)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"logs": []string{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (h *Handler) GetVMStatus(c *gin.Context) {
	vmID := c.Param("id")

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	memUsage := 0.0
	if node.MemoryTotal > 0 {
		memUsage = float64(node.MemoryUsed) / float64(node.MemoryTotal) * 100
	}

	status := map[string]interface{}{
		"vm_id":       vmID,
		"name":        node.Name,
		"status":      node.Status,
		"cpu_usage":   node.CPUUsage,
		"memory_used": node.MemoryUsed,
		"memory_total": node.MemoryTotal,
		"memory_usage": memUsage,
		"disk_read_iops": node.DiskReadIOPS,
		"disk_write_iops": node.DiskWriteIOPS,
		"last_heartbeat": node.LastHeartbeat,
	}

	c.JSON(http.StatusOK, status)
}

func (h *Handler) ReinstallVM(c *gin.Context) {
	vmID := c.Param("id")

	var req struct {
		OSImage string `json:"os_image" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	node, exists := h.NodeRegistry.GetNode(vmID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	config := map[string]interface{}{
		"ide0": fmt.Sprintf("local:iso/%s.iso,media=cdrom", req.OSImage),
	}

	if err := h.PVEScheduler.PVEClient.SetVMConfig(node.ID, 100, config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reinstall VM: " + err.Error()})
		return
	}

	if err := h.PVEScheduler.PVEClient.RebootVM(node.ID, 100); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reboot VM: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "reinstalling", "vm_id": vmID, "os_image": req.OSImage})
}
