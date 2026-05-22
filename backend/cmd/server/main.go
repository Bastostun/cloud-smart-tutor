package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud-smart-tutor-backend/internal/cluster"
	"cloud-smart-tutor-backend/internal/config"
	"cloud-smart-tutor-backend/internal/domain"
	"cloud-smart-tutor-backend/internal/handlers"
	"cloud-smart-tutor-backend/internal/monitor"
	"cloud-smart-tutor-backend/internal/probe"
	"cloud-smart-tutor-backend/internal/pve"
	"cloud-smart-tutor-backend/internal/sandbox"
	"cloud-smart-tutor-backend/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("===========================================")
	log.Println("  云端智师 - 微型 IDC 高并发中枢启动中...")
	log.Println("===========================================")

	cfg := config.LoadConfig()

	gin.SetMode(cfg.Mode)

	hub := websocket.NewHub()
	go hub.Run()
	log.Println("WebSocket Hub started")

	pveConfig := domain.PVEClusterConfig{
		APIEndpoint:   cfg.PVEEndpoint,
		Username:      cfg.PVEUsername,
		Password:      cfg.PVEPassword,
		Realm:         cfg.PVERealm,
		TLSSkipVerify: cfg.PVETLSSkip,
		TimeoutSeconds: cfg.PVETimeout,
	}
	pveClient := pve.NewPVEClient(pveConfig)

	if err := pveClient.Authenticate(); err != nil {
		log.Printf("PVE authentication failed (will retry): %v", err)
	}

	pveScheduler := pve.NewPVEScheduler(pveClient, hub, cfg.MaxConcurrentClones)
	pveScheduler.Start(cfg.PVEWorkers)
	log.Printf("PVE Scheduler started with %d workers", cfg.PVEWorkers)

	nodeRegistry := monitor.NewNodeRegistry(hub, pveClient)
	nodeRegistry.StartPolling(cfg.NodePollInterval)
	log.Println("Node Registry started with polling")

	errorEngine := cluster.NewErrorClusteringEngine(hub, cfg.ErrorAlertThreshold)
	go func() {
		ticker := time.NewTicker(cfg.ErrorCleanupMaxAge)
		defer ticker.Stop()
		for range ticker.C {
			errorEngine.CleanupExpired(cfg.ErrorCleanupMaxAge)
		}
	}()
	log.Println("Error Clustering Engine started")

	auditStream := sandbox.NewAuditStream(hub, pveClient)
	log.Println("Sandbox Audit Stream started")

	probeIngestor := probe.NewProbeIngestor(
		hub,
		errorEngine,
		auditStream,
		cfg.DefaultTotalStudents,
	)
	probeIngestor.SetTotalStudents(cfg.DefaultTotalStudents)
	log.Println("Probe Ingestor started")

	handler := handlers.NewHandler(
		hub,
		nodeRegistry,
		pveScheduler,
		errorEngine,
		auditStream,
		probeIngestor,
	)

	r := gin.Default()
	r.Use(corsMiddleware())
	handler.RegisterRoutes(r)
	log.Println("HTTP routes registered")

	registerDemoNodes(nodeRegistry, cfg)

	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.ServerAddress, cfg.ServerPort)
		log.Printf("Server starting on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	log.Println("云端智师 backend stopped")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func registerDemoNodes(registry *monitor.NodeRegistry, cfg *config.Config) {
	demoNodes := []domain.EdgeNode{
		{ID: "node-teacher-01", Name: "教师机", Type: domain.NodeTeacherServer, Status: domain.NodeOnline, IPAddress: "192.168.1.1", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 45, MemoryTotal: 32000000000, MemoryUsed: 20000000000},
		{ID: "node-student-01", Name: "学生-01", Type: domain.NodeStudentVM, Status: domain.NodeOnline, IPAddress: "192.168.1.101", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 78, MemoryTotal: 4000000000, MemoryUsed: 2500000000},
		{ID: "node-student-02", Name: "学生-02", Type: domain.NodeStudentVM, Status: domain.NodeOnline, IPAddress: "192.168.1.102", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 32, MemoryTotal: 4000000000, MemoryUsed: 2000000000},
		{ID: "node-student-03", Name: "学生-03", Type: domain.NodeStudentVM, Status: domain.NodeWarning, IPAddress: "192.168.1.103", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 87, MemoryTotal: 4000000000, MemoryUsed: 3200000000},
		{ID: "node-student-04", Name: "学生-04", Type: domain.NodeStudentVM, Status: domain.NodeOnline, IPAddress: "192.168.1.104", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 23, MemoryTotal: 4000000000, MemoryUsed: 1500000000},
		{ID: "node-student-05", Name: "学生-05", Type: domain.NodeStudentVM, Status: domain.NodeOnline, IPAddress: "192.168.1.105", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 56, MemoryTotal: 4000000000, MemoryUsed: 2800000000},
		{ID: "node-student-06", Name: "学生-06", Type: domain.NodeStudentVM, Status: domain.NodeOffline, IPAddress: "192.168.1.106", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 0, MemoryTotal: 4000000000, MemoryUsed: 0},
		{ID: "node-student-07", Name: "学生-07", Type: domain.NodeStudentVM, Status: domain.NodeOnline, IPAddress: "192.168.1.107", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 67, MemoryTotal: 4000000000, MemoryUsed: 2400000000},
		{ID: "node-student-08", Name: "学生-08", Type: domain.NodeStudentVM, Status: domain.NodeWarning, IPAddress: "192.168.1.108", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 82, MemoryTotal: 4000000000, MemoryUsed: 3100000000},
		{ID: "node-student-09", Name: "学生-09", Type: domain.NodeStudentVM, Status: domain.NodeOnline, IPAddress: "192.168.1.109", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 41, MemoryTotal: 4000000000, MemoryUsed: 2100000000},
		{ID: "node-student-10", Name: "学生-10", Type: domain.NodeStudentVM, Status: domain.NodeOnline, IPAddress: "192.168.1.110", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 29, MemoryTotal: 4000000000, MemoryUsed: 1800000000},
		{ID: "node-student-11", Name: "学生-11", Type: domain.NodeStudentVM, Status: domain.NodeOnline, IPAddress: "192.168.1.111", ClassroomID: cfg.DefaultClassroomID, CPUUsage: 53, MemoryTotal: 4000000000, MemoryUsed: 2600000000},
	}

	for _, node := range demoNodes {
		registry.RegisterNode(node)
	}

	log.Printf("Registered %d demo nodes for classroom %s", len(demoNodes), cfg.DefaultClassroomID)
}
