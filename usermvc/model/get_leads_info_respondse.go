package model

type GetLeadInfoResponse struct {
	Status  int
	Payload interface{}
}
type LeadInfo struct {
	Leadid           string `json:"leadid"`
	Accountname      string `json:"accountname"`
	Aliases          string `json:"aliases"`
	Contactfirstname string `json:"contactfirstname"`
	Contactlastname  string `json:"contactlastname"`
	Contact_Mobile   string `json:"contact_mobile"`
	Email            string `json:"email"`
	Leadscore        int    `json:"leadscore"`
	Masterstatus     string `json:"masterstatus"`
}

type LeadInfoInLeadToAccount struct {
	LeadName    string `json:"leadname"`
	ContactName string `json:"contact"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	UserName    string `json:"username"`
}

type ContactInfoFromLeadAndMaster struct {
	LeadDetails          LeadInfoInLeadToAccount `json:"leaddetails"`
	ContactMasterDetails Contact                 `json:"erpcontacts"`
}

type Contact struct {
	CustomerName string `json:"customername"`
	ContactName  string `json:"contactname"`
	Country      string `json:"country"`
}
