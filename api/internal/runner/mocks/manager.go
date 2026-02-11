package mocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/georgi-georgiev/testmesh/internal/storage/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Manager manages mock server instances
type Manager struct {
	repo     *repository.MockRepository
	logger   *zap.Logger
	servers  map[uuid.UUID]*ServerInstance
	portPool *PortPool
	mu       sync.RWMutex
}

// ServerInstance represents a running mock server
type ServerInstance struct {
	ID       uuid.UUID
	Server   *http.Server
	Matcher  *EndpointMatcher
	State    *StateManager
	Listener net.Listener
}

// NewManager creates a new mock server manager
func NewManager(repo *repository.MockRepository, logger *zap.Logger) *Manager {
	return &Manager{
		repo:     repo,
		logger:   logger,
		servers:  make(map[uuid.UUID]*ServerInstance),
		portPool: NewPortPool(10000, 20000), // Port range 10000-20000
	}
}

// StartServer starts a new mock server
func (m *Manager) StartServer(ctx context.Context, serverID uuid.UUID, name string, executionID *uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if server already exists
	if _, exists := m.servers[serverID]; exists {
		return fmt.Errorf("server %s already running", serverID)
	}

	// Allocate port
	port, err := m.portPool.Allocate(serverID)
	if err != nil {
		return fmt.Errorf("failed to allocate port: %w", err)
	}

	// Create server record
	server := &models.MockServer{
		ID:          serverID,
		ExecutionID: executionID,
		Name:        name,
		Port:        port,
		BaseURL:     fmt.Sprintf("http://localhost:%d", port),
		Status:      models.MockServerStatusStarting,
	}

	if err := m.repo.CreateServer(server); err != nil {
		m.portPool.Release(serverID)
		return fmt.Errorf("failed to create server record: %w", err)
	}

	// Load endpoints
	endpoints, err := m.repo.ListEndpoints(serverID)
	if err != nil {
		return fmt.Errorf("failed to load endpoints: %w", err)
	}

	// Create matcher and state manager
	matcher := NewEndpointMatcher(endpoints, m.logger)
	stateManager := NewStateManager(serverID, m.repo, m.logger)

	// Create HTTP server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		m.portPool.Release(serverID)
		return fmt.Errorf("failed to start listener: %w", err)
	}

	handler := m.createHandler(serverID, matcher, stateManager)
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Store instance
	instance := &ServerInstance{
		ID:       serverID,
		Server:   httpServer,
		Matcher:  matcher,
		State:    stateManager,
		Listener: listener,
	}
	m.servers[serverID] = instance

	// Start server in goroutine
	go func() {
		if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			m.logger.Error("Mock server error", zap.Error(err), zap.String("server_id", serverID.String()))
		}
	}()

	// Update server status
	now := time.Now()
	server.Status = models.MockServerStatusRunning
	server.StartedAt = &now
	if err := m.repo.UpdateServer(server); err != nil {
		return fmt.Errorf("failed to update server status: %w", err)
	}

	m.logger.Info("Mock server started",
		zap.String("server_id", serverID.String()),
		zap.String("name", name),
		zap.Int("port", port),
		zap.Int("endpoints", len(endpoints)),
	)

	return nil
}

// StopServer stops a running mock server
func (m *Manager) StopServer(serverID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance, exists := m.servers[serverID]
	if !exists {
		return fmt.Errorf("server %s not running", serverID)
	}

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := instance.Server.Shutdown(ctx); err != nil {
		m.logger.Warn("Error shutting down mock server", zap.Error(err))
	}

	// Close listener
	if err := instance.Listener.Close(); err != nil {
		m.logger.Warn("Error closing listener", zap.Error(err))
	}

	// Release port
	m.portPool.Release(serverID)

	// Update server record
	server, err := m.repo.GetServerByID(serverID)
	if err != nil {
		return fmt.Errorf("failed to get server record: %w", err)
	}

	now := time.Now()
	server.Status = models.MockServerStatusStopped
	server.StoppedAt = &now
	if err := m.repo.UpdateServer(server); err != nil {
		return fmt.Errorf("failed to update server status: %w", err)
	}

	// Remove from instances
	delete(m.servers, serverID)

	m.logger.Info("Mock server stopped", zap.String("server_id", serverID.String()))

	return nil
}

// AddEndpoint adds a new endpoint to a running mock server
func (m *Manager) AddEndpoint(serverID uuid.UUID, endpoint *models.MockEndpoint) error {
	m.mu.RLock()
	instance, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server %s not running", serverID)
	}

	// Create endpoint in database
	if err := m.repo.CreateEndpoint(endpoint); err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	// Add to matcher
	instance.Matcher.AddEndpoint(endpoint)

	m.logger.Info("Endpoint added to mock server",
		zap.String("server_id", serverID.String()),
		zap.String("method", endpoint.Method),
		zap.String("path", endpoint.Path),
	)

	return nil
}

