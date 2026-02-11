'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useFlows, useDeleteFlow } from '@/lib/hooks/useFlows';
import { useCreateExecution } from '@/lib/hooks/useExecutions';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Play, Trash2, Eye, Plus } from 'lucide-react';

export default function FlowsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [suiteFilter, setSuiteFilter] = useState('');

  const { data, isLoading, error } = useFlows({
    suite: suiteFilter || undefined,
  });

  const deleteFlow = useDeleteFlow();
  const createExecution = useCreateExecution();

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this flow?')) {
      deleteFlow.mutate(id);
    }
  };

  const handleRun = async (flowId: string) => {
    createExecution.mutate({
      flow_id: flowId,
      environment: 'development',
    });
  };

  const flows = data?.flows || [];
  const filteredFlows = flows.filter((flow) => {
    const matchesSearch =
      flow.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      flow.description?.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  });

  if (error) {
    return (
      <div className="container mx-auto py-8">
        <Card>
          <CardHeader>
            <CardTitle>Error Loading Flows</CardTitle>
            <CardDescription>
              {error instanceof Error ? error.message : 'An error occurred'}
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold">Test Flows</h1>
          <p className="text-muted-foreground mt-1">
            Manage and execute your test flows
          </p>
        </div>
        <Link href="/flows/new">
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            Create Flow
          </Button>
        </Link>
      </div>

      <Card className="mb-6">
        <CardContent className="pt-6">
          <div className="flex gap-4">
            <Input
              placeholder="Search flows..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="flex-1"
            />
            <Input
              placeholder="Filter by suite..."
              value={suiteFilter}
              onChange={(e) => setSuiteFilter(e.target.value)}
              className="w-64"
            />
          </div>
        </CardContent>
      </Card>

      {isLoading ? (
        <Card>
          <CardContent className="py-12">
            <div className="text-center text-muted-foreground">
              Loading flows...
            </div>
          </CardContent>
        </Card>
      ) : filteredFlows.length === 0 ? (
        <Card>
          <CardContent className="py-12">
            <div className="text-center">
              <p className="text-muted-foreground mb-4">
                {flows.length === 0 ? 'No flows found' : 'No flows match your search'}
              </p>
              <Link href="/flows/new">
                <Button>
                  <Plus className="w-4 h-4 mr-2" />
                  Create Your First Flow
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Suite</TableHead>
                <TableHead>Tags</TableHead>
                <TableHead>Steps</TableHead>
                <TableHead>Environment</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredFlows.map((flow) => (
                <TableRow key={flow.id}>
                  <TableCell>
                    <Link
                      href={`/flows/${flow.id}`}
                      className="font-medium hover:underline"
                    >
                      {flow.name}
                    </Link>
                    {flow.description && (
                      <p className="text-sm text-muted-foreground mt-1">
                        {flow.description}
                      </p>
                    )}
                  </TableCell>
                  <TableCell>
                    {flow.suite && (
                      <Badge variant="outline">{flow.suite}</Badge>
                    )}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-1 flex-wrap">
                      {flow.tags?.map((tag) => (
                        <Badge key={tag} variant="secondary">
                          {tag}
                        </Badge>
                      ))}
                    </div>
                  </TableCell>
                  <TableCell>{flow.definition.steps?.length || 0}</TableCell>
                  <TableCell>
                    <Badge variant="outline">{flow.environment}</Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleRun(flow.id)}
                        disabled={createExecution.isPending}
                      >
                        <Play className="w-4 h-4" />
                      </Button>
                      <Link href={`/flows/${flow.id}`}>
                        <Button variant="ghost" size="sm">
                          <Eye className="w-4 h-4" />
                        </Button>
                      </Link>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDelete(flow.id)}
                        disabled={deleteFlow.isPending}
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      )}

      {data && data.total > 0 && (
        <div className="mt-4 text-sm text-muted-foreground text-center">
          Showing {filteredFlows.length} of {data.total} flows
        </div>
      )}
    </div>
  );
}
