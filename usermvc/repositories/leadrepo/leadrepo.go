package leadrepo

import (
	"context"
	"usermvc/entity"
	"usermvc/model"
	"usermvc/repositories"

	"github.com/jinzhu/gorm"
)

type LeadRepo interface {
	GetCmsLeads(ctx context.Context, req model.GetLeadDetailsRequestBody) (entity.CmsLeadsMaster, error)

	GetSalutation(ctx context.Context, salutationid int64) (*entity.CmsSalutationMaster, error)
	GetCmsLeadsShippingAddress(ctx context.Context, leadId string) (*entity.CmsLeadsShippingAddressMaster, error)
	GetCmsAccountProductSegment(ctx context.Context, productsegmentid int) (*entity.CmsAccountProductSegmentMaster, error)
	GetCmsPhonecodes(ctx context.Context, contact_extid int64) (*entity.CmsPhonecodesMaster, error)
	GetCmsAccountType(ctx context.Context, accountID int64) (*entity.CmsAccountTypeMaster, error)
	GetcmsCoffeetype(ctx context.Context, Id int64) (*entity.CmsCoffeetypeMaster, error)
	GetCmsLeadsBillingAddress(ctx context.Context, leadID string) (*entity.CmsLeadsBillingAddressMaster, error)

	GetLeadCreationInfo(ctx context.Context, req *model.GetLeadCreationInfoRequest) (interface{}, error)

	LeadExists(ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest) (bool, error)

	CreateNewLead(ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest) error

	InsertLeadRecord(ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (string, error)

	AddBillingAddressOnNewLead(leadId string, ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error)

	AddShippingAddressOnNewLead(leadId string, ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error)

	LogToAuditLogOnNewLead(leadId string, ctx context.Context,
		equestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error)

	LogNotificationOnLeadCreateUpdate(leadId string, status string, ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error)

	UpdateLead(ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest) error

	UpdateLeadsShippingaddress(ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error)

	UpdateLeadsBillingaddress(ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error)

	EditAuditLogEntryOnLeadAmend(ctx context.Context,
		requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error)

	ProvideLeadsData(ctx context.Context,
		reqParams model.ProvideLeadsInfoReqContext) ([]model.LeadInfo, error)

	GetDetailsOfLead(ctx context.Context,
		leadDetailsRequest model.GetLeadDetailsRequestBody) (model.LeadDetails, error)

	GetAuditLogEntriesOfLead(ctx context.Context,
		leadId string) ([]model.AuditLogGCPO, error)

	ChangeLeadCreator(ctx context.Context,
		newCreatorUserId string, leadId string) error

	ProvideContactInfoForLead2Aaccount(ctx context.Context,
		leadInfoForL2A model.LeadInfoInLeadToAccount) (model.ContactInfoFromLeadAndMaster, error)
}

type leadRepo struct {
	db *gorm.DB
}

func NewLeadRepo() LeadRepo {
	newDb, err := repositories.NewDb()
	if err != nil {
		panic(err)
	}
	newDb.AutoMigrate(&entity.CmsLeadsMaster{})
	return &leadRepo{
		db: newDb,
	}
}
