package leadrepo

import (
	"context"
	"usermvc/entity"
	logger2 "usermvc/utility/logger"
)

/*
*
 */
func (leadRepoRef leadRepo) GetCmsLeadsShippingAddress(ctx context.Context, leadId string) (*entity.CmsLeadsShippingAddressMaster, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("getting record from cms_account_product_segment_master where leadId  ", leadId)
	var cmsLeadsShippingAddressMaster entity.CmsLeadsShippingAddressMaster
	if err := leadRepoRef.db.Table("dbo.cms_leads_shipping_address_master").Model(&entity.CmsLeadsShippingAddressMaster{}).Where("leadid=?", leadId).Find(&cmsLeadsShippingAddressMaster).Error; err != nil {
		logger.Error("error while getting data from cms_salutation_master")
		return nil, err
	}
	return &cmsLeadsShippingAddressMaster, nil
}

/*
*
 */
func (leadRepoRef leadRepo) GetCmsLeadsBillingAddress(ctx context.Context, leadID string) (*entity.CmsLeadsBillingAddressMaster, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("getting record from cms_leads_billing_address_master where leadId  ", leadID)
	var cmsLeadsBillingAddressMaster entity.CmsLeadsBillingAddressMaster
	if err := leadRepoRef.db.Table("dbo.cms_leads_billing_address_master").Model(&entity.CmsLeadsBillingAddressMaster{}).Where("leadid=?", leadID).Find(&cmsLeadsBillingAddressMaster).Error; err != nil {
		logger.Error("error while getting data from cms_coffeetype_master")
		return nil, err
	}
	return &cmsLeadsBillingAddressMaster, nil
}
