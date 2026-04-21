export const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ||
  process.env.API_BASE_URL ||
  "http://127.0.0.1:8082/api/v1";

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, { cache: "no-store" });
  if (!res.ok) throw new Error(`API error: ${res.status} for ${API_BASE}${path}`);
  return res.json() as Promise<T>;
}

export type ImportSummary = {
  batch_id: string;
  total_rows: number;
  imported_rows: number;
  duplicate_rows: number;
  error_rows: number;
  file_hash: string;
  status: string;
  dry_run: boolean;
};

/** Multipart CSV upload (main or follow-ups). Pass FormData with `file` and optional `uploader`. */
export async function apiUploadCsv(
  path: "/uploads/main" | "/uploads/followups",
  formData: FormData,
  opts?: { dry_run?: boolean }
): Promise<ImportSummary> {
  const q = new URLSearchParams();
  if (opts?.dry_run === true) q.set("dry_run", "true");
  const qs = q.toString();
  const url = `${API_BASE}${path}${qs ? `?${qs}` : ""}`;
  const res = await fetch(url, { method: "POST", body: formData });
  const text = await res.text();
  if (!res.ok) {
    let msg = text;
    try {
      const j = JSON.parse(text) as { error?: string };
      if (typeof j?.error === "string") msg = j.error;
    } catch {
      /* keep body */
    }
    throw new Error(msg || `Upload failed (${res.status})`);
  }
  return JSON.parse(text) as ImportSummary;
}

export type Summary = {
  overall_compliance: number;
  facilities_assessed: number;
  unresolved_followups: number;
};

export type RankingItem = {
  name: string;
  compliance: number;
};

export type DomainScoreItem = {
  domain: string;
  compliance: number;
};

export type TrendItem = {
  period: string;
  compliance: number;
};

export type QuestionPerformanceItem = {
  question: string;
  compliance: number;
  observations?: number;
};

export type UploadBatchItem = {
  id: string;
  file_type: string;
  file_name: string;
  status: string;
  uploader: string;
  total_rows: number;
  imported_rows: number;
  duplicate_rows: number;
  error_rows: number;
  created_at: string;
};

export type FollowupItem = {
  id: string;
  domain: string;
  challenge: string;
  intervention: string;
  responsibility: string;
  timelines: string;
  status: string;
  parent_key: string;
  created_at: string;
};
