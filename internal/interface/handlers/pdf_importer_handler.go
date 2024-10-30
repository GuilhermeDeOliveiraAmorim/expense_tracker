package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
	"github.com/gin-gonic/gin"
)

type PdfImporterHandler struct {
	pdfImporterFactory *factory.PdfImporterFactory
}

func NewPdfImporterHandler(factory *factory.PdfImporterFactory) *PdfImporterHandler {
	return &PdfImporterHandler{
		pdfImporterFactory: factory,
	}
}

func (h *PdfImporterHandler) ImportPdf(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	defaultCategoryID := c.Query("default_category_id")
	if defaultCategoryID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing start date",
			Status:   http.StatusBadRequest,
			Detail:   "Start date is required",
			Instance: util.RFC400,
		}})
		return
	}

	pdfData, errReadAll := io.ReadAll(c.Request.Body)
	if errReadAll != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}

	tempFile, errTemp := os.CreateTemp("", "uploaded-*.pdf")
	if errTemp != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(pdfData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write PDF to file"})
		return
	}
	tempFile.Close()

	input := usecases.PdfImporterInputDto{
		UserID:            userID,
		DefaultCategoryID: defaultCategoryID,
		Path:              tempFile.Name(),
	}

	output, errs := h.pdfImporterFactory.PdfImporter.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}
