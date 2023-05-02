// Copyright 2023 Nordcloud Oy or its affiliates. All Rights Reserved.

package errors

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	errorKey           = "error"
	errorCtxKey        = "error_context"
	errorStackKey      = "error_stack"
	awsErrorCodeKey    = "aws_error_code"
	awsErrorMessageKey = "aws_error_message"
)

type contextBuilder func(nce *NCError) Fields

func getLogger(err error, builder contextBuilder) *logrus.Entry {
	return logrus.WithFields(buildLogFields(err, builder))
}

// LogWithSeverity uses severity stored in the error to select appropriate log level.
func LogWithSeverity(err error) {
	switch GetErrorSeverity(err) {
	case ERROR:
		LogError(err)
	case WARN:
		LogWarning(err)
	case INFO:
		LogInfo(err)
	case DEBUG:
		LogDebug(err)
	default:
		LogError(err)
	}
}

// LogError logs err with `logrus.Error` method (level=error).
// (uses a newly created `logrus.Entry`)
func LogError(err error) {
	getLogger(err, (*NCError).GetContext).Error(err.Error())
}

// LogWarning logs err at level = warning.
// (uses a newly created `logrus.Entry`)
func LogWarning(err error) {
	getLogger(err, (*NCError).GetContext).Warn(err.Error())
}

// LogInfo logs err at level = info.
// (uses a newly created `logrus.Entry`)
func LogInfo(err error) {
	getLogger(err, (*NCError).GetContext).Info(err.Error())
}

// LogDebug logs err at level = debug.
// (uses a newly created `logrus.Entry`)
func LogDebug(err error) {
	getLogger(err, (*NCError).GetMergedFieldsContext).Debug(err.Error())
}

// LogErrorMerged logs err with `logrus.Error` method (level=error) and merged fields as context.
// (uses a newly created `logrus.Entry`)
func LogErrorMerged(err error) {
	getLogger(err, (*NCError).GetMergedFieldsContext).Error(err.Error())
}

// LogWarningMerged logs err at level = warning. and merged fields as context.
// (uses a newly created `logrus.Entry`)
func LogWarningMerged(err error) {
	getLogger(err, (*NCError).GetMergedFieldsContext).Warn(err.Error())
}

// LogInfoMerged logs err at level = info. and merged fields as context.
// (uses a newly created `logrus.Entry`)
func LogInfoMerged(err error) {
	getLogger(err, (*NCError).GetMergedFieldsContext).Info(err.Error())
}

// LogDebugMerged logs err at level = debug. and merged fields as context.
// (uses a newly created `logrus.Entry`)
func LogDebugMerged(err error) {
	getLogger(err, (*NCError).GetMergedFieldsContext).Debug(err.Error())
}

// GetLogFields converts an error into `logrus.Fields`. It will set an `error` field so you don't have to use the
// `WithError()` method on your own. Additionally it will also fill the output with an error context under the
// `error context` field.
func GetLogFields(err error) logrus.Fields {
	return buildLogFields(err, (*NCError).GetContext)
}

// GetLogFields converts an error into `logrus.Fields`. It will set an `error` field so you don't have to use the
// `WithError()` method on your own. Additionally it will also fill the output with an error context
// with merged causes' fields under the `error context` field.
func GetMergedLogFields(err error) logrus.Fields {
	return buildLogFields(err, (*NCError).GetMergedFieldsContext)
}

// LogErrorPlain logs error with its merged fields and stack at level = Error.
func LogErrorPlain(err error) {
	logrus.WithFields(buildPlainLogFields(err)).Error(err.Error())
}

// LogWarningPlain logs error with its merged fields and stack at level = Warning.
func LogWarningPlain(err error) {
	logrus.WithFields(buildPlainLogFields(err)).Warning(err.Error())
}

// LogInfoPlain logs error with its merged fields and stack at level = Info.
func LogInfoPlain(err error) {
	logrus.WithFields(buildPlainLogFields(err)).Info(err.Error())
}

// LogDebugPlain logs error with its merged fields and stack at level = Debug.
func LogDebugPlain(err error) {
	logrus.WithFields(buildPlainLogFields(err)).Debug(err.Error())
}

func buildLogFields(err error, buildContext contextBuilder) logrus.Fields {
	nativeError := errors.Cause(err)
	if ncError, ok := nativeError.(NCError); ok {
		//rootError is AWS error
		rootError := errors.Cause(ncError.RootError)
		if awsErr, ok := rootError.(awserr.Error); ok {
			return logrus.Fields{
				errorKey:           ncError.Error(),
				errorCtxKey:        buildContext(&ncError),
				awsErrorCodeKey:    awsErr.Code(),
				awsErrorMessageKey: awsErr.Message(),
			}
		}

		return logrus.Fields{
			errorKey:    ncError.Error(),
			errorCtxKey: buildContext(&ncError),
		}
	}
	//error is not NCError but still it is AWS error
	if awsErr, ok := nativeError.(awserr.Error); ok {
		errKey := awsErr.Error()
		if awsErr.OrigErr() != nil {
			errKey = awsErr.OrigErr().Error()
		}
		return logrus.Fields{
			errorKey:           errKey,
			awsErrorCodeKey:    awsErr.Code(),
			awsErrorMessageKey: awsErr.Message(),
		}
	}
	if err != nil {
		return logrus.Fields{errorKey: err.Error()}
	}
	return logrus.Fields{errorKey: nil}
}

func buildPlainLogFields(err error) logrus.Fields {
	nativeError := errors.Cause(err)
	if ncError, ok := nativeError.(NCError); ok {
		logFields := logrus.Fields(ncError.GetMergedFields())
		logFields[errorKey] = ncError.Error()
		logFields[errorStackKey] = ncError.Stack

		//rootError is AWS error
		rootError := errors.Cause(ncError.RootError)
		if awsErr, ok := rootError.(awserr.Error); ok {
			logFields[awsErrorCodeKey] = awsErr.Code()
			logFields[awsErrorMessageKey] = awsErr.Message()
		}

		return logFields
	}
	//error is not NCError but still it is AWS error
	if awsErr, ok := nativeError.(awserr.Error); ok {
		return logrus.Fields{
			errorKey:           awsErr.OrigErr().Error(),
			awsErrorCodeKey:    awsErr.Code(),
			awsErrorMessageKey: awsErr.Message(),
		}
	}
	if err != nil {
		return logrus.Fields{errorKey: err.Error()}
	}
	return logrus.Fields{errorKey: nil}
}
