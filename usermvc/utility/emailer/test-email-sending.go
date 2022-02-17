package emailer

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
)

/*
*
*	Executed through SES API based - simpleMessage
*
 */
func TestSendingSimpleEmail_TextBody(ctx context.Context) {

	sendEmailRequestInput := SendEmailRequestVO{
		Type: EmailType{
			CallerITInfraRegion: "ap-south-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "itsupport@continental.coffee",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			BodyText: Content{
				CharSet: "iso-8859-1",
				Data:    "TestSendingSimpleEmail_TextBody - Test Body Text",
			},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingSimpleEmail_TextBody",
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}
}

/*
*
*	Executed through SES API based - simpleMessage
*
 */
func TestSendingSimpleEmail_HtmlBody(ctx context.Context) {

	sendEmailRequestInput := SendEmailRequestVO{
		Type: EmailType{
			CallerITInfraRegion: "ap-south-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "itsupport@continental.coffee",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			BodyHtml: Content{
				Data: "<!DOCTYPE html><html><head><img src='https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>",
			},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingSimpleEmail_HtmlBody",
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}
}

/*
*	Executed through SES API based - sendTemplatedMessage
*
*	** can provide IMAGE URLs directly in the template **
 */
func TestSendingTemplatedEmailWithoutAttachments(ctx context.Context) {

	dataFeedToTemplate := make(map[string]string)
	dataFeedToTemplate["MessageToUser"] = "TestSendingTemplatedEmailWithoutAttachments - SES Api for Templated Emails"
	dataFeedToTemplate["AccountName"] = "Prospect Name"
	dataFeedToTemplate["AccountCountry"] = "India"
	dataFeedToTemplate["AccountOwner"] = "Banani"

	sendEmailRequestInput := SendEmailRequestVO{
		Type: EmailType{
			CallerITInfraRegion: "ap-south-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "itsupport@continental.coffee",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			TemplateForHtmlBody: EmailTemplate{
				TemplateRef:           "EmailOnLeadCreation",
				DynamicDataOfTemplate: dataFeedToTemplate,
			},
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success -> ", emailRequestProcessingSuccess)
	}
}

/*
* Tested with Both RawMessage and ses api based simpleMessage
*
Should build raw message like below. ** But it's happening through SES API => sendSimpleMessage
*
*
Message-ID: <4078219069@1644820124930425000>
From: 'deb.work.related@gmail.com' <deb.work.related@gmail.com>
Subject: TestSendingRawMessege_TextBodyOnly_Without_Attachments
To: debdas.sinha@gmail.com,debdas.sinha@outlook.com
MIME-Version: 1.0
Content-Type: text/plain; charset=iso-8859-1
Content-Transfer-Encoding: quoted-printable

Text body of TestSendingRawMessege_TextBodyOnly_Without_Attachments

*
*/
func TestSendingRawMessege_TextBodyOnly_Without_Attachments(ctx context.Context) {

	sendEmailRequestInput := SendEmailRequestVO{
		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "deb.work.related@gmail.com",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			BodyText: Content{
				CharSet: "iso-8859-1",
				Data:    "Text body of TestSendingRawMessege_TextBodyOnly_Without_Attachments",
			},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingRawMessege_TextBodyOnly_Without_Attachments",
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}
}

/* ******* Needs logic to extract image urls and treat as embedded / inline ones *****
*
*
* Tested with Both RawMessage and ses api based simpleMessage
*
Should build raw message like below. ** But it's happening through SES API => sendSimpleMessage
*
*

Message-ID: <2585050284@1644819750540919000>
From: 'deb.work.related@gmail.com' <deb.work.related@gmail.com>
Subject: TestSendingRawMessege_HtmlBodyOnly_Without_Attachments
To: debdas.sinha@gmail.com,debdas.sinha@outlook.com
MIME-Version: 1.0
Content-Type: text/html; charset=iso-8859-1
Content-Transfer-Encoding: quoted-printable

<!DOCTYPE html><html><head><img src='https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>

*
*/
func TestSendingRawMessege_HtmlBodyOnly_Without_Attachments(ctx context.Context) {

	sendEmailRequestInput := SendEmailRequestVO{
		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "deb.work.related@gmail.com",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			BodyHtml: Content{
				CharSet: "iso-8859-1",
				Data:    "<!DOCTYPE html><html><head><img src='https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>",
			},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingRawMessege_HtmlBodyOnly_Without_Attachments",
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}
}

