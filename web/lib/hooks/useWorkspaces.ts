import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  listWorkspaces,
  getWorkspace,
  getWorkspaceBySlug,
  createWorkspace,
  updateWorkspace,
  deleteWorkspace,
  getPersonalWorkspace,
  getUserRole,
  listMembers,
  addMember,
  updateMember,
  removeMember,
  listInvitations,
  inviteMember,
  revokeInvitation,
  acceptInvitation,
  type Workspace,
  type WorkspaceType,
  type CreateWorkspaceRequest,
  type UpdateWorkspaceRequest,
  type AddMemberRequest,
  type InviteMemberRequest,
  type WorkspaceRole,
  type AcceptInvitationRequest,
} from '@/lib/api/workspaces';

// Query keys
export const workspaceKeys = {
  all: ['workspaces'] as const,
  lists: () => [...workspaceKeys.all, 'list'] as const,
  list: (filters: Record<string, any>) => [...workspaceKeys.lists(), filters] as const,
  details: () => [...workspaceKeys.all, 'detail'] as const,
  detail: (id: string) => [...workspaceKeys.details(), id] as const,
  bySlug: (slug: string) => [...workspaceKeys.all, 'slug', slug] as const,
  personal: () => [...workspaceKeys.all, 'personal'] as const,
  role: (workspaceId: string) => [...workspaceKeys.all, workspaceId, 'role'] as const,
  members: (workspaceId: string) => [...workspaceKeys.all, workspaceId, 'members'] as const,
  invitations: (workspaceId: string) => [...workspaceKeys.all, workspaceId, 'invitations'] as const,
};

// List workspaces
export function useWorkspaces(params?: {
  type?: WorkspaceType;
  search?: string;
  sort_by?: string;
  sort_desc?: boolean;
  limit?: number;
  offset?: number;
}) {
  return useQuery({
    queryKey: workspaceKeys.list(params ?? {}),
    queryFn: () => listWorkspaces(params),
  });
}

// Get single workspace
export function useWorkspace(id: string | undefined) {
  return useQuery({
    queryKey: workspaceKeys.detail(id!),
    queryFn: () => getWorkspace(id!),
    enabled: !!id,
  });
}

// Get workspace by slug
export function useWorkspaceBySlug(slug: string | undefined) {
  return useQuery({
    queryKey: workspaceKeys.bySlug(slug!),
    queryFn: () => getWorkspaceBySlug(slug!),
    enabled: !!slug,
  });
}

// Get personal workspace
export function usePersonalWorkspace(userName?: string) {
  return useQuery({
    queryKey: workspaceKeys.personal(),
    queryFn: () => getPersonalWorkspace(userName),
  });
}

// Get user role in workspace
export function useUserRole(workspaceId: string | undefined) {
  return useQuery({
    queryKey: workspaceKeys.role(workspaceId!),
    queryFn: () => getUserRole(workspaceId!),
    enabled: !!workspaceId,
  });
}

// Get workspace members
export function useWorkspaceMembers(workspaceId: string | undefined) {
  return useQuery({
    queryKey: workspaceKeys.members(workspaceId!),
    queryFn: () => listMembers(workspaceId!),
    enabled: !!workspaceId,
  });
}

// Get workspace invitations
export function useWorkspaceInvitations(workspaceId: string | undefined) {
  return useQuery({
    queryKey: workspaceKeys.invitations(workspaceId!),
    queryFn: () => listInvitations(workspaceId!),
    enabled: !!workspaceId,
  });
}

// Create workspace mutation
export function useCreateWorkspace() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateWorkspaceRequest) => createWorkspace(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.lists() });
    },
  });
}

// Update workspace mutation
export function useUpdateWorkspace() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateWorkspaceRequest }) =>
      updateWorkspace(id, data),
    onSuccess: (updatedWorkspace) => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.lists() });
      queryClient.setQueryData(workspaceKeys.detail(updatedWorkspace.id), updatedWorkspace);
    },
  });
}

// Delete workspace mutation
export function useDeleteWorkspace() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteWorkspace(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.lists() });
    },
  });
}

// Add member mutation
export function useAddMember(workspaceId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AddMemberRequest) => addMember(workspaceId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.members(workspaceId) });
      queryClient.invalidateQueries({ queryKey: workspaceKeys.detail(workspaceId) });
    },
  });
}

// Update member mutation
export function useUpdateMember(workspaceId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: WorkspaceRole }) =>
      updateMember(workspaceId, userId, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.members(workspaceId) });
    },
  });
}

// Remove member mutation
export function useRemoveMember(workspaceId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userId: string) => removeMember(workspaceId, userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.members(workspaceId) });
      queryClient.invalidateQueries({ queryKey: workspaceKeys.detail(workspaceId) });
    },
  });
}

// Invite member mutation
export function useInviteMember(workspaceId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: InviteMemberRequest) => inviteMember(workspaceId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.invitations(workspaceId) });
    },
  });
}

// Revoke invitation mutation
export function useRevokeInvitation(workspaceId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (invitationId: string) => revokeInvitation(workspaceId, invitationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.invitations(workspaceId) });
    },
  });
}

// Accept invitation mutation
export function useAcceptInvitation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AcceptInvitationRequest) => acceptInvitation(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.lists() });
    },
  });
}
