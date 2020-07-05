package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"image/jpeg"
	"image/png"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"encoding/base64"

	mimetype "github.com/gabriel-vasile/mimetype"
	resizer "github.com/nfnt/resize"
)

type imagResultFormat struct {
	headerType   string
	base64String string
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	imgUrl, urlFound := request.QueryStringParameters["url"]
	imgWidthString, widthFound := request.QueryStringParameters["width"]
	imgHeightString, heightFound := request.QueryStringParameters["height"]
	var response events.APIGatewayProxyResponse
	if urlFound && (widthFound || heightFound) {

		if !widthFound && heightFound {
			//we must have a height but no width
			imgWidthString = "0"
		} else if !heightFound && widthFound {
			//we must have a width but no height
			imgHeightString = "0"
		}

		var errImg, errWidth, errHeight error

		// query parameters can be URL encoded
		imgUrl, errImg = url.QueryUnescape(imgUrl)
		imgWidthString, errWidth = url.QueryUnescape(imgWidthString)
		imgHeightString, errHeight = url.QueryUnescape(imgHeightString)

		if nil != errImg || nil != errWidth || nil != errHeight {
			log.Fatal(errImg, errWidth, errImg)
		} else {
			imgWidth, errCWidth := strconv.Atoi(imgWidthString)
			imgHeight, errCHeight := strconv.Atoi(imgHeightString)
			if nil != errCWidth && nil != errCHeight {
				log.Fatal(errCWidth, errCHeight)
			}
			newImage := ResizeImage(imgUrl, uint(imgWidth), uint(imgHeight))
			response = prepareResponse(newImage)
		}
	} else {
		// Some parameters were not passed in.
		response = prepareResponse(makeError())
	}

	return response, nil
}

func makeError() imagResultFormat {
	buf := bytes.NewBuffer(nil)
	//Exist's in the base as an asset ( ~75k image size )
	f, err := os.Open("error.jpg")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(buf,f)
	f.Close()
	var imageTypeResponse imagResultFormat
	imageTypeResponse.headerType = "image/jpeg"
	imageTypeResponse.base64String = base64.StdEncoding.EncodeToString(buf.Bytes())
	return imageTypeResponse
}

func prepareResponse(imageResult imagResultFormat) events.APIGatewayProxyResponse {
	var Headers = make(map[string]string)
	Headers["content-type"] = imageResult.headerType
	Headers["Access-Control-Allow-Origin"] = "*"
	result := events.APIGatewayProxyResponse{
		Body:            fmt.Sprintf("%v", imageResult.base64String),
		StatusCode:      200,
		Headers:         Headers,
		IsBase64Encoded: true,
	}
	return result
}

func ResizeImage(Path string, Width uint, Height uint) imagResultFormat {

	// Downlaod the file
	res, err := http.Get(Path)
	if err != nil || res.StatusCode != 200 {
		log.Fatal(err)
	}
	defer res.Body.Close()

	//Clone th reader
	var buf bytes.Buffer
	readerClone := io.TeeReader(res.Body, &buf)

	mimeObject, err := mimetype.DetectReader(readerClone)
	mimeType := mimeObject.String()

	/*
		Because mimetype.DetectReader only reads far enough into the reader to get the info
		require we actually need a full-copy of the stream
	 */
	cp, _ := ioutil.ReadAll(readerClone)
	if cp == nil {
		//yay!
	}
	var imageTypeResponse imagResultFormat
	var imageBuf bytes.Buffer
	switch mimeType {
	case "image/jpg", "image/jpeg":
		imageTypeResponse.headerType = "image/jpeg"
		img, err := jpeg.Decode(&buf)
		if err != nil {
			log.Fatal(err)
		}
		newImage := resizer.Resize(Width, Height, img, resizer.Lanczos3)
		err = jpeg.Encode(&imageBuf, newImage, nil)
		if err != nil {
			log.Fatal(err)
		}
	case "image/png":
		imageTypeResponse.headerType = "image/png"
		img, err := png.Decode(&buf)
		if err != nil {
			log.Fatal(err)
		}
		newImage := resizer.Resize(Width, Height, img, resizer.Lanczos3)
		err = png.Encode(&imageBuf, newImage)
		if err != nil {
			log.Fatal(err)
		}
		default:
			log.Fatal(err)
	}

	imageTypeResponse.base64String = base64.StdEncoding.EncodeToString(imageBuf.Bytes())
	return imageTypeResponse
}
