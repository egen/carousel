package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/starkandwayne/carousel/store"

	"github.com/grantae/certinfo"

	humanize "github.com/dustin/go-humanize"
)

func (a *Application) viewDetails() *tview.Flex {
	return tview.NewFlex()
}

func (a *Application) actionShowDetails(ref interface{}) {
	a.layout.details.Clear().AddItem(a.renderDetailsFor(ref), 0, 1, false)
}

func (a *Application) renderDetailsFor(ref interface{}) tview.Primitive {
	switch v := ref.(type) {
	case *store.Cert:
		return a.renderCertDetail(v)
	case *store.CertVersion:
		return a.renderCertVersionDetail(v)
	default:
		return a.renderWelcome()
	}
}

func (a *Application) renderCertDetail(c *store.Cert) tview.Primitive {
	t := tview.NewTable()
	t.SetBorder(true)
	addSimpleRow(t, "ID", c.Id)
	addSimpleRow(t, "Name", c.Name)

	a.layout.tree.SetInputCapture(a.nextFocusIncputCaptureHandler(t))
	t.SetInputCapture(a.nextFocusIncputCaptureHandler(a.layout.tree))
	return t
}

func (a *Application) renderCertVersionDetail(cv *store.CertVersion) tview.Primitive {
	t := tview.NewTable()
	t.SetBorder(true)
	t.SetTitle("Credhub & BOSH")

	addSimpleRow(t, "ID", cv.Id)
	addSimpleRow(t, "Expiry", fmt.Sprintf("%s (%s)",
		cv.Expiry.Format(time.RFC3339),
		humanize.RelTime(cv.Expiry, time.Now(), "ago", "from now")))
	addSimpleRow(t, "Transitional", strconv.FormatBool(cv.Transitional))
	addSimpleRow(t, "Certificate Authority", strconv.FormatBool(cv.CertificateAuthority))
	addSimpleRow(t, "Self Signed", strconv.FormatBool(cv.SelfSigned))

	addSimpleRow(t, "Deployments", renderDeployments(cv.Deployments))

	i, err := certinfo.CertificateText(cv.Certificate)
	if err != nil {
		panic(err)
	}

	info := tview.NewTextView().SetText(i).
		SetTextColor(tcell.Color102)

	info.SetBorder(true)
	info.SetTitle("Raw Certificate")

	a.layout.tree.SetInputCapture(a.nextFocusIncputCaptureHandler(t))
	t.SetInputCapture(a.nextFocusIncputCaptureHandler(info))
	info.SetInputCapture(a.nextFocusIncputCaptureHandler(a.layout.tree))

	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t, 8, 1, false).
		AddItem(info, 0, 1, true)
}

func (a *Application) renderWelcome() tview.Primitive {
	h := tview.NewBox().SetBorder(true).SetTitle("help")

	a.layout.tree.SetInputCapture(a.nextFocusIncputCaptureHandler(h))
	h.SetInputCapture(a.nextFocusIncputCaptureHandler(a.layout.tree))
	return h
}

func addSimpleRow(t *tview.Table, lbl, val string) {
	if val == "" {
		return
	}
	row := t.GetRowCount()
	t.SetCell(row, 0, tview.NewTableCell(lbl).SetStyle(tcell.Style{}.Bold(true)))
	t.SetCellSimple(row, 1, val)
}

func renderDeployments(deployments []*store.Deployment) string {
	tmp := make([]string, 0)
	for _, d := range deployments {
		tmp = append(tmp, d.Name)
	}

	return strings.Join(tmp, ", ")
}