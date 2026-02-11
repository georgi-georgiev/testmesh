import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export default function Page() {
  return (
    <div className="container mx-auto py-16">
      <div className="max-w-4xl mx-auto text-center space-y-8">
        <div className="space-y-4">
          <h1 className="text-5xl font-bold">TestMesh</h1>
          <p className="text-xl text-muted-foreground">
            E2E Integration Testing Platform
          </p>
          <p className="text-muted-foreground">
            Write tests in YAML and execute them across multiple protocols
          </p>
        </div>

        <div className="flex gap-4 justify-center">
          <Link href="/flows">
            <Button size="lg">View Flows</Button>
          </Link>
          <Link href="/flows/new">
            <Button size="lg" variant="outline">Create Flow</Button>
          </Link>
        </div>

        <div className="grid md:grid-cols-3 gap-6 mt-12">
          <Card>
            <CardHeader>
              <CardTitle>HTTP Requests</CardTitle>
              <CardDescription>
                Test REST APIs with GET, POST, PUT, DELETE
              </CardDescription>
            </CardHeader>
            <CardContent>
              Full support for headers, body, and assertions
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Database Queries</CardTitle>
              <CardDescription>
                Execute SQL queries against PostgreSQL
              </CardDescription>
            </CardHeader>
            <CardContent>
              INSERT, SELECT, UPDATE, DELETE with parameterized queries
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Assertions</CardTitle>
              <CardDescription>
                Validate responses with expressions
              </CardDescription>
            </CardHeader>
            <CardContent>
              JSONPath, boolean expressions, and status code checks
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
