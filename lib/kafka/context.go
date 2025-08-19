package kafka

import (
	"context"
	"time"
)

type (
	topicCtx     struct{}
	partitionCtx struct{}
	offsetCtx    struct{}
	timestampCtx struct{}
)

func TopicFromContext(ctx context.Context) string {
	v, _ := ctx.Value(topicCtx{}).(string)
	return v
}

func TopicToContext(ctx context.Context, topic string) context.Context {
	return context.WithValue(ctx, topicCtx{}, topic)
}

func PartitionFromContext(ctx context.Context) int32 {
	v, _ := ctx.Value(partitionCtx{}).(int32)
	return v
}

func PartitionToContext(ctx context.Context, partition int32) context.Context {
	return context.WithValue(ctx, partitionCtx{}, partition)
}

func OffsetFromContext(ctx context.Context) int64 {
	v, _ := ctx.Value(offsetCtx{}).(int64)
	return v
}

func OffsetToContext(ctx context.Context, offset int64) context.Context {
	return context.WithValue(ctx, offsetCtx{}, offset)
}

func TimestampToContext(ctx context.Context, timestamp time.Time) context.Context {
	return context.WithValue(ctx, timestampCtx{}, timestamp)
}

func TimestampFromContext(ctx context.Context) time.Time {
	v, _ := ctx.Value(timestampCtx{}).(time.Time)
	return v
}
