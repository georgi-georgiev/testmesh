'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useCreateFlow } from '@/lib/hooks/useFlows';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Textarea } from '@/components/ui/textarea';
import { ArrowLeft, Save } from 'lucide-react';
import { Alert, AlertDescription } from '@/components/ui/alert';

const EXAMPLE_FLOW = `name: "Example Flow"
description: "A sample test flow"
suite: "smoke-tests"
tags: ["api", "smoke"]

env:
  BASE_URL: "https://jsonplaceholder.typicode.com"

steps:
  - id: get_user
    action: http_request
    name: "Get user details"
    config:
      method: GET
      url: "\${BASE_URL}/users/1"
    assert:
      - status == 200
      - body.name != null
    output:
      user_id: "$.id"
      user_name: "$.name"

  - id: get_posts
    action: http_request
    name: "Get user posts"
    config:
      method: GET
      url: "\${BASE_URL}/users/\${get_user.user_id}/posts"
    assert:
      - status == 200`;

export default function NewFlowPage() {
  const router = useRouter();
  const [yaml, setYaml] = useState('');
  const [error, setError] = useState<string | null>(null);
  const createFlow = useCreateFlow();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!yaml.trim()) {
      setError('Please enter a flow definition');
      return;
    }

    try {
      const flow = await createFlow.mutateAsync({ yaml });
      router.push(`/flows/${flow.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create flow');
    }
  };

  const loadExample = () => {
    setYaml(EXAMPLE_FLOW);
  };

  return (
    <div className="container mx-auto py-8 max-w-4xl">
      <div className="mb-6">
        <Link href="/flows">
          <Button variant="ghost" size="sm" className="mb-4">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Flows
          </Button>
        </Link>

        <h1 className="text-3xl font-bold">Create New Flow</h1>
        <p className="text-muted-foreground mt-2">
          Define your test flow using YAML syntax
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
              <div>
                <CardTitle>Flow Definition</CardTitle>
                <CardDescription>
                  Write your flow in YAML format
                </CardDescription>
              </div>
              <Button type="button" variant="outline" onClick={loadExample}>
                Load Example
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {error && (
              <Alert variant="destructive" className="mb-4">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <Textarea
              value={yaml}
              onChange={(e) => setYaml(e.target.value)}
              placeholder={EXAMPLE_FLOW}
              className="font-mono text-sm min-h-[500px]"
            />

            <div className="mt-4 p-4 bg-muted rounded-lg">
              <h4 className="font-medium mb-2">Quick Reference</h4>
              <ul className="text-sm text-muted-foreground space-y-1">
                <li>• <strong>name:</strong> Flow name (required)</li>
                <li>• <strong>description:</strong> Flow description</li>
                <li>• <strong>suite:</strong> Test suite name</li>
                <li>• <strong>tags:</strong> Array of tags</li>
                <li>• <strong>env:</strong> Environment variables</li>
                <li>• <strong>steps:</strong> Array of test steps (required)</li>
                <li>• <strong>action:</strong> http_request, database_query</li>
                <li>• <strong>assert:</strong> Array of assertions</li>
                <li>• <strong>output:</strong> Extract values for later use</li>
              </ul>
            </div>
          </CardContent>
        </Card>

        <div className="flex justify-end gap-4">
          <Link href="/flows">
            <Button type="button" variant="outline">
              Cancel
            </Button>
          </Link>
          <Button type="submit" disabled={createFlow.isPending}>
            <Save className="w-4 h-4 mr-2" />
            {createFlow.isPending ? 'Creating...' : 'Create Flow'}
          </Button>
        </div>
      </form>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Supported Actions</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <h4 className="font-medium mb-2">HTTP Request</h4>
            <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
              {`- action: http_request
  config:
    method: GET|POST|PUT|DELETE
    url: "https://api.example.com/endpoint"
    headers:
      Content-Type: application/json
    body:
      key: value
  assert:
    - status == 200
    - body.data != null`}
            </pre>
          </div>

          <div>
            <h4 className="font-medium mb-2">Database Query</h4>
            <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
              {`- action: database_query
  config:
    connection: "postgresql://user:pass@host:5432/db"
    query: "SELECT * FROM users WHERE id = $1"
    params: [1]
  assert:
    - row_count > 0
    - first_row.name == "John"`}
            </pre>
          </div>

          <div>
            <h4 className="font-medium mb-2">Variables</h4>
            <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
              {`# Use environment variables
url: "\${BASE_URL}/endpoint"

# Use step outputs
url: "/users/\${get_user.user_id}"

# Built-in variables
id: "\${RANDOM_ID}"
timestamp: "\${TIMESTAMP}"`}
            </pre>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
