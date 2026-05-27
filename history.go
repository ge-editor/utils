package utils

// Primitive プリミティブ型の型制約インターフェース
type Primitive interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string | ~bool
}

// History 履歴を保持するジェネリクス構造体
type History[T Primitive] struct {
	Items []T
}

// Add 要素を先頭に追加し、既存の同一要素を削除（重複除去）する
func (h *History[T]) Add(item T) {
	newItems := make([]T, 0, len(h.Items)+1)
	newItems = append(newItems, item) // 先頭に追加

	for _, v := range h.Items {
		if v != item {
			newItems = append(newItems, v) // 既存の重複分を除外
		}
	}
	h.Items = newItems
}

/*
Usage:

func main() {
	// string型での使用例
	history := History[string]{}

	history.Add("apple")
	history.Add("banana")
	history.Add("apple") // "apple" が先頭に移動し、後方の古いものは削除される
	history.Add("orange")

	fmt.Println(history.Items) // 出力: [orange apple banana]
}
*/
