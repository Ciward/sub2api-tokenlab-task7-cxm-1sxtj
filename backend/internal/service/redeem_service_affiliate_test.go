package service

import "testing"

func TestComputeRedeemAffiliateBaseAmount_BalancePositive(t *testing.T) {
	got := computeRedeemAffiliateBaseAmount(RedeemTypeBalance, 25)
	if got != 25 {
		t.Fatalf("expected 25, got %v", got)
	}
}

func TestComputeRedeemAffiliateBaseAmount_BalanceNonPositive(t *testing.T) {
	cases := []float64{0, -3}
	for _, v := range cases {
		got := computeRedeemAffiliateBaseAmount(RedeemTypeBalance, v)
		if got != 0 {
			t.Fatalf("expected 0 for value %v, got %v", v, got)
		}
	}
}

func TestComputeRedeemAffiliateBaseAmount_SubscriptionNoRebate(t *testing.T) {
	got := computeRedeemAffiliateBaseAmount(RedeemTypeSubscription, 88)
	if got != 0 {
		t.Fatalf("expected 0 for subscription redeem, got %v", got)
	}
}
