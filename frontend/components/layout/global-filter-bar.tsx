"use client";

import { useEffect, useMemo, useState } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import type { Route } from "next";
import { FILTER_KEYS } from "@/lib/filters";
import { apiGet } from "@/lib/api";

type FilterState = Record<string, string>;
type OptionItem = { value: string };

const STORAGE_KEY = "nisst-global-filters";

function readFromSearch(sp: URLSearchParams): FilterState {
  const out: FilterState = {};
  for (const k of FILTER_KEYS) out[k] = sp.get(k) ?? "";
  return out;
}

export function GlobalFilterBar() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [values, setValues] = useState<FilterState>(() => readFromSearch(new URLSearchParams()));
  const [options, setOptions] = useState<Record<string, string[]>>({
    period: [],
    region: [],
    district: [],
    facility: [],
    level: [],
  });

  const activeCount = useMemo(() => FILTER_KEYS.filter((k) => values[k]?.trim()).length, [values]);

  useEffect(() => {
    setValues(readFromSearch(searchParams));
  }, [searchParams]);

  useEffect(() => {
    const timer = window.setTimeout(async () => {
      const regionQ = encodeURIComponent(values.region ?? "");
      const districtQ = encodeURIComponent(values.district ?? "");
      const facilityQ = encodeURIComponent(values.facility ?? "");
      const periodQ = encodeURIComponent(values.period ?? "");
      const levelQ = encodeURIComponent(values.level ?? "");

      const [regions, districts, facilities, periods, levels] = await Promise.allSettled([
        apiGet<{ items: OptionItem[] }>(`/metadata/regions?q=${regionQ}&limit=20`),
        apiGet<{ items: OptionItem[] }>(
          `/metadata/districts?q=${districtQ}&region=${encodeURIComponent(values.region ?? "")}&limit=20`,
        ),
        apiGet<{ items: OptionItem[] }>(
          `/metadata/facilities?q=${facilityQ}&region=${encodeURIComponent(values.region ?? "")}&district=${encodeURIComponent(values.district ?? "")}&limit=20`,
        ),
        apiGet<{ items: OptionItem[] }>(`/metadata/periods?q=${periodQ}&limit=20`),
        apiGet<{ items: OptionItem[] }>(`/metadata/levels?q=${levelQ}&limit=20`),
      ]);

      setOptions({
        region: regions.status === "fulfilled" ? regions.value.items.map((i) => i.value) : [],
        district: districts.status === "fulfilled" ? districts.value.items.map((i) => i.value) : [],
        facility: facilities.status === "fulfilled" ? facilities.value.items.map((i) => i.value) : [],
        period: periods.status === "fulfilled" ? periods.value.items.map((i) => i.value) : [],
        level: levels.status === "fulfilled" ? levels.value.items.map((i) => i.value) : [],
      });
    }, 220);

    return () => window.clearTimeout(timer);
  }, [values.region, values.district, values.facility, values.period, values.level]);

  useEffect(() => {
    const hasQueryFilters = FILTER_KEYS.some((k) => (searchParams.get(k) ?? "").trim() !== "");
    if (!hasQueryFilters) {
      const raw = typeof window !== "undefined" ? window.localStorage.getItem(STORAGE_KEY) : null;
      if (!raw) return;
      try {
        const saved = JSON.parse(raw) as FilterState;
        const params = new URLSearchParams(searchParams.toString());
        let changed = false;
        for (const k of FILTER_KEYS) {
          const v = saved[k];
          if (v && v.trim()) {
            params.set(k, v.trim());
            changed = true;
          }
        }
        if (changed) router.replace(`${pathname}?${params.toString()}` as Route);
      } catch {
        // ignore invalid local storage payload
      }
    }
  }, [pathname, router, searchParams]);

  function applyFilters() {
    const params = new URLSearchParams(searchParams.toString());
    for (const k of FILTER_KEYS) {
      const v = values[k]?.trim() ?? "";
      if (v) params.set(k, v);
      else params.delete(k);
    }
    const next = `${pathname}${params.toString() ? `?${params.toString()}` : ""}`;
    if (typeof window !== "undefined") window.localStorage.setItem(STORAGE_KEY, JSON.stringify(values));
    router.push(next as Route);
  }

  function clearFilters() {
    const params = new URLSearchParams(searchParams.toString());
    for (const k of FILTER_KEYS) params.delete(k);
    const empty: FilterState = {};
    for (const k of FILTER_KEYS) empty[k] = "";
    setValues(empty);
    if (typeof window !== "undefined") window.localStorage.removeItem(STORAGE_KEY);
    router.push(pathname as Route);
  }

  return (
    <div className="sticky top-0 z-20 bg-[#f4f7fb]/95 backdrop-blur border-b border-slate-200 mb-5">
      <div className="py-3 flex flex-wrap gap-2 items-end">
        <input list="period-options" className="bg-white border rounded px-2 py-1 text-sm w-32" placeholder="Period" value={values.period ?? ""} onChange={(e) => setValues((p) => ({ ...p, period: e.target.value }))} />
        <datalist id="period-options">{options.period.map((o) => <option key={o} value={o} />)}</datalist>

        <input list="region-options" className="bg-white border rounded px-2 py-1 text-sm w-32" placeholder="Region" value={values.region ?? ""} onChange={(e) => setValues((p) => ({ ...p, region: e.target.value, district: "", facility: "" }))} />
        <datalist id="region-options">{options.region.map((o) => <option key={o} value={o} />)}</datalist>

        <input list="district-options" className="bg-white border rounded px-2 py-1 text-sm w-36" placeholder="District" value={values.district ?? ""} onChange={(e) => setValues((p) => ({ ...p, district: e.target.value, facility: "" }))} />
        <datalist id="district-options">{options.district.map((o) => <option key={o} value={o} />)}</datalist>

        <input list="facility-options" className="bg-white border rounded px-2 py-1 text-sm w-44" placeholder="Facility" value={values.facility ?? ""} onChange={(e) => setValues((p) => ({ ...p, facility: e.target.value }))} />
        <datalist id="facility-options">{options.facility.map((o) => <option key={o} value={o} />)}</datalist>

        <input list="level-options" className="bg-white border rounded px-2 py-1 text-sm w-32" placeholder="Level" value={values.level ?? ""} onChange={(e) => setValues((p) => ({ ...p, level: e.target.value }))} />
        <datalist id="level-options">{options.level.map((o) => <option key={o} value={o} />)}</datalist>
        <button onClick={applyFilters} className="bg-navy text-white rounded px-3 py-1.5 text-sm">Apply</button>
        <button onClick={clearFilters} className="bg-white border rounded px-3 py-1.5 text-sm">Clear</button>
        <span className="text-xs text-slate-500 ml-1">{activeCount} active</span>
      </div>
    </div>
  );
}
