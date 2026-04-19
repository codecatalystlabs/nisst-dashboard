import Link from "next/link";

const nav = [
  ["Overview", "/"],
  ["Leadership & Governance", "/domains/leadership-governance"],
  ["Human Resources", "/domains/human-resources"],
  ["Medicines & Supplies", "/domains/medicines-supplies"],
  ["Health Financing", "/domains/health-financing"],
  ["Health Information", "/domains/health-information"],
  ["Infrastructure", "/domains/infrastructure"],
  ["Service Delivery", "/domains/service-delivery"],
  ["Quality of Care & Safety", "/domains/quality-care-safety"],
  ["Follow-up Actions", "/followups"],
  ["Data Uploads", "/uploads"],
  ["Metadata / Configuration", "/metadata"]
] as const;

export function SidebarNav() {
  return (
    <aside className="bg-navy text-white p-6 sticky top-0 h-screen">
      <h1 className="text-lg font-semibold tracking-wide">NISST MIDAS</h1>
      <nav className="mt-6 space-y-2">
        {nav.map(([label, href]) => (
          <Link key={href} href={href} className="block rounded-md px-3 py-2 text-sm hover:bg-white/10">{label}</Link>
        ))}
      </nav>
    </aside>
  );
}
