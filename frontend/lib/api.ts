const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ||
  process.env.API_BASE_URL ||
  "http://127.0.0.1:8082/api/v1";

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, { cache: "no-store" });
  if (!res.ok) throw new Error(`API error: ${res.status} for ${API_BASE}${path}`);
  return res.json() as Promise<T>;
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
