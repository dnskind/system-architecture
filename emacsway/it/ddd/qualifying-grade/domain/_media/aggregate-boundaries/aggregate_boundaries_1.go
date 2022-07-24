package qualifying_grade_1

import (
    "errors"
    "time"
)

type EndorsedId uint64
type MemberId uint64
type Grade uint
type AvailableEndorsementCount uint
type ReceivedEndorsementCount uint
type RecognizerId uint64
type EndorsementId uint64
type ArtifactDescription string

type Weight uint8

const (
    PeerWeight   = 1
    HigherWeight = 2
)

const (
    WithoutGrade = Grade(0)
    Grade3       = Grade(1)
    Grade2       = Grade(2)
    Grade1       = Grade(3)
    Candidate    = Grade(4)
    Expert       = Grade(5)
)

type Endorsed struct {
    id                       EndorsedId
    memberId                 MemberId
    grade                    Grade
    receivedEndorsementCount ReceivedEndorsementCount
    gradeLogEntries          []GradeLogEntry
    version                  uint
    createdAt                time.Time
}

func (e Endorsed) GetId() EndorsedId {
    return e.id
}

func (e Endorsed) GetGrade() Grade {
    return e.grade
}

func (e Endorsed) GetVersion() uint {
    return e.version
}

func (e *Endorsed) IncreaseReceivedEndorsementCount(w Weight, dt time.Time) {
    e.receivedEndorsementCount += ReceivedEndorsementCount(w)
    if e.grade == WithoutGrade && e.receivedEndorsementCount >= 6 {
        e.setGrade(Grade3, dt)
        e.receivedEndorsementCount = 0
    } else if e.grade == Grade3 && e.receivedEndorsementCount >= 10 {
        e.setGrade(Grade2, dt)
        e.receivedEndorsementCount = 0
    } else if e.grade == Grade2 && e.receivedEndorsementCount >= 14 {
        e.setGrade(Grade1, dt)
        e.receivedEndorsementCount = 0
    } else if e.grade == Grade1 && e.receivedEndorsementCount >= 20 {
        e.setGrade(Candidate, dt)
        e.receivedEndorsementCount = 0
    } else if e.grade == Candidate && e.receivedEndorsementCount >= 10 {
        e.setGrade(Expert, dt)
        e.receivedEndorsementCount = 0
    }
}

func (e *Endorsed) setGrade(g Grade, dt time.Time) {
    e.gradeLogEntries = append(e.gradeLogEntries, GradeLogEntry{
        e.id, e.version, g, dt,
    })
    e.grade = g
}

func (e *Endorsed) IncreaseVersion() {
    e.version += 1
}

type GradeLogEntry struct {
    endorsedId          EndorsedId
    endorsedVersion     uint
    assignedGrade       Grade
    createdAt           time.Time
}

type Recognizer struct {
    id                        RecognizerId
    memberId                  MemberId
    grade                     Grade
    availableEndorsementCount AvailableEndorsementCount
    version                   uint
    createdAt                 time.Time
}

func (r Recognizer) Endorse(e Endorsed, desc ArtifactDescription, t time.Time) (Endorsement, error) {
    if r.grade < e.grade {
        return Endorsement{}, errors.New(
            "it is allowed to endorse only members with equal or lower grade",
        )
    }
    if r.availableEndorsementCount == 0 {
        return Endorsement{}, errors.New(
            "you have reached the limit of available recommendations this year",
        )
    }
    return Endorsement{
        0, r.id, r.grade, r.version,
        e.id, e.grade, e.GetVersion(),
        desc, t,
    }, nil
}

func (r *Recognizer) DecreaseAvailableEndorsementCount() error {
    if r.availableEndorsementCount == 0 {
        return errors.New("no endorsement is available")
    }
    r.availableEndorsementCount -= 1
    return nil
}

func (r *Recognizer) IncreaseVersion() {
    r.version += 1
}

type Endorsement struct {
    id                  EndorsementId
    recognizerId        RecognizerId
    recognizerGrade     Grade
    recognizerVersion   uint
    endorsedId          EndorsedId
    endorsedGrade       Grade
    endorsedVersion     uint
    artifactDescription ArtifactDescription
    createdAt           time.Time
}
