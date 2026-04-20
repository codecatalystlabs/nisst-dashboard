export const FILTER_KEYS = ["period", "region", "district", "facility", "level", "domain", "subdomain"] as const;

export type FilterKey = (typeof FILTER_KEYS)[number];
export type SearchParamValue = string | string[] | undefined;
export type SearchParamMap = Record<string, SearchParamValue>;

export function buildFilterQuery(search: SearchParamMap | undefined): string {
  if (!search) return "";
  const params = new URLSearchParams();
  for (const key of FILTER_KEYS) {
    const raw = search[key];
    const value = Array.isArray(raw) ? raw[0] : raw;
    if (typeof value === "string" && value.trim() !== "") {
      params.set(key, value.trim());
    }
  }
  const q = params.toString();
  return q ? `?${q}` : "";
}

/** Shallow merge URL search params with explicit filter overrides (e.g. domain from route slug). */
export function mergeFilterSearch(
  base: SearchParamMap | undefined,
  patch: Partial<Record<FilterKey, string>>
): SearchParamMap {
  const out: SearchParamMap = { ...(base ?? {}) };
  for (const key of FILTER_KEYS) {
    if (!(key in patch)) continue;
    const trimmed = patch[key]?.trim() ?? "";
    if (trimmed === "") delete out[key];
    else out[key] = trimmed;
  }
  return out;
}
