package db

import "context"

type Database interface {
	Close(ctx context.Context) error
	Insert(ctx context.Context, object any) error
}