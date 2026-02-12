package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

// PluginType represents the type of plugin
type PluginType string

const (
	PluginTypeAction   PluginType = "action"   // Custom action handlers
	PluginTypeAuth     PluginType = "auth"     // Authentication providers
	PluginTypeExporter PluginType = "exporter" // Export formats
	PluginTypeImporter PluginType = "importer" // Import formats
	PluginTypeReporter PluginType = "reporter" // Report generators
)

// PluginManifest describes a plugin
type PluginManifest struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Homepage    string            `json:"homepage,omitempty"`
	Type        PluginType        `json:"type"`
	EntryPoint  string            `json:"entry_point"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Permissions []string          `json:"permissions,omitempty"`
}

// Plugin represents a loaded plugin
type Plugin struct {
	Manifest *PluginManifest `json:"manifest"`
	Path     string          `json:"path"`
	Enabled  bool            `json:"enabled"`
	Loaded   bool            `json:"loaded"`
	Error    string          `json:"error,omitempty"`
}

// ActionPlugin interface for action handlers
type ActionPlugin interface {
	Name() string
	Execute(ctx context.Context, config map[string]interface{}) (map[string]interface{}, error)
}

// Registry manages plugins
type Registry struct {
	mu        sync.RWMutex
	plugins   map[string]*Plugin
	actions   map[string]ActionPlugin
	pluginDir string
	logger    *zap.Logger
}

// NewRegistry creates a new plugin registry
func NewRegistry(pluginDir string, logger *zap.Logger) *Registry {
	return &Registry{
		plugins:   make(map[string]*Plugin),
		actions:   make(map[string]ActionPlugin),
		pluginDir: pluginDir,
		logger:    logger,
	}
}

// Discover finds all plugins in the plugin directory
func (r *Registry) Discover() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(r.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Scan for plugin directories
	entries, err := os.ReadDir(r.pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(r.pluginDir, entry.Name())
		manifestPath := filepath.Join(pluginPath, "manifest.json")

		// Check if manifest exists
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		// Load manifest
		manifest, err := r.loadManifest(manifestPath)
		if err != nil {
			r.logger.Warn("Failed to load plugin manifest",
				zap.String("path", manifestPath),
				zap.Error(err))
			continue
		}

		plugin := &Plugin{
			Manifest: manifest,
			Path:     pluginPath,
			Enabled:  true,
		}

		r.plugins[manifest.ID] = plugin
		r.logger.Info("Discovered plugin",
			zap.String("id", manifest.ID),
			zap.String("name", manifest.Name),
			zap.String("version", manifest.Version))
	}

	return nil
}

// loadManifest loads a plugin manifest from file
func (r *Registry) loadManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// Load loads a specific plugin
func (r *Registry) Load(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, ok := r.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	if plugin.Loaded {
		return nil // Already loaded
	}

	// Load based on type
	switch plugin.Manifest.Type {
	case PluginTypeAction:
		if err := r.loadActionPlugin(plugin); err != nil {
			plugin.Error = err.Error()
			return err
		}
	default:
		return fmt.Errorf("unsupported plugin type: %s", plugin.Manifest.Type)
	}

	plugin.Loaded = true
	plugin.Error = ""
	return nil
}

// loadActionPlugin loads an action plugin
func (r *Registry) loadActionPlugin(plugin *Plugin) error {
	// For now, this is a placeholder
	// In a full implementation, this would:
	// 1. Load a Go plugin (.so file)
	// 2. Or start a subprocess for non-Go plugins
	// 3. Or connect to a plugin via RPC

	r.logger.Info("Loaded action plugin",
		zap.String("id", plugin.Manifest.ID),
		zap.String("entry_point", plugin.Manifest.EntryPoint))

	return nil
}

// Unload unloads a specific plugin
func (r *Registry) Unload(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, ok := r.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	if !plugin.Loaded {
		return nil // Already unloaded
	}

	// Remove from action registry if it's an action plugin
	if plugin.Manifest.Type == PluginTypeAction {
		delete(r.actions, plugin.Manifest.ID)
	}

	plugin.Loaded = false
	return nil
}

// Enable enables a plugin
func (r *Registry) Enable(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, ok := r.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	plugin.Enabled = true
	return nil
}

// Disable disables a plugin
func (r *Registry) Disable(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, ok := r.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	plugin.Enabled = false

	// Also unload if loaded
	if plugin.Loaded {
		r.mu.Unlock()
		err := r.Unload(id)
		r.mu.Lock()
		return err
	}

	return nil
}

// List returns all plugins
func (r *Registry) List() []*Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]*Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// Get returns a specific plugin
func (r *Registry) Get(id string) (*Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, ok := r.plugins[id]
	if !ok {
		return nil, fmt.Errorf("plugin not found: %s", id)
	}

	return plugin, nil
}

// GetAction returns a registered action plugin
func (r *Registry) GetAction(name string) (ActionPlugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	action, ok := r.actions[name]
	return action, ok
}

// RegisterAction registers a custom action plugin
func (r *Registry) RegisterAction(name string, plugin ActionPlugin) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.actions[name] = plugin
	r.logger.Info("Registered action plugin", zap.String("name", name))
}

// Install installs a plugin from a source (URL, path, etc.)
func (r *Registry) Install(source string) (*Plugin, error) {
	// This would:
	// 1. Download from URL or copy from path
	// 2. Verify signature/checksum
	// 3. Extract to plugin directory
	// 4. Load manifest
	// 5. Run any setup

	return nil, fmt.Errorf("plugin installation not yet implemented")
}

// Uninstall removes a plugin
func (r *Registry) Uninstall(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, ok := r.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	// Unload first
	if plugin.Loaded {
		r.mu.Unlock()
		if err := r.Unload(id); err != nil {
			r.mu.Lock()
			return err
		}
		r.mu.Lock()
	}

	// Remove from registry
	delete(r.plugins, id)

	// Remove plugin directory
	if err := os.RemoveAll(plugin.Path); err != nil {
		return fmt.Errorf("failed to remove plugin directory: %w", err)
	}

	r.logger.Info("Uninstalled plugin", zap.String("id", id))

	return nil
}

// LoadAll loads all enabled plugins
func (r *Registry) LoadAll() error {
	for id, plugin := range r.plugins {
		if plugin.Enabled && !plugin.Loaded {
			if err := r.Load(id); err != nil {
				r.logger.Error("Failed to load plugin",
					zap.String("id", id),
					zap.Error(err))
			}
		}
	}
	return nil
}
