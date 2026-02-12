import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { aiApi } from '../api/ai';
import type {
  GenerateFlowRequest,
  ImportOpenAPIRequest,
  AnalyzeCoverageRequest,
  SuggestionStatus,
  AIProviderType,
} from '../api/types';

// Query keys
export const aiKeys = {
  all: ['ai'] as const,
  providers: () => [...aiKeys.all, 'providers'] as const,
  usage: () => [...aiKeys.all, 'usage'] as const,
  suggestions: () => [...aiKeys.all, 'suggestions'] as const,
  suggestionsList: (flowId: string, status?: SuggestionStatus) =>
    [...aiKeys.suggestions(), 'list', flowId, status] as const,
  suggestionDetail: (id: string) => [...aiKeys.suggestions(), 'detail', id] as const,
};

// Provider hooks
export function useAIProviders() {
  return useQuery({
    queryKey: aiKeys.providers(),
    queryFn: () => aiApi.getProviders(),
  });
}

export function useAIUsage() {
  return useQuery({
    queryKey: aiKeys.usage(),
    queryFn: () => aiApi.getUsage(),
  });
}

// Generation hooks
export function useGenerateFlow() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: GenerateFlowRequest) => aiApi.generate(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: aiKeys.usage() });
    },
  });
}

// Import hooks
export function useImportOpenAPI() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ImportOpenAPIRequest) => aiApi.importOpenAPI(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['flows'] });
      queryClient.invalidateQueries({ queryKey: aiKeys.usage() });
    },
  });
}

export function useImportPostman() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: { collection: string; provider?: AIProviderType; model?: string; create_flows?: boolean }) =>
      aiApi.importPostman(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['flows'] });
      queryClient.invalidateQueries({ queryKey: aiKeys.usage() });
    },
  });
}

export function useImportPact() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: { contract: string; provider?: AIProviderType; model?: string; create_flows?: boolean }) =>
      aiApi.importPact(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['flows'] });
      queryClient.invalidateQueries({ queryKey: aiKeys.usage() });
    },
  });
}

// Coverage hooks
export function useAnalyzeCoverage() {
  return useMutation({
    mutationFn: (data: AnalyzeCoverageRequest) => aiApi.analyzeCoverage(data),
  });
}

// Self-healing hooks
export function useAnalyzeFailure() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (executionId: string) => aiApi.analyzeFailure(executionId),
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: aiKeys.suggestionsList(data.flow_id),
      });
    },
  });
}

export function useSuggestions(flowId: string, status?: SuggestionStatus) {
  return useQuery({
    queryKey: aiKeys.suggestionsList(flowId, status),
    queryFn: () => aiApi.listSuggestions({ flow_id: flowId, status }),
    enabled: !!flowId,
  });
}

export function useSuggestion(id: string) {
  return useQuery({
    queryKey: aiKeys.suggestionDetail(id),
    queryFn: () => aiApi.getSuggestion(id),
    enabled: !!id,
  });
}

export function useApplySuggestion() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => aiApi.applySuggestion(id),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: aiKeys.suggestions() });
      queryClient.invalidateQueries({ queryKey: ['flows', 'detail', data.flow_id] });
    },
  });
}

export function useAcceptSuggestion() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => aiApi.acceptSuggestion(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: aiKeys.suggestions() });
    },
  });
}

export function useRejectSuggestion() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => aiApi.rejectSuggestion(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: aiKeys.suggestions() });
    },
  });
}
