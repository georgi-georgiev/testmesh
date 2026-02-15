'use client';

import { useState, useEffect } from 'react';
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
  Plus,
  Trash2,
  X,
  HelpCircle,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type { FlowNode, FlowNodeData, ActionType } from './types';
import { isFlowNodeData } from './types';

// Icon mapping
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

interface PropertiesPanelProps {
  node: FlowNode | null;
  onNodeUpdate: (nodeId: string, data: Partial<FlowNodeData>) => void;
  onClose?: () => void;
  className?: string;
}

export default function PropertiesPanel({
  node,
  onNodeUpdate,
  onClose,
  className,
}: PropertiesPanelProps) {
  const [localData, setLocalData] = useState<FlowNodeData | null>(null);

  // Sync local state with node prop
  useEffect(() => {
    if (node && isFlowNodeData(node.data)) {
      setLocalData({ ...node.data });
    } else {
      setLocalData(null);
    }
  }, [node]);

  // Update local state and propagate to parent
  const updateData = (updates: Partial<FlowNodeData>) => {
    if (!node || !localData) return;

    const newData = { ...localData, ...updates };
    setLocalData(newData);
    onNodeUpdate(node.id, updates);
  };

  // Update config property
  const updateConfig = (key: string, value: any) => {
    if (!localData) return;
    updateData({
      config: { ...localData.config, [key]: value },
    });
  };

  if (!node || !localData) {
    return (
      <div className={cn('w-80 border-l bg-muted/30 p-4', className)}>
        <div className="text-center py-12 text-muted-foreground">
          <p className="text-sm">Select a node to edit its properties</p>
        </div>
      </div>
    );
  }

  const Icon = actionIcons[localData.action] || HelpCircle;

  return (
    <div className={cn('w-80 border-l bg-muted/30 flex flex-col h-full', className)}>
      {/* Header */}
      <div className="p-3 border-b flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="p-1.5 rounded bg-muted">
            <Icon className="w-4 h-4 text-muted-foreground" />
          </div>
          <span className="font-semibold text-sm">Properties</span>
        </div>
        {onClose && (
          <Button variant="ghost" size="sm" onClick={onClose} className="h-6 w-6 p-0">
            <X className="w-4 h-4" />
          </Button>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        <Tabs defaultValue="general" className="w-full">
          <TabsList className="w-full justify-start rounded-none border-b px-3">
            <TabsTrigger value="general" className="text-xs">General</TabsTrigger>
            <TabsTrigger value="config" className="text-xs">Config</TabsTrigger>
            <TabsTrigger value="assert" className="text-xs">Assert</TabsTrigger>
            <TabsTrigger value="output" className="text-xs">Output</TabsTrigger>
          </TabsList>

          {/* General Tab */}
          <TabsContent value="general" className="p-3 space-y-4">
            <div className="space-y-2">
              <Label htmlFor="step-id" className="text-xs">Step ID</Label>
              <Input
                id="step-id"
                value={localData.stepId}
                onChange={(e) => updateData({ stepId: e.target.value })}
                placeholder="step_id"
                className="h-8 text-sm font-mono"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="step-name" className="text-xs">Name</Label>
              <Input
                id="step-name"
                value={localData.name || ''}
                onChange={(e) => updateData({ name: e.target.value, label: e.target.value })}
                placeholder="Step name"
                className="h-8 text-sm"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="step-description" className="text-xs">Description</Label>
              <Textarea
                id="step-description"
                value={localData.description || ''}
                onChange={(e) => updateData({ description: e.target.value })}
                placeholder="Optional description"
                className="text-sm resize-none"
                rows={2}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="step-timeout" className="text-xs">Timeout</Label>
              <Input
                id="step-timeout"
                value={localData.timeout || ''}
                onChange={(e) => updateData({ timeout: e.target.value })}
                placeholder="e.g., 30s"
                className="h-8 text-sm"
              />
            </div>
          </TabsContent>

          {/* Config Tab - Action-specific */}
          <TabsContent value="config" className="p-3 space-y-4">
            <ActionConfig
              action={localData.action}
              config={localData.config}
              onConfigChange={updateConfig}
            />
          </TabsContent>

          {/* Assert Tab */}
          <TabsContent value="assert" className="p-3 space-y-4">
            <AssertionsEditor
              assertions={localData.assert || []}
              onChange={(assertions) => updateData({ assert: assertions })}
            />
          </TabsContent>

          {/* Output Tab */}
          <TabsContent value="output" className="p-3 space-y-4">
            <OutputEditor
              output={localData.output || {}}
              onChange={(output) => updateData({ output })}
            />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}

// Action-specific configuration components
function ActionConfig({
  action,
  config,
  onConfigChange,
}: {
  action: ActionType;
  config: Record<string, any>;
  onConfigChange: (key: string, value: any) => void;
}) {
  switch (action) {
    case 'http_request':
      return <HTTPRequestConfig config={config} onChange={onConfigChange} />;
    case 'database_query':
      return <DatabaseQueryConfig config={config} onChange={onConfigChange} />;
    case 'log':
      return <LogConfig config={config} onChange={onConfigChange} />;
    case 'delay':
      return <DelayConfig config={config} onChange={onConfigChange} />;
    case 'assert':
      return <AssertConfig config={config} onChange={onConfigChange} />;
    case 'transform':
      return <TransformConfig config={config} onChange={onConfigChange} />;
    case 'mock_server_start':
      return <MockServerStartConfig config={config} onChange={onConfigChange} />;
    case 'mock_server_stop':
      return <MockServerStopConfig config={config} onChange={onConfigChange} />;
    case 'contract_generate':
      return <ContractGenerateConfig config={config} onChange={onConfigChange} />;
    case 'contract_verify':
      return <ContractVerifyConfig config={config} onChange={onConfigChange} />;
    default:
      return (
        <div className="text-sm text-muted-foreground">
          No configuration options for this action
        </div>
      );
  }
}

// HTTP Request Config
function HTTPRequestConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Method</Label>
        <Select value={config.method || 'GET'} onValueChange={(v) => onChange('method', v)}>
          <SelectTrigger className="h-8 text-sm">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="GET">GET</SelectItem>
            <SelectItem value="POST">POST</SelectItem>
            <SelectItem value="PUT">PUT</SelectItem>
            <SelectItem value="DELETE">DELETE</SelectItem>
            <SelectItem value="PATCH">PATCH</SelectItem>
            <SelectItem value="HEAD">HEAD</SelectItem>
            <SelectItem value="OPTIONS">OPTIONS</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label className="text-xs">URL</Label>
        <Input
          value={config.url || ''}
          onChange={(e) => onChange('url', e.target.value)}
          placeholder="https://api.example.com/endpoint"
          className="h-8 text-sm font-mono"
        />
        <p className="text-[10px] text-muted-foreground">
          Use {'${VAR}'} for variables, e.g., {'${BASE_URL}/users'}
        </p>
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Headers (JSON)</Label>
        <Textarea
          value={JSON.stringify(config.headers || {}, null, 2)}
          onChange={(e) => {
            try {
              onChange('headers', JSON.parse(e.target.value));
            } catch {
              // Invalid JSON, ignore
            }
          }}
          placeholder='{"Content-Type": "application/json"}'
          className="text-xs font-mono resize-none"
          rows={3}
        />
      </div>

      {['POST', 'PUT', 'PATCH'].includes(config.method) && (
        <div className="space-y-2">
          <Label className="text-xs">Body (JSON)</Label>
          <Textarea
            value={typeof config.body === 'object' ? JSON.stringify(config.body, null, 2) : config.body || ''}
            onChange={(e) => {
              try {
                onChange('body', JSON.parse(e.target.value));
              } catch {
                onChange('body', e.target.value);
              }
            }}
            placeholder='{"key": "value"}'
            className="text-xs font-mono resize-none"
            rows={4}
          />
        </div>
      )}
    </div>
  );
}

// Database Query Config
function DatabaseQueryConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Connection String</Label>
        <Input
          value={config.connection || ''}
          onChange={(e) => onChange('connection', e.target.value)}
          placeholder="postgresql://user:pass@host:5432/db"
          className="h-8 text-sm font-mono"
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">SQL Query</Label>
        <Textarea
          value={config.query || ''}
          onChange={(e) => onChange('query', e.target.value)}
          placeholder="SELECT * FROM users WHERE id = $1"
          className="text-xs font-mono resize-none"
          rows={4}
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Parameters (JSON Array)</Label>
        <Input
          value={JSON.stringify(config.params || [])}
          onChange={(e) => {
            try {
              onChange('params', JSON.parse(e.target.value));
            } catch {
              // Invalid JSON
            }
          }}
          placeholder="[1, 'value']"
          className="h-8 text-sm font-mono"
        />
      </div>
    </div>
  );
}

