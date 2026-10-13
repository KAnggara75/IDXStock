package services

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestStockService_ParseExcel(t *testing.T) {
	svc := NewStockService()

	t.Run("invalid_excel_file", func(t *testing.T) {
		buf := bytes.NewBufferString("not an excel file")
		stocks, err := svc.ParseExcel(buf)
		assert.Error(t, err)
		assert.Nil(t, stocks)
	})

	t.Run("valid_excel_file", func(t *testing.T) {
		f := excelize.NewFile()
		sheet := f.GetSheetName(0)

		// Header row
		f.SetSheetRow(sheet, "A1", &[]interface{}{"No", "Kode", "Nama", "Tanggal Pencatatan", "Saham", "Papan Pencatatan"})
		// Data row 1
		f.SetSheetRow(sheet, "A2", &[]interface{}{1, "BBCA", "Bank Central Asia", "1990-01-01", "11,766,313,488", "UTAMA"})
		// Data row 2 with less columns (should be skipped)
		f.SetSheetRow(sheet, "A3", &[]interface{}{2, "INCOMPLETE"})

		buf, err := f.WriteToBuffer()
		assert.NoError(t, err)

		stocks, err := svc.ParseExcel(buf)
		assert.NoError(t, err)
		assert.Len(t, stocks, 1)
		assert.Equal(t, "BBCA", stocks[0].Code)
		assert.Equal(t, "Bank Central Asia", stocks[0].CompanyName)
		assert.Equal(t, int64(11766313488), stocks[0].Shares)
		assert.Equal(t, "Main", stocks[0].ListingBoard)
	})
}
