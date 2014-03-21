package lpad

import (
	"strconv"
	"strings"
)

// A BugStub holds details necessary for creating a new bug in Launchpad.
type BugStub struct {
	Title           string   // Required
	Description     string   // Required
	Target          AnyValue // Project, source package, or distribution
	Private         bool
	SecurityRelated bool
	Tags            []string
}

// CreateBug creates a new bug with an appropriate bug task and returns it.
func (root *Root) Bug(id int) (*Bug, error) {
	v, err := root.Location("/bugs/" + strconv.Itoa(id)).Get(nil)
	if err != nil {
		return nil, err
	}
	return &Bug{v}, nil
}

// CreateBug creates a new bug with an appropriate bug task and returns it.
func (root *Root) CreateBug(stub *BugStub) (*Bug, error) {
	params := Params{
		"ws.op":       "createBug",
		"title":       stub.Title,
		"description": stub.Description,
		"target":      stub.Target.AbsLoc(),
	}
	if len(stub.Tags) > 0 {
		params["tags"] = strings.Join(stub.Tags, " ")
	}
	if stub.Private {
		params["private"] = "true"
	}
	if stub.SecurityRelated {
		params["security_related"] = "true"
	}
	v, err := root.Location("/bugs").Post(params)
	if err != nil {
		return nil, err
	}
	return &Bug{v}, nil
}

// The Bug type represents a bug in Launchpad.
type Bug struct {
	*Value
}

// Id returns the bug numeric identifier (the bug # itself).
func (bug *Bug) Id() int {
	return bug.IntField("id")
}

// Title returns the short bug summary.
func (bug *Bug) Title() string {
	return bug.StringField("title")
}

// Description returns the main bug description.
func (bug *Bug) Description() string {
	return bug.StringField("description")
}

// Tags returns the set of tags associated with the bug.
func (bug *Bug) Tags() []string {
	return bug.StringListField("tags")
}

// Private returns true if the bug is flagged as private.
func (bug *Bug) Private() bool {
	return bug.BoolField("private")
}

// SecurityRelated returns true if the bug describes sensitive
// information about a security vulnerability.
func (bug *Bug) SecurityRelated() bool {
	return bug.BoolField("security_related")
}

// WebPage returns the URL for accessing this bug in a browser.
func (bug *Bug) WebPage() string {
	return bug.StringField("web_link")
}

// SetTitle changes the bug title.
// Patch must be called to commit all changes.
func (bug *Bug) SetTitle(title string) {
	bug.SetField("title", title)
}

// SetDescription changes the bug description.
// Patch must be called to commit all changes.
func (bug *Bug) SetDescription(description string) {
	bug.SetField("description", description)
}

// SetTags changes the bug tags.
// Patch must be called to commit all changes.
func (bug *Bug) SetTags(tags []string) {
	bug.SetField("tags", tags)
}

// SetPrivate changes the bug private flag.
// Patch must be called to commit all changes.
func (bug *Bug) SetPrivate(private bool) {
	bug.SetField("private", private)
}

// SetSecurityRelated sets to related the flag that tells if
// a bug is security sensitive or not.
// Patch must be called to commit all changes.
func (bug *Bug) SetSecurityRelated(related bool) {
	bug.SetField("security_related", related)
}

// LinkBranch associates a branch with this bug.
func (bug *Bug) LinkBranch(branch *Branch) error {
	params := Params{
		"ws.op":  "linkBranch",
		"branch": branch.AbsLoc(),
	}
	_, err := bug.Post(params)
	return err
}

