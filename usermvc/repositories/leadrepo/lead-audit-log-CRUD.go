package leadrepo

import (
	"context"
	"time"
	"usermvc/model"
	logger2 "usermvc/utility/logger"

	"github.com/jinzhu/gorm"
)

/*
*
 */
func (leadRepoRef leadRepo) LogToAuditLogOnNewLead(leadId string, ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error) {

	insertAuditLogRecordSQL := `INSERT INTO dbo.auditlog_cms_leads_master_newpg(
		leadid, createdby, created_date, description)
		VALUES($1, $2, $3, $4)`

	_, err := db.Raw(insertAuditLogRecordSQL, leadId,
		requestPayload.CreatedUserid, requestPayload.CreatedDate,
		`Lead Created`).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while adding AUDIT LOG record on New Lead Creation in dbo.auditlog_cms_leads_master_newpg ", err.Error())
		return false, err
	}
	return true, nil
}

/*
*
 */
func (leadRepoRef leadRepo) EditAuditLogEntryOnLeadAmend(ctx context.Context,
	requestPayload model.InsertLeadDetailsRequest, db *gorm.DB) (bool, error) {

	updateAuditLogEntrySQL := `update dbo.auditlog_cms_leads_master_newpg
		set	description=$1, modifiedby=$2, modified_date=$3 where leadid=$4`

	_, err := db.Raw(updateAuditLogEntrySQL,
		"Lead Details Modified",
		requestPayload.ModifiedUserid,
		time.Now().Format("2006-01-02"),
		requestPayload.LeadId).Rows()

	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while updating AUDIT LOG entry in dbo.auditlog_cms_leads_master_newpg on Lead Details Update ", err.Error())
		return false, err
	}
	return true, nil
}

/*
*
*
 */
func (leadRepoRef leadRepo) GetAuditLogEntriesOfLead(ctx context.Context,
	leadId string) ([]model.AuditLogGCPO, error) {

	var auditLogEntries []model.AuditLogGCPO
	query := `select u.username as createduser, ld.created_date, ld.description, v.username as modifieduser, ld.modified_date
	   from dbo.auditlog_cms_leads_master_newpg ld inner join dbo.users_master_newpg u on ld.createdby=u.userid
	   left join dbo.users_master_newpg v on ld.modifiedby=v.userid where ld.leadid=$1 order by logid desc`

	rows, err := leadRepoRef.db.Raw(query, leadId).Rows()
	defer rows.Close()
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while fetching audit log entries of Lead: ", leadId, err.Error())
		return nil, err
	}
	for rows.Next() {
		var oneEntry model.AuditLogGCPO
		err = rows.Scan(&oneEntry.CreatedUserid, &oneEntry.CreatedDate, &oneEntry.Description, &oneEntry.ModifiedUserid, &oneEntry.ModifiedDate)
		if err == nil {
			auditLogEntries = append(auditLogEntries, oneEntry)
		}
	}
	return auditLogEntries, err
}
