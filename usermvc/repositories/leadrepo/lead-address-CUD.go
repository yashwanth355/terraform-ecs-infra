package leadrepo

import (
	"context"
	"usermvc/model"
	apputils "usermvc/utility/apputils"
	logger2 "usermvc/utility/logger"

	"github.com/jinzhu/gorm"
)

/*
*
 */
func (leadRepoRef leadRepo) AddBillingAddressOnNewLead(leadId string, ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error) {

	insertBAddressRecordSQL := `INSERT INTO dbo.cms_leads_billing_address_master 
		(leadid, billingid, street, city, stateprovince, postalcode, country) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := leadRepoRef.db.Raw(insertBAddressRecordSQL, leadId, leadId,
		requestPayload.BillingStreetAddress, requestPayload.BillingCity,
		requestPayload.BillingState, requestPayload.BillingPostalCode,
		requestPayload.BillingCountry).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while creating Billing Address on New Lead Creation in dbo.cms_leads_billing_address_master ", err.Error())
		return false, err
	}
	return true, nil
}

/*
*
 */
func (leadRepoRef leadRepo) AddShippingAddressOnNewLead(leadId string, ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error) {

	insertShipAddressRecordSQL := `INSERT INTO dbo.cms_leads_shipping_address_master 
		(leadid, shippingid, street, city, stateprovince, postalcode, country) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := leadRepoRef.db.Raw(insertShipAddressRecordSQL, leadId, leadId,
		requestPayload.ContactStreetAddress, apputils.NullColumnValue(requestPayload.ContactCity),
		requestPayload.ContactState, requestPayload.ContactPostalCode,
		requestPayload.ContactCountrycode).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while creating Shipping Address on New Lead Creation in dbo.cms_leads_shipping_address_master ", err.Error())
		return false, err
	}
	return true, nil
}

/*
*
 */
func (leadRepoRef leadRepo) UpdateLeadsBillingaddress(ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error) {

	updateLeadBARecordSQL := `UPDATE dbo.cms_leads_billing_address_master 
		SET street=$1, city=$2, stateprovince=$3, 
		postalcode=$4, country=$5 where billingid=$6`

	_, err := db.Raw(updateLeadBARecordSQL, requestPayload.BillingStreetAddress,
		requestPayload.BillingCity, requestPayload.BillingState, requestPayload.BillingPostalCode,
		requestPayload.BillingCountry, requestPayload.LeadId).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while updating Leads Billing Address in dbo.cms_leads_billing_address_master ", err.Error())
		return false, err
	}
	return true, nil
}

/*
*
 */
func (leadRepoRef leadRepo) UpdateLeadsShippingaddress(ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error) {

	updateLeadShpAdrsRecordSQL := `UPDATE dbo.cms_leads_shipping_address_master 
			SET street=$1, city=$2, stateprovince=$3, 
			postalcode=$4, country=$5 where shippingid=$6`

	_, err := db.Raw(updateLeadShpAdrsRecordSQL, requestPayload.ContactStreetAddress,
		apputils.NullColumnValue(requestPayload.ContactCity), requestPayload.ContactState, requestPayload.ContactPostalCode,
		requestPayload.ContactCountrycode, requestPayload.LeadId).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while updating Leads Shipping Address in dbo.cms_leads_billing_address_master ", err.Error())
		return false, err
	}
	return true, nil
}
