// Copyright (c) 2017 Zededa, Inc.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/lf-edge/eve/pkg/pillar/base"
	uuid "github.com/satori/go.uuid"
)

// This is what we assume will come from the ZedControl for base OS.
// Note that we can have different versions  configured for the
// same UUID, hence the key is the UUIDandVersion  We assume the
// elements in ContentTreeConfig should be installed, but activation
// is driven by the Activate attribute.

type BaseOsConfig struct {
	UUIDandVersion        UUIDandVersion
	BaseOsVersion         string // From GetShortVersion
	ContentTreeConfigList []ContentTreeConfig
	RetryCount            int32
	Activate              bool
}

func (config BaseOsConfig) Key() string {
	return config.UUIDandVersion.UUID.String()
}

// LogCreate :
func (config BaseOsConfig) LogCreate(logBase *base.LogObject) {
	logObject := base.NewLogObject(logBase, base.BaseOsConfigLogType, config.BaseOsVersion,
		config.UUIDandVersion.UUID, config.LogKey())
	if logObject == nil {
		return
	}
	logObject.CloneAndAddField("activate", config.Activate).
		Noticef("BaseOs config create")
}

// LogModify :
func (config BaseOsConfig) LogModify(logBase *base.LogObject, old interface{}) {
	logObject := base.EnsureLogObject(logBase, base.BaseOsConfigLogType, config.BaseOsVersion,
		config.UUIDandVersion.UUID, config.LogKey())

	oldConfig, ok := old.(BaseOsConfig)
	if !ok {
		logObject.Clone().Fatalf("LogModify: Old object interface passed is not of BaseOsConfig type")
	}
	if oldConfig.Activate != config.Activate {

		logObject.CloneAndAddField("activate", config.Activate).
			AddField("old-activate", oldConfig.Activate).
			Noticef("BaseOs config modify")
	} else {
		// XXX remove?
		logObject.CloneAndAddField("diff", cmp.Diff(oldConfig, config)).
			Noticef("BaseOs config modify other change")
	}

}

// LogDelete :
func (config BaseOsConfig) LogDelete(logBase *base.LogObject) {
	logObject := base.EnsureLogObject(logBase, base.BaseOsConfigLogType, config.BaseOsVersion,
		config.UUIDandVersion.UUID, config.LogKey())
	logObject.CloneAndAddField("activate", config.Activate).
		Noticef("BaseOs config delete")

	base.DeleteLogObject(logBase, config.LogKey())
}

// LogKey :
func (config BaseOsConfig) LogKey() string {
	return string(base.BaseOsConfigLogType) + "-" + config.BaseOsVersion
}

// Indexed by UUIDandVersion as above
type BaseOsStatus struct {
	UUIDandVersion        UUIDandVersion
	BaseOsVersion         string
	Activated             bool
	TooEarly              bool // Failed since previous was inprogress/test
	ContentTreeStatusList []ContentTreeStatus
	PartitionLabel        string
	PartitionDevice       string // From zboot
	PartitionState        string // From zboot
	// Mininum state across all steps/StorageStatus.
	// Error* set implies error.
	State SwState
	// error strings across all steps/StorageStatus
	// ErrorAndTime provides SetErrorNow() and ClearError()
	ErrorAndTime
}

func (status BaseOsStatus) Key() string {
	return status.UUIDandVersion.UUID.String()
}

// LogCreate :
func (status BaseOsStatus) LogCreate(logBase *base.LogObject) {
	logObject := base.NewLogObject(logBase, base.BaseOsStatusLogType, status.BaseOsVersion,
		status.UUIDandVersion.UUID, status.LogKey())
	if logObject == nil {
		return
	}
	logObject.CloneAndAddField("state", status.State.String()).
		Noticef("BaseOs status create")
}

// LogModify :
func (status BaseOsStatus) LogModify(logBase *base.LogObject, old interface{}) {
	logObject := base.EnsureLogObject(logBase, base.BaseOsStatusLogType, status.BaseOsVersion,
		status.UUIDandVersion.UUID, status.LogKey())

	oldStatus, ok := old.(BaseOsStatus)
	if !ok {
		logObject.Clone().Fatalf("LogModify: Old object interface passed is not of BaseOsStatus type")
	}
	if oldStatus.State != status.State {

		logObject.CloneAndAddField("state", status.State.String()).
			AddField("old-state", oldStatus.State.String()).
			Noticef("BaseOs status modify")
	} else {
		// XXX remove?
		logObject.CloneAndAddField("diff", cmp.Diff(oldStatus, status)).
			Noticef("BaseOs status modify other change")
	}

	if status.HasError() {
		errAndTime := status.ErrorAndTime
		logObject.CloneAndAddField("state", status.State.String()).
			AddField("error", errAndTime.Error).
			AddField("error-time", errAndTime.ErrorTime).
			Noticef("BaseOs status modify")
	}
}

