package kokizami

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
				ID:        1,
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

		if diff := cmp.Diff(ret, tc.wantKizami); diff != "" {
			t.Fatalf("unexpected result: (-got +want) %s", diff)
		}
	}
}

func TestGet(t *testing.T) {
	k, teardown := setup(t)
	defer teardown()

	tcs := []struct {
		inDesc     string
		wantErr    bool
		wantKizami *Kizami
	}{
		{
			inDesc:  "hoge",
			wantErr: false,
			wantKizami: &Kizami{
				ID:        1,
				Desc:      "hoge",
				StartedAt: mockedNow(),
				StoppedAt: initialTime,
			},
		},
		{
			inDesc:  "fuga",
			wantErr: false,
			wantKizami: &Kizami{
				ID:        2,
				Desc:      "fuga",
				StartedAt: mockedNow(),
				StoppedAt: initialTime,
			},
		},
		{
			inDesc:  "piyo",
			wantErr: false,
			wantKizami: &Kizami{
				ID:        3,
				Desc:      "piyo",
				StartedAt: mockedNow(),
				StoppedAt: initialTime,
			},
		},
	}

	for i, tc := range tcs {
		ki, err := k.Start(tc.inDesc)
		if err != nil {
			t.Fatalf("unexpected result: [got] %v [want] nil", err)
		}

		ret, err := k.Get(ki.ID)
		if err != nil {
			t.Fatalf("unexpected result: [got] %v [want] nil", err)
		}

		if diff := cmp.Diff(ki, ret); diff != "" {
			t.Fatalf("unexpected result: (-got +want) %s", diff)
		}
		if diff := cmp.Diff(ret, tc.wantKizami); diff != "" {
			t.Fatalf("unexpected result: (-got +want) %s", diff)
		}
	}
}
