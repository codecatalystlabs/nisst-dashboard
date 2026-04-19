import Link from "next/link";
import type { Route } from "next";
import { serviceDeliveryTabs } from "@/lib/service-delivery";

export default function ServiceDeliveryPage() {
  return (
    <div className="space-y-5">
      <h2 className="text-2xl font-semibold text-navy">Service Delivery</h2>
      <p className="text-slate-600">Select a Service Delivery sub-section to view scoped indicators and comparisons.</p>
      <div className="flex flex-wrap gap-2">
        {serviceDeliveryTabs.map((tab) => (
          <Link
            key={tab.slug}
            href={`/domains/service-delivery/${tab.slug}` as Route}
            className="rounded-full bg-white px-4 py-2 shadow-midas text-sm hover:bg-slate-50"
          >
            {tab.label}
          </Link>
        ))}
      </div>
    </div>
  );
}
