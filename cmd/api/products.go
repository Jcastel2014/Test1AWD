package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jcastel2014/test1/internal/data"
	"github.com/jcastel2014/test1/internal/validator"
)

func (a *appDependencies) createProduct(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Image_url   string  `json:"image_url"`
		Price       float64 `json:"price"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	product := &data.Product{
		Name:        incomingData.Name,
		Description: incomingData.Description,
		Category:    incomingData.Category,
		Image_url:   incomingData.Image_url,
		Price:       incomingData.Price,
	}

	v := validator.New()

	// one sent to identify handler
	data.ValidateProduct(v, product, 1)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.productModel.Insert(product)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/product/%d", product.ID))

	data := envelope{
		"product": product,
	}

	err = a.writeJSON(w, http.StatusCreated, data, headers)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
}

func (a *appDependencies) displayProduct(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	product, err := a.productModel.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrResponse(w, r, err)
		}

		return
	}

	data := envelope{
		"product": product,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}
}

// func (a *appDependencies) updateProduct(w http.ResponseWriter, r *http.Request) {
// 	id, err := a.readIDParam(r)

// 	if err != nil {
// 		a.notFoundResponse(w, r)
// 		return
// 	}

// 	product, err := a.productModel.Get(id)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			a.notFoundResponse(w, r)
// 		default:
// 			a.serverErrResponse(w, r, err)
// 		}

// 		return
// 	}

// 	var incomingData struct {
// 		Name        string `json:"name"`
// 		Description string `json:"description"`
// 		Category    string `json:"category"`
// 		Image_url   string `json:"image_url"`
// 	}

// 	err = a.readJSON(w, r, &incomingData)

// 	if err != nil {
// 		a.badRequestResponse(w, r, err)
// 		return
// 	}

// 	if incomingData.Content != nil {
// 		comment.Content = *incomingData.Content
// 	}

// 	err = a.commentModel.Update(comment)

// 	if err != nil {
// 		a.serverErrResponse(w, r, err)
// 		return
// 	}

// 	data := envelope{
// 		"comment": comment,
// 	}

// 	err = a.writeJSON(w, http.StatusOK, data, nil)
// 	if err != nil {
// 		a.serverErrResponse(w, r, err)
// 		return
// 	}
// }

// func (a *appDependencies) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := a.readIDParam(r)

// 	if err != nil {
// 		a.notFoundResponse(w, r)
// 		return
// 	}

// 	err = a.commentModel.Delete(id)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			a.notFoundResponse(w, r)
// 		default:
// 			a.serverErrResponse(w, r, err)
// 		}

// 		return
// 	}

// 	data := envelope{
// 		"message": "comment successfully deleted",
// 	}

// 	err = a.writeJSON(w, http.StatusOK, data, nil)
// 	if err != nil {
// 		a.serverErrResponse(w, r, err)
// 	}
// }

// func (a *appDependencies) listCommentsHandler(w http.ResponseWriter, r *http.Request) {
// 	var queryParametersData struct {
// 		Content string
// 		Author  string
// 		data.Filters
// 	}

// 	queryParameters := r.URL.Query()
// 	queryParametersData.Content = a.getSingleQueryParameters(queryParameters, "content", "")
// 	queryParametersData.Author = a.getSingleQueryParameters(queryParameters, "author", "")

// 	v := validator.New()

// 	queryParametersData.Filters.Page = a.getSingleIntegerParameters(queryParameters, "page", 1, v)
// 	queryParametersData.Filters.PageSize = a.getSingleIntegerParameters(queryParameters, "page_size", 10, v)

// 	data.ValidateFilters(v, queryParametersData.Filters)
// 	if !v.IsEmpty() {
// 		a.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	comments, err := a.commentModel.GetAll(queryParametersData.Content, queryParametersData.Author, queryParametersData.Filters)

// 	if err != nil {
// 		a.serverErrResponse(w, r, err)
// 		return
// 	}

// 	data := envelope{
// 		"comments": comments,
// 	}

// 	err = a.writeJSON(w, http.StatusOK, data, nil)

// 	if err != nil {
// 		a.serverErrResponse(w, r, err)
// 	}
// }
