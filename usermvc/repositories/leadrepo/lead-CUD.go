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
func (leadRepoRef leadRepo) CreateNewLead(ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest) error {

	db := leadRepoRef.db
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	var opError error
	newLeadid, insertLeadErr := leadRepoRef.InsertLeadRecord(ctx, requestPayload, db)
	if insertLeadErr == nil {

		var success bool
		success, opError = leadRepoRef.LogNotificationOnLeadCreateUpdate(newLeadid, `Lead Created`, ctx,
			requestPayload, db)
		if success && opError == nil {

			success, opError = leadRepoRef.AddBillingAddressOnNewLead(newLeadid, ctx,
				requestPayload, db)
		}
		if success && opError == nil {

			success, opError = leadRepoRef.AddShippingAddressOnNewLead(newLeadid, ctx,
				requestPayload, db)
		}
		if success && opError == nil {

			success, opError = leadRepoRef.LogToAuditLogOnNewLead(newLeadid, ctx,
				requestPayload, db)
		}
		if success && opError == nil {

			return tx.Commit().Error
		}
		if !success || opError != nil {
			tx.Rollback()
		}
	} else {
		tx.Rollback()
		opError = insertLeadErr
	}
	return opError
}

/*
*
 */
func (leadRepoRef leadRepo) UpdateLead(ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest) error {

	db := leadRepoRef.db
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	var opError error
	var success bool

	success, opError = doUpdateLead(ctx, requestPayload, db)

	if opError == nil && success {

		success, opError = leadRepoRef.UpdateLeadsBillingaddress(ctx,
			requestPayload, db)
	}
	if success && opError == nil {

		success, opError = leadRepoRef.UpdateLeadsShippingaddress(ctx,
			requestPayload, db)
	}
	if success && opError == nil {

		success, opError = leadRepoRef.EditAuditLogEntryOnLeadAmend(ctx,
			requestPayload, db)
	}
	if success && opError == nil {

		success, opError = leadRepoRef.LogNotificationOnLeadCreateUpdate(requestPayload.LeadId,
			`Lead Updated`, ctx, requestPayload, db)
	}
	if success && opError == nil {
		return tx.Commit().Error
	}
	if !success && opError != nil {
		tx.Rollback()
	}
	return opError
}

/*
*
 */
func (leadRepoRef leadRepo) InsertLeadRecord(ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (string, error) {

	newLeadNum, err := getNextLeadNumber(ctx, db)
	if err == nil && newLeadNum != 0 {
		newLeadId, err := doInsertLead(newLeadNum, ctx, requestPayload, db)
		if err == nil && newLeadId != "" {
			return newLeadId, nil
		}
	}
	return "", err
}

/*
*
 */
func doInsertLead(newLeadNum int64, ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (string, error) {

	newLeadId := makeLeadId(newLeadNum)

	insertLeadRecordSQL := `INSERT INTO dbo.cms_leads_master ( leadid, autogencode,
		legacyid, accountname, accounttypeid, phone, email, createddate, createduserid,
		shipping_continentid, countryid, approxannualrev, website, productsegmentid,
		leadscore, masterstatus,contactfirstname, contactlastname, manfacunit, instcoffee,
		price, approvalstatus, contact_salutationid, contact_position, contact_mobile,
		shipping_continent, shipping_country, coffeetypeid, aliases, isactive,
		otherinformation, contact_ext_id ) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 
		$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, 
		$28, $29, $30, $31, $32)`

	_, err := db.Raw(insertLeadRecordSQL,
		newLeadId, newLeadId, newLeadNum,
		requestPayload.Accountname, requestPayload.Accounttypeid,
		requestPayload.ContactPhone, requestPayload.ContactEmail, requestPayload.CreatedDate,
		requestPayload.CreatedUserid, requestPayload.ShippingContinentid,
		requestPayload.ShippingCountryid, requestPayload.Approximativeannualrevenue,
		requestPayload.Website, requestPayload.Productsegmentid,
		requestPayload.Leadscore, requestPayload.Masterstatus,
		requestPayload.Contactfirstname, requestPayload.Contactlastname,
		requestPayload.Manfacunit, requestPayload.Instcoffee,
		requestPayload.Price, requestPayload.Approvalstatus,
		requestPayload.ContactSalutationid, requestPayload.ContactPosition,
		requestPayload.ContactMobile, requestPayload.ShippingContinent,
		requestPayload.ShippingCountry,
		apputils.NullColumnValue(requestPayload.Coffeetypeid),
		requestPayload.Aliases, requestPayload.Isactive,
		apputils.NullColumnValue(requestPayload.OtherInformation),
		requestPayload.Contact_ext).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while inserting New Lead Record to dbo.cms_leads_master ", err.Error())
		return "", err
	}
	return newLeadId, nil
}

/*
*
 */
func doUpdateLead(ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error) {

	updateLeadRecordSQL := `UPDATE dbo.cms_leads_master SET 
	accountname=$1, accounttypeid=$2, contact_mobile=$3,
	email=$4, phone=$5, modifieddate=$6, modifieduserid=$7,
	shipping_continentid=$8, countryid=$9, approxannualrev=$10,
	website=$11, productsegmentid=$12, leadscore=$13,
	contactfirstname=$14, contactlastname=$15, manfacunit=$16,
	instcoffee=$17, price=$18, contact_salutationid=$19, contact_position=$20,
	shipping_continent=$21, shipping_country=$22,
	coffeetypeid=$23, aliases=$24, otherinformation=$25, contact_ext_id=$26
	where leadid=$27`

	_, err := db.Raw(updateLeadRecordSQL,
		requestPayload.Accountname, requestPayload.Accounttypeid,
		requestPayload.ContactMobile, requestPayload.ContactEmail,
		requestPayload.ContactPhone, requestPayload.ModifiedDate,
		requestPayload.ModifiedUserid, requestPayload.ShippingContinentid,
		requestPayload.ShippingCountryid, requestPayload.Approximativeannualrevenue,
		requestPayload.Website, requestPayload.Productsegmentid,
		requestPayload.Leadscore, requestPayload.Contactfirstname,
		requestPayload.Contactlastname, requestPayload.Manfacunit,
		requestPayload.Instcoffee, requestPayload.Price,
		requestPayload.ContactSalutationid, requestPayload.ContactPosition,
		requestPayload.ShippingContinent, requestPayload.ShippingCountry,
		requestPayload.Coffeetypeid, requestPayload.Aliases,
		requestPayload.OtherInformation, requestPayload.Contact_ext,
		requestPayload.LeadId).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while updating Lead Record in dbo.cms_leads_master ", err.Error())
		return false, err
	}
	return true, nil
}

/*
*
 */
func (leadRepoRef leadRepo) LogNotificationOnLeadCreateUpdate(leadId string, status string, ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error) {

	insertNotifRecordSQL := `insert into dbo.notifications_master_newpg
	(userid, objid, status, feature_category)
	values($1, $2, $3, 'Lead')`

	_, err := leadRepoRef.db.Raw(insertNotifRecordSQL,
		requestPayload.CreatedUserid, leadId, status).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while Loggin Notification Record on New Lead Creation in dbo.notifications_master_newpg ", err.Error())
		return false, err
	}
	return true, nil
}

/*
*
 */
func (leadRepoRef leadRepo) ChangeLeadCreator(ctx context.Context,
	newCreatorUserId string, leadId string) error {

	updateQuery := `update dbo.cms_leads_master 
		set createduserid=$1 where leadid=$2`

	_, err := leadRepoRef.db.Raw(updateQuery, newCreatorUserId, leadId).Rows()

	return err
}
