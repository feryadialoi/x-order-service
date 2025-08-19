package kafka

import "github.com/IBM/sarama"

type Header map[string]string

func (Header) FromSarama(headers []*sarama.RecordHeader) Header {
	header := make(Header, len(headers))
	for _, h := range headers {
		header[string(h.Key)] = string(h.Value)
	}
	return header
}

func (h Header) Get(key string) string {
	return h[key]
}

func (h Header) Set(key string, value string) {
	h[key] = value
}

func (h Header) Del(key string) {
	delete(h, key)
}

func (h Header) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

func (h Header) Values() []string {
	values := make([]string, 0, len(h))
	for _, v := range h {
		values = append(values, v)
	}
	return values
}
