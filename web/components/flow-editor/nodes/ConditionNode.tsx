'use client';

import { memo } from 'react';
import { Handle, Position, type NodeProps } from 'reactflow';
import { GitBranch, CheckCircle2, XCircle, Loader2, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { ConditionNodeData } from '../types';

// Status icon component
function StatusIcon({ status }: { status?: ConditionNodeData['status'] }) {
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

/**
 * ConditionNode - A diamond-shaped decision node for conditional branching
 *
 * Visual design:
 * - Diamond shape rotated 45 degrees
 * - Input handle at top
 * - "Then" output handle on the right (true branch)
 * - "Else" output handle on the left (false branch)
 * - Shows the condition expression
 */
function ConditionNode({ data, selected }: NodeProps<ConditionNodeData>) {
  const expression = data.config?.expression || 'condition';

  return (
    <div className="relative">
      {/* Input Handle (top) */}
      <Handle
        type="target"
        position={Position.Top}
        id="input"
        className="!w-3 !h-3 !bg-muted-foreground !border-2 !border-background"
        style={{ top: -6 }}
      />

      {/* Diamond shape container */}
      <div
        className={cn(
          'w-[140px] h-[140px] flex items-center justify-center',
          'transform rotate-45',
          'rounded-lg border-2 shadow-sm transition-all',
          'bg-cyan-50 dark:bg-cyan-950',
          'border-cyan-200 dark:border-cyan-800',
          selected && 'ring-2 ring-primary ring-offset-2',
          data.status === 'failed' && 'border-red-500',
          data.status === 'running' && 'border-blue-500 animate-pulse'
        )}
      >
        {/* Inner content (counter-rotate to keep text upright) */}
        <div className="transform -rotate-45 text-center p-2 max-w-[100px]">
          <div className="flex items-center justify-center gap-1 mb-1">
            <GitBranch className="w-4 h-4 text-cyan-500" />
            <StatusIcon status={data.status} />
          </div>
          <div className="text-xs font-medium truncate">
            {data.name || 'Condition'}
          </div>
          <div className="text-[10px] text-muted-foreground font-mono truncate mt-1">
            {expression.length > 15 ? expression.substring(0, 15) + '...' : expression}
          </div>
        </div>
      </div>

      {/* Then handle (right) - True branch */}
      <Handle
        type="source"
        position={Position.Right}
        id="then"
        className="!w-3 !h-3 !bg-green-500 !border-2 !border-background"
        style={{ right: -6, top: '50%' }}
      />
      <div
        className="absolute text-[10px] font-medium text-green-600 dark:text-green-400"
        style={{ right: -30, top: '50%', transform: 'translateY(-50%)' }}
      >
        true
      </div>

      {/* Else handle (left) - False branch */}
      <Handle
        type="source"
        position={Position.Left}
        id="else"
        className="!w-3 !h-3 !bg-red-500 !border-2 !border-background"
        style={{ left: -6, top: '50%' }}
      />
      <div
        className="absolute text-[10px] font-medium text-red-600 dark:text-red-400"
        style={{ left: -32, top: '50%', transform: 'translateY(-50%)' }}
      >
        false
      </div>

      {/* Continuation handle (bottom) - After both branches merge */}
      <Handle
        type="source"
        position={Position.Bottom}
        id="next"
        className="!w-3 !h-3 !bg-muted-foreground !border-2 !border-background"
        style={{ bottom: -6 }}
      />
    </div>
  );
}

export default memo(ConditionNode);
