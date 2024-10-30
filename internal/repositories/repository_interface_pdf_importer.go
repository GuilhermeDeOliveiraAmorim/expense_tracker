package repositories

type PdfImporterRepositoryInterface interface {
	ImportPdf(userID string, defaultCategoryID string, path string) error
}