/* ******* Needs logic to extract image urls and treat as embedded / inline ones *****
*
*
* Tested with Both RawMessage and ses api based simpleMessage (happens through simpleMessage)
*
* Should build raw message like below. ** But it's happening through SES API => sendSimpleMessage
*
Message-ID: <3609724441@1644819388552066000>
From: 'deb.work.related@gmail.com' <deb.work.related@gmail.com>
Subject: TestSendingRawMessege_BothTextAndHtmlBody_Without_Attachments
To: debdas.sinha@gmail.com,debdas.sinha@outlook.com
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="1867535975"

--1867535975
Content-Type: text/plain; charset=iso-8859-1
Content-Transfer-Encoding: quoted-printable

Text body of TestSendingRawMessege_BothTextAndHtmlBody_Without_Attachments
--1867535975
Content-Type: text/html; charset=iso-8859-1
Content-Transfer-Encoding: quoted-printable

<!DOCTYPE html><html><head><img src='https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>

--1867535975--


*
*/
func TestSendingRawMessege_BothTextAndHtmlBody_Without_Attachments(ctx context.Context) {

	sendEmailRequestInput := SendEmailRequestVO{
		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "deb.work.related@gmail.com",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			BodyText: Content{
				CharSet: "iso-8859-1",
				Data:    "Text body of TestSendingRawMessege_BothTextAndHtmlBody_Without_Attachments",
			},
			BodyHtml: Content{
				CharSet: "iso-8859-1",
				Data:    "<!DOCTYPE html><html><head><img src='https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>",
			},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingRawMessege_BothTextAndHtmlBody_Without_Attachments",
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}
}

/* ******* Needs logic to extract image urls and treat as embedded / inline ones *****
*
Message-ID: <2577986150@1644819062510008000>
From: 'deb.work.related@gmail.com' <deb.work.related@gmail.com>
Subject: TestSendingRawMessege_BothTextAndHtmlBody_With_Attachments
To: debdas.sinha@gmail.com,debdas.sinha@outlook.com
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="3330259215"

--3330259215
Content-Type: multipart/alternative; boundary="alt_3330259215"

--alt_3330259215
Content-Type: text/plain; charset=iso-8859-1
Content-Transfer-Encoding: quoted-printable

Text body of TestSendingRawMessege_BothTextAndHtmlBody_With_Attachments
--alt_3330259215
Content-Type: text/html; charset=iso-8859-1
Content-Transfer-Encoding: quoted-printable

<!DOCTYPE html><html><head><img src='https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>

--alt_3330259215--
--3330259215
Content-Type: text/plain; name="test2.txt"
Content-Disposition: attachment;filename="test2.txt"
Content-Transfer-Encoding: base64
X-Attachment-Id: 1589160316

RGV2T3BzIDIKCgpQdWJsaWMgUmVwb3MgMg==

--3330259215--

*
*
*/
func TestSendingRawMessege_BothTextAndHtmlBody_With_Attachments(ctx context.Context) {

	fmt.Println("----------------------\n")
	dirOfAttachedFile, err := os.Getwd()
	if err != nil {
		fmt.Println("FATAL: Error Getting Directory of current executable: ", err)
	}
	fmt.Println("Directory where attachment file is put for testing: ", dirOfAttachedFile)
	fmt.Println("----------------------\n")

	contentBytes, err := ioutil.ReadFile(dirOfAttachedFile + "/test2.txt")
	if err != nil {
		panic(err)
		return
	}
	attachment1 := Attachment{
		Filename:     "test2.txt",
		ContentType:  ATTACHMENT_CONTENT_TYPE_PLAIN,
		ContentBytes: contentBytes,
	}
	attachments := []Attachment{attachment1}

	sendEmailRequestInput := SendEmailRequestVO{

		Attachments: attachments,

		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "deb.work.related@gmail.com",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			BodyText: Content{
				CharSet: "iso-8859-1",
				Data:    "Text body of TestSendingRawMessege_BothTextAndHtmlBody_With_Attachments",
			},
			BodyHtml: Content{
				CharSet: "iso-8859-1",
				Data:    "<!DOCTYPE html><html><head><img src='https://storage.googleapis.com/ustat-ww/ws/mindlounge.co.in/img/ID/ml.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>",
			},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingRawMessege_BothTextAndHtmlBody_With_Attachments",
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}
}