// Log Config
function LogConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Message</Label>
        <Textarea
          value={config.message || ''}
          onChange={(e) => onChange('message', e.target.value)}
          placeholder="Log message with ${variables}"
          className="text-sm resize-none"
          rows={3}
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Level</Label>
        <Select value={config.level || 'info'} onValueChange={(v) => onChange('level', v)}>
          <SelectTrigger className="h-8 text-sm">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="debug">Debug</SelectItem>
            <SelectItem value="info">Info</SelectItem>
            <SelectItem value="warn">Warning</SelectItem>
            <SelectItem value="error">Error</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}

// Delay Config
function DelayConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Duration</Label>
        <Input
          value={config.duration || ''}
          onChange={(e) => onChange('duration', e.target.value)}
          placeholder="1s, 500ms, 2m"
          className="h-8 text-sm"
        />
        <p className="text-[10px] text-muted-foreground">
          Supports: ms (milliseconds), s (seconds), m (minutes)
        </p>
      </div>
    </div>
  );
}

// Assert Config
function AssertConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Expression</Label>
        <Input
          value={config.expression || ''}
          onChange={(e) => onChange('expression', e.target.value)}
          placeholder="status == 200"
          className="h-8 text-sm font-mono"
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Message (on failure)</Label>
        <Input
          value={config.message || ''}
          onChange={(e) => onChange('message', e.target.value)}
          placeholder="Expected status 200"
          className="h-8 text-sm"
        />
      </div>
    </div>
  );
}

