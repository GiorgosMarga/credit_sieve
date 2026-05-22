package main

// type ARC struct {
// 	T1 *LRU
// 	B1 *LRU

// 	T2 *LRU
// 	B2 *LRU

// 	p int
// 	c int
// }

// func NewARC(c int) *ARC {
// 	return &ARC{
// 		p:  0,
// 		T1: NewLRU(c),
// 		T2: NewLRU(c),
// 		B1: NewLRU(c),
// 		B2: NewLRU(c),
// 		c:  c,
// 	}
// }

// func (arc *ARC) Get(key string) (string, bool) {
// 	// case 1
// 	if val, ok := arc.T1.Remove(key); ok {
// 		arc.T2.Insert(key, val)
// 		return val, true
// 	}
// 	if val, ok := arc.T2.Get(key); ok {
// 		return val, true
// 	}

// 	// case 2
// 	if arc.B1.Contains(key) {
// 		arc.p = min(arc.p+max(1, arc.B2.Len()/arc.B1.Len()), arc.c)
// 		arc.replace(key)
// 		arc.B1.Remove(key)

// 		// fetch from actual storage
// 		arc.T2.Insert(key, "RESTORED")
// 		return "RESTORED", true
// 	}

// 	// case 3
// 	if arc.B2.Contains(key) {
// 		arc.p = max(arc.p-max(1, arc.B1.Len()/arc.B2.Len()), arc.c)
// 		arc.replace(key)
// 		arc.B2.Remove(key)
// 		arc.T2.Insert(key, "RESTORED")
// 		return "RESTORED", true
// 	}

// 	// case 4 (miss)
// 	// case A
// 	if arc.T1.Len()+arc.B1.Len() == arc.c {
// 		if arc.T1.Len() < arc.c {
// 			arc.B1.Evict()
// 			arc.replace(key)
// 		} else {
// 			arc.T1.Evict()
// 		}
// 	} else if arc.T1.Len()+arc.B1.Len() < arc.c {
// 		if arc.T1.Len()+arc.B1.Len()+arc.T2.Len()+arc.B2.Len() >= arc.c {
// 			if arc.T1.Len()+arc.B1.Len()+arc.T2.Len()+arc.B2.Len() == arc.c {
// 				arc.B2.Evict()
// 				arc.replace(key)
// 			}
// 		}
// 	}
// 	arc.T1.Insert(key, "FETCHED")
// 	return "", false
// }

// func (arc *ARC) replace(key string) {

// 	if arc.T1.Len() > 0 && (arc.T1.Len() > arc.p || (arc.B2.Contains(key) && arc.T1.Len() == arc.p)) {
// 		arc.T1.Evict()
// 		arc.B1.Insert(key, "")
// 	} else {
// 		arc.T2.Evict()
// 		arc.B2.Insert(key, "")
// 	}
// }
