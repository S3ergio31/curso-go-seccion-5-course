package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/S3ergio31/curso-go-seccion-5-course/internal/course"
	"github.com/S3ergio31/curso-go-seccion-5-response/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewCourseHttpServer(endpoints course.Endpoints) http.Handler {
	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse,
		encodeResponse,
		opts...,
	)).Methods("POST")

	router.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	return router
}

func decodeCreateCourse(_ context.Context, r *http.Request) (any, error) {
	var request course.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return request, nil
}

func decodeUpdateCourse(_ context.Context, r *http.Request) (any, error) {
	var request course.UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	path := mux.Vars(r)

	request.ID = path["id"]

	return request, nil
}

func decodeGetCourse(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	req := course.GetRequest{ID: p["id"]}

	return req, nil
}

func decodeDeleteCourse(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	req := course.DeleteRequest{ID: p["id"]}

	return req, nil
}

func decodeGetAllCourse(_ context.Context, r *http.Request) (any, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := course.GetAllRequest{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, res any) error {
	r := res.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
