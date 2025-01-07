package printer

import (
	"fmt"
	"strings"

	"github.com/to77e/paslok/internal/models"
)

func PrintResources(in []models.Resource) string {
	out := strings.Builder{}
	for _, i := range in {
		out.WriteString("-----------------------------------------\n")
		out.WriteString(fmt.Sprintf("[%d] Service: %s\n", i.Id, i.Service))
		out.WriteString(fmt.Sprintf("    Username: %s\n", i.Username))
		out.WriteString(fmt.Sprintf("    Created: %s\n", i.CreatedAt))
		if i.Comment != "" {
			out.WriteString(fmt.Sprintf("    Comment: %s\n", i.Comment))
		}
	}
	if len(in) > 0 {
		out.WriteString("-----------------------------------------\n")
	}
	return out.String()
}
