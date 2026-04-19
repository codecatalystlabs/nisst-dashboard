import { DomainComplianceChart } from "@/components/charts/domain-compliance-chart";
import { HorizontalRanking } from "@/components/charts/horizontal-ranking";
import { QuestionPerformanceTable } from "@/components/tables/question-performance-table";
import { KpiCard } from "@/components/ui/kpi-card";
import { apiGet, type DomainScoreItem, type QuestionPerformanceItem, type RankingItem, type Summary } from "@/lib/api";
import { serviceDeliveryTabs, type ServiceDeliverySubtab } from "@/lib/service-delivery";

type Props = { params: Promise<{ subtab: ServiceDeliverySubtab }> };

async function loadData() {
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

export default async function ServiceDeliverySubtabPage({ params }: Props) {
  const { subtab } = await params;
  const tab = serviceDeliveryTabs.find((t) => t.slug === subtab) ?? serviceDeliveryTabs[0];
  const filter = `domain=service%20delivery&subdomain=${encodeURIComponent(tab.slug)}`;
  const data = await (async () => {
    try {
      const [summary, domains, facilities, questions] = await Promise.all([
        apiGet<Summary>(`/analytics/summary?${filter}`),
        apiGet<{ items: DomainScoreItem[] }>(`/analytics/domain-scores?${filter}`),
        apiGet<{ items: RankingItem[] }>(`/analytics/facility-ranking?${filter}`),
        apiGet<{ items: QuestionPerformanceItem[] }>(`/analytics/question-performance?${filter}`),
      ]);
      return { summary, domains: domains.items, facilities: facilities.items, questions: questions.items };
    } catch {
      return loadData();
    }
  })();

  const scopedQuestions = data.questions.filter((q) => {
    const hay = q.question.toLowerCase();
    return tab.keywords.some((k) => hay.includes(k));
  });
  const questions = scopedQuestions.length > 0 ? scopedQuestions : data.questions;
  const topFacilities = data.facilities.slice(0, 10).map((x) => ({ name: x.name, score: x.compliance }));

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold text-navy">{tab.label}</h2>
        <p className="text-slate-600">Scoped performance view for Service Delivery subsection analytics.</p>
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
          rows={questions.slice(0, 30).map((q) => ({ question: q.question, domain: tab.label, score: q.compliance }))}
        />
      </section>
    </div>
  );
}
