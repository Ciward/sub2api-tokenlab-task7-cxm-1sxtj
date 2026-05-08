package service

import (
	"context"
	"encoding/json"
	"strconv"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentauditlog"
	"github.com/Wei-Shaw/sub2api/internal/payment"
)

func (s *PaymentService) maybeHandleSubscriptionKeySwitch(ctx context.Context, order *dbent.PaymentOrder) {
	if s == nil || s.subscriptionKeySwitch == nil || order == nil || order.SubscriptionGroupID == nil {
		return
	}

	result := s.subscriptionKeySwitch.HandleNewSubscription(ctx, order.UserID, *order.SubscriptionGroupID)
	if result == nil || !result.NeedsUserNotice() {
		return
	}

	s.writeAuditLog(ctx, order.ID, PaymentAuditActionSubscriptionKeySwitchResult, "system", map[string]any{
		"action":                 result.Action,
		"target_group_id":        result.TargetGroupID,
		"target_group_name":      result.TargetGroupName,
		"api_key_count":          result.APIKeyCount,
		"auto_switch_enabled":    result.AutoSwitchEnabled,
		"auto_switched_key_id":   result.AutoSwitchedKeyID,
		"auto_switched_key_name": result.AutoSwitchedKeyName,
	})
}

func (s *PaymentService) GetSubscriptionKeySwitchResult(ctx context.Context, orderID, userID int64) (*SubscriptionKeySwitchResult, error) {
	order, err := s.GetOrder(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	return s.GetSubscriptionKeySwitchResultForOrder(ctx, order), nil
}

func (s *PaymentService) GetSubscriptionKeySwitchResultForOrder(ctx context.Context, order *dbent.PaymentOrder) *SubscriptionKeySwitchResult {
	result := &SubscriptionKeySwitchResult{Action: SubscriptionKeySwitchActionNone}
	if order == nil || order.OrderType != payment.OrderTypeSubscription || order.SubscriptionGroupID == nil {
		return result
	}
	if order.Status != OrderStatusCompleted {
		return result
	}

	if persisted := s.GetPersistedSubscriptionKeySwitchResultForOrder(ctx, order); persisted.Action != SubscriptionKeySwitchActionNone {
		return persisted
	}

	return s.describeSubscriptionKeySwitchFallback(ctx, order, result)
}

func (s *PaymentService) GetPersistedSubscriptionKeySwitchResultForOrder(ctx context.Context, order *dbent.PaymentOrder) *SubscriptionKeySwitchResult {
	result := &SubscriptionKeySwitchResult{Action: SubscriptionKeySwitchActionNone}
	if order == nil || order.OrderType != payment.OrderTypeSubscription || order.SubscriptionGroupID == nil {
		return result
	}
	if order.Status != OrderStatusCompleted {
		return result
	}

	logEntry, err := s.entClient.PaymentAuditLog.Query().
		Where(
			paymentauditlog.OrderIDEQ(strconv.FormatInt(order.ID, 10)),
			paymentauditlog.ActionEQ(PaymentAuditActionSubscriptionKeySwitchResult),
		).
		Order(dbent.Desc(paymentauditlog.FieldCreatedAt), dbent.Desc(paymentauditlog.FieldID)).
		First(ctx)
	if err != nil {
		return result
	}

	var parsed SubscriptionKeySwitchResult
	if err := json.Unmarshal([]byte(logEntry.Detail), &parsed); err != nil {
		return result
	}
	if parsed.Action == "" {
		parsed.Action = SubscriptionKeySwitchActionNone
	}
	return &parsed
}

func (s *PaymentService) describeSubscriptionKeySwitchFallback(ctx context.Context, order *dbent.PaymentOrder, defaultResult *SubscriptionKeySwitchResult) *SubscriptionKeySwitchResult {
	if s == nil || s.subscriptionKeySwitch == nil || order == nil || order.SubscriptionGroupID == nil {
		return defaultResult
	}
	return s.subscriptionKeySwitch.DescribeCurrentState(ctx, order.UserID, *order.SubscriptionGroupID)
}
