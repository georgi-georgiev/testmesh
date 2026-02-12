package repository

import (
	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CollectionRepository handles collection database operations
type CollectionRepository struct {
	db *gorm.DB
}

// NewCollectionRepository creates a new collection repository
func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

// Create creates a new collection
func (r *CollectionRepository) Create(collection *models.Collection) error {
	return r.db.Create(collection).Error
}

// GetByID retrieves a collection by ID
func (r *CollectionRepository) GetByID(id uuid.UUID) (*models.Collection, error) {
	var collection models.Collection
	if err := r.db.First(&collection, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

// GetByIDWithFlows retrieves a collection with its flows
func (r *CollectionRepository) GetByIDWithFlows(id uuid.UUID) (*models.Collection, error) {
	var collection models.Collection
	if err := r.db.Preload("Flows", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&collection, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

// GetByIDWithChildren retrieves a collection with its children
func (r *CollectionRepository) GetByIDWithChildren(id uuid.UUID) (*models.Collection, error) {
	var collection models.Collection
	if err := r.db.Preload("Children", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&collection, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

// List retrieves all root collections (no parent)
func (r *CollectionRepository) List(limit, offset int) ([]models.Collection, int64, error) {
	var collections []models.Collection
	var total int64

	query := r.db.Model(&models.Collection{}).Where("parent_id IS NULL")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Limit(limit).Offset(offset).Order("sort_order ASC, created_at DESC").Find(&collections).Error; err != nil {
		return nil, 0, err
	}

	return collections, total, nil
}

// ListAll retrieves all collections (including nested)
func (r *CollectionRepository) ListAll() ([]models.Collection, error) {
	var collections []models.Collection
	if err := r.db.Order("sort_order ASC, created_at DESC").Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}

// ListChildren retrieves children of a collection
func (r *CollectionRepository) ListChildren(parentID uuid.UUID) ([]models.Collection, error) {
	var collections []models.Collection
	if err := r.db.Where("parent_id = ?", parentID).Order("sort_order ASC").Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}

// Update updates a collection
func (r *CollectionRepository) Update(collection *models.Collection) error {
	return r.db.Save(collection).Error
}

// Delete deletes a collection (soft delete)
func (r *CollectionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Collection{}, "id = ?", id).Error
}

// Move moves a collection to a new parent
func (r *CollectionRepository) Move(id uuid.UUID, newParentID *uuid.UUID) error {
	return r.db.Model(&models.Collection{}).Where("id = ?", id).Update("parent_id", newParentID).Error
}

// Reorder updates the sort order of a collection
func (r *CollectionRepository) Reorder(id uuid.UUID, sortOrder int) error {
	return r.db.Model(&models.Collection{}).Where("id = ?", id).Update("sort_order", sortOrder).Error
}

// AddFlow adds a flow to a collection
func (r *CollectionRepository) AddFlow(collectionID, flowID uuid.UUID) error {
	return r.db.Model(&models.Flow{}).Where("id = ?", flowID).Update("collection_id", collectionID).Error
}

// RemoveFlow removes a flow from a collection
func (r *CollectionRepository) RemoveFlow(flowID uuid.UUID) error {
	return r.db.Model(&models.Flow{}).Where("id = ?", flowID).Update("collection_id", nil).Error
}

// GetFlows retrieves flows in a collection
func (r *CollectionRepository) GetFlows(collectionID uuid.UUID) ([]models.Flow, error) {
	var flows []models.Flow
	if err := r.db.Where("collection_id = ?", collectionID).Order("sort_order ASC").Find(&flows).Error; err != nil {
		return nil, err
	}
	return flows, nil
}

// ReorderFlow updates the sort order of a flow within a collection
func (r *CollectionRepository) ReorderFlow(flowID uuid.UUID, sortOrder int) error {
	return r.db.Model(&models.Flow{}).Where("id = ?", flowID).Update("sort_order", sortOrder).Error
}

// GetTree builds the full collection tree
func (r *CollectionRepository) GetTree() ([]models.CollectionTreeNode, error) {
	// Get all collections
	var collections []models.Collection
	if err := r.db.Order("sort_order ASC").Find(&collections).Error; err != nil {
		return nil, err
	}

	// Get all flows with collection assignments
	var flows []models.Flow
	if err := r.db.Where("collection_id IS NOT NULL").Order("sort_order ASC").Find(&flows).Error; err != nil {
		return nil, err
	}

	// Build tree
	return buildTree(collections, flows), nil
}

// buildTree constructs the collection tree recursively
func buildTree(collections []models.Collection, flows []models.Flow) []models.CollectionTreeNode {
	// Create map of collections by ID
	collectionMap := make(map[uuid.UUID]models.Collection)
	for _, c := range collections {
		collectionMap[c.ID] = c
	}

	// Create map of flows by collection ID
	flowMap := make(map[uuid.UUID][]models.Flow)
	for _, f := range flows {
		if f.CollectionID != nil {
			flowMap[*f.CollectionID] = append(flowMap[*f.CollectionID], f)
		}
	}

	// Build tree starting from root collections
	var rootNodes []models.CollectionTreeNode
	for _, c := range collections {
		if c.ParentID == nil {
			node := buildCollectionNode(c, collections, flowMap)
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes
}

// buildCollectionNode builds a tree node for a collection
func buildCollectionNode(collection models.Collection, allCollections []models.Collection, flowMap map[uuid.UUID][]models.Flow) models.CollectionTreeNode {
	node := models.CollectionTreeNode{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		Icon:        collection.Icon,
		Color:       collection.Color,
		Type:        "collection",
		SortOrder:   collection.SortOrder,
	}

	// Add child collections
	for _, c := range allCollections {
		if c.ParentID != nil && *c.ParentID == collection.ID {
			childNode := buildCollectionNode(c, allCollections, flowMap)
			node.Children = append(node.Children, childNode)
		}
	}

	// Add flows
	if flows, ok := flowMap[collection.ID]; ok {
		for _, f := range flows {
			flowID := f.ID
			flowNode := models.CollectionTreeNode{
				ID:        f.ID,
				Name:      f.Name,
				Type:      "flow",
				SortOrder: f.SortOrder,
				FlowID:    &flowID,
			}
			node.Children = append(node.Children, flowNode)
		}
	}

	return node
}

// Search searches collections by name
func (r *CollectionRepository) Search(query string, limit int) ([]models.Collection, error) {
	var collections []models.Collection
	if err := r.db.Where("name ILIKE ?", "%"+query+"%").Limit(limit).Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}

// GetAncestors retrieves all ancestors of a collection (for breadcrumb)
func (r *CollectionRepository) GetAncestors(id uuid.UUID) ([]models.Collection, error) {
	var ancestors []models.Collection

	current, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	for current.ParentID != nil {
		parent, err := r.GetByID(*current.ParentID)
		if err != nil {
			break
		}
		ancestors = append([]models.Collection{*parent}, ancestors...)
		current = parent
	}

	return ancestors, nil
}

// Duplicate duplicates a collection with all its contents
func (r *CollectionRepository) Duplicate(id uuid.UUID, newName string) (*models.Collection, error) {
	original, err := r.GetByIDWithFlows(id)
	if err != nil {
		return nil, err
	}

	duplicate := &models.Collection{
		Name:        newName,
		Description: original.Description,
		Icon:        original.Icon,
		Color:       original.Color,
		ParentID:    original.ParentID,
		Variables:   original.Variables,
		Auth:        original.Auth,
	}

	if err := r.Create(duplicate); err != nil {
		return nil, err
	}

	// Note: Flows are not duplicated, just the collection structure
	// Flows would need to be duplicated separately if needed

	return duplicate, nil
}
