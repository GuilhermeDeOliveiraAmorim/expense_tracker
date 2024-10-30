package pdfservice

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/gen2brain/go-fitz"
	"github.com/otiai10/gosseract/v2"
	"gorm.io/gorm"
)

type PdfImporterRepository struct {
	gorm *gorm.DB
}

func NewPdfImporterRepository(gorm *gorm.DB) *PdfImporterRepository {
	return &PdfImporterRepository{
		gorm: gorm,
	}
}

func (p *PdfImporterRepository) ImportPdf(userID string, defaultCategoryID string, rootPath string) error {
	files, err := p.getPdfFiles(rootPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		folder := strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))
		imgDir := filepath.Join("img", folder)

		if err := p.convertPdfToImages(file, imgDir); err != nil {
			return err
		}

		text, err := p.extractTextFromImages(imgDir)
		if err != nil {
			return err
		}

		padrao := regexp.MustCompile(`(\d{1,2} de \w{3}\. \d{4})\s+([\w\s\*\(\)\.\-]+)\s+[-+]?\s*R\$ ([\d,.]+)`)

		resultados := padrao.FindAllStringSubmatch(text, -1)

		meses := map[string]string{
			"jan.": "01",
			"fev.": "02",
			"mar.": "03",
			"abr.": "04",
			"mai.": "05",
			"jun.": "06",
			"jul.": "07",
			"ago.": "08",
			"set.": "09",
			"out.": "10",
			"nov.": "11",
			"dez.": "12",
		}

		expenses := []entities.Expense{}

		for _, resultado := range resultados {
			if len(resultado) == 4 {
				dataStr := resultado[1]
				beneficiario := strings.TrimSpace(resultado[2])
				valorStr := resultado[3]

				parteData := strings.Split(dataStr, " ")
				if len(parteData) == 4 {
					mesStr := parteData[2]
					if mesNum, ok := meses[mesStr]; ok {
						dataStr = fmt.Sprintf("%s-%s-%s", parteData[0], mesNum, parteData[3])
					}
				}

				data, err := time.Parse("02-01-2006", dataStr)
				if err != nil {
					fmt.Printf("Erro ao converter data: %v\n", err)
					continue
				}

				valor, err := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(valorStr, ".", ""), ",", "."), 64)
				if err != nil {
					fmt.Printf("Erro ao converter valor: %v\n", err)
					continue
				}

				newExpense, errNewExpense := entities.NewExpense(userID, valor, data, defaultCategoryID, beneficiario)
				if errNewExpense != nil {
					fmt.Printf("Erro ao criar nova despesa: %v\n", errNewExpense)
					continue
				}

				expenses = append(expenses, *newExpense)
			}
		}

		errCreatesMultipleExpenses := p.CreatesMultipleExpenses(expenses)
		if errCreatesMultipleExpenses != nil {
			return errCreatesMultipleExpenses
		}

		if err := p.cleanUp(imgDir); err != nil {
			return err
		}
	}

	return nil
}

func (p *PdfImporterRepository) CreatesMultipleExpenses(expenses []entities.Expense) error {
	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	for _, expense := range expenses {
		if err := tx.Create(&repositoriesgorm.Expenses{
			ID:            expense.ID,
			Active:        expense.Active,
			CreatedAt:     expense.CreatedAt,
			UpdatedAt:     expense.UpdatedAt,
			DeactivatedAt: expense.DeactivatedAt,
			UserID:        expense.UserID,
			Amount:        expense.Amount,
			ExpanseDate:   expense.ExpenseDate,
			CategoryID:    expense.CategoryID,
			Notes:         expense.Notes,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (i *PdfImporterRepository) getPdfFiles(rootPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".pdf" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (i *PdfImporterRepository) convertPdfToImages(pdfPath, imgDir string) error {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return err
	}
	defer doc.Close()

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			return err
		}

		err = os.MkdirAll(imgDir, 0755)
		if err != nil {
			return err
		}

		imgPath := filepath.Join(imgDir, fmt.Sprintf("image-%05d.png", n))
		f, err := os.Create(imgPath)
		if err != nil {
			return err
		}

		err = png.Encode(f, img)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

func (i *PdfImporterRepository) extractTextFromImages(imgDir string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	var extractedText strings.Builder

	err := filepath.Walk(imgDir, func(imgPath string, info os.FileInfo, err error) error {
		if filepath.Ext(imgPath) == ".png" {
			client.SetImage(imgPath)
			text, err := client.Text()
			if err != nil {
				return err
			}
			extractedText.WriteString(text + "\n")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return extractedText.String(), nil
}

func (i *PdfImporterRepository) cleanUp(imgDir string) error {
	return os.RemoveAll(imgDir)
}
