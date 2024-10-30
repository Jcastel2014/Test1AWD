package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *appDependencies) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.notAllowedResponse)

	//create product
	router.HandlerFunc(http.MethodPost, "/createProduct", a.createProduct)

	//display product
	router.HandlerFunc(http.MethodGet, "/displayProduct/:id", a.displayProduct)

	//update
	// router.HandlerFunc(http.MethodPatch, "/updateProduct/:id", a.updateProduct)

	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	// router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.updateCommentHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	// router.HandlerFunc(http.MethodDelete, "/v1/comments/:id", a.deleteCommentHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/comments", a.listCommentsHandler)
	return a.recoverPanic(router)
}
