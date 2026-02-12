'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Plus, FolderTree, Search, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { CollectionTree, CollectionDialog } from '@/components/collections';
import {
  useCollectionTree,
  useCreateCollection,
  useUpdateCollection,
  useDeleteCollection,
  useDuplicateCollection,
  useCollection,
} from '@/lib/hooks/useCollections';
import type { CollectionTreeNode, CreateCollectionRequest, UpdateCollectionRequest } from '@/lib/api/types';

export default function CollectionsPage() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedNodeId, setSelectedNodeId] = useState<string>();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingCollectionId, setEditingCollectionId] = useState<string | null>(null);
  const [parentIdForCreate, setParentIdForCreate] = useState<string | undefined>();

  // Queries
  const { data: treeData, isLoading } = useCollectionTree();
  const { data: editingCollection } = useCollection(editingCollectionId || '');

  // Mutations
  const createCollection = useCreateCollection();
  const updateCollection = useUpdateCollection();
  const deleteCollection = useDeleteCollection();
  const duplicateCollection = useDuplicateCollection();

  const handleSelect = (node: CollectionTreeNode) => {
    setSelectedNodeId(node.id);
  };

  const handleCreateCollection = (parentId?: string) => {
    setEditingCollectionId(null);
    setParentIdForCreate(parentId);
    setDialogOpen(true);
  };

  const handleEditCollection = (id: string) => {
    setEditingCollectionId(id);
    setParentIdForCreate(undefined);
    setDialogOpen(true);
  };

  const handleDeleteCollection = async (id: string) => {
    if (confirm('Are you sure you want to delete this collection?')) {
      await deleteCollection.mutateAsync(id);
    }
  };

  const handleDuplicateCollection = async (id: string) => {
    const name = prompt('Enter name for the duplicate:');
    if (name) {
      await duplicateCollection.mutateAsync({ id, name });
    }
  };

  const handleMoveCollection = (id: string) => {
    // TODO: Implement move dialog
    alert('Move functionality coming soon!');
  };

  const handleRunFlow = (flowId: string) => {
    router.push(`/flows/${flowId}/run`);
  };

  const handleEditFlow = (flowId: string) => {
    router.push(`/flows/${flowId}/edit`);
  };

  const handleDeleteFlow = async (flowId: string) => {
    // This removes the flow from the collection, not deletes it
    // TODO: Implement remove from collection
    alert('Remove from collection functionality coming soon!');
  };

  const handleDialogSubmit = async (data: CreateCollectionRequest | UpdateCollectionRequest) => {
    if (editingCollectionId) {
      await updateCollection.mutateAsync({ id: editingCollectionId, data });
    } else {
      await createCollection.mutateAsync(data as CreateCollectionRequest);
    }
  };

  // Filter tree by search query
  const filteredTree = treeData?.tree || [];
  // TODO: Implement proper tree filtering

  return (
    <div className="flex h-[calc(100vh-4rem)]">
      {/* Sidebar with tree */}
      <div className="w-72 border-r flex flex-col bg-muted/30">
        <div className="p-3 border-b space-y-2">
          <div className="flex items-center gap-2">
            <FolderTree className="w-5 h-5 text-primary" />
            <h2 className="font-semibold">Collections</h2>
          </div>
          <div className="relative">
            <Search className="absolute left-2 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search collections..."
              className="pl-8 h-8 text-sm"
            />
          </div>
        </div>

        {isLoading ? (
          <div className="flex-1 flex items-center justify-center">
            <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
          </div>
        ) : (
          <CollectionTree
            tree={filteredTree}
            selectedId={selectedNodeId}
            onSelect={handleSelect}
            onCreateCollection={handleCreateCollection}
            onEditCollection={handleEditCollection}
            onDeleteCollection={handleDeleteCollection}
            onDuplicateCollection={handleDuplicateCollection}
            onMoveCollection={handleMoveCollection}
            onRunFlow={handleRunFlow}
            onEditFlow={handleEditFlow}
            onDeleteFlow={handleDeleteFlow}
            className="flex-1"
          />
        )}
      </div>

      {/* Main content area */}
      <div className="flex-1 p-6 overflow-auto">
        {selectedNodeId ? (
          <SelectedNodeDetails
            nodeId={selectedNodeId}
            tree={filteredTree}
            onEdit={handleEditCollection}
            onDelete={handleDeleteCollection}
            onAddFlow={() => {
              // TODO: Implement add flow dialog
              alert('Add flow functionality coming soon!');
            }}
          />
        ) : (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <FolderTree className="w-16 h-16 text-muted-foreground/30 mb-4" />
            <h3 className="text-lg font-medium mb-2">Organize Your Flows</h3>
            <p className="text-muted-foreground mb-6 max-w-md">
              Collections help you organize your test flows into logical groups.
              Create nested folders, set collection-level variables, and manage authentication
              settings that apply to all flows within.
            </p>
            <Button onClick={() => handleCreateCollection()}>
              <Plus className="w-4 h-4 mr-2" />
              Create Collection
            </Button>
          </div>
        )}
      </div>

      {/* Collection dialog */}
      <CollectionDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        collection={editingCollection || null}
        parentId={parentIdForCreate}
        onSubmit={handleDialogSubmit}
        isLoading={createCollection.isPending || updateCollection.isPending}
      />
    </div>
  );
}

// Helper component to display selected node details
function SelectedNodeDetails({
  nodeId,
  tree,
  onEdit,
  onDelete,
  onAddFlow,
}: {
  nodeId: string;
  tree: CollectionTreeNode[];
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
  onAddFlow: () => void;
}) {
  // Find node in tree
  const findNode = (nodes: CollectionTreeNode[], id: string): CollectionTreeNode | null => {
    for (const node of nodes) {
      if (node.id === id) return node;
      if (node.children) {
        const found = findNode(node.children, id);
        if (found) return found;
      }
    }
    return null;
  };

  const node = findNode(tree, nodeId);

  if (!node) {
    return (
      <div className="text-center text-muted-foreground">
        Collection not found
      </div>
    );
  }

  if (node.type === 'flow') {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            {node.name}
          </CardTitle>
          <CardDescription>
            Test Flow
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Click "Edit Flow" to view and modify this flow.
          </p>
        </CardContent>
      </Card>
    );
  }

  const childCollections = node.children?.filter((c) => c.type === 'collection') || [];
  const childFlows = node.children?.filter((c) => c.type === 'flow') || [];

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              {node.icon && <span className="text-2xl">{node.icon}</span>}
              <div>
                <CardTitle>{node.name}</CardTitle>
                {node.description && (
                  <CardDescription>{node.description}</CardDescription>
                )}
              </div>
            </div>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" onClick={() => onEdit(node.id)}>
                Edit
              </Button>
              <Button variant="outline" size="sm" onClick={onAddFlow}>
                Add Flow
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-muted-foreground">Sub-collections:</span>
              <span className="ml-2 font-medium">{childCollections.length}</span>
            </div>
            <div>
              <span className="text-muted-foreground">Flows:</span>
              <span className="ml-2 font-medium">{childFlows.length}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {childFlows.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Flows in this Collection</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {childFlows.map((flow) => (
                <div
                  key={flow.id}
                  className="flex items-center justify-between p-2 rounded-md hover:bg-muted/50"
                >
                  <span className="text-sm">{flow.name}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