/* ******* Needs logic to extract image urls and treat as embedded / inline ones *****
*

Message-ID: <1922798977@1644818388584705000>
From: 'deb.work.related@gmail.com' <deb.work.related@gmail.com>
Subject: TestSendingRawMessege_HtmlBody_With_Attachments
To: debdas.sinha@gmail.com,debdas.sinha@outlook.com
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="3330259215"

--3330259215
Content-Type: text/html; charset=iso-8859-1
Content-Transfer-Encoding: quoted-printable

<!DOCTYPE html><html><head><img src='https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>
--3330259215
Content-Type: application/pdf; name="test2.txt"
Content-Disposition: attachment;filename="test2.txt"
Content-Transfer-Encoding: base64
X-Attachment-Id: 3250962856

RGV2T3BzIDIKCgpQdWJsaWMgUmVwb3MgMg==

--3330259215--

*
*/
func TestSendingRawMessege_HtmlBody_With_Attachments(ctx context.Context) {

	fmt.Println("----------------------\n")
	dirOfAttachedFile, err := os.Getwd()
	if err != nil {
		fmt.Println("FATAL: Error Getting Directory of current executable: ", err)
	}
	fmt.Println("Directory where attachment file is put for testing: ", dirOfAttachedFile)
	fmt.Println("----------------------\n")

	contentBytes, err := ioutil.ReadFile(dirOfAttachedFile + "/test2.txt")
	if err != nil {
		panic(err)
		return
	}
	attachment1 := Attachment{
		Filename:     "test2.txt",
		ContentType:  ATTACHMENT_CONTENT_TYPE_PLAIN,
		ContentBytes: contentBytes,
	}
	attachments := []Attachment{attachment1}

	sendEmailRequestInput := SendEmailRequestVO{

		Attachments: attachments,

		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "deb.work.related@gmail.com",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			BodyHtml: Content{
				CharSet: "iso-8859-1",
				Data:    "<!DOCTYPE html><html><head><img src='https://storage.googleapis.com/ustat-ww/ws/mindlounge.co.in/img/ID/ml.png'><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>New lead has been created.</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Prospect Name</td><td>India</td><td>Banani</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>",
			},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingRawMessege_HtmlBody_With_Attachments",
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}
}

/*
*
*
Message-ID: <3355164019@1644817933493076000>
From: 'deb.work.related@gmail.com' <deb.work.related@gmail.com>
Subject: TestSendingRawMessege_TextBody_With_Attachments - SUBJECT
To: debdas.sinha@gmail.com,debdas.sinha@outlook.com
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="3330259215"

--3330259215
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

Text Body of TestSendingRawMessege_TextBody_With_Attachments
--3330259215
Content-Type: text/plain; name="test2.txt"
Content-Disposition: attachment;filename="test2.txt"
Content-Transfer-Encoding: base64
X-Attachment-Id: 687242461

RGV2T3BzIDIKCgpQdWJsaWMgUmVwb3MgMg==

--3330259215--
*
*/
func TestSendingRawMessege_TextBody_With_Attachments(ctx context.Context) {

	fmt.Println("----------------------\n")
	dirOfAttachedFile, err := os.Getwd()
	if err != nil {
		fmt.Println("FATAL: Error Getting Directory of current executable: ", err)
	}
	fmt.Println("Directory where attachment file is put for testing: ", dirOfAttachedFile)
	fmt.Println("----------------------\n")

	contentBytes, err := ioutil.ReadFile(dirOfAttachedFile + "/test2.txt")
	if err != nil {
		panic(err)
		return
	}
	attachment1 := Attachment{
		Filename:     "test2.txt",
		ContentType:  ATTACHMENT_CONTENT_TYPE_PLAIN,
		ContentBytes: contentBytes,
	}
	attachments := []Attachment{attachment1}

	sendEmailRequestInput := SendEmailRequestVO{

		Attachments: attachments,

		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "deb.work.related@gmail.com",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Message: EmailBody{
			BodyText: Content{
				Data: "Text Body of TestSendingRawMessege_TextBody_With_Attachments",
			},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingRawMessege_TextBody_With_Attachments - SUBJECT",
		},
	}
	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)
	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}
}

