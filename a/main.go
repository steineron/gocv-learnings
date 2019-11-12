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
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
	"os"
)

type Contour [][]image.Point

func (a Contour) Len() int           { return len(a) }
func (a Contour) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Contour) Less(i, j int) bool { return gocv.ContourArea(a[i]) < gocv.ContourArea(a[j]) }

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
	defer mat.Close()

	orig := mat.Clone()
	defer orig.Close()

	gocv.Resize(mat, &mat, image.Point{X:0,Y:0}, 0.75, 0.75, gocv.InterpolationLinear)

	// convert the image to grayscale, blur it, and find edges
	// in the image
	gray := gocv.NewMat()
	defer gray.Close()

	gocv.CvtColor(mat, &gray, gocv.ColorBGRToGray)

	gocv.GaussianBlur(gray, &gray, image.Point{X: 35, Y: 35}, 0, 0, gocv.BorderDefault)
	gocv.Threshold(gray, &gray, 25, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)
	// remaining cleanup of the image to use for finding contours.
	// first use threshold
	//gocv.Threshold(gray, &gray, 25, 255, gocv.ThresholdBinary)

	window := displayMat(gray, "grayed")

	edged := gocv.NewMat()
	gocv.Canny(gray, &edged, 75.0, 200.0)
	gocv.Canny(gray, &edged, 150.0, 175.0)

	window = displayMat(edged, "edged once")

	// then dilate
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 5))
	defer kernel.Close()
	gocv.Dilate(edged, &edged, kernel)
	window = displayMat(edged, "dilated")

	//gocv.Canny(edged, &edged, 75.0, 200.0)
	//window = displayMat(edged, "edged again")

	//
	//hull := gocv.NewMat()
	//defer hull.Close()


	/*
	   cnts = cv2.findContours(edged.copy(), cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)
	   cnts = imutils.grab_contours(cnts)
	   cnts = sorted(cnts, key = cv2.contourArea, reverse = True)[:5]

	   # loop over the contours
	   for c in cnts:
	   	# approximate the contour
	   	peri = cv2.arcLength(c, True)
	   	approx = cv2.approxPolyDP(c, 0.02 * peri, True)

	   	# if our approximated contour has four points, then we
	   	# can assume that we have found our screen
	   	if len(approx) == 4:
	   		screenCnt = approx
	   		break

	   # show the contour (outline) of the piece of paper
	   print("STEP 2: Find contours of paper")
	   cv2.drawContours(image, [screenCnt], -1, (0, 255, 0), 2)
	   cv2.imshow("Outline", image)
	   cv2.waitKey(0)
	   cv2.destroyAllWindows()
	*/

	contours := gocv.FindContours(edged.Clone(), gocv.RetrievalTree, gocv.ChainApproxSimple)
	//contours=contours[0]
	//sort.Sort(Contour(contours))
	var screenContour [][]image.Point = make([][]image.Point, 1)
	for _, c := range contours {
		peri := gocv.ArcLength(c, true)
		approx := gocv.ApproxPolyDP(c, 0.02*peri, true)

		print("peri:", peri, "\n")
		print("approx:", approx, "\n")

		//defects := gocv.NewMat()
		//defer defects.Close()
		//
		//gocv.ConvexHull(c, &hull, true, false)
		//gocv.ConvexityDefects(c, hull, &defects)
		// if our approximated contour has four points, then we
		// can assume that we have found our screen
		if len(approx) == 4 {
			screenContour[0] = approx
			break
		}
		//screenContour[0] = contours[0]

		//
		//status := "Motion detected"
		//statusColor := color.RGBA{0, 255, 0, 0}
		//gocv.DrawContours(&img, contours, i, statusColor, 2)
		//
		//rect := gocv.BoundingRect(c)
		//gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
	}

	for i := 0; i < len(contours); i++ {
		display := mat.Clone()
		gocv.DrawContours(&display, contours, i, color.RGBA{R: 0, G: 255, B: 0, A: 50}, 2)
		window.IMShow(display)
		if window.WaitKey(0) >= 0 {
			continue
		}
		//display.Close()//
	}

	// detect lines example

	mat = edged//orig
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
		gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 2)
	}

	for {
		window.IMShow(mat)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}

func displayMat(mat gocv.Mat, windowName string) *gocv.Window {
	window := gocv.NewWindow(windowName)
	defer window.Close()
	for {
		///window2.IMShow(mat)
		window.IMShow(mat)
		if window.WaitKey(2) >= 0 {
			break
		}
	}
	return window
}
