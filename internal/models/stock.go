package models

type Stock struct {
	Id            *int    `json:"id,omitempty"`
	Code          string  `json:"code"`
	CompanyName   string  `json:"company_name"`
	ListingDate   *string `json:"listing_date,omitempty"`
	DelistingDate *string `json:"delisting_date,omitempty"`
	ListingBoard  string  `json:"listing_board"`
	Shares        int64   `json:"shares"`
	SectorId      *int    `json:"sector_id,omitempty"`
	SectorSubId   *int    `json:"sector_sub_id,omitempty"`
	IndustryId    *int    `json:"industry_id,omitempty"`
	IndustrySubId *int    `json:"industry_sub_id,omitempty"`
}

type PasardanaStockDetail struct {
	Id                 int      `json:"Id"`
	Code               string   `json:"Code"`
	Name               string   `json:"Name"`
	TotalEmployees     *string  `json:"TotalEmployees"`
	ListingDate        *string  `json:"ListingDate"`
	AnnualDividend     *float64 `json:"AnnualDividend"`
	GeneralInformation *string  `json:"GeneralInformation"`
	FoundingDate       *string  `json:"FoundingDate"`
	FkNewSectorId      *int     `json:"fkNewSectorId"`
	FkNewSubSectorId   *int     `json:"fkNewSubSectorId"`
	FkNewIndustryId    *int     `json:"fkNewIndustryId"`
	FkNewSubIndustryId *int     `json:"fkNewSubIndustryId"`
}

type PasardanaStock struct {
	Id   int    `json:"Id"`
	Code string `json:"Code"`
	Name string `json:"Name"`
}

type ListingBoard struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	NameEn *string `json:"name_en,omitempty"`
}

type StockResponse struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type GetAllStocksFilter struct {
	Board  string `query:"board"`
	Search string `query:"search"`
}

type SyncSectorIndustryResponse struct {
	TotalChecked int `json:"total_checked"`
	TotalUpdated int `json:"total_updated"`
}

type StockDetailResponse struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	ListingDate   string  `json:"listing_date"`
	DelistingDate *string `json:"delisting_date,omitempty"`
	Shares        int64   `json:"shares"`
	Board         string  `json:"board"`
	Sector        *string `json:"sector,omitempty"`
	SubSector     *string `json:"sub_sector,omitempty"`
	Industry      *string `json:"industry,omitempty"`
	SubIndustry   *string `json:"sub_industry,omitempty"`
}
