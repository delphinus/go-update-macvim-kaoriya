package gumk

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/cheggaaa/pb.v1"
)

func (g *Gumk) fetch(url string, w io.Writer, bar chan *pb.ProgressBar) error {
	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(g.context)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error in Get")
	}

	if resp.StatusCode/100 != 2 {
		return errors.Errorf("response is not OK. Status: %d %s", resp.StatusCode, resp.Status)
	}

	total, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	b := pb.New64(total).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10).SetWidth(80)
	bar <- b

	go func() {
		<-g.context.Done()
		resp.Body.Close()
		b.Finish()
	}()

	_, err = io.Copy(w, b.NewProxyReader(resp.Body))
	return err
}
