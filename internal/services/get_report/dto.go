package get_report

import (
	"fmt"
	"github.com/frutonanny/wallet-service/pkg"
)

func getServiceName(serviceID int64) string {
	if name, ok := pkg.Services[serviceID]; ok {
		return name
	}
	return fmt.Sprintf("unknown serviceID: %d", serviceID)
}
