package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type subscriptionKeySwitchAPIKeyRepoStub struct {
	APIKeyRepository
	keys        []APIKey
	updatedKeys []APIKey
}

func (s *subscriptionKeySwitchAPIKeyRepoStub) ListByUserID(context.Context, int64, pagination.PaginationParams, APIKeyListFilters) ([]APIKey, *pagination.PaginationResult, error) {
	return append([]APIKey(nil), s.keys...), &pagination.PaginationResult{Total: int64(len(s.keys)), Page: 1, PageSize: len(s.keys), Pages: 1}, nil
}

func (s *subscriptionKeySwitchAPIKeyRepoStub) Update(_ context.Context, key *APIKey) error {
	s.updatedKeys = append(s.updatedKeys, *key)
	for i := range s.keys {
		if s.keys[i].ID == key.ID {
			s.keys[i] = *key
		}
	}
	return nil
}

type subscriptionKeySwitchGroupRepoStub struct {
	GroupRepository
	group *Group
}

func (s *subscriptionKeySwitchGroupRepoStub) GetByID(context.Context, int64) (*Group, error) {
	return s.group, nil
}

type subscriptionKeySwitchSettingRepoStub struct {
	SettingRepository
	value string
}

func (s *subscriptionKeySwitchSettingRepoStub) GetValue(context.Context, string) (string, error) {
	return s.value, nil
}

type subscriptionKeySwitchInvalidatorStub struct {
	calls []string
}

func (s *subscriptionKeySwitchInvalidatorStub) InvalidateAuthCacheByKey(_ context.Context, key string) {
	s.calls = append(s.calls, key)
}

func (s *subscriptionKeySwitchInvalidatorStub) InvalidateAuthCacheByUserID(context.Context, int64) {}

func (s *subscriptionKeySwitchInvalidatorStub) InvalidateAuthCacheByGroupID(context.Context, int64) {}

func TestSubscriptionKeySwitchCoordinatorAutoSwitchesSingleKey(t *testing.T) {
	keyGroupID := int64(1)
	apiKeyRepo := &subscriptionKeySwitchAPIKeyRepoStub{
		keys: []APIKey{{ID: 9, Key: "sk-test", Name: "default", GroupID: &keyGroupID}},
	}
	invalidator := &subscriptionKeySwitchInvalidatorStub{}
	coordinator := NewSubscriptionKeySwitchCoordinator(
		apiKeyRepo,
		&subscriptionKeySwitchGroupRepoStub{group: &Group{ID: 42, Name: "Claude Pro", SubscriptionType: SubscriptionTypeSubscription}},
		&subscriptionKeySwitchSettingRepoStub{value: "true"},
		invalidator,
	)

	result := coordinator.HandleNewSubscription(context.Background(), 7, 42)
	require.NotNil(t, result)
	require.Equal(t, SubscriptionKeySwitchActionAutoSwitched, result.Action)
	require.Equal(t, int64(42), result.TargetGroupID)
	require.Equal(t, "Claude Pro", result.TargetGroupName)
	require.Equal(t, 1, result.APIKeyCount)
	require.NotNil(t, result.AutoSwitchedKeyID)
	require.Equal(t, int64(9), *result.AutoSwitchedKeyID)
	require.Equal(t, "default", result.AutoSwitchedKeyName)
	require.Len(t, apiKeyRepo.updatedKeys, 1)
	require.NotNil(t, apiKeyRepo.updatedKeys[0].GroupID)
	require.Equal(t, int64(42), *apiKeyRepo.updatedKeys[0].GroupID)
	require.Equal(t, []string{"sk-test"}, invalidator.calls)
}

