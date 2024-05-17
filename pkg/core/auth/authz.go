package auth

import (
	"github.com/samber/lo"
)

// Role ..
type Role string

// ToString ..
func (u Role) ToString() string {
	return string(u)
}

func ValidRole(role Role) bool {
	return lo.Contains(_validRoles, role)
}

// roles ..
const (
	RoleAdmin = Role("admin")
	RoleUser  = Role("user")
)

var _validRoles = []Role{
	RoleAdmin,
	RoleUser,
}

const _ok = true

type Permission int

const (
	CreateAssignment Permission = iota
	UpdateAssignment
	ViewAssignment
	ViewAnyAssignments
	DeleteAssignment
	GradeAssignment
	ViewAnySubmissions

	CreateSubmission
	CreateSubmissionForOther
	UpdateSubmission
	ViewSubmission
	DeleteSubmission
	DeleteSubmissionForOther

	ViewAnyUsers
	CreateAnyUser
	UpdateUser
	CreateUser

	CreateMedia

	ViewCourse
	CreateCourse
	UpdateCourse
	ViewAdminCoursesStudents
	ViewAdminCourseDetail
	ViewStudentEnrolledCourses
	EnrollStudentCourse
)

var policy = map[Role]map[Permission]bool{
	RoleAdmin: {
		ViewAnyUsers: _ok,
	},
	RoleUser: {
		ViewAssignment:             _ok,
		UpdateSubmission:           _ok,
		DeleteSubmission:           _ok,
		UpdateUser:                 _ok,
		CreateMedia:                _ok,
		CreateSubmission:           _ok,
		ViewCourse:                 _ok,
		ViewStudentEnrolledCourses: _ok,
		EnrollStudentCourse:        _ok,
	},
}

// Granted check if role is granted with a permission
func (r Role) Granted(perm Permission) bool {
	role, ok := policy[r]
	if !ok {
		return false
	}

	return role[perm]
}

func (r Role) Can(perms ...Permission) bool {
	for _, perm := range perms {
		if !r.Granted(perm) {
			return false
		}
	}

	return true
}
