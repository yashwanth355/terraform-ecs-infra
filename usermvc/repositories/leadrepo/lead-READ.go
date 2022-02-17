package leadrepo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"usermvc/entity"
	"usermvc/model"
	logger2 "usermvc/utility/logger"

	"github.com/jinzhu/gorm"
)

/*
*	To get/list all, multiple leads summary info in a grid or something
*
*	controller.GetLeadsInfo -> calls leadRepo.ProvideLeadsData
*
 */
func (leadRepoRef leadRepo) ProvideLeadsData(ctx context.Context,
	reqParams model.ProvideLeadsInfoReqContext) ([]model.LeadInfo, error) {

	var leads []model.LeadInfo
	queryPart := `ORDER BY createddate desc`
	if reqParams.Filter != "" {
		queryPart = "where " + reqParams.Filter + " ORDER BY createddate desc"
	}
	query := fmt.Sprintf(`SELECT leadid, accountname, aliases, contactfirstname, 
		contactlastname, contact_mobile, email, leadscore,masterstatus 
		FROM dbo.LeadsGrid %s`, queryPart)

	db := leadRepoRef.db
	err := db.Raw(query).Scan(&leads).Error
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error from func: Lead Repository -> ProvideLeadsData in querying dbo.LeadsGrid ", err.Error())
		return nil, err
	}
	return leads, nil
}

/*
*
 */
func (leadRepoRef leadRepo) LeadExists(ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest) (bool, error) {

	query := `select accountname from dbo.cms_leads_master where accountname = $1`
	rows, err := leadRepoRef.db.Raw(query, requestPayload.Accountname).Rows()
	defer rows.Close()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while executing query to check if record exists in dbo.cms_leads_master with given Lead Name ", err.Error())
		return true, err
	}
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

/*
*
 */
func getNextLeadNumber(ctx context.Context, db *gorm.DB) (int64, error) {
	type LeadId struct{ LeadNumber int64 }
	var leadIdStruct LeadId
	query := `SELECT idsno as LeadNumber FROM dbo.cms_leads_master where idsno is not null ORDER BY idsno DESC LIMIT 1`
	rows, err := db.Raw(query).Rows()
	defer rows.Close()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("error reading last lead number from cms_leads_master", err.Error())
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&leadIdStruct.LeadNumber)
	}
	return (leadIdStruct.LeadNumber + 1), nil
}

/*
*
 */
func (leadRepoRef leadRepo) GetLeadCreationInfo(ctx context.Context,
	reqParams *model.GetLeadCreationInfoRequest) (interface{}, error) {

	var resultObj interface{}
	var opError error
	resultObj = nil
	opError = nil

	requestedInfoType := reqParams.Type
	const (
		ACCOUNTDETAILS  = "accountDetails"
		PRODUCTSEGMENTS = "productsegments"
		PHONECODES      = "phonecodes"
		COUNTRIES       = "countries"
		COFFEETYPES     = "coffeetypes"
		SALUTATIONS     = "salutations"
		CONTINENTS      = "continents"
	)
	switch {

	case requestedInfoType == ACCOUNTDETAILS:
		resultObj, opError = getAccountTypeInfo(ctx, leadRepoRef)
	case requestedInfoType == PRODUCTSEGMENTS:
		resultObj, opError = getAccountProductSegments(ctx, leadRepoRef)
	case requestedInfoType == PHONECODES:
		resultObj, opError = getPhoneCodes(ctx, leadRepoRef)
	case requestedInfoType == COUNTRIES:
		resultObj, opError = getContinentCountries(reqParams.ContinentName, ctx, leadRepoRef)
	case requestedInfoType == COFFEETYPES:
		resultObj, opError = getCoffeeTypes(ctx, leadRepoRef)
	case requestedInfoType == SALUTATIONS:
		resultObj, opError = getSalutations(ctx, leadRepoRef)
	case requestedInfoType == CONTINENTS:
		resultObj, opError = getContinentsInfo(ctx, leadRepoRef)
	default:
		errMsg := fmt.Sprintf("Requested Information Type should be either of %s, %s, %s ,%s, %s, %s, %s", ACCOUNTDETAILS, PRODUCTSEGMENTS, PHONECODES, COUNTRIES, COFFEETYPES, SALUTATIONS, CONTINENTS)
		return nil, errors.New(errMsg)
	}
	return resultObj, opError
}

