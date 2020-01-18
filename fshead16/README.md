# fshead32

一个固定长度16字节的协议头。

```
// FsHead 协议头
type FsHead struct {
	// 魔法变量 用于校验协议是否匹配    [0-4) bytes
	// 若为0，则使用默认值
	MagicNum uint32

	// 调用方名称,取前8个字节           [4-10) bytes
	ClientName string

	// 后面的元数据长度                [10-12) bytes
	// 消息完整格式为：
	// {FsHead:固定长度}{Meta}{Body}
	MetaLen uint16

	// 消息体的长度                    [12-16) bytes
	BodyLen uint32
}
```