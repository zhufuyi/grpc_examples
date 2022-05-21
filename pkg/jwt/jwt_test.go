package jwt

import (
	"fmt"
	"testing"
	"time"
)

func TestVerifyTokenStandard(t *testing.T) {
	token, err := GenerateTokenStandard()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)

	err = VerifyTokenStandard(token)
	if err != nil {
		t.Error(err)
	}
}

func TestVerifyTokenCustom(t *testing.T) {
	uid := "123"
	role := "admin"

	// 正常验证
	token, err := GenerateTokenWithCustom(uid, role)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
	v, err := VerifyTokenCustom(token)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)

	// 无效token格式
	token2 := "xxx.xxx.xxx"
	v, err = VerifyTokenCustom(token2)
	if !compareErr(err, formatErr) {
		t.Fatal(err)
	}

	// 签名失败
	token3 := token + "xxx"
	v, err = VerifyTokenCustom(token3)
	if !compareErr(err, signatureErr) {
		t.Fatal(err)
	}

	// token已过期
	SetExpire(time.Second)
	token, err = GenerateTokenWithCustom(uid, role)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	v, err = VerifyTokenCustom(token)
	if !compareErr(err, expiredErr) {
		t.Fatal(err)
	}
}

func compareErr(err1, err2 error) bool {
	return err1.Error() == err2.Error()
}
