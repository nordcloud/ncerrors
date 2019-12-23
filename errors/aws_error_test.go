package errors

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/stretchr/testify/assert"
)

func TestGetAWSErrorCode_Standard(t *testing.T) {
	cause := errors.New("std except")
	errCode := GetAWSErrorCode(cause)
	assert.Equal(t, "", errCode)
}

func TestGetAWSErrorCode_AWSErr(t *testing.T) {
	cause := awserr.New("code", "aws error", nil)
	errCode := GetAWSErrorCode(cause)
	assert.Equal(t, "code", errCode)
}

func TestGetAWSErrorCode_Context(t *testing.T) {
	cause := WithContext(awserr.New("code1", "aws error", nil), "context1", nil)
	errCode := GetAWSErrorCode(cause)
	assert.Equal(t, "code1", errCode)
}

func TestGetAWSErrorCode_Wrap(t *testing.T) {
	cause := WithContext(errors.Wrap(awserr.New("code1", "aws error", nil), "wrap"), "context1", nil)
	errCode := GetAWSErrorCode(cause)
	assert.Equal(t, "code1", errCode)
}

func TestGetAWSErrorCode_MultilevelContext(t *testing.T) {
	cause := WithContext(WithContext(errors.Wrap(awserr.New("code1", "aws error", nil), "wrap"), "context0", nil), "context1", nil)
	errCode := GetAWSErrorCode(cause)
	assert.Equal(t, "code1", errCode)
}
