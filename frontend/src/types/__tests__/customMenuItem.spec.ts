import { describe, it, expect } from 'vitest'
import type { CustomMenuItem, CustomMenuOpenMode } from '@/types'

describe('CustomMenuItem types', () => {
  describe('CustomMenuOpenMode', () => {
    it('supports iframe mode', () => {
      const mode: CustomMenuOpenMode = 'iframe'
      expect(mode).toBe('iframe')
    })

    it('supports redirect mode', () => {
      const mode: CustomMenuOpenMode = 'redirect'
      expect(mode).toBe('redirect')
    })
  })

  describe('CustomMenuItem', () => {
    it('can have iframe open_mode', () => {
      const item: CustomMenuItem = {
        id: 'test-1',
        label: 'Test Menu',
        icon_svg: '<svg></svg>',
        url: 'https://example.com/page',
        visibility: 'user',
        sort_order: 0,
        open_mode: 'iframe',
      }
      expect(item.open_mode).toBe('iframe')
    })

    it('can have redirect open_mode', () => {
      const item: CustomMenuItem = {
        id: 'test-2',
        label: 'External Site',
        icon_svg: '<svg></svg>',
        url: 'https://external-site.com',
        visibility: 'user',
        sort_order: 1,
        open_mode: 'redirect',
      }
      expect(item.open_mode).toBe('redirect')
    })

    it('open_mode is optional for backward compatibility', () => {
      // Legacy items without open_mode field should still be valid
      const item: CustomMenuItem = {
        id: 'legacy',
        label: 'Legacy Menu',
        icon_svg: '<svg></svg>',
        url: 'https://example.com',
        visibility: 'admin',
        sort_order: 0,
        // Notice: no open_mode field
      }
      // open_mode should be undefined but treated as 'iframe' by default
      expect(item.open_mode).toBeUndefined()
    })
  })

  describe('open_mode semantics', () => {
    it('empty/missing open_mode defaults to iframe behavior', () => {
      const items: CustomMenuItem[] = [
        { id: '1', label: 'A', icon_svg: '', url: 'https://a.com', visibility: 'user', sort_order: 0 },
        { id: '2', label: 'B', icon_svg: '', url: 'https://b.com', visibility: 'user', sort_order: 1, open_mode: 'iframe' },
      ]

      // Both should use iframe embedding
      for (const item of items) {
        const usesIframe = item.open_mode !== 'redirect'
        expect(usesIframe).toBe(true)
      }
    })

    it('explicit redirect mode bypasses iframe', () => {
      const item: CustomMenuItem = {
        id: 'ext',
        label: 'External',
        icon_svg: '',
        url: 'https://external.com',
        visibility: 'user',
        sort_order: 0,
        open_mode: 'redirect',
      }

      const usesRedirect = item.open_mode === 'redirect'
      expect(usesRedirect).toBe(true)
    })
  })
})
