package kokizami

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type mockRepo struct {
	kizamis  map[string]*Kizami
	relation map[int][]int
	tags     map[string]*Tag
}

type mockKizamiRepo struct {
	now  func() time.Time
	repo *mockRepo
}

type mockTagRepo struct {
	repo *mockRepo
}

type mockSummaryRepo struct{}

func (m *mockKizamiRepo) AllKizami() ([]*Kizami, error) {
	ks := make([]Kizami, len(m.repo.kizamis))
	c := 0
	for k := range m.repo.kizamis {
		ks[c] = *m.repo.kizamis[k]
		c++
	}

	ret := make([]*Kizami, len(ks))
	for i := range ks {
		ret[i] = &ks[i]
	}

	return ret, nil
}

func (m *mockKizamiRepo) Insert(desc string) (*Kizami, error) {
	id := len(m.repo.kizamis) + 1
	k := &Kizami{
		ID:        id,
		Desc:      desc,
		StartedAt: m.now(),
		StoppedAt: initialTime(),
	}
	m.repo.kizamis[strconv.Itoa(id)] = k
	return k, nil
}

func (m *mockKizamiRepo) Update(k *Kizami) error {
	m.repo.kizamis[strconv.Itoa(k.ID)] = k
	return nil
}

func (m *mockKizamiRepo) Delete(k *Kizami) error {
	delete(m.repo.kizamis, strconv.Itoa(k.ID))
	return nil
}

func (m *mockKizamiRepo) KizamiByID(id int) (*Kizami, error) {
	if k, ok := m.repo.kizamis[strconv.Itoa(id)]; ok {
		return k, nil
	}
	return nil, fmt.Errorf("Kizami that has id [%d] is not found", id)
}

func (m *mockKizamiRepo) KizamisByStoppedAt(t time.Time) ([]*Kizami, error) {
	ret := []*Kizami{}
	for k, v := range m.repo.kizamis {
		if v.StoppedAt == t {
			ret = append(ret, m.repo.kizamis[k])
		}
	}
	return ret, nil
}

func (m *mockKizamiRepo) Tagging(kizamiID int, tagIDs []int) error {
	m.repo.relation[kizamiID] = tagIDs
	return nil
}

func (m *mockKizamiRepo) Untagging(kizamiID int) error {
	delete(m.repo.relation, kizamiID)
	return nil
}

func (m *mockTagRepo) FindTagByID(id int) (*Tag, error) {
	panic("not implemented")
}

func (m *mockTagRepo) FindAllTags() ([]*Tag, error) {
	tags := make([]Tag, len(m.repo.tags))
	c := 0
	for k := range m.repo.tags {
		tags[c] = *m.repo.tags[k]
		c++
	}

	ret := make([]*Tag, len(tags))
	for i := range tags {
		ret[i] = &tags[i]
	}

	return ret, nil
}

func (m *mockTagRepo) FindTagsByKizamiID(kizamiID int) ([]*Tag, error) {
	ret := []*Tag{}
	tids := m.repo.relation[kizamiID]
	for _, v := range tids {
		if t, ok := m.repo.tags[strconv.Itoa(v)]; ok {
			ret = append(ret, t)
		}
	}

	return ret, nil
}

func (m *mockTagRepo) FindTagsByLabels(labels []string) ([]*Tag, error) {
	ret := []*Tag{}
	for _, label := range labels {
		for k, v := range m.repo.tags {
			if v.Label == label {
				ret = append(ret, m.repo.tags[k])
			}
		}
	}

	return ret, nil
}

func (m *mockTagRepo) InsertTags(labels []string) error {
	for i := range labels {
		id := len(m.repo.tags) + 1
		t := &Tag{
			ID:    id,
			Label: labels[i],
		}
		m.repo.tags[strconv.Itoa(id)] = t
	}

	return nil
}

func (m *mockTagRepo) Delete(id int) error {
	delete(m.repo.tags, strconv.Itoa(id))
	return nil
}

