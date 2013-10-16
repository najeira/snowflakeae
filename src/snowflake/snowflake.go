package snowflake

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"math/rand"
	"time"
)

const Epoch int64 = 1288834974657
const RandomBit uint8 = 22
const KindSnowflake = "Snowflake"

var RandomMax int32
var TimeMax int64

type Snowflake struct {
	Updated int64 `datastore:"u,noindex"`
}

func init() {
	RandomMax = (1 << RandomBit) - 1
	TimeMax = (1 << (64 - RandomBit)) - 1
	rand.Seed(time.Now().Unix())
}

func Generate(c appengine.Context) (int64, error) {
	r := rand.Int31n(RandomMax) + 1
	t := Now()
	if err := update(c, t, r); err != nil {
		return 0, err
	}
	return combine(t, r), nil
}

func Now() int64 {
	now := time.Now()
	nows := int64(now.Unix() * 1000)
	nowm := int64(now.Nanosecond()) / time.Millisecond.Nanoseconds()
	nowi := nows + nowm
	return nowi - Epoch
}

func Parse(i int64) time.Time {
	n := (i >> RandomBit) + Epoch
	sec := n / 1000
	nano := (n % 1000) * time.Millisecond.Nanoseconds()
	t := time.Unix(sec, nano)
	return t
}

func update(c appengine.Context, t int64, r int32) error {
	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		s := Snowflake{}
		key := datastore.NewKey(c, KindSnowflake, "", int64(r), nil)
		if err := datastore.Get(c, key, &s); err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}
			//fmt.Println("new key: ", key)
		} else {
			if s.Updated >= t {
				// conflict
				return fmt.Errorf("conflict")
			}
			//fmt.Println("exists key: ", key)
		}
		s.Updated = t
		if _, err := datastore.Put(c, key, &s); err != nil {
			return err
		}
		return nil
	}, nil)
}

func combine(t int64, r int32) int64 {
	if r > RandomMax {
		panic("r")
	}
	if t > TimeMax {
		panic("t")
	}
	i := t << RandomBit
	i += int64(r)
	return i
}
