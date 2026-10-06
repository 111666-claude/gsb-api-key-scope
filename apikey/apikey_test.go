package apikey

import "testing"

func TestIssueAndCheck(t *testing.T) {
	st := NewState()
	c := &Counter{}
	st.Issue("k1", "o", []string{"read"}, 100, c)
	if !st.Check("k1", "read", 0, c) {
		t.Fatal("k1 应能 read")
	}
}

func TestUnknownKeyDenied(t *testing.T) {
	st := NewState()
	if st.Check("nope", "read", 0, &Counter{}) {
		t.Fatal("未知 key 应拒绝")
	}
}

func TestMissingScopeDenied(t *testing.T) {
	st := NewState()
	c := &Counter{}
	st.Issue("k1", "o", []string{"read"}, 100, c)
	if st.Check("k1", "write", 0, c) {
		t.Fatal("没有 write 权限应拒绝")
	}
}

func TestRevokeDenies(t *testing.T) {
	st := NewState()
	c := &Counter{}
	st.Issue("k1", "o", []string{"read"}, 100, c)
	st.Revoke("k1", c)
	if st.Check("k1", "read", 0, c) {
		t.Fatal("吊销后应拒绝")
	}
}

func TestCounterCounts(t *testing.T) {
	st := NewState()
	c := &Counter{}
	st.Issue("k1", "o", []string{"read"}, 100, c)
	st.Check("k1", "read", 0, c)
	if c.Scanned == 0 {
		t.Fatal("应统计比较次数")
	}
}

func TestGrantStored(t *testing.T) {
	st := NewState()
	st.Grant("o", "write")
	if len(st.Overrides["o"]) != 1 {
		t.Fatalf("覆盖层应记下 write：%+v", st.Overrides)
	}
}
