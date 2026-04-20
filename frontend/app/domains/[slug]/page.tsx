import { DomainComplianceChart } from "@/components/charts/domain-compliance-chart";
import { HorizontalRanking } from "@/components/charts/horizontal-ranking";
import { TrendLineChart } from "@/components/charts/trend-line-chart";
import { QuestionPerformanceTable } from "@/components/tables/question-performance-table";
import { KpiCard } from "@/components/ui/kpi-card";
import { apiGet, type DomainScoreItem, type QuestionPerformanceItem, type RankingItem, type Summary } from "@/lib/api";
import { domainFilterForSlug } from "@/lib/domain-slugs";
import { buildFilterQuery, mergeFilterSearch, type SearchParamMap } from "@/lib/filters";

export const dynamic = "force-dynamic";

type Props = { params: Promise<{ slug: string }>; searchParams: Promise<SearchParamMap> };

type TrendItem = { period: string; compliance: number };

async function loadDomainData(slug: string, search?: SearchParamMap) {
  const domain = domainFilterForSlug(slug);
  const merged = domain ? mergeFilterSearch(search, { domain }) : { ...search };
  const q = buildFilterQuery(merged);

  const empty = {
    summary: { overall_compliance: 0, facilities_assessed: 0, unresolved_followups: 0 } satisfies Summary,
    domains: [] as DomainScoreItem[],
    facilities: [] as RankingItem[],
    questions: [] as QuestionPerformanceItem[],
    trends: [] as TrendItem[],
    errors: [] as string[],
    domainLabel: domain ?? slug.replace(/-/g, " "),
  };

  const jobs = await Promise.allSettled([
    apiGet<Summary>(`/analytics/summary${q}`),
    apiGet<{ items: DomainScoreItem[] }>(`/analytics/domain-scores${q}`),
    apiGet<{ items: RankingItem[] }>(`/analytics/facility-ranking${q}`),
    apiGet<{ items: QuestionPerformanceItem[] }>(`/analytics/question-performance${q}`),
    apiGet<{ items: TrendItem[] }>(`/analytics/trends${q}`),
  ]);

  if (jobs[0].status === "fulfilled") empty.summary = jobs[0].value;
  else empty.errors.push("summary");
  if (jobs[1].status === "fulfilled") {
    const items = jobs[1].value.items;
    empty.domains = Array.isArray(items) ? items : [];
  } else empty.errors.push("domain-scores");
  if (jobs[2].status === "fulfilled") {
    const items = jobs[2].value.items;
    empty.facilities = Array.isArray(items) ? items : [];
  } else empty.errors.push("facility-ranking");
  if (jobs[3].status === "fulfilled") {
    const items = jobs[3].value.items;
    empty.questions = Array.isArray(items) ? items : [];
  } else empty.errors.push("question-performance");
  if (jobs[4].status === "fulfilled") {
    const items = jobs[4].value.items;
    empty.trends = Array.isArray(items) ? items : [];
  } else empty.errors.push("trends");

  return empty;
}

export default async function DomainPage({ params, searchParams }: Props) {
  const { slug } = await params;
  const sp = await searchParams;
  const data = await loadDomainData(slug, sp);
  const topFacilities = (Array.isArray(data.facilities) ? data.facilities : [])
    .slice(0, 10)
    .map((x) => ({ name: x.name, score: x.compliance }));
  const domainTitle = data.domainLabel;

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold text-navy capitalize">{domainTitle}</h2>
        <p className="text-slate-600">Performance scoped to this domain. Adjust global filters (period, region, etc.) and click Apply.</p>
      </div>

      {data.errors.length > 0 ? (
        <div className="rounded-md border border-amber-300 bg-amber-50 px-4 py-2 text-sm text-amber-800">
          Some analytics sources failed to load: {data.errors.join(", ")}.
        </div>
      ) : null}

      <section className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <KpiCard title="Compliance (this domain)" value={`${(data.summary.overall_compliance * 100).toFixed(1)}%`} />
        <KpiCard title="Facilities Assessed" value={`${data.summary.facilities_assessed}`} />
        <KpiCard title="Open follow-ups (all domains)" value={`${data.summary.unresolved_followups}`} />
      </section>

      <section className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <div className="space-y-2">
          <h3 className="text-sm font-medium text-slate-700">Domain scorecard</h3>
          <DomainComplianceChart
            data={(Array.isArray(data.domains) ? data.domains : []).map((x) => ({ domain: x.domain, score: x.compliance }))}
          />
        </div>
        <div className="space-y-2">
          <h3 className="text-sm font-medium text-slate-700">Compliance by period</h3>
          <TrendLineChart
            data={(Array.isArray(data.trends) ? data.trends : []).map((t) => ({ period: t.period, score: t.compliance }))}
          />
        </div>
      </section>

      <section className="space-y-2">
        <h3 className="text-sm font-medium text-slate-700">Facility ranking (this domain)</h3>
        <HorizontalRanking data={topFacilities} title="Top facilities" />
      </section>

      <section>
        <h3 className="text-sm font-medium text-slate-700 mb-2">Question performance</h3>
        <QuestionPerformanceTable
          rows={(Array.isArray(data.questions) ? data.questions : [])
            .slice(0, 30)
            .map((q) => ({ question: q.question, domain: domainTitle, score: q.compliance }))}
        />
      </section>
    </div>
  );
}
