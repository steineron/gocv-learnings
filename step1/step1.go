// What it does:
//
// This example shows how to find lines in an image using Hough transform.
//
// How to run:
//
// 		go run ./cmd/find-lines/main.go lines.jpg
//
// build example

package step1

import (
	"gocv.io/x/gocv"
	"image"
)

type Contour [][]image.Point

func (a Contour) Len() int           { return len(a) }
func (a Contour) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Contour) Less(i, j int) bool { return gocv.ContourArea(a[i]) < gocv.ContourArea(a[j]) }

func main() {
	filename := "~/Desktop/img1.png" //os.Args[1]

	mat := gocv.IMRead(filename, gocv.IMReadColor)

	defer mat.Close()
	//ratio := mat.Size()[0] / 50.0
	orig := mat.Clone()
	defer orig.Close()

	// convert the image to grayscale, blur it, and find edges
	// in the image
	gray := gocv.NewMat()
	defer gray.Close()

	gocv.CvtColor(mat, &gray, gocv.ColorBGRToGray)

	gocv.GaussianBlur(gray, &gray, image.Point{X: 35, Y: 35}, 0, 0, gocv.BorderDefault)
	gocv.Threshold(gray, &gray, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)
	// remaining cleanup of the image to use for finding contours.
	// first use threshold
	//gocv.Threshold(gray, &gray, 25, 255, gocv.ThresholdBinary)

	edged := gocv.NewMat()
	//gocv.Canny(gray, &edged, 75.0, 200.0)

	// then dilate
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	defer kernel.Close()
	gocv.Dilate(edged, &edged, kernel)

	gocv.Canny(gray, &edged, 75.0, 200.0)



	window := gocv.NewWindow("out-lines")
	for {
		window.IMShow(edged)
		if window.WaitKey(2) >= 0 {
			break
		}
	}

}
