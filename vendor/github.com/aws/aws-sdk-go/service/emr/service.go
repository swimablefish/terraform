// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package emr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
	"github.com/aws/aws-sdk-go/private/signer/v4"
)

// Amazon Elastic MapReduce (Amazon EMR) is a web service that makes it easy
// to process large amounts of data efficiently. Amazon EMR uses Hadoop processing
// combined with several AWS products to do tasks such as web indexing, data
// mining, log file analysis, machine learning, scientific simulation, and data
// warehousing.
//The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
type EMR struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// A ServiceName is the name of the service the client will make API calls to.
const ServiceName = "elasticmapreduce"

// New creates a new instance of the EMR client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a EMR client from just a session.
//     svc := emr.New(mySession)
//
//     // Create a EMR client with additional configuration
//     svc := emr.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *EMR {
	c := p.ClientConfig(ServiceName, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion string) *EMR {
	svc := &EMR{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2009-03-31",
				JSONVersion:   "1.1",
				TargetPrefix:  "ElasticMapReduce",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBack(v4.Sign)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a EMR operation and runs any
// custom request initialization.
func (c *EMR) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