/*
*
*
Message-ID: <2495239802@1644817308314355000>
From: 'itsupport@continental.coffee' <itsupport@continental.coffee>
Subject: Lead created
To: debdas.sinha@outlook.com,debdas.sinha@gmail.com
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="3330259215"

--3330259215
Content-Type: multipart/related; boundary="rel_3330259215"

--rel_3330259215
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

<!DOCTYPE html><html><head><img src=3D=22cid:3553515852@1644817308314418000=22><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>TestSendingRawMessege_From_Template_And_With_Attachments</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>Mixed</td><td>Related</td><td>Raw</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>
--rel_3330259215
Content-Type: image/png; name="83517792.png"
Content-Transfer-Encoding: base64
Content-Disposition: inline; filename="83517792.png"
Content-ID: <3553515852@1644817308314418000>
X-Attachment-Id: 3553515852@1644817308314418000

iVBORw0KGgoAAAANSUhEUgAAAMgAAABzCAMAAADqi70aAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOTTPK/5S8DSU5xCEYx8gAAAABJRU5ErkJggg==

--rel_3330259215--
--3330259215
Content-Type: text/plain; name="test2.txt"
Content-Disposition: attachment;filename="test2.txt"
Content-Transfer-Encoding: base64
X-Attachment-Id: 2980185960

RGV2T3BzIDIKCgpQdWJsaWMgUmVwb3MgMg==

--3330259215--

*
*/
func TestSendingRawMessege_From_Template_And_With_Attachments(ctx context.Context) {

	fmt.Println("----------------------\n")
	dirOfAttachedFile, err := os.Getwd()
	if err != nil {
		fmt.Println("FATAL: Error Getting Directory of current executable: ", err)
	}
	fmt.Println("Directory where attachment file is put for testing: ", dirOfAttachedFile)
	fmt.Println("----------------------\n")

	//contentBytes, err := ioutil.ReadFile(dirOfAttachedFile + "/TestAttachDoc.pdf")
	contentBytes, err := ioutil.ReadFile(dirOfAttachedFile + "/test2.txt")
	if err != nil {
		panic(err)
		return
	}
	attachment2 := Attachment{
		Filename:     "test2.txt",
		ContentType:  ATTACHMENT_CONTENT_TYPE_PLAIN,
		ContentBytes: contentBytes,
	}
	/*contentBytes, err = ioutil.ReadFile(dirOfAttachedFile + "/test1.txt")
	if err != nil {
		panic(err)
		return
	}
	attachment1 := Attachment{
		Filename:     "test1.txt",
		ContentType:  ATTACHMENT_CONTENT_TYPE_PLAIN,
		ContentBytes: contentBytes,
	}*/
	attachments := []Attachment{attachment2}

	dataFeedToTemplate := make(map[string]string)

	dataFeedToTemplate["MessageToUser"] = "TestSendingRawMessege_From_Template_And_With_Attachments"
	dataFeedToTemplate["AccountName"] = "Mixed"
	dataFeedToTemplate["AccountCountry"] = "Related"
	dataFeedToTemplate["AccountOwner"] = "Raw"

	mediaDataToEmplate := make(map[string]Content)

	mediaContentImage := Content{
		ContentType: CONTENT_TYPE_IMAGE_PNG,
		//Url:         "https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png",
		Url: "https://storage.googleapis.com/ustat-ww/ws/mindlounge.co.in/img/ID/ml.png",
	}
	mediaDataToEmplate["TOP_WIDE_N_SHORT_BANNER_CCL"] = mediaContentImage

	sendEmailRequestInput := SendEmailRequestVO{

		Attachments: attachments,
		Message: EmailBody{
			TemplateForHtmlBody: EmailTemplate{
				TemplateRef:           "OnNewLead-DynamicInlineImageURLs-EscapeChars",
				DynamicDataOfTemplate: dataFeedToTemplate,
			},
			DynamicDataOfMediaInBody: mediaDataToEmplate,
		},
		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "itsupport@continental.coffee",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@outlook.com", "debdas.sinha@gmail.com"},
		},
	}
	fmt.Println("[[Email Sender]] Tester TestSendingRawMessege_From_Template_And_With_Attachments -> Num Attach -> ", len(sendEmailRequestInput.Attachments))

	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)

	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}

}

