package SSE

import (
	"bytes"
	component "cloud/internal/template/components"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	datastar "github.com/starfederation/datastar/sdk/go"
)

func Feed(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	var buf bytes.Buffer
	for {
		buf.Reset()
		id := uuid.New()
		replaceHTML := component.Replace("feed", id.String(), time.Now().UTC())
		replaceHTML.Render(context.Background(), &buf)

		sse.MergeFragments(
			fmt.Sprint(buf.String()),
			// to replace the entire parent div, use withSelector. Otherwise it will simply replace the span with id feed that is calling this endpoint
			// datastar.WithSelector("div"),
		)

		// you can do send a uuid (or message or notification) once, or at a specified number of times or every second
		// in this example the feed is called on load, but you could also do it on click, etc
		time.Sleep(1 * time.Second)
	}

}
