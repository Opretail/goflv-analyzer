package core

import (
	"fmt"
	"goflv-analyzer/flv"
	"goflv-analyzer/flv/tag"
	syslog "goflv-analyzer/helper"

	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Work ...
func Work() {
	i := newInstance()
	defer i.close()
	{
		var flvSrc FlvSrc
		var err error
		if len(os.Args) < 2 {
			syslog.Clog.Traceln("Please add parameters to the command line! such as [./flvanalyzer ***://***.flv]")
			return
		}
		i.inputParam = os.Args[1]
		i.realTimeStreaming = strings.Contains(i.inputParam, "://")

		if !i.realTimeStreaming {
			flvSrc, err = NewNormalFile(i.inputParam)
			if err != nil {
				syslog.Clog.Errorln("NewNormalFile err", err)
				return
			}
		} else {
			flvSrc, err = NewFlvHTTP(i.inputParam)
			if err != nil {
				syslog.Clog.Errorln("NewFlvHTTP err", err)
				return
			}
		}
		defer flvSrc.Close()
		i.dec, err = flv.NewDecoder(flvSrc.Open())
		if err != nil {
			syslog.Clog.Errorln("Failed to create decoder:", err)
			return
		}
		if i.dec.Header().Flags == flv.FlagsVideo {
			i.audioStatus = false
		} else {
			i.audioStatus = true
		}
	}
	go i.parseFLV()

	err := i.Draw()
	if err != nil {
		syslog.Clog.Errorln("i.Draw error:", err)
		return
	}
	syslog.Clog.Traceln(fmt.Sprintf("Video duration: %.1fs", float32(i.lastVideoTimeStamp)/1000))
	if !i.audioStatus {
		syslog.Clog.Traceln("No audio")
		syslog.Clog.Traceln(fmt.Sprintf("Maximum increment of video frame timestamp: %dms", i.maxInterval))
	} else {
		syslog.Clog.Traceln(fmt.Sprintf("Maximum difference of audio and video timestamps: %dms", i.maxInterval))
	}
}

type instance struct {
	dec *flv.Decoder

	intervalSlice      []float64
	maxInterval        int32
	dataSize           int64
	lastVideoTimeStamp uint32
	lastAudioTimeStamp uint32
	inputParam         string
	realTimeStreaming  bool
	audioStatus        bool
	signalCh           chan struct {
		title string
	}
	endSignal chan struct{}
}

func newInstance() *instance {
	i := &instance{
		dec: new(flv.Decoder),

		intervalSlice:      make([]float64, 0),
		maxInterval:        0,
		dataSize:           0,
		lastVideoTimeStamp: 0,
		lastAudioTimeStamp: 0,

		inputParam:        "",
		realTimeStreaming: false,
		audioStatus:       false,
		signalCh:          make(chan struct{ title string }, 10),
		endSignal:         make(chan struct{}, 1),
	}
	i.intervalSlice = append(i.intervalSlice, 200)
	return i
}

func (i *instance) close() {
	close(i.signalCh)
	close(i.endSignal)
}

func (i *instance) parseFLV() {
	for {
		var flvTag tag.FlvTag
		if err := i.dec.Decode(&flvTag); err != nil {
			if err == io.EOF {
				i.endSignal <- struct{}{}
				break
			}
		} else {
			i.HandleFlvTag(&flvTag)
		}
	}
}

// Draw ...
func (i *instance) Draw() error {
	n := 0
	for {
		if len(i.intervalSlice) > 2 {
			break
		}
		time.Sleep(time.Millisecond * 10)
		n++
		if n >= 1000 {
			if i.lastVideoTimeStamp > 5000 && i.lastAudioTimeStamp == 0 {
				syslog.Clog.Traceln("The stream has audio head, no audio frame")
			}
		}
	}
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer func() {
		ui.Close()
	}()
	lc := widgets.NewPlot()
	if !i.audioStatus {
		lc.Title = "Increment of video frame time stamp (ms)"
		lc.LineColors[0] = ui.ColorGreen
	} else {
		lc.Title = "Time difference of audio and video frames, Always positive (ms)"
		lc.LineColors[0] = ui.ColorYellow
	}
	lc.Data = append(lc.Data, i.intervalSlice)
	lc.AxesColor = ui.ColorWhite
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1,
			ui.NewCol(1, lc),
		),
	)
	ui.Render(grid)
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q":
				return nil
			}
		case <-time.After(time.Second * 5):
			{
				syslog.Clog.Errorln("If there is no difference or increment in 5S, it will exit automatically\n")
				return nil
			}
		case <-i.endSignal:
			time.Sleep(time.Second * 1)
			return nil
		case sig := <-i.signalCh:
			if len(i.intervalSlice) >= termWidth-5 { // -5, So that the latest difference can be shown on the chart
				i.intervalSlice = i.intervalSlice[1:]
				lc.Data[0] = i.intervalSlice
			} else {
				lc.Data[0] = i.intervalSlice
			}
			lc.Title = sig.title
			lc.Data[0][0] = 200
			ui.Render(grid)
			if !i.realTimeStreaming {
				time.Sleep(time.Millisecond * 5)
			}
		}
	}
}

// HandleFlvTag ...
func (i *instance) HandleFlvTag(flvTag *tag.FlvTag) {
	defer flvTag.Close()
	switch flvTag.Type {
	case tag.TypeScriptData:
	case tag.TypeAudio:
		{
			i.lastAudioTimeStamp = flvTag.Timestamp
			if i.lastAudioTimeStamp != 0 && i.lastVideoTimeStamp != 0 {
				title := "[audio > video](ms)"
				dValue := int32(i.lastAudioTimeStamp) - int32(i.lastVideoTimeStamp)
				if dValue < 0 {
					dValue *= -1
					title = "[video > audio](ms)"
				}
				if i.maxInterval < dValue {
					i.maxInterval = dValue
				}
				d := float64(dValue)
				i.intervalSlice = append(i.intervalSlice, d)
				i.signalCh <- struct{ title string }{
					title: title,
				}
			}
		}
	case tag.TypeVideo:
		{
			if !i.audioStatus {
				title := "[video only](ms)"
				dValue := int32(flvTag.Timestamp) - int32(i.lastVideoTimeStamp)
				if i.maxInterval < dValue {
					i.maxInterval = dValue
				}
				d := float64(dValue)
				i.intervalSlice = append(i.intervalSlice, d)
				i.signalCh <- struct{ title string }{
					title: title,
				}
			}
			i.lastVideoTimeStamp = flvTag.Timestamp
		}
	}
}

// FlvSrc ...
type FlvSrc interface {
	Open() io.Reader
	Close()
}

// NormalFile ...
type NormalFile struct {
	F *os.File
}

// NewNormalFile ...
func NewNormalFile(fileName string) (*NormalFile, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return &NormalFile{F: file}, nil
}

// Open ...
func (s *NormalFile) Open() io.Reader {
	return s.F
}

// Close ...
func (s *NormalFile) Close() {
	s.F.Close()
}

// FlvHTTP ...
type FlvHTTP struct {
	R *http.Response
}

// NewFlvHTTP ...
func NewFlvHTTP(url string) (*FlvHTTP, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return &FlvHTTP{R: resp}, nil
}

// Open ...
func (s *FlvHTTP) Open() io.Reader {
	return s.R.Body
}

// Close ...
func (s *FlvHTTP) Close() {
	s.R.Body.Close()
}
