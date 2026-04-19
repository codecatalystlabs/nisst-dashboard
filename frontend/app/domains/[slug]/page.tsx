type Props = { params: Promise<{ slug: string }> };

import { DomainComplianceChart } from "@/components/charts/domain-compliance-chart";
import { HorizontalRanking } from "@/components/charts/horizontal-ranking";
import { QuestionPerformanceTable } from "@/components/tables/question-performance-table";
import { KpiCard } from "@/components/ui/kpi-card";
import { apiGet, type DomainScoreItem, type QuestionPerformanceItem, type RankingItem, type Summary } from "@/lib/api";

async function loadDomainData() {
  try {
    const [summary, domains, facilities, questions] = await Promise.all([
      apiGet<Summary>("/analytics/summary"),
      apiGet<{ items: DomainScoreItem[] }>("/analytics/domain-scores"),
      apiGet<{ items: RankingItem[] }>("/analytics/facility-ranking"),
      apiGet<{ items: QuestionPerformanceItem[] }>("/analytics/question-performance"),
    ]);
    return { summary, domains: domains.items, facilities: facilities.items, questions: questions.items };
  } catch {
    return {
      summary: { overall_compliance: 0, facilities_assessed: 0, unresolved_followups: 0 },
      domains: [] as DomainScoreItem[],
      facilities: [] as RankingItem[],
      questions: [] as QuestionPerformanceItem[],
    };
  }
}

export default async function DomainPage({ params }: Props) {
  const { slug } = await params;
  const data = await loadDomainData();
  const slugLabel = slug.replace(/-/g, " ");
  const topFacilities = data.facilities.slice(0, 10).map((x) => ({ name: x.name, score: x.compliance }));

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold text-navy capitalize">{slugLabel}</h2>
        <p className="text-slate-600">Domain performance, weakest indicators, and facility comparisons.</p>
      </div>

      <section className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <KpiCard title="Overall Compliance" value={`${(data.summary.overall_compliance * 100).toFixed(1)}%`} />
        <KpiCard title="Facilities Assessed" value={`${data.summary.facilities_assessed}`} />
        <KpiCard title="Open Follow-ups" value={`${data.summary.unresolved_followups}`} />
      </section>

      <section className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <DomainComplianceChart data={data.domains.map((x) => ({ domain: x.domain, score: x.compliance }))} />
        <HorizontalRanking data={topFacilities} title="Top Facility Compliance" />
      </section>

      <section>
        <QuestionPerformanceTable
          rows={data.questions.slice(0, 30).map((q) => ({ question: q.question, domain: slugLabel, score: q.compliance }))}
        />
      </section>
    </div>
  );
}
