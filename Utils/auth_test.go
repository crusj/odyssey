package Utils

import (
	"testing"
)

var (
	token *Token = &Token{}
)

func  TestEncodeAndDecode(t *testing.T){
	token.SetSecret("abcdejl")
	tokenString, err := token.CreateToken(1)
	if err != nil {
		t.Fatal(err,token.secret)
	}
	id, err := token.DecodeToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatal("ID解析错误")
	}
}