// LogDelete :
func (status BaseOsStatus) LogDelete(logBase *base.LogObject) {
	logObject := base.EnsureLogObject(logBase, base.BaseOsStatusLogType, status.BaseOsVersion,
		status.UUIDandVersion.UUID, status.LogKey())
	logObject.CloneAndAddField("state", status.State.String()).
		Noticef("BaseOs status delete")

	base.DeleteLogObject(logBase, status.LogKey())
}

// LogKey :
func (status BaseOsStatus) LogKey() string {
	return string(base.BaseOsStatusLogType) + "-" + status.BaseOsVersion
}

// return value holder
type RetStatus struct {
	Changed   bool
	MinState  SwState
	AllErrors string
	ErrorTime time.Time
}

// Mirrors proto definition for ConfigItem
// The value can be bool, float, uint, or string
type ConfigItem struct {
	Key   string
	Value interface{}
}

// Mirrors proto definition for MetricItem
// The value can be bool, float, uint, or string
type MetricItem struct {
	Key   string
	Type  MetricItemType
	Value interface{}
}

type MetricItemType uint8

const (
	MetricItemOther   MetricItemType = iota // E.g., a string like an ESSID
	MetricItemGauge                         // Goes up and down over time
	MetricItemCounter                       // Monotonically increasing (until reboot)
	MetricItemState                         // Toggles on and off; count transitions
)

type DatastoreConfig struct {
	UUID      uuid.UUID
	DsType    string
	Fqdn      string
	ApiKey    string // XXX: to be deprecated, use CipherBlockStatus instead
	Password  string // XXX: to be deprecated, use CipherBlockStatus instead
	Dpath     string // depending on DsType, it could be bucket or path
	Region    string
	DsCertPEM [][]byte // cert chain used for the datastore

	// CipherBlockStatus, for encrypted credentials
	CipherBlockStatus
}

// Key is the key in pubsub
func (config DatastoreConfig) Key() string {
	return config.UUID.String()
}

// LogCreate :
func (config DatastoreConfig) LogCreate(logBase *base.LogObject) {
	logObject := base.NewLogObject(logBase, base.DatastoreConfigLogType, "",
		config.UUID, config.LogKey())
	if logObject == nil {
		return
	}
	logObject.Noticef("Datastore config create")
}

// LogModify :
func (config DatastoreConfig) LogModify(logBase *base.LogObject, old interface{}) {
	logObject := base.EnsureLogObject(logBase, base.DatastoreConfigLogType, "",
		config.UUID, config.LogKey())

	oldConfig, ok := old.(DatastoreConfig)
	if !ok {
		logObject.Clone().Fatalf("LogModify: Old object interface passed is not of DatastoreConfig type")
	}
	// XXX remove?
	logObject.CloneAndAddField("diff", cmp.Diff(oldConfig, config)).
		Noticef("Datastore config modify")
}

// LogDelete :
func (config DatastoreConfig) LogDelete(logBase *base.LogObject) {
	logObject := base.EnsureLogObject(logBase, base.DatastoreConfigLogType, "",
		config.UUID, config.LogKey())
	logObject.Noticef("Datastore config delete")

	base.DeleteLogObject(logBase, config.LogKey())
}

// LogKey :
func (config DatastoreConfig) LogKey() string {
	return string(base.DatastoreConfigLogType) + "-" + config.Key()
}

// BootReason captures our best guess of why the device (re)booted
type BootReason uint8

// BootReasonNone is the initial value, followed by three normal reasons
// to boot/reboot, and then different error reasons
// Must match the values in api/proto/info/info.proto.BootReason
const (
	BootReasonNone BootReason = iota

	BootReasonFirst              // Normal - was not yet onboarded
	BootReasonRebootCmd          // Normal - result of a reboot command in the API
	BootReasonUpdate             // Normal - from an EVE image update in the API
	BootReasonFallback           // Fallback from a failed EVE image update
	BootReasonDisconnect         // Disconnected from controller for too long
	BootReasonFatal              // Fatal error causing log.Fatal
	BootReasonOOM                // OOM causing process to be killed
	BootReasonWatchdogHung       // Software watchdog due stuck agent
	BootReasonWatchdogPid        // Software watchdog due to e.g., golang panic
	BootReasonKernel             // TBD how we detect this
	BootReasonPowerFail          // Known power failure e.g., from disk controller S.M.A.R.T counter increase
	BootReasonUnknown            // Could be power failure, kernel panic, or hardware watchdog
	BootReasonVaultFailure       // Vault was not ready within the expected time
	BootReasonParseFail    = 255 // BootReasonFromString didn't find match
)