/*
*
 */
func getAccountTypeInfo(ctx context.Context, leadRepoRef leadRepo) (interface{}, error) {

	var accountTypeInfo []model.AccountsInformation
	query := `SELECT accounttypeid, accounttype FROM dbo.cms_account_type_master`
	db := leadRepoRef.db
	err := db.Raw(query).Scan(&accountTypeInfo).Error
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error from func: getAccountTypeInfo in querying dbo.cms_account_type_master ", err.Error())
		return nil, err
	}
	return accountTypeInfo, err
}

/*
*
 */
func getAccountProductSegments(ctx context.Context, leadRepoRef leadRepo) (interface{}, error) {

	var productSegments []model.ProductSegments
	query := `SELECT productsegmentid,productsegment FROM dbo.cms_account_product_segment_master`
	db := leadRepoRef.db
	err := db.Raw(query).Scan(&productSegments).Error
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error from func: getAccountProductSegments in querying dbo.cms_account_product_segment_master ", err.Error())
		return nil, err
	}
	return productSegments, err
}

/*
*
 */
func getPhoneCodes(ctx context.Context, leadRepoRef leadRepo) (interface{}, error) {

	var allPhoneCodes []model.PhoneCodes
	db := leadRepoRef.db
	query := `SELECT id, Country_Name, Dial FROM dbo.cms_phonecodes_master`
	rows, err := db.Raw(query).Rows()
	defer rows.Close()
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error from func: getPhoneCodes in querying dbo.cms_phonecodes_master ", err.Error())
		return nil, err
	}
	var onePhoneCodeEntry model.PhoneCodes
	for rows.Next() {
		err = rows.Scan(&onePhoneCodeEntry.Id, &onePhoneCodeEntry.Countryname, &onePhoneCodeEntry.Dialcode)
		allPhoneCodes = append(allPhoneCodes, onePhoneCodeEntry)
	}
	return allPhoneCodes, err
}

/*
*
 */
func getContinentCountries(continent string, ctx context.Context, leadRepoRef leadRepo) (interface{}, error) {

	var countriesInfo []model.AccountsInformation
	query := `select countryname from dbo.continents_countries_master where continentname=$1`
	db := leadRepoRef.db
	err := db.Raw(query, continent).Scan(&countriesInfo).Error
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error from func: getContinentCountries in querying dbo.continents_countries_master ", err.Error())
		return nil, err
	}
	return countriesInfo, err
}

/*
*
 */
func getCoffeeTypes(ctx context.Context, leadRepoRef leadRepo) (interface{}, error) {

	var allCoffeeTypes []model.CoffeeTypes
	query := `SELECT id, coffeetype FROM dbo.cms_coffeetype_master`
	db := leadRepoRef.db
	err := db.Raw(query).Scan(&allCoffeeTypes).Error
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error from func: getCoffeeTypes in querying dbo.cms_coffeetype_master ", err.Error())
		return nil, err
	}
	return allCoffeeTypes, err
}

/*
*
 */
func getSalutations(ctx context.Context, leadRepoRef leadRepo) (interface{}, error) {

	var salutionsInfo []model.Salutations
	query := `SELECT salutationid, salutation FROM dbo.cms_salutation_master where isactive=$1`
	db := leadRepoRef.db
	err := db.Raw(query, 1).Scan(&salutionsInfo).Error
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error from func: getSalutations in querying dbo.cms_salutation_master ", err.Error())
		return nil, err
	}
	return salutionsInfo, err
}

/*
*
 */
