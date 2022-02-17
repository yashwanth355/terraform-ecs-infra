package emailer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	logger2 "usermvc/utility/logger"

	"go.uber.org/zap"
)

type (
	SendingInfraServiceStub struct {
		SendServiceRef                string                            //AWS_SES
		ServiceConsumptionInputParams *InputConfigForServiceConsumption // input params
		ServiceLogger                 ServiceLogger
	}
	InputConfigForServiceConsumption struct {
		AllParams          map[string]string
		CredentialsInfoMap map[string]string
	}
	ServiceLogger struct {
		Logger *zap.SugaredLogger
	}
)

/*
*	Internal Known Errors ( implementation related )
*	that may occur
*
 */
var ServiceStubBuildingErr = errors.New("Fatal - Could not build a Stub for the Email Sending Infrastructre Service.")

/*
*
 */
func (sendingServiceRef SendingInfraServiceStub) SendEmail(ctx context.Context,
	sendEmailRequestVO SendEmailRequestVO) (bool, error) {

	var isSuccess bool = false
	var processingErr error
	var serviceStub SendingInfraServiceStub
	serviceStub, processingErr = getServiceStubInstance(ctx, sendEmailRequestVO)
	// log info for observing, troubleshooting
	logDebugInfo(sendEmailRequestVO, serviceStub)
	if processingErr == nil {
		processingErr = validateSanitize(&sendEmailRequestVO, serviceStub)
		if processingErr == nil {
			// SEND, picking the right EMAIL SENDING INFRA SERVICE PROVIDER e.g. AWS-SES
			isSuccess, processingErr = sendUsingRelevantProvider(sendEmailRequestVO, serviceStub)
		}
	}
	return isSuccess, processingErr
}

/*
*
 */
func getServiceStubInstance(ctx context.Context, sendEmailRequestVO SendEmailRequestVO) (SendingInfraServiceStub, error) {

	var stub SendingInfraServiceStub
	sendThruServiceRef := sendEmailRequestVO.SendUsingService
	emailTypeObj := &sendEmailRequestVO.Type

	if sendThruServiceRef == "" {
		if emailTypeObj != nil {
			sendThruServiceRef = emailTypeObj.SendThoughServiceRef
		}
		if sendThruServiceRef == "" {
			sendThruServiceRef = EMAIL_SENDER_SERVICE_AWS_SES
		}
	}
	consumptionConfig, configErr := buildConfigForServiceConsumption(sendEmailRequestVO, sendThruServiceRef)

	if configErr != nil {
		return stub, ServiceStubBuildingErr
	}

	stub = SendingInfraServiceStub{
		SendServiceRef:                sendThruServiceRef,
		ServiceConsumptionInputParams: consumptionConfig,
		ServiceLogger: ServiceLogger{
			Logger: logger2.GetLoggerWithContext(ctx),
		},
	}
	return stub, nil
}

/*
*
 */
func sendUsingRelevantProvider(sendEmailRequestVO SendEmailRequestVO, serviceStub SendingInfraServiceStub) (bool, error) {

	var sendSucceeded bool = false
	var processingErr error

	switch serviceStub.SendServiceRef {

	case EMAIL_SENDER_SERVICE_AWS_SES:

		sendSucceeded, processingErr = sendWithAwsSes(sendEmailRequestVO, serviceStub.ServiceConsumptionInputParams)
	default:
		processingErr = UnknownEmailSendingInfraServiceErr
	}
	return sendSucceeded, processingErr
}

/*
*
 */
func buildConfigForServiceConsumption(sendEmailRequestVO SendEmailRequestVO, sendThruServiceRef string) (*InputConfigForServiceConsumption, error) {

	var inputParamsToConsumeService *InputConfigForServiceConsumption

	var credentialsInfoMap = make(map[string]string)
	var allParamsMap = make(map[string]string)

	switch sendThruServiceRef {

	case EMAIL_SENDER_SERVICE_AWS_SES:

		inputParamsToConsumeService, _ = buildConsumptionConfigForAwsSes(sendEmailRequestVO, allParamsMap, credentialsInfoMap)

	default:
		return nil, UnknownEmailSendingInfraServiceErr
	}
	return inputParamsToConsumeService, nil
}

/*
*
 */
func logDebugInfo(sendEmailRequestVO SendEmailRequestVO, serviceStub SendingInfraServiceStub) {

	var servicesLogger ServiceLogger = serviceStub.ServiceLogger

	servicesLogger.Logger.Info("[[Email Sender]] service base impl -> will send through Provider -> ", serviceStub.SendServiceRef)

	fmt.Println("[[Email Sender]] service base impl -> will send through Provider -> ", serviceStub.SendServiceRef)
}

/*
*
*	to do: attachments total size validation (nice to have)
*
*
 */
