'use client';

import { memo } from 'react';
import { Handle, Position, type NodeProps } from 'reactflow';
import {
  Globe,
  Database,
  FileText,
  Clock,
  CheckCircle,
  Wand2,
  GitBranch,
  Repeat,
  Server,
  ServerOff,
  FileCode,
  FileCheck,
  Box,
  AlertCircle,
  CheckCircle2,
  Loader2,
  XCircle,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import type { FlowNodeData, ActionType } from '../types';

// Icon mapping for action types
const actionIcons: Record<ActionType, React.ElementType> = {
  http_request: Globe,
  database_query: Database,
  log: FileText,
  delay: Clock,
  assert: CheckCircle,
  transform: Wand2,
  condition: GitBranch,
  for_each: Repeat,
  mock_server_start: Server,
  mock_server_stop: ServerOff,
  contract_generate: FileCode,
  contract_verify: FileCheck,
};

// Color mapping for action types
const actionColors: Record<ActionType, { bg: string; border: string; icon: string }> = {
  http_request: {
    bg: 'bg-blue-50 dark:bg-blue-950',
    border: 'border-blue-200 dark:border-blue-800',
    icon: 'text-blue-500',
  },
  database_query: {
    bg: 'bg-purple-50 dark:bg-purple-950',
    border: 'border-purple-200 dark:border-purple-800',
    icon: 'text-purple-500',
  },
  log: {
    bg: 'bg-gray-50 dark:bg-gray-900',
    border: 'border-gray-200 dark:border-gray-700',
    icon: 'text-gray-500',
  },
  delay: {
    bg: 'bg-yellow-50 dark:bg-yellow-950',
    border: 'border-yellow-200 dark:border-yellow-800',
    icon: 'text-yellow-500',
  },
  assert: {
    bg: 'bg-green-50 dark:bg-green-950',
    border: 'border-green-200 dark:border-green-800',
    icon: 'text-green-500',
  },
  transform: {
    bg: 'bg-orange-50 dark:bg-orange-950',
    border: 'border-orange-200 dark:border-orange-800',
    icon: 'text-orange-500',
  },
  condition: {
    bg: 'bg-cyan-50 dark:bg-cyan-950',
    border: 'border-cyan-200 dark:border-cyan-800',
    icon: 'text-cyan-500',
  },
  for_each: {
    bg: 'bg-indigo-50 dark:bg-indigo-950',
    border: 'border-indigo-200 dark:border-indigo-800',
    icon: 'text-indigo-500',
  },
  mock_server_start: {
    bg: 'bg-pink-50 dark:bg-pink-950',
    border: 'border-pink-200 dark:border-pink-800',
    icon: 'text-pink-500',
  },
  mock_server_stop: {
    bg: 'bg-pink-50 dark:bg-pink-950',
    border: 'border-pink-200 dark:border-pink-800',
    icon: 'text-pink-500',
  },
  contract_generate: {
    bg: 'bg-teal-50 dark:bg-teal-950',
    border: 'border-teal-200 dark:border-teal-800',
    icon: 'text-teal-500',
  },
  contract_verify: {
    bg: 'bg-teal-50 dark:bg-teal-950',
    border: 'border-teal-200 dark:border-teal-800',
    icon: 'text-teal-500',
  },
};

// Status icon component
function StatusIcon({ status }: { status?: FlowNodeData['status'] }) {
  switch (status) {
    case 'running':
      return <Loader2 className="w-4 h-4 text-blue-500 animate-spin" />;
    case 'completed':
      return <CheckCircle2 className="w-4 h-4 text-green-500" />;
    case 'failed':
      return <XCircle className="w-4 h-4 text-red-500" />;
    case 'skipped':
      return <AlertCircle className="w-4 h-4 text-gray-400" />;
    default:
      return null;
  }
}

// Get display info for HTTP method
function getHttpMethodBadge(method: string) {
  const colors: Record<string, string> = {
    GET: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
    POST: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
    PUT: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300',
    DELETE: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
    PATCH: 'bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300',
    HEAD: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
    OPTIONS: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
  };

  return colors[method] || colors.GET;
}

// Node content based on action type
function NodeContent({ data }: { data: FlowNodeData }) {
  const { action, config } = data;

  switch (action) {
    case 'http_request':
      return (
        <div className="flex items-center gap-2 text-xs truncate">
          <span
            className={cn(
              'px-1.5 py-0.5 rounded font-mono font-medium',
              getHttpMethodBadge(config.method)
            )}
          >
            {config.method}
          </span>
          <span className="truncate text-muted-foreground font-mono">
            {config.url || 'No URL'}
          </span>
        </div>
      );

    case 'database_query':
      return (
        <div className="text-xs text-muted-foreground truncate font-mono">
          {config.query?.substring(0, 40) || 'No query'}
          {config.query?.length > 40 ? '...' : ''}
        </div>
      );

    case 'log':
      return (
        <div className="text-xs text-muted-foreground truncate">
          {config.message || 'No message'}
        </div>
      );

    case 'delay':
      return (
        <div className="text-xs text-muted-foreground">
          Wait {config.duration || '0s'}
        </div>
      );

    case 'assert':
      return (
        <div className="text-xs text-muted-foreground font-mono truncate">
          {config.expression || 'No expression'}
        </div>
      );

    case 'transform':
      return (
        <div className="text-xs text-muted-foreground truncate">
          → {config.output_var || 'output'}
        </div>
      );

    case 'condition':
      return (
        <div className="text-xs text-muted-foreground font-mono truncate">
          if {config.expression || '...'}
        </div>
      );

    case 'for_each':
      return (
        <div className="text-xs text-muted-foreground truncate">
          foreach {config.item_var || 'item'} in {config.items || '[]'}
        </div>
      );

    case 'mock_server_start':
      return (
        <div className="text-xs text-muted-foreground truncate">
          Start: {config.name || 'unnamed'}
          {config.port && ` :${config.port}`}
        </div>
      );

    case 'mock_server_stop':
      return (
        <div className="text-xs text-muted-foreground truncate">
          Stop: {config.name || 'unnamed'}
        </div>
      );

    case 'contract_generate':
      return (
        <div className="text-xs text-muted-foreground truncate">
          {config.consumer || '?'} → {config.provider || '?'}
        </div>
      );

    case 'contract_verify':
      return (
        <div className="text-xs text-muted-foreground truncate">
          Verify: {config.contract_id?.substring(0, 8) || 'no contract'}
        </div>
      );

    default:
      return null;
  }
}

// Main FlowNode component
function FlowNode({ data, selected }: NodeProps<FlowNodeData>) {
  const Icon = actionIcons[data.action] || Box;
  const colors = actionColors[data.action] || actionColors.log;

  return (
    <div
      className={cn(
        'rounded-lg border-2 shadow-sm transition-all',
        'min-w-[240px] max-w-[320px]',
        colors.bg,
        colors.border,
        selected && 'ring-2 ring-primary ring-offset-2',
        data.status === 'failed' && 'border-red-500',
        data.status === 'running' && 'border-blue-500 animate-pulse'
      )}
    >
      {/* Input Handle */}
      <Handle
        type="target"
        position={Position.Top}
        className="!w-3 !h-3 !bg-muted-foreground !border-2 !border-background"
      />

      <div className="p-3">
        {/* Header */}
        <div className="flex items-center gap-2 mb-1">
          <div className={cn('p-1 rounded', colors.icon)}>
            <Icon className="w-4 h-4" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="font-medium text-sm truncate">
              {data.name || data.label}
            </div>
          </div>
          <StatusIcon status={data.status} />
        </div>

        {/* Action-specific content */}
        <NodeContent data={data} />

        {/* Badges */}
        <div className="flex items-center gap-1 mt-2 flex-wrap">
          {data.assert && data.assert.length > 0 && (
            <span className="text-[10px] px-1.5 py-0.5 rounded bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300">
              {data.assert.length} assertion{data.assert.length > 1 ? 's' : ''}
            </span>
          )}
          {data.output && Object.keys(data.output).length > 0 && (
            <span className="text-[10px] px-1.5 py-0.5 rounded bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300">
              {Object.keys(data.output).length} output{Object.keys(data.output).length > 1 ? 's' : ''}
            </span>
          )}
          {data.retry && (
            <span className="text-[10px] px-1.5 py-0.5 rounded bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300">
              retry ×{data.retry.max_attempts}
            </span>
          )}
        </div>
      </div>

      {/* Output Handle */}
      <Handle
        type="source"
        position={Position.Bottom}
        className="!w-3 !h-3 !bg-muted-foreground !border-2 !border-background"
      />
    </div>
  );
}

export default memo(FlowNode);
