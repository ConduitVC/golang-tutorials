package identicon

import (
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
)

const scale = 10
const width = 5
const height = 5
const margin = 2

var bgColor = color.RGBA{50, 50, 50, 255}
var fgColor = color.White

// pointsToImage takes a multidimensional array of points and converts into a png image
func pointsToImage(points [][]int) image.Image {
	identiconLeft := image.Point{0, 0}
	identiconRight := image.Point{width*scale + margin*scale*2, height*scale + margin*scale*2}
	identicon := image.NewRGBA(image.Rectangle{identiconLeft, identiconRight})

	draw.Draw(identicon, identicon.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)

	backgroundLeft := image.Point{margin * scale, margin * scale}
	backgroundRight := image.Point{identiconRight.X - margin*scale, identiconRight.Y - margin*scale}
	background := image.Rectangle{backgroundLeft, backgroundRight}

	draw.Draw(identicon, background, &image.Uniform{fgColor}, image.ZP, draw.Src)

	for i := 0; i < len(points); i++ {
		x, y := points[i][0]%width, points[i][1]%height

		xMargin, yMargin := x*scale+margin*scale, y*scale+margin*scale
		scaledX, scaledY := xMargin+scale, yMargin+scale

		pointBottomLeft, pointUpperRight := image.Point{xMargin, yMargin}, image.Point{scaledX, scaledY}

		draw.Draw(identicon, image.Rectangle{pointBottomLeft, pointUpperRight}, &image.Uniform{bgColor}, image.ZP, draw.Src)
	}

	return identicon
}

// hashToPoints takes hashed array of bytes and converts it to a multidimensional array of integer values representing x and y
func hashToPoints(hash []byte) [][]int {
	var hashPoints [][]int

	oneFourthHashLength := len(hash) / 4
	currentIndex := 0

	for i := 0; i < oneFourthHashLength; i++ {
		var decodedHashPoints []int

		bytesX, bytesY := string(hash[currentIndex:currentIndex+2]), string(hash[currentIndex+2:currentIndex+4])

		decodedHexX, _ := hex.DecodeString(bytesX)
		decodedHexY, _ := hex.DecodeString(bytesY)

		decodedHashPoints = append(decodedHashPoints, int(decodedHexX[0]))
		decodedHashPoints = append(decodedHashPoints, int(decodedHexY[0]))

		hashPoints = append(hashPoints, decodedHashPoints)
		currentIndex += 4
	}

	return hashPoints
}

// FromHash takes hashed array of bytes and converts to a unique identicon
func FromHash(hash []byte) image.Image {
	pointsArray := hashToPoints(hash)
	return pointsToImage(pointsArray)
}
