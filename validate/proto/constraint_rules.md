## 约束规则

[提供的约束](https://github.com/envoyproxy/protoc-gen-validate/blob/main/validate/validate.proto)在 JSON Schema 中更大程度地建模。同一字段可以混合使用PGV规则；该插件确保应用于字段的规则在代码生成之前不会矛盾。

检查[约束规则比较矩阵](https://github.com/envoyproxy/protoc-gen-validate/blob/main/rule_comparison.md)以了解特定于语言的约束功能。

### [](https://github.com/envoyproxy/protoc-gen-validate#numerics)数字

> 所有数字类型（`float`, `double`, `int32`, `int64`, `uint32`, `uint64`, `sint32`, `sint64`, `fixed32`, `fixed64`, `sfixed32`, `sfixed64`）共享相同的规则。

-   **const**：该字段必须_完全_是指定的值。
    
```protobuf
    // x must equal 1.23 exactly
    float x = 1 [(validate.rules).float.const = 1.23];
```
    
-   **lt/lte/gt/gte**：这些不等式（分别为`<`、`<=`、`>`、`>=`）允许导出字段必须驻留的范围。
    
```protobuf
    // x must be less than 10
    int32 x = 1 [(validate.rules).int32.lt = 10];
    
    // x must be greater than or equal to 20
    uint64 x = 1 [(validate.rules).uint64.gte = 20];
    
    // x must be in the range [30, 40)
    fixed32 x = 1 [(validate.rules).fixed32 = {gte:30, lt: 40}];
    
    反转 和 的值`lt(e)`是`gt(e)`有效的，并创建一个独占范围。
    
    // x must be outside the range [30, 40)
    double x = 1 [(validate.rules).double = {lt:30, gte:40}];
    
-   **in/not_in**：这两个规则允许为字段的值指定允许/拒绝列表。
    
    // x must be either 1, 2, or 3
    uint32 x = 1 [(validate.rules).uint32 = {in: [1,2,3]}];
    
    // x cannot be 0 nor 0.99
    float x = 1 [(validate.rules).float = {not_in: [0, 0.99]}];
```
    
-   **ignore_empty**：此规则指定如果字段为空或设置为默认值，则忽略任何验证规则。如果能够在更新请求中取消设置字段，或者在切换到 WKT 不可行的情况下跳过可选字段的验证，这些通常很有用。
    
```protobuf
    unint32 x = 1 [(validate.rules).uint32 = {ignore_empty: true, gte: 200}];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#bools)布尔值

-   **const**：该字段必须_完全_是指定的值。
    
```protobuf
    // x must be set to true
    bool x = 1 [(validate.rules).bool.const = true];
    
    // x cannot be set to true
    bool x = 1 [(validate.rules).bool.const = false];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#strings)字符串

-   **const**：该字段必须_完全_是指定的值。
    
```protobuf
    // x must be set to "foo"
    string x = 1 [(validate.rules).string.const = "foo"];
```
    
-   **len/min_len/max_len**：这些规则限制字段中的字符数（Unicode 代码点）。请注意，字符数可能与字符串中的字节数不同。该字符串被认为是原样，并且不规范化。
    
```protobuf
    // x must be exactly 5 characters long
    string x = 1 [(validate.rules).string.len = 5];
    
    // x must be at least 3 characters long
    string x = 1 [(validate.rules).string.min_len = 3];
    
    // x must be between 5 and 10 characters, inclusive
    string x = 1 [(validate.rules).string = {min_len: 5, max_len: 10}];
    
-   **min_bytes/max_bytes**：这些规则限制了字段中的字节数。
    
    // x must be at most 15 bytes long
    string x = 1 [(validate.rules).string.max_bytes = 15];
    
    // x must be between 128 and 1024 bytes long
    string x = 1 [(validate.rules).string = {min_bytes: 128, max_bytes: 1024}];
```
    
-   **pattern**：该字段必须匹配指定[的符合 RE2 的](https://github.com/google/re2/wiki/Syntax)正则表达式。包含的表达式应该省略任何分隔符（即，`/\d+/`应该只是`\d+`）。
    
```protobuf
    // x must be a non-empty, case-insensitive hexadecimal string
    string x = 1 [(validate.rules).string.pattern = "(?i)^[0-9a-f]+$"];
```
    
-   **prefix/suffix/contains/not_contains**：该字段必须在可选的显式位置包含指定的子字符串，或者不包含指定的子字符串。
    
```protobuf
    // x must begin with "foo"
    string x = 1 [(validate.rules).string.prefix = "foo"];
    
    // x must end with "bar"
    string x = 1 [(validate.rules).string.suffix = "bar"];
    
    // x must contain "baz" anywhere inside it
    string x = 1 [(validate.rules).string.contains = "baz"];
    
    // x cannot contain "baz" anywhere inside it
    string x = 1 [(validate.rules).string.not_contains = "baz"];
    
    // x must begin with "fizz" and end with "buzz"
    string x = 1 [(validate.rules).string = {prefix: "fizz", suffix: "buzz"}];
    
    // x must end with ".proto" and be less than 64 characters
    string x = 1 [(validate.rules).string = {suffix: ".proto", max_len:64}];
    
-   **in/not_in**：这两个规则允许为字段的值指定允许/拒绝列表。
    
    // x must be either "foo", "bar", or "baz"
    string x = 1 [(validate.rules).string = {in: ["foo", "bar", "baz"]}];
    
    // x cannot be "fizz" nor "buzz"
    string x = 1 [(validate.rules).string = {not_in: ["fizz", "buzz"]}];
```
    
-   **ignore_empty**：此规则指定如果字段为空或设置为默认值，则忽略任何验证规则。如果能够在更新请求中取消设置字段，或者在切换到 WKT 不可行的情况下跳过可选字段的验证，这些通常很有用。
    
```protobuf
    string CountryCode = 1 [(validate.rules).string = {ignore_empty: true, len: 2}];
```
    
-   **众所周知的格式**：这些规则为常见的字符串模式提供了高级约束。这些约束通常比等效的正则表达式模式更宽松、更高效，同时提供更多解释性的故障描述。
    
```protobuf
    // x must be a valid email address (via RFC 1034)
    string x = 1 [(validate.rules).string.email = true];
    
    // x must be a valid address (IP or Hostname).
    string x = 1 [(validate.rules).string.address = true];
    
    // x must be a valid hostname (via RFC 1034)
    string x = 1 [(validate.rules).string.hostname = true];
    
    // x must be a valid IP address (either v4 or v6)
    string x = 1 [(validate.rules).string.ip = true];
    
    // x must be a valid IPv4 address
    // eg: "192.168.0.1"
    string x = 1 [(validate.rules).string.ipv4 = true];
    
    // x must be a valid IPv6 address
    // eg: "fe80::3"
    string x = 1 [(validate.rules).string.ipv6 = true];
    
    // x must be a valid absolute URI (via RFC 3986)
    string x = 1 [(validate.rules).string.uri = true];
    
    // x must be a valid URI reference (either absolute or relative)
    string x = 1 [(validate.rules).string.uri_ref = true];
    
    // x must be a valid UUID (via RFC 4122)
    string x = 1 [(validate.rules).string.uuid = true];
    
    // x must conform to a well known regex for HTTP header names (via RFC 7230)
    string x = 1 [(validate.rules).string.well_known_regex = HTTP_HEADER_NAME]
    
    // x must conform to a well known regex for HTTP header values (via RFC 7230) 
    string x = 1 [(validate.rules).string.well_known_regex = HTTP_HEADER_VALUE];
    
    // x must conform to a well known regex for headers, disallowing \r\n\0 characters.
    string x = 1 [(validate.rules).string {well_known_regex: HTTP_HEADER_VALUE, strict: false}];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#bytes)字节

> 文字值应该用字符串表示，必要时使用转义。

-   **const**：该字段必须_完全_是指定的值。
    
```protobuf
    // x must be set to "foo" ("\x66\x6f\x6f")
    bytes x = 1 [(validate.rules).bytes.const = "foo"];
    
    // x must be set to "\xf0\x90\x28\xbc"
    bytes x = 1 [(validate.rules).bytes.const = "\xf0\x90\x28\xbc"];
```
    
-   **len/min_len/max_len**：这些规则限制字段中的字节数。
    
```protobuf
    // x must be exactly 3 bytes
    bytes x = 1 [(validate.rules).bytes.len = 3];
    
    // x must be at least 3 bytes long
    bytes x = 1 [(validate.rules).bytes.min_len = 3];
    
    // x must be between 5 and 10 bytes, inclusive
    bytes x = 1 [(validate.rules).bytes = {min_len: 5, max_len: 10}];
```
    
-   **pattern**：该字段必须匹配指定[的符合 RE2 的](https://github.com/google/re2/wiki/Syntax)正则表达式。包含的表达式应该省略任何分隔符（即，`/\d+/`应该只是`\d+`）。
    
```protobuf
    // x must be a non-empty, ASCII byte sequence
    bytes x = 1 [(validate.rules).bytes.pattern = "^[\x00-\x7F]+$"];
```
    
-   **prefix/suffix/contains**：该字段必须在可选的显式位置包含指定的字节序列。
    
```protobuf
    // x must begin with "\x99"
    bytes x = 1 [(validate.rules).bytes.prefix = "\x99"];
    
    // x must end with "buz\x7a"
    bytes x = 1 [(validate.rules).bytes.suffix = "buz\x7a"];
    
    // x must contain "baz" anywhere inside it
    bytes x = 1 [(validate.rules).bytes.contains = "baz"];
    
-   **in/not_in**：这两个规则允许为字段的值指定允许/拒绝列表。
    
    // x must be either "foo", "bar", or "baz"
    bytes x = 1 [(validate.rules).bytes = {in: ["foo", "bar", "baz"]}];
    
    // x cannot be "fizz" nor "buzz"
    bytes x = 1 [(validate.rules).bytes = {not_in: ["fizz", "buzz"]}];
```
    
-   **ignore_empty**：此规则指定如果字段为空或设置为默认值，则忽略任何验证规则。如果能够在更新请求中取消设置字段，或者在切换到 WKT 不可行的情况下跳过可选字段的验证，这些通常很有用。
    
```protobuf
    bytes x = 1 [(validate.rules).bytes = {ignore_empty: true, in: ["foo", "bar", "baz"]}];
```
    
-   **众所周知的格式**：这些规则为常见模式提供了高级约束。这些约束通常比等效的正则表达式模式更宽松、更高效，同时提供更多解释性的故障描述。
    
```protobuf
    // x must be a valid IP address (either v4 or v6) in byte format
    bytes x = 1 [(validate.rules).bytes.ip = true];
    
    // x must be a valid IPv4 address in byte format
    // eg: "\xC0\xA8\x00\x01"
    bytes x = 1 [(validate.rules).bytes.ipv4 = true];
    
    // x must be a valid IPv6 address in byte format
    // eg: "\x20\x01\x0D\xB8\x85\xA3\x00\x00\x00\x00\x8A\x2E\x03\x70\x73\x34"
    bytes x = 1 [(validate.rules).bytes.ipv6 = true];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#enums)枚举

> 所有文字值都应使用枚举描述符中定义的数字 (int32) 值。

以下示例使用此`State`枚举

```protobuf
enum State {
  INACTIVE = 0;
  PENDING  = 1;
  ACTIVE   = 2;
}
```

-   **const**：该字段必须_完全_是指定的值。
    
```protobuf
    // x must be set to ACTIVE (2)
    State x = 1 [(validate.rules).enum.const = 2];
```
    
-   **defined_only**：该字段必须是枚举描述符中的指定值之一。
    
```protobuf
    // x can only be INACTIVE, PENDING, or ACTIVE
    State x = 1 [(validate.rules).enum.defined_only = true];
```
    
-   **in/not_in**：这两个规则允许为字段的值指定允许/拒绝列表。
    
```protobuf
    // x must be either INACTIVE (0) or ACTIVE (2)
    State x = 1 [(validate.rules).enum = {in: [0,2]}];
    
    // x cannot be PENDING (1)
    State x = 1 [(validate.rules).enum = {not_in: [1]}];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#messages)留言

> 如果一个字段包含一条消息并且该消息是使用 PGV 生成的，则将递归执行验证。未使用 PGV 生成的消息将被跳过。

```protobuf
// if Person was generated with PGV and x is set,
// x's fields will be validated.
Person x = 1;
```

-   **skip**：此规则指定不应评估此字段的验证规则。
    
```protobuf
    // The fields on Person x will not be validated.
    Person x = 1 [(validate.rules).message.skip = true];
```
    
-   **required**：此规则指定该字段不能取消设置。
    
```protobuf
    // x cannot be unset
    Person x = 1 [(validate.rules).message.required = true];
    
    // x cannot be unset, but the validations on x will not be performed
    Person x = 1 [(validate.rules).message = {required: true, skip: true}];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#repeated)重复

-   **min_items/max_items**：这些规则控制字段中包含多少元素
    
```protobuf
    // x must contain at least 3 elements
    repeated int32 x = 1 [(validate.rules).repeated.min_items = 3];
    
    // x must contain between 5 and 10 Persons, inclusive
    repeated Person x = 1 [(validate.rules).repeated = {min_items: 5, max_items: 10}];
    
    // x must contain exactly 7 elements
    repeated double x = 1 [(validate.rules).repeated = {min_items: 7, max_items: 7}];
```
    
-   **unique**：此规则要求字段中的所有元素都必须是唯一的。此规则不支持重复消息。
    
```protobuf
    // x must contain unique int64 values
    repeated int64 x = 1 [(validate.rules).repeated.unique = true];
```
    
-   **items**：此规则指定应应用于字段中每个元素的约束。`skip`除非在此约束中指定，否则重复消息字段也应用其验证规则。
    
```protobuf
    // x must contain positive float values
    repeated float x = 1 [(validate.rules).repeated.items.float.gt = 0];
    
    // x must contain Persons but don't validate them
    repeated Person x = 1 [(validate.rules).repeated.items.message.skip = true];
```
    
-   **ignore_empty**：此规则指定如果字段为空或设置为默认值，则忽略任何验证规则。如果能够在更新请求中取消设置字段，或者在切换到 WKT 不可行的情况下跳过可选字段的验证，这些通常很有用。
    
```protobuf
    repeated int64 x = 1 [(validate.rules).repeated = {ignore_empty: true, items: {int64: {gt: 200}}}];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#maps)地图

-   **min_pairs/max_pairs**：这些规则控制该字段中包含多少 KV 对
    
```protobuf
    // x must contain at most 3 KV pairs
    map<string, uint64> x = 1 [(validate.rules).map.min_pairs = 3];
    
    // x must contain between 5 and 10 KV pairs
    map<string, string> x = 1 [(validate.rules)].map = {min_pairs: 5, max_pairs: 10}];
    
    // x must contain exactly 7 KV pairs
    map<string, Person> x = 1 [(validate.rules)].map = {min_pairs: 7, max_pairs: 7}];
```
    
-   **no_sparse**：对于具有消息值的映射字段，将此规则设置为 true 会禁止具有未设置值的键。
    
```protobuf
    // all values in x must be set
    map<uint64, Person> x = 1 [(validate.rules).map.no_sparse = true];
```
    
-   **keys**：此规则指定应用于字段中键的约束。
    
```protobuf
    // x's keys must all be negative
    <sint32, string> x = [(validate.rules).map.keys.sint32.lt = 0];
```
    
-   **values**：此规则指定应用于字段中每个值的约束。`skip`除非在此约束中指定，否则重复消息字段也应用其验证规则。
    
```protobuf
    // x must contain strings of at least 3 characters
    map<string, string> x = 1 [(validate.rules).map.values.string.min_len = 3];
    
    // x must contain Persons but doesn't validate them
    map<string, Person> x = 1 [(validate.rules).map.values.message.skip = true];
```
    
-   **ignore_empty**：此规则指定如果字段为空或设置为默认值，则忽略任何验证规则。如果能够在更新请求中取消设置字段，或者在切换到 WKT 不可行的情况下跳过可选字段的验证，这些通常很有用。
    
```protobuf
    map<string, string> x = 1 [(validate.rules).map = {ignore_empty: true, values: {string: {min_len: 3}}}];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#well-known-types-wkts)知名类型 (WKT)

一组 WKT 与在许多领域有用的[protoc](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf)和通用消息模式打包在一起。

#### [](https://github.com/envoyproxy/protoc-gen-validate#scalar-value-wrappers)标量值包装器

在`proto3`语法中，没有办法区分未设置和标量字段的零值。值 WKT 通过将它们包装在消息中来允许这种区分。PGV 允许使用包装器封装的相同标量规则。

```protobuf
// if it is set, x must be greater than 3
google.protobuf.Int32Value x = 1 [(validate.rules).int32.gt = 3];

消息规则也可以与标量知名类型 (WKT) 一起使用：

// Ensures that if a value is not set for age, it would not pass the validation despite its zero value being 0.
message X { google.protobuf.Int32Value age = 1 [(validate.rules).int32.gt = -1, (validate.rules).message.required = true]; }
```

#### [](https://github.com/envoyproxy/protoc-gen-validate#anys)任何

```protobuf
-   **required**：此规则指定必须设置该字段
    
    // x cannot be unset
    google.protobuf.Any x = 1 [(validate.rules).any.required = true];
    
-   **in/not_in**：这两个规则允许`type_url`为此字段中的值指定允许/拒绝列表。如果可能，请考虑使用`oneof`联合而不是。`in`
    
    // x must not be the Duration or Timestamp WKT
    google.protobuf.Any x = 1 [(validate.rules).any = {not_in: [
        "type.googleapis.com/google.protobuf.Duration",
        "type.googleapis.com/google.protobuf.Timestamp"
      ]}];
```
    

#### [](https://github.com/envoyproxy/protoc-gen-validate#durations)持续时间

```protobuf
-   **required**：此规则指定必须设置该字段
    
    // x cannot be unset
    google.protobuf.Duration x = 1 [(validate.rules).duration.required = true];
    
-   **const**：该字段必须_完全_是指定的值。
    
    // x must equal 1.5s exactly
    google.protobuf.Duration x = 1 [(validate.rules).duration.const = {
        seconds: 1,
        nanos:   500000000
      }];
```
    
-   **lt/lte/gt/gte**：这些不等式（分别为`<`、`<=`、`>`、`>=`）允许导出字段必须驻留的范围。
    
```protobuf
    // x must be less than 10s
    google.protobuf.Duration x = 1 [(validate.rules).duration.lt.seconds = 10];
    
    // x must be greater than or equal to 20ns
    google.protobuf.Duration x = 1 [(validate.rules).duration.gte.nanos = 20];
    
    // x must be in the range [0s, 1s)
    google.protobuf.Duration x = 1 [(validate.rules).duration = {
        gte: {},
        lt:  {seconds: 1}
      }];
    
    反转 和 的值`lt(e)`是`gt(e)`有效的，并创建一个独占范围。
    
    // x must be outside the range [0s, 1s)
    google.protobuf.Duration x = 1 [(validate.rules).duration = {
        lt:  {},
        gte: {seconds: 1}
      }];
```
    
-   **in/not_in**：这两个规则允许为字段的值指定允许/拒绝列表。
    
```protobuf
    // x must be either 0s or 1s
    google.protobuf.Duration x = 1 [(validate.rules).duration = {in: [
        {},
        {seconds: 1}
      ]}];
    
    // x cannot be 20s nor 500ns
    google.protobuf.Duration x = 1 [(validate.rules).duration = {not_in: [
        {seconds: 20},
        {nanos: 500}
      ]}];
```
    

#### [](https://github.com/envoyproxy/protoc-gen-validate#timestamps)时间戳

-   **required**：此规则指定必须设置该字段
    
```protobuf
    // x cannot be unset
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp.required = true];
```
    
-   **const**：该字段必须_完全_是指定的值。
    
```protobuf
    // x must equal 2009/11/10T23:00:00.500Z exactly
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp.const = {
        seconds: 63393490800,
        nanos:   500000000
      }];
```
    
-   **lt/lte/gt/gte**：这些不等式（分别为`<`、`<=`、`>`、`>=`）允许导出字段必须驻留的范围。
    
```protobuf
    // x must be less than the Unix Epoch
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp.lt.seconds = 0];
    
    // x must be greater than or equal to 2009/11/10T23:00:00Z
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp.gte.seconds = 63393490800];
    
    // x must be in the range [epoch, 2009/11/10T23:00:00Z)
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp = {
        gte: {},
        lt:  {seconds: 63393490800}
      }];
    
    反转 和 的值`lt(e)`是`gt(e)`有效的，并创建一个独占范围。
    
    // x must be outside the range [epoch, 2009/11/10T23:00:00Z)
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp = {
        lt:  {},
        gte: {seconds: 63393490800}
      }];
```
    
-   **lt_now/gt_now**：这些不等式允许相对于当前时间的范围。这些规则不能与上述绝对规则一起使用。
    
```protobuf
    // x must be less than the current timestamp
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp.lt_now = true];
```
    
-   **within**：此规则指定字段的值应在当前时间的持续时间内。此规则可以与这些范围结合使用`lt_now`并`gt_now`控制这些范围。
    
```protobuf
    // x must be within ±1s of the current time
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp.within.seconds = 1];
    
    // x must be within the range (now, now+1h)
    google.protobuf.Timestamp x = 1 [(validate.rules).timestamp = {
        gt_now: true,
        within: {seconds: 3600}
      }];
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#message-global)消息-全球

-   **disabled**：消息上的字段的所有验证规则都可以被取消，包括任何支持自身验证的消息字段。
    
```protobuf
    message Person {
      option (validate.disabled) = true;
    
      // x will not be required to be greater than 123
      uint64 x = 1 [(validate.rules).uint64.gt = 123];
    
      // y's fields will not be validated
      Person y = 2;
    }
```
    
-   **忽略**：不为此消息生成验证方法或任何相关验证代码。
    
```protobuf
    message Person {
      option (validate.ignored) = true;
    
      // x will not be required to be greater than 123
      uint64 x = 1 [(validate.rules).uint64.gt = 123];
    
      // y's fields will not be validated
      Person y = 2;
    }
```
    

### [](https://github.com/envoyproxy/protoc-gen-validate#oneofs)一个人

-   **required**`oneof` ：要求必须设置a 中的字段之一。默认情况下，可以不设置或设置联合字段之一。启用此规则不允许将它们全部取消。

 
```protobuf
oneof id {
  // either x, y, or z must be set.
  option (validate.required) = true;

  string x = 1;
  int32  y = 2;
  Person z = 3;
}
```
