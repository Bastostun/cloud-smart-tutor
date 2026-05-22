package pve

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"cloud-smart-tutor-backend/internal/domain"
	"cloud-smart-tutor-backend/internal/websocket"
)

type PVEClient struct {
	Endpoint   string
	Username   string
	Password   string
	Realm      string
	Ticket     string
	CSRFToken  string
	HTTPClient *http.Client
	mu         sync.RWMutex
}

type PVECredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Realm    string `json:"realm"`
}

type PVETicketResponse struct {
	Data struct {
		Ticket    string `json:"ticket"`
		CSRFToken string `json:"CSRFPreventionToken"`
	} `json:"data"`
}

func NewPVEClient(config domain.PVEClusterConfig) *PVEClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.TLSSkipVerify},
	}

	return &PVEClient{
		Endpoint: config.APIEndpoint,
		Username: config.Username,
		Password: config.Password,
		Realm:    config.Realm,
		HTTPClient: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(config.TimeoutSeconds) * time.Second,
		},
	}
}

func (c *PVEClient) Authenticate() error {
	creds := PVECredentials{
		Username: c.Username,
		Password: c.Password,
		Realm:    c.Realm,
	}

	body, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("marshal credentials: %w", err)
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/api2/json/access/ticket", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed with status %d: %s", resp.StatusCode, string(body))
	}

	var ticketResp PVETicketResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticketResp); err != nil {
		return fmt.Errorf("decode ticket: %w", err)
	}

	c.mu.Lock()
	c.Ticket = ticketResp.Data.Ticket
	c.CSRFToken = ticketResp.Data.CSRFToken
	c.mu.Unlock()

	log.Println("PVE authentication successful")
	return nil
}

