package htmltopdf

import (
	"fmt"

	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
)

func GenerateRevenueSharingContract(v models.Video, user models.User) error {
	pdf := NewPDFService()
	data, err := fetchDataForRevenueSharingContract(v, user)
	if err != nil {
		return err
	}
	contractURI, err := pdf.GenerateRevenueSharingContract(data)
	if err != nil {
		fmt.Printf("error in generating pdf: %v\n", err.Error())
		return err
	}

	// update the database record for this video with the address of the contract
	newVideo := v
	newVideo.RevenueSharingContractLink = contractURI

	tx := initializers.DB.Model(&newVideo).Where("id = ?", v.ID).UpdateColumns(&newVideo)
	if tx.Error != nil {
		fmt.Printf("something went wrong while saving contract link: %v\n", tx.Error.Error())
		return tx.Error
	}
	return nil
}

func fetchDataForRevenueSharingContract(v models.Video, user models.User) (*models.RevenueSharingPdfData, error) {
	revenueShared := fmt.Sprintf("%.3v", v.RevenueShared/float32(v.TokensOffered))
	tokensReleased := fmt.Sprintf("%v", v.TokensOffered)
	tokenPrice := fmt.Sprintf("$%v", v.TokenPrice)
	data := &models.RevenueSharingPdfData{
		VideoID:        v.ID.String(),
		CreatorID:      v.CreatorID.String(),
		CreatorName:    user.Name,
		VideoName:      v.Name,
		ReleaseDate:    v.ReleaseDate.String(),
		RevenueShared:  revenueShared,
		TokensReleased: tokensReleased,
		TokenPrice:     tokenPrice,
	}
	return data, nil
}
