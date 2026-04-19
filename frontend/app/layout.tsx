import "./globals.css";
import type { ReactNode } from "react";
import { SidebarNav } from "@/components/layout/sidebar-nav";

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <div className="min-h-screen grid grid-cols-[280px_1fr]">
          <SidebarNav />
          <main className="p-6 lg:p-8">{children}</main>
        </div>
      </body>
    </html>
  );
}
