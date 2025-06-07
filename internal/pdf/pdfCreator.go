package pdf

import (
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf/v2"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
	"log/slog"
)

func CreatePDF(db *gorm.DB, scores []models.Score) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("DejaVu", "", "static/fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "static/fonts/DejaVuSans-Bold.ttf")
	pdf.AddUTF8Font("DejaVu", "I", "static/fonts/DejaVuSans-Oblique.ttf")
	pdf.AddUTF8Font("DejaVu", "BI", "static/fonts/DejaVuSans-BoldOblique.ttf")
	pdf.SetFont("DejaVu", "", 12)
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("DejaVu", "I", 8)
		pdf.CellFormat(0, 10, "© 2025 BST - BS Thalwil, Martin Knoller Stocker. https://github.com/linuxswords/hat", "", 0, "C", false, 0, "")
	})
	pdf.AddPage()
	pdf.ImageOptions("static/images/hat-logo.png", 170, 10, 30, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	if len(scores) == 0 {
		return nil, fmt.Errorf("no scores available")
	}

	var tournament models.Tournament
	if err := db.Preload("HandicapSet").First(&tournament, scores[0].TournamentID).Error; err != nil {
		slog.Error("error loading handicapset", "error", err)
		return nil, err
	}

	pdf.SetFont("DejaVu", "", 16)
	pdf.Cell(40, 10, "Tournament Results")
	pdf.Ln(10)
	pdf.SetFont("DejaVu", "", 12)
	pdf.Cell(40, 10, "Tournament Name: "+tournament.Name)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Date: "+tournament.Date.Format("2006-01-02"))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Type: "+tournament.TournamentType)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Venue: "+tournament.Venue)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Handicap Set: "+tournament.HandicapSet.Name)
	pdf.Ln(20)

	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(15, 7, "Rank", "1", 0, "C", false, 0, "")
	pdf.CellFormat(65, 7, "Name", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Bow Class Code", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Score", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Handicap", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Total Score", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	for i, score := range scores {
		var archer models.Archer
		if err := db.Preload("BowClass").First(&archer, score.ArcherID).Error; err != nil {
			slog.Error("error loading bowclass", "error", err)
			return nil, err
		}

		var handicapEntry models.HandicapEntry
		if err := db.Where("handicap_set_id = ? AND bow_class_id = ?", tournament.HandicapSetID, archer.BowClassID).First(&handicapEntry).Error; err != nil {
			slog.Error("error loading handicap entry", "error", err)
			fmt.Printf("Archer %s %s, and bowclass %s", archer.FirstName, archer.LastName, archer.BowClass.Code)
			return nil, err
		}
		if i < 3 {
			pdf.SetFont("DejaVu", "B", 10)
		} else {
			pdf.SetFont("DejaVu", "", 10)
		}
		pdf.CellFormat(15, 6, fmt.Sprintf("%d", score.Ranking), "1", 0, "C", false, 0, "")
		pdf.CellFormat(65, 6, fmt.Sprintf("%s %s", archer.FirstName, archer.LastName), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, archer.BowClass.Code, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", score.Score), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", handicapEntry.Value), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", score.TotalScore), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		slog.Error("error byte buffer", "error", err)
		return nil, err
	}
	pdf.SetY(-30)
	pdf.SetFont("DejaVu", "I", 8)
	pdf.CellFormat(0, 10, "© 2025 BST - BS Thalwil, Martin Knoller Stocker. https://github.com/linuxswords/hat", "", 0, "C", false, 0, "")
	return buf.Bytes(), nil
}
