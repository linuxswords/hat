package pdf

import (
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf/v2"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

func CreatePDF(db *gorm.DB, scores []models.Score) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("DejaVu", "", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf")
	pdf.SetFont("DejaVu", "", 12)
	if len(scores) == 0 {
		return nil, fmt.Errorf("no scores available")
	}

	var tournament models.Tournament
	if err := db.Preload("HandycapSet").First(&tournament, scores[0].TournamentID).Error; err != nil {
		return nil, err
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Tournament Results")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Tournament Name: "+tournament.Name)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Date: "+tournament.Date.Format("2006-01-02"))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Type: "+tournament.TournamentType)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Venue: "+tournament.Venue)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Handycap Set: "+tournament.HandycapSet.Name)
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(15, 7, "Rank", "1", 0, "C", false, 0, "")
	pdf.CellFormat(65, 7, "Name", "1", 0, "C", false, 0, "")
	pdf.CellFormat(65, 7, "Bow Class", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 7, "Score", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	for i, score := range scores {
		var archer models.Archer
		if err := db.Preload("BowClass").First(&archer, score.ArcherID).Error; err != nil {
			return nil, err
		}
		if i < 3 {
			pdf.SetFont("Arial", "B", 10)
		} else {
			pdf.SetFont("Arial", "", 10)
		}
		pdf.CellFormat(15, 6, fmt.Sprintf("%d", score.Ranking), "1", 0, "C", false, 0, "")
		pdf.CellFormat(65, 6, fmt.Sprintf("%s %s", archer.FirstName, archer.LastName), "1", 0, "L", false, 0, "")
		pdf.CellFormat(65, 6, archer.BowClass.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("%.2f", score.TotalScore), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
