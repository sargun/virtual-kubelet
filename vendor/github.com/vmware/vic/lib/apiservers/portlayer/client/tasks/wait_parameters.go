package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
)

// NewWaitParams creates a new WaitParams object
// with the default values initialized.
func NewWaitParams() *WaitParams {
	var ()
	return &WaitParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewWaitParamsWithTimeout creates a new WaitParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewWaitParamsWithTimeout(timeout time.Duration) *WaitParams {
	var ()
	return &WaitParams{

		timeout: timeout,
	}
}

// NewWaitParamsWithContext creates a new WaitParams object
// with the default values initialized, and the ability to set a context for a request
func NewWaitParamsWithContext(ctx context.Context) *WaitParams {
	var ()
	return &WaitParams{

		Context: ctx,
	}
}

// NewWaitParamsWithHTTPClient creates a new WaitParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewWaitParamsWithHTTPClient(client *http.Client) *WaitParams {
	var ()
	return &WaitParams{
		HTTPClient: client,
	}
}

/*WaitParams contains all the parameters to send to the API endpoint
for the wait operation typically these are written to a http.Request
*/
type WaitParams struct {

	/*Config*/
	Config *models.TaskWaitConfig

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the wait params
func (o *WaitParams) WithTimeout(timeout time.Duration) *WaitParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the wait params
func (o *WaitParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the wait params
func (o *WaitParams) WithContext(ctx context.Context) *WaitParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the wait params
func (o *WaitParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the wait params
func (o *WaitParams) WithHTTPClient(client *http.Client) *WaitParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the wait params
func (o *WaitParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithConfig adds the config to the wait params
func (o *WaitParams) WithConfig(config *models.TaskWaitConfig) *WaitParams {
	o.SetConfig(config)
	return o
}

// SetConfig adds the config to the wait params
func (o *WaitParams) SetConfig(config *models.TaskWaitConfig) {
	o.Config = config
}

// WriteToRequest writes these params to a swagger request
func (o *WaitParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	r.SetTimeout(o.timeout)
	var res []error

	if o.Config == nil {
		o.Config = new(models.TaskWaitConfig)
	}

	if err := r.SetBodyParam(o.Config); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}