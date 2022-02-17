package emailer

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	apputils "usermvc/utility/apputils"
)

/*
*
 */
func buildRawMessage(sendEmailRequestVO SendEmailRequestVO) (string, error) {

	var buildErr error = nil
	var builder strings.Builder
	var hasInlineMedia bool
	var hasAttachments bool
	var msgPart string
	var rootBoundaryId string = ""
	builder.WriteString(buildTopHeaders(sendEmailRequestVO))
	if len(sendEmailRequestVO.Attachments) > 0 {
		rootBoundaryId, msgPart = buildMixedHeader()
		builder.WriteString(msgPart)
		hasAttachments = true
	}
	if emptyBody, emptyBodyPart := processEmptyBody(sendEmailRequestVO, rootBoundaryId); emptyBody {
		builder.WriteString(emptyBodyPart)
		//emptyBody, emptyBodyPart := processEmptyBody(sendEmailRequestVO, rootBoundaryId)
		//if emptyBody {
		//builder.WriteString(emptyBodyPart)
	} else {
		keyAndUrlMap, keyAndCidMap := extractInlineImageInfoFromRequest(sendEmailRequestVO)
		if keyAndUrlMap == nil {
			keyAndUrlMap, keyAndCidMap = pepareAndCheckInlineImageInfo(&sendEmailRequestVO)
		}
		if keyAndUrlMap != nil {
			hasInlineMedia = true
		}
		var altBoundary string
		var relBoundary string
		if sendEmailRequestVO.Message.HasBothTextAndHtml {
			msgPart, altBoundary = buildAlternativeHeader(rootBoundaryId)
			builder.WriteString(msgPart)
		}
		if altBoundary == "" {
			builder.WriteString(buildTextBodyPart(sendEmailRequestVO, rootBoundaryId, keyAndCidMap, hasInlineMedia))
		} else {
			builder.WriteString(buildTextBodyPart(sendEmailRequestVO, altBoundary, keyAndCidMap, hasInlineMedia))
		}
		msgPart, relBoundary = processHtmlBody(sendEmailRequestVO, hasInlineMedia, rootBoundaryId, altBoundary, keyAndCidMap, keyAndUrlMap)
		builder.WriteString(msgPart)
		msgPart, buildErr = processInlineImages(sendEmailRequestVO, hasInlineMedia, relBoundary, keyAndCidMap, keyAndUrlMap)
		if buildErr == nil && msgPart != "" {
			builder.WriteString(msgPart)
		}
		if altBoundary != "" && buildErr == nil {
			builder.WriteString("\n--" + altBoundary + "--\n")
		}
	}
	if hasAttachments && buildErr == nil {
		msgPart, buildErr = buildAttachmentsPart(sendEmailRequestVO, "--"+rootBoundaryId)
		if buildErr == nil {
			builder.WriteString(msgPart)
		}
	}
	if rootBoundaryId != "" && buildErr == nil {
		builder.WriteString("\n--" + rootBoundaryId + "--")
	}
	if buildErr != nil {
		return "", buildErr
	}
	return builder.String(), nil
}

/*
*
 */
func processEmptyBody(sendEmailRequestVO SendEmailRequestVO, parentBoundary string) (bool, string) {

	body := sendEmailRequestVO.Message
	if (body.BodyHtml == Content{} && body.BodyText == Content{}) {
		var builder strings.Builder
		if parentBoundary != "" {
			builder.WriteString("--" + parentBoundary + "\n")
		}
		builder.WriteString("Content-Type: text/plain; charset=" + CHARSET_UTF_8 + "\n")
		builder.WriteString("Content-Length: 0\n\n")
		return true, builder.String()
	}
	return false, ""
}

/*
*
 */
func processHtmlBody(sendEmailRequestVO SendEmailRequestVO,
	hasInlineMedia bool, rootBoundary string, altBoundary string,
	keyAndCidMap map[string]string, keyAndUrlMap map[string]string) (string, string) {

	var builder strings.Builder
	var msgPart string
	var relBoundary = ""
	if hasInlineMedia {
		if altBoundary == "" {
			msgPart, relBoundary = buildRelatedHeader(rootBoundary)
		} else {
			msgPart, relBoundary = buildRelatedHeader(altBoundary)
		}
		builder.WriteString(msgPart)
		builder.WriteString(buildHtmlBodyPart(sendEmailRequestVO, relBoundary, keyAndCidMap, hasInlineMedia))
	} else if altBoundary != "" {
		builder.WriteString(buildHtmlBodyPart(sendEmailRequestVO, altBoundary, keyAndCidMap, hasInlineMedia))
	} else {
		builder.WriteString(buildHtmlBodyPart(sendEmailRequestVO, rootBoundary, keyAndCidMap, hasInlineMedia))
	}
	return builder.String(), relBoundary
}

