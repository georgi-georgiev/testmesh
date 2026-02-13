'use client';

import { useState, useCallback, useEffect } from 'react';
import {
  Save,
  Undo,
  Redo,
  Play,
  Code,
  LayoutGrid,
  Settings2,
  ChevronLeft,
  ChevronRight,
  AlertCircle,
  CheckCircle2,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Textarea } from '@/components/ui/textarea';
import { Alert, AlertDescription } from '@/components/ui/alert';
import type { FlowDefinition } from '@/lib/api/types';
import type { FlowNode, FlowNodeData, PaletteItem } from './types';
import { flowDefinitionToYaml } from './utils';

import FlowCanvas from './FlowCanvas';
import NodePalette from './NodePalette';
import PropertiesPanel from './PropertiesPanel';

interface FlowEditorProps {
  initialDefinition?: FlowDefinition;
  initialYaml?: string;
  onSave?: (yaml: string, definition: FlowDefinition) => void;
  onRun?: (definition: FlowDefinition) => void;
  isSaving?: boolean;
  isRunning?: boolean;
  className?: string;
}

// Simple YAML parser for flow definitions
// Supports both root-level and `flow:` wrapped structures
function parseYaml(yaml: string): FlowDefinition | null {
  try {
    const lines = yaml.split('\n');
    const definition: Partial<FlowDefinition> = {
      name: '',
      description: '',
      suite: '',
      tags: [],
      steps: [],
    };

    let currentSection: 'root' | 'env' | 'setup' | 'steps' | 'teardown' = 'root';
    let currentStep: any = null;
    let currentKey = '';

    // Detect if YAML uses `flow:` wrapper (indent offset)
    let baseIndent = 0;
    for (const line of lines) {
      const trimmed = line.trim();
      if (!trimmed || trimmed.startsWith('#')) continue;
      if (trimmed === 'flow:') {
        baseIndent = 2; // Properties are nested under flow:
        break;
      }
      // If first non-comment line is not `flow:`, assume root-level
      break;
    }

    for (const line of lines) {
      const trimmed = line.trim();
      if (!trimmed || trimmed.startsWith('#')) continue;
      if (trimmed === 'flow:') continue; // Skip the wrapper line

      // Calculate indent level relative to base
      const lineIndent = line.search(/\S/) - baseIndent;

      // Parse key-value pairs
      const kvMatch = trimmed.match(/^(\w+):\s*(.*)$/);
      if (kvMatch) {
        const [, key, value] = kvMatch;

        if (lineIndent === 0) {
          // Root level (flow properties)
          switch (key) {
            case 'name':
              definition.name = value.replace(/^["']|["']$/g, '');
              break;
            case 'description':
              definition.description = value.replace(/^["']|["']$/g, '');
              break;
            case 'suite':
              definition.suite = value.replace(/^["']|["']$/g, '');
              break;
            case 'tags':
              if (value.startsWith('[')) {
                definition.tags = JSON.parse(value.replace(/'/g, '"'));
              }
              break;
            case 'env':
              currentSection = 'env';
              definition.env = {};
              break;
            case 'setup':
              currentSection = 'setup';
              definition.setup = [];
              break;
            case 'steps':
              currentSection = 'steps';
              definition.steps = [];
              break;
            case 'teardown':
              currentSection = 'teardown';
              definition.teardown = [];
              break;
          }
        } else if (currentSection === 'env' && lineIndent === 2) {
          definition.env = definition.env || {};
          definition.env[key] = value.replace(/^["']|["']$/g, '');
        } else if (currentStep && ['setup', 'steps', 'teardown'].includes(currentSection)) {
          // Step properties
          if (key === 'config') {
            currentStep.config = {};
            currentKey = 'config';
          } else if (key === 'assert') {
            currentStep.assert = [];
            currentKey = 'assert';
          } else if (key === 'output') {
            currentStep.output = {};
            currentKey = 'output';
          } else if (currentKey === 'config') {
            currentStep.config[key] = value.replace(/^["']|["']$/g, '');
          } else if (currentKey === 'output') {
            currentStep.output[key] = value.replace(/^["']|["']$/g, '');
          } else {
            currentStep[key] = value.replace(/^["']|["']$/g, '');
          }
        }
      }

      // Parse list items
      const listMatch = trimmed.match(/^-\s*(.+)$/);
      if (listMatch) {
        const listValue = listMatch[1];

        if (['setup', 'steps', 'teardown'].includes(currentSection)) {
          if (listValue.startsWith('id:')) {
            // New step
            if (currentStep) {
              const targetArray = definition[currentSection as 'setup' | 'steps' | 'teardown'];
              if (targetArray) {
                targetArray.push(currentStep);
              }
            }
            currentStep = {
              id: listValue.replace('id:', '').trim(),
              config: {},
            };
            currentKey = '';
          } else if (currentKey === 'assert') {
            currentStep.assert = currentStep.assert || [];
            currentStep.assert.push(listValue);
          }
        }
      }
    }

    // Add last step
    if (currentStep && ['setup', 'steps', 'teardown'].includes(currentSection)) {
      const targetArray = definition[currentSection as 'setup' | 'steps' | 'teardown'];
      if (targetArray) {
        targetArray.push(currentStep);
      }
    }

    return definition as FlowDefinition;
  } catch (error) {
    console.error('YAML parsing error:', error);
    return null;
  }
}

export default function FlowEditor({
  initialDefinition,
  initialYaml,
  onSave,
  onRun,
  isSaving = false,
  isRunning = false,
  className,
}: FlowEditorProps) {
  // Editor mode: 'visual' or 'yaml'
  const [mode, setMode] = useState<'visual' | 'yaml'>('visual');

  // Flow definition state
  const [definition, setDefinition] = useState<FlowDefinition>(
    initialDefinition || {
      name: 'Untitled Flow',
      description: '',
      suite: '',
      tags: [],
      steps: [],
    }
  );

  // YAML state (for yaml mode)
  const [yaml, setYaml] = useState(initialYaml || '');
  const [yamlError, setYamlError] = useState<string | null>(null);
  const [validationSuccess, setValidationSuccess] = useState<string | null>(null);

  // UI state
  const [selectedNode, setSelectedNode] = useState<FlowNode | null>(null);
  const [showPalette, setShowPalette] = useState(true);
  const [showProperties, setShowProperties] = useState(true);
  const [isDirty, setIsDirty] = useState(false);

  // History for undo/redo
  const [history, setHistory] = useState<FlowDefinition[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);

  // Sync YAML when switching modes
  useEffect(() => {
    if (mode === 'yaml' && definition) {
      setYaml(flowDefinitionToYaml(definition));
    }
  }, [mode]);

  // Handle mode switch
  const handleModeChange = useCallback((newMode: string) => {
    setValidationSuccess(null);
    if (newMode === 'yaml') {
      // Visual → YAML: Generate YAML from definition
      setYaml(flowDefinitionToYaml(definition));
      setYamlError(null);
    } else {
      // YAML → Visual: Parse YAML to definition
      const parsed = parseYaml(yaml);
      if (parsed) {
        setDefinition(parsed);
        setYamlError(null);
      } else {
        setYamlError('Invalid YAML syntax. Please fix before switching to visual mode.');
        return; // Don't switch modes
      }
    }
    setMode(newMode as 'visual' | 'yaml');
  }, [definition, yaml]);

  // Handle definition changes (from visual editor)
  const handleDefinitionChange = useCallback((newDefinition: FlowDefinition) => {
    setDefinition(newDefinition);
    setIsDirty(true);

    // Add to history
    setHistory((prev) => [...prev.slice(0, historyIndex + 1), newDefinition]);
    setHistoryIndex((prev) => prev + 1);
  }, [historyIndex]);

  // Handle YAML changes
  const handleYamlChange = useCallback((newYaml: string) => {
    setYaml(newYaml);
    setIsDirty(true);
    setYamlError(null);

    // Try to parse for validation
    const parsed = parseYaml(newYaml);
    if (!parsed) {
      setYamlError('Invalid YAML syntax');
    }
    setValidationSuccess(null);
  }, []);

  // Handle validate YAML
  const handleValidate = useCallback(() => {
    setYamlError(null);
    setValidationSuccess(null);

    const parsed = parseYaml(yaml);
    if (!parsed) {
      setYamlError('Invalid YAML syntax');
      return;
    }

    // Check required fields
    const errors: string[] = [];
    if (!parsed.name || parsed.name === 'Untitled Flow') {
      errors.push('Flow name is required');
    }
    if (!parsed.steps || parsed.steps.length === 0) {
      errors.push('At least one step is required');
    }

    // Check each step has required fields
    parsed.steps?.forEach((step, index) => {
      if (!step.id) {
        errors.push(`Step ${index + 1}: missing 'id'`);
      }
      if (!step.action) {
        errors.push(`Step ${index + 1} (${step.id || 'unnamed'}): missing 'action'`);
      }
    });

    if (errors.length > 0) {
      setYamlError(`Validation failed:\n• ${errors.join('\n• ')}`);
    } else {
      setValidationSuccess(`Valid! Flow "${parsed.name}" with ${parsed.steps?.length || 0} steps`);
    }
  }, [yaml]);

  // Handle node selection
  const handleNodeSelect = useCallback((node: FlowNode | null) => {
    setSelectedNode(node);
    if (node) {
      setShowProperties(true);
    }
  }, []);

  // Handle node update from properties panel
  const handleNodeUpdate = useCallback((nodeId: string, data: Partial<FlowNodeData>) => {
    // Update the definition based on node changes
    setDefinition((prev) => {
      const updateSteps = (steps: any[] | undefined) => {
        if (!steps) return steps;
        return steps.map((step) => {
          if (step.id === data.stepId || step.id === nodeId) {
            return {
              ...step,
              ...data,
              id: data.stepId || step.id,
            };
          }
          return step;
        });
      };

      return {
        ...prev,
        setup: updateSteps(prev.setup),
        steps: updateSteps(prev.steps) || [],
        teardown: updateSteps(prev.teardown),
      };
    });
    setIsDirty(true);
  }, []);

  // Handle save
  const handleSave = useCallback(() => {
    let finalYaml = yaml;
    let finalDefinition = definition;

    if (mode === 'visual') {
      finalYaml = flowDefinitionToYaml(definition);
    } else {
      const parsed = parseYaml(yaml);
      if (parsed) {
        finalDefinition = parsed;
      } else {
        setYamlError('Cannot save: Invalid YAML syntax');
        return;
      }
    }

    onSave?.(finalYaml, finalDefinition);
    setIsDirty(false);
  }, [mode, yaml, definition, onSave]);

  // Handle run
  const handleRun = useCallback(() => {
    let finalDefinition = definition;

    if (mode === 'yaml') {
      const parsed = parseYaml(yaml);
      if (parsed) {
        finalDefinition = parsed;
      } else {
        setYamlError('Cannot run: Invalid YAML syntax');
        return;
      }
    }

    onRun?.(finalDefinition);
  }, [mode, yaml, definition, onRun]);

  // Undo/Redo
  const handleUndo = useCallback(() => {
    if (historyIndex > 0) {
      setHistoryIndex((prev) => prev - 1);
      setDefinition(history[historyIndex - 1]);
    }
  }, [history, historyIndex]);

  const handleRedo = useCallback(() => {
    if (historyIndex < history.length - 1) {
      setHistoryIndex((prev) => prev + 1);
      setDefinition(history[historyIndex + 1]);
    }
  }, [history, historyIndex]);

  const canUndo = historyIndex > 0;
  const canRedo = historyIndex < history.length - 1;

  return (
    <div className={cn('flex flex-col h-full bg-background', className)}>
      {/* Toolbar */}
      <div className="flex items-center justify-between px-4 py-2 border-b bg-muted/30">
        <div className="flex items-center gap-2">
          {/* Mode Tabs */}
          <Tabs value={mode} onValueChange={handleModeChange}>
            <TabsList className="h-8">
              <TabsTrigger value="visual" className="text-xs px-3 h-7">
                <LayoutGrid className="w-3.5 h-3.5 mr-1.5" />
                Visual
              </TabsTrigger>
              <TabsTrigger value="yaml" className="text-xs px-3 h-7">
                <Code className="w-3.5 h-3.5 mr-1.5" />
                YAML
              </TabsTrigger>
            </TabsList>
          </Tabs>

          <div className="w-px h-6 bg-border mx-2" />

          {/* Undo/Redo */}
          <Button
            variant="ghost"
            size="sm"
            onClick={handleUndo}
            disabled={!canUndo}
            className="h-8 w-8 p-0"
          >
            <Undo className="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={handleRedo}
            disabled={!canRedo}
            className="h-8 w-8 p-0"
          >
            <Redo className="w-4 h-4" />
          </Button>

          {/* Panel toggles (visual mode) */}
          {mode === 'visual' && (
            <>
              <div className="w-px h-6 bg-border mx-2" />
              <Button
                variant={showPalette ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => setShowPalette(!showPalette)}
                className="h-8 text-xs"
              >
                <ChevronLeft className={cn('w-4 h-4 mr-1 transition-transform', !showPalette && 'rotate-180')} />
                Actions
              </Button>
              <Button
                variant={showProperties ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => setShowProperties(!showProperties)}
                className="h-8 text-xs"
              >
                Properties
                <ChevronRight className={cn('w-4 h-4 ml-1 transition-transform', !showProperties && 'rotate-180')} />
              </Button>
            </>
          )}

          {/* Validate button (yaml mode) */}
          {mode === 'yaml' && (
            <>
              <div className="w-px h-6 bg-border mx-2" />
              <Button
                variant="outline"
                size="sm"
                onClick={handleValidate}
                className="h-8 text-xs"
              >
                <CheckCircle2 className="w-3.5 h-3.5 mr-1.5" />
                Validate
              </Button>
            </>
          )}
        </div>

        <div className="flex items-center gap-2">
          {isDirty && (
            <span className="text-xs text-muted-foreground">Unsaved changes</span>
          )}

          <Button
            variant="outline"
            size="sm"
            onClick={handleSave}
            disabled={isSaving}
            className="h-8"
          >
            <Save className="w-4 h-4 mr-1.5" />
            {isSaving ? 'Saving...' : 'Save'}
          </Button>

          {onRun && (
            <Button
              size="sm"
              onClick={handleRun}
              disabled={isRunning}
              className="h-8"
            >
              <Play className="w-4 h-4 mr-1.5" />
              {isRunning ? 'Running...' : 'Run'}
            </Button>
          )}
        </div>
      </div>

      {/* Error Alert */}
      {yamlError && (
        <Alert variant="destructive" className="mx-4 mt-2">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription className="whitespace-pre-wrap">{yamlError}</AlertDescription>
        </Alert>
      )}

      {/* Success Alert */}
      {validationSuccess && !yamlError && (
        <Alert className="mx-4 mt-2 border-green-500/50 text-green-700 dark:text-green-400">
          <CheckCircle2 className="h-4 w-4 text-green-600" />
          <AlertDescription>{validationSuccess}</AlertDescription>
        </Alert>
      )}

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {mode === 'visual' ? (
          <>
            {/* Node Palette */}
            {showPalette && <NodePalette />}

            {/* Canvas */}
            <div className="flex-1 relative">
              <FlowCanvas
                definition={definition}
                onDefinitionChange={handleDefinitionChange}
                onNodeSelect={handleNodeSelect}
                selectedNodeId={selectedNode?.id}
              />
            </div>

            {/* Properties Panel */}
            {showProperties && (
              <PropertiesPanel
                node={selectedNode}
                onNodeUpdate={handleNodeUpdate}
                onClose={() => setShowProperties(false)}
              />
            )}
          </>
        ) : (
          /* YAML Editor */
          <div className="flex-1 p-4">
            <Textarea
              value={yaml}
              onChange={(e) => handleYamlChange(e.target.value)}
              placeholder="Enter your flow definition in YAML format..."
              className="w-full h-full font-mono text-sm resize-none"
            />
          </div>
        )}
      </div>
    </div>
  );
}
