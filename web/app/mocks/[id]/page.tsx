'use client';

import { useState } from 'react';
import { use } from 'react';
import Link from 'next/link';
import {
  useMockServer,
  useMockServerEndpoints,
  useMockServerRequests,
  useMockServerStates,
} from '@/lib/hooks/useMockServers';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Server, ArrowLeft, ExternalLink, Clock, Database } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';
import type { MockServerStatus } from '@/lib/api/types';

type PageParams = Promise<{ id: string }>;

export default function MockServerDetailPage({ params }: { params: PageParams }) {
  const resolvedParams = use(params);
  const serverId = resolvedParams.id;
  const [activeTab, setActiveTab] = useState<'endpoints' | 'requests' | 'state'>('endpoints');

  const { data: server, isLoading: serverLoading, error: serverError } = useMockServer(serverId);
  const { data: endpointsData } = useMockServerEndpoints(serverId);
  const { data: requestsData } = useMockServerRequests(serverId, { limit: 50 });
  const { data: statesData } = useMockServerStates(serverId);

  const getStatusBadge = (status: MockServerStatus) => {
    const colors = {
      starting: 'bg-yellow-500',
      running: 'bg-green-500',
      stopped: 'bg-gray-500',
      failed: 'bg-red-500',
    };

    return (
      <Badge variant="default" className="capitalize">
        <span className={`w-2 h-2 rounded-full ${colors[status]} mr-2`} />
        {status}
      </Badge>
    );
  };

  if (serverError) {
    return (
      <div className="container mx-auto py-8">
        <Card>
          <CardHeader>
            <CardTitle>Error Loading Mock Server</CardTitle>
            <CardDescription>
              {serverError instanceof Error ? serverError.message : 'An error occurred'}
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    );
  }

  if (serverLoading || !server) {
    return (
      <div className="container mx-auto py-8">
        <div className="text-center text-muted-foreground">Loading mock server...</div>
      </div>
    );
  }

  const endpoints = endpointsData?.endpoints || [];
  const requests = requestsData?.requests || [];
  const states = statesData?.states || [];

  return (
    <div className="container mx-auto py-8">
      <div className="mb-6">
        <Link href="/mocks">
          <Button variant="ghost" size="sm" className="mb-4">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Mock Servers
          </Button>
        </Link>

        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold flex items-center gap-2">
              <Server className="w-8 h-8" />
              {server.name}
            </h1>
            <div className="flex items-center gap-4 mt-2">
              <code className="text-sm bg-muted px-2 py-1 rounded">{server.base_url}</code>
              {server.status === 'running' && (
                <a
                  href={server.base_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-muted-foreground hover:text-foreground flex items-center gap-1"
                >
                  Open in browser <ExternalLink className="w-3 h-3" />
                </a>
              )}
            </div>
          </div>
          <div className="text-right">
            {getStatusBadge(server.status)}
            <p className="text-sm text-muted-foreground mt-2">Port: {server.port}</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium">Endpoints</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{endpoints.length}</div>
            <p className="text-xs text-muted-foreground mt-1">Configured endpoints</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium">Requests</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{requests.length}</div>
            <p className="text-xs text-muted-foreground mt-1">Total requests received</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium">Uptime</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {server.started_at
                ? formatDistanceToNow(new Date(server.started_at))
                : '-'}
            </div>
            <p className="text-xs text-muted-foreground mt-1">Since start</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <div className="flex space-x-4 border-b">
            <button
              onClick={() => setActiveTab('endpoints')}
              className={`pb-2 px-1 ${
                activeTab === 'endpoints'
                  ? 'border-b-2 border-primary font-medium'
                  : 'text-muted-foreground'
              }`}
            >
              Endpoints ({endpoints.length})
            </button>
            <button
              onClick={() => setActiveTab('requests')}
              className={`pb-2 px-1 ${
                activeTab === 'requests'
                  ? 'border-b-2 border-primary font-medium'
                  : 'text-muted-foreground'
              }`}
            >
              Requests ({requests.length})
            </button>
            <button
              onClick={() => setActiveTab('state')}
              className={`pb-2 px-1 ${
                activeTab === 'state'
                  ? 'border-b-2 border-primary font-medium'
                  : 'text-muted-foreground'
              }`}
            >
              State ({states.length})
            </button>
          </div>
        </CardHeader>
        <CardContent>
          {activeTab === 'endpoints' && (
            <div>
              {endpoints.length === 0 ? (
                <p className="text-center text-muted-foreground py-8">No endpoints configured</p>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Method</TableHead>
                      <TableHead>Path</TableHead>
                      <TableHead>Response Status</TableHead>
                      <TableHead>Priority</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {endpoints.map((endpoint) => (
                      <TableRow key={endpoint.id}>
                        <TableCell>
                          <Badge variant="outline">{endpoint.method}</Badge>
                        </TableCell>
                        <TableCell>
                          <code className="text-sm">{endpoint.path}</code>
                        </TableCell>
                        <TableCell>{endpoint.response_config.status_code}</TableCell>
                        <TableCell>{endpoint.priority}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </div>
          )}

          {activeTab === 'requests' && (
            <div>
              {requests.length === 0 ? (
                <p className="text-center text-muted-foreground py-8">No requests received yet</p>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Time</TableHead>
                      <TableHead>Method</TableHead>
                      <TableHead>Path</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Matched</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {requests.map((request) => (
                      <TableRow key={request.id}>
                        <TableCell>
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <Clock className="w-3 h-3" />
                            {formatDistanceToNow(new Date(request.received_at), {
                              addSuffix: true,
                            })}
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline">{request.method}</Badge>
                        </TableCell>
                        <TableCell>
                          <code className="text-sm">{request.path}</code>
                        </TableCell>
                        <TableCell>{request.response_code}</TableCell>
                        <TableCell>
                          {request.matched ? (
                            <Badge variant="default" className="bg-green-500">
                              Matched
                            </Badge>
                          ) : (
                            <Badge variant="secondary">Unmatched</Badge>
                          )}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </div>
          )}

          {activeTab === 'state' && (
            <div>
              {states.length === 0 ? (
                <p className="text-center text-muted-foreground py-8">No state variables</p>
              ) : (
                <div className="space-y-4">
                  {states.map((state) => (
                    <Card key={state.id}>
                      <CardHeader>
                        <CardTitle className="text-sm flex items-center gap-2">
                          <Database className="w-4 h-4" />
                          {state.state_key}
                        </CardTitle>
                        <CardDescription className="text-xs">
                          Updated {formatDistanceToNow(new Date(state.updated_at), { addSuffix: true })}
                        </CardDescription>
                      </CardHeader>
                      <CardContent>
                        <pre className="text-sm bg-muted p-3 rounded overflow-auto">
                          {JSON.stringify(state.state_value, null, 2)}
                        </pre>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
