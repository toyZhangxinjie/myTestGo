package canvas

import (
	"myTestGo/exception"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/softwarebackend"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"image/jpeg"
	"math"
	"os"
)

func CanvasDrawWindow() error {

	wnd, cv, err := sdlcanvas.CreateWindow(1280, 720, "Hello")
	if err != nil {
		return exception.NewException(500, err.Error())
	}
	defer wnd.Destroy()

	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)

		for r := 0.0; r < math.Pi*2; r += math.Pi * 0.1 {
			cv.SetFillStyle(int(r*10), int(r*20), int(r*40))
			cv.BeginPath()
			cv.MoveTo(w*0.5, h*0.5)
			cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.4, r, r+0.1*math.Pi, false)
			cv.ClosePath()
			cv.Fill()
		}

		cv.SetFillStyle("#FFF")
		cv.SetLineWidth(10)
		cv.BeginPath()
		cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.4, 0, math.Pi*2, false)
		cv.Stroke()
	})

	return nil
}

func CanvasDrawImg() error {
	backend := softwarebackend.New(500, 240)
	cv := canvas.New(backend)

	w, h := float64(cv.Width()), float64(cv.Height())
	cv.SetFillStyle("#FFF")
	cv.SetStrokeStyle("#FFF")
	cv.SetShadowOffset(0, 0)
	cv.FillRect(0, 0, w, h)

	cv.BeginPath()
	font, err := cv.LoadFont("PingFang.ttc")
	if err != nil {
		return exception.NewException(500, err.Error())
	}
	cv.SetTextAlign(canvas.Center)
	cv.SetTextBaseline(canvas.Middle)
	cv.SetFont(font, 50)
	cv.SetStrokeStyle("#000")
	cv.StrokeText("张新杰", w / 2, h / 2)

	//for r := 0.0; r < math.Pi*2; r += math.Pi * 0.1 {
	//	cv.SetFillStyle(int(r*10), int(r*20), int(r*40))
	//	cv.BeginPath()
	//	cv.MoveTo(w*0.5, h*0.5)
	//	cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.4, r, r+0.1*math.Pi, false)
	//	cv.ClosePath()
	//	cv.Fill()
	//}
	//
	//cv.SetStrokeStyle("#FFF")
	//cv.SetLineWidth(10)
	//cv.BeginPath()
	//cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.4, 0, math.Pi*2, false)


	cv.Stroke()

	f1, err := os.Create("/Users/zhangxinjie/Downloads/result.jpeg")
	if err != nil {
		return exception.NewException(500, err.Error())
	}
	err = jpeg.Encode(f1, backend.Image, nil)
	if err != nil {
		return exception.NewException(500, err.Error())
	}

	return nil
}
