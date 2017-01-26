package store4

// TODO(js) stringManager no longer manages solely strings :/

type stringManager struct {
	// strToID maps strings to IDs.
	strToID map[interface{}]uint64
	// idToStrInfo maps IDs to string info.
	idToStrInfo map[uint64]*strInfo
	// nextID holds the next term ID to issue.
	nextID uint64
}

func newStringManager() *stringManager {
	s := &stringManager{}
	s.strToID = make(map[interface{}]uint64)
	s.idToStrInfo = make(map[uint64]*strInfo)
	// TODO(js) change wildcard assumption?

	// We use "*" as a wildcard, so we give it ID 0
	// to make things easy elsewhere.
	s.strToID["*"] = 0
	// Start IDs from 1.
	s.nextID = 1
	return s
}

// strInfo holds details for each string.
type strInfo struct {
	str      interface{} // The string itself.
	refCount uint64      // Reference count.
}

// idToString returns the term for a given ID.
// The given ID must exist.
func (s *stringManager) idToString(id uint64) string {
	return s.idToStrInfo[id].str.(string)
}

// stringToID returns the ID for a given string and true
// if the string exists, and 0 and false if it does not.
func (s *stringManager) stringToID(str interface{}) (uint64, bool) {
	id, ok := s.strToID[str]
	return id, ok
}

// getOrCreateID returns an ID for a given string.
// If no existing ID is present, it creates a new one.
// For any existing string, it also increments the reference count.
func (s *stringManager) getOrCreateID(str interface{}) uint64 {
	id, ok := s.strToID[str]
	if ok {
		if id != 0 {
			s.idToStrInfo[id].refCount++
		}
	} else {
		id = s.nextID
		s.nextID++
		s.strToID[str] = id
		s.idToStrInfo[id] = &strInfo{
			str:      str,
			refCount: 1,
		}
	}
	return id
}

// releaseRef decrements a string's reference count.
// When a string is no longer referenced, it is removed
// from all maps.
// The given id must exist or releaseRef will aspolde.
func (s *stringManager) releaseRef(id uint64) {
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