/*
*
 */
func processInlineImages(sendEmailRequestVO SendEmailRequestVO, hasInlineMedia bool,
	relBoundary string, keyAndCidMap map[string]string, keyAndUrlMap map[string]string) (string, error) {

	var processingErr error = nil
	if hasInlineMedia {

		var builder strings.Builder
		var oneInlineImgPart string
		for key, media := range sendEmailRequestVO.Message.DynamicDataOfMediaInBody {
			oneInlineImgPart, processingErr = buildPartForOneInlineImage(media, keyAndCidMap[key], keyAndUrlMap[key], relBoundary)
			if processingErr != nil {
				break
			}
			builder.WriteString(oneInlineImgPart)
		}
		if processingErr == nil {
			builder.WriteString("\n--" + relBoundary + "--\n")
			return builder.String(), nil
		}
	}
	return "", processingErr
}

/*
*
 */
func buildPartForOneInlineImage(mediaInfo Content, contentId string,
	imgUrl string, boundary string) (string, error) {

	var builder strings.Builder
	var processingErr error = nil
	var imgAsBase64 string
	imgAsBase64, processingErr = apputils.Base64OfImageFromUrl(mediaInfo.Url)
	if processingErr == nil {
		var mediaType string = mediaInfo.ContentType // image/png, image/jpeg
		fileName := fmt.Sprint(apputils.Crc32OfString(contentId)) + "." + mediaType[6:]

		builder.WriteString("--" + boundary + "\n")
		builder.WriteString("Content-Type: " + mediaType + "; name=\"" + fileName + "\"\n")
		builder.WriteString("Content-Transfer-Encoding: base64\n")
		builder.WriteString("Content-Disposition: inline; filename=\"" + fileName + "\"\n")
		builder.WriteString("Content-ID: <" + contentId + ">\n")
		builder.WriteString("X-Attachment-Id: " + contentId + "\n\n")
		//imgAsBase64 = "imgAsBase64"
		builder.WriteString(imgAsBase64)
		builder.WriteString("\n")
	}
	return builder.String(), processingErr
}

/*
*
 */
func buildHtmlBodyPart(sendEmailRequestVO SendEmailRequestVO, boundary string,
	mediaContentKeyAndIdMap map[string]string, hasInlineMedia bool) string {

	var builder strings.Builder
	htmlBody := sendEmailRequestVO.Message.BodyHtml
	if (htmlBody != Content{}) {
		if boundary != "" {
			builder.WriteString("--" + boundary + "\n")
		}
		builder.WriteString("Content-Type: " + CONTENT_TYPE_TEXT_HTML + "; charset=" + htmlBody.CharSet + "\n")
		builder.WriteString("Content-Transfer-Encoding: quoted-printable\n\n")
		if hasInlineMedia {
			htmlWithContentIds := apputils.ReplaceKeyWithPrefixedValues(htmlBody.Data,
				"{{", "}}", mediaContentKeyAndIdMap, "cid:")
			builder.WriteString(htmlWithContentIds + "\n")
		} else {
			builder.WriteString(htmlBody.Data + "\n")
		}
		return builder.String()
	}
	return ""
}

/*
*
 */
func buildTextBodyPart(sendEmailRequestVO SendEmailRequestVO, boundary string,
	mediaContentKeyAndIdMap map[string]string, hasInlineMedia bool) string {

	var builder strings.Builder
	textBody := sendEmailRequestVO.Message.BodyText
	if (textBody != Content{}) {
		if boundary != "" {
			builder.WriteString("--" + boundary + "\n")
		}
		builder.WriteString("Content-Type: " + CONTENT_TYPE_TEXT_PLAIN + "; charset=" + textBody.CharSet + "\n")
		builder.WriteString("Content-Transfer-Encoding: quoted-printable\n\n")
		builder.WriteString(textBody.Data + "\n")
		return builder.String()
	}
	return ""
}