// String returns the string name
func (br BootReason) String() string {
	switch br {
	case BootReasonNone:
		return "BootReasonNone"
	case BootReasonFirst:
		return "BootReasonFirst"
	case BootReasonRebootCmd:
		return "BootReasonRebootCmd"
	case BootReasonUpdate:
		return "BootReasonUpdate"
	case BootReasonFallback:
		return "BootReasonFallback"
	case BootReasonDisconnect:
		return "BootReasonDisconnect"
	case BootReasonFatal:
		return "BootReasonFatal"
	case BootReasonOOM:
		return "BootReasonOOM"
	case BootReasonWatchdogHung:
		return "BootReasonWatchdogHung"
	case BootReasonWatchdogPid:
		return "BootReasonWatchdogPid"
	case BootReasonKernel:
		return "BootReasonKernel"
	case BootReasonPowerFail:
		return "BootReasonPowerFail"
	case BootReasonUnknown:
		return "BootReasonUnknown"
	case BootReasonVaultFailure:
		return "BootReasonVaultFailure"
	default:
		return fmt.Sprintf("Unknown BootReason %d", br)
	}
}

// StartWithSavedConfig indicates a normal reboot where we should immediately
// start the applications.
// Note that on most platforms we get Unknown for a power cycle
func (br BootReason) StartWithSavedConfig() bool {
	switch br {
	case BootReasonNone:
		return false
	case BootReasonFirst:
		return false
	case BootReasonRebootCmd:
		return true
	case BootReasonUpdate:
		return true
	case BootReasonFallback:
		return false
	case BootReasonDisconnect:
		return true
	case BootReasonFatal:
		return false
	case BootReasonOOM:
		return false
	case BootReasonWatchdogHung:
		return false
	case BootReasonWatchdogPid:
		return false
	case BootReasonKernel:
		return true // XXX get false Kernel for power cycle events?
	case BootReasonPowerFail:
		return true
	case BootReasonUnknown:
		return true
	case BootReasonVaultFailure:
		return false
	default:
		return false
	}
}

// BootReasonFromString parses what above String produced
// Empty string is returned as None
func BootReasonFromString(str string) BootReason {
	str = strings.TrimSuffix(str, "\n")
	str = strings.TrimSpace(str)
	switch str {
	case "", "BootReasonNone":
		return BootReasonNone
	case "BootReasonFirst":
		return BootReasonFirst
	case "BootReasonRebootCmd":
		return BootReasonRebootCmd
	case "BootReasonUpdate":
		return BootReasonUpdate
	case "BootReasonFallback":
		return BootReasonFallback
	case "BootReasonDisconnect":
		return BootReasonDisconnect
	case "BootReasonFatal":
		return BootReasonFatal
	case "BootReasonOOM":
		return BootReasonOOM
	case "BootReasonWatchdogHung":
		return BootReasonWatchdogHung
	case "BootReasonWatchdogPid":
		return BootReasonWatchdogPid
	case "BootReasonKernel":
		return BootReasonKernel
	case "BootReasonPowerFail":
		return BootReasonPowerFail
	case "BootReasonUnknown":
		return BootReasonUnknown
	case "BootReasonVaultFailure":
		return BootReasonVaultFailure
	default:
		return BootReasonParseFail
	}
}

// MaintenanceModeReason captures reason for entering into maintenance mode
type MaintenanceModeReason uint8

// MaintenanceModeReason codes for storing reason for getting into maintenance mode
const (
	MaintenanceModeReasonNone MaintenanceModeReason = iota
	MaintenanceModeReasonUserRequested
	MaintenanceModeReasonVaultLockedUp
)

// String returns the verbose equivalent of MaintenanceModeReason code
func (mmr MaintenanceModeReason) String() string {
	switch mmr {
	case MaintenanceModeReasonNone:
		return "MaintenanceModeReasonNone"
	case MaintenanceModeReasonUserRequested:
		return "MaintenanceModeReasonUserRequested"
	case MaintenanceModeReasonVaultLockedUp:
		return "MaintenanceModeReasonVaultLockedUp"
	default:
		return fmt.Sprintf("Unknown MaintenanceModeReason %d", mmr)
	}
}

// NodeAgentStatus :
type NodeAgentStatus struct {
	Name                       string
	CurPart                    string
	UpdateInprogress           bool
	RemainingTestTime          time.Duration
	DeviceReboot               bool
	RebootReason               string     // From last reboot
	BootReason                 BootReason // From last reboot
	RebootStack                string     // From last reboot
	RebootTime                 time.Time  // From last reboot
	RestartCounter             uint32
	RebootImage                string
	LocalMaintenanceMode       bool                  //enter Maintenance Mode
	LocalMaintenanceModeReason MaintenanceModeReason //reason for Maintenance Mode
}