func (m *mockSummaryRepo) ElapsedOfMonthByDesc(yyyymm string) ([]*Elapsed, error) {
	return nil, nil
}

func (m *mockSummaryRepo) ElapsedOfMonthByTag(yyyymm string) ([]*Elapsed, error) {
	return nil, nil
}

func setup() *Kokizami {
	mockNow := time.Now()
	repo := &mockRepo{
		kizamis:  map[string]*Kizami{},
		relation: map[int][]int{},
		tags:     map[string]*Tag{},
	}
	return &Kokizami{
		now: func() time.Time { return mockNow },

		KizamiRepo: &mockKizamiRepo{
			now:  func() time.Time { return mockNow },
			repo: repo,
		},
		TagRepo: &mockTagRepo{
			repo: repo,
		},
		SummaryRepo: &mockSummaryRepo{},
	}
}

func TestStart(t *testing.T) {
	k := setup()

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
	k := setup()

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
	k := setup()

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
	k := setup()

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
	k := setup()

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
	k := setup()

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
	k := setup()

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

func TestAddTags(t *testing.T) {
	tcs := [][]string{
		{"hoge", "fuga", "piyo"},
		{},
	}

	for i, tags := range tcs {
		k := setup()

		err := k.AddTags(tags)
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		ret, err := k.Tags()
		if err != nil {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] nil", i, err)
		}

		m := make(map[string]struct{})
		for _, v := range ret {
			m[v.Label] = struct{}{}
		}

		if len(ret) != len(tags) {
			t.Fatalf("[No.%d] unexpected result: [got] %v [want] %v", i, len(ret), len(tags))
		}

		for j := range tags {
			if _, ok := m[tags[j]]; !ok {
				t.Fatalf("[No.%d] unexpected result: %v is missing", i, tags[j])
			}
		}
	}
}

func TestDeleteTag(t *testing.T) {
	k := setup()

	inTags := []string{"hoge", "fuga", "piyo"}

	err := k.AddTags(inTags)
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	ret, err := k.TagsByLabels([]string{"fuga"})
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	if len(ret) != 1 {
		t.Fatalf("unexpected result: [got] %v [want] %v", len(ret), 1)
	}

	err = k.DeleteTag(ret[0].ID)
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	ret, err = k.Tags()
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	if len(ret) != 2 {
		t.Fatalf("unexpected result: [got] %v [want] %v", len(ret), 2)
	}
}

func TestTagging(t *testing.T) {
	k := setup()

	// prepare tags
	inTags := []string{"hoge", "fuga", "piyo"}
	err := k.AddTags(inTags)
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	tags, err := k.Tags()
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	// prepare kizami
	ki, err := k.Start("foo")
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	tids := make([]int, len(tags))
	for i := range tags {
		tids[i] = tags[i].ID
	}

	err = k.Tagging(ki.ID, tids)
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	tagged, err := k.TagsByKizamiID(ki.ID)
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	m := make(map[string]struct{})
	for _, v := range tagged {
		m[v.Label] = struct{}{}
	}

	for i := range tags {
		if _, ok := m[tags[i].Label]; !ok {
			t.Fatalf("unexpected result: %v is missing", tags[i])
		}
	}

	if len(tagged) != len(tags) {
		t.Fatalf("unexpected result: [got] %v [want] %v", len(tagged), len(tags))
	}

	err = k.Untagging(ki.ID)
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	tagged, err = k.TagsByKizamiID(ki.ID)
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}

	if len(tagged) != 0 {
		t.Fatalf("unexpected result: [got] %v [want] %v", len(tagged), 0)
	}
}

func TestSummaryByTag(t *testing.T) {
	k := setup()

	_, err := k.SummaryByTag("2019-05")
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}
}

func TestSummaryByDesc(t *testing.T) {
	k := setup()

	_, err := k.SummaryByDesc("2019-05")
	if err != nil {
		t.Fatalf("unexpected result: [got] %v [want] nil", err)
	}
}
