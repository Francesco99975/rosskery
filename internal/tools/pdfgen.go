package tools

import (
	"fmt"
	"strings"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func GenerateInvoice(order *models.Order) (string, error) {
	cfg := config.NewBuilder().Build()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader())
	if err != nil {
		return "", err
	}

	m.AddRows(text.NewRow(10, fmt.Sprintf("Invoice %s", order.Id), props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	m.AddRows(text.NewRow(10, fmt.Sprintf("Pickup Date and Time %s", order.Pickuptime.Format("2006-01-02 03:04 PM")), props.Text{
		Top:   1,
		Style: fontstyle.Italic,
		Align: align.Center,
	}))

	m.AddRow(7,
		text.NewCol(3, "Transactions", props.Text{
			Top:   1.5,
			Size:  9,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: getDarkGrayColor()})

	m.AddRows(getTransactions(order.Purchases)...)

	m.AddRow(40,
		code.NewQrCol(6, order.Id, props.Rect{
			Center:  true,
			Percent: 75,
		}),
	)

	document, err := m.Generate()
	if err != nil {
		return "", err
	}

	filename := strings.ReplaceAll(fmt.Sprintf("%s+%s+%s.pdf", order.Id, order.Method, order.Created.Format("2006-01-02 03:04 PM")), " ", "_")

	err = document.Save(filename)
	if err != nil {
		return "", err
	}

	return filename, err
}

func getPageHeader() core.Row {
	return row.New(20).Add(
		image.NewFromFileCol(3, "static/images/logo.png", props.Rect{
			Center:  true,
			Percent: 80,
		}),
		col.New(6),
		col.New(3).Add(
			text.New("rosadmin@rosskery.com", props.Text{
				Size:  8,
				Align: align.Right,
				Color: getRedColor(),
			}),
			text.New("Tel: 55 024 12345-1234", props.Text{
				Top:   12,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Right,
				Color: getBlueColor(),
			}),
			text.New("www.rosskery.com", props.Text{
				Top:   15,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Right,
				Color: getBlueColor(),
			}),
		),
	)
}

func getTransactions(purchases []models.Purchase) []core.Row {
	rows := []core.Row{
		row.New(5).Add(
			col.New(3),
			text.NewCol(4, "Product", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Quantity", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(3, "Price", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
		),
	}

	var contentsRow []core.Row
	contents := make([][]string, 0)
	for _, purchase := range purchases {
		var rPrice float64
		if purchase.Product.Weighed {
			rPrice = float64(purchase.Product.Price * purchase.Quantity / 10)
		} else {
			rPrice = float64(purchase.Product.Price * purchase.Quantity)
		}
		contents = append(contents, []string{purchase.Product.Name, fmt.Sprint(purchase.Quantity), helpers.FormatPrice(rPrice / 100)})
	}

	for i, content := range contents {
		r := row.New(4).Add(
			col.New(3),
			text.NewCol(4, content[0], props.Text{Size: 8, Align: align.Center}),
			text.NewCol(2, content[1], props.Text{Size: 8, Align: align.Center}),
			text.NewCol(3, content[2], props.Text{Size: 8, Align: align.Center}),
		)
		if i%2 == 0 {
			gray := getGrayColor()
			r.WithStyle(&props.Cell{BackgroundColor: gray})
		}

		contentsRow = append(contentsRow, r)
	}

	rows = append(rows, contentsRow...)

	rows = append(rows, row.New(20).Add(
		col.New(7),
		text.NewCol(2, "Total:", props.Text{
			Top:   5,
			Style: fontstyle.Bold,
			Size:  8,
			Align: align.Right,
		}),
		text.NewCol(3, helpers.FormatPrice(float64(helpers.FoldSlice[models.Purchase, func(models.Purchase, int) int, int](purchases, func(prev models.Purchase, cur int) int {
			if prev.Product.Weighed {
				return prev.Product.Price*prev.Quantity/10 + cur
			} else {
				return prev.Product.Price*prev.Quantity + cur
			}
		}, 0))/100.0), props.Text{
			Top:   5,
			Style: fontstyle.Bold,
			Size:  8,
			Align: align.Center,
		}),
	))

	return rows
}

func getDarkGrayColor() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func getGrayColor() *props.Color {
	return &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getBlueColor() *props.Color {
	return &props.Color{
		Red:   10,
		Green: 10,
		Blue:  150,
	}
}

func getRedColor() *props.Color {
	return &props.Color{
		Red:   150,
		Green: 10,
		Blue:  10,
	}
}
