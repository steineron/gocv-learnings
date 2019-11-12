// What it does:
//
// This example shows how to find lines in an image using Hough transform.
//
// How to run:
//
// 		go run ./cmd/find-lines/main.go lines.jpg
//
// build example

package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	filename := os.Args[1]
	/*
		ratio = image.shape[0] / 500.0
		orig = image.copy()
		image = imutils.resize(image, height = 500)

		# convert the image to grayscale, blur it, and find edges
		# in the image
		gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
		gray = cv2.GaussianBlur(gray, (5, 5), 0)
		edged = cv2.Canny(gray, 75, 200)

		# show the original image and the edge detected image
		print("STEP 1: Edge Detection")
		cv2.imshow("Image", image)
		cv2.imshow("Edged", edged)
		cv2.waitKey(0)
		cv2.destroyAllWindows()
	*/

	mat := gocv.IMRead(filename, gocv.IMReadColor)

	//ratio := mat.Size()[0] / 10.0
	orig := mat.Clone()
	//gocv.Resize(mat, &mat, image.Point{X: ratio, Y: 10.0}, 10, 10, gocv.InterpolationLinear)

	// convert the image to grayscale, blur it, and find edges
	// in the image
	gray := gocv.NewMat()
	gocv.CvtColor(mat, &gray, gocv.ColorBGRToGray)

	gocv.GaussianBlur(gray,&gray, image.Point{X:5, Y:5}, 0,0,gocv.BorderDefault)
	edged:= gocv.NewMat()
	gocv.Canny(gray, &edged, 75.0, 200.0)

	window := gocv.NewWindow("out-lines")
	window2 := gocv.NewWindow("origin")
	for {
		window2.IMShow(mat)
		window.IMShow(edged)
		if window.WaitKey(10) >= 0 {
			break
		}
	}


	mat = orig
	matCanny := gocv.NewMat()
	matLines := gocv.NewMat()

	window = gocv.NewWindow("detected lines")

	gocv.Canny(mat, &matCanny, 50, 200)
	gocv.HoughLinesP(matCanny, &matLines, 1, math.Pi/180, 80)

	fmt.Println(matLines.Cols())
	fmt.Println(matLines.Rows())
	for i := 0; i < matLines.Rows(); i++ {
		pt1 := image.Pt(int(matLines.GetVeciAt(i, 0)[0]), int(matLines.GetVeciAt(i, 0)[1]))
		pt2 := image.Pt(int(matLines.GetVeciAt(i, 0)[2]), int(matLines.GetVeciAt(i, 0)[3]))
		gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 10)
	}

	for {
		window.IMShow(mat)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}
