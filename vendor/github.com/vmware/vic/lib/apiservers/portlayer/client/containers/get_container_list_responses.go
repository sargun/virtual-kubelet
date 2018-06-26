package containers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
)

// GetContainerListReader is a Reader for the GetContainerList structure.
type GetContainerListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetContainerListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetContainerListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 500:
		result := NewGetContainerListInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetContainerListOK creates a GetContainerListOK with default headers values
func NewGetContainerListOK() *GetContainerListOK {
	return &GetContainerListOK{}
}

/*GetContainerListOK handles this case with default header values.

OK
*/
type GetContainerListOK struct {
	Payload []*models.ContainerInfo
}

func (o *GetContainerListOK) Error() string {
	return fmt.Sprintf("[GET /containers/list][%d] getContainerListOK  %+v", 200, o.Payload)
}

func (o *GetContainerListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetContainerListInternalServerError creates a GetContainerListInternalServerError with default headers values
func NewGetContainerListInternalServerError() *GetContainerListInternalServerError {
	return &GetContainerListInternalServerError{}
}

/*GetContainerListInternalServerError handles this case with default header values.

server error
*/
type GetContainerListInternalServerError struct {
	Payload *models.Error
}

func (o *GetContainerListInternalServerError) Error() string {
	return fmt.Sprintf("[GET /containers/list][%d] getContainerListInternalServerError  %+v", 500, o.Payload)
}

func (o *GetContainerListInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}