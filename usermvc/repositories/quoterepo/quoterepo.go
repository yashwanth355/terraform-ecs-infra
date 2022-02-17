package quoterepo

import (
	"context"
	"strings"
	"usermvc/entity"
	"usermvc/model"
	"usermvc/repositories"
	logger2 "usermvc/utility/logger"

	"github.com/jinzhu/gorm"
)

type QuoteRepo interface {
	GetQuoatotionCreateInfoReq(ctx context.Context, req model.GetQuoatotionCreateInfoReq) (interface{}, error)
	GetProdPackcategoryName(ctx context.Context, categoryId int) (*string, error)
	GetProdPackCategoryWeight(ctx context.Context, weightId int) (*string, error)
}

type quoteRepo struct {
	db *gorm.DB
}

/*
*
 */
func NewLeadRepo() QuoteRepo {
	newDb, err := repositories.NewDb()
	if err != nil {
		panic(err)
	}
	return &quoteRepo{
		db: newDb,
	}
}

/*
*
 */
func (quoteRepoRef quoteRepo) GetQuoatotionCreateInfoReq(ctx context.Context, req model.GetQuoatotionCreateInfoReq) (interface{}, error) {
	const (
		INCOTERMS        = "incoterms"
		CURRENCIES       = "currencies"
		LOADINGPORTS     = "loadingports"
		DESTINATIONPORTS = "destinationports"
		ACCOUNTDETAILS   = "accountdetails"
		VIEWQUOTE        = "viewquote"
		REQUESTPRICE     = "requestprice"
	)
	logger := logger2.GetLoggerWithContext(ctx)
	if req.Type == INCOTERMS {
		logger.Info("going to fetch all records from cms_incoterms_master")
		var response entity.CmsIncotermsMaster
		err := quoteRepoRef.db.Table("cms_incoterms_master").Model(&entity.CmsIncotermsMaster{}).Find(&response).Error
		if err != nil {
			logger.Error("error while get all account details from cms_leads_master ", err.Error())
			return nil, err
		}
		return response, err
	}
	if req.Type == CURRENCIES {
		sqlStatement := `SELECT "currencyid", "currencyname", "currencycode" FROM "project_currency_master"`
		rows, err := quoteRepoRef.db.Table("project_currency_master").Raw(sqlStatement).Rows()
		if err != nil {
			logger.Error("error while get all account details from project_currency_master ", err.Error())
			return nil, err
		}
		var allCurrencies []model.Currencies
		defer rows.Close()
		for rows.Next() {
			var currency model.Currencies
			err = rows.Scan(&currency.Currencyid, &currency.Currencyname, &currency.Currencycode)
			allCurrencies = append(allCurrencies, currency)
		}
		return allCurrencies, nil
	}
	if req.Type == LOADINGPORTS {
		logger.Info("going to fetch all records from cms_portloading_master")
		var response model.Loadingports
		err := quoteRepoRef.db.Table("cms_portloading_master").Model(&entity.CmsPortloadingMaster{}).Find(&response).Error
		if err != nil {
			logger.Error("error while get all account details from cms_portloading_master ", err.Error())
			return nil, err
		}
		return response, err
	}
	if req.Type == DESTINATIONPORTS {
		logger.Info("going to fetch all records from cms_destination_master")
		var response model.Destinationports
		err := quoteRepoRef.db.Table("cms_destination_master").Model(&entity.CmsDestinationMaster{}).Find(&response).Error
		if err != nil {
			logger.Error("error while get all account details from cms_destination_master ", err.Error())
			return nil, err
		}
		return response, err
	}
	if req.Type == ACCOUNTDETAILS {
		var allAccounts []model.AccountDetails
		sqlStatement := `select a.accountid, a.accountname, a.accounttypeid, concat(b.street,' ', b.city,' ', b.stateprovince, ' ', b.postalcode) as address
		from 
	   accounts_master a
		INNER JOIN accounts_billing_address_master b on b.accountid =a.accountid
		INNER Join cms_leads_master l on l.accountid = a.accountid`

		rows, err := quoteRepoRef.db.Raw(sqlStatement).Rows()

		defer rows.Close()
		for rows.Next() {
			var account model.AccountDetails
			err = rows.Scan(&account.Accountid, &account.Accountname, &account.Accounttypeid, &account.Address)

			var accounttypes []string
			if account.Accounttypeid != "" {
				z := strings.Split(account.Accounttypeid, ",")
				for i, z := range z {
					logger.Info("get account name", i, z)
					sqlStatement := `SELECT accounttype FROM "cms_account_type_master" where accounttypeid=$1`
					rows1, err1 := quoteRepoRef.db.Raw(sqlStatement, z).Rows()

					if err1 != nil {
						logger.Error(err, "unable to add account names")
					}

					for rows1.Next() {
						var accounttype string
						err = rows1.Scan(&accounttype)
						accounttypes = append(accounttypes, accounttype)
					}
				}
			}
			account.Accounttypename = strings.Join(accounttypes, ",")
			if account.Accountid != 0 {
				sqlStatement := `SELECT contactid, contactfirstname FROM "contacts_master" where accountid=$1`
				rows2, err2 := quoteRepoRef.db.Raw(sqlStatement, account.Accountid).Rows()

				if err2 != nil {
					logger.Error(err)
					logger.Info("unable to add account type", account.Accountid)
				}

				for rows2.Next() {
					var contact model.ContactDetails
					err = rows2.Scan(&contact.Contactid, &contact.Contactname)
					allAccounts := append(account.Contacts, contact)
					account.Contacts = allAccounts
					logger.Info("added one account", allAccounts)
				}
			}
			allAccounts = append(allAccounts, account)
		}
		return allAccounts, nil
	}
	if req.Type == VIEWQUOTE {
		logger.Info("get account details", req.Type)
		sqlStatement := `SELECT q.accountid,a.accountname, q.accounttypename, q.contactid, t.contactfirstname, q.createddate, u.firstname as createdby , 
        r.currencyname,
        q.currencycode,
        q.currencyid,
        q.fromdate,
		q.todate,
		q.paymentterm,
		q.otherspecification,
		q.remarks,
		q.destinationcountryid,
		q.destination,
		q.finalaccountid,
		concat(b.street,' ', b.city,' ', b.stateprovince, ' ', b.postalcode) as address,
		c.incoterms,
		q.incotermsid,
		s.status,
		q.portloading,
        q.portloadingid,
		q.destinationid,
		q.remarksfromgmc from crm_quote_master q
        INNER JOIN accounts_master a on q.accountid = a.accountid 
		INNER JOIN accounts_billing_address_master b on q.accountid = b.accountid
		INNER JOIN cms_incoterms_master c on q.incotermsid = c.incotermsid
        INNER JOIN cms_allstatus_master s ON q.statusid = s.id
		INNER JOIN userdetails_master u ON q.createdby = u.userid
        INNER JOIN contacts_master t ON q.contactid = t.contactid
		INNER JOIN project_currency_master r ON q.currencyid = r.currencyid
		where q.quoteid =$1`
		rows, err := quoteRepoRef.db.Raw(sqlStatement, req.QuoteId).Rows()
		if err != nil {
			logger.Error("error while getting data ", err.Error())
			return nil, err
		}
		var quote model.GetQuoteDetails
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&quote.Accountid,
				&quote.Accountname,
				&quote.Accounttypename,
				&quote.Contactid,
				&quote.Contactname,
				&quote.CreatedDate,
				&quote.Createdby,
				&quote.Currencyname,
				&quote.Currencycode,
				&quote.Currencyid,
				&quote.Fromdate,
				&quote.Todate,
				&quote.Paymentterms,
				&quote.Otherspecifications,
				&quote.Remarksfrommarketing,
				&quote.Destinationcountryid,
				&quote.Destination,
				&quote.Finalclientaccountid,
				&quote.Billingaddress,
				&quote.Incoterms,
				&quote.Incotermsid,
				&quote.Status,
				&quote.PortLoading,
				&quote.Portloadingid,
				&quote.Portdestinationid,
				&quote.Remarksfromgmc)

		}

		if quote.Finalclientaccountid != "" {

			sqlStatement := `select accountname as address from accounts_master b where accountid=$1`

			rows, err := quoteRepoRef.db.Raw(sqlStatement, quote.Finalclientaccountid).Rows()

			if err != nil {
				logger.Info(err, "unable to add final account name", quote.Finalclientaccountid)
				return nil, err
			}

			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&quote.Finalclientaccountname)
				logger.Info("added final account name")
			}
		}
		return quote, nil
	}

	if req.Type == REQUESTPRICE {
		logger.Info("update request price status", req.Type)
		sqlStatement := `UPDATE crm_quote_master SET statusid=2 where quoteid=$1`
		rows, err := quoteRepoRef.db.Raw(sqlStatement, req.QuoteId).Rows()
		if err != nil {
			logger.Info("error while getting data ", err.Error())
			return err, nil
		}
		return rows, nil
	}
	return nil, nil
}

func (quoteRepoRef quoteRepo) GetProdPackcategoryName(ctx context.Context, categoryId int) (*string, error) {
	type categoryName struct {
		Categorytypename string
	}
	var result categoryName
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("going to fetch all leadAccounts")
	sqlStatement := `select categorytypename from cms_prod_pack_category_type where categorytypeid=$1`
	if err := quoteRepoRef.db.Raw(sqlStatement, categoryId).Scan(&result).Error; err != nil {
		logger.Error("error while get all account details from Accont master ", err.Error())
		return nil, err
	}

	return &result.Categorytypename, nil
}

func (quoteRepoRef quoteRepo) GetProdPackCategoryWeight(ctx context.Context, weightId int) (*string, error) {
	type Weight struct {
		Weightname string
	}
	var result Weight
	logger := logger2.GetLoggerWithContext(ctx)
	logger.Info("going to fetch all leadAccounts")
	sqlStatement := `select weightname from cms_prod_pack_category_weight where weightid=$1`
	if err := quoteRepoRef.db.Raw(sqlStatement, weightId).Scan(&result).Error; err != nil {
		logger.Error("error while get all account details from Accont master ", err.Error())
		return nil, err
	}
	return &result.Weightname, nil
}
