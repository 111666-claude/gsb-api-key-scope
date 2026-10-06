package apikey

import "fmt"

// Sample 是一个固定场景，用来复现题面里的锚点。
type Sample struct {
	Name string
	Run  func() map[string]any
}

// Samples 返回全部场景。
func Samples() []Sample {
	return []Sample{
		{Name: "idem", Run: func() map[string]any {
			st := NewState()
			c := &Counter{}
			st.Issue("k1", "o", []string{"read"}, 100, c)
			st.Issue("k1", "o", []string{"write"}, 100, c)
			return map[string]any{"read": st.Check("k1", "read", 0, c)}
		}},
		{Name: "cascade", Run: func() map[string]any {
			st := NewState()
			c := &Counter{}
			st.Issue("k1", "o", []string{"read"}, 100, c)
			st.Issue("k2", "o", []string{"read"}, 100, c)
			st.RevokeOwner("o", c)
			return map[string]any{"allowed": st.Check("k1", "read", 0, c)}
		}},
		{Name: "expiry", Run: func() map[string]any {
			st := NewState()
			c := &Counter{}
			st.Issue("k1", "o", []string{"read"}, 10, c)
			return map[string]any{"allowed": st.Check("k1", "read", 20, c)}
		}},
		{Name: "override", Run: func() map[string]any {
			st := NewState()
			c := &Counter{}
			st.Issue("k1", "o", nil, 100, c)
			st.Grant("o", "write")
			return map[string]any{"allowed": st.Check("k1", "write", 0, c)}
		}},
	}
}

// Work 跑规模线场景。
func Work(n int) map[string]any {
	st := NewState()
	c := &Counter{}
	for i := 0; i < n; i++ {
		st.Issue(fmt.Sprintf("k%04d", i), "o", []string{"read"}, 1000000, c)
	}
	for i := 0; i < n; i++ {
		st.Check(fmt.Sprintf("k%04d", i), "read", 0, c)
	}
	return map[string]any{"keys": n, "scanned": c.Scanned}
}
