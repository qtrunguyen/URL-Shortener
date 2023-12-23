package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"

	helpers "github.com/qtrunguyen/URL-Shortener/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`

	// XRateRemaining to prevent excessive amount of request
	XRateRemaining int `json:"rate_limit"`

	// For each user IP, reset rate limit for the time duration of XRateLimitReset
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(context *fiber.Ctx) error {
	body := new(request)

	if err := context.BodyParser(&body); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	if !govalidator.IsURL(body.URL) {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}

	if helpers.RemoveDomainError(body.URL) {
		return context.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "invalid domain"})
	}

	body.URL = helpers.EnforceHTTP(body.URL)
}
