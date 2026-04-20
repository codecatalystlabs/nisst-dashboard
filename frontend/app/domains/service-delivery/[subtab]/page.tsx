import { DomainComplianceChart } from "@/components/charts/domain-compliance-chart";
import { HorizontalRanking } from "@/components/charts/horizontal-ranking";
import { QuestionPerformanceTable } from "@/components/tables/question-performance-table";
import { KpiCard } from "@/components/ui/kpi-card";
import { apiGet, type DomainScoreItem, type QuestionPerformanceItem, type RankingItem, type Summary } from "@/lib/api";
import { buildFilterQuery, type SearchParamMap } from "@/lib/filters";
import { serviceDeliveryTabs, type ServiceDeliverySubtab } from "@/lib/service-delivery";
import { ServiceDeliveryTabNav } from "@/components/layout/service-delivery-tab-nav";
import { Suspense } from "react";

type Props = { params: Promise<{ subtab: ServiceDeliverySubtab }>; searchParams: Promise<SearchParamMap> };

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

export default async function ServiceDeliverySubtabPage({ params, searchParams }: Props) {
  const { subtab } = await params;
  const baseSearch = await searchParams;
  const tab = serviceDeliveryTabs.find((t) => t.slug === subtab) ?? serviceDeliveryTabs[0];
  const merged: SearchParamMap = { ...baseSearch, domain: "service delivery", subdomain: tab.slug };
  const filter = buildFilterQuery(merged);
  const data = await (async () => {
    try {
      const [summary, domains, facilities, questions] = await Promise.all([
        apiGet<Summary>(`/analytics/summary${filter}`),
        apiGet<{ items: DomainScoreItem[] }>(`/analytics/domain-scores${filter}`),
        apiGet<{ items: RankingItem[] }>(`/analytics/facility-ranking${filter}`),
        apiGet<{ items: QuestionPerformanceItem[] }>(`/analytics/question-performance${filter}`),
      ]);
      return {
        summary,
        domains: Array.isArray(domains?.items) ? domains.items : [],
        facilities: Array.isArray(facilities?.items) ? facilities.items : [],
        questions: Array.isArray(questions?.items) ? questions.items : [],
      };
    } catch {
      return loadData();
    }
  })();

  const allQuestions = Array.isArray(data.questions) ? data.questions : [];
  const scopedQuestions = allQuestions.filter((q) => {
    const hay = q.question.toLowerCase();
    return tab.keywords.some((k) => hay.includes(k));
  });
  const questions = scopedQuestions.length > 0 ? scopedQuestions : allQuestions;
  const topFacilities = (Array.isArray(data.facilities) ? data.facilities : []).slice(0, 10).map((x) => ({ name: x.name, score: x.compliance }));

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold text-navy">{tab.label}</h2>
        <p className="text-slate-600">Scoped performance view for Service Delivery subsection analytics.</p>
      </div>
      <Suspense fallback={null}>
        <ServiceDeliveryTabNav activeSlug={tab.slug} />
      </Suspense>

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
