package views

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
)

const (
	CHANNELS         = "channels"
	CHANNELS_COLUMNS = "channels_columns"
)

type Channels struct {
	*gocui.View
	items []*models.Channel
}

func (c *Channels) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	columns, err := g.SetView(CHANNELS_COLUMNS, x0, y0, x1, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	columns.Frame = false
	columns.BgColor = gocui.ColorGreen
	columns.FgColor = gocui.ColorBlack | gocui.AttrBold
	displayChannelsColumns(columns)

	c.View, err = g.SetView(CHANNELS, x0, y0+1, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		_, err = g.SetCurrentView(CHANNELS)
		if err != nil {
			return err
		}
	}
	c.View.Frame = false
	c.View.Autoscroll = true
	c.View.SelBgColor = gocui.ColorCyan
	c.View.SelFgColor = gocui.ColorBlack
	c.Highlight = true

	c.display()
	return nil
}

func displayChannelsColumns(v *gocui.View) {
	fmt.Fprintln(v, fmt.Sprintf("%-9s %-26s %12s %12s %5s",
		"Status",
		"Gauge",
		"Local",
		"Capacity",
		"pHTLC",
	))
}

func (c *Channels) Update(items []*models.Channel) {
	c.items = items
}

func (c *Channels) display() {
	c.Clear()
	for _, item := range c.items {
		line := fmt.Sprintf("%s %s %s %12d %5d",
			active(item),
			gauge(item),
			color.Cyan(fmt.Sprintf("%12d", item.LocalBalance)),
			item.Capacity,
			len(item.PendingHTLC),
		)
		fmt.Fprintln(c.View, line)
	}
}

func active(c *models.Channel) string {
	if c.Active {
		return color.Green(fmt.Sprintf("%-9s", "active"))
	}
	return color.Red(fmt.Sprintf("%-9s", "inactive"))
}

func gauge(c *models.Channel) string {
	index := int(c.LocalBalance * int64(20) / c.Capacity)
	var buffer bytes.Buffer
	for i := 0; i < 20; i++ {
		if i < index {
			buffer.WriteString(color.Cyan("|"))
			continue
		}
		buffer.WriteString(" ")
	}
	return fmt.Sprintf("[%s] %2d%%", buffer.String(), c.LocalBalance*100/c.Capacity)
}

func NewChannels() *Channels {
	return &Channels{}
}
