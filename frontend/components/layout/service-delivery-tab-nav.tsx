"use client";

import Link from "next/link";
import type { Route } from "next";
import { useSearchParams } from "next/navigation";
import { serviceDeliveryTabs, type ServiceDeliverySubtab } from "@/lib/service-delivery";

export function ServiceDeliveryTabNav({ activeSlug }: { activeSlug?: ServiceDeliverySubtab }) {
  const searchParams = useSearchParams();
  const query = searchParams.toString();

  return (
    <div className="flex flex-wrap gap-2">
      {serviceDeliveryTabs.map((item) => (
        <Link
          key={item.slug}
          href={`/domains/service-delivery/${item.slug}${query ? `?${query}` : ""}` as Route}
          className={`rounded-full px-4 py-2 shadow-midas text-sm ${
            item.slug === activeSlug ? "bg-navy text-white" : "bg-white hover:bg-slate-50"
          }`}
        >
          {item.label}
        </Link>
      ))}
    </div>
  );
}