func (bug *Bug) Owner() (*Person, error) {
	v, err := bug.Link("owner_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &Person{v}, nil
}

//
func (bug *Bug) DateCreated() string {
	return bug.StringField("date_created")
}

//
func (bug *Bug) DateLastMessage() string {
	return bug.StringField("date_last_message")
}

//
func (bug *Bug) DateLastUpdated() string {
	return bug.StringField("date_last_updated")
}

//
func (bug *Bug) DateLatestPatchUploaded() string {
	return bug.StringField("date_latest_patch_uploaded")
}

//
func (bug *Bug) IsExpirable() bool {
	return bug.BoolField("is_expirable")
}

//
func (bug *Bug) Heat() int {
	return bug.IntField("heat")
}

// A BugTask represents the association of a bug with a project
// or source package, and the related information.
type BugTask struct {
	*Value
}

type BugImportance string

const (
	ImUnknown   BugImportance = "Unknown"
	ImCritical  BugImportance = "Critical"
	ImHigh      BugImportance = "High"
	ImMedium    BugImportance = "Medium"
	ImLow       BugImportance = "Low"
	ImWishlist  BugImportance = "Wishlist"
	ImUndecided BugImportance = "Undecided"
)

type BugStatus string

const (
	StUnknown      BugStatus = "Unknown"
	StNew          BugStatus = "New"
	StIncomplete   BugStatus = "Incomplete"
	StOpinion      BugStatus = "Opinion"
	StInvalid      BugStatus = "Invalid"
	StWontFix      BugStatus = "Won't fix"
	StExpired      BugStatus = "Expired"
	StConfirmed    BugStatus = "Confirmed"
	StTriaged      BugStatus = "Triaged"
	StInProgress   BugStatus = "In Progress"
	StFixCommitted BugStatus = "Fix Committed"
	StFixReleased  BugStatus = "Fix Released"
)

// Status returns the current status for the bug task. See
// the Status type for supported values.
func (task *BugTask) Status() BugStatus {
	return BugStatus(task.StringField("status"))
}

// Importance returns the current importance for the bug task. See
// the Importance type for supported values.
func (task *BugTask) Importance() BugImportance {
	return BugImportance(task.StringField("importance"))
}

// Assignee returns the person currently assigned to work on the task.
func (task *BugTask) Assignee() (*Person, error) {
	v, err := task.Link("assignee_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &Person{v}, nil
}

//
func (task *BugTask) BugTargetName() string {
	return task.StringField("bug_target_name")
}

//
func (task *BugTask) BugTargetDisplayName() string {
	return task.StringField("bug_target_display_name")
}

// Milestone returns the milestone the task is currently targeted at.
func (task *BugTask) Milestone() (*Milestone, error) {
	v, err := task.Link("milestone_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &Milestone{v}, nil
}

//
func (task *BugTask) DateCreated() string {
	return task.StringField("date_created")
}

//
func (task *BugTask) DateConfirmed() string {
	return task.StringField("date_confirmed")
}

//
func (task *BugTask) DateAssigned() string {
	return task.StringField("date_assigned")
}

//
func (task *BugTask) DateClosed() string {
	return task.StringField("date_closed")
}

//
func (task *BugTask) DateFixCommitted() string {
	return task.StringField("date_fix_committed")
}

//
func (task *BugTask) DateFixReleased() string {
	return task.StringField("date_fix_released")
}

//
func (task *BugTask) DateInProgress() string {
	return task.StringField("date_in_progress")
}

//
func (task *BugTask) DateIncomplete() string {
	return task.StringField("date_incomplete")
}

//
func (task *BugTask) DateLeftClosed() string {
	return task.StringField("date_left_closed")
}

//
func (task *BugTask) DateLeftNew() string {
	return task.StringField("date_left_new")
}

//
func (task *BugTask) DateTriaged() string {
	return task.StringField("date_triaged")
}

//
func (task *BugTask) IsComplete() bool {
	return task.BoolField("date_is_complete")
}

//
func (task *BugTask) Owner() (*Person, error) {
	v, err := task.Link("owner_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &Person{v}, nil
}

//
func (task *BugTask) Title() string {
	return task.StringField("title")
}

// SetStatus changes the current status for the bug task. See
// the Status type for supported values.
func (task *BugTask) SetStatus(status BugStatus) {
	task.SetField("status", string(status))
}

// Importance changes the current importance for the bug task. See
// the Importance type for supported values.
func (task *BugTask) SetImportance(importance BugImportance) {
	task.SetField("importance", string(importance))
}

// SetAssignee changes the person currently assigned to work on the task.
func (task *BugTask) SetAssignee(person *Person) {
	task.SetField("assignee_link", person.AbsLoc())
}

// SetMilestone changes the milestone the task is currently targeted at.
func (task *BugTask) SetMilestone(ms *Milestone) {
	task.SetField("milestone_link", ms.AbsLoc())
}

// BugTaskList represents a list of BugTasks for iteration.
type BugTaskList struct {
	*Value
}

// For iterates over the list of bug tasks and calls f for each one.
// If f returns a non-nil error, iteration will stop and the error will
// be returned as the result of For.
func (list *BugTaskList) For(f func(bt *BugTask) error) error {
	return list.Value.For(func(v *Value) error {
		f(&BugTask{v})
		return nil
	})
}

// Tasks returns the list of bug tasks associated with the bug.
func (bug *Bug) Tasks() (*BugTaskList, error) {
	v, err := bug.Link("bug_tasks_collection_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &BugTaskList{v}, nil
}
