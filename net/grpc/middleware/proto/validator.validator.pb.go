// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: validator.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	math "math"
	regexp "regexp"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var _regex_ValidateRequest_Name = regexp.MustCompile(`^[a-z]{2,5}$`)

func (this *ValidateRequest) Validate() error {
	if !(this.Id > 0) {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must be greater than '0'`, this.Id))
	}
	if !(this.Id < 100) {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must be less than '100'`, this.Id))
	}
	if !_regex_ValidateRequest_Name.MatchString(this.Name) {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]{2,5}$"`, this.Name))
	}
	return nil
}
func (this *InnerMessage) Validate() error {
	if !(this.SomeInteger > 0) {
		return github_com_mwitkow_go_proto_validators.FieldError("SomeInteger", fmt.Errorf(`value '%v' must be greater than '0'`, this.SomeInteger))
	}
	if !(this.SomeInteger < 100) {
		return github_com_mwitkow_go_proto_validators.FieldError("SomeInteger", fmt.Errorf(`value '%v' must be less than '100'`, this.SomeInteger))
	}
	if !(this.SomeFloat >= 0) {
		return github_com_mwitkow_go_proto_validators.FieldError("SomeFloat", fmt.Errorf(`value '%v' must be greater than or equal to '0'`, this.SomeFloat))
	}
	if !(this.SomeFloat <= 1) {
		return github_com_mwitkow_go_proto_validators.FieldError("SomeFloat", fmt.Errorf(`value '%v' must be lower than or equal to '1'`, this.SomeFloat))
	}
	return nil
}

var _regex_ValidateResponse_Message = regexp.MustCompile(`^[a-z]{2,5}$`)

func (this *ValidateResponse) Validate() error {
	if !(this.Code > 100) {
		return github_com_mwitkow_go_proto_validators.FieldError("Code", fmt.Errorf(`value '%v' must be greater than '100'`, this.Code))
	}
	if !(this.Code < 900) {
		return github_com_mwitkow_go_proto_validators.FieldError("Code", fmt.Errorf(`value '%v' must be less than '900'`, this.Code))
	}
	if !_regex_ValidateResponse_Message.MatchString(this.Message) {
		return github_com_mwitkow_go_proto_validators.FieldError("Message", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]{2,5}$"`, this.Message))
	}
	if nil == this.Inner {
		return github_com_mwitkow_go_proto_validators.FieldError("Inner", fmt.Errorf("message must exist"))
	}
	if this.Inner != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Inner); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Inner", err)
		}
	}
	return nil
}