// GetServer retrieves a server instance
func (m *Manager) GetServer(serverID uuid.UUID) (*ServerInstance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	instance, exists := m.servers[serverID]
	if !exists {
		return nil, fmt.Errorf("server %s not running", serverID)
	}

	return instance, nil
}

// StopAllServers stops all running mock servers
func (m *Manager) StopAllServers() error {
	m.mu.Lock()
	serverIDs := make([]uuid.UUID, 0, len(m.servers))
	for id := range m.servers {
		serverIDs = append(serverIDs, id)
	}
	m.mu.Unlock()

	var lastErr error
	for _, id := range serverIDs {
		if err := m.StopServer(id); err != nil {
			lastErr = err
			m.logger.Error("Failed to stop server", zap.Error(err), zap.String("server_id", id.String()))
		}
	}

	return lastErr
}

// createHandler creates the HTTP handler for a mock server
func (m *Manager) createHandler(serverID uuid.UUID, matcher *EndpointMatcher, stateManager *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log request
		reqBody, _ := io.ReadAll(r.Body)
		r.Body.Close()

		// Convert headers to map
		headers := make(map[string]interface{})
		for k, v := range r.Header {
			if len(v) == 1 {
				headers[k] = v[0]
			} else {
				headers[k] = v
			}
		}

		// Convert query params to map
		queryParams := make(map[string]interface{})
		for k, v := range r.URL.Query() {
			if len(v) == 1 {
				queryParams[k] = v[0]
			} else {
				queryParams[k] = v
			}
		}

		// Match endpoint
		endpoint, matched := matcher.Match(r.Method, r.URL.Path, headers, queryParams, reqBody)

		mockRequest := &models.MockRequest{
			MockServerID: serverID,
			Method:       r.Method,
			Path:         r.URL.Path,
			Headers:      headers,
			QueryParams:  queryParams,
			Body:         string(reqBody),
			Matched:      matched,
			ResponseCode: http.StatusNotFound,
		}

		var response *models.ResponseConfig
		if matched && endpoint != nil {
			mockRequest.EndpointID = &endpoint.ID
			response = &endpoint.ResponseConfig

			// Handle state updates
			if endpoint.StateConfig != nil {
				if err := stateManager.UpdateState(endpoint.StateConfig); err != nil {
					m.logger.Error("Failed to update state", zap.Error(err))
				}
			}
		} else {
			// Default 404 response
			response = &models.ResponseConfig{
				StatusCode: http.StatusNotFound,
				BodyText:   "No matching endpoint found",
			}
		}

		// Apply response delay
		if response.DelayMs > 0 {
			time.Sleep(time.Duration(response.DelayMs) * time.Millisecond)
		}

		// Set response headers
		for k, v := range response.Headers {
			w.Header().Set(k, v)
		}

		// Set status code
		mockRequest.ResponseCode = response.StatusCode
		w.WriteHeader(response.StatusCode)

		// Write response body
		if response.BodyJSON != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response.BodyJSON)
		} else if response.BodyText != "" {
			w.Write([]byte(response.BodyText))
		} else if response.Body != nil {
			if bodyBytes, err := json.Marshal(response.Body); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write(bodyBytes)
			}
		}

		// Log request to database (async)
		go func() {
			if err := m.repo.CreateRequest(mockRequest); err != nil {
				m.logger.Error("Failed to log mock request", zap.Error(err))
			}
		}()

		m.logger.Debug("Mock request handled",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Bool("matched", matched),
			zap.Int("status", response.StatusCode),
		)
	}
}

// PortPool manages port allocation
type PortPool struct {
	min       int
	max       int
	allocated map[uuid.UUID]int
	used      map[int]bool
	mu        sync.Mutex
}

// NewPortPool creates a new port pool
func NewPortPool(min, max int) *PortPool {
	return &PortPool{
		min:       min,
		max:       max,
		allocated: make(map[uuid.UUID]int),
		used:      make(map[int]bool),
	}
}

// Allocate allocates a port for a server
func (p *PortPool) Allocate(serverID uuid.UUID) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if already allocated
	if port, exists := p.allocated[serverID]; exists {
		return port, nil
	}

	// Find available port
	for port := p.min; port <= p.max; port++ {
		if !p.used[port] {
			p.allocated[serverID] = port
			p.used[port] = true
			return port, nil
		}
	}

	return 0, fmt.Errorf("no available ports in range %d-%d", p.min, p.max)
}

// Release releases a port
func (p *PortPool) Release(serverID uuid.UUID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if port, exists := p.allocated[serverID]; exists {
		delete(p.allocated, serverID)
		delete(p.used, port)
	}
}
