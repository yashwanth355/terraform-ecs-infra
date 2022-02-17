package leadrepo

import (
	"context"
	"database/sql"
	"fmt"
	"usermvc/model"
)

/*
*
 */
func (leadRepoRef leadRepo) ProvideContactInfoForLead2Aaccount(ctx context.Context, leadInfoForL2A model.LeadInfoInLeadToAccount) (model.ContactInfoFromLeadAndMaster, error) {

	var contactsFromLeadAndMaster model.ContactInfoFromLeadAndMaster = model.ContactInfoFromLeadAndMaster{}

	query := `select l.accountname, concat(l.contactfirstname,',',l.contactlastname) as contact,
	l.phone, l.email, l.website,u.username from dbo.cms_leads_master as l LEFT join 
	dbo.users_master_newpg u on u.userid = l.createduserid where l.accountname=$1`
	rows, err := leadRepoRef.db.Raw(query, leadInfoForL2A.LeadName).Rows()
	defer rows.Close()

	if err == nil {
		var username sql.NullString
		for rows.Next() {
			err = rows.Scan(&contactsFromLeadAndMaster.LeadDetails.LeadName,
				&contactsFromLeadAndMaster.LeadDetails.ContactName,
				&contactsFromLeadAndMaster.LeadDetails.Phone,
				&contactsFromLeadAndMaster.LeadDetails.Email,
				&contactsFromLeadAndMaster.LeadDetails.Website, &username)
			contactsFromLeadAndMaster.LeadDetails.UserName = username.String
		}
		if err == nil {

			var likeFilterValue = "'%" + leadInfoForL2A.LeadName + "%'"
			query = `select con.custname, concat(con.contactfirstname,' ',con.contactlastname) as contactname,
			abm.country from dbo.contacts_master con left join dbo.accounts_billing_address_master abm
			on abm.accountid=con.accountid where con.custname ilike %s`

			rows, err = leadRepoRef.db.Raw(fmt.Sprintf(query, likeFilterValue)).Rows()
			defer rows.Close()
			if err == nil {
				var country sql.NullString
				for rows.Next() {
					err = rows.Scan(&contactsFromLeadAndMaster.ContactMasterDetails.CustomerName,
						&contactsFromLeadAndMaster.ContactMasterDetails.ContactName, &country)
					contactsFromLeadAndMaster.ContactMasterDetails.Country = country.String
				}
				return contactsFromLeadAndMaster, err
			}
		}
	}
	return contactsFromLeadAndMaster, err
}

/*
var leadname LeadDetails
	err := json.Unmarshal([]byte(request.Body), &leadname)




	var rows *sql.Rows
	log.Println("Fetch the lead details in confirmation popup")
	var le LeadandERPContacts
	var userName, country sql.NullString
	sqlStatement1 := `select l.accountname,concat(l.contactfirstname,',',l.contactlastname) as contact,
					  l.phone,l.email,l.website,u.username
					  from dbo.cms_leads_master as l
					  LEFT join dbo.users_master_newpg u on u.userid=l.createduserid
					  where l.accountname=$1`
	rows, err = db.Query(sqlStatement1, leadname.LeadName)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{500, headers, nil, err.Error(), false}, nil
	}
	// defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&le.LeadDetails.LeadName, &le.LeadDetails.Contact, &le.LeadDetails.Phone,
			&le.LeadDetails.Email, &le.LeadDetails.Website, &userName)
		le.LeadDetails.UserName = userName.String
	}
	var parameter string
	parameter = "'%" + leadname.LeadName + "%'"
	log.Println("Fetch the existing ERP Contacts in confirmation Popup")
	sqlStatementERPCon := `select con.custname,
							concat(con.contactfirstname,' ',con.contactlastname) as contactname,
							abm.country
							from dbo.contacts_master con
							left join dbo.accounts_billing_address_master abm
							on abm.accountid=con.accountid
							where con.custname ilike %s`
	rows, err = db.Query(fmt.Sprintf(sqlStatementERPCon, parameter))
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{500, headers, nil, err.Error(), false}, nil
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&le.ERPContactDetails.CustomerName, &le.ERPContactDetails.ContactName, &country)
		le.ERPContactDetails.Country = country.String
	}
	res, _ := json.Marshal(le) */
