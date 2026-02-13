// Flow Editor Types
import type { Node, Edge } from 'reactflow';
import type { Step, FlowDefinition } from '@/lib/api/types';

// Action types supported by the flow editor
export type ActionType =
  | 'http_request'
  | 'database_query'
  | 'log'
  | 'delay'
  | 'assert'
  | 'transform'
  | 'condition'
  | 'for_each'
  | 'mock_server_start'
  | 'mock_server_stop'
  | 'contract_generate'
  | 'contract_verify';

// Simplified node data structure - uses Record<string, any> for config
// This allows flexibility while maintaining structure for common properties
export interface FlowNodeData {
  label: string;
  stepId: string;
  action: ActionType;
  name?: string;
  description?: string;
  config: Record<string, any>;
  assert?: string[];
  output?: Record<string, string>;
  retry?: {
    max_attempts: number;
    delay: string;
    backoff?: string;
  };
  timeout?: string;
  // UI state
  isSelected?: boolean;
  isRunning?: boolean;
  status?: 'pending' | 'running' | 'completed' | 'failed' | 'skipped';
}

// Condition node data - extends FlowNodeData with branching info
export interface ConditionNodeData extends FlowNodeData {
  action: 'condition';
  // IDs of nodes in the then branch
  thenBranchIds?: string[];
  // IDs of nodes in the else branch
  elseBranchIds?: string[];
}

// ForEach node data - extends FlowNodeData with loop info
export interface ForEachNodeData extends FlowNodeData {
  action: 'for_each';
  // IDs of nested step nodes
  nestedStepIds?: string[];
  // Count of nested steps (for display when collapsed)
  nestedStepCount?: number;
  // Whether the container is expanded
  isExpanded?: boolean;
}

// Section header node data (for visual grouping)
export interface SectionHeaderData {
  label: string;
  section: FlowSection;
}

// Custom node type for React Flow - supports flow nodes, section headers, and special nodes
export type FlowNode = Node<FlowNodeData | SectionHeaderData | ConditionNodeData | ForEachNodeData>;
export type FlowEdge = Edge;

// Labeled edge for condition branches
export interface LabeledEdgeData {
  label?: string;
  labelStyle?: React.CSSProperties;
}

// Editor state
export interface FlowEditorState {
  nodes: FlowNode[];
  edges: FlowEdge[];
  selectedNodeId: string | null;
  isDirty: boolean;
}

// Palette item for draggable nodes
export interface PaletteItem {
  type: ActionType;
  label: string;
  description: string;
  icon: string;
  category: 'http' | 'database' | 'control' | 'mock' | 'contract' | 'utility';
  defaultConfig: Record<string, any>;
}

// Section types for flow phases
export type FlowSection = 'setup' | 'steps' | 'teardown';

// Conversion utilities types
export interface ConversionOptions {
  generateIds?: boolean;
  includeComments?: boolean;
}

// Type guard to check if node data is FlowNodeData (not section header)
export function isFlowNodeData(data: FlowNodeData | SectionHeaderData): data is FlowNodeData {
  return 'action' in data && 'stepId' in data;
}

// Type guard to check if node data is SectionHeaderData
export function isSectionHeaderData(data: FlowNodeData | SectionHeaderData): data is SectionHeaderData {
  return 'section' in data && !('action' in data);
}

// Type guard to check if node data is ConditionNodeData
export function isConditionNodeData(data: FlowNodeData | SectionHeaderData | ConditionNodeData | ForEachNodeData): data is ConditionNodeData {
  return 'action' in data && data.action === 'condition';
}

// Type guard to check if node data is ForEachNodeData
export function isForEachNodeData(data: FlowNodeData | SectionHeaderData | ConditionNodeData | ForEachNodeData): data is ForEachNodeData {
  return 'action' in data && data.action === 'for_each';
}
