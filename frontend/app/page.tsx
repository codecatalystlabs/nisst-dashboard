import { KpiCard } from "@/components/ui/kpi-card";
import { DomainComplianceChart } from "@/components/charts/domain-compliance-chart";
import { TrendLineChart } from "@/components/charts/trend-line-chart";
import { QuestionPerformanceTable } from "@/components/tables/question-performance-table";
import { apiGet } from "@/lib/api";
import { buildFilterQuery, type SearchParamMap } from "@/lib/filters";

export const dynamic = "force-dynamic";

type Summary = {
  overall_compliance: number;
  facilities_assessed: number;
  unresolved_followups: number;
};

type DomainScore = {
  domain: string;
  compliance: number;
};

type Trend = {
  period: string;
  compliance: number;
};

type QuestionPerf = {
  question: string;
  compliance: number;
};

async function loadOverview(searchParams?: SearchParamMap) {
  const q = buildFilterQuery(searchParams);
  const base = {
    summary: { overall_compliance: 0, facilities_assessed: 0, unresolved_followups: 0 },
    domainRes: { items: [] as DomainScore[] },
    trendRes: { items: [] as Trend[] },
    questionRes: { items: [] as QuestionPerf[] },
    gapsRes: { items: [] as QuestionPerf[] },
    errors: [] as string[],
  };

  const jobs = await Promise.allSettled([
    apiGet<Summary>(`/analytics/summary${q}`),
    apiGet<{ items: DomainScore[] }>(`/analytics/domain-scores${q}`),
    apiGet<{ items: Trend[] }>(`/analytics/trends${q}`),
    apiGet<{ items: QuestionPerf[] }>(`/analytics/question-performance${q}`),
    apiGet<{ items: QuestionPerf[] }>(`/analytics/gaps${q}`),
  ]);

  if (jobs[0].status === "fulfilled") base.summary = jobs[0].value;
  else base.errors.push("summary");
  if (jobs[1].status === "fulfilled") base.domainRes = jobs[1].value;
  else base.errors.push("domain-scores");
  if (jobs[2].status === "fulfilled") base.trendRes = jobs[2].value;
  else base.errors.push("trends");
  if (jobs[3].status === "fulfilled") base.questionRes = jobs[3].value;
  else base.errors.push("question-performance");
  if (jobs[4].status === "fulfilled") base.gapsRes = jobs[4].value;
  else base.errors.push("gaps");

  return base;
}

export default async function OverviewPage({ searchParams }: { searchParams: Promise<SearchParamMap> }) {
  const resolved = await searchParams;
  const { summary, domainRes, trendRes, questionRes, gapsRes, errors } = await loadOverview(resolved);
  const weakest = gapsRes.items[0];

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-semibold text-navy">Executive Overview</h2>
      {errors.length > 0 ? (
        <div className="rounded-md border border-amber-300 bg-amber-50 px-4 py-2 text-sm text-amber-800">
          Some analytics sources failed to load: {errors.join(", ")}.
        </div>
      ) : null}
      <section className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <KpiCard title="Overall Compliance" value={`${(summary.overall_compliance * 100).toFixed(1)}%`} />
        <KpiCard title="Facilities Assessed" value={`${summary.facilities_assessed}`} />
        <KpiCard title="Unresolved Follow-ups" value={`${summary.unresolved_followups}`} />
        <KpiCard
          title="Weakest Indicator"
          value={weakest ? `${weakest.question} (${(weakest.compliance * 100).toFixed(1)}%)` : "No data"}
        />
      </section>
      <section className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <DomainComplianceChart data={domainRes.items.map((d) => ({ domain: d.domain, score: d.compliance }))} />
        <TrendLineChart data={trendRes.items.map((t) => ({ period: t.period, score: t.compliance }))} />
      </section>
      <QuestionPerformanceTable
        rows={questionRes.items.slice(0, 20).map((q) => ({ question: q.question, domain: "Cross-domain", score: q.compliance }))}
      />
    </div>
  );
}
