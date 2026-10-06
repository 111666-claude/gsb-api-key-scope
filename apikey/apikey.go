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

	index   map[string]int
	revoked map[string]bool
}

// NewState 建一份台账。
func NewState() *State {
	return &State{
		Overrides: map[string][]string{},
		index:     map[string]int{},
		revoked:   map[string]bool{},
	}
}

func (st *State) init() {
	if st.Overrides == nil {
		st.Overrides = map[string][]string{}
	}
	if st.index == nil {
		st.index = map[string]int{}
		for i := range st.Keys {
			if _, ok := st.index[st.Keys[i].ID]; !ok {
				st.index[st.Keys[i].ID] = i
			}
		}
	}
	if st.revoked == nil {
		st.revoked = map[string]bool{}
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

// Issue 签发一个 key。
func (st *State) Issue(id, owner string, scopes []string, expires int, c *Counter) {
	st.init()
	if _, ok := st.index[id]; ok {
		return
	}
	st.index[id] = len(st.Keys)
	st.Keys = append(st.Keys, Key{ID: id, Owner: owner, Scopes: scopes, Expires: expires})
}

// Revoke 吊销一个 key。
func (st *State) Revoke(id string, c *Counter) {
	st.init()
	if i, ok := st.index[id]; ok {
		st.Keys[i].Revoked = true
	}
}

// RevokeOwner 按所有者级联吊销。
func (st *State) RevokeOwner(owner string, c *Counter) {
	st.init()
	st.revoked[owner] = true
}

// Grant 给某个所有者的覆盖层加一个动作。
func (st *State) Grant(owner, action string) {
	st.init()
	st.Overrides[owner] = append(st.Overrides[owner], action)
}

// Check 判断某个 key 能否执行某个动作。
func (st *State) Check(id, action string, now int, c *Counter) bool {
	st.init()
	i, ok := st.index[id]
	if !ok {
		return false
	}
	key := &st.Keys[i]
	if key.Revoked || st.revoked[key.Owner] {
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
