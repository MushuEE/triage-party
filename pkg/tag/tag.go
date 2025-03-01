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

package tag

import "fmt"

// Tag is used for automatically labelling issues
type Tag struct {
	ID   string `json:"id"`
	Desc string `json:"description"`

	NeedsComments bool
	NeedsReviews  bool
	NeedsTimeline bool
}

var (
	// Simple tags
	Assigned      = Tag{ID: "assigned", Desc: "Someone is assigned"}
	Closed        = Tag{ID: "closed", Desc: "This item has been closed"}
	OpenMilestone = Tag{ID: "open-milestone", Desc: "The issue is associated to an open milestone"}
	Similar       = Tag{ID: "similar", Desc: "Title appears similar to another PR or issue"}
	Merged        = Tag{ID: "merged", Desc: "PR was merged"}
	Draft         = Tag{ID: "draft", Desc: "Draft PR"}

	// Comment-based tags
	Commented       = Tag{ID: "commented", Desc: "A project member has commented on this", NeedsComments: true}
	Send            = Tag{ID: "send", Desc: "A project member commented more recently than the author", NeedsComments: true}
	Recv            = Tag{ID: "recv", Desc: "The author commented more recently than a project member", NeedsComments: true}
	RecvQ           = Tag{ID: "recv-q", Desc: "The author has asked a question since the last project member commented", NeedsComments: true}
	AuthorLast      = Tag{ID: "author-last", Desc: "The last commenter was the original author", NeedsComments: true}
	AssigneeUpdated = Tag{ID: "assignee-updated", Desc: "Issue has been updated by its assignee", NeedsComments: true}

	// Timeline-based tags
	XrefApproved            = Tag{ID: "pr-approved", Desc: "Last review was an approval", NeedsTimeline: true}
	XrefReviewedWithComment = Tag{ID: "pr-reviewed-with-comment", Desc: "Last review was a comment", NeedsTimeline: true}
	XrefChangesRequested    = Tag{ID: "pr-changes-requested", Desc: "Last review was a request for changes", NeedsTimeline: true}
	XrefNewCommits          = Tag{ID: "pr-new-commits", Desc: "PR has commits since the last review", NeedsTimeline: true}
	XrefPushedAfterApproval = Tag{ID: "pr-pushed-after-approval", Desc: "PR was pushed to after approval", NeedsTimeline: true}
	XrefUnreviewed          = Tag{ID: "pr-unreviewed", Desc: "PR has never been reviewed", NeedsTimeline: true}

	// Review-based tags
	Approved            = Tag{ID: "approved", Desc: "Last review was an approval", NeedsReviews: true}
	ReviewedWithComment = Tag{ID: "reviewed-with-comment", Desc: "Last review was a comment", NeedsReviews: true}
	ChangesRequested    = Tag{ID: "changes-requested", Desc: "Last review was a request for changes", NeedsReviews: true}
	NewCommits          = Tag{ID: "new-commits", Desc: "PR has commits since the last review", NeedsReviews: true}
	PushedAfterApproval = Tag{ID: "pushed-after-approval", Desc: "PR was pushed to after approval", NeedsReviews: true}
	Unreviewed          = Tag{ID: "unreviewed", Desc: "PR has never been reviewed", NeedsReviews: true}

	// Special
	None = Tag{ID: "none", Desc: "No tag matched", NeedsComments: true, NeedsReviews: true, NeedsTimeline: true}
)

var Tags = []Tag{
	Assigned,
	Closed,
	OpenMilestone,
	Similar,
	Merged,
	Draft,
	Commented,
	Send,
	Recv,
	RecvQ,
	AuthorLast,
	AssigneeUpdated,
	Approved,
	ReviewedWithComment,
	ChangesRequested,
	NewCommits,
	PushedAfterApproval,
	Unreviewed,
	XrefApproved,
	XrefReviewedWithComment,
	XrefChangesRequested,
	XrefNewCommits,
	XrefPushedAfterApproval,
	XrefUnreviewed,
}

func RoleLast(role string) Tag {
	return Tag{
		ID:   fmt.Sprintf("%s-last", role),
		Desc: fmt.Sprintf("The last commenter was a project %s", role),
	}
}

func Dedup(tags []Tag) []Tag {
	deduped := []Tag{}
	seen := map[string]bool{}

	for _, t := range tags {
		if seen[t.ID] {
			continue
		}
		deduped = append(deduped, t)
		seen[t.ID] = true
	}

	return deduped
}