// Transform Config
function TransformConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Input</Label>
        <Input
          value={config.input || ''}
          onChange={(e) => onChange('input', e.target.value)}
          placeholder="${step_output.data}"
          className="h-8 text-sm font-mono"
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Expression</Label>
        <Textarea
          value={config.expression || ''}
          onChange={(e) => onChange('expression', e.target.value)}
          placeholder="Transform expression"
          className="text-sm font-mono resize-none"
          rows={3}
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Output Variable</Label>
        <Input
          value={config.output_var || ''}
          onChange={(e) => onChange('output_var', e.target.value)}
          placeholder="transformed_data"
          className="h-8 text-sm font-mono"
        />
      </div>
    </div>
  );
}

// Mock Server Start Config
function MockServerStartConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Server Name</Label>
        <Input
          value={config.name || ''}
          onChange={(e) => onChange('name', e.target.value)}
          placeholder="mock-api"
          className="h-8 text-sm"
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Port</Label>
        <Input
          type="number"
          value={config.port || 5016}
          onChange={(e) => onChange('port', parseInt(e.target.value) || 5016)}
          placeholder="5016"
          className="h-8 text-sm"
        />
      </div>

      <div className="text-xs text-muted-foreground">
        Configure endpoints in the advanced settings or edit YAML directly
      </div>
    </div>
  );
}

// Mock Server Stop Config
function MockServerStopConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Server Name</Label>
        <Input
          value={config.name || ''}
          onChange={(e) => onChange('name', e.target.value)}
          placeholder="mock-api"
          className="h-8 text-sm"
        />
      </div>
    </div>
  );
}

// Contract Generate Config
function ContractGenerateConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Consumer</Label>
        <Input
          value={config.consumer || ''}
          onChange={(e) => onChange('consumer', e.target.value)}
          placeholder="frontend-app"
          className="h-8 text-sm"
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Provider</Label>
        <Input
          value={config.provider || ''}
          onChange={(e) => onChange('provider', e.target.value)}
          placeholder="user-service"
          className="h-8 text-sm"
        />
      </div>

      <div className="text-xs text-muted-foreground">
        Configure interactions in the advanced settings or edit YAML directly
      </div>
    </div>
  );
}

