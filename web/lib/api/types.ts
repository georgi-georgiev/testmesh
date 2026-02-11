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

// Mock Server Types
export type MockServerStatus = 'starting' | 'running' | 'stopped' | 'failed';

export interface MockServer {
  id: string;
  execution_id?: string;
  name: string;
  port: number;
  base_url: string;
  status: MockServerStatus;
  started_at?: string;
  stopped_at?: string;
  created_at: string;
  updated_at: string;
}

export interface MockEndpoint {
  id: string;
  mock_server_id: string;
  path: string;
  method: string;
  match_config: MatchConfig;
  response_config: ResponseConfig;
  state_config?: StateConfig;
  priority: number;
  created_at: string;
  updated_at: string;
}

export interface MatchConfig {
  path_pattern?: string;
  headers?: Record<string, string>;
  query_params?: Record<string, string>;
  body_pattern?: string;
  body_json?: Record<string, any>;
}

export interface ResponseConfig {
  status_code: number;
  headers?: Record<string, string>;
  body?: any;
  body_json?: Record<string, any>;
  body_text?: string;
  delay_ms?: number;
  template?: boolean;
  template_vars?: Record<string, any>;
}

export interface StateConfig {
  state_key: string;
  initial_value?: any;
  update_rule?: string;
  update_value?: any;
  condition?: Record<string, any>;
}

export interface MockRequest {
  id: string;
  mock_server_id: string;
  endpoint_id?: string;
  method: string;
  path: string;
  headers: Record<string, any>;
  query_params: Record<string, any>;
  body: string;
  matched: boolean;
  response_code: number;
  received_at: string;
}

export interface MockState {
  id: string;
  mock_server_id: string;
  state_key: string;
  state_value: Record<string, any>;
  updated_at: string;
}

// Contract Testing Types
export type VerificationStatus = 'pending' | 'passed' | 'failed';
export type BreakingChangeSeverity = 'critical' | 'major' | 'minor';

export interface Contract {
  id: string;
  consumer: string;
  provider: string;
  version: string;
  pact_version: string;
  contract_data: ContractData;
  flow_id?: string;
  created_at: string;
  updated_at: string;
}

export interface ContractData {
  consumer: { name: string };
  provider: { name: string };
  interactions: Interaction[];
  metadata: {
    pactSpecification: { version: string };
    client?: { name: string; version: string };
  };
}

export interface Interaction {
  id: string;
  contract_id: string;
  description: string;
  provider_state?: string;
  request: HTTPRequest;
  response: HTTPResponse;
  interaction_type: string;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface HTTPRequest {
  method: string;
  path: string;
  query?: Record<string, any>;
  headers?: Record<string, any>;
  body?: any;
}

export interface HTTPResponse {
  status: number;
  headers?: Record<string, any>;
  body?: any;
}

export interface Verification {
  id: string;
  contract_id: string;
  provider_version: string;
  status: VerificationStatus;
  verified_at: string;
  results: VerificationResults;
  execution_id?: string;
  created_at: string;
  updated_at: string;
}

export interface VerificationResults {
  total_interactions: number;
  passed_interactions: number;
  failed_interactions: number;
  details: InteractionResult[];
  summary: string;
}

export interface InteractionResult {
  interaction_id: string;
  description: string;
  passed: boolean;
  mismatches?: Mismatch[];
  actual_request?: Record<string, any>;
  actual_response?: Record<string, any>;
}

export interface Mismatch {
  type: string;
  expected: any;
  actual: any;
  path?: string;
  message: string;
}

export interface BreakingChange {
  id: string;
  old_contract_id: string;
  new_contract_id: string;
  change_type: string;
  severity: BreakingChangeSeverity;
  description: string;
  details: ChangeDetails;
  detected_at: string;
  created_at: string;
}

export interface ChangeDetails {
  interaction_id?: string;
  field?: string;
  old_value?: any;
  new_value?: any;
  impact: string;
  suggestion?: string;
  metadata?: Record<string, any>;
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

// Mock Server API Responses
export interface ListMockServersResponse {
  servers: MockServer[];
  total: number;
  limit: number;
  offset: number;
}

export interface ListMockEndpointsResponse {
  endpoints: MockEndpoint[];
  total: number;
}

export interface ListMockRequestsResponse {
  requests: MockRequest[];
  total: number;
  limit: number;
  offset: number;
}

export interface ListMockStatesResponse {
  states: MockState[];
  total: number;
}

// Contract Testing API Responses
export interface ListContractsResponse {
  contracts: Contract[];
  total: number;
  limit: number;
  offset: number;
}

export interface GetContractResponse {
  contract: Contract;
  interactions: Interaction[];
}

export interface ListVerificationsResponse {
  verifications: Verification[];
  total: number;
  limit: number;
  offset: number;
}

export interface ListBreakingChangesResponse {
  changes: BreakingChange[];
  total: number;
}

export interface DetectBreakingChangesResponse {
  changes: BreakingChange[];
  summary: {
    total: number;
    critical: number;
    major: number;
    minor: number;
  };
}
