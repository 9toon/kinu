package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/9toon/kinu/logger"
	"github.com/9toon/kinu/resource"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(0)

	imageType := r.FormValue("name")
	if len(imageType) == 0 {
		RespondBadRequest(w, "required name parameter")
		return
	}

	imageId := r.FormValue("id")
	if len(imageId) == 0 {
		RespondBadRequest(w, "required id parameter")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		RespondBadRequest(w, "invalid file")
		return
	}

	err = resource.New(imageType, imageId).Store(file)
	if err != nil {
		if _, ok := err.(*ErrInvalidRequest); ok {
			RespondBadRequest(w, err.Error())
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}

	RespondImageUploadSuccessJson(w, imageType, imageId)

	logger.WithFields(logrus.Fields{
		"path":   r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
	}).Info("success")
}
