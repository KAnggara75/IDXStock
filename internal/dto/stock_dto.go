package dto

type Stock struct {
	Code         string `json:"code"`
	CompanyName  string `json:"company_name"`
	ListingDate  string `json:"listing_date"`
	Shares       int64  `json:"shares"`
	ListingBoard string `json:"listing_board"`
}
