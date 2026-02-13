package repository

import (
	"fmt"

	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnvironmentRepository handles environment database operations
type EnvironmentRepository struct {
	db *gorm.DB
}

// NewEnvironmentRepository creates a new environment repository
func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

// Create creates a new environment
func (r *EnvironmentRepository) Create(env *models.Environment) error {
	// If this is set as default, unset other defaults
	if env.IsDefault {
		if err := r.db.Model(&models.Environment{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
	}
	return r.db.Create(env).Error
}

// GetByID retrieves an environment by ID
func (r *EnvironmentRepository) GetByID(id uuid.UUID) (*models.Environment, error) {
	var env models.Environment
	if err := r.db.First(&env, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

// GetByName retrieves an environment by name
func (r *EnvironmentRepository) GetByName(name string) (*models.Environment, error) {
	var env models.Environment
	if err := r.db.First(&env, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

// GetDefault retrieves the default environment
func (r *EnvironmentRepository) GetDefault() (*models.Environment, error) {
	var env models.Environment
	if err := r.db.First(&env, "is_default = ?", true).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

// Update updates an environment
func (r *EnvironmentRepository) Update(env *models.Environment) error {
	// If setting as default, unset other defaults
	if env.IsDefault {
		if err := r.db.Model(&models.Environment{}).
			Where("is_default = ? AND id != ?", true, env.ID).
			Update("is_default", false).Error; err != nil {
			return err
		}
	}
	return r.db.Save(env).Error
}

// Delete soft-deletes an environment
func (r *EnvironmentRepository) Delete(id uuid.UUID) error {
	// Check if this is the default environment
	env, err := r.GetByID(id)
	if err != nil {
		return err
	}
	if env.IsDefault {
		return fmt.Errorf("cannot delete the default environment")
	}
	return r.db.Delete(&models.Environment{}, "id = ?", id).Error
}

// List retrieves all environments with optional filtering
func (r *EnvironmentRepository) List(params *ListEnvironmentsParams) ([]*models.Environment, int64, error) {
	var environments []*models.Environment
	var total int64

	query := r.db.Model(&models.Environment{})

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	// Count total
	query.Count(&total)

	// Apply sorting
	if params.SortBy != "" {
		order := "ASC"
		if params.SortDesc {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", params.SortBy, order))
	} else {
		// Default: show default env first, then by name
		query = query.Order("is_default DESC, name ASC")
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	if err := query.Find(&environments).Error; err != nil {
		return nil, 0, err
	}

	return environments, total, nil
}

// ListEnvironmentsParams contains parameters for listing environments
type ListEnvironmentsParams struct {
	Search   string
	SortBy   string
	SortDesc bool
	Limit    int
	Offset   int
}

// SetDefault sets an environment as the default
func (r *EnvironmentRepository) SetDefault(id uuid.UUID) error {
	// Unset current default
	if err := r.db.Model(&models.Environment{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
		return err
	}
	// Set new default
	return r.db.Model(&models.Environment{}).Where("id = ?", id).Update("is_default", true).Error
}

// Duplicate creates a copy of an environment with a new name
func (r *EnvironmentRepository) Duplicate(id uuid.UUID, newName string) (*models.Environment, error) {
	original, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Copy variables
	varsCopy := make(models.EnvironmentVariables, len(original.Variables))
	copy(varsCopy, original.Variables)

	newEnv := &models.Environment{
		Name:        newName,
		Description: original.Description + " (copy)",
		Color:       original.Color,
		IsDefault:   false,
		Variables:   varsCopy,
	}

	if err := r.Create(newEnv); err != nil {
		return nil, err
	}

	return newEnv, nil
}

// EnsureDefaultExists creates a default environment if none exists
func (r *EnvironmentRepository) EnsureDefaultExists() error {
	var count int64
	r.db.Model(&models.Environment{}).Count(&count)

	if count == 0 {
		defaultEnv := &models.Environment{
			Name:        "Default",
			Description: "Default environment",
			Color:       "#3B82F6",
			IsDefault:   true,
			Variables:   models.EnvironmentVariables{},
		}
		return r.Create(defaultEnv)
	}
	return nil
}