func (c *PVEClient) doRequest(method, path string, body io.Reader) ([]byte, error) {
	c.mu.RLock()
	ticket := c.Ticket
	csrftoken := c.CSRFToken
	c.mu.RUnlock()

	url := c.Endpoint + "/api2/json" + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", "PVEAuthCookie="+ticket)
	req.Header.Set("CSRFPreventionToken", csrftoken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("PVE API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

type PVEClusterStatus struct {
	Nodes []PVENodeInfo `json:"nodes"`
}

type PVENodeInfo struct {
	Node   string  `json:"node"`
	Status string  `json:"status"`
	CPU    float64 `json:"cpu"`
	Mem    uint64  `json:"mem"`
	MaxMem uint64  `json:"maxmem"`
	Disk   uint64  `json:"disk"`
	MaxDisk uint64 `json:"maxdisk"`
}

type PVEVMInfo struct {
	VMID   int    `json:"vmid"`
	Name   string `json:"name"`
	Status string `json:"status"`
	CPU    float64 `json:"cpu"`
	Mem    uint64  `json:"mem"`
	MaxMem uint64  `json:"maxmem"`
}

func (c *PVEClient) GetClusterStatus() (*PVEClusterStatus, error) {
	body, err := c.doRequest("GET", "/cluster/status", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []PVENodeInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &PVEClusterStatus{Nodes: response.Data}, nil
}

func (c *PVEClient) GetNodeResources(node string) (*PVENodeInfo, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/status", node), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data PVENodeInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *PVEClient) GetVMs(node string) ([]PVEVMInfo, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/qemu", node), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []PVEVMInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *PVEClient) CloneVM(node string, vmID int, newID int, name string, cloneType string) (string, error) {
	payload := map[string]interface{}{
		"newid":  newID,
		"name":   name,
		"format": "qcow2",
	}

	if cloneType == "linked" {
		payload["full"] = 0
	} else {
		payload["full"] = 1
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	respBody, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/clone", node, vmID), bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	var response struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	return response.Data, nil
}

func (c *PVEClient) StartVM(node string, vmID int) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/status/start", node, vmID), nil)
	return err
}

func (c *PVEClient) StopVM(node string, vmID int) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/status/stop", node, vmID), nil)
	return err
}

func (c *PVEClient) SuspendVM(node string, vmID int) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/status/suspend", node, vmID), nil)
	return err
}

func (c *PVEClient) GetTaskStatus(node string, upid string) (string, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/tasks/%s/status", node, upid), nil)
	if err != nil {
		return "", err
	}

	var response struct {
		Data struct {
			Status string `json:"status"`
			ExitStatus string `json:"exitstatus"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.Data.Status, nil
}

type PVEScheduler struct {
	PVEClient *PVEClient
	Hub       *websocket.Hub
	TaskPool  chan CloneTask
	Results   map[string]*domain.CloneBatchStatus
	mu        sync.RWMutex
}

type CloneTask struct {
	BatchID      string
	Node         string
	TemplateVMID int
	NewVMID      int
	Name         string
	CloneType    string
	ResultChan   chan domain.CloneResult
}

func NewPVEScheduler(client *PVEClient, hub *websocket.Hub, maxConcurrent int) *PVEScheduler {
	return &PVEScheduler{
		PVEClient: client,
		Hub:       hub,
		TaskPool:  make(chan CloneTask, maxConcurrent),
		Results:   make(map[string]*domain.CloneBatchStatus),
	}
}

func (s *PVEScheduler) Start(workers int) {
	for i := 0; i < workers; i++ {
		go s.cloneWorker(i)
	}
	log.Printf("PVE scheduler started with %d workers", workers)
}

func (s *PVEScheduler) cloneWorker(id int) {
	log.Printf("Clone worker %d started", id)
	for task := range s.TaskPool {
		log.Printf("Worker %d processing task: node=%s, vmid=%d->%d", id, task.Node, task.TemplateVMID, task.NewVMID)

		result := domain.CloneResult{
			TaskID:       fmt.Sprintf("task-%d-%d", id, task.NewVMID),
			TargetVMID:   fmt.Sprintf("%d", task.NewVMID),
			TargetNodeID: task.Node,
		}

		upid, err := s.PVEClient.CloneVM(task.Node, task.TemplateVMID, task.NewVMID, task.Name, task.CloneType)
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
			result.CompletedAt = time.Now()
		} else {
			result.Status = "pending"
			result.TaskID = upid

			// Poll task status
			for j := 0; j < 60; j++ {
				time.Sleep(2 * time.Second)
				status, err := s.PVEClient.GetTaskStatus(task.Node, upid)
				if err != nil {
					continue
				}
				if status == "stopped" {
					result.Status = "success"
					result.CompletedAt = time.Now()
					break
				}
			}
			if result.Status == "pending" {
				result.Status = "timeout"
				result.Error = "clone task timed out"
				result.CompletedAt = time.Now()
			}
		}

		s.updateBatchStatus(task.BatchID, result)
		task.ResultChan <- result
	}
}

func (s *PVEScheduler) CloneBatch(request domain.VMCloneRequest) string {
	batchID := fmt.Sprintf("batch-%d", time.Now().UnixNano())

	s.mu.Lock()
	s.Results[batchID] = &domain.CloneBatchStatus{
		BatchID:    batchID,
		TotalCount: request.TargetCount,
		StartedAt:  time.Now(),
		Results:    make([]domain.CloneResult, 0, request.TargetCount),
	}
	s.mu.Unlock()

	resultChan := make(chan domain.CloneResult, request.TargetCount)

	for i := 0; i < request.TargetCount; i++ {
		newVMID := 10000 + i
		targetNode := ""
		if len(request.TargetNodeIDs) > 0 {
			targetNode = request.TargetNodeIDs[i%len(request.TargetNodeIDs)]
		}

		task := CloneTask{
			BatchID:      batchID,
			Node:         targetNode,
			TemplateVMID: 100,
			NewVMID:      newVMID,
			Name:         fmt.Sprintf("clone-%s-%d", request.ClassroomID, i),
			CloneType:    request.CloneType,
			ResultChan:   resultChan,
		}

		select {
		case s.TaskPool <- task:
		default:
			log.Printf("Task pool full, dropping clone task for VM %d", newVMID)
		}
	}

	go s.collectBatchResults(batchID, resultChan, request.TargetCount)

	return batchID
}

func (s *PVEScheduler) collectBatchResults(batchID string, resultChan chan domain.CloneResult, totalCount int) {
	collected := 0
	for result := range resultChan {
		s.mu.Lock()
		if batch, exists := s.Results[batchID]; exists {
			batch.Results = append(batch.Results, result)
			if result.Status == "success" {
				batch.SuccessCount++
			} else {
				batch.FailedCount++
			}
			batch.PendingCount = batch.TotalCount - batch.SuccessCount - batch.FailedCount
		}
		s.mu.Unlock()

		collected++
		if collected >= totalCount {
			s.mu.Lock()
			if batch, exists := s.Results[batchID]; exists {
				now := time.Now()
				batch.CompletedAt = &now
			}
			s.mu.Unlock()
			close(resultChan)
			break
		}
	}

	s.mu.RLock()
	if batch, exists := s.Results[batchID]; exists {
		s.Hub.SendToRoom("admin", &websocket.WSMessage{
			Type:      websocket.WSCloneStatus,
			Timestamp: time.Now(),
			Payload:   batch,
		})
	}
	s.mu.RUnlock()
}

func (s *PVEScheduler) GetBatchStatus(batchID string) (*domain.CloneBatchStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	status, exists := s.Results[batchID]
	return status, exists
}
