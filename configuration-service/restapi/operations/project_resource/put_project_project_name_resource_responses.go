// Code generated by go-swagger; DO NOT EDIT.

package project_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/configuration-service/models"
)

// PutProjectProjectNameResourceCreatedCode is the HTTP code returned for type PutProjectProjectNameResourceCreated
const PutProjectProjectNameResourceCreatedCode int = 201

/*PutProjectProjectNameResourceCreated Success. Project resources have been updated. The version of the new configuration is returned.

swagger:response putProjectProjectNameResourceCreated
*/
type PutProjectProjectNameResourceCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Version `json:"body,omitempty"`
}

// NewPutProjectProjectNameResourceCreated creates PutProjectProjectNameResourceCreated with default headers values
func NewPutProjectProjectNameResourceCreated() *PutProjectProjectNameResourceCreated {

	return &PutProjectProjectNameResourceCreated{}
}

// WithPayload adds the payload to the put project project name resource created response
func (o *PutProjectProjectNameResourceCreated) WithPayload(payload *models.Version) *PutProjectProjectNameResourceCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name resource created response
func (o *PutProjectProjectNameResourceCreated) SetPayload(payload *models.Version) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameResourceCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutProjectProjectNameResourceBadRequestCode is the HTTP code returned for type PutProjectProjectNameResourceBadRequest
const PutProjectProjectNameResourceBadRequestCode int = 400

/*PutProjectProjectNameResourceBadRequest Failed. Project resources could not be updated.

swagger:response putProjectProjectNameResourceBadRequest
*/
type PutProjectProjectNameResourceBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutProjectProjectNameResourceBadRequest creates PutProjectProjectNameResourceBadRequest with default headers values
func NewPutProjectProjectNameResourceBadRequest() *PutProjectProjectNameResourceBadRequest {

	return &PutProjectProjectNameResourceBadRequest{}
}

// WithPayload adds the payload to the put project project name resource bad request response
func (o *PutProjectProjectNameResourceBadRequest) WithPayload(payload *models.Error) *PutProjectProjectNameResourceBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name resource bad request response
func (o *PutProjectProjectNameResourceBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameResourceBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PutProjectProjectNameResourceDefault Error

swagger:response putProjectProjectNameResourceDefault
*/
type PutProjectProjectNameResourceDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutProjectProjectNameResourceDefault creates PutProjectProjectNameResourceDefault with default headers values
func NewPutProjectProjectNameResourceDefault(code int) *PutProjectProjectNameResourceDefault {
	if code <= 0 {
		code = 500
	}

	return &PutProjectProjectNameResourceDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put project project name resource default response
func (o *PutProjectProjectNameResourceDefault) WithStatusCode(code int) *PutProjectProjectNameResourceDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put project project name resource default response
func (o *PutProjectProjectNameResourceDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put project project name resource default response
func (o *PutProjectProjectNameResourceDefault) WithPayload(payload *models.Error) *PutProjectProjectNameResourceDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name resource default response
func (o *PutProjectProjectNameResourceDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameResourceDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}