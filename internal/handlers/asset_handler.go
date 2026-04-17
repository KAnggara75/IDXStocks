package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AssetHandler struct{}

func NewAssetHandler() *AssetHandler {
	return &AssetHandler{}
}

func (h *AssetHandler) GetCompanyLogo(c fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "code is required",
		})
	}

	// Remove .png suffix if present (case insensitive) and normalize to uppercase
	code = strings.TrimSuffix(strings.ToUpper(code), ".PNG")

	targetURL := fmt.Sprintf("https://assets.stockbit.com/logos/companies/%s.png", code)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequestWithContext(c.Context(), "GET", targetURL, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create request",
		})
	}

	// Add realistic user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Failed to proxy logo for %s: %v", code, err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "failed to fetch logo from source",
		})
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "logo not found",
			})
		}
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": fmt.Sprintf("source returned status %d", resp.StatusCode),
		})
	}

	// Set response headers from the source
	c.Set("Content-Type", resp.Header.Get("Content-Type"))
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		c.Set("Content-Length", contentLength)
	}
	// Add caching headers to reduce future proxy requests
	c.Set("Cache-Control", "public, max-age=86400")

	// c.SendStream will read from resp.Body and close it when done
	return c.SendStream(resp.Body)
}