func validateSanitize(sendEmailRequestVO *SendEmailRequestVO, serviceStub SendingInfraServiceStub) error {

	var sanityCheckErr error
	if sendEmailRequestVO.MimeVersion == "" {
		sendEmailRequestVO.MimeVersion = "1.0"
	}
	//template := sendEmailRequestVO.Template
	htmlTemplate := sendEmailRequestVO.Message.TemplateForHtmlBody
	//if template.TemplateRef == "" {
	if htmlTemplate.TemplateRef == "" {
		checkBody(&sendEmailRequestVO.Message)

		//} else if !template.DoesNotNeedsDataPopulating && len(template.TemplateDataFeed) == 0 {
	} else if !htmlTemplate.DoesNotNeedsDataPopulating && len(htmlTemplate.DynamicDataOfTemplate) == 0 {
		sanityCheckErr = NoTemplatePopulatingDataErr
	}
	if len(sendEmailRequestVO.TargetRecipients.ToList) == 0 {
		sanityCheckErr = NoTargetRecipientsErr
	}
	senderDetails := &sendEmailRequestVO.SenderDetails
	sanityCheckErr = checkSenderDetails(senderDetails)
	if sanityCheckErr != nil {
		return sanityCheckErr
	}
	return nil
}

/*
*
 */
func checkSenderDetails(senderDetails *Sender) error {

	var sanityCheckErr error
	if senderDetails.SendFromIdentity == "" {
		sanityCheckErr = NoSendFromIdentityErr
	}
	if senderDetails.ReplyToIdentity == "" {
		senderDetails.ReplyToIdentity = senderDetails.SendFromIdentity
	}
	if senderDetails.DisplayFromName == "" {
		senderDetails.DisplayFromName = senderDetails.SendFromIdentity
	}
	if sanityCheckErr != nil {
		return sanityCheckErr
	}
	return nil
}

/*
*
 */
func checkBody(body *EmailBody) error {

	var sanityCheckErr error = nil

	if body.BodyText.Data != "" && body.BodyText.CharSet == "" {
		body.BodyText.CharSet = CHARSET_UTF_8
		body.BodyText.ContentType = CONTENT_TYPE_TEXT_PLAIN
	}
	if body.BodyHtml.Data != "" && body.BodyHtml.CharSet == "" {
		body.BodyText.CharSet = CHARSET_UTF_8
		body.BodyHtml.ContentType = CONTENT_TYPE_TEXT_HTML
	}
	var containsBothTextNHtml bool = false
	if (body.BodyText != Content{} && body.BodyHtml != Content{}) || (body.BodyText.Data != "" && body.BodyHtml.Data != "") {
		containsBothTextNHtml = true
	}
	body.HasBothTextAndHtml = containsBothTextNHtml
	return sanityCheckErr
}

/*
*
 */
func createDummyBodyIfEmpty(body EmailBody, createEmptyTextBody bool) EmailBody {

	emptyTxtBody := Content{
		CharSet:     CHARSET_UTF_8,
		ContentType: CONTENT_TYPE_TEXT_PLAIN,
		Data:        "",
	}
	if (body.BodyHtml == Content{} && body.BodyText == Content{}) {
		fmt.Println("Empty Email Body..")
		body = EmailBody{
			BodyText: emptyTxtBody,
		}
	} else if (createEmptyTextBody && body.BodyText == Content{}) {
		fmt.Println("Create Empty Text Body..")
		body.BodyText = emptyTxtBody
	}
	return body
}

/*
*
 */
func makeEmailBodyHtml(html string, charset string) EmailBody {

	if charset == "" {
		charset = CHARSET_UTF_8
	}
	var body EmailBody
	body = EmailBody{
		BodyHtml: Content{
			CharSet: charset,
			Data:    html,
		},
	}
	return body
}

/*
*
 */
func extractInlineImageInfoFromRequest(sendEmailRequestVO SendEmailRequestVO) (map[string]string, map[string]string) {

	var keyUrlMap map[string]string = nil
	var keyContentIdMap map[string]string = nil
	//mediaDataContentMap := sendEmailRequestVO.Template.TemplateMediaDataFeed
	mediaDataContentMap := sendEmailRequestVO.Message.DynamicDataOfMediaInBody
	if len(mediaDataContentMap) > 0 {
		keyUrlMap = make(map[string]string)
		keyContentIdMap = make(map[string]string)
		for k, v := range mediaDataContentMap {
			//fmt.Println("extractInlineImageInfoFromRequest -> mediaDataContentMap -> key: ", k)
			//fmt.Println("extractInlineImageInfoFromRequest -> mediaDataContentMap -> ContentType of value: ", v.ContentType)
			//fmt.Println("extractInlineImageInfoFromRequest -> mediaDataContentMap -> Url of value: ", v.Url)
			if strings.HasPrefix(v.ContentType, "image") {
				if v.Url != "" {
					keyUrlMap[k] = v.Url
					keyContentIdMap[k] = generateContentId(k)
				}
			}
		}
	}
	return keyUrlMap, keyContentIdMap
}
