"use client";

import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

export function HorizontalRanking({ data, title }: { data: Array<{ name: string; score: number }>; title: string }) {
  return (
    <div className="bg-white rounded-xl shadow-midas p-4 h-[360px]">
      <h3 className="text-sm font-medium text-slate-600 mb-3">{title}</h3>
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data} layout="vertical" margin={{ left: 8, right: 12, top: 8, bottom: 8 }}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis type="number" domain={[0, 1]} />
          <YAxis type="category" dataKey="name" width={140} />
          <Tooltip />
          <Bar dataKey="score" fill="#0B1A33" radius={[0, 6, 6, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
