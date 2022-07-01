package hash

import (
	"testing"
)

func TestCreateEncodeFunc(t *testing.T) {
	type args struct {
		secretKey string
		msg       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty secret key",
			args: args{
				secretKey: "",
				msg:       "test message",
			},
			want: "",
		},
		{
			name: "test data #1",
			args: args{
				secretKey: "/tmp/secretKey",
				msg:       "test message",
			},
			want: "6a98f91fe9f64e7e0491ed7c84bddd1673b1a76841446c2786a15c85cd3d1f14",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodeFunc := CreateEncodeFunc(tt.args.secretKey)
			hashResult := encodeFunc(tt.args.msg)
			if hashResult != tt.want {
				t.Errorf("CreateEncodeFunc() = %v, want %v", hashResult, tt.want)
			}
		})
	}
}

func TestCheckSign(t *testing.T) {
	type args struct {
		secretKey     string
		msg           string
		correctResult string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty secretKey",
			args: args{
				secretKey:     "",
				msg:           "not important",
				correctResult: "not important",
			},
			want: true,
		},
		{
			name: "test data #1",
			args: args{
				secretKey:     "/tmp/secretKey",
				msg:           "test message",
				correctResult: "not correct result",
			},
			want: false,
		},
		{
			name: "test data #1",
			args: args{
				secretKey:     "/tmp/secretKey",
				msg:           "correct result",
				correctResult: "correct result",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSign(tt.args.secretKey, tt.args.msg, tt.args.correctResult); got != tt.want {
				t.Errorf("CheckSign() = %v, want %v", got, tt.want)
			}
		})
	}
}
