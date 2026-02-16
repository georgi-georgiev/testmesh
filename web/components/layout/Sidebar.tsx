'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import {
  LayoutDashboard,
  FileText,
  Play,
  FolderTree,
  Server,
  FileCode,
  Calendar,
  Gauge,
  BarChart3,
  FileBarChart,
  Sparkles,
  Puzzle,
  Activity,
  Building2,
  Upload,
  Terminal,
  History,
  ChevronLeft,
  ChevronRight,
  Settings,
  Globe,
  Plug,
  Users,
} from 'lucide-react';
import { EnvironmentSelector } from '@/components/environments/EnvironmentSelector';
import { WorkspaceSwitcher } from '@/components/workspaces/WorkspaceSwitcher';
import { useState } from 'react';

interface NavItem {
  title: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  badge?: string;
}

interface NavSection {
  title: string;
  items: NavItem[];
}

const navigation: NavSection[] = [
  {
    title: 'Overview',
    items: [
      { title: 'Dashboard', href: '/', icon: LayoutDashboard },
      { title: 'Activity', href: '/activity', icon: Activity },
    ],
  },
  {
    title: 'Admin',
    items: [
      { title: 'Dashboard', href: '/admin', icon: Settings },
      { title: 'Users', href: '/admin/users', icon: Users },
      { title: 'Integrations', href: '/admin/integrations', icon: Plug },
      { title: 'Health', href: '/admin/health', icon: Activity },
    ],
  },
  {
    title: 'Testing',
    items: [
      { title: 'Flows', href: '/flows', icon: FileText },
      { title: 'Executions', href: '/executions', icon: Play },
      { title: 'Collections', href: '/collections', icon: FolderTree },
      { title: 'Schedules', href: '/schedules', icon: Calendar },
      { title: 'Runner', href: '/runner', icon: Terminal },
    ],
  },
  {
    title: 'Infrastructure',
    items: [
      { title: 'Mock Servers', href: '/mocks', icon: Server },
      { title: 'Contracts', href: '/contracts', icon: FileCode },
      { title: 'Load Testing', href: '/load-testing', icon: Gauge },
    ],
  },
  {
    title: 'Insights',
    items: [
      { title: 'Analytics', href: '/analytics', icon: BarChart3 },
      { title: 'Reports', href: '/reports', icon: FileBarChart },
      { title: 'History', href: '/history', icon: History },
    ],
  },
  {
    title: 'AI & Integrations',
    items: [
      { title: 'AI Features', href: '/ai', icon: Sparkles },
      { title: 'Plugins', href: '/plugins', icon: Puzzle },
      { title: 'Import', href: '/import', icon: Upload },
    ],
  },
  {
    title: 'Settings',
    items: [
      { title: 'Environments', href: '/environments', icon: Globe },
      { title: 'Workspaces', href: '/workspaces', icon: Building2 },
    ],
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const [collapsed, setCollapsed] = useState(false);

  const isActive = (href: string) => {
    if (href === '/') {
      return pathname === '/';
    }
    return pathname.startsWith(href);
  };

  return (
    <div
      className={cn(
        'flex flex-col h-full border-r bg-background transition-all duration-300',
        collapsed ? 'w-16' : 'w-64'
      )}
    >
      {/* Header */}
      <div className="flex h-14 items-center border-b px-4">
        {!collapsed && (
          <Link href="/" className="flex items-center gap-2 font-semibold">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
              TM
            </div>
            <span>TestMesh</span>
          </Link>
        )}
        {collapsed && (
          <Link href="/" className="flex items-center justify-center w-full">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
              TM
            </div>
          </Link>
        )}
      </div>

      {/* Workspace Switcher */}
      {!collapsed && (
        <div className="px-3 py-2 border-b">
          <WorkspaceSwitcher className="w-full" />
        </div>
      )}
      {collapsed && (
        <div className="px-2 py-2 border-b flex justify-center">
          <WorkspaceSwitcher compact />
        </div>
      )}

      {/* Environment Selector */}
      {!collapsed && (
        <div className="px-3 py-2 border-b">
          <EnvironmentSelector className="w-full" />
        </div>
      )}
      {collapsed && (
        <div className="px-2 py-2 border-b flex justify-center">
          <EnvironmentSelector compact />
        </div>
      )}

      {/* Navigation */}
      <ScrollArea className="flex-1 min-h-0 overflow-hidden">
        <nav className="space-y-1 px-2 py-2">
          {navigation.map((section, sectionIdx) => (
            <div key={section.title}>
              {sectionIdx > 0 && <Separator className="my-2" />}
              {!collapsed && (
                <h4 className="mb-1 px-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                  {section.title}
                </h4>
              )}
              <div className="space-y-1">
                {section.items.map((item) => {
                  const Icon = item.icon;
                  const active = isActive(item.href);
                  return (
                    <Link key={item.href} href={item.href}>
                      <Button
                        variant={active ? 'secondary' : 'ghost'}
                        className={cn(
                          'w-full',
                          collapsed ? 'justify-center px-2' : 'justify-start',
                          active && 'bg-secondary'
                        )}
                        title={collapsed ? item.title : undefined}
                      >
                        <Icon className={cn('h-4 w-4', !collapsed && 'mr-2')} />
                        {!collapsed && <span>{item.title}</span>}
                        {!collapsed && item.badge && (
                          <span className="ml-auto text-xs bg-primary/10 text-primary px-1.5 py-0.5 rounded">
                            {item.badge}
                          </span>
                        )}
                      </Button>
                    </Link>
                  );
                })}
              </div>
            </div>
          ))}
        </nav>
      </ScrollArea>

      {/* Collapse Toggle */}
      <div className="border-t p-2">
        <Button
          variant="ghost"
          size="sm"
          className="w-full justify-center"
          onClick={() => setCollapsed(!collapsed)}
        >
          {collapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <>
              <ChevronLeft className="h-4 w-4 mr-2" />
              <span>Collapse</span>
            </>
          )}
        </Button>
      </div>
    </div>
  );
}