// Key :
func (status NodeAgentStatus) Key() string {
	return status.Name
}

// LogCreate :
func (status NodeAgentStatus) LogCreate(logBase *base.LogObject) {
	logObject := base.NewLogObject(logBase, base.NodeAgentStatusLogType, status.Name,
		nilUUID, status.LogKey())
	if logObject == nil {
		return
	}
	logObject.Noticef("Nodeagent status create")
}

// LogModify :
func (status NodeAgentStatus) LogModify(logBase *base.LogObject, old interface{}) {
	logObject := base.EnsureLogObject(logBase, base.NodeAgentStatusLogType, status.Name,
		nilUUID, status.LogKey())

	oldStatus, ok := old.(NodeAgentStatus)
	if !ok {
		logObject.Clone().Fatalf("LogModify: Old object interface passed is not of NodeAgentStatus type")
	}
	// XXX remove?
	logObject.CloneAndAddField("diff", cmp.Diff(oldStatus, status)).
		Noticef("Nodeagent status modify")
}

// LogDelete :
func (status NodeAgentStatus) LogDelete(logBase *base.LogObject) {
	logObject := base.EnsureLogObject(logBase, base.NodeAgentStatusLogType, status.Name,
		nilUUID, status.LogKey())
	logObject.Noticef("Nodeagent status delete")

	base.DeleteLogObject(logBase, status.LogKey())
}

// LogKey :
func (status NodeAgentStatus) LogKey() string {
	return string(base.NodeAgentStatusLogType) + "-" + status.Key()
}

// ConfigGetStatus : Config Get Status from Controller
type ConfigGetStatus uint8

// ConfigGetSuccess : Config get is successful
const (
	ConfigGetSuccess ConfigGetStatus = iota + 1
	ConfigGetFail
	ConfigGetTemporaryFail
	ConfigGetReadSaved
)

// ZedAgentStatus :
type ZedAgentStatus struct {
	Name                 string
	ConfigGetStatus      ConfigGetStatus
	RebootCmd            bool
	RebootReason         string     // Current reason to reboot
	BootReason           BootReason // Current reason to reboot
	MaintenanceMode      bool       // Don't run apps etc
	ForceFallbackCounter int        // Try image fallback when counter changes
	CurrentProfile       string     // Current profile
}

// Key :
func (status ZedAgentStatus) Key() string {
	return status.Name
}

// LogCreate :
func (status ZedAgentStatus) LogCreate(logBase *base.LogObject) {
	logObject := base.NewLogObject(logBase, base.ZedAgentStatusLogType, status.Name,
		nilUUID, status.LogKey())
	if logObject == nil {
		return
	}
	logObject.Noticef("Zedagent status create")
}

// LogModify :
func (status ZedAgentStatus) LogModify(logBase *base.LogObject, old interface{}) {
	logObject := base.EnsureLogObject(logBase, base.ZedAgentStatusLogType, status.Name,
		nilUUID, status.LogKey())

	oldStatus, ok := old.(ZedAgentStatus)
	if !ok {
		logObject.Clone().Fatalf("LogModify: Old object interface passed is not of ZedAgentStatus type")
	}
	// XXX remove?
	logObject.CloneAndAddField("diff", cmp.Diff(oldStatus, status)).
		Noticef("Zedagent status modify")
}

// LogDelete :
func (status ZedAgentStatus) LogDelete(logBase *base.LogObject) {
	logObject := base.EnsureLogObject(logBase, base.ZedAgentStatusLogType, status.Name,
		nilUUID, status.LogKey())
	logObject.Noticef("Zedagent status delete")

	base.DeleteLogObject(logBase, status.LogKey())
}

// LogKey :
func (status ZedAgentStatus) LogKey() string {
	return string(base.ZedAgentStatusLogType) + "-" + status.Key()
}

// DeviceOpsCmd - copy of zconfig.DeviceOpsCmd
type DeviceOpsCmd struct {
	Counter      uint32
	DesiredState bool
	OpsTime      string
}

// BaseOSMgrStatus : for sending from baseosmgr
type BaseOSMgrStatus struct {
	CurrentRetryUpdateCounter uint32 // CurrentRetryUpdateCounter from baseosmgr
}

// BaseOs : copy of zconfig.BaseOS
type BaseOs struct {
	ContentTreeUUID          string
	ConfigRetryUpdateCounter uint32
}
