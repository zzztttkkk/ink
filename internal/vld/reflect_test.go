package vld

import (
	"fmt"
	"github.com/zzztttkkk/ink/internal/utils"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type AB struct {
	A string
	B string
}

func (ab *AB) FromRequest(_ *http.Request) *Error {
	ab.A = fmt.Sprintf("%d", time.Now().UnixNano())
	ab.B = fmt.Sprintf("%d", time.Now().UnixNano())
	return nil
}

type C string

func (c *C) FromRequest(_ *http.Request) *Error {
	*c = C(fmt.Sprintf("%d", time.Now().UnixNano()))
	return nil
}

type User struct {
	Name      string    `vld:"name;RuneCountRange=1-20"`
	Email     string    `vld:"email"`
	Age       int       `vld:"age"`
	CreatedAt time.Time `vld:"created_at"`
	Nums      []int     `vld:"nums;NumRange=1-30;LenRange=4-4"`
	AB1       AB        `vld:"ab"`
	AB2       *AB       `vld:"ab_ptr"`
	C2        *C        `vld:"c_ptr"`
}

func TestGetRules(t *testing.T) {
	rules := GetRules(reflect.TypeOf(User{}))
	req := utils.Must(http.NewRequest("Post", "/", nil))
	req.PostForm = map[string][]string{}
	req.PostForm["name"] = []string{"ztk<Spk>"}
	req.PostForm["email"] = []string{"ztk@local.dev"}
	req.PostForm["age"] = []string{"123"}
	req.PostForm["created_at"] = []string{"189123000"}
	req.PostForm["nums"] = []string{"1", "21", "3", "4"}

	var user User
	e := rules.BindAndValidate(req, reflect.ValueOf(&user).Elem())
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Printf("%#v\n", user)
	}
}

func TestValidate(t *testing.T) {
	v := User{
		Name:  "ztk",
		Email: "ztk@local.dev",
	}

	if err := GetRules(reflect.TypeOf(v)).Validate(&v); err != nil {
		fmt.Println(err)
	}
}