func getContinentsInfo(ctx context.Context, leadRepoRef leadRepo) (interface{}, error) {

	var continentsInfo []model.Continents
	query := `SELECT distinct "continent_name", "continent_code" FROM "continents"`
	db := leadRepoRef.db
	err := db.Raw(query).Scan(&continentsInfo).Error
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error from func: getContinentsInfo in querying public.continents ", err.Error())
		return nil, err
	}
	return continentsInfo, err
}

/*
*
 */
func (lr leadRepo) GetCmsLeads(ctx context.Context, req model.GetLeadDetailsRequestBody) (entity.CmsLeadsMaster, error) {

	var cmsLeadsMaster entity.CmsLeadsMaster
	if err := lr.db.Table("dbo.cms_leads_master").Where(&entity.CmsLeadsMaster{Leadid: req.Id}).Model(&entity.CmsLeadsMaster{}).Find(&cmsLeadsMaster); err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("error while getting record from cms_leads_master where req is ", req)
	}
	return cmsLeadsMaster, nil
}

/*
*
 */
func (lr leadRepo) GetSalutation(ctx context.Context, salutationid int64) (*entity.CmsSalutationMaster, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("getting record from cms_salutation_master where SalutationId  ", salutationid)
	var salutation entity.CmsSalutationMaster
	if err := lr.db.Table("dbo.cms_salutation_master").Where("salutationid=?", salutationid).Find(&salutation).Error; err != nil {
		logger.Error("error while getting data from cms_salutation_master")
		return nil, err
	}
	return &salutation, nil
}

/*
*
 */
func (leadRepoRef leadRepo) GetcmsCoffeetype(ctx context.Context, Id int64) (*entity.CmsCoffeetypeMaster, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("getting record from cms_coffeetype_master where leadId  ", Id)
	var cmsCoffeetypeMaster entity.CmsCoffeetypeMaster
	if err := leadRepoRef.db.Table("dbo.cms_coffeetype_master").Model(&entity.CmsCoffeetypeMaster{}).Where("id=?", Id).Find(&cmsCoffeetypeMaster).Error; err != nil {
		logger.Error("error while getting data from cms_coffeetype_master")
		return nil, err
	}
	return &cmsCoffeetypeMaster, nil
}

/*
*
 */
func (lr leadRepo) GetCmsAccountProductSegment(ctx context.Context, productsegmentid int) (*entity.CmsAccountProductSegmentMaster, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("getting record from cms_account_product_segment_master where leadId  ", productsegmentid)
	var cmsAccountProductSegmentMaster entity.CmsAccountProductSegmentMaster
	if err := lr.db.Table("dbo.cms_account_product_segment_master").Where("productsegmentid=?", productsegmentid).Find(&cmsAccountProductSegmentMaster).Error; err != nil {
		logger.Error("error while getting data from cms_salutation_master")
		return nil, err
	}
	return &cmsAccountProductSegmentMaster, nil
}

/*
*
 */
func (lr leadRepo) GetCmsPhonecodes(ctx context.Context, contact_extid int64) (*entity.CmsPhonecodesMaster, error) {

	var cmsPhonecodesMaster entity.CmsPhonecodesMaster
	sqlStatement := `SELECT * FROM "dbo.cms_phonecodes_master" where id=$1`
	rows, err := lr.db.Raw(sqlStatement, contact_extid).Rows()
	defer rows.Close()
	if err != nil {
		fmt.Println("SELECT * FROM dbo.cms_phonecodes_master ereor -->", err)
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Info(err, "unable to add contact ext code", contact_extid)
	}
	if err == nil {
		for rows.Next() {
			err = rows.Scan(&cmsPhonecodesMaster.Id, &cmsPhonecodesMaster.Countryname, &cmsPhonecodesMaster.Dial)
		}
		return &cmsPhonecodesMaster, nil
	}
	return nil, err
}

/*
*
 */
func (lr leadRepo) GetCmsAccountType(ctx context.Context, accountTypeID int64) (*entity.CmsAccountTypeMaster, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("getting record from cms_account_type_master where accountId ", accountTypeID)
	var cmsAccountTypeMaster entity.CmsAccountTypeMaster
	if err := lr.db.Table("dbo.cms_account_type_master").Model(&entity.CmsAccountTypeMaster{}).Where("accounttypeid=?", accountTypeID).Find(&cmsAccountTypeMaster).Error; err != nil {
		logger.Error("error while getting data from cms_phonecodes_master")
		return nil, err
	}
	return &cmsAccountTypeMaster, nil
}