/*
*

func buildPartForInlineImages(mediaInfo map[string]Content,
	keyContentIdMap map[string]string,
	keyUrlMap map[string]string, boundary string) (string, error) {

	var builder strings.Builder
	var processingErr error = nil
	var oneInlineImgPart string
	for key, media := range mediaInfo {
		oneInlineImgPart, processingErr = buildPartForOneInlineImage(media, keyContentIdMap[key], keyUrlMap[key], boundary)
		if processingErr != nil {
			break
		}
		builder.WriteString(oneInlineImgPart)
	}
	//builder.WriteString(boundary + "--\n\n")
	return builder.String(), processingErr
} */

/*
*
 */
func buildAttachmentsPart(sendEmailRequestVO SendEmailRequestVO, boundary string) (string, error) {

	var builder strings.Builder
	var processingErr error = nil
	var oneAttachmentPartOfRawMsg string
	attachments := sendEmailRequestVO.Attachments
	for _, attachment := range attachments {
		oneAttachmentPartOfRawMsg, processingErr = buildPartForOneAttachment(attachment, boundary)
		if processingErr != nil {
			break
		}
		builder.WriteString(oneAttachmentPartOfRawMsg)
	}
	return builder.String(), processingErr
}

/*
*
 */
func buildPartForOneAttachment(attachment Attachment, boundary string) (string, error) {

	var ioErr error
	var attachmentContentString string
	if attachment.FQPath != "" {
		attachmentContentString, ioErr = apputils.Base64OfFileContent(attachment.FQPath)
		if ioErr != nil {
			return "", ioErr
		}
	} else if len(attachment.ContentBytes) > 0 {
		attachmentContentString = base64.StdEncoding.EncodeToString(attachment.ContentBytes)
	} else if attachment.ContentString != "" {
		attachmentContentString = base64.StdEncoding.EncodeToString([]byte(attachment.ContentString))
	} else {
		return "", AttachmentContentAccessErr
	}
	var builder strings.Builder
	builder.WriteString(boundary + "\n")
	builder.WriteString("Content-Type: " + attachment.ContentType + "; name=\"" + attachment.Filename + "\"\n")
	builder.WriteString("Content-Disposition: attachment;filename=\"" + attachment.Filename + "\"\n")
	builder.WriteString("Content-Transfer-Encoding: base64\n")
	builder.WriteString("X-Attachment-Id: " + fmt.Sprint(apputils.Crc32OfString(attachment.Filename+fmt.Sprint(time.Now().UnixNano()))) + "\n\n")
	//attachmentContentString = "attachmentContentString"
	builder.WriteString(attachmentContentString)
	builder.WriteString("\n")
	return builder.String(), nil
}

/*
*
 */
func buildTopHeaders(sendEmailRequestVO SendEmailRequestVO) string {

	var builder strings.Builder
	builder.WriteString("Message-ID: <" + generateMessageId(sendEmailRequestVO.TargetRecipients.ToList[0]) + ">\n")
	builder.WriteString("From: '" + sendEmailRequestVO.SenderDetails.DisplayFromName + "' <" + sendEmailRequestVO.SenderDetails.SendFromIdentity + ">\n")
	builder.WriteString("Subject: " + sendEmailRequestVO.Subject.Data + "\n")
	builder.WriteString("To: " + strings.Join(sendEmailRequestVO.TargetRecipients.ToList[:], ",") + "\n")
	builder.WriteString("MIME-Version: " + sendEmailRequestVO.MimeVersion + "\n")
	return builder.String()
}

/*
*
 */
func buildMixedHeader() (string, string) {

	var builder strings.Builder
	rootBoundaryId := generateBoundaryId("MESSAGE-WITH-ATTACHMENTs")
	builder.WriteString("Content-Type: multipart/mixed; boundary=\"" + rootBoundaryId + "\"\n\n")
	//builder.WriteString("--" + rootBoundaryId + "\n")
	return rootBoundaryId, builder.String()
}

/*
*
 */
