package adrepo

import (
	"homework6/internal/ads"
	"homework6/internal/app"
)

type SliceRepository struct {
	ads []ads.Ad
}

func (r *SliceRepository) CreateAd(title string, text string, userID int64) int64 {
	newID := int64(len(r.ads))
	r.ads = append(r.ads, ads.Ad{ID: newID, Title: title, Text: text, AuthorID: userID, Published: false})
	return newID
}

func (r *SliceRepository) ChangeAdStatus(adID int64, status bool) {
	r.ads[adID].Published = status
}

func (r *SliceRepository) UpdateAd(adID int64, title string, text string) {
	ad := r.ads[adID]
	ad.Title = title
	ad.Text = text
	r.ads[adID] = ad
}

func (r *SliceRepository) GetAd(adID int64) ads.Ad {
	return r.ads[adID]
}

func New() app.Repository {
	return &SliceRepository{make([]ads.Ad, 0)}
}
