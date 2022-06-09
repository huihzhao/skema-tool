syntax = "proto3";
// Generated by skemabuild. DO NOT EDIT.
package {{ .Package }};
{{ .Options }}

import "google/api/annotations.proto";
import "validate/validate.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

// User Defined Protobuf Code Section
service {{ .Service }} {
    rpc Heathcheck (HealthcheckRequest) returns (HealthcheckResponse){
        option (google.api.http) = {
          get: "/api/healthcheck"
        };
    }

    rpc Helloworld (HelloworldRequest) returns (HelloworldResponse) {
        option (google.api.http) = {
            post: "/api/helloworld"
            body: "*"
        };
    };
}

message HealthcheckRequest {
}

message HealthcheckResponse {
    string result = 1;
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
        example: "{\"result\": \"ok\"}"
    };
}

message HelloworldRequest {
    string msg = 1 [(validate.rules).string.min_len = 3];
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
        example: "{\"msg\": \"hello world\"}"
    };
}

message HelloworldResponse {
    string msg = 1;
    string code = 2;
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
        example: "{\"msg\": \"hello world from server\", \"code\":\"0\"}"
    };
}