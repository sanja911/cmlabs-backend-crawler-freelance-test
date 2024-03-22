package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	// Fungsi untuk melakukan crawling dan menyimpan hasilnya dalam file HTML
	crawlAndSaveHTML := func(c *fiber.Ctx, url string) error {
		hostName := strings.TrimPrefix(url, "http://")
		hostName = strings.TrimPrefix(hostName, "https://")
		hostName = strings.TrimSuffix(hostName, "/")

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}

		// Membuat HTTP Client baru
		client := &http.Client{}

		// Membuat HTTP Request ke URL target
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Gagal membuat HTTP request")
		}

		// Melakukan HTTP request
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Gagal melakukan HTTP request")
		}

		// Membaca isi dari response body
		htmlContent, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Gagal membaca response body")
		}

		// Menutup response body setelah selesai digunakan
		defer resp.Body.Close()

		err = os.MkdirAll(hostName, os.ModePerm)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error saat membuat folder")
		}

		// Menyimpan isi HTML ke dalam file
		fileName := fmt.Sprintf("%s/result.html", hostName)
		err = os.WriteFile(fileName, htmlContent, 0644)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error saat menyimpan hasil ke file")
		}
		// Mengirimkan file HTML sebagai respon
		return c.SendFile(fileName)
	}

	// Handler untuk melakukan crawling dan menyimpan hasilnya dalam file HTML untuk website pertama
	app.Get("/crawl/cmlabs", func(c *fiber.Ctx) error {
		url := "https://cmlabs.co"
		return crawlAndSaveHTML(c, url)
	})

	// Handler untuk melakukan crawling dan menyimpan hasilnya dalam file HTML untuk website kedua
	app.Get("/crawl/sequence", func(c *fiber.Ctx) error {
		url := "https://sequence.day"
		return crawlAndSaveHTML(c, url)
	})

	// Handler untuk melakukan crawling dan menyimpan hasilnya dalam file HTML untuk website ketiga (masukkan website ke tautan yang sesuai)
	app.Get("/crawl/chickin", func(c *fiber.Ctx) error {
		url := "https://chickin.id/"
		return crawlAndSaveHTML(c, url)
	})

	// Jalankan server pada port 3000
	log.Fatal(app.Listen(":3000"))
}
