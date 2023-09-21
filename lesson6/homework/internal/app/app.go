package app

import (
	"fmt"
	"homework6/internal/ads"
)

var ErrBadRequest = fmt.Errorf("bad request")
var ErrForbidden = fmt.Errorf("forbidden")

type App interface {
	CreateAd(title string, text string, userID int64) (ads.Ad, error)
	ChangeAdStatus(adID int64, userID int64, status bool) (ads.Ad, error)
	UpdateAd(adID int64, userID int64, title string, text string) (ads.Ad, error)
}

type Repository interface {
	CreateAd(title string, text string, userID int64) int64
	ChangeAdStatus(adID int64, status bool)
	UpdateAd(adID int64, title string, text string)
	GetAd(adID int64) ads.Ad
}

type ValidatingApp struct {
	repo Repository
}

func (a ValidatingApp) CreateAd(title string, text string, userID int64) (ads.Ad, error) {
	err := ValidateTitleAndText(title, text)
	if err != nil {
		return ads.Ad{}, err
	}
	adID := a.repo.CreateAd(title, text, userID)
	return a.repo.GetAd(adID), nil
}

func ValidateTitle(title string) error {
	if len(title) == 0 {
		return fmt.Errorf("%w: empty title", ErrBadRequest)
	}
	if len(title) > 100 {
		return fmt.Errorf("%w: too big title title", ErrBadRequest)
	}
	return nil
}

func ValidateText(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("%w: empty text", ErrBadRequest)
	}
	if len(text) > 500 {
		return fmt.Errorf("%w: too big text", ErrBadRequest)
	}
	return nil
}

func ValidateTitleAndText(title string, text string) error {
	err := ValidateTitle(title)
	if err != nil {
		return err
	}
	err = ValidateText(text)
	if err != nil {
		return err
	}
	return nil
}

func (a ValidatingApp) ChangeAdStatus(adID int64, userID int64, status bool) (ads.Ad, error) {
	if a.repo.GetAd(adID).AuthorID != userID {
		return ads.Ad{}, fmt.Errorf("%w: can't change other's adds", ErrForbidden)
	}
	a.repo.ChangeAdStatus(adID, status)
	return a.repo.GetAd(adID), nil
}

func (a ValidatingApp) UpdateAd(adID int64, userID int64, title string, text string) (ads.Ad, error) {
	if a.repo.GetAd(adID).AuthorID != userID {
		return ads.Ad{}, fmt.Errorf("%w: can't change other's adds", ErrForbidden)
	}
	err := ValidateTitleAndText(title, text)
	if err != nil {
		return ads.Ad{}, err
	}
	a.repo.UpdateAd(adID, title, text)
	return a.repo.GetAd(adID), nil
}

func NewApp(repo Repository) App {
	return ValidatingApp{repo}
}
