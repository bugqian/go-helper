package fn

import (
	"fmt"
	"strings"
	"testing"
)

type User struct {
	Name   string
	Age    int
	Gender int
}

var users = []User{
	{"张三", 20, 1},         // 按字段顺序初始化（需和结构体字段顺序完全一致）
	{Name: "李四", Age: 25}, // 按字段名初始化（推荐，可读性高，可省略部分字段）
	{Name: "王五", Age: 30, Gender: 2},
	{Name: "赵六", Age: 30, Gender: 1},
	{Name: "陈7", Age: 30, Gender: 0},
}

func TestTransform(t *testing.T) {
	list := []string{"a", "b", "c"}
	fmt.Println(Transform(list, func(s string) string {
		return strings.ToUpper(s)
	}))
}

// 重写list
func TestAddIfNotExist(t *testing.T) {
	list := []string{"a", "b", "c"}
	fmt.Println(AddIfNotExist(list, "d"))
}

func TestContains(t *testing.T) {
	list := []string{"a", "b", "c"}
	fmt.Println(Contains(list, "a"))
	fmt.Println(Contains(list, "d"))
}

func TestContainsAll(t *testing.T) {
	list := []string{"a", "b", "c"}
	list2 := []string{"a", "b"}
	list3 := []string{"b", "c", "d"}
	fmt.Println(ContainsAll(list, list2))
	fmt.Println(ContainsAll(list, list3))
}

func TestContainsFunc(t *testing.T) {
	list := []string{"a", "b", "c"}
	fmt.Println(ContainsFunc(list, func(s string) bool {
		return strings.ToLower("D") == s
	}))

	fmt.Println(ContainsFunc(users, func(s User) bool {
		return s.Age >= 20 && s.Gender == 1
	}))
}

func TestFilter(t *testing.T) {
	fmt.Println(Filter(users, func(u User) bool {
		return u.Age >= 20 && u.Gender == 1
	}))
}

func TestAsMap(t *testing.T) {
	resp := AsMap(users, func(u User) string {
		return fmt.Sprintf("%s,%d", u.Name, u.Age)
	})
	if len(resp) > 0 {
		for k, v := range resp {
			fmt.Println(k, v)
		}
	}
}

func TestAsMap2(t *testing.T) {
	resp := AsMap2(users, func(u User) (int, int) {
		return u.Age, u.Gender
	})
	if len(resp) > 0 {
		for k, v := range resp {
			if len(v) > 0 {
				fmt.Println(k, len(v))
				for kk, vv := range v {
					fmt.Println(k, kk, vv)
				}
			}
		}
	}
}
