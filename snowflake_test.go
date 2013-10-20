package snowflake

import (
	"appengine"
	"appengine/aetest"
	_ "fmt"
	"strconv"
	"testing"
	"time"
)

func TestFoo(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if Epoch != 1288834974657 {
		t.Errorf("invalid Epoch: %d", Epoch)
	}
	if RandomBit != 22 {
		t.Errorf("invalid RandomBit: %d", RandomBit)
	}

	tm := strconv.FormatInt(TimeMax, 2)
	if len(tm) != int(63-RandomBit) {
		t.Errorf("invalid TimeMax: %s", tm)
	}
	if 1<<(63-RandomBit) != TimeMax+1 {
		t.Errorf("invalid TimeMax: %d", TimeMax)
	}

	rm := strconv.FormatInt(int64(RandomMax), 2)
	if len(rm) != int(RandomBit) {
		t.Errorf("invalid RandomMax: %s", rm)
	}
	if 1<<RandomBit != RandomMax+1 {
		t.Errorf("invalid RandomMax: %d", RandomMax)
	}

	now := Now()
	if now < 0 || TimeMax < now {
		t.Errorf("invalid Now(): %d", now)
	}

	nowt := time.Now()
	sid, err := Generate(c)
	if err != nil {
		t.Errorf("Generate: %v", err)
	}

	parset := Parse(sid)
	diffi := nowt.Unix() - parset.Unix()
	if diffi <= -2 || 2 <= diffi {
		t.Errorf("invalid Parse(): %v", parset)
	}

	var lastSnow int64
	for i := 0; i < 100; i++ {
		sid, err = Generate(c)
		if err != nil {
			t.Errorf("Generate: %v", err)
		}
		if lastSnow != 0 && lastSnow > sid {
			t.Errorf("Generate: %d > %d", lastSnow, sid)
		}
		lastSnow = sid
		//fmt.Println(sid)
	}

	sid, err = GenerateSameKey(c)
	if err != nil {
		t.Errorf("GenerateSameKey: %v", err)
	}
	sid2, err := GenerateSameKey(c)
	if err != nil {
		t.Errorf("GenerateSameKey: %v", err)
	} else if sid > sid2 {
		t.Errorf("Generate: %d > %d", sid, sid2)
	}
}

func GenerateSameKey(c appengine.Context) (int64, error) {
	t := Now()
	if err := update(c, t, 1); err != nil {
		return 0, err
	}
	return combine(t, 1), nil
}