func buildRelatedHeader(parentBoundary string) (string, string) {

	var builder strings.Builder
	var relBoundary string = ""
	if parentBoundary == "" {
		relBoundary = generateBoundaryId("multipart/related")
	} else {
		builder.WriteString("--" + parentBoundary + "\n")
		relBoundary = "rel_" + parentBoundary
	}
	builder.WriteString("Content-Type: multipart/related; boundary=\"" + relBoundary + "\"\n\n")
	//builder.WriteString("--" + relBoundary + "\n")
	return builder.String(), relBoundary
}

/*
*
 */
func buildAlternativeHeader(parentBoundary string) (string, string) {

	var builder strings.Builder
	var altBoundary string
	if parentBoundary == "" {
		altBoundary = generateBoundaryId("multipart/alternative")
	} else {
		builder.WriteString("--" + parentBoundary + "\n")
		altBoundary = "alt_" + parentBoundary
	}
	builder.WriteString("Content-Type: multipart/alternative; boundary=\"" + altBoundary + "\"\n\n")
	//builder.WriteString("--" + altBoundary + "\n")
	return builder.String(), altBoundary
}

/*
*
 */
func generateBoundaryId(inputHint string) string {
	return fmt.Sprint(apputils.Crc32OfString(inputHint))
}

/*
*
 */
func generateContentId(inputHint string) string {
	nanoS := fmt.Sprint(time.Now().UnixNano())
	return fmt.Sprint(apputils.Crc32OfString(inputHint+nanoS)) + "@" + nanoS
}

/*
*
 */
func generateMessageId(inputHint string) string {
	nanoS := fmt.Sprint(time.Now().UnixNano())
	return fmt.Sprint(apputils.Crc32OfString(inputHint+nanoS)) + "@" + nanoS
}

/*
*
 */
func pepareAndCheckInlineImageInfo(sendEmailRequestVO *SendEmailRequestVO) (map[string]string, map[string]string) {

	body := sendEmailRequestVO.Message
	htmlBody := body.BodyHtml
	if (htmlBody != Content{}) {
		var html = htmlBody.Data
		imsgSrcs := apputils.AllMatchesWithRegex(`\bsrc=["']([^"']+)["']`, html)
		var processedImgUrl = make(map[string]string)
		var processedImgSrcMap = make(map[string]string)
		var oneImgUrl = ""
		var keyForImgUrl = ""
		mediaDataMap := make(map[string]Content)
		for _, oneImgSrc := range imsgSrcs {

			if processedImgSrcMap[oneImgSrc] == "" {
				oneImgUrl = apputils.AllSubMatchesWithRegex(`src=["']+([^"']+)["']`, oneImgSrc)[0]
				if processedImgUrl[oneImgUrl] == "" {
					keyForImgUrl = "KEY_" + fmt.Sprint(apputils.Crc32OfString(oneImgUrl))
					//fmt.Println("pepareAndCheckInlineImageInfo -> replace ->", oneImgSrc)
					//fmt.Println("pepareAndCheckInlineImageInfo -> replace with ->", "src=3D=22{{"+keyForImgUrl+"}}=22")
					html = strings.Replace(html, oneImgSrc, "src=3D=22{{"+keyForImgUrl+"}}=22", -1)

					mediaContentImage := Content{
						ContentType: "image/" + oneImgUrl[strings.LastIndex(oneImgUrl, ".")+1:],
						Url:         oneImgUrl,
					}
					mediaDataMap[keyForImgUrl] = mediaContentImage
					processedImgUrl[oneImgUrl] = "Y"
				}
				processedImgSrcMap[oneImgSrc] = "Y"
			}
		}
		if len(imsgSrcs) > 0 {
			htmlBody.Data = html
			//fmt.Println("pepareAndCheckInlineImageInfo -> html ->", html)
			body.BodyHtml = htmlBody
			body.DynamicDataOfMediaInBody = mediaDataMap
			sendEmailRequestVO.Message = body
			//fmt.Println("pepareAndCheckInlineImageInfo -> vo html ->", sendEmailRequestVO.Message.BodyHtml.Data)
			//sendEmailRequestVO.Template.TemplateMediaDataFeed = mediaDataToEmplate
			//sendEmailRequestVO.Message.DynamicDataOfMediaInBody = mediaDataMap
			return extractInlineImageInfoFromRequest(*sendEmailRequestVO)
		}
	}
	return nil, nil
}
