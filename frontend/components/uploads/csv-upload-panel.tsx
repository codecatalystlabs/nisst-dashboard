"use client";

import { useRouter } from "next/navigation";
import { useCallback, useState } from "react";
import { apiUploadCsv, type ImportSummary } from "@/lib/api";

function formatSummary(s: ImportSummary) {
  const id = s.batch_id ? `${s.batch_id.slice(0, 8)}…` : "(no id)";
  return `Batch ${id} · rows ${s.total_rows} · imported ${s.imported_rows} · dup ${s.duplicate_rows} · errors ${s.error_rows} · ${s.status}${s.dry_run ? " (dry run)" : ""}`;
}

function UploadCard({
  title,
  hint,
  path,
}: {
  title: string;
  hint: string;
  path: "/uploads/main" | "/uploads/followups";
}) {
  const router = useRouter();
  const [busy, setBusy] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const onSubmit = useCallback(
    async (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();
      setMessage(null);
      setError(null);
      const form = e.currentTarget;
      const fileInput = form.elements.namedItem("file") as HTMLInputElement;
      const uploaderInput = form.elements.namedItem("uploader") as HTMLInputElement;
      const dryRunInput = form.elements.namedItem("dry_run") as HTMLInputElement;
      const file = fileInput.files?.[0];
      if (!file) {
        setError("Choose a CSV file.");
        return;
      }
      const fd = new FormData();
      fd.append("file", file);
      const uploader = uploaderInput.value.trim();
      if (uploader) fd.append("uploader", uploader);
      const dryRun = dryRunInput.checked;
      setBusy(true);
      try {
        const summary = await apiUploadCsv(path, fd, { dry_run: dryRun });
        setMessage(formatSummary(summary));
        fileInput.value = "";
        if (!dryRun) router.refresh();
      } catch (err) {
        setError(err instanceof Error ? err.message : "Upload failed");
      } finally {
        setBusy(false);
      }
    },
    [path, router]
  );

  return (
    <div className="bg-white rounded-xl shadow-midas p-5 border border-slate-100">
      <h3 className="text-lg font-semibold text-navy">{title}</h3>
      <p className="text-sm text-slate-600 mt-1">{hint}</p>
      <form className="mt-4 space-y-3" onSubmit={onSubmit}>
        <div>
          <label className="block text-xs font-medium text-slate-500 mb-1">CSV file</label>
          <input
            name="file"
            type="file"
            accept=".csv,text/csv"
            className="block w-full text-sm text-slate-700 file:mr-3 file:rounded file:border-0 file:bg-navy file:px-3 file:py-1.5 file:text-white"
            disabled={busy}
          />
        </div>
        <div>
          <label className="block text-xs font-medium text-slate-500 mb-1">Uploader (optional)</label>
          <input
            name="uploader"
            type="text"
            placeholder="Name or role"
            className="w-full border rounded-md px-3 py-2 text-sm"
            disabled={busy}
          />
        </div>
        <label className="flex items-center gap-2 text-sm text-slate-700">
          <input name="dry_run" type="checkbox" className="rounded border-slate-300" disabled={busy} />
          Dry run only (validate, do not persist)
        </label>
        <button
          type="submit"
          disabled={busy}
          className="bg-navy text-white text-sm font-medium rounded-md px-4 py-2 disabled:opacity-50"
        >
          {busy ? "Uploading…" : "Upload"}
        </button>
      </form>
      {message ? <p className="mt-3 text-sm text-emerald-800 bg-emerald-50 border border-emerald-200 rounded-md px-3 py-2">{message}</p> : null}
      {error ? <p className="mt-3 text-sm text-red-800 bg-red-50 border border-red-200 rounded-md px-3 py-2">{error}</p> : null}
    </div>
  );
}

export function CsvUploadPanel() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <UploadCard
        title="Main supervision export"
        hint="NISST main CSV (wide table). Required columns include SubmissionDate, facility_name, level, region, district, period, meta-instanceID."
        path="/uploads/main"
      />
      <UploadCard
        title="Follow-up actions export"
        hint="Child CSV linked by PARENT_KEY → main meta-instanceID. Columns: domain, challenge, intervention, responsibility, timelines, PARENT_KEY, KEY."
        path="/uploads/followups"
      />
    </div>
  );
}
