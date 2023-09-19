/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package table

import (
	"errors"
	"fmt"
)

// Event table's primary key reserved 500(as is [1, 500]) ids as
// reserved id for system when event table is initialized.
// these ids are used to custom the special event.
// such as CursorReminder resource as follows.

// EventCursorReminderPrimaryID represent the primary id which represent
// the cursor reminder's record.
const EventCursorReminderPrimaryID = 1

// Event defines a resource's changes details, which is used to
// do caching operations periodically.
// an event can not be edited after created.
type Event struct {
	// ID is an auto-increased value, which is a unique identity
	// of an event.
	ID         uint32           `json:"id" gorm:"primaryKey"`
	Spec       *EventSpec       `json:"spec" gorm:"embedded"`
	Attachment *EventAttachment `json:"attachment" gorm:"embedded"`
	State      *EventState      `json:"state" gorm:"embedded"`
	Revision   *CreatedRevision `json:"revision" gorm:"embedded"`
}

// TableName is resource change event's database table name.
func (e *Event) TableName() string {
	return "events"
}

// AppID AuditRes interface
func (e *Event) AppID() uint32 {
	return 0
}

// ResID AuditRes interface
func (e *Event) ResID() uint32 {
	return e.ID
}

// ResType AuditRes interface
func (e *Event) ResType() string {
	return "event"
}

// ValidateCreate the event is valid or not when create it.
func (e *Event) ValidateCreate() error {
	if e.ID > 0 {
		return errors.New("id should not be set")
	}

	if e.Spec == nil {
		return errors.New("spec not set")
	}

	if err := e.Spec.Validate(); err != nil {
		return err
	}

	if err := e.Attachment.Validate(); err != nil {
		return err
	}

	if e.Revision == nil {
		return errors.New("revision not set")
	}

	if err := e.Revision.Validate(); err != nil {
		return err
	}

	return nil
}

// EventType is the operation type of event
type EventType string

// Validate the event type
func (e EventType) Validate() error {
	switch e {
	case InsertOp:
	case UpdateOp:
	case DeleteOp:
	default:
		return fmt.Errorf("unknown event type: %s", e)
	}

	return nil
}

const (
	// InsertOp means a resource is inserted
	InsertOp EventType = "insert"
	// UpdateOp means a resource is updated
	UpdateOp EventType = "update"
	// DeleteOp means a resource is deleted
	DeleteOp EventType = "delete"
)

// EventResource defines all the resources which can fire an event
type EventResource string

// Validate an event resource is valid or not.
func (er EventResource) Validate() error {
	switch er {
	case CursorReminder:
		return errors.New("event reminder resource is not allowed to be created")
	case Publish:
	case Application:
	case CredentialEvent:
	default:
		return fmt.Errorf("unsupported event resource: %s", er)
	}

	return nil
}

const (
	// CursorReminder is a special event resource which is used
	// to store where we have already consumed the event.
	// Note:
	// this event resource can not be generated by user, it
	// can be used only by the system itself.
	CursorReminder EventResource = "cursorReminder"
	// Publish means this is an event which represent a strategy has been published.
	Publish EventResource = "Publish"
	// Application means this is an event which represent an application resource.
	Application EventResource = "application"
	// CredentialEvent means this is an event which represent a credential resource.
	CredentialEvent EventResource = "credential"
)

// EventSpec defines the specifics of event
type EventSpec struct {
	// Resource defines what kind of resource an event belongs to
	Resource EventResource `json:"resource" gorm:"column:resource"`
	// ResourceID is the identity of this changed resource with uint32 type.
	ResourceID uint32 `json:"resource_id" gorm:"column:resource_id"`
	// ResourceUid is the identity of this changed resource with string type.
	ResourceUid string    `json:"resource_uid" gorm:"column:resource_uid"`
	OpType      EventType `json:"op_type" gorm:"column:op_type"`
}

// Validate event specifics
func (e *EventSpec) Validate() error {
	if err := e.Resource.Validate(); err != nil {
		return fmt.Errorf("validate resource failed, err: %v", err)
	}

	// at least one kind of resource id is set.
	if e.ResourceID <= 0 && len(e.ResourceUid) == 0 {
		return errors.New("invalid resource id or uid")
	}

	if err := e.OpType.Validate(); err != nil {
		return err
	}

	return nil
}

// EventAttachment is the attachment of an event.
type EventAttachment struct {
	BizID uint32 `json:"biz_id" gorm:"column:biz_id"`
	AppID uint32 `json:"app_id" gorm:"column:app_id"`
}

// Validate the event attachment is valid or not.
func (ea *EventAttachment) Validate() error {
	if ea.BizID <= 0 {
		return errors.New("invalid biz id")
	}

	return nil
}

// EventFinalStatus defines the generated event's final status type.
// Note: this status do not describe whether the event is consumed or not.
type EventFinalStatus uint

const (
	// UnknownFS means this event's related business logic transaction's
	// *Final State* is not known because of some unexpected reasons.
	// we don't know it was success or failed. This is usually caused by
	// the operation of update event status failed.
	// This is an event's default status.
	UnknownFS EventFinalStatus = 0
	// SuccessFS means the event's related business logic transaction
	// is finished with success state.
	SuccessFS EventFinalStatus = 1
	// FailedFS means the event's related business logic transaction
	// is finally finished with success state.
	FailedFS EventFinalStatus = 2
)

// EventState defines the generated event's related state infos.
type EventState struct {
	// As is known, event is inserted before the previous business logic
	// db operation is committed(such as publish a release), and it is
	// allowed that an event can be inserted success but the related
	// resource's operation can be failed or be rollback(such as the published
	// release event can be inserted success to one db sharding, but the
	// real published release can be rollback(failed) on another db sharding).
	// This status is updated after the sharding transaction is finished with
	// a success or failed state.
	FinalStatus EventFinalStatus `json:"final_status" gorm:"column:final_status"`
}
