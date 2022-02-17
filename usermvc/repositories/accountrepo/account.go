package accountrepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	// "database/sql"
	// "reflect"
	"usermvc/entity"
	"usermvc/model"
	"usermvc/repositories"
	logger2 "usermvc/utility/logger"
)

type AccountRepo interface {
	//Create(context context.Context, accountDetailsRequest entity.AccountDetails) ()
	CvuAccountContactDetails(ctx context.Context, req *model.Input) (interface{}, error)
	// GetAllAccountDetails(ctx context.Context) ([]*model.GetAllAccountsResponseBody, error)
	// GetAllQuoteLineItems(ctx context.Context, req *model.GetAllQuoteLineRequest) ([]model.GetAllQouteLineItemsResponseBody, error)
	// GetAllQoutes(ctx context.Context, req *model.GetAllQoutesRequestBody) ([]*model.QuoteDetails, error)
	//GetLeadCreationInfo(ctx context.Context, req *model.GetLeadCreationInfoRequest) (interface{}, error)
	//GetQuotationCreationInfo(ctx context.Context, req *model.GetQuoatotionCreateInfoReq) (interface{}, error)

}

const (
	MarketingMExecutive = "Marketing Executive"
	ManagingDirector    = "Managing Director"
	PENDINGQUOTES       = "pendingquotes"
	MYQOUTES            = "myquotes"
)

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo() AccountRepo {
	newDb, err := repositories.NewDb()
	if err != nil {
		panic(err)
	}
	newDb.AutoMigrate(&entity.User{})
	return &accountRepo{
		db: newDb,
	}
}

