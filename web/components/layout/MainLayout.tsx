'use client';

import { Sidebar } from './Sidebar';
import { usePathname } from 'next/navigation';

// Pages that should NOT show the sidebar (e.g., login)
const noSidebarPaths = ['/login', '/login/callback'];

interface MainLayoutProps {
  children: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  const pathname = usePathname();
  const showSidebar = !noSidebarPaths.some(path => pathname.startsWith(path));

  if (!showSidebar) {
    return <>{children}</>;
  }

  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar />
      <main className="flex-1 overflow-auto px-4">
        {children}
      </main>
    </div>
  );
}
