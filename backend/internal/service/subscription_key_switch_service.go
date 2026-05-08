package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type SubscriptionKeySwitchAction string

const (
	SubscriptionKeySwitchActionNone                 SubscriptionKeySwitchAction = "none"
	SubscriptionKeySwitchActionAutoSwitched         SubscriptionKeySwitchAction = "auto_switched"
	SubscriptionKeySwitchActionManualSwitchRequired SubscriptionKeySwitchAction = "manual_switch_required"
	SubscriptionKeySwitchActionNoAPIKeys            SubscriptionKeySwitchAction = "no_api_keys"
)

const PaymentAuditActionSubscriptionKeySwitchResult = "SUBSCRIPTION_KEY_SWITCH_RESULT"

type SubscriptionKeySwitchResult struct {
	Action              SubscriptionKeySwitchAction `json:"action"`
	TargetGroupID       int64                       `json:"target_group_id,omitempty"`
	TargetGroupName     string                      `json:"target_group_name,omitempty"`
	APIKeyCount         int                         `json:"api_key_count,omitempty"`
	AutoSwitchEnabled   bool                        `json:"auto_switch_enabled"`
	AutoSwitchedKeyID   *int64                      `json:"auto_switched_key_id,omitempty"`
	AutoSwitchedKeyName string                      `json:"auto_switched_key_name,omitempty"`
}

func (r *SubscriptionKeySwitchResult) NeedsUserNotice() bool {
	if r == nil {
		return false
	}
	return r.Action != SubscriptionKeySwitchActionNone
}

type SubscriptionKeySwitchCoordinator struct {
	apiKeyRepo           APIKeyRepository
	groupRepo            GroupRepository
	settingRepo          SettingRepository
	authCacheInvalidator APIKeyAuthCacheInvalidator
}

func NewSubscriptionKeySwitchCoordinator(
	apiKeyRepo APIKeyRepository,
	groupRepo GroupRepository,
	settingRepo SettingRepository,
	authCacheInvalidator APIKeyAuthCacheInvalidator,
) *SubscriptionKeySwitchCoordinator {
	return &SubscriptionKeySwitchCoordinator{
		apiKeyRepo:           apiKeyRepo,
		groupRepo:            groupRepo,
		settingRepo:          settingRepo,
		authCacheInvalidator: authCacheInvalidator,
	}
}

func (s *SubscriptionKeySwitchCoordinator) HandleNewSubscription(ctx context.Context, userID, groupID int64) *SubscriptionKeySwitchResult {
	return s.evaluate(ctx, userID, groupID, true)
}

func (s *SubscriptionKeySwitchCoordinator) DescribeCurrentState(ctx context.Context, userID, groupID int64) *SubscriptionKeySwitchResult {
	return s.evaluate(ctx, userID, groupID, false)
}

func (s *SubscriptionKeySwitchCoordinator) evaluate(ctx context.Context, userID, groupID int64, allowMutation bool) *SubscriptionKeySwitchResult {
	result := &SubscriptionKeySwitchResult{
		Action:        SubscriptionKeySwitchActionNone,
		TargetGroupID: groupID,
	}

	if s == nil || userID <= 0 || groupID <= 0 {
		return result
	}

	if s.groupRepo != nil {
		if group, err := s.groupRepo.GetByID(ctx, groupID); err == nil && group != nil {
			result.TargetGroupName = group.Name
		}
	}

	result.AutoSwitchEnabled = s.subscriptionAutoSwitchSingleKeyEnabled(ctx)
	if s.apiKeyRepo == nil {
		return result
	}

	keys, _, err := s.apiKeyRepo.ListByUserID(ctx, userID, pagination.PaginationParams{
		Page:     1,
		PageSize: 1000,
	}, APIKeyListFilters{})
	if err != nil {
		slog.Warn("subscription_key_switch_list_keys_failed", "user_id", userID, "group_id", groupID, "error", err)
		return result
	}
	result.APIKeyCount = len(keys)

	if len(keys) == 0 {
		result.Action = SubscriptionKeySwitchActionNoAPIKeys
		return result
	}

	allKeysAlreadyOnTarget := true
	for i := range keys {
		if keys[i].GroupID == nil || *keys[i].GroupID != groupID {
			allKeysAlreadyOnTarget = false
			break
		}
	}
	if allKeysAlreadyOnTarget {
		return result
	}

	switch len(keys) {
	case 1:
		if !result.AutoSwitchEnabled || !allowMutation {
			result.Action = SubscriptionKeySwitchActionManualSwitchRequired
			return result
		}

		key := keys[0]
		if allowMutation {
			key.GroupID = &groupID
			key.Group = nil
			if err := s.apiKeyRepo.Update(ctx, &key); err != nil {
				slog.Warn("subscription_key_switch_auto_update_failed", "user_id", userID, "key_id", key.ID, "group_id", groupID, "error", err)
				result.Action = SubscriptionKeySwitchActionManualSwitchRequired
				return result
			}
			if s.authCacheInvalidator != nil {
				s.authCacheInvalidator.InvalidateAuthCacheByKey(ctx, key.Key)
			}
		}
		result.Action = SubscriptionKeySwitchActionAutoSwitched
		result.AutoSwitchedKeyID = &key.ID
		result.AutoSwitchedKeyName = strings.TrimSpace(key.Name)
		return result
	default:
		result.Action = SubscriptionKeySwitchActionManualSwitchRequired
		return result
	}
}

func (s *SubscriptionKeySwitchCoordinator) subscriptionAutoSwitchSingleKeyEnabled(ctx context.Context) bool {
	if s == nil || s.settingRepo == nil {
		return false
	}
	value, err := s.settingRepo.GetValue(ctx, SettingKeySubscriptionAutoSwitchSingleKey)
	if err != nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(value), "true")
}
