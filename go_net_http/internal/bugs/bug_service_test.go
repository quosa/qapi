package bugs_test

import (
	"testing"

	"github.com/quosa/qapi/internal/bugs"
)

func TestNewBugService(t *testing.T) {
	var bs bugs.IBugService

	if bs = bugs.NewBugService(); bs == nil {
		t.Error("NewBugService() should not return nil")
	}
	bugs, err := bs.GetAllBugs()
	if err != nil {
		t.Errorf("GetAllBugs() should not return error, got %v", err)
	}
	if len(bugs) != 0 {
		t.Errorf("GetAllBugs() should return empty list, got %v", bugs)
	}
}

func TestAddBug(t *testing.T) {
	var bs bugs.IBugService

	if bs = bugs.NewBugService(); bs == nil {
		t.Error("NewBugService() should not return nil")
	}
	bug := bugs.Bug{
		Title:       "Test Bug",
		Description: "Test Description",
		Status:      "New",
	}
	newBug, err := bs.CreateBug(bug)
	if err != nil {
		t.Errorf("CreateBug() should not return error, got %v", err)
	}
	if newBug.ID == 0 {
		t.Error("CreateBug() should not return zero ID")
	}
	if newBug.Title != bug.Title {
		t.Errorf("CreateBug() should return bug with title %s, got %s", bug.Title, newBug.Title)
	}
	if newBug.Description != bug.Description {
		t.Errorf("CreateBug() should return bug with description %s, got %s", bug.Description, newBug.Description)
	}
	if newBug.Status != bug.Status {
		t.Errorf("CreateBug() should return bug with status %s, got %s", bug.Status, newBug.Status)
	}

	bugs, err := bs.GetAllBugs()
	if err != nil {
		t.Errorf("GetAllBugs() should not return error, got %v", err)
	}
	if len(bugs) != 1 {
		t.Errorf("GetAllBugs() should return empty list, got %v", bugs)
	}

	gotBug, err := bs.GetBugByID(newBug.ID)
	if err != nil {
		t.Errorf("GetBugByID() should not return error, got %v", err)
	}
	if gotBug.ID == 0 {
		t.Error("GetBugByID() should not return zero ID")
	}
	if gotBug.Title != bug.Title {
		t.Errorf("GetBugByID() should return bug with title %s, got %s", bug.Title, gotBug.Title)
	}
	if gotBug.Description != bug.Description {
		t.Errorf("GetBugByID() should return bug with description %s, got %s", bug.Description, gotBug.Description)
	}
}

func TestBugCycle(t *testing.T) {
	bs := bugs.NewBugService()

	// CREATE a bug
	test_bug := bugs.Bug{
		Title:       "Test Bug",
		Description: "Test Description",
		Status:      "New",
	}
	b_created, err := bs.CreateBug(test_bug)
	if err != nil {
		t.Errorf("CreateBug() should not return error, got %v", err)
	}
	b_got, err := bs.GetBugByID(b_created.ID)
	if err != nil {
		t.Errorf("GetBugByID() should not return error, got %v", err)
	}
	if b_got.ID != b_created.ID {
		t.Errorf("GetBugByID() should return created bug with ID %d, got %d", b_created.ID, b_got.ID)
	}

	// UPDATE the bug
	b_got.Title = "Test Bug Updated"
	b_updated, err := bs.UpdateBug(b_got)
	if err != nil {
		t.Errorf("UpdateBug() should not return error, got %v", err)
	}
	if b_updated.ID != b_got.ID {
		t.Errorf("UpdateBug() should return bug with ID %d, got %d", b_created.ID, b_updated.ID)
	}
	b_updated_got, err := bs.GetBugByID(b_created.ID)
	if err != nil {
		t.Errorf("GetBugByID() should not return error, got %v", err)
	}
	if b_updated_got.ID != b_created.ID {
		t.Errorf("GetBugByID() should return updated bug with ID %d, got %d", b_created.ID, b_updated_got.ID)
	}

	// DELETE the bug
	err = bs.DeleteBug(b_updated.ID)
	if err != nil {
		t.Errorf("DeleteBug() should not return error, got %v", err)

	}
	b_deleted_got, err := bs.GetBugByID(b_updated.ID)
	if err == nil {
		t.Errorf("GetBugByID() should return error, got %v", err)
	}
	if b_deleted_got.ID != 0 {
		t.Errorf("GetBugByID() should return empty Bug with error, got %v", b_deleted_got)
	}
}

func TestDoubleAddBug(t *testing.T) {
	var bs bugs.IBugService

	if bs = bugs.NewBugService(); bs == nil {
		t.Error("NewBugService() should not return nil")
	}
	bug := bugs.Bug{
		Title:       "Test Bug",
		Description: "Test Description",
		Status:      "New",
	}
	newBug1, err := bs.CreateBug(bug)
	if err != nil {
		t.Errorf("CreateBug() should not return error, got %v", err)
	}
	newBug2, err := bs.CreateBug(bug)
	if err != nil {
		t.Errorf("CreateBug() should not return error, got %v", err)
	}
	if newBug1.ID == newBug2.ID {
		t.Errorf("CreateBug() should return different IDs, got %d and %d", newBug1.ID, newBug2.ID)
	}
}

func TestGetNonexistingBug(t *testing.T) {
	var bs bugs.IBugService
	if bs = bugs.NewBugService(); bs == nil {
		t.Error("NewBugService() should not return nil")
	}
	_, err := bs.GetBugByID(123)
	if err == nil {
		t.Error("GetBugByID() should return error for non-existing bug")
	}
}

func TestUpdateNonexistingBug(t *testing.T) {
	var bs bugs.IBugService
	if bs = bugs.NewBugService(); bs == nil {
		t.Error("NewBugService() should not return nil")
	}
	_, err := bs.UpdateBug(bugs.Bug{ID: 123})
	if err == nil {
		t.Error("GetBugByID() should return error for non-existing bug")
	}
}

func TestDeleteNonexistingBug(t *testing.T) {
	var bs bugs.IBugService
	if bs = bugs.NewBugService(); bs == nil {
		t.Error("NewBugService() should not return nil")
	}
	err := bs.DeleteBug(123)
	if err == nil {
		t.Error("GetBugByID() should return error for non-existing bug")
	}
}
