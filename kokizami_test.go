package kokizami

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func setup(t *testing.T) (*Kokizami, func()) {
	mockNow := time.Now()

	k := &Kokizami{
		DBPath: "file::testdb?mode=memory",
		now:    func() time.Time { return mockNow },
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
				StartedAt: k.now(),
				StoppedAt: initialTime(),
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
			t.Fatalf("[No.%d] unexpected result: (-got +want) %s", i, diff)
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
				StartedAt: k.now(),
				StoppedAt: initialTime(),
			},
		},
		{
			inDesc:  "fuga",
			wantErr: false,
			wantKizami: &Kizami{
				ID:        2,
				Desc:      "fuga",
				StartedAt: k.now(),
				StoppedAt: initialTime(),
			},
		},
		{
			inDesc:  "piyo",
			wantErr: false,
			wantKizami: &Kizami{
				ID:        3,
				Desc:      "piyo",
				StartedAt: k.now(),
				StoppedAt: initialTime(),
			},
		},
	}

	for i, tc := range tcs {
		ki, err := k.Start(tc.inDesc)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		ret, err := k.Get(ki.ID)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		if diff := cmp.Diff(ki, ret); diff != "" {
			t.Fatalf("[No.%d] unexpected result: (-got +want) %s", i, diff)
		}
		if diff := cmp.Diff(ret, tc.wantKizami); diff != "" {
			t.Fatalf("[No.%d] unexpected result: (-got +want) %s", i, diff)
		}
	}
}

func TestEdit(t *testing.T) {
	k, teardown := setup(t)
	defer teardown()

	tcs := []struct {
		inDesc     string
		wantKizami *Kizami
	}{
		{
			inDesc: "hoge",
			wantKizami: &Kizami{
				Desc:      "edited hoge",
				StartedAt: k.now(),
				StoppedAt: initialTime(),
			},
		},
		{
			inDesc: "fuga",
			wantKizami: &Kizami{
				Desc:      "edited fuga",
				StartedAt: k.now(),
				StoppedAt: initialTime(),
			},
		},
		{
			inDesc: "piyo",
			wantKizami: &Kizami{
				Desc:      "edited piyo",
				StartedAt: k.now(),
				StoppedAt: initialTime(),
			},
		},
	}

	for i, tc := range tcs {
		ki, err := k.Start(tc.inDesc)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		tc.wantKizami.ID = ki.ID
		ret, err := k.Edit(tc.wantKizami)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		edited, err := k.Get(ki.ID)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		if diff := cmp.Diff(ret, edited); diff != "" {
			t.Fatalf("[No.%d] unexpected result: (-got +want) %s", i, diff)
		}
	}
}

func TestStop(t *testing.T) {
	k, teardown := setup(t)
	defer teardown()

	tcs := []struct {
		inDesc     string
		wantKizami *Kizami
	}{
		{
			inDesc: "hoge",
			wantKizami: &Kizami{
				Desc:      "hoge",
				StartedAt: k.now(),
				StoppedAt: k.now(),
			},
		},
		{
			inDesc: "fuga",
			wantKizami: &Kizami{
				Desc:      "fuga",
				StartedAt: k.now(),
				StoppedAt: k.now(),
			},
		},
		{
			inDesc: "piyo",
			wantKizami: &Kizami{
				Desc:      "piyo",
				StartedAt: k.now(),
				StoppedAt: k.now(),
			},
		},
	}

	for i, tc := range tcs {
		ki, err := k.Start(tc.inDesc)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		err = k.Stop(ki.ID)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		ret, err := k.Get(ki.ID)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		tc.wantKizami.ID = ki.ID
		if diff := cmp.Diff(ret, tc.wantKizami); diff != "" {
			t.Fatalf("[No.%d] unexpected result: (-got +want) %s", i, diff)
		}
	}
}

func TestStopAll(t *testing.T) {
	k, teardown := setup(t)
	defer teardown()

	tcs := []struct {
		inDesc     string
		wantKizami *Kizami
	}{
		{
			inDesc: "hoge",
			wantKizami: &Kizami{
				Desc:      "hoge",
				StartedAt: k.now(),
				StoppedAt: k.now(),
			},
		},
		{
			inDesc: "fuga",
			wantKizami: &Kizami{
				Desc:      "fuga",
				StartedAt: k.now(),
				StoppedAt: k.now(),
			},
		},
		{
			inDesc: "piyo",
			wantKizami: &Kizami{
				Desc:      "piyo",
				StartedAt: k.now(),
				StoppedAt: k.now(),
			},
		},
	}

	for i, tc := range tcs {
		ki, err := k.Start(tc.inDesc)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}
		tc.wantKizami.ID = ki.ID
	}

	err := k.StopAll()
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	for i, tc := range tcs {
		ret, err := k.Get(tc.wantKizami.ID)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}
		if diff := cmp.Diff(ret, tc.wantKizami); diff != "" {
			t.Fatalf("[No.%d] unexpected result: (-got +want) %s", i, diff)
		}
	}
}

func TestDelete(t *testing.T) {
	k, teardown := setup(t)
	defer teardown()

	ret, err := k.Start("hoge")
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	err = k.Delete(ret.ID)
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	_, err = k.Get(ret.ID)
	if err == nil {
		t.Fatalf("unexpected result: [got] nil [want] some error")
	}
}

func TestList(t *testing.T) {
	k, teardown := setup(t)
	defer teardown()

	expectedLen := 3
	for i := 0; i < expectedLen; i++ {
		_, err := k.Start("hoge")
		if err != nil {
			t.Fatalf("unexpected result: [got] %v [want] nil", err)
		}
	}

	ret, err := k.List()
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	if len(ret) != expectedLen {
		t.Fatalf("unexpected result: [got] %v [want] %v", len(ret), expectedLen)
	}
}
