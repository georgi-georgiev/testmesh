// Flow Editor Utilities - YAML ↔ Visual Conversion
import type { Node, Edge } from 'reactflow';
import type { Step, FlowDefinition } from '@/lib/api/types';
import type {
  FlowNode,
  FlowEdge,
  FlowNodeData,
  SectionHeaderData,
  ActionType,
  FlowSection,
  ConversionOptions,
} from './types';
import { isFlowNodeData, isSectionHeaderData } from './types';

// Generate a unique ID for nodes
export function generateNodeId(): string {
  return `node_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

// Generate a unique step ID
export function generateStepId(): string {
  return `step_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

// Default configurations for each action type
export const defaultConfigs: Record<ActionType, Record<string, any>> = {
  http_request: {
    method: 'GET',
    url: '',
    headers: {},
    body: null,
  },
  database_query: {
    connection: '',
    query: '',
    params: [],
  },
  log: {
    message: '',
    level: 'info',
  },
  delay: {
    duration: '1s',
  },
  assert: {
    expression: '',
    message: '',
  },
  transform: {
    input: '',
    expression: '',
    output_var: '',
  },
  condition: {
    expression: '',
    then_steps: [],
    else_steps: [],
  },
  for_each: {
    items: '',
    item_var: 'item',
    steps: [],
  },
  mock_server_start: {
    name: '',
    port: 8080,
    endpoints: [],
  },
  mock_server_stop: {
    name: '',
  },
  contract_generate: {
    consumer: '',
    provider: '',
    interactions: [],
  },
  contract_verify: {
    contract_id: '',
    provider_base_url: '',
  },
};

// Node layout constants
const NODE_WIDTH = 280;
const NODE_HEIGHT = 80;
const NODE_VERTICAL_GAP = 100;
const SECTION_GAP = 150;

// Convert a Step to a FlowNode
export function stepToNode(
  step: Step,
  position: { x: number; y: number },
  section: FlowSection
): FlowNode {
  const nodeId = step.id || generateNodeId();

  const data: FlowNodeData = {
    label: step.name || step.id || step.action,
    stepId: step.id || nodeId,
    action: step.action as ActionType,
    name: step.name,
    description: step.description,
    config: step.config || defaultConfigs[step.action as ActionType] || {},
    assert: step.assert,
    output: step.output,
    retry: step.retry,
    timeout: step.timeout,
  };

  return {
    id: nodeId,
    type: 'flowNode',
    position,
    data,
  };
}

// Convert a FlowNode back to a Step
export function nodeToStep(node: FlowNode): Step {
  const data = node.data as FlowNodeData;

  const step: Step = {
    id: data.stepId,
    action: data.action,
    config: data.config,
  };

  if (data.name) step.name = data.name;
  if (data.description) step.description = data.description;
  if (data.assert && data.assert.length > 0) step.assert = data.assert;
  if (data.output && Object.keys(data.output).length > 0) step.output = data.output;
  if (data.retry) step.retry = data.retry;
  if (data.timeout) step.timeout = data.timeout;

  return step;
}

// Create edges between sequential nodes
export function createSequentialEdges(nodeIds: string[]): FlowEdge[] {
  const edges: FlowEdge[] = [];

  for (let i = 0; i < nodeIds.length - 1; i++) {
    edges.push({
      id: `edge_${nodeIds[i]}_${nodeIds[i + 1]}`,
      source: nodeIds[i],
      target: nodeIds[i + 1],
      type: 'smoothstep',
      animated: false,
    });
  }

  return edges;
}

// Convert FlowDefinition to nodes and edges
export function flowDefinitionToNodesAndEdges(
  definition: FlowDefinition
): { nodes: FlowNode[]; edges: FlowEdge[] } {
  const nodes: FlowNode[] = [];
  const edges: FlowEdge[] = [];
  let currentY = 50;

  // Add section label nodes (visual only)
  const sections: { name: FlowSection; steps: Step[] | undefined }[] = [
    { name: 'setup', steps: definition.setup },
    { name: 'steps', steps: definition.steps },
    { name: 'teardown', steps: definition.teardown },
  ];

  sections.forEach(({ name, steps }) => {
    if (!steps || steps.length === 0) return;

    // Add section header node
    const sectionHeaderId = `section_${name}`;
    const sectionData: SectionHeaderData = {
      label: name.charAt(0).toUpperCase() + name.slice(1),
      section: name,
    };

    nodes.push({
      id: sectionHeaderId,
      type: 'sectionHeader',
      position: { x: 300, y: currentY },
      data: sectionData,
      draggable: false,
      selectable: false,
    });

    currentY += 60;
    const nodeIdsInSection: string[] = [];

    // Convert steps to nodes
    steps.forEach((step) => {
      const position = { x: 300, y: currentY };
      const node = stepToNode(step, position, name);
      nodes.push(node);
      nodeIdsInSection.push(node.id);
      currentY += NODE_HEIGHT + NODE_VERTICAL_GAP;
    });

    // Create edges within section
    const sectionEdges = createSequentialEdges(nodeIdsInSection);
    edges.push(...sectionEdges);

    currentY += SECTION_GAP - NODE_VERTICAL_GAP;
  });

  return { nodes, edges };
}

// Convert nodes and edges back to FlowDefinition
export function nodesAndEdgesToFlowDefinition(
  nodes: FlowNode[],
  edges: FlowEdge[],
  existingDefinition: Partial<FlowDefinition> = {}
): FlowDefinition {
  // Filter out section headers and sort nodes by Y position
  const stepNodes = nodes
    .filter((n) => n.type === 'flowNode')
    .sort((a, b) => a.position.y - b.position.y);

  // Group nodes by section (stored in node data or determine by position)
  const setupNodes: FlowNode[] = [];
  const mainNodes: FlowNode[] = [];
  const teardownNodes: FlowNode[] = [];

  // Find section boundaries from section header nodes
  const sectionHeaders = nodes.filter((n) => n.type === 'sectionHeader');
  const setupHeader = sectionHeaders.find((n) => {
    const data = n.data as SectionHeaderData;
    return data.section === 'setup';
  });
  const stepsHeader = sectionHeaders.find((n) => {
    const data = n.data as SectionHeaderData;
    return data.section === 'steps';
  });
  const teardownHeader = sectionHeaders.find((n) => {
    const data = n.data as SectionHeaderData;
    return data.section === 'teardown';
  });

  stepNodes.forEach((node) => {
    const nodeY = node.position.y;

    // Determine section by position relative to headers
    if (teardownHeader && nodeY > teardownHeader.position.y) {
      teardownNodes.push(node);
    } else if (stepsHeader && nodeY > stepsHeader.position.y) {
      mainNodes.push(node);
    } else if (setupHeader && nodeY > setupHeader.position.y) {
      setupNodes.push(node);
    } else {
      // Default to main steps
      mainNodes.push(node);
    }
  });

  // Convert nodes to steps
  const setup = setupNodes.length > 0 ? setupNodes.map(nodeToStep) : undefined;
  const steps = mainNodes.map(nodeToStep);
  const teardown = teardownNodes.length > 0 ? teardownNodes.map(nodeToStep) : undefined;

  return {
    name: existingDefinition.name || 'Untitled Flow',
    description: existingDefinition.description || '',
    suite: existingDefinition.suite || '',
    tags: existingDefinition.tags || [],
    env: existingDefinition.env,
    setup,
    steps,
    teardown,
  };
}

// Convert FlowDefinition to YAML string
export function flowDefinitionToYaml(definition: FlowDefinition): string {
  const lines: string[] = [];

  // Header
  lines.push(`name: "${definition.name}"`);
  if (definition.description) {
    lines.push(`description: "${definition.description}"`);
  }
  if (definition.suite) {
    lines.push(`suite: "${definition.suite}"`);
  }
  if (definition.tags && definition.tags.length > 0) {
    lines.push(`tags: [${definition.tags.map((t) => `"${t}"`).join(', ')}]`);
  }

  // Environment variables
  if (definition.env && Object.keys(definition.env).length > 0) {
    lines.push('');
    lines.push('env:');
    Object.entries(definition.env).forEach(([key, value]) => {
      if (typeof value === 'string') {
        lines.push(`  ${key}: "${value}"`);
      } else {
        lines.push(`  ${key}: ${JSON.stringify(value)}`);
      }
    });
  }

  // Helper to convert steps to YAML
  const stepsToYaml = (steps: Step[], indent: string = '  '): string[] => {
    const stepLines: string[] = [];

    steps.forEach((step) => {
      stepLines.push(`${indent}- id: ${step.id}`);
      stepLines.push(`${indent}  action: ${step.action}`);

      if (step.name) {
        stepLines.push(`${indent}  name: "${step.name}"`);
      }
      if (step.description) {
        stepLines.push(`${indent}  description: "${step.description}"`);
      }

      // Config
      if (step.config && Object.keys(step.config).length > 0) {
        stepLines.push(`${indent}  config:`);
        Object.entries(step.config).forEach(([key, value]) => {
          if (value === null || value === undefined) return;
          if (typeof value === 'object') {
            stepLines.push(`${indent}    ${key}: ${JSON.stringify(value)}`);
          } else if (typeof value === 'string') {
            // Handle strings with special characters
            if (value.includes('${') || value.includes('"') || value.includes('\n')) {
              stepLines.push(`${indent}    ${key}: "${value.replace(/"/g, '\\"')}"`);
            } else {
              stepLines.push(`${indent}    ${key}: "${value}"`);
            }
          } else {
            stepLines.push(`${indent}    ${key}: ${value}`);
          }
        });
      }

      // Assertions
      if (step.assert && step.assert.length > 0) {
        stepLines.push(`${indent}  assert:`);
        step.assert.forEach((assertion) => {
          stepLines.push(`${indent}    - ${assertion}`);
        });
      }

      // Output
      if (step.output && Object.keys(step.output).length > 0) {
        stepLines.push(`${indent}  output:`);
        Object.entries(step.output).forEach(([key, value]) => {
          stepLines.push(`${indent}    ${key}: "${value}"`);
        });
      }

      // Retry
      if (step.retry) {
        stepLines.push(`${indent}  retry:`);
        stepLines.push(`${indent}    max_attempts: ${step.retry.max_attempts}`);
        stepLines.push(`${indent}    delay: "${step.retry.delay}"`);
        if (step.retry.backoff) {
          stepLines.push(`${indent}    backoff: "${step.retry.backoff}"`);
        }
      }

      // Timeout
      if (step.timeout) {
        stepLines.push(`${indent}  timeout: "${step.timeout}"`);
      }

      stepLines.push('');
    });

    return stepLines;
  };

  // Setup steps
  if (definition.setup && definition.setup.length > 0) {
    lines.push('');
    lines.push('setup:');
    lines.push(...stepsToYaml(definition.setup));
  }

  // Main steps
  lines.push('');
  lines.push('steps:');
  lines.push(...stepsToYaml(definition.steps));

  // Teardown steps
  if (definition.teardown && definition.teardown.length > 0) {
    lines.push('');
    lines.push('teardown:');
    lines.push(...stepsToYaml(definition.teardown));
  }

  return lines.join('\n');
}

// Validate a node's configuration
export function validateNodeConfig(node: FlowNode): string[] {
  const errors: string[] = [];

  if (!isFlowNodeData(node.data)) {
    return errors;
  }

  const data = node.data;

  switch (data.action) {
    case 'http_request':
      if (!data.config.url) {
        errors.push('URL is required');
      }
      if (!data.config.method) {
        errors.push('HTTP method is required');
      }
      break;

    case 'database_query':
      if (!data.config.connection) {
        errors.push('Database connection string is required');
      }
      if (!data.config.query) {
        errors.push('SQL query is required');
      }
      break;

    case 'log':
      if (!data.config.message) {
        errors.push('Log message is required');
      }
      break;

    case 'delay':
      if (!data.config.duration) {
        errors.push('Delay duration is required');
      }
      break;

    case 'assert':
      if (!data.config.expression) {
        errors.push('Assertion expression is required');
      }
      break;

    case 'transform':
      if (!data.config.expression) {
        errors.push('Transform expression is required');
      }
      if (!data.config.output_var) {
        errors.push('Output variable name is required');
      }
      break;

    case 'mock_server_start':
      if (!data.config.name) {
        errors.push('Mock server name is required');
      }
      break;

    case 'mock_server_stop':
      if (!data.config.name) {
        errors.push('Mock server name is required');
      }
      break;

    case 'contract_generate':
      if (!data.config.consumer) {
        errors.push('Consumer name is required');
      }
      if (!data.config.provider) {
        errors.push('Provider name is required');
      }
      break;

    case 'contract_verify':
      if (!data.config.contract_id) {
        errors.push('Contract ID is required');
      }
      if (!data.config.provider_base_url) {
        errors.push('Provider base URL is required');
      }
      break;
  }

  return errors;
}

// Get action category color
export function getActionColor(action: ActionType): string {
  switch (action) {
    case 'http_request':
      return 'bg-blue-500';
    case 'database_query':
      return 'bg-purple-500';
    case 'log':
      return 'bg-gray-500';
    case 'delay':
      return 'bg-yellow-500';
    case 'assert':
      return 'bg-green-500';
    case 'transform':
      return 'bg-orange-500';
    case 'condition':
      return 'bg-cyan-500';
    case 'for_each':
      return 'bg-indigo-500';
    case 'mock_server_start':
    case 'mock_server_stop':
      return 'bg-pink-500';
    case 'contract_generate':
    case 'contract_verify':
      return 'bg-teal-500';
    default:
      return 'bg-gray-500';
  }
}

// Get action icon name
export function getActionIcon(action: ActionType): string {
  switch (action) {
    case 'http_request':
      return 'Globe';
    case 'database_query':
      return 'Database';
    case 'log':
      return 'FileText';
    case 'delay':
      return 'Clock';
    case 'assert':
      return 'CheckCircle';
    case 'transform':
      return 'Wand2';
    case 'condition':
      return 'GitBranch';
    case 'for_each':
      return 'Repeat';
    case 'mock_server_start':
      return 'Server';
    case 'mock_server_stop':
      return 'ServerOff';
    case 'contract_generate':
      return 'FileCode';
    case 'contract_verify':
      return 'FileCheck';
    default:
      return 'Box';
  }
}
