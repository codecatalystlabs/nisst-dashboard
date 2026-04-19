"use client";
import { Line, LineChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

export function TrendLineChart({ data }: { data: Array<{ period: string; score: number }> }) {
  return <div className="bg-white rounded-xl shadow-midas p-4 h-[300px]"><ResponsiveContainer width="100%" height="100%"><LineChart data={data}><CartesianGrid strokeDasharray="3 3" /><XAxis dataKey="period" /><YAxis domain={[0,1]} /><Tooltip /><Line type="monotone" dataKey="score" stroke="#0B1A33" strokeWidth={3} /></LineChart></ResponsiveContainer></div>;
}
