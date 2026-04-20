import { Suspense } from "react";
import { ServiceDeliveryTabNav } from "@/components/layout/service-delivery-tab-nav";

export default function ServiceDeliveryPage() {
  return (
    <div className="space-y-5">
      <h2 className="text-2xl font-semibold text-navy">Service Delivery</h2>
      <p className="text-slate-600">Select a Service Delivery sub-section to view scoped indicators and comparisons.</p>
      <Suspense fallback={null}>
        <ServiceDeliveryTabNav />
      </Suspense>
    </div>
  );
}
