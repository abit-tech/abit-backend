package htmltopdf

import (
	"fmt"

	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
)

func GenerateTokenOwnershipContract(t models.Token, v models.Video, c models.User, p models.User) error {
	pdf := NewPDFService()
	data, err := fetchDataForTokenOwnershipContract(t, v, c, p)
	if err != nil {
		return err
	}
	contractURI, err := pdf.GenerateOwnershipContract(data)
	if err != nil {
		fmt.Printf("error in generating pdf: %v\n", err.Error())
		return err
	}

	// update the database record for this token with the address of the contract
	newToken := t
	newToken.OwnershipContractLink = contractURI

	tx := initializers.DB.Model(&newToken).Where("id = ?", t.ID).UpdateColumns(&newToken)
	if tx.Error != nil {
		fmt.Printf("something went wrong while saving contract link: %v\n", tx.Error.Error())
		return tx.Error
	}
	return nil
}

func fetchDataForTokenOwnershipContract(t models.Token, v models.Video, c models.User, p models.User) (*models.OwnershipPdfData, error) {
	revenueShared := fmt.Sprintf("%.3v", v.RevenueShared/float32(v.TokensOffered))
	tokenPrice := fmt.Sprintf("$%v", v.TokenPrice)
	data := &models.OwnershipPdfData{
		VideoID:      v.ID.String(),
		TokenID:      t.ID.String(),
		OwnerID:      p.ID.String(),
		CreatorName:  c.Name,
		OwnerName:    p.Name,
		VideoName:    v.Name,
		ReleaseDate:  v.ReleaseDate.String(),
		RevenueShare: revenueShared,
		TokenPrice:   tokenPrice,
	}
	return data, nil
}
