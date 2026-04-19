INSERT INTO domains (code, name, display_order) VALUES
('D1','Leadership and Governance',1),
('D2','Human Resources for Health',2),
('D3','Medicines and Health Supplies',3),
('D4','Health Financing',4),
('D5','Health Information Management',5),
('D6','Health Infrastructure',6),
('D7','Service Delivery',7),
('D8','Quality of Care and Safety',8),
('D9','Follow-up Action Areas',9)
ON CONFLICT (code) DO NOTHING;

INSERT INTO indicator_definitions (code, name, expression, entity_scope) VALUES
('overall_compliance','Overall Compliance','sum(scoreable_yes)/sum(scoreable_total)','global'),
('domain_compliance','Domain Compliance','sum(domain_yes)/sum(domain_total)','domain'),
('subdomain_compliance','Subdomain Compliance','sum(subdomain_yes)/sum(subdomain_total)','subdomain'),
('question_compliance','Question Compliance','sum(question_yes)/sum(question_total)','question'),
('followup_count','Follow-up Count','count(follow_up_actions)','domain')
ON CONFLICT (code) DO NOTHING;
