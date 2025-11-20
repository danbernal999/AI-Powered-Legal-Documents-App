package repository

func (r *Repository) SeedTemplates() error {
	templates := []struct {
		id          string
		name        string
		description string
		typeStr     string
		content     string
		variables   string
	}{
		{
			id:          "nda",
			name:        "Non-Disclosure Agreement",
			description: "Protect confidential information with a professional NDA",
			typeStr:     "nda",
			content:     `NON-DISCLOSURE AGREEMENT

This Non-Disclosure Agreement ("Agreement") is entered into as of {{execution_date}}, between {{disclosing_party}}, a {{disclosing_entity_type}} ("Disclosing Party"), and {{receiving_party}}, a {{receiving_entity_type}} ("Receiving Party").

WHEREAS, the Disclosing Party possesses certain confidential information and desires to disclose it to the Receiving Party under the terms and conditions set forth in this Agreement; and

NOW, THEREFORE, in consideration of the mutual covenants and agreements contained herein, the parties agree as follows:

1. DEFINITIONS
   1.1 "Confidential Information" means all non-public technical, business, financial, and other information disclosed by the Disclosing Party to the Receiving Party, including but not limited to trade secrets, patents, patent applications, research, product plans, designs, source code, and customer lists.

2. OBLIGATIONS
   2.1 The Receiving Party agrees to maintain the Confidential Information in strict confidence using the same degree of care it uses to protect its own confidential information, but no less than reasonable care.
   2.2 The Receiving Party shall limit access to the Confidential Information to employees and contractors who need to know and who are bound by confidentiality obligations.

3. PERMITTED DISCLOSURES
   3.1 The Receiving Party may disclose Confidential Information only as required by law, court order, or regulatory authority, provided that it gives the Disclosing Party prompt notice to allow the Disclosing Party to seek protective measures.

4. TERM
   4.1 This Agreement shall commence on the date hereof and continue for a period of {{confidentiality_period}} years unless earlier terminated by written consent of both parties.

5. RETURN OR DESTRUCTION
   5.1 Upon request or termination of this Agreement, the Receiving Party shall return or destroy all Confidential Information and provide written certification of such destruction within thirty (30) days.

6. GOVERNING LAW
   6.1 This Agreement shall be governed by and construed in accordance with the laws of {{governing_jurisdiction}}, without regard to its conflict of laws principles.

7. ENTIRE AGREEMENT
   7.1 This Agreement constitutes the entire agreement between the parties and supersedes all prior and contemporaneous agreements.

IN WITNESS WHEREOF, the parties have executed this Agreement as of the date first written above.

DISCLOSING PARTY:
Signature: _______________________
Name: _______________________
Date: _______________________

RECEIVING PARTY:
Signature: _______________________
Name: _______________________
Date: _______________________`,
			variables: `[{"name":"disclosing_party","description":"Name of the party disclosing information","type":"string","required":true},{"name":"receiving_party","description":"Name of the party receiving information","type":"string","required":true},{"name":"execution_date","description":"Date of agreement execution","type":"string","required":true},{"name":"disclosing_entity_type","description":"Type of entity (e.g., corporation, LLC)","type":"string","required":true},{"name":"receiving_entity_type","description":"Type of entity (e.g., corporation, LLC)","type":"string","required":true},{"name":"confidentiality_period","description":"Duration of confidentiality in years","type":"number","required":true},{"name":"governing_jurisdiction","description":"State or jurisdiction governing this agreement","type":"string","required":true}]`,
		},
		{
			id:          "employment",
			name:        "Employment Contract",
			description: "Create a professional employment agreement",
			typeStr:     "employment",
			content:     `EMPLOYMENT AGREEMENT

This Employment Agreement ("Agreement") is entered into as of {{start_date}}, between {{company_name}}, a {{entity_type}} ("Employer"), and {{employee_name}} ("Employee").

1. POSITION AND DUTIES
   1.1 The Employee shall be employed in the position of {{job_title}}, located in {{work_location}}.
   1.2 The Employee shall perform such duties and responsibilities as are customary for this position and as directed by the Employer.

2. COMPENSATION AND BENEFITS
   2.1 Annual Salary: {{annual_salary}}
   2.2 The Employee shall receive {{pay_frequency}} compensation for services rendered.
   2.3 The Employee shall be entitled to {{vacation_days}} days of paid vacation annually, {{sick_days}} days of sick leave, and {{holidays_count}} paid holidays per year.

3. EMPLOYMENT TERM
   3.1 The Employee's employment shall commence on {{start_date}} and shall be on an {{employment_type}} basis.
   3.2 Either party may terminate this employment relationship at any time, for any reason, with {{notice_period}} days' written notice.

4. CONFIDENTIALITY
   4.1 The Employee acknowledges that during employment, they will have access to confidential information, trade secrets, and proprietary materials of the Employer.
   4.2 The Employee agrees to maintain the confidentiality of such information during and after employment.

5. INTELLECTUAL PROPERTY
   5.1 Any work product, inventions, or intellectual property created by the Employee during the course of employment shall be the sole property of the Employer.

6. NON-COMPETITION
   6.1 During employment and for {{non_compete_period}} years after termination, the Employee agrees not to engage in any competing business within {{competition_radius}}.

7. BENEFITS AND INSURANCE
   7.1 The Employer shall provide the Employee with health insurance, retirement benefits, and other benefits as defined in the Employer's benefits plan.

8. TERMINATION
   8.1 The Employee may be terminated for cause for any serious violation of company policies or law.
   8.2 Upon termination, the Employee shall be entitled to {{severance_package}}.

9. GOVERNING LAW
   9.1 This Agreement shall be governed by the laws of {{governing_jurisdiction}}.

IN WITNESS WHEREOF:

EMPLOYER:
Signature: _______________________
Name/Title: _______________________
Date: _______________________

EMPLOYEE:
Signature: _______________________
Name: _______________________
Date: _______________________`,
			variables: `[{"name":"company_name","description":"Name of the employer company","type":"string","required":true},{"name":"employee_name","description":"Full name of the employee","type":"string","required":true},{"name":"job_title","description":"Position title","type":"string","required":true},{"name":"annual_salary","description":"Annual salary amount","type":"string","required":true},{"name":"work_location","description":"Work location","type":"string","required":true},{"name":"entity_type","description":"Entity type (corporation, LLC, etc.)","type":"string","required":true},{"name":"start_date","description":"Employment start date","type":"string","required":true},{"name":"employment_type","description":"Full-time or Part-time","type":"string","required":true},{"name":"pay_frequency","description":"Payment frequency (weekly, biweekly, monthly)","type":"string","required":true},{"name":"vacation_days","description":"Annual vacation days","type":"number","required":true},{"name":"sick_days","description":"Annual sick days","type":"number","required":true},{"name":"holidays_count","description":"Number of paid holidays","type":"number","required":true},{"name":"notice_period","description":"Termination notice period in days","type":"number","required":true},{"name":"non_compete_period","description":"Non-compete duration in years","type":"number","required":true},{"name":"competition_radius","description":"Geographic radius for non-competition","type":"string","required":true},{"name":"severance_package","description":"Severance details","type":"string","required":true},{"name":"governing_jurisdiction","description":"Governing jurisdiction","type":"string","required":true}]`,
		},
		{
			id:          "rental",
			name:        "Lease Agreement",
			description: "Professional property rental agreement",
			typeStr:     "rental",
			content:     `LEASE AGREEMENT

This Lease Agreement ("Lease") is made and entered into as of {{lease_date}}, between {{landlord_name}}, a {{landlord_type}} ("Landlord"), and {{tenant_name}}, a {{tenant_type}} ("Tenant").

1. PROPERTY
   1.1 The Landlord hereby leases to the Tenant, and the Tenant leases from the Landlord, the property located at:
   {{property_address}}
   {{property_description}}

2. LEASE TERM
   2.1 The lease term shall commence on {{lease_start_date}} and shall continue through {{lease_end_date}}, unless earlier terminated pursuant to the terms herein.
   2.2 This lease shall automatically renew for successive periods unless either party provides written notice of non-renewal at least {{renewal_notice_days}} days prior to expiration.

3. RENT AND PAYMENT
   3.1 Base Monthly Rent: {{monthly_rent}}
   3.2 Rent shall be payable on the {{rent_due_date}} of each month.
   3.3 Rent shall be paid to: {{payment_address}}
   3.4 Late fee: {{late_fee}} if rent is not received by the {{late_fee_date}} of the month.

4. SECURITY DEPOSIT
   4.1 Security Deposit Amount: {{security_deposit}}
   4.2 The security deposit shall be held by the Landlord and shall be returned within {{deposit_return_days}} days after lease termination, less any deductions for damages or unpaid rent.

5. UTILITIES AND SERVICES
   5.1 The Tenant shall be responsible for: {{tenant_utilities}}
   5.2 The Landlord shall be responsible for: {{landlord_utilities}}

6. MAINTENANCE AND REPAIRS
   6.1 The Tenant shall maintain the property in good condition and immediately report any damage or repairs needed.
   6.2 The Landlord shall be responsible for major structural repairs and maintaining common areas.

7. USE OF PROPERTY
   7.1 The property shall be used solely for residential purposes.
   7.2 The Tenant shall not make alterations without the Landlord's written consent.

8. OCCUPANCY
   8.1 Maximum number of occupants: {{max_occupants}}
   8.2 No subleasing without the Landlord's written consent.

9. TERMINATION
   9.1 Either party may terminate with {{termination_notice_days}} days' written notice.
   9.2 If the Tenant leaves before the lease end date, the Tenant shall be liable for remaining rent.

10. GOVERNING LAW
    10.1 This Lease shall be governed by the laws of {{governing_jurisdiction}}.

IN WITNESS WHEREOF:

LANDLORD:
Signature: _______________________
Name: _______________________
Date: _______________________

TENANT:
Signature: _______________________
Name: _______________________
Date: _______________________`,
			variables: `[{"name":"landlord_name","description":"Landlord's full name","type":"string","required":true},{"name":"tenant_name","description":"Tenant's full name","type":"string","required":true},{"name":"property_address","description":"Complete property address","type":"string","required":true},{"name":"property_description","description":"Description of the property","type":"string","required":true},{"name":"lease_date","description":"Date lease is signed","type":"string","required":true},{"name":"lease_start_date","description":"Lease commencement date","type":"string","required":true},{"name":"lease_end_date","description":"Lease end date","type":"string","required":true},{"name":"monthly_rent","description":"Monthly rent amount","type":"string","required":true},{"name":"rent_due_date","description":"Day of month rent is due","type":"number","required":true},{"name":"security_deposit","description":"Security deposit amount","type":"string","required":true},{"name":"landlord_type","description":"Landlord entity type","type":"string","required":true},{"name":"tenant_type","description":"Tenant entity type","type":"string","required":true},{"name":"payment_address","description":"Address to send rent payments","type":"string","required":true},{"name":"late_fee","description":"Late payment fee amount","type":"string","required":true},{"name":"late_fee_date","description":"Date late fee applies","type":"number","required":true},{"name":"deposit_return_days","description":"Days to return deposit after lease ends","type":"number","required":true},{"name":"tenant_utilities","description":"Utilities tenant pays","type":"string","required":true},{"name":"landlord_utilities","description":"Utilities landlord pays","type":"string","required":true},{"name":"max_occupants","description":"Maximum number of occupants allowed","type":"number","required":true},{"name":"termination_notice_days","description":"Days notice required for termination","type":"number","required":true},{"name":"renewal_notice_days","description":"Days notice required for non-renewal","type":"number","required":true},{"name":"governing_jurisdiction","description":"Governing jurisdiction","type":"string","required":true}]`,
		},
		{
			id:          "freelance",
			name:        "Freelance Agreement",
			description: "Independent contractor service agreement",
			typeStr:     "freelance",
			content:     `INDEPENDENT CONTRACTOR AGREEMENT

This Independent Contractor Agreement ("Agreement") is made and entered into as of {{agreement_date}}, between {{client_name}}, a {{client_type}} ("Client"), and {{contractor_name}}, a {{contractor_type}} ("Contractor").

1. SERVICES
   1.1 The Contractor agrees to provide the following services:
   {{scope_of_work}}
   
   1.2 The Contractor shall deliver the work in a professional manner and in accordance with industry standards.

2. COMPENSATION
   2.1 Total Project Fee: {{project_fee}}
   2.2 Payment Terms: {{payment_terms}}
   2.3 Payment shall be made to: {{payment_method}}
   
3. PROJECT TIMELINE
   3.1 Project Start Date: {{project_start_date}}
   3.2 Project Completion Date: {{project_completion_date}}
   3.3 Deliverables: {{deliverables}}

4. INTELLECTUAL PROPERTY RIGHTS
   4.1 Upon receipt of full payment, all intellectual property rights, including copyrights, patents, and trade secrets created by the Contractor for this project shall become the exclusive property of the Client.
   4.2 The Contractor grants the Client all rights needed to use the work product for any purpose.

5. CONFIDENTIALITY
   5.1 The Contractor acknowledges that during the course of this engagement, the Contractor may have access to confidential information of the Client.
   5.2 The Contractor agrees to keep all such information confidential and not disclose it to any third party.

6. INDEPENDENT CONTRACTOR STATUS
   6.1 The Contractor is an independent contractor and not an employee of the Client.
   6.2 The Contractor is responsible for all taxes, insurance, and other statutory obligations.
   6.3 The Contractor shall provide their own equipment and tools.

7. LIABILITY LIMITATIONS
   7.1 The Client and Contractor agree to limit liability to the amount paid for services under this Agreement.

8. TERMINATION
   8.1 Either party may terminate this Agreement with {{termination_notice_days}} days' written notice.
   8.2 Upon termination, the Contractor shall deliver all completed work and return all Client materials.

9. GOVERNING LAW
   9.1 This Agreement shall be governed by the laws of {{governing_jurisdiction}}.

IN WITNESS WHEREOF:

CLIENT:
Signature: _______________________
Name/Title: _______________________
Date: _______________________

CONTRACTOR:
Signature: _______________________
Name: _______________________
Date: _______________________`,
			variables: `[{"name":"client_name","description":"Client's business name","type":"string","required":true},{"name":"contractor_name","description":"Contractor's full name","type":"string","required":true},{"name":"scope_of_work","description":"Detailed description of work to be performed","type":"string","required":true},{"name":"project_fee","description":"Total project compensation","type":"string","required":true},{"name":"payment_terms","description":"Payment schedule and terms","type":"string","required":true},{"name":"payment_method","description":"Payment method (bank transfer, check, etc.)","type":"string","required":true},{"name":"project_start_date","description":"When work begins","type":"string","required":true},{"name":"project_completion_date","description":"Expected project completion date","type":"string","required":true},{"name":"deliverables","description":"Specific deliverables expected","type":"string","required":true},{"name":"client_type","description":"Client entity type","type":"string","required":true},{"name":"contractor_type","description":"Contractor entity type","type":"string","required":true},{"name":"agreement_date","description":"Date agreement is signed","type":"string","required":true},{"name":"termination_notice_days","description":"Days notice for termination","type":"number","required":true},{"name":"governing_jurisdiction","description":"Governing jurisdiction","type":"string","required":true}]`,
		},
	}

	for _, t := range templates {
		query := `
		INSERT INTO templates (id, name, description, type, content, variables)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err := r.db.Exec(query, t.id, t.name, t.description, t.typeStr, t.content, t.variables)
		if err != nil {
			continue
		}
	}

	return nil
}

func (r *Repository) SeedTemplatesIfEmpty() error {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM templates").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return r.SeedTemplates()
	}

	return nil
}
