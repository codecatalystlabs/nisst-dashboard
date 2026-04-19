import type { FollowupItem } from "@/lib/api";

export function FollowupsTable({ rows }: { rows: FollowupItem[] }) {
  return (
    <div className="bg-white rounded-xl shadow-midas p-4 overflow-auto">
      <table className="min-w-full text-sm">
        <thead>
          <tr className="text-left text-slate-500">
            <th className="py-2 pr-3">Domain</th>
            <th className="py-2 pr-3">Challenge</th>
            <th className="py-2 pr-3">Intervention</th>
            <th className="py-2 pr-3">Responsibility</th>
            <th className="py-2 pr-3">Timeline</th>
            <th className="py-2 pr-3">Status</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((r) => (
            <tr key={r.id} className="border-t align-top">
              <td className="py-2 pr-3">{r.domain || "Unknown"}</td>
              <td className="py-2 pr-3">{r.challenge || "-"}</td>
              <td className="py-2 pr-3">{r.intervention || "-"}</td>
              <td className="py-2 pr-3">{r.responsibility || "-"}</td>
              <td className="py-2 pr-3">{r.timelines || "-"}</td>
              <td className="py-2 pr-3">{r.status || "open"}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
