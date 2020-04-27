// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package triage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/triage-party/pkg/hubbub"
	"k8s.io/klog"
)

// Collection represents a fully loaded YAML configuration
type Collection struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description,omitempty"`
	RuleIDs      []string `yaml:"rules"`
	Dedup        bool     `yaml:"dedup,omitempty"`
	Hidden       bool     `yaml:"hidden,omitempty"`
	UsedForStats bool     `yaml:"used_for_statistics,omitempty"`
}

// The result of Execute
type CollectionResult struct {
	Time        time.Time
	RuleResults []*RuleResult

	Total             int
	TotalPullRequests int
	TotalIssues       int

	AvgAge             time.Duration
	AvgCurrentHold     time.Duration
	AvgAccumulatedHold time.Duration

	TotalAgeDays             float64
	TotalCurrentHoldDays     float64
	TotalAccumulatedHoldDays float64
}

// ExecuteCollection executes a collection.
func (p *Party) ExecuteCollection(ctx context.Context, s Collection) (*CollectionResult, error) {
	klog.Infof(">>> Executing collection %q: %s", s.ID, s.RuleIDs)
	start := time.Now()

	os := []*RuleResult{}
	seen := map[string]*Rule{}
	seenRule := map[string]bool{}

	for _, tid := range s.RuleIDs {
		if seenRule[tid] {
			klog.Errorf("collection %q has a duplicate rule: %q - ignoring", s.ID, tid)
			continue
		}

		seenRule[tid] = true

		t, err := p.LookupRule(tid)
		if err != nil {
			return nil, err
		}

		ro, err := p.ExecuteRule(ctx, t, seen)
		if err != nil {
			return nil, fmt.Errorf("rule %q: %v", t.Name, err)
		}

		os = append(os, ro)
	}

	r := SummarizeCollectionResult(os)
	r.Time = time.Now()
	klog.Infof("<<< Collection %q took %s to execute", s.ID, time.Since(start))
	return r, nil
}

// SummarizeCollectionResult adds together statistics about collection results {
func SummarizeCollectionResult(os []*RuleResult) *CollectionResult {
	klog.Infof("Summarizing collection result with %s rules...", len(os))

	r := &CollectionResult{}

	for _, oc := range os {
		klog.Infof("total age is %.1f days", r.TotalAgeDays)

		r.Total += len(oc.Items)
		if oc.Rule.Type == hubbub.PullRequest {
			r.TotalPullRequests += len(oc.Items)
		} else {
			r.TotalIssues += len(oc.Items)
		}

		r.RuleResults = append(r.RuleResults, oc)

		r.TotalAgeDays += oc.TotalAgeDays
		r.TotalCurrentHoldDays += oc.TotalCurrentHoldDays
		r.TotalAccumulatedHoldDays += oc.TotalAccumulatedHoldDays

	}
	if r.Total == 0 {
		klog.Warningf("no summary, total=0")
		return r
	}

	r.AvgAge = avgDayDuration(r.TotalAgeDays, r.Total)
	r.AvgCurrentHold = avgDayDuration(r.TotalCurrentHoldDays, r.Total)
	r.AvgAccumulatedHold = avgDayDuration(r.TotalAccumulatedHoldDays, r.Total)
	return r
}

func avgDayDuration(total float64, count int) time.Duration {
	return time.Duration(int64(total/float64(count)*24)) * time.Hour
}

// Flush the search cache for a collection
func (p *Party) FlushSearchCache(id string, minAge time.Duration) error {
	s, err := p.LookupCollection(id)
	if err != nil {
		return err
	}

	flushed := map[string]bool{}
	for _, tid := range s.RuleIDs {
		t, err := p.LookupRule(tid)
		if err != nil {
			return err
		}
		for _, r := range t.Repos {
			if !flushed[r] {
				klog.Infof("Flushing search cache for %s ...", r)
				org, project, err := parseRepo(r)
				if err != nil {
					return err
				}
				if err := p.engine.FlushSearchCache(org, project, minAge); err != nil {
					klog.Warningf("flush for %s/%s: %v", org, project, err)
				}
				flushed[r] = true
			}
		}
	}
	return nil
}

// ListCollections a fully resolved collections
func (p *Party) ListCollections() ([]Collection, error) {
	return p.collections, nil
}

// Return a fully resolved collection
func (p *Party) LookupCollection(id string) (Collection, error) {
	for _, s := range p.collections {
		if s.ID == id {
			return s, nil
		}
	}
	return Collection{}, fmt.Errorf("%q not found", id)
}
