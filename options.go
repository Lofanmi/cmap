package cmap

// Options 创建Map的配置选项
type Options struct {
	ShardCount uint32
	Serializer *SerializerFunc
}

// Option 配置选项函数
type Option func(*Options)

// WithShardCount 设置分片数量
func WithShardCount(count uint32) Option {
	return func(o *Options) {
		o.ShardCount = count
	}
}

// WithSerializer 设置序列化格式
func WithSerializer(v *SerializerFunc) Option {
	return func(o *Options) {
		o.Serializer = v
	}
}
