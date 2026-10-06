// Package apikey 管理 API key 的权限与吊销。
package apikey

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// Key 是一个 API key。
type Key struct {
	ID      string
	Owner   string
	Scopes  []string
	Expires int
	Revoked bool
}

// State 是 key 台账。
type State struct {
	Keys      []Key
	Overrides map[string][]string

	byID    map[string]int
	byOwner map[string][]int
}

// NewState 建一份台账。
func NewState() *State {
	return &State{
		Overrides: map[string][]string{},
		byID:      map[string]int{},
		byOwner:   map[string][]int{},
	}
}

func contains(scopes []string, action string, c *Counter) bool {
	for _, s := range scopes {
		c.Scanned++
		if s == action {
			return true
		}
	}
	return false
}

// Issue 签发一个 key；同一个 key 再次签发不改变既有登记。
func (st *State) Issue(id, owner string, scopes []string, expires int, c *Counter) {
	c.Scanned++
	if _, ok := st.byID[id]; ok {
		return
	}
	idx := len(st.Keys)
	st.Keys = append(st.Keys, Key{ID: id, Owner: owner, Scopes: scopes, Expires: expires})
	st.byID[id] = idx
	st.byOwner[owner] = append(st.byOwner[owner], idx)
}

// Revoke 吊销一个 key。
func (st *State) Revoke(id string, c *Counter) {
	c.Scanned++
	if idx, ok := st.byID[id]; ok {
		st.Keys[idx].Revoked = true
	}
}

// RevokeOwner 按所有者级联吊销。
func (st *State) RevokeOwner(owner string, c *Counter) {
	for _, idx := range st.byOwner[owner] {
		c.Scanned++
		st.Keys[idx].Revoked = true
	}
}

// Grant 给某个所有者的覆盖层加一个动作。
func (st *State) Grant(owner, action string) {
	st.Overrides[owner] = append(st.Overrides[owner], action)
}

// Check 判断某个 key 能否执行某个动作。
func (st *State) Check(id, action string, now int, c *Counter) bool {
	c.Scanned++
	idx, ok := st.byID[id]
	if !ok {
		return false
	}
	key := &st.Keys[idx]
	if key.Revoked {
		return false
	}
	if key.Expires > 0 && now >= key.Expires {
		return false
	}
	if contains(key.Scopes, action, c) {
		return true
	}
	return contains(st.Overrides[key.Owner], action, c)
}