/*
*
 */
func (leadRepoRef leadRepo) GetDetailsOfLead(ctx context.Context,
	leadDetailsRequest model.GetLeadDetailsRequestBody) (model.LeadDetails, error) {

	logger := logger2.GetLoggerWithContext(ctx)
	var processingErr error = nil
	var retLeadDetails model.LeadDetails
	var leadMasterEntity entity.CmsLeadsMaster
	leadMasterEntity, processingErr = leadRepoRef.GetCmsLeads(ctx, leadDetailsRequest)
	if processingErr == nil {
		retLeadDetails = cmsLeadMasterToLeadDetail(leadMasterEntity)
	}
	if processingErr == nil {
		var salutationEntity *entity.CmsSalutationMaster
		salutationEntity, processingErr = leadRepoRef.GetSalutation(ctx, retLeadDetails.ContactSalutationid)
		retLeadDetails.Salutations = model.Salutations{
			Salutationid: salutationEntity.Salutationid,
			Salutation:   salutationEntity.Salutation,
		}
	}
	if processingErr == nil {
		var billingAddressEntity *entity.CmsLeadsBillingAddressMaster
		billingAddressEntity, processingErr = leadRepoRef.GetCmsLeadsBillingAddress(ctx, leadDetailsRequest.Id)
		retLeadDetails.BillingStreetAddress = billingAddressEntity.Street
		retLeadDetails.BillingCity = billingAddressEntity.City
		retLeadDetails.BillingState = billingAddressEntity.Stateprovince
		retLeadDetails.BillingPostalCode = billingAddressEntity.Postalcode
		retLeadDetails.BillingCountry = billingAddressEntity.Country
	}
	if processingErr == nil {
		var shippingAddressEntity *entity.CmsLeadsShippingAddressMaster
		shippingAddressEntity, processingErr = leadRepoRef.GetCmsLeadsShippingAddress(ctx, leadDetailsRequest.Id)

		retLeadDetails.ContactStreetAddress = shippingAddressEntity.Street
		retLeadDetails.ContactCity = shippingAddressEntity.City
		retLeadDetails.ContactState = shippingAddressEntity.Stateprovince
		retLeadDetails.ContactPostalCode = shippingAddressEntity.Postalcode
		retLeadDetails.ContactCountry = shippingAddressEntity.Country
	}
	if processingErr == nil && retLeadDetails.Coffeetypeid != "" {

		coffeIds := strings.Split(retLeadDetails.Coffeetypeid, ",")
		var coffeeTypeEntity *entity.CmsCoffeetypeMaster
		for _, coffeId := range coffeIds {
			cofeeIdInt, _ := strconv.ParseInt(coffeId, 10, 64)
			coffeeTypeEntity, processingErr = leadRepoRef.GetcmsCoffeetype(ctx, cofeeIdInt)
			if processingErr != nil {
				logger.Error("Error from GetcmsCoffeetype", processingErr.Error())
				break
			}
			retLeadDetails.CoffeeTypes = append(retLeadDetails.CoffeeTypes, model.CoffeeTypes{
				CoffeeType:   coffeeTypeEntity.Coffeetype,
				CoffeeTypeId: strconv.FormatInt(coffeeTypeEntity.ID, 10),
			})
		}
	}
	if processingErr == nil && retLeadDetails.Accounttypeid != "" {

		accountTypeIds := strings.Split(retLeadDetails.Accounttypeid, ",")
		var accountTypeIdInt int64
		var accountTypeEntity *entity.CmsAccountTypeMaster
		for _, Accounttypeid := range accountTypeIds {

			accountTypeIdInt, _ = strconv.ParseInt(Accounttypeid, 10, 64)
			accountTypeEntity, processingErr = leadRepoRef.GetCmsAccountType(ctx, accountTypeIdInt)
			if processingErr != nil {
				logger.Error("Error from GetCmsAccountType", processingErr.Error())
				break
			}
			retLeadDetails.AccountTypes = append(retLeadDetails.AccountTypes, model.AccountsInformation{
				Accounttypeid: strconv.FormatInt(accountTypeEntity.Accounttypeid, 10),
				Accounttype:   accountTypeEntity.Accounttype,
			})
		}
	}
	if processingErr == nil && retLeadDetails.Productsegmentid != "" {

		var productSegmentEntity *entity.CmsAccountProductSegmentMaster
		segmentIds := strings.Split(retLeadDetails.Productsegmentid, ",")
		for _, segmentId := range segmentIds {

			segmentIdToInt, _ := strconv.Atoi(segmentId)
			productSegmentEntity, processingErr = leadRepoRef.GetCmsAccountProductSegment(ctx, segmentIdToInt)
			if processingErr != nil {
				logger.Error("Error from GetCmsAccountProductSegment", processingErr.Error())
				break
			}
			retLeadDetails.Productsegment = append(retLeadDetails.Productsegment, model.ProductSegments{
				Productsegmentid: productSegmentEntity.Productsegmentid,
				Productsegment:   productSegmentEntity.Productsegment,
			})
		}
	}
	if processingErr == nil {
		var auditLogEntries []model.AuditLogGCPO
		auditLogEntries, processingErr = leadRepoRef.GetAuditLogEntriesOfLead(ctx, leadDetailsRequest.Id)
		if processingErr == nil {
			retLeadDetails.AuditLogEntries = auditLogEntries
		}
	}
	retLeadDetails.ContactExt = strconv.FormatInt(retLeadDetails.Contact_extid, 10)
	/*if leadDetailResp.Contact_extid != 0 {
		contactId := leadDetailResp.Contact_extid
		if err != nil {
			logger.Info("invalid contact id")
		} else {
			contactResp, err := lc.leadRepo.GetCmsPhonecodes(ctx, contactId)
			if err != nil {
				logger.Error("error while getting contact detals", err.Error())
			} else {
				leadDetailResp.Contact_ext = model.PhoneCodes{}
				leadDetailResp.Contact_ext.Id = int(contactResp.Id)
				leadDetailResp.Contact_ext.Countryname = contactResp.Countryname
				leadDetailResp.Contact_ext.Dialcode = contactResp.Dial
			}
		}

	}*/
	return retLeadDetails, processingErr
}

