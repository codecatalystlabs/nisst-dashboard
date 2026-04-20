import "./globals.css";
import { Suspense, type ReactNode } from "react";
import { SidebarNav } from "@/components/layout/sidebar-nav";
import { GlobalFilterBar } from "@/components/layout/global-filter-bar";

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <div className="min-h-screen grid grid-cols-[280px_1fr]">
          <Suspense fallback={<aside className="bg-navy text-white p-6 sticky top-0 h-screen" />}>
            <SidebarNav />
          </Suspense>
          <main className="p-6 lg:p-8">
            <Suspense fallback={null}>
              <GlobalFilterBar />
            </Suspense>
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
