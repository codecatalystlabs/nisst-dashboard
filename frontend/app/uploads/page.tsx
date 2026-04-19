import { apiGet, type UploadBatchItem } from "@/lib/api";

async function loadUploads() {
  try {
    const res = await apiGet<{ items: UploadBatchItem[] }>("/uploads?limit=100");
    return res.items;
  } catch {
    return [] as UploadBatchItem[];
  }
}

export default async function UploadsPage() {
  const rows = await loadUploads();
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold text-navy">Data Uploads</h2>
        <p className="text-slate-600">Main and follow-up import batches with audit outcomes.</p>
      </div>

      <div className="bg-white rounded-xl shadow-midas p-4 overflow-auto">
        <table className="min-w-full text-sm">
          <thead>
            <tr className="text-left text-slate-500">
              <th className="py-2 pr-3">Type</th>
              <th className="py-2 pr-3">File</th>
              <th className="py-2 pr-3">Status</th>
              <th className="py-2 pr-3">Uploader</th>
              <th className="py-2 pr-3">Rows</th>
              <th className="py-2 pr-3">Imported</th>
              <th className="py-2 pr-3">Duplicates</th>
              <th className="py-2 pr-3">Errors</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => (
              <tr key={r.id} className="border-t">
                <td className="py-2 pr-3">{r.file_type}</td>
                <td className="py-2 pr-3">{r.file_name}</td>
                <td className="py-2 pr-3">{r.status}</td>
                <td className="py-2 pr-3">{r.uploader}</td>
                <td className="py-2 pr-3">{r.total_rows}</td>
                <td className="py-2 pr-3">{r.imported_rows}</td>
                <td className="py-2 pr-3">{r.duplicate_rows}</td>
                <td className="py-2 pr-3">{r.error_rows}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
