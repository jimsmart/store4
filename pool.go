package store4

type pool struct {
	// strToID maps strings to IDs.
	strToID map[string]uint64
	// idToStrInfo maps IDs to string info.
	idToStrInfo map[uint64]*strInfo
	// nextStrID holds the next string ID to issue.
	nextStrID uint64
	// itemToID maps non-string items to IDs.
	itemToID map[interface{}]uint64
	// idToItemInfo maps IDs to item info.
	idToItemInfo map[uint64]*itemInfo
	// nextItemID holds the next non-string ID to issue.
	nextItemID uint64
}

func newPool() *pool {
	s := &pool{}
	s.strToID = make(map[string]uint64)
	s.idToStrInfo = make(map[uint64]*strInfo)
	s.itemToID = make(map[interface{}]uint64)
	s.idToItemInfo = make(map[uint64]*itemInfo)

	// We use "*" as a wildcard, so we give it ID 0
	// to make things easy elsewhere.
	s.strToID["*"] = 0
	// Start string IDs from 1.
	s.nextStrID = 1
	// Start non-string IDs from 0 - with the highest bit set.
	s.nextItemID = 1 << 63
	return s
}

// strInfo holds details for each string.
type strInfo struct {
	str      string // The string itself.
	refCount uint64 // Reference count.
}

// ifaceInfo holds details for each non-string.
type itemInfo struct {
	item     interface{} // The item itself.
	refCount uint64      // Reference count.
}

// idToString returns the string for a given ID.
// The given ID must exist.
func (s *pool) idToString(id uint64) string {
	return s.idToStrInfo[id].str
}

// idToAny returns the item for a given ID.
// The given ID must exist.
func (s *pool) idToAny(id uint64) interface{} {
	if id&(1<<63) == 0 {
		return s.idToStrInfo[id].str
	}
	return s.idToItemInfo[id].item
}

// stringToID returns the ID for a given string and true
// if the string exists, and 0 and false if it does not.
func (s *pool) stringToID(str string) (uint64, bool) {
	id, ok := s.strToID[str]
	return id, ok
}

// anyToID returns the ID for a given item and true
// if the item exists, and 0 and false if it does not.
func (s *pool) anyToID(item interface{}) (uint64, bool) {
	if str, sok := item.(string); sok {
		return s.stringToID(str)
	}
	id, ok := s.itemToID[item]
	return id, ok
}

// getOrCreateIDString returns an ID for a given string.
// If no existing ID is present, it creates a new one.
// For any existing string, it also increments the reference count.
func (s *pool) getOrCreateIDString(str string) uint64 {
	id, ok := s.strToID[str]
	if ok {
		if id != 0 {
			s.idToStrInfo[id].refCount++
		}
	} else {
		id = s.nextStrID
		s.nextStrID++
		s.strToID[str] = id
		s.idToStrInfo[id] = &strInfo{
			str:      str,
			refCount: 1,
		}
	}
	return id
}

// getOrCreateIDAny returns an ID for a given item.
func (s *pool) getOrCreateIDAny(item interface{}) uint64 {
	if str, sok := item.(string); sok {
		return s.getOrCreateIDString(str)
	}
	id, ok := s.itemToID[item]
	if ok {
		if id != 0 {
			s.idToItemInfo[id].refCount++
		}
	} else {
		id = s.nextItemID
		s.nextItemID++
		s.itemToID[item] = id
		s.idToItemInfo[id] = &itemInfo{
			item:     item,
			refCount: 1,
		}
	}
	return id
}

// releaseRefString decrements a string's reference count.
// When a string is no longer referenced, it is removed
// from all maps.
// The given id must exist or releaseStringRef will aspolde.
func (s *pool) releaseRefString(id uint64) {
	info := s.idToStrInfo[id]
	c := info.refCount
	c--
	if c == 0 {
		delete(s.strToID, info.str)
		delete(s.idToStrInfo, id)
		return
	}
	info.refCount = c
}

func (s *pool) releaseRefAny(id uint64) {
	if id&(1<<63) == 0 {
		s.releaseRefString(id)
		return
	}
	info := s.idToItemInfo[id]
	c := info.refCount
	c--
	if c == 0 {
		delete(s.itemToID, info.item)
		delete(s.idToItemInfo, id)
		return
	}
	info.refCount = c
}
