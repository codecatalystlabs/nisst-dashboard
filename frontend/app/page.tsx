import { KpiCard } from "@/components/ui/kpi-card";
import { DomainComplianceChart } from "@/components/charts/domain-compliance-chart";
import { TrendLineChart } from "@/components/charts/trend-line-chart";
import { QuestionPerformanceTable } from "@/components/tables/question-performance-table";
import { apiGet } from "@/lib/api";

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

async function loadOverview() {
  try {
    const [summary, domainRes, trendRes, questionRes, gapsRes] = await Promise.all([
      apiGet<Summary>("/analytics/summary"),
      apiGet<{ items: DomainScore[] }>("/analytics/domain-scores"),
      apiGet<{ items: Trend[] }>("/analytics/trends"),
      apiGet<{ items: QuestionPerf[] }>("/analytics/question-performance"),
      apiGet<{ items: QuestionPerf[] }>("/analytics/gaps"),
    ]);

    return { summary, domainRes, trendRes, questionRes, gapsRes };
  } catch {
    return {
      summary: { overall_compliance: 0, facilities_assessed: 0, unresolved_followups: 0 },
      domainRes: { items: [] as DomainScore[] },
      trendRes: { items: [] as Trend[] },
      questionRes: { items: [] as QuestionPerf[] },
      gapsRes: { items: [] as QuestionPerf[] },
    };
  }
}

export default async function OverviewPage() {
  const { summary, domainRes, trendRes, questionRes, gapsRes } = await loadOverview();
  const weakest = gapsRes.items[0];

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-semibold text-navy">Executive Overview</h2>
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
