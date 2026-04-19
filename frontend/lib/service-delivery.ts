export const serviceDeliveryTabs = [
  { slug: "integrated-health-services", label: "Integrated Health Services", keywords: ["integrated", "service"] },
  { slug: "emergency-critical-care", label: "Emergency and Critical Care Services", keywords: ["emergency", "critical"] },
  { slug: "inpatient-medical", label: "Inpatient Medical Services", keywords: ["inpatient", "medical"] },
  { slug: "surgical-services", label: "Surgical Services", keywords: ["surgical", "surgery"] },
  { slug: "nursing-midwifery", label: "Nursing and Midwifery Services", keywords: ["nursing", "midwifery"] },
  { slug: "mncah", label: "Maternal, Newborn, Child and Adolescent Health", keywords: ["maternal", "newborn", "child", "adolescent"] },
  { slug: "community-health", label: "Community Health Interventions / Engagement", keywords: ["community", "engagement"] },
  { slug: "oral-health", label: "Oral Health Care Services", keywords: ["oral", "dental"] },
  { slug: "palliative-care", label: "Palliative Care Services", keywords: ["palliative"] },
  { slug: "laboratory-blood", label: "Laboratory and Blood Supply", keywords: ["laboratory", "lab", "blood"] },
  { slug: "radiology", label: "Radiology Services", keywords: ["radiology", "imaging"] },
] as const;

export type ServiceDeliverySubtab = (typeof serviceDeliveryTabs)[number]["slug"];
