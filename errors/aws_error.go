package errors

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

const (
	AWSAccessDenied          = "AccessDenied"
	AWSAccessDeniedException = "AccessDeniedException"

	AWSDynamoTableNotFound = dynamodb.ErrCodeTableNotFoundException
	AWSS3BucketNotFound    = s3.ErrCodeNoSuchBucket

	AWSRedshiftClusterSnapsotQuotaExceeded = redshift.ErrCodeClusterSnapshotQuotaExceededFault
)

//GetAWSErrorCode returns the underlying AWS error code from the error.
func GetAWSErrorCode(err error) string {
	nativeError := errors.Cause(err)
	if nativeError == nil {
		return ""
	}
	if ncError, ok := nativeError.(NCError); ok {
		nativeError = errors.Cause(ncError.RootError)
	}
	if nativeError == nil {
		return ""
	}

	if awsError, ok := nativeError.(awserr.Error); ok {
		return awsError.Code()
	}

	return ""
}
