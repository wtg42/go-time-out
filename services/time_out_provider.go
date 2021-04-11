package services

import (
	"time"

	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

type TimeoutProvider struct {
	// 每過多久需要休息
	Name      string
	TotalTime int
	progress  *mpb.Progress
	bar       *mpb.Bar
	d         *time.Ticker
}

func (bp *TimeoutProvider) NewProvider() {
	bp.newProgress()
	bp.newBar()
}

func (bp *TimeoutProvider) newProgress() {
	// initialize progress container, with custom width
	bp.progress = mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(10*time.Millisecond),
	)
}

func (bp *TimeoutProvider) newBar() {
	// adding a single bar, which will inherit container's width
	bp.bar = bp.progress.Add(int64(bp.TotalTime),
		// progress bar filler with customized style
		mpb.NewBarFiller("[=>-]"),
		mpb.PrependDecorators(
			// display our name with one space on the right
			decor.Name(bp.Name, decor.WC{W: len(bp.Name) + 10, C: decor.DidentRight}),
			// replace ETA decorator with "done" message, OnComplete event
			decor.OnComplete(
				decor.Counters(0, "% d / % d"), "done",
			),
		),
		mpb.AppendDecorators(decor.Spinner(nil)),
	)
}

func (bp *TimeoutProvider) Incr() {
	bp.bar.Increment()
}

func (bp *TimeoutProvider) Completed() bool {
	return bp.bar.Completed()
}

func (bp *TimeoutProvider) WaitProgress() {
	bp.progress.Wait()
}

func (bp *TimeoutProvider) StartTicker() {
	bp.d = time.NewTicker(time.Second * 1)
	for {
		<-bp.d.C
		bp.Incr()
		if bp.Completed() {
			break
		}
	}
}