// Contract Verify Config
function ContractVerifyConfig({
  config,
  onChange,
}: {
  config: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label className="text-xs">Contract ID</Label>
        <Input
          value={config.contract_id || ''}
          onChange={(e) => onChange('contract_id', e.target.value)}
          placeholder="contract-uuid"
          className="h-8 text-sm font-mono"
        />
      </div>

      <div className="space-y-2">
        <Label className="text-xs">Provider Base URL</Label>
        <Input
          value={config.provider_base_url || ''}
          onChange={(e) => onChange('provider_base_url', e.target.value)}
          placeholder="http://localhost:5016"
          className="h-8 text-sm"
        />
      </div>
    </div>
  );
}

// Assertions Editor
function AssertionsEditor({
  assertions,
  onChange,
}: {
  assertions: string[];
  onChange: (assertions: string[]) => void;
}) {
  const addAssertion = () => {
    onChange([...assertions, '']);
  };

  const updateAssertion = (index: number, value: string) => {
    const newAssertions = [...assertions];
    newAssertions[index] = value;
    onChange(newAssertions);
  };

  const removeAssertion = (index: number) => {
    onChange(assertions.filter((_, i) => i !== index));
  };

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <Label className="text-xs">Assertions</Label>
        <Button variant="ghost" size="sm" onClick={addAssertion} className="h-6 px-2 text-xs">
          <Plus className="w-3 h-3 mr-1" />
          Add
        </Button>
      </div>

      {assertions.length === 0 ? (
        <p className="text-xs text-muted-foreground">No assertions defined</p>
      ) : (
        <div className="space-y-2">
          {assertions.map((assertion, index) => (
            <div key={index} className="flex items-center gap-2">
              <Input
                value={assertion}
                onChange={(e) => updateAssertion(index, e.target.value)}
                placeholder="status == 200"
                className="h-7 text-xs font-mono flex-1"
              />
              <Button
                variant="ghost"
                size="sm"
                onClick={() => removeAssertion(index)}
                className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
              >
                <Trash2 className="w-3 h-3" />
              </Button>
            </div>
          ))}
        </div>
      )}

      <p className="text-[10px] text-muted-foreground">
        Examples: status == 200, body.id != null, body.items.length {'>'} 0
      </p>
    </div>
  );
}

// Output Editor
function OutputEditor({
  output,
  onChange,
}: {
  output: Record<string, string>;
  onChange: (output: Record<string, string>) => void;
}) {
  const entries = Object.entries(output);

  const addOutput = () => {
    onChange({ ...output, '': '' });
  };

  const updateOutput = (oldKey: string, newKey: string, value: string) => {
    const newOutput = { ...output };
    if (oldKey !== newKey) {
      delete newOutput[oldKey];
    }
    newOutput[newKey] = value;
    onChange(newOutput);
  };

  const removeOutput = (key: string) => {
    const newOutput = { ...output };
    delete newOutput[key];
    onChange(newOutput);
  };

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <Label className="text-xs">Output Variables</Label>
        <Button variant="ghost" size="sm" onClick={addOutput} className="h-6 px-2 text-xs">
          <Plus className="w-3 h-3 mr-1" />
          Add
        </Button>
      </div>

      {entries.length === 0 ? (
        <p className="text-xs text-muted-foreground">No output variables defined</p>
      ) : (
        <div className="space-y-2">
          {entries.map(([key, value], index) => (
            <div key={index} className="flex items-center gap-2">
              <Input
                value={key}
                onChange={(e) => updateOutput(key, e.target.value, value)}
                placeholder="var_name"
                className="h-7 text-xs font-mono w-24"
              />
              <span className="text-muted-foreground">=</span>
              <Input
                value={value}
                onChange={(e) => updateOutput(key, key, e.target.value)}
                placeholder="$.path.to.value"
                className="h-7 text-xs font-mono flex-1"
              />
              <Button
                variant="ghost"
                size="sm"
                onClick={() => removeOutput(key)}
                className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
              >
                <Trash2 className="w-3 h-3" />
              </Button>
            </div>
          ))}
        </div>
      )}

      <p className="text-[10px] text-muted-foreground">
        Use JSONPath to extract values: $.id, $.data[0].name
      </p>
    </div>
  );
}
