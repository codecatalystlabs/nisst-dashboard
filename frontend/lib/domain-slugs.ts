/**
 * URL segment under /domains/[slug] → value passed as analytics `domain` query param.
 * Must align with `domains.name` in DB seeds (ILIKE match in backend).
 */
export const DOMAIN_SLUG_TO_FILTER: Record<string, string> = {
  "leadership-governance": "Leadership and Governance",
  "human-resources": "Human Resources for Health",
  "medicines-supplies": "Medicines and Health Supplies",
  "health-financing": "Health Financing",
  "health-information": "Health Information Management",
  infrastructure: "Health Infrastructure",
  "service-delivery": "Service Delivery",
  "quality-care-safety": "Quality of Care and Safety",
};

export function domainFilterForSlug(slug: string): string | undefined {
  return DOMAIN_SLUG_TO_FILTER[slug];
}
