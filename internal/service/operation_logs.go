package service

import "context"

type OperationLogInput struct {
	StaffId       int64
	StaffName     string
	Operation     string
	Resource      string
	ResourceId    string
	Detail        string
	Result        int
	FailureReason string
}

type IOperationLogs interface {
	Log(ctx context.Context, in *OperationLogInput) error
}

var localOperationLogs IOperationLogs

func OperationLogs() IOperationLogs {
	if localOperationLogs == nil {
		panic("implement not found for interface IOperationLogs, forgot register?")
	}
	return localOperationLogs
}

func RegisterOperationLogs(i IOperationLogs) {
	localOperationLogs = i
}
