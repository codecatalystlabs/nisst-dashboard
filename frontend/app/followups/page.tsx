import { DomainComplianceChart } from "@/components/charts/domain-compliance-chart";
import { FollowupsTable } from "@/components/tables/followups-table";
import { apiGet, type FollowupItem } from "@/lib/api";
import { buildFilterQuery, type SearchParamMap } from "@/lib/filters";

type LabelCount = { label: string; count: number };

async function loadFollowups(search?: SearchParamMap) {
  const q = buildFilterQuery(search);
  try {
    const [followupsRes, byDomainRes, byRespRes] = await Promise.all([
      apiGet<{ items: FollowupItem[] }>(`/analytics/followups${q ? `${q}&limit=200` : "?limit=200"}`),
      apiGet<{ items: LabelCount[] }>(`/analytics/followups/by-domain${q}`),
      apiGet<{ items: LabelCount[] }>(`/analytics/followups/by-responsibility${q}`),
    ]);
    return { followups: followupsRes.items, byDomain: byDomainRes.items, byResp: byRespRes.items };
  } catch {
    return { followups: [] as FollowupItem[], byDomain: [] as LabelCount[], byResp: [] as LabelCount[] };
  }
}

export default async function FollowupsPage({ searchParams }: { searchParams: Promise<SearchParamMap> }) {
  const { followups, byDomain, byResp } = await loadFollowups(await searchParams);

  const domainChartData = byDomain.map((x) => ({ domain: x.label, score: x.count }));
  const maxResp = byResp.reduce((m, x) => Math.max(m, x.count), 0) || 1;
  const respChartData = byResp.map((x) => ({ domain: x.label, score: x.count / maxResp }));

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold text-navy">Follow-up Actions</h2>
        <p className="text-slate-600">Action burden, responsibilities, and detailed intervention tracking.</p>
      </div>

      <section className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <DomainComplianceChart data={domainChartData} />
        <DomainComplianceChart data={respChartData} />
      </section>

      <section>
        <FollowupsTable rows={followups} />
      </section>
    </div>
  );
}
