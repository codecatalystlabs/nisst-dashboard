"use client";
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

export function DomainComplianceChart({ data }: { data: Array<{ domain: string; score: number }> }) {
  return <div className="bg-white rounded-xl shadow-midas p-4 h-[320px]"><ResponsiveContainer width="100%" height="100%"><BarChart data={data}><CartesianGrid strokeDasharray="3 3" /><XAxis dataKey="domain" /><YAxis /><Tooltip /><Bar dataKey="score" fill="#C7A64B" /></BarChart></ResponsiveContainer></div>;
}
