# fshead

一个固定长度32字节的协议头。

```
// FsHead 协议头
type FsHead struct {

	// 协议版本
	Version uint16

	// 调用方名称,前8个字节
	ClientName string

	// 调用方ID，若不需要，可以传0
	// server端也可以依次做身份校验
	UserID uint32

	// 日志ID
	LogID uint32

	// 预留字段，业务可以扩展使用
	Reserve uint32

	// 后面的元数据长度
	// 消息完整格式为：{FsHead:固定长度}{Meta}{Body}
	MetaLen uint16

	// 消息体的长度
	BodyLen uint32

	// 魔法变量 用于校验协议是否匹配
	// 若为0，则使用默认值
	MagicNum uint32
}
```