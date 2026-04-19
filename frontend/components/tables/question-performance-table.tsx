type Row = { question: string; domain: string; score: number };

export function QuestionPerformanceTable({ rows }: { rows: Row[] }) {
  return (
    <div className="bg-white rounded-xl shadow-midas p-4 overflow-auto">
      <table className="min-w-full text-sm">
        <thead><tr className="text-left text-slate-500"><th className="py-2">Question</th><th className="py-2">Domain</th><th className="py-2">Score</th></tr></thead>
        <tbody>{rows.map((r) => <tr key={r.question} className="border-t"><td className="py-2">{r.question}</td><td className="py-2">{r.domain}</td><td className="py-2">{(r.score*100).toFixed(1)}%</td></tr>)}</tbody>
      </table>
    </div>
  );
}