func TestSubscriptionKeySwitchCoordinatorReturnsManualSwitchForMultipleKeys(t *testing.T) {
	groupOne := int64(1)
	groupTwo := int64(2)
	apiKeyRepo := &subscriptionKeySwitchAPIKeyRepoStub{
		keys: []APIKey{
			{ID: 1, Key: "sk-1", Name: "k1", GroupID: &groupOne},
			{ID: 2, Key: "sk-2", Name: "k2", GroupID: &groupTwo},
		},
	}
	coordinator := NewSubscriptionKeySwitchCoordinator(
		apiKeyRepo,
		&subscriptionKeySwitchGroupRepoStub{group: &Group{ID: 42, Name: "GPT Sub", SubscriptionType: SubscriptionTypeSubscription}},
		&subscriptionKeySwitchSettingRepoStub{value: "true"},
		nil,
	)

	result := coordinator.HandleNewSubscription(context.Background(), 7, 42)
	require.Equal(t, SubscriptionKeySwitchActionManualSwitchRequired, result.Action)
	require.Equal(t, 2, result.APIKeyCount)
	require.Empty(t, apiKeyRepo.updatedKeys)
}

func TestSubscriptionKeySwitchCoordinatorStillWarnsWhenOnlySomeKeysAlreadyUseTargetGroup(t *testing.T) {
	targetGroupID := int64(42)
	apiKeyRepo := &subscriptionKeySwitchAPIKeyRepoStub{
		keys: []APIKey{
			{ID: 1, Key: "sk-1", Name: "already-right", GroupID: &targetGroupID},
			{ID: 2, Key: "sk-2", Name: "stale", GroupID: nil},
		},
	}
	coordinator := NewSubscriptionKeySwitchCoordinator(
		apiKeyRepo,
		&subscriptionKeySwitchGroupRepoStub{group: &Group{ID: 42, Name: "GPT Sub", SubscriptionType: SubscriptionTypeSubscription}},
		&subscriptionKeySwitchSettingRepoStub{value: "true"},
		nil,
	)

	result := coordinator.HandleNewSubscription(context.Background(), 7, 42)
	require.Equal(t, SubscriptionKeySwitchActionManualSwitchRequired, result.Action)
	require.Equal(t, 2, result.APIKeyCount)
	require.Empty(t, apiKeyRepo.updatedKeys)
}

func TestSubscriptionKeySwitchCoordinatorSkipsWhenAlreadyBound(t *testing.T) {
	targetGroupID := int64(42)
	apiKeyRepo := &subscriptionKeySwitchAPIKeyRepoStub{
		keys: []APIKey{{ID: 9, Key: "sk-test", Name: "default", GroupID: &targetGroupID}},
	}
	coordinator := NewSubscriptionKeySwitchCoordinator(
		apiKeyRepo,
		&subscriptionKeySwitchGroupRepoStub{group: &Group{ID: 42, Name: "Claude Pro", SubscriptionType: SubscriptionTypeSubscription}},
		&subscriptionKeySwitchSettingRepoStub{value: "true"},
		nil,
	)

	result := coordinator.HandleNewSubscription(context.Background(), 7, 42)
	require.Equal(t, SubscriptionKeySwitchActionNone, result.Action)
	require.Equal(t, 1, result.APIKeyCount)
	require.Empty(t, apiKeyRepo.updatedKeys)
}

func TestSubscriptionKeySwitchCoordinatorReturnsNoAPIKeysWhenUserHasNone(t *testing.T) {
	apiKeyRepo := &subscriptionKeySwitchAPIKeyRepoStub{keys: []APIKey{}}
	coordinator := NewSubscriptionKeySwitchCoordinator(
		apiKeyRepo,
		&subscriptionKeySwitchGroupRepoStub{group: &Group{ID: 42, Name: "Claude Pro", SubscriptionType: SubscriptionTypeSubscription}},
		&subscriptionKeySwitchSettingRepoStub{value: "false"},
		nil,
	)

	result := coordinator.HandleNewSubscription(context.Background(), 7, 42)
	require.Equal(t, SubscriptionKeySwitchActionNoAPIKeys, result.Action)
	require.Zero(t, result.APIKeyCount)
}
