# protoc 插件实现说明

## 预备知识

### 1. [Extensions](https://protobuf.dev/programming-guides/proto2/#extensions)

```protobuff
// example
// file kittens/video_ext.proto

import "kittens/video.proto";
import "media/user_content.proto";

package kittens;

// This extension allows kitten videos in a media.UserContent message.
extend media.UserContent {
  // Video is a message imported from kittens/video.proto
  repeated Video kitten_videos = 126;
}

// file media/user_content.proto

package media;

// A container message to hold stuff that a user has created.
message UserContent {
  extensions 100 to 199;
}
```

### 2. [Custom Options](https://protobuf.dev/programming-guides/proto2/#customoptions)

Protocol Buffers允许定义和使用自己的选项。但是这是大多数人不需要的高级功能。
由于选项是由 google/protobuf/descriptor.proto 中定义的Message（如 FileOptions 或 FieldOptions），    
定义自己的选项都只是扩展这些Message。例如：

```protobuf
import "google/protobuf/descriptor.proto";

extend google.protobuf.FileOptions {
  optional string my_file_option = 50000;
}
extend google.protobuf.MessageOptions {
  optional int32 my_message_option = 50001;
}
extend google.protobuf.FieldOptions {
  optional float my_field_option = 50002;
}
extend google.protobuf.OneofOptions {
  optional int64 my_oneof_option = 50003;
}
extend google.protobuf.EnumOptions {
  optional bool my_enum_option = 50004;
}
extend google.protobuf.EnumValueOptions {
  optional uint32 my_enum_value_option = 50005;
}
extend google.protobuf.ServiceOptions {
  optional MyEnum my_service_option = 50006;
}
extend google.protobuf.MethodOptions {
  optional MyMessage my_method_option = 50007;
}

option (my_file_option) = "Hello world!";

message MyMessage {
  option (my_message_option) = 1234;

  optional int32 foo = 1 [(my_field_option) = 4.5];
  optional string bar = 2;
  oneof qux {
    option (my_oneof_option) = 42;

    string quux = 3;
  }
}

enum MyEnum {
  option (my_enum_option) = true;

  FOO = 1 [(my_enum_value_option) = 321];
  BAR = 2;
}

message RequestType {}
message ResponseType {}

service MyService {
  option (my_service_option) = FOO;

  rpc MyMethod(RequestType) returns(ResponseType) {
    // Note:  my_method_option has type MyMessage.  We can set each field
    //   within it using a separate "option" line.
    option (my_method_option).foo = 567;
    option (my_method_option).bar = "Some string";
  }

  rpc OtherMethod(RequestType) returns(ResponseType) {
    option (my_method_option) = {
      foo: 567,
      bar: "Some string"
    };
  }
}

// ---------------------------------------------------------
// 请注意，如果您想在定义它的包之外的包中使用自定义选项，则必须在选项名称前加上包名称前缀，就像您在类型名称中所做的那样。例如
// foo.proto
import "google/protobuf/descriptor.proto";
package foo;
extend google.protobuf.MessageOptions {
  optional string my_option = 51234;
}

// bar.proto
import "foo.proto";
package bar;
message OtherMessage {
  option (foo.my_option) = "Hello world!";
}
```

### 3. protoc 插件名命名规范

plugin必须命名为"protoc-gen-\$NAME"，然后在标志"--\${NAME}_out"被传递给protoc命令。

### 4. Go编写插件引用的库

```shell
go get google.golang.org/protobuf/compiler/protogen
```

## reference

[protobuf](https://github.com/protocolbuffers/protobuf)  protoc 命令行工具
[proto buffers document](https://protobuf.dev) 
[protobuf-go](https://github.com/protocolbuffers/protobuf-go)  go代码生成 (protoc 插件)

https://rotemtam.com/2021/03/22/creating-a-protoc-plugin-to-gen-go-code/