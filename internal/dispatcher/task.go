package dispatcher

import (
	"bytes"
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
	vlib "github.com/tsenart/vegeta/lib"

	"io"
	"io/ioutil"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"
)

// ITask defines an interface for attack tasks
type ITask interface {
	ITaskGetter
	ITaskActions
}

type ITaskGetter interface {
	// ID returns the attack task ID
	ID() string
	// Status returns the attack task status
	Status() models.AttackStatus
	// Params returns the attack task params
	Params() models.AttackParams
	// CreatedAt returns the created at timestamp
	CreatedAt() time.Time
	// UpdatedAt returns the updated at timestamp
	UpdatedAt() time.Time
	// Result returns the result as a byte array
	Result() io.Reader
}

type ITaskActions interface {
	// Run the attack using the configured attack function.
	Run(vegeta.AttackFunc) error
	// Complete changes task status to completed
	Complete(io.Reader) error
	// Cancel changes task status to canceled
	Cancel() error
	// Fail changes task status to failed
	Fail() error
	// SendUpdate sends an update on the update chan to the caller
	SendUpdate()
}

type UpdateMessage struct {
	ID     string
	Status models.AttackStatus
}

type attackContext struct {
	context.Context
	cancelFn context.CancelFunc
}

func newAttackContext() attackContext {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return attackContext{
		ctx,
		cancel,
	}
}

type task struct {
	ctx attackContext

	id     string
	params models.AttackParams
	status models.AttackStatus
	result *bytes.Buffer

	createdAt time.Time
	updatedAt time.Time

	updateCh chan UpdateMessage
}

// NewTask returns a new instance of a task object
func NewTask(updateCh chan UpdateMessage, params models.AttackParams) *task { //nolint: golint
	id := uuid.NewV4().String()
	t := &task{
		newAttackContext(),

		id,
		params,
		models.AttackResponseStatusScheduled,
		bytes.NewBuffer(make([]byte, 0)),

		time.Now(),
		time.Now(),

		updateCh,
	}

	t.log(nil).Debug("creating new task")

	return t
}

// Run an attack task using the passed in attack function
func (t *task) Run(fn vegeta.AttackFunc) error {
	if t.status != models.AttackResponseStatusScheduled {
		return fmt.Errorf("cannot run task %s with status %s", t.id, t.status)
	}

	t.log(nil).Debug("running")

	go run(t, fn) //nolint: errcheck

	t.status = models.AttackResponseStatusRunning

	t.SendUpdate()

	return nil
}

// Complete marks a task as completed
func (t *task) Complete(result io.Reader) error {
	if t.status != models.AttackResponseStatusRunning {
		return fmt.Errorf("cannot mark completed for task %s with status %s", t.id, t.status)
	}

	buf, err := ioutil.ReadAll(result)
	if err != nil {
		return err
	}

	t.status = models.AttackResponseStatusCompleted
	t.result = bytes.NewBuffer(buf)

	t.SendUpdate()

	t.log(nil).Debug("completed")

	return nil
}

// Cancel invokes the context cancel and marks a task as canceled
func (t *task) Cancel() error {
	if t.status == models.AttackResponseStatusCompleted || t.status == models.AttackResponseStatusFailed || t.status == models.AttackResponseStatusCanceled { // nolint: lll
		return fmt.Errorf("cannot cancel task %s with status %s", t.id, t.status)
	}

	t.ctx.cancelFn()

	t.status = models.AttackResponseStatusCanceled

	t.SendUpdate()

	t.log(nil).Debug("canceled")

	return nil
}

// Fail marks a task as failed
func (t *task) Fail() error {
	t.status = models.AttackResponseStatusFailed

	t.SendUpdate()

	t.log(nil).Error("failed")
	return nil
}

// SendUpdate to send a status update on the update channel
func (t *task) SendUpdate() {
	t.updatedAt = time.Now()
	t.updateCh <- UpdateMessage{
		t.id,
		t.status,
	}
}

// ID returns the task identifier
func (t *task) ID() string {
	return t.id
}

// Status returns the latest task status
func (t *task) Status() models.AttackStatus {
	return t.status
}

// Params returns a the configured attack params
func (t *task) Params() models.AttackParams {
	return t.params
}

// CreatedAt returns the created at timestamp
func (t *task) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the created at timestamp
func (t *task) UpdatedAt() time.Time {
	return t.updatedAt
}

// Result returns the result as a io.Reader
func (t *task) Result() io.Reader {
	return t.result
}

// TODO: Remove dependency on vegeta lib. Move functionality to pkg/vegeta package.
func run(t *task, fn vegeta.AttackFunc) error {
	opts, err := vegeta.NewAttackOptsFromAttackParams(t.id, t.params)
	if err != nil {
		return err
	}

	result := fn(opts)
	if result == nil {
		return fmt.Errorf("empty channel returned")
	}

	buf := bytes.NewBuffer(nil)
	enc := vlib.NewEncoder(buf)
loop:
	for {
		select {
		case r, ok := <-result:
			if !ok {
				break loop
			}
			if err = enc.Encode(r); err != nil {
				_ = t.Fail()
			}
		case <-t.ctx.Done():
			t.log(nil).Warnf("task %s was canceled", t.id)
			return nil
		}
	}

	// Mark attack as completed
	err = t.Complete(buf)
	if err != nil {
		log.WithError(err).Error("Failed to Complete")
		_ = t.Fail()
	}

	return nil
}

func (t *task) log(fields map[string]interface{}) *log.Entry {
	l := log.WithField("component", "task")

	l = l.WithFields(log.Fields{
		"ID":     t.id,
		"Status": t.status,
	})

	if fields != nil {
		l = l.WithFields(fields)
	}
	return l
}

func attackDetailFromTask(t ITaskGetter) models.AttackDetails {
	details := models.AttackDetails{
		AttackInfo: models.AttackInfo{
			ID:        t.ID(),
			Status:    t.Status(),
			Params:    t.Params(),
			CreatedAt: t.CreatedAt().Format(time.RFC1123),
			UpdatedAt: t.UpdatedAt().Format(time.RFC1123),
		},
	}

	if t.Status() == models.AttackResponseStatusCompleted {
		result := t.Result()
		buf, _ := ioutil.ReadAll(result)
		details.Result = buf
	}

	return details
}
