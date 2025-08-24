package comps

import (
	"fmt"

	gc "github.com/abdheshnayak/gohtmlx/example/dist/gohtmlxc"
	"github.com/abdheshnayak/gohtmlx/example/src/types"
	. "github.com/abdheshnayak/gohtmlx/pkg/element"
)

func Home() Element {
	items := []types.Table{}
	for i := range 50 {
		items = append(items, types.Table{
			Id: fmt.Sprintf("%d", i+1),
		})
	}
	return gc.Home{
		Tables: items,
		Attrs: Attrs{
			"Items": items,
		},
	}.Get()
}
