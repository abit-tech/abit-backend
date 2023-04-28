package htmltopdf

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
)

// PDFService represents the interface of a pdf generation service
type PDFService interface {
	GenerateRevenueSharingContract(data *models.RevenueSharingPdfData) (string, error)
	GenerateOwnershipContract(data *models.OwnershipPdfData) (string, error)
}

type PDFServiceImpl struct{}

func NewPDFService() PDFService {
	return PDFServiceImpl{}
}

func (p PDFServiceImpl) GenerateRevenueSharingContract(data *models.RevenueSharingPdfData) (string, error) {
	var templ *template.Template
	var err error

	conf := initializers.AppConf
	revenueSharingTemplateURI := conf.RevenueSharingContractTemplateURI

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	folderPath := filepath.Dir(ex)

	fmt.Println("current path: ", folderPath)
	templatePath := path.Join(folderPath, revenueSharingTemplateURI)
	fmt.Println("template path: ", templatePath)

	// use Go's default HTML template generation tools to generate your HTML
	if templ, err = template.ParseFiles(templatePath); err != nil {
		return "", err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		return "", err
	}

	// initalize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return "", err
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)

	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)

	// magic
	err = pdfg.Create()
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("/contracts/revenue_sharing/%s.pdf", data.VideoID)
	savePath := path.Join(folderPath, fileName)
	fmt.Println("save path: ", savePath)

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(savePath)
	if err != nil {
		log.Fatal(err)
	}

	return savePath, nil
}

func (p PDFServiceImpl) GenerateOwnershipContract(data *models.OwnershipPdfData) (string, error) {
	var templ *template.Template
	var err error

	conf := initializers.AppConf
	tokenOwnershipContractTemplateURI := conf.TokenOwnershipContractTemplateURI

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	folderPath := filepath.Dir(ex)

	fmt.Println("current path: ", folderPath)
	templatePath := path.Join(folderPath, tokenOwnershipContractTemplateURI)
	fmt.Println("template path: ", templatePath)

	// use Go's default HTML template generation tools to generate your HTML
	if templ, err = template.ParseFiles(templatePath); err != nil {
		return "", err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		return "", err
	}

	// initalize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return "", err
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)

	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)

	// magic
	err = pdfg.Create()
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("/contracts/token_ownership/%s.pdf", data.TokenID)
	savePath := path.Join(folderPath, fileName)
	fmt.Println("save path: ", savePath)

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(savePath)
	if err != nil {
		log.Fatal(err)
	}

	return savePath, nil
}