/*
*
*
Message-ID: <1565591125@1644816638015527000>
From: 'itsupport@continental.coffee' <itsupport@continental.coffee>
Subject: Lead created
To: debdas.sinha@gmail.com,debdas.sinha@outlook.com
MIME-Version: 1.0
Content-Type: multipart/related; boundary="1198242597"

--1198242597
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

<!DOCTYPE html><html><head><img src=3D=22cid:1937804623@1644816638015546000=22><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;} td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body><h3>Hi,</h3><p>TestSendingRawMessage_Templated_Without_Attachments</p><table><tr><th>Account Name</th><th>Account Country</th><th>Account Owner</th></tr><tr><td>RawMessage</td><td>Dynamic</td><td>ImageUrlHandling</td></tr></table><p>Regards,</p><p>a2z.cclproducts</p></body></html>
--1198242597
Content-Type: image/png; name="3930379605.png"
Content-Transfer-Encoding: base64
Content-Disposition: inline; filename="3930379605.png"
Content-ID: <1937804623@1644816638015546000>
X-Attachment-Id: 1937804623@1644816638015546000

iVBORw0KGgoAAAANSUhEUgAAAMgAAABzCAMAAADqi70aAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAR6aVRYdFhNTDpjb20uYWRvYmUueG1wAAAAAAA8eDp4bXBtZXRhIHhtbG5zOng9ImFkb2JlOm5zOm1ldGEvIiB4OnhtcHRrP
tdOMskkk0wyySSTTPK/5S8DSU5xCEYx8gAAAABJRU5ErkJggg==

--1198242597--

*
*
*/
func TestSendingRawMessage_Templated_Without_Attachments(ctx context.Context) {

	dataFeedToTemplate := make(map[string]string)
	dataFeedToTemplate["MessageToUser"] = "TestSendingRawMessage_Templated_Without_Attachments"
	dataFeedToTemplate["AccountName"] = "RawMessage"
	dataFeedToTemplate["AccountCountry"] = "Dynamic"
	dataFeedToTemplate["AccountOwner"] = "ImageUrlHandling"

	mediaDataToEmplate := make(map[string]Content)

	mediaContentImage := Content{
		ContentType: CONTENT_TYPE_IMAGE_PNG,
		Url:         "https://storage.googleapis.com/ustat-ww/ws/mindlounge.co.in/img/ID/ml.png",
		//"https://s3.ap-south-1.amazonaws.com/beta-a2z.cclproducts.com/static/media/CCLEmailTemplate.png",
	}
	mediaDataToEmplate["TOP_WIDE_N_SHORT_BANNER_CCL"] = mediaContentImage

	sendEmailRequestInput := SendEmailRequestVO{

		Message: EmailBody{
			TemplateForHtmlBody: EmailTemplate{
				TemplateRef:           "OnNewLead-DynamicInlineImageURLs-EscapeChars",
				DynamicDataOfTemplate: dataFeedToTemplate,
			},
			DynamicDataOfMediaInBody: mediaDataToEmplate,
		},
		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "itsupport@continental.coffee",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
	}

	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)

	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success -> ", emailRequestProcessingSuccess)
	}

}

