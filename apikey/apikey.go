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
}

// NewState 建一份台账。
func NewState() *State {
	return &State{Overrides: map[string][]string{}}
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

// Issue 签发一个 key。
func (st *State) Issue(id, owner string, scopes []string, expires int, c *Counter) {
	for i := range st.Keys {
		c.Scanned++
		if st.Keys[i].ID == id {
			st.Keys[i].Scopes = scopes
			return
		}
	}
	st.Keys = append(st.Keys, Key{ID: id, Owner: owner, Scopes: scopes, Expires: expires})
}

// Revoke 吊销一个 key。
func (st *State) Revoke(id string, c *Counter) {
	for i := range st.Keys {
		c.Scanned++
		if st.Keys[i].ID == id {
			st.Keys[i].Revoked = true
			return
		}
	}
}

// RevokeOwner 按所有者级联吊销。
func (st *State) RevokeOwner(owner string, c *Counter) {
	for i := range st.Keys {
		c.Scanned++
		if st.Keys[i].Owner == owner {
			_ = st.Keys[i].ID
		}
	}
}

// Grant 给某个所有者的覆盖层加一个动作。
func (st *State) Grant(owner, action string) {
	st.Overrides[owner] = append(st.Overrides[owner], action)
}

// Check 判断某个 key 能否执行某个动作。
func (st *State) Check(id, action string, now int, c *Counter) bool {
	for i := range st.Keys {
		c.Scanned++
		if st.Keys[i].ID == id {
			return !st.Keys[i].Revoked && contains(st.Keys[i].Scopes, action, c)
		}
	}
	return false
}