/*func (lr leadRepo) GetAllLeadDetails(ctx context.Context) ([]*model.GetAlleadsResonseBody, error) {
	var result []*model.GetAlleadsResonseBody
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("going to fetch all leadAccounts")
	rows, err := lr.db.Raw(`SELECT "accountname","aliases","contactfirstname","contactlastname","phone","email","approvalstatus" FROM "cms_leads_master"`).Rows()

	if err != nil {
		logger.Error(err.Error())
	}
	for rows.Next() {
		var account model.GetAlleadsResonseBody
		err = rows.Scan(&account.AccountName, &account.Aliases, &account.Contactfirstname, &account.Contactlastname, &account.Phone, &account.Email, &account.ApprovalStatus)
		result = append(result, &account)
	}
	res, _ := json.Marshal(result)
	logger.Info("response from cms_leads_master", string(res))
	return result, nil
}*/

/*func (lr leadRepo) GetLeadDetails(ctx context.Context, req model.GetLeadDetailsRequestBody) (interface{}, error) {

	var cmsLeadsMaster []*entity.CmsLeadsMaster
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("going to fetch all leadAccounts")
	if err := lr.db.Table("cms_leads_master").Where(&entity.CmsLeadsMaster{Leadid: "12"}).Model(&entity.CmsLeadsMaster{}).Find(&cmsLeadsMaster); err != nil {
		logger.Error("error while getting record from cms_leads_master where req is ", req)
	}

	return cmsLeadsMaster, nil
}*/
