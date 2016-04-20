package signform

import (
	"dfss/dfssc/sign"
	"github.com/visualfc/goqt/ui"
)

var icons map[sign.SignerStatus]*ui.QIcon
var icons_labels = map[sign.SignerStatus]string{
	sign.StatusWaiting:    "Waiting",
	sign.StatusConnecting: "Connecting",
	sign.StatusConnected:  "Connected",
	sign.StatusError:      "Error",
}

var iconsLoaded = false

func loadIcons() {
	if iconsLoaded {
		return
	}

	icons = map[sign.SignerStatus]*ui.QIcon{
		sign.StatusWaiting:    ui.NewIconWithFilename(":/images/time.png"),
		sign.StatusConnecting: ui.NewIconWithFilename(":/images/time.png"),
		sign.StatusConnected:  ui.NewIconWithFilename(":/images/connected.png"),
		sign.StatusError:      ui.NewIconWithFilename(":/images/error.png"),
	}

	iconsLoaded = true
}
