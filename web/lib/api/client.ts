import axios, { AxiosInstance, AxiosError } from 'axios';
import type {
  Flow,
  Execution,
  CreateFlowRequest,
  UpdateFlowRequest,
  ListFlowsResponse,
  CreateExecutionRequest,
  ListExecutionsResponse,
  GetStepsResponse,
  GetLogsResponse,
  HealthResponse,
  ExecutionStatus,
} from './types';

// API Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Create axios instance with default configuration
const createAxiosInstance = (): AxiosInstance => {
  const instance = axios.create({
    baseURL: API_BASE_URL,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Request interceptor
  instance.interceptors.request.use(
    (config) => {
      // Add any auth tokens here if needed
      // const token = localStorage.getItem('token');
      // if (token) {
      //   config.headers.Authorization = `Bearer ${token}`;
      // }
      return config;
    },
    (error) => Promise.reject(error)
  );

  // Response interceptor
  instance.interceptors.response.use(
    (response) => response,
    (error: AxiosError) => {
      // Handle common errors
      if (error.response?.status === 401) {
        // Handle unauthorized - redirect to login, etc.
      }
      return Promise.reject(error);
    }
  );

  return instance;
};

// Export for use by other API modules
export const apiClient = createAxiosInstance();

// Health API
export const healthApi = {
  check: async (): Promise<HealthResponse> => {
    const response = await apiClient.get<HealthResponse>('/health');
    return response.data;
  },
};

// Flow API
export const flowApi = {
  create: async (data: CreateFlowRequest): Promise<Flow> => {
    const response = await apiClient.post<Flow>('/api/v1/flows', data);
    return response.data;
  },

  list: async (params?: {
    suite?: string;
    tags?: string[];
    limit?: number;
    offset?: number;
  }): Promise<ListFlowsResponse> => {
    const response = await apiClient.get<ListFlowsResponse>('/api/v1/flows', {
      params,
    });
    return response.data;
  },

  get: async (id: string): Promise<Flow> => {
    const response = await apiClient.get<Flow>(`/api/v1/flows/${id}`);
    return response.data;
  },

  update: async (id: string, data: UpdateFlowRequest): Promise<Flow> => {
    const response = await apiClient.put<Flow>(`/api/v1/flows/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/api/v1/flows/${id}`);
  },
};

// Execution API
export const executionApi = {
  create: async (data: CreateExecutionRequest): Promise<Execution> => {
    const response = await apiClient.post<Execution>('/api/v1/executions', data);
    return response.data;
  },

  list: async (params?: {
    flow_id?: string;
    status?: ExecutionStatus;
    limit?: number;
    offset?: number;
  }): Promise<ListExecutionsResponse> => {
    const response = await apiClient.get<ListExecutionsResponse>('/api/v1/executions', {
      params,
    });
    return response.data;
  },

  get: async (id: string): Promise<Execution> => {
    const response = await apiClient.get<Execution>(`/api/v1/executions/${id}`);
    return response.data;
  },

  cancel: async (id: string): Promise<Execution> => {
    const response = await apiClient.post<Execution>(`/api/v1/executions/${id}/cancel`);
    return response.data;
  },

  getLogs: async (id: string): Promise<GetLogsResponse> => {
    const response = await apiClient.get<GetLogsResponse>(`/api/v1/executions/${id}/logs`);
    return response.data;
  },

  getSteps: async (id: string): Promise<GetStepsResponse> => {
    const response = await apiClient.get<GetStepsResponse>(`/api/v1/executions/${id}/steps`);
    return response.data;
  },

  getStep: async (executionId: string, stepId: string): Promise<import('./types').ExecutionStep> => {
    const response = await apiClient.get<import('./types').ExecutionStep>(
      `/api/v1/executions/${executionId}/steps/${stepId}`
    );
    return response.data;
  },
};

// Export a default API object with all endpoints
const api = {
  health: healthApi,
  flows: flowApi,
  executions: executionApi,
};

export default api;
