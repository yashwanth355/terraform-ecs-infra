package controller

//update
import (
	"usermvc/model"
	"usermvc/repositories/accountrepo"
	logger2 "usermvc/utility/logger"

	"github.com/gin-gonic/gin"
)
type AccountController interface {
	CvuAccountContactDetails(ctx *gin.Context)
}
type accountController struct {
	accountRepo accountrepo.AccountRepo
}

func NewAccountController() AccountController {
	return accountController{
		accountRepo: accountrepo.NewAccountRepo(),
	}
}
func (account accountController) CvuAccountContactDetails(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	var cvuaccountcontactdetailsrequest *model.Input
	logger.Info("started handling poController.cvuAccountContactDetails")
	if err := ctx.ShouldBindJSON(&cvuaccountcontactdetailsrequest); err != nil {
		logger.Error("error while parsing the cvuaccount contact details request body", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := account.accountRepo.CvuAccountContactDetails(ctx, cvuaccountcontactdetailsrequest)
	if err != nil {
		ctx.JSON(503, "error while getting cvuAccount Contact Details from repo")
		return
	}
	ctx.JSON(200, res)
	return
}

// func (ctr controller) GetAllAccountDetails(c *gin.Context) {
// 	res, err := ctr.AccountRepo.GetAllAccountDetails(context.Background())

// 	if err != nil {
// 		zap.S().Error("error from the getappAccountDetails ", err.Error())
// 		c.JSON(500, err.Error())
// 		return
// 	}

// 	c.JSON(200, res)
// }

// //
// //func (ctr controller) GetAllLeadAccounts(c *gin.Context)   {
// //	res, err := ctr.accountSvc.GetAllAccountDetails(context.Background())
// //	if err != nil {
// //		zap.S().Error("error from the getappAccountDetails ", err.Error())
// //		c.JSON(500, err.Error())
// //		return
// //	}
// //	c.JSON(200, &model.GetAllAccountDetailsResponse{
// //		Status:  200,
// //		Payload: res,
// //	})
// //}
// func validateRequest(ctx gin.Context) error {
// 	//write validation here
// 	return nil
// }

// func (ctr controller) TestPing(c *gin.Context) {
// 	c.Status(200)
// }