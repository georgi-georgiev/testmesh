'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useMockServers, useDeleteMockServer } from '@/lib/hooks/useMockServers';
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
import { Eye, Trash2, X, Search, Server, ExternalLink } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';
import type { MockServerStatus } from '@/lib/api/types';

export default function MockServersPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<MockServerStatus | ''>('');

  const { data, isLoading, error } = useMockServers({
    status: statusFilter || undefined,
  });

  const deleteMockServer = useDeleteMockServer();

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this mock server?')) {
      deleteMockServer.mutate(id);
    }
  };

  const servers = data?.servers || [];
  const filteredServers = servers.filter((server) => {
    const matchesSearch =
      server.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      server.base_url.toLowerCase().includes(searchQuery.toLowerCase());

    return matchesSearch;
  });

  const hasActiveFilters = searchQuery || statusFilter;
  const clearFilters = () => {
    setSearchQuery('');
    setStatusFilter('');
  };

  const getStatusBadge = (status: MockServerStatus) => {
    const variants = {
      starting: 'default',
      running: 'default',
      stopped: 'secondary',
      failed: 'destructive',
    };

    const colors = {
      starting: 'bg-yellow-500',
      running: 'bg-green-500',
      stopped: 'bg-gray-500',
      failed: 'bg-red-500',
    };

    return (
      <Badge variant={variants[status] as any} className="capitalize">
        <span className={`w-2 h-2 rounded-full ${colors[status]} mr-2`} />
        {status}
      </Badge>
    );
  };

  if (error) {
    return (
      <div className="container mx-auto py-8">
        <Card>
          <CardHeader>
            <CardTitle>Error Loading Mock Servers</CardTitle>
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
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Server className="w-8 h-8" />
            Mock Servers
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage mock API servers for isolated testing
          </p>
        </div>
      </div>

      <Card className="mb-6">
        <CardContent className="pt-6">
          <div className="space-y-4">
            <div className="flex gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  placeholder="Search mock servers by name or URL..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value as MockServerStatus | '')}
                className="px-3 py-2 border rounded-md bg-background"
              >
                <option value="">All Statuses</option>
                <option value="running">Running</option>
                <option value="stopped">Stopped</option>
                <option value="starting">Starting</option>
                <option value="failed">Failed</option>
              </select>
              {hasActiveFilters && (
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={clearFilters}
                  title="Clear filters"
                >
                  <X className="w-4 h-4" />
                </Button>
              )}
            </div>

            {hasActiveFilters && (
              <div className="flex gap-2 items-center text-sm text-muted-foreground">
                <span>Active filters:</span>
                {searchQuery && (
                  <Badge variant="secondary" className="gap-1">
                    Search: {searchQuery}
                    <button
                      onClick={() => setSearchQuery('')}
                      className="ml-1 hover:text-foreground"
                    >
                      <X className="w-3 h-3" />
                    </button>
                  </Badge>
                )}
                {statusFilter && (
                  <Badge variant="secondary" className="gap-1">
                    Status: {statusFilter}
                    <button
                      onClick={() => setStatusFilter('')}
                      className="ml-1 hover:text-foreground"
                    >
                      <X className="w-3 h-3" />
                    </button>
                  </Badge>
                )}
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {isLoading ? (
        <Card>
          <CardContent className="py-12">
            <div className="text-center text-muted-foreground">
              Loading mock servers...
            </div>
          </CardContent>
        </Card>
      ) : filteredServers.length === 0 ? (
        <Card>
          <CardContent className="py-12">
            <div className="text-center">
              <Server className="w-12 h-12 mx-auto mb-4 text-muted-foreground" />
              <p className="text-muted-foreground mb-4">
                {servers.length === 0
                  ? 'No mock servers found'
                  : 'No mock servers match your search'}
              </p>
              <p className="text-sm text-muted-foreground">
                Mock servers are created automatically when you run flows with{' '}
                <code className="bg-muted px-1 py-0.5 rounded">mock_server_start</code> actions
              </p>
            </div>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Base URL</TableHead>
                <TableHead>Port</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Uptime</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredServers.map((server) => (
                <TableRow key={server.id}>
                  <TableCell>
                    <Link
                      href={`/mocks/${server.id}`}
                      className="font-medium hover:underline"
                    >
                      {server.name}
                    </Link>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <code className="text-sm bg-muted px-2 py-1 rounded">
                        {server.base_url}
                      </code>
                      {server.status === 'running' && (
                        <a
                          href={server.base_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-muted-foreground hover:text-foreground"
                        >
                          <ExternalLink className="w-3 h-3" />
                        </a>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>{server.port}</TableCell>
                  <TableCell>{getStatusBadge(server.status)}</TableCell>
                  <TableCell>
                    {server.started_at
                      ? formatDistanceToNow(new Date(server.started_at), {
                          addSuffix: true,
                        })
                      : '-'}
                  </TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Link href={`/mocks/${server.id}`}>
                        <Button variant="ghost" size="sm">
                          <Eye className="w-4 h-4" />
                        </Button>
                      </Link>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDelete(server.id)}
                        disabled={deleteMockServer.isPending}
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
          Showing {filteredServers.length} of {data.total} mock servers
        </div>
      )}
    </div>
  );
}
