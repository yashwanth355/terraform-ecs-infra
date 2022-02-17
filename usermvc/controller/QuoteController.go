package controller

// import (
// 	"usermvc/model"
// 	"usermvc/repositories/accountrepo"
// 	"usermvc/repositories/quoterepo"
// 	logger2 "usermvc/utility/logger"

// 	"github.com/gin-gonic/gin"
// )

// type (
// 	QuoteController interface {
// 		GetAllQuoteLineItems(ctx *gin.Context)
// 		GetAllQuotes(ctx *gin.Context)
// 		GetQuotationCreationInfo(ctx *gin.Context)
// 		GetQuoteInformation(ctx *gin.Context)
// 	}

// 	quoteController struct {
// 		accountrepo accountrepo.AccountRepo
// 		quoteRepo   quoterepo.QuoteRepo
// 	}
// )

// /*
// *
//  */
// func NewQuoteController() QuoteController {
// 	return &quoteController{
// 		accountrepo: accountrepo.NewAccountRepo(),
// 		quoteRepo:   quoterepo.NewLeadRepo(),
// 	}
// }

// func (qControllerRef quoteController) GetAllQuoteLineItems(ctx *gin.Context) {
// 	logger := logger2.GetLoggerWithContext(ctx)
// 	var getAllQoutedRequestBody *model.GetAllQuoteLineRequest
// 	logger.Info("start hadling leadController.leadController")
// 	if err := ctx.ShouldBindJSON(&getAllQoutedRequestBody); err != nil {
// 		logger.Error("error while parsing the qoutationItems", err.Error())
// 		ctx.JSON(403, err.Error())
// 		return
// 	}
// 	res, err := qControllerRef.accountrepo.GetAllQuoteLineItems(ctx, getAllQoutedRequestBody)
// 	if err != nil {
// 		logger.Error("error while getting all QuoteLineItens", err.Error())
// 		ctx.JSON(503, err.Error())
// 		return
// 	}
// 	for k, QuoteLineItem := range res {
// 		if QuoteLineItem.Packcategorytypeid != nil {
// 			categoryName, err := qControllerRef.quoteRepo.GetProdPackcategoryName(ctx, *QuoteLineItem.Packcategorytypeid)
// 			if err != nil {
// 				logger.Error("error while getting category name where category_id is ", QuoteLineItem.Packcategorytypeid)
// 				ctx.JSON(503, err.Error())
// 				return
// 			}
// 			res[k].CategoryType = *categoryName
// 		}
// 		if QuoteLineItem.Packweightid != nil {
// 			weight, err := qControllerRef.quoteRepo.GetProdPackCategoryWeight(ctx, *QuoteLineItem.Packweightid)
// 			if err != nil {
// 				logger.Error("error while getting  weight ", QuoteLineItem.Packweightid)
// 				ctx.JSON(503, err.Error())
// 				return
// 			}
// 			res[k].Weight = *weight
// 		}
// 	}
// 	ctx.JSON(200, res)
// }

// func (qControllerRef quoteController) GetAllQuotes(ctx *gin.Context) {
// 	logger := logger2.GetLoggerWithContext(ctx)
// 	validate, err := qControllerRef.validateGetAllQuotes(ctx)
// 	if err != nil {
// 		logger.Error("error while validating request", err.Error())
// 		ctx.JSON(503, "error while validating request")
// 		return
// 	}
// 	if !validate {
// 		logger.Error("request couls not be validate")
// 		ctx.JSON(404, "invalid request")
// 		return
// 	}
// 	var getAllQoutedRequestBody *model.GetAllQoutesRequestBody
// 	if err := ctx.ShouldBindJSON(&getAllQoutedRequestBody); err != nil {
// 		logger.Error("error while parsing the GetAllQuoteRequest", err.Error())
// 		ctx.JSON(403, err.Error())
// 		return
// 	}
// 	res, err := qControllerRef.accountrepo.GetAllQoutes(ctx, getAllQoutedRequestBody)
// 	if err != nil {
// 		logger.Error("error while getting allQuoteItems", err.Error())
// 		ctx.JSON(503, err.Error())
// 		return
// 	}
// 	ctx.JSON(200, res)
// }

// func (qControllerRef quoteController) GetQuotationCreationInfo(ctx *gin.Context) {
// 	logger := logger2.GetLoggerWithContext(ctx)
// 	var getQuoatotionCreateInfoReq model.GetQuoatotionCreateInfoReq
// 	if err := ctx.ShouldBindJSON(&getQuoatotionCreateInfoReq); err != nil {
// 		logger.Error("error while parsing the getQuoatotionCreateInfoReq", err.Error())
// 		ctx.JSON(403, err.Error())
// 		return
// 	}

// 	LeadsInfo, err := qControllerRef.quoteRepo.GetQuoatotionCreateInfoReq(ctx, getQuoatotionCreateInfoReq)
// 	if err != nil {
// 		logger.Error("error while getting GetCmsLeadsBillingAddressMaster", err.Error())
// 		ctx.JSON(503, err.Error())
// 		return
// 	}

// 	ctx.JSON(200, LeadsInfo)
// }

// func (qControllerRef quoteController) GetQuoteInformation(ctx *gin.Context) {

// }

// func (qControllerRef quoteController) validateGetAllQuotes(ctx *gin.Context) (bool, error) {
// 	return true, nil
// }

// //func (lc leadController) 	GetQuotationCreationInfo(ctx *gin.Context) {
// ////	1) looger method
// //	logger := logger2.GetLoggerWithContext(ctx)
// //	validate, err := lc.validateGetLeadCreationInfo(ctx)
// //	if err != nil {
// //		logger.Error("error while validating request", err.Error())
// //		ctx.JSON(503, "error while validating request")
// //		return
// //	}
// //	if !validate {
// //		logger.Error("request could not be validate")
// //		ctx.JSON(404, "invalid request")
// //		return
// //	}
// //	var getAllQoutedRequestBody *model.GetQuoatotionReqBody
// //	if err := ctx.ShouldBindJSON(&getAllQoutedRequestBody); err != nil {
// //		logger.Error("error while parsing the GetAllQuoteRequest", err.Error())
// //		ctx.JSON(403, err.Error())
// //		return
// //	}
// //	res, err := lc.accountrepo.GetQuotationCreationInfo(ctx, getAllQoutedRequestBody)
// //	if err != nil {
// //		logger.Error("error while getting GetLeadCreationInfo", err.Error())
// //		ctx.JSON(503, err.Error())
// //		return
// //	}
// //	ctx.JSON(200, &model.GetQuoationResp{
// //		Status:  200,
// //		Payload: res,
// //	})
// ////}
