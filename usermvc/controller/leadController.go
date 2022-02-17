package controller

import (
	"fmt"
	"usermvc/model"
	"usermvc/repositories/accountrepo"
	"usermvc/repositories/leadrepo"
	"usermvc/repositories/userrepo"
	emailer "usermvc/utility/emailer"
	logger2 "usermvc/utility/logger"

	"github.com/gin-gonic/gin"
)

type (
	LeadController interface {
		InsertLeadDetails(ctx *gin.Context)
		GetLeadCreationInfo(ctx *gin.Context)
		GetLeadsInfo(ctx *gin.Context)
		GetLeadDetails(ctx *gin.Context)
		ReassignLead(ctx *gin.Context)
		GetContactsInLead2AccConvert(ctx *gin.Context)
	}

	leadController struct {
		accountrepo accountrepo.AccountRepo
		leadRepo    leadrepo.LeadRepo
	}
)

const (
	ROLE_MANAGING_DIRECTOR          = "Managing Director"
	LEAD_CREATE_REASSIGN_EMAIL_FROM = "itsupport@continental.coffee"
)

/*
*
 */
func NewLeadController() LeadController {
	return &leadController{
		accountrepo: accountrepo.NewAccountRepo(),
		leadRepo:    leadrepo.NewLeadRepo(),
	}
}

/*
*	Handles both CREATE & UPDATE
 */
func (lc leadController) InsertLeadDetails(ctx *gin.Context) {

	var insertLeadDetailsRequestBody model.InsertLeadDetailsRequest
	if err := ctx.ShouldBindJSON(&insertLeadDetailsRequestBody); err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while parsing InsertLeadDetailsRequest", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	var persistenceOpErr error
	if insertLeadDetailsRequestBody.Update {

		persistenceOpErr = lc.leadRepo.UpdateLead(ctx, insertLeadDetailsRequestBody)

	} else {
		var leadExists bool
		leadExists, persistenceOpErr = lc.leadRepo.LeadExists(ctx,
			insertLeadDetailsRequestBody)
		if !leadExists && persistenceOpErr == nil {
			persistenceOpErr = lc.leadRepo.CreateNewLead(ctx, insertLeadDetailsRequestBody)
			go emailConfirmUserAboutTheCreatedLead(ctx, insertLeadDetailsRequestBody)
		} else if leadExists && persistenceOpErr == nil {
			ctx.JSON(230, "Lead Name already exists.")
		}
	}
	if persistenceOpErr != nil {
		ctx.JSON(500, persistenceOpErr.Error())
		return
	}
	ctx.JSON(200, "SUCCESS")
}

/*
*
 */
func emailConfirmUserAboutTheCreatedLead(ctx *gin.Context, requestPayload model.InsertLeadDetailsRequest) {

	userRepo := userrepo.NewUserRepo()
	userInfo, err := userRepo.GetUserInfoByUserId(ctx, requestPayload.CreatedUserid)

	if err == nil {

		dataFeedToTemplate := make(map[string]string)
		dataFeedToTemplate["MessageToUser"] = "New lead has been created."
		dataFeedToTemplate["AccountName"] = requestPayload.Accountname
		dataFeedToTemplate["AccountCountry"] = requestPayload.ShippingCountry
		dataFeedToTemplate["AccountOwner"] = userInfo.Firstname + " " + userInfo.Lastname + " (" + userInfo.Username + ") "

		sendEmailRequestInput := emailer.SendEmailRequestVO{
			SenderDetails: emailer.Sender{
				SendFromIdentity: LEAD_CREATE_REASSIGN_EMAIL_FROM,
			},
			TargetRecipients: emailer.Recipients{
				ToList: []string{requestPayload.CreatorsEmail},
			},
			Message: emailer.EmailBody{
				TemplateForHtmlBody: emailer.EmailTemplate{
					TemplateRef:           "EmailOnLeadCreation",
					DynamicDataOfTemplate: dataFeedToTemplate,
				},
			},
		}
		//var emailRequestAccepted bool
		_, err = emailer.Send(ctx, sendEmailRequestInput)
		//fmt.Println("Email Send Request on New Lead OK? ", emailRequestAccepted)
	}
	if err != nil {
		//fmt.Println("Error while trying to send email confirmation to 'New Lead Creator' ", err.Error())
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while trying to send email confirmation to 'New Lead Creator' ", err.Error())
	}
}

/*
*
 */
func (lc leadController) ReassignLead(ctx *gin.Context) {

	var processingErr error
	var request model.ReassignLeadRequest
	if processingErr = ctx.ShouldBindJSON(&request); processingErr != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Error while parsing InsertLeadDetailsRequest", processingErr.Error())
		ctx.JSON(403, processingErr.Error())
		return
	}
	var assigner model.UserDetails
	userRepo := userrepo.NewUserRepo()
	assigner, processingErr = userRepo.GetUserInfoByUserId(ctx,
		request.ReassignedByUserId)
	if processingErr == nil && assigner.Role == ROLE_MANAGING_DIRECTOR {

		processingErr = lc.leadRepo.ChangeLeadCreator(ctx, request.AssignToUserId, request.LeadId)
		if processingErr == nil {
			var assignee model.UserDetails
			assignee, processingErr = userRepo.GetUserInfoByUserId(ctx,
				request.AssignToUserId)
			if processingErr == nil {
				emailAssigneeOnLeadReassign(ctx, assignee, request.LeadId)
			}
		}
	}
	if processingErr != nil {
		ctx.JSON(500, processingErr.Error())
		return
	} else {
		ctx.JSON(200, "SUCCESS")
	}
}

/*
*
 */
func emailAssigneeOnLeadReassign(ctx *gin.Context,
	user model.UserDetails, leadId string) {

	sendEmailRequestInput := emailer.SendEmailRequestVO{
		SenderDetails: emailer.Sender{
			SendFromIdentity: LEAD_CREATE_REASSIGN_EMAIL_FROM,
		},
		TargetRecipients: emailer.Recipients{
			ToList: []string{user.Emailid},
		},
		Message: emailer.EmailBody{
			BodyText: emailer.Content{
				Data: "Lead with id: " + leadId + " is assigned to " + user.Firstname + " " + user.Lastname + " ( Username: " + user.Username + " ).",
			},
		},
		Subject: emailer.Content{
			Data: "Lead: '" + leadId + "' is reassigned",
		},
	}
	emailRequestAccepted, err := emailer.Send(ctx, sendEmailRequestInput)
	fmt.Println("Reassign Email Request OK? ", emailRequestAccepted)
	if err != nil {
		logger := logger2.GetLoggerWithContext(ctx)
		logger.Error("Lead Reassign Email Send Error: ", err.Error())
	}
}

/*
*
 */
func (controllerRef leadController) GetLeadsInfo(ctx *gin.Context) {

	logger := logger2.GetLoggerWithContext(ctx)
	isAuthorised, _ := func(ctx *gin.Context) (bool, error) {
		authorised := true
		return authorised, nil
	}(ctx)

	if isAuthorised {
		var getLeadsReqPayload model.ProvideLeadsInfoReqContext
		if err := ctx.ShouldBindJSON(&getLeadsReqPayload); err != nil {

			logger.Error("Error while mapping / binding Request Payload to  model.ProvideLeadsInfoReqContext in LeadController ", err.Error())
			ctx.JSON(500, err.Error())
			return
		}
		leadsData, err := controllerRef.leadRepo.ProvideLeadsData(ctx, getLeadsReqPayload)
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.JSON(200, leadsData)
	}
}

/*
* 	may need more* developer testing
*
 */
func (lc leadController) GetLeadCreationInfo(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	validate, err := lc.validateGetLeadCreationInfo(ctx)
	if err != nil {
		logger.Error("error while validating request", err.Error())
		ctx.JSON(503, "error while validating request")
		return
	}
	if !validate {
		logger.Error("request could not be validate")
		ctx.JSON(404, "invalid request")
		return
	}
	var getLeadCreationInfoRequest *model.GetLeadCreationInfoRequest
	if err := ctx.ShouldBindJSON(&getLeadCreationInfoRequest); err != nil {
		logger.Error("error while parsing the GetLeadCreationInfoRequest", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := lc.leadRepo.GetLeadCreationInfo(ctx, getLeadCreationInfoRequest)
	if err != nil {
		logger.Error("error while getting GetLeadCreationInfo", err.Error())
		ctx.JSON(503, err.Error())
		return
	}
	ctx.JSON(200, res)
}

/*
*
 */
func (lc leadController) GetLeadDetails(ctx *gin.Context) {

	logger := logger2.GetLoggerWithContext(ctx)
	validateResult, err := lc.validateGetLeadDetails(ctx)
	if err != nil || !validateResult {
		logger.Error("Request Validation failed: ", err.Error())
		ctx.JSON(403, "Invalid Request.")
		return
	}
	var getLeadDetailsRequestBody model.GetLeadDetailsRequestBody
	if err := ctx.ShouldBindJSON(&getLeadDetailsRequestBody); err != nil {
		logger.Error("Error while parsing getLeadDetailsRequest", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	leadDetailResp, err := lc.leadRepo.GetDetailsOfLead(ctx, getLeadDetailsRequestBody)
	if err != nil {
		logger.Error("Error occured in LeadController -> GetLeadDetails:  ", err.Error())
		ctx.JSON(500, err.Error())
		return
	}
	ctx.JSON(200, leadDetailResp)
}

/*
*
 */
func (lc leadController) GetContactsInLead2AccConvert(ctx *gin.Context) {

	logger := logger2.GetLoggerWithContext(ctx)
	var leadInfoForL2AConvert model.LeadInfoInLeadToAccount
	if err := ctx.ShouldBindJSON(&leadInfoForL2AConvert); err != nil {
		logger.Error("Error while binding request params to LeadInfoInLeadToAccount in LeadController->GetContactsInLead2AccConvert ", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	contactsFromLeadAndMaster, err := lc.leadRepo.ProvideContactInfoForLead2Aaccount(ctx, leadInfoForL2AConvert)
	if err != nil {
		logger.Error("Error occured in LeadController -> GetContactsInLead2AccConvert: ", err.Error())
		ctx.JSON(500, err.Error())
		return
	} else {
		ctx.JSON(200, contactsFromLeadAndMaster)
	}
}

/*
func (lc leadController) GetAllLeadsDetails(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	res, err := lc.leadRepo.GetAllLeadDetails(ctx)
	if err != nil {
		logger.Error("error while getting all leadaccounts", err.Error())
		ctx.JSON(503, err.Error())
	}
	ctx.JSON(200, res)
}*/

func (lc leadController) validateGetLeadsInfo(ctx *gin.Context) (bool, error) {
	return true, nil
}

func (lc leadController) validateGetLeadCreationInfo(ctx *gin.Context) (bool, error) {
	return true, nil
}

func (lc leadController) validateGetLeadDetails(ctx *gin.Context) (bool, error) {
	return true, nil
}
