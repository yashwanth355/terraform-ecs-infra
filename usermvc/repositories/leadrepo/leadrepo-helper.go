package leadrepo

import (
	"strconv"
	"usermvc/entity"
	"usermvc/model"
)

/*
*
 */
func makeLeadId(newLeadNumber int64) string {
	return "Lead-" + strconv.FormatInt(newLeadNumber, 10)
}

/*
*
 */
func cmsLeadMasterToLeadDetail(cmsLeads entity.CmsLeadsMaster) model.LeadDetails {

	contactExtId, _ := strconv.ParseInt(cmsLeads.ContactExt, 10, 64)

	return model.LeadDetails{
		Accountname:                cmsLeads.Accountname,
		Accounttypeid:              cmsLeads.Accounttypeid,
		ContactMobile:              cmsLeads.ContactMobile,
		ContactEmail:               cmsLeads.Email,
		Approximativeannualrevenue: cmsLeads.Approxannualrev,
		Website:                    cmsLeads.Website,
		Productsegmentid:           cmsLeads.Productsegmentid,
		Contactfirstname:           cmsLeads.Contactfirstname,
		Contactlastname:            cmsLeads.Contactlastname,
		Contact_Firstname:          cmsLeads.Contactfirstname,
		Contact_Lastname:           cmsLeads.Contactlastname,
		Manfacunit:                 cmsLeads.Manfacunit,
		Instcoffee:                 cmsLeads.Instcoffee,
		Price:                      cmsLeads.Price,
		ContactSalutationid:        cmsLeads.ContactSalutationid,
		ContactPosition:            cmsLeads.ContactPosition,
		ContactPhone:               cmsLeads.ContactMobile,
		Phone:                      cmsLeads.ContactMobile,
		ShippingContinent:          cmsLeads.ShippingContinent,
		ShippingCountry:            cmsLeads.ShippingCountry,
		Coffeetypeid:               cmsLeads.Coffeetypeid,
		Aliases:                    cmsLeads.Aliases,
		OtherInformation:           cmsLeads.Otherinformation,
		Status:                     cmsLeads.Masterstatus,
		Leadscore:                  cmsLeads.Leadscore,
		Contact_extid:              contactExtId,
		ContactExt:                 cmsLeads.ContactExt,
	}
}