// type AccountRepo interface {
// 	CvuAccountContactDetails(ctx context.Context, req *model.Input) (interface{}, error)
// }
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
func (account accountRepo) CvuAccountContactDetails(ctx context.Context, req *model.Input) (interface{}, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	var con *model.Contactdetails
	var bi model.BillingInfo
	var filter string
	if con.Update {
		sqlStatementuc1 := `UPDATE dbo.contacts_master SET 
			salutationid=$1,
			contactfirstname=$2,
			contactlastname=$3,
			contactemail=$4,
			position=$5,
			contactphone=$6,	
			contactmobilenumber=$7,
			contactmodifiedby=$8
			where contactid=$9`
		rows, err := account.db.Raw(sqlStatementuc1,
			con.Salutationid,
			NewNullString(con.Firstname),
			NewNullString(con.Lastname),
			NewNullString(con.Email),
			NewNullString(con.Position),
			NewNullString(con.Phone),
			NewNullString(con.Mobile),
			NewNullString(con.Loggedinuserid),
			con.ContactID).Rows()

		if err != nil {
			logger.Error("error while fetching records from dbo.contacts_master SET ", err.Error())
			logger.Error(err.Error())
		}
		defer rows.Close()
	} else if con.Create {
		log.Println("Creating Contact for the account")
		log.Println("Finding latest contactid")
		var Pkey string
		// Pkey := strconv.Itoa(FindLatestSerial("contactid", "dbo.contacts_master", "contactid", "contactid"))
		sqlStatementc1 := `INSERT INTO dbo.contacts_master (
			contactid,
			salutationid,
			contactfirstname,
			contactlastname,
			contactemail,
			position,
			contactphone,	
			contactmobilenumber,
			accountid,
			contactcreatedby)
			VALUES($1, $2, $3, $4, $5, $6, $7,$8,$9,$10)`
		// billingaddressid,
		// shippingaddressid
		// (select billingid from dbo.accounts_billing_address_master where accountid=$10),
		// (select shippingid from dbo.accounts_shipping_address_master where accountid=$11)
		rows, err := account.db.Raw(sqlStatementc1,
			Pkey,
			con.Salutationid,
			con.Firstname,
			con.Lastname,
			con.Email,
			con.Position,
			con.Phone,
			con.Mobile,
			con.AccountID,
			con.Loggedinuserid).Rows()
		if err != nil {
			log.Println("Insert to Contacts table failed")
			log.Println(err.Error())
		}
		defer rows.Close()
	} else if con.View {
		var accname, accowner, cfname, clname, cemail, cphone, cmobile, csalesnum, cpastsalesnum sql.NullString
		log.Println("Account Contacts View Module entered")
		sqlStatementv1 := `select 
								 c.contactid,
								 acc.accountid,
								 acc.accountname,
								 acc.account_owner,
								 c.contactfirstname,
								 c.contactlastname,
								 c.contactemail,
								 c.contactphone,
								 c.contactmobilenumber,
								 c.currentsalesnumber,
								 c.pastsalesnumber
								 from dbo.contacts_master c
								 inner join
								 dbo.accounts_master acc on c.accountid=acc.accountid
								 where c.contactid=$1`
		rows, err := account.db.Raw(sqlStatementv1, con.ContactID).Rows()
		if err != nil {
			log.Println("Fetching Contact from DB failed")
			log.Println(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&con.ContactID, &con.AccountID, &accname, &accowner, &cfname, &clname, &cemail, &cphone, &cmobile, &csalesnum, &cpastsalesnum)
		}
		con.Accountname = accname.String
		con.ContactOwner = accowner.String
		con.Firstname = cfname.String
		con.Lastname = clname.String
		con.Email = cemail.String
		con.Phone = cphone.String
		con.Mobile = cmobile.String
		con.CurrentSalesNo = csalesnum.String
		con.PastSalesNo = cpastsalesnum.String
		//Billing Info
		if con.BillingAdvancedFilter {
			log.Println("Advanced Filter Selected for billing address")
			filter = "where accountid=" + con.AccountID + " and" + con.FilterParam
		} else {
			filter = ""
		}
		log.Println("get Billing Info")
		sqlStatementbi1 := `select * from (SELECT billingid as billing_id,
							street as billing_street,city as billing_city,
							stateprovince as billing_state,postalcode as billing_postalcode,
							country as billing_country,primary_address as billing_primary
							from dbo.accounts_billing_address_master where accountid=%s and contactid=%s) a %s`
		rows, err = account.db.Raw(fmt.Sprintf(sqlStatementbi1, con.AccountID, con.ContactID, filter)).Rows()
		if err != nil {
			log.Println(err)
			log.Println("unable to find billing address for the account")
		}
		var bstreet, bcity, bstate, bpcode, bcountry sql.NullString
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&bi.B_BillingID, &bstreet, &bcity, &bstate, &bpcode, &bcountry, &bi.B_Primary)
			bi.B_Street = bstreet.String
			bi.B_City = bcity.String
			bi.B_State = bstate.String
			bi.B_PostalCode = bpcode.String
			bi.B_Country = bcountry.String
			billInfo := append(con.BillingInfo, bi)
			con.BillingInfo = billInfo
			log.Println("added billing info")
			log.Println(con.BillingInfo)
		}
		//Get Shipping Info
		if con.ShippingAdvancedFilter {
			log.Println("Advanced Filter Selected for shipping address")
			filter = "where " + con.FilterParam
		} else {
			filter = ""
		}
		log.Println("get shipping Info")
		sqlStatementsi1 := `select * from (SELECT shippingid as shipping_id,
							street as shipping_street,city as shipping_city,
							stateprovince as shipping_state,postalcode as shipping_postalcode,
							country as shipping_country,primary_address as shipping_primary
							from dbo.accounts_shipping_address_master where accountid=%s and contactid=%s) a %s`
		rows, err = account.db.Raw(fmt.Sprintf(sqlStatementsi1, con.AccountID, con.ContactID, filter)).Rows()

		if err != nil {
			log.Println(err)
			log.Println("unable to find shipping address for the account")
		}
		defer rows.Close()
		var sstreet, scity, sstate, spcode, scountry sql.NullString
		for rows.Next() {
			var si model.ShippingInfo
			err = rows.Scan(&si.S_ShippingID, &sstreet, &scity, &sstate, &spcode, &scountry, &si.S_Primary)
			si.S_Street = sstreet.String
			si.S_City = scity.String
			si.S_State = sstate.String
			si.S_PostalCode = spcode.String
			si.S_Country = scountry.String
			shipInfo := append(con.ShippingInfo, si)
			con.ShippingInfo = shipInfo
			log.Println("added shipping info")
			log.Println(con.ShippingInfo)
		}
	} else if con.UpdateConBilling {
		log.Println("Update Account Contact Billing Module Entered")

		for _, bInfo := range con.BillingInfo {
			sqlStatementcbu1 := `UPDATE dbo.accounts_billing_address_master
									SET 
									street=$1, 
									city=$2, 
									stateprovince=$3, 
									postalcode=$4, 
									country=$5,
									primary_address=$6 
									WHERE
									accountid=$7
									and
									billingid=$8`

			rows1, err := account.db.Raw(sqlStatementcbu1,
				NewNullString(bInfo.B_Street),
				NewNullString(bInfo.B_City),
				NewNullString(bInfo.B_State),
				NewNullString(bInfo.B_PostalCode),
				NewNullString(bInfo.B_Country),
				bInfo.B_Primary,
				con.AccountID,
				bInfo.B_BillingID).Rows()

			if err != nil {
				log.Println("Insert billing info into table failed", rows1)
				log.Println(err.Error(), rows1)
			}
		}
	} else if con.UpdateConShipping {
		log.Println("Update Account Contact Shipping Module Entered", con.ShippingInfo)
		for _, sInfo := range con.ShippingInfo {
			log.Println("Started updating", sInfo, " in AccountId-", con.AccountID)
			sqlStatementcsu1 := `UPDATE dbo.accounts_shipping_address_master
									SET 
									street=$1, 
									city=$2, 
									stateprovince=$3, 
									postalcode=$4, 
									country=$5,
									primary_address=$6 
									WHERE
									accountid=$7
									and
									shippingid=$8`
			rows1, err := account.db.Raw(sqlStatementcsu1,
				NewNullString(sInfo.S_Street),
				NewNullString(sInfo.S_City),
				NewNullString(sInfo.S_State),
				NewNullString(sInfo.S_PostalCode),
				NewNullString(sInfo.S_Country),
				sInfo.S_Primary,
				con.AccountID,
				sInfo.S_ShippingID).Rows()
			if err != nil {
				log.Println("Insert Shipping info into table failed", rows1)
				log.Println(err.Error(), rows1)
			}
		}
	} else if con.CreateConBilling {
		log.Println("Entered Create Contact Billing")
		var Pkey string
		//Generating PO NOs----------------
		// Pkey := "BillID-" + strconv.Itoa(FindLatestSerial("idsno", "dbo.accounts_billing_address_master", "idsno", "idsno"))
		log.Println(Pkey)
		for _, bInfo := range con.BillingInfo {
			log.Println("Started updating", bInfo, " in AccountId-", con.AccountID)
			sqlStatementCCB1 := `INSERT INTO dbo.accounts_billing_address_master(
										 accountid,
										 contactid,
										 street, 
										 city, 
										 stateprovince,
										 postalcode, 
										 country, 
										 billingid, 
										 primary_address)
									VALUES($1, $2, $3, $4, $5, $6, $7,$8,$9)`
			rows1, err := account.db.Raw(sqlStatementCCB1,
				con.AccountID,
				con.ContactID,
				bInfo.B_Street,
				bInfo.B_City,
				bInfo.B_State,
				bInfo.B_PostalCode, bInfo.B_Country, Pkey, bInfo.B_Primary).Rows()
			if err != nil {
				log.Println("Insert billing info into table failed", rows1)
				log.Println(err.Error())
			}
		}

		sqlStatementCCB2 := `update dbo.accounts_billing_address_master
								set
								primary_address=false
								where
								billingid !=$1 
								and
								accountid=$2`
		_, err2 := account.db.Raw(sqlStatementCCB2, Pkey, con.AccountID).Rows()
		if err2 != nil {
			log.Println("Setting primary address as false for the accountid failed")
			log.Println(err2.Error())
		}
		log.Println("Updated primary billing address to contact")
		sqlStatementCCB3 := `update dbo.contacts_master
								set
								billingaddressid=$1
								where
								contactid=$2`
		_, err3 := account.db.Raw(sqlStatementCCB3, Pkey, con.ContactID).Rows()
		if err3 != nil {
			log.Println("Updated primary address to contact for accountid failed")
			log.Println(err3.Error())
		}
	} else if con.CreateConShipping {
		// var si ShippingInfo
		// log.Println(con)

		// Pkey := "ShipID-" + strconv.Itoa(FindLatestSerial("idsno", "dbo.accounts_shipping_address_master", "idsno", "idsno"))
		var Pkey string
		log.Println("Entered Create Contact Shipping")
		for _, sInfo := range con.ShippingInfo {
			log.Println("Started updating", sInfo, " in AccountId-", con.AccountID)
			sqlStatementCCS1 := `INSERT INTO dbo.accounts_shipping_address_master(
								accountid, contactid,shippingid,street, city, stateprovince,
								postalcode, country,primary_address)
							VALUES($1, $2, $3, $4, $5, $6, $7,$8,$9)`

			_, err := account.db.Raw(sqlStatementCCS1,
				con.AccountID,
				con.ContactID,
				Pkey,
				sInfo.S_Street,
				sInfo.S_City,
				sInfo.S_State,
				sInfo.S_PostalCode,
				sInfo.S_Country,
				sInfo.S_Primary).Rows()
			if err != nil {
				log.Println("Insert Shipping info into table failed")
				log.Println(err.Error())
			}
		}
		sqlStatementCCS2 := `update dbo.accounts_shipping_address_master
								set
								primary_address=false
								where
								shippingid !=$1 
								and
								accountid=$2`
		_, err2 := account.db.Raw(sqlStatementCCS2, Pkey, con.AccountID).Rows()
		if err2 != nil {
			log.Println("Setting primary address as false for the accountid failed")
			log.Println(err2.Error())
		}
		log.Println("Updated primary Shipping address to contact")
		sqlStatementCCS3 := `update dbo.contacts_master
								set
								shippingaddressid=$1
								where
								contactid=$2`
		_, err3 := account.db.Raw(sqlStatementCCS3, Pkey, con.ContactID).Rows()
		if err3 != nil {
			log.Println("Updated primary address to contact for accountid failed")
			log.Println(err3.Error())
		}
		// return nil, nil

		// func (account accountRepo) FindLatestSerial(param1, param2, param3, param4, string) (ids, int){
		// 	log.Println("Finding latest serial num")
		// 	// db, _ := sql.Open("postgres", PsqlInfo)

		// 	var rows *sql.Rows
		// 	sqlStatement1 := fmt.Sprintf("SELECT %s FROM %s where %s is not null ORDER BY %s DESC LIMIT 1", param1, param2, param3, param4)
		// 	rows, err := account.db.Raw(sqlStatement1).Rows()
		// 	for rows.Next() {
		// 		err = rows.Scan(&ids)
		// 	}
		// 	if err != nil {
		// 	log.Println(err)
		// 	}
		// 	return ids + 1
		// }
	}
	return nil, nil
}