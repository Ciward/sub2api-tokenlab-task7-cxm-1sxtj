export type SubscriptionKeySwitchAction =
  | 'none'
  | 'auto_switched'
  | 'manual_switch_required'
  | 'no_api_keys'

export interface SubscriptionKeySwitchResult {
  action: SubscriptionKeySwitchAction
  target_group_id?: number
  target_group_name?: string
  api_key_count?: number
  auto_switch_enabled?: boolean
  auto_switched_key_id?: number
  auto_switched_key_name?: string
}
