package fn

import (
	"sort"
)

func Transform[S any, T any, List ~[]S](li List, transform func(S) T) []T {
	res := make([]T, len(li))
	for i, v := range li {
		res[i] = transform(v)
	}
	return res
}

// Unique slice去重
func Unique[V comparable](list []V) []V {
	result := make([]V, 0, len(list))
	tmpMap := make(map[V]int, len(list))
	for _, v := range list {
		if _, ok := tmpMap[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// AsMap 将切片转换为map，key由keyFunc生成，value为切片元素本身
// K: map的key类型（必须可比较）
// V: 切片元素和map的value类型
func AsMap[K comparable, V any, List ~[]V](slice List, keyFunc func(V) K) map[K]V {
	// 初始化map，容量设为切片长度减少扩容
	result := make(map[K]V, len(slice))
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = item
	}
	return result
}

// AsMap2 将切片转换为二级嵌套map
// K1: 第一层map的key类型（可比较）
// K2: 第二层map的key类型（可比较）
// V: 切片元素类型（任意）
// 参数：
//
//	slice: 源切片
//	keyFunc: 从元素提取k1和k2的函数（一个函数返回两个key）
//
// 返回值：二级嵌套map（map[K1]map[K2]V）
func AsMap2[K1 comparable, K2 comparable, V any](slice []V, keyFunc func(V) (K1, K2)) map[K1]map[K2]V {
	// 初始化外层map，容量设为切片长度减少扩容
	result := make(map[K1]map[K2]V, len(slice))

	for _, item := range slice {
		// 单个函数提取k1和k2，集中管理key逻辑
		k1, k2 := keyFunc(item)

		// 外层map无当前k1时，初始化内层map
		if _, exists := result[k1]; !exists {
			result[k1] = make(map[K2]V)
		}
		// 赋值到对应层级
		result[k1][k2] = item
	}

	return result
}

// Contains 判断切片中是否包含指定元素（适用于可比较类型）
// T: 必须是可比较类型（comparable）
func Contains[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// ContainsFunc 自定义判断逻辑（支持不可比较类型，如切片、map等）
// T: 任意类型
func ContainsFunc[T any](slice []T, match func(T) bool) bool {
	for _, item := range slice {
		if match(item) {
			return true
		}
	}
	return false
}

// ContainsAll 优化版：判断parent切片是否包含child切片的所有元素（非连续）
// 核心逻辑：先将父切片转map，再遍历子切片快速查询（O(n+m)）
func ContainsAll[T comparable](parent, child []T) bool {
	// 边界条件快速判断
	if len(child) == 0 {
		return true
	}
	if len(parent) < len(child) {
		return false
	}
	// 父切片转map（用空结构体节省内存）
	parentMap := make(map[T]struct{}, len(parent))
	for _, v := range parent {
		parentMap[v] = struct{}{}
	}
	// 遍历子切片，检查每个元素是否在map中
	for _, v := range child {
		if _, exists := parentMap[v]; !exists {
			return false
		}
	}
	return true
}

// Intersection 计算两个切片的交集（去重），基于map优化（O(n+m)）
// T: 可比较类型（comparable）
func Intersection[T comparable](a, b []T) []T {
	// 初始化map存储第一个切片的元素（空结构体节省内存）
	elementMap := make(map[T]struct{}, len(a))
	for _, v := range a {
		elementMap[v] = struct{}{}
	}
	// 遍历第二个切片，收集共同元素（自动去重）
	var intersection []T
	for _, v := range b {
		if _, exists := elementMap[v]; exists {
			intersection = append(intersection, v)
			// 删除已找到的元素，避免重复添加（如b中有重复元素时）
			delete(elementMap, v)
		}
	}
	return intersection
}

// Filter 泛型筛选函数：根据自定义条件筛选切片元素，返回新切片
// T: 切片元素的任意类型
// filterFunc: 筛选条件函数，返回true表示该元素符合条件，会被保留
func Filter[T any](slice []T, filterFunc func(item T) bool) []T {
	// 初始化新切片，容量设为原切片长度（减少扩容）
	var result []T = make([]T, 0, len(slice))

	// 遍历原切片，筛选符合条件的元素
	for _, item := range slice {
		if filterFunc(item) {
			result = append(result, item)
		}
	}

	return result
}

// Find 查找切片中第一个匹配条件的元素
// T: 切片元素的任意类型
// matchFunc: 匹配条件函数，返回true表示元素符合要求
// 返回值：第一个匹配的元素、是否找到（bool）
func Find[T any](slice []T, matchFunc func(item T) bool) (res T, found bool) {
	// 遍历切片，找到第一个匹配的元素立即返回
	for _, item := range slice {
		if matchFunc(item) {
			return item, true
		}
	}
	return
}

// FindAll 查找切片中所有匹配条件的元素
// T: 切片元素的任意类型
// matchFunc: 匹配条件函数，返回true表示元素符合要求
// 返回值：1. 所有匹配的元素组成的新切片 2. 是否找到匹配元素（bool）
func FindAll[T any](slice []T, matchFunc func(item T) bool) ([]T, bool) {
	// 初始化结果切片，容量设为原切片长度（减少扩容）
	result := make([]T, 0, len(slice))

	// 遍历切片，收集所有匹配的元素
	for _, item := range slice {
		if matchFunc(item) {
			result = append(result, item)
		}
	}

	// 判断是否找到：结果切片长度>0则为true，否则false
	found := len(result) > 0
	return result, found
}

// SliceSort 通用切片排序函数：仅需传入切片和自定义比较函数
// compareFunc返回true表示a应该排在b前面（即a < b）
func SliceSort[T any](list []T, compareFunc func(a, b T) bool) {
	sort.Slice(list, func(i, j int) bool {
		return compareFunc(list[i], list[j])
	})
}

// GroupToMap 将切片按指定key分组，转换为 map[K][]V（map[k]slice）结构
// K: 分组key的类型（必须可比较）
// V: 切片元素的类型（任意）
// 参数：
//
//	srcSlice: 源切片（待分组的原始列表）
//	keyFunc: 从元素提取分组key的函数
//
// 返回值：map[K][]V —— 每个key对应一个包含所有同key元素的切片（slice）
func GroupToMap[K comparable, V any, List ~[]V](list List, keyFunc func(V) K) map[K]List { // map[k]slice 核心结构：key对应一个切片
	// 初始化map，容量设为源切片长度（减少扩容）
	groupMap := make(map[K]List, len(list))

	// 遍历源切片，按key分组聚合到对应slice中
	for _, item := range list {
		key := keyFunc(item)
		// 将元素追加到对应key的slice中（自动初始化空slice）
		groupMap[key] = append(groupMap[key], item)
	}
	return groupMap
}

// GroupToMap2 双层分组：将List按k1、k2两级key分组为 map[K1]map[K2]List
// K1: 第一层分组key类型（可比较）
// K2: 第二层分组key类型（可比较）
// V: 切片元素基础类型
// List: 自定义切片类型（约束为~[]V，兼容任意[]V类型）
// 参数：
//
//	list: 源列表（自定义切片类型）
//	keyFunc: 单个函数返回k1和k2（集中管理两级key提取逻辑）
//
// 返回值：map[K1]map[K2]List —— 双层key对应元素列表
func GroupToMap2[K1 comparable, K2 comparable, V any, List ~[]V](list List, keyFunc func(V) (K1, K2)) map[K1]map[K2]List {
	// 初始化外层map，容量设为源列表长度减少扩容
	groupMap := make(map[K1]map[K2]List, len(list))

	// 遍历源列表，按两级key分组聚合
	for _, item := range list {
		// 提取当前元素的k1和k2
		k1, k2 := keyFunc(item)

		// 外层map无k1时，初始化内层map
		if _, exists := groupMap[k1]; !exists {
			groupMap[k1] = make(map[K2]List)
		}
		// 将元素追加到 k1→k2 对应的List中（自动初始化空List）
		groupMap[k1][k2] = append(groupMap[k1][k2], item)
	}

	return groupMap
}

// GroupToMap3 三级分组：将List按k1、k2、k3三级key分组为 map[K1]map[K2]map[K3]List
// K1/K2/K3: 各级分组key类型（均需可比较）
// V: 切片元素基础类型
// List: 自定义切片类型（约束为~[]V，兼容任意[]V类型）
// 参数：
//
//	list: 源列表（自定义切片类型）
//	keyFunc: 单个函数返回k1、k2、k3（集中管理三级key提取逻辑）
//
// 返回值：map[K1]map[K2]map[K3]List —— 三级key对应元素列表
func GroupToMap3[K1 comparable, K2 comparable, K3 comparable, V any, List ~[]V](list List, keyFunc func(V) (K1, K2, K3)) map[K1]map[K2]map[K3]List {
	// 初始化外层map，容量设为源列表长度减少扩容
	groupMap := make(map[K1]map[K2]map[K3]List, len(list))

	// 遍历源列表，按三级key分组聚合
	for _, item := range list {
		// 提取当前元素的k1、k2、k3
		k1, k2, k3 := keyFunc(item)

		// 外层map无k1时，初始化第二层map
		if _, exists := groupMap[k1]; !exists {
			groupMap[k1] = make(map[K2]map[K3]List)
		}
		// 第二层map无k2时，初始化第三层map
		if _, exists := groupMap[k1][k2]; !exists {
			groupMap[k1][k2] = make(map[K3]List)
		}
		// 将元素追加到 k1→k2→k3 对应的List中（自动初始化空List）
		groupMap[k1][k2][k3] = append(groupMap[k1][k2][k3], item)
	}

	return groupMap
}

// AddIfNotExist 检查切片中是否存在元素，不存在则新增
// T 约束为 comparable，确保可以用 == 比较元素
func AddIfNotExist[T comparable](list []T, elem ...T) (res []T) {
	res = list
	if len(elem) == 0 {
		return
	}
	// 遍历切片检查元素是否存在
	for _, item := range elem {
		if !Contains(list, item) {
			res = append(res, item)
		}
	}
	return
}
