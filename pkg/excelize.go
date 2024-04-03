package pkg

import (
	"unicode/utf8"

	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/xuri/excelize/v2"
)

type ExcelizePkg interface {
	InitExcelize(sheetName string) (*excelize.File, string)
	GetCellWidth(value string, margin int) int
	StyleBold(f *excelize.File) int
	StyleHeader(f *excelize.File, align string, colorFill string) int
	StyleCurrencyRpIndo(f *excelize.File, bold bool, colorFont string, alignFont string, fill bool, colorFill string) int
	StyleFill(f *excelize.File, colorFill string) int
}
type excelizePkg struct{}

func NewExcelizePkg() ExcelizePkg {
	return new(excelizePkg)
}

func (e *excelizePkg) InitExcelize(sheetName string) (*excelize.File, string) {
	f := excelize.NewFile()
	f.SetSheetName(f.GetSheetName(0), sheetName)
	return f, sheetName
}
func (e *excelizePkg) GetCellWidth(value string, margin int) int {
	return utf8.RuneCountInString(value) + margin
}

func (e *excelizePkg) StyleBold(f *excelize.File) int {
	id, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		helper.LogsError(err)
		return 0
	}
	return id
}
func (e *excelizePkg) StyleHeader(f *excelize.File, align string, colorFill string) int {
	id, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: align,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{colorFill},
			Pattern: 1,
		},
	})
	if err != nil {
		helper.LogsError(err)
		return 0
	}
	return id
}
func (e *excelizePkg) StyleCurrencyRpIndo(f *excelize.File, bold bool, colorFont string, alignFont string, fill bool, colorFill string) int {
	style := &excelize.Style{
		NumFmt: 359,
		Font: &excelize.Font{
			Bold:  bold,
			Color: colorFont,
		},
		Alignment: &excelize.Alignment{
			Horizontal: alignFont,
		},
	}

	if fill {
		style.Fill = excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
		}
		if colorFill != "" {
			style.Fill.Color = []string{colorFill}
		}
	}

	id, err := f.NewStyle(style)
	if err != nil {
		helper.LogsError(err)
		return 0
	}
	return id
}
func (e *excelizePkg) StyleFill(f *excelize.File, colorFill string) int {
	style := &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{colorFill},
		},
	}
	id, err := f.NewStyle(style)
	if err != nil {
		helper.LogsError(err)
		return 0
	}
	return id
}
