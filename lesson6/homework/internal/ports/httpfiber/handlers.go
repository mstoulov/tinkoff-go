package httpfiber

import (
	"errors"
	"homework6/internal/ads"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"homework6/internal/app"
)

func setStatus(err error, c *fiber.Ctx) {
	if errors.Is(err, app.ErrBadRequest) {
		c.Status(http.StatusBadRequest)
	}
	if errors.Is(err, app.ErrForbidden) {
		c.Status(http.StatusForbidden)
	}
}

// Метод для создания объявления (ad)
func createAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody createAdRequest
		err := c.BodyParser(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		var ad ads.Ad
		ad, err = a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)

		if err != nil {
			setStatus(err, c)
			return c.JSON(AdErrorResponse(err))
		}
		return c.JSON(AdSuccessResponse(&ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		var ad ads.Ad
		ad, err = a.ChangeAdStatus(int64(adID), reqBody.UserID, reqBody.Published)

		if err != nil {
			setStatus(err, c)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(&ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		var ad ads.Ad
		ad, err = a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			setStatus(err, c)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(&ad))
	}
}
