// API Types for TestMesh

export interface Flow {
  id: string;
  name: string;
  description: string;
  suite: string;
  tags: string[];
  definition: FlowDefinition;
  environment: string;
  created_at: string;
  updated_at: string;
}

export interface FlowDefinition {
  name: string;
  description: string;
  suite: string;
  tags: string[];
  env?: Record<string, any>;
  setup?: Step[];
  steps: Step[];
  teardown?: Step[];
}

export interface Step {
  id?: string;
  action: string;
  name?: string;
  description?: string;
  config: Record<string, any>;
  assert?: string[];
  output?: Record<string, string>;
  retry?: RetryConfig;
  timeout?: string;
}

export interface RetryConfig {
  max_attempts: number;
  delay: string;
  backoff?: string;
}

export type ExecutionStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';

export interface Execution {
  id: string;
  flow_id: string;
  flow?: Flow;
  status: ExecutionStatus;
  environment: string;
  started_at: string | null;
  finished_at: string | null;
  duration_ms: number;
  total_steps: number;
  passed_steps: number;
  failed_steps: number;
  error?: string;
  created_at: string;
  updated_at: string;
}

export type StepStatus = 'pending' | 'running' | 'completed' | 'failed' | 'skipped';

export interface ExecutionStep {
  id: string;
  execution_id: string;
  step_id: string;
  step_name: string;
  action: string;
  status: StepStatus;
  started_at: string | null;
  finished_at: string | null;
  duration_ms: number;
  output: Record<string, any>;
  error_message?: string;
  attempt: number;
  created_at: string;
  updated_at: string;
}

// API Request/Response Types

export interface CreateFlowRequest {
  yaml: string;
}

export interface UpdateFlowRequest {
  yaml: string;
}

export interface ListFlowsResponse {
  flows: Flow[];
  total: number;
  limit: number;
  offset: number;
}

export interface CreateExecutionRequest {
  flow_id: string;
  environment?: string;
  variables?: Record<string, string>;
}

export interface ListExecutionsResponse {
  executions: Execution[];
  total: number;
  limit: number;
  offset: number;
}

export interface GetStepsResponse {
  steps: ExecutionStep[];
}

export interface GetLogsResponse {
  logs: string[];
}

export interface HealthResponse {
  status: string;
  database: string;
  service: string;
  version: string;
}
