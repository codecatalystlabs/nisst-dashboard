type Props = { title: string; value: string; delta?: string };

export function KpiCard({ title, value, delta }: Props) {
  return (
    <div className="bg-white rounded-xl shadow-midas p-5">
      <p className="text-xs uppercase text-slate-500">{title}</p>
      <p className="text-2xl font-semibold mt-1 text-navy">{value}</p>
      {delta ? <p className="text-sm mt-1 text-gold">{delta}</p> : null}
    </div>
  );
}