/*
Tested with Both RawMessage and ses api based simpleMessage
*
* Should build raw message like below. ** But it's happening through SES API => sendSimpleMessage
*
*
Message-ID: <1499091322@1644814485809859000>
From: 'deb.work.related@gmail.com' <deb.work.related@gmail.com>
Subject: TestSendingRawMessege_EmptyBody_Without_Attachments 2 - SUBJECT
To: debdas.sinha@gmail.com,debdas.sinha@outlook.com
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Content-Length: 0

*/
func TestSendingRawMessege_EmptyBody_Without_Attachments(ctx context.Context) {

	sendEmailRequestInput := SendEmailRequestVO{
		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "deb.work.related@gmail.com",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com", "debdas.sinha@outlook.com"},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingRawMessege_EmptyBody_Without_Attachments 2 - SUBJECT",
		},
	}

	fmt.Println("[[Email Sender]] Tester TestSendingRawMessege_EmptyBody_Without_Attachments -> Num Attach -> ", len(sendEmailRequestInput.Attachments))

	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)

	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}

}

/*
*
*
*
Message-ID: <2526350408@1644813697725928000>
From: 'deb.work.related@gmail.com' <deb.work.related@gmail.com>
Subject: TestSendingRawMessege_EmptyBody_With_Attachments - SUBJECT
To: debdas.sinha@gmail.com
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="3330259215"

--3330259215
Content-Type: text/plain; charset=UTF-8
Content-Length: 0

--3330259215
Content-Type: text/plain; name="test.txt"
Content-Disposition: attachment;filename="test.txt"
Content-Transfer-Encoding: base64
X-Attachment-Id: 3163740384

RGV2T3BzIDIKCgpQdWJsaWMgUmVwb3MgMg==

--3330259215--
*
*/
func TestSendingRawMessege_EmptyBody_With_Attachments(ctx context.Context) {

	fmt.Println("----------------------\n")
	dirOfAttachedFile, err := os.Getwd()
	if err != nil {
		fmt.Println("FATAL: Error Getting Directory of current executable: ", err)
	}
	fmt.Println("Directory where attachment file is put for testing: ", dirOfAttachedFile)
	fmt.Println("----------------------\n")

	//contentBytes, err := ioutil.ReadFile(dirOfAttachedFile + "/TestAttachment.pdf")
	contentBytes, err := ioutil.ReadFile(dirOfAttachedFile + "/test2.txt")
	if err != nil {
		panic(err)
		return
	}
	attachment1 := Attachment{
		Filename:    "test2.txt",
		ContentType: ATTACHMENT_CONTENT_TYPE_PLAIN,
		//ContentType:  ATTACHMENT_CONTENT_TYPE_PDF,
		ContentBytes: contentBytes,
	}
	attachments := []Attachment{attachment1}

	sendEmailRequestInput := SendEmailRequestVO{

		Attachments: attachments,

		Type: EmailType{
			CallerITInfraRegion: "us-east-1",
			IsSystemEmail:       true,
		},
		SenderDetails: Sender{
			SendFromIdentity: "deb.work.related@gmail.com",
		},
		TargetRecipients: Recipients{
			ToList: []string{"debdas.sinha@gmail.com"},
		},
		Subject: Content{
			CharSet: "iso-8859-1",
			Data:    "TestSendingRawMessege_EmptyBody_With_Attachments - SUBJECT",
		},
	}

	fmt.Println("[[Email Sender]] Tester TestSendingRawMessege_EmptyBody_With_Attachments -> Num Attach -> ", len(sendEmailRequestInput.Attachments))

	emailRequestProcessingSuccess, err := Send(ctx, sendEmailRequestInput)

	if err != nil {
		fmt.Println("[[Email Sender]] Processing Error -> ", err.Error())
	} else {
		fmt.Println("[[Email Sender]] Processing Success ->  ", emailRequestProcessingSuccess)
	}

}
