package kokizami

import (
	"testing"
	"time"
)

var mockNow = time.Now()

func mockedNow() time.Time {
	return mockNow
}

func setup(t *testing.T) (*Kokizami, func()) {
	k := &Kokizami{
		DBPath: "file::testdb?mode=memory",
		now:    mockedNow,
	}
	err := k.Initialize()
	if err != nil {
		t.Fatalf("unexpected result: %v", err)
	}

	return k, func() {
		err = k.Finalize()
		if err != nil {
			t.Fatalf("unexpected result. [got] %v [want] nil", err)
		}
	}
}

func TestInitialize(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()
}

func TestStart(t *testing.T) {
	k, teardown := setup(t)
	defer teardown()

	tcs := []struct {
		inDesc     string
		wantErr    bool
		wantKizami *Kizami
	}{
		{
			inDesc:     "",
			wantErr:    true,
			wantKizami: nil,
		},
		{
			inDesc:  "hoge",
			wantErr: false,
			wantKizami: &Kizami{
				Desc:      "hoge",
				StartedAt: mockedNow(),
				StoppedAt: initialTime,
			},
		},
	}

	for i, tc := range tcs {
		ret, err := k.Start(tc.inDesc)
		if !tc.wantErr && err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		if err != nil {
			continue
		}

		if ret.Desc != tc.wantKizami.Desc {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] %v", i, ret.Desc, tc.wantKizami.Desc)
		}
		if !ret.StartedAt.Equal(tc.wantKizami.StartedAt) {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] %v", i, ret.StartedAt, tc.wantKizami.StartedAt)
		}
		if !ret.StoppedAt.Equal(tc.wantKizami.StoppedAt) {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] %v", i, ret.StoppedAt, tc.wantKizami.StoppedAt)
		}
	}
}
