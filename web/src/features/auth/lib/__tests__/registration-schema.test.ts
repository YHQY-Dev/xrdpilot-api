/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { describe, expect, test } from 'vitest'

import {
  ACADEMIC_IDENTITY_CODES,
  createRegisterFormSchema,
} from '../registration-profile'

const validProfile = {
  email: 'jane.doe@xrdpilot.com',
  verificationCode: '123456',
  realName: '张三',
  organization: '某大学化学系',
  academicIdentity: 'phd' as const,
  supervisorName: '李老师',
  researchDirection: '计算化学',
  usagePurpose: '用于课题第一性原理计算',
  password: 'password12',
  confirmPassword: 'password12',
}

describe('createRegisterFormSchema', () => {
  test('accepts a complete registration payload without a username field', () => {
    const parsed = createRegisterFormSchema({
      emailVerificationRequired: true,
    }).parse(validProfile)

    expect(parsed).not.toHaveProperty('username')
    expect(parsed.academicIdentity).toBe('phd')
    expect(ACADEMIC_IDENTITY_CODES).toContain(parsed.academicIdentity)
  })

  test('rejects a free-text academic identity that is not a stable code', () => {
    const result = createRegisterFormSchema({
      emailVerificationRequired: false,
    }).safeParse({
      ...validProfile,
      academicIdentity: '研究生',
    })

    expect(result.success).toBe(false)
  })

  test('requires verification code only when email verification is enabled', () => {
    const required = createRegisterFormSchema({
      emailVerificationRequired: true,
    }).safeParse({
      ...validProfile,
      verificationCode: '',
    })
    expect(required.success).toBe(false)

    const optional = createRegisterFormSchema({
      emailVerificationRequired: false,
    }).safeParse({
      ...validProfile,
      verificationCode: '',
    })
    expect(optional.success).toBe(true)
  })

  test('always requires email even when verification is disabled', () => {
    const result = createRegisterFormSchema({
      emailVerificationRequired: false,
    }).safeParse({
      ...validProfile,
      email: '',
    })

    expect(result.success).toBe(false)
  })
})
