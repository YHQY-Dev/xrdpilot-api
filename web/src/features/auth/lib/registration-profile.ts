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
import { z } from 'zod'

const PASSWORD_MIN_LENGTH = 8
const PASSWORD_MAX_LENGTH = 20

export const ACADEMIC_IDENTITIES = [
  { code: 'undergraduate', labelKey: 'Undergraduate' },
  { code: 'master', labelKey: "Master's student" },
  { code: 'phd', labelKey: 'PhD student' },
  { code: 'postdoc', labelKey: 'Postdoctoral researcher' },
  { code: 'lecturer', labelKey: 'Lecturer' },
  { code: 'associate_professor', labelKey: 'Associate professor' },
  { code: 'professor', labelKey: 'Professor' },
  { code: 'researcher', labelKey: 'Researcher' },
  { code: 'engineer', labelKey: 'Engineer' },
  { code: 'other', labelKey: 'Other' },
] as const

export type AcademicIdentityCode = (typeof ACADEMIC_IDENTITIES)[number]['code']

export const ACADEMIC_IDENTITY_CODES = ACADEMIC_IDENTITIES.map(
  (item) => item.code
) as [AcademicIdentityCode, ...AcademicIdentityCode[]]

export const ACADEMIC_IDENTITY_LABEL_KEYS: Record<AcademicIdentityCode, string> =
  ACADEMIC_IDENTITIES.reduce(
    (labels, item) => {
      labels[item.code] = item.labelKey
      return labels
    },
    {} as Record<AcademicIdentityCode, string>
  )

const requiredText = (message: string, max: number) =>
  z
    .string()
    .trim()
    .min(1, message)
    .max(max, message)

export function createRegisterFormSchema(options: {
  emailVerificationRequired: boolean
}) {
  const verificationCode = options.emailVerificationRequired
    ? z.string().trim().min(1, 'Please enter the verification code')
    : z.string().optional()

  return z
    .object({
      email: z
        .string()
        .trim()
        .min(1, 'Please enter your email')
        .email('Please enter a valid email address')
        .max(50, 'Please enter a valid email address'),
      verificationCode,
      realName: requiredText('Please enter your real name', 64),
      organization: requiredText(
        'Please enter your research organization',
        128
      ),
      academicIdentity: z
        .string()
        .refine(
          (value): value is AcademicIdentityCode =>
            ACADEMIC_IDENTITY_CODES.includes(value as AcademicIdentityCode),
          { message: 'Please select your academic identity' }
        ),
      supervisorName: requiredText('Please enter the supervisor name', 64),
      researchDirection: requiredText(
        'Please enter your research direction',
        255
      ),
      usagePurpose: requiredText(
        'Please enter the intended use of the program',
        2000
      ),
      password: z
        .string()
        .min(1, 'Please enter your password')
        .min(
          PASSWORD_MIN_LENGTH,
          'Password must be between 8 and 20 characters'
        )
        .max(
          PASSWORD_MAX_LENGTH,
          'Password must be at most 20 characters long'
        ),
      confirmPassword: z.string().min(1, 'Please confirm your password'),
    })
    .refine((data) => data.password === data.confirmPassword, {
      message: "Passwords don't match.",
      path: ['confirmPassword'],
    })
}

export type RegisterFormValues = z.input<
  ReturnType<typeof createRegisterFormSchema>
>

export const REGISTER_FORM_DEFAULT_VALUES = {
  email: '',
  verificationCode: '',
  realName: '',
  organization: '',
  academicIdentity: '',
  supervisorName: '',
  researchDirection: '',
  usagePurpose: '',
  password: '',
  confirmPassword: '',
}

export function academicIdentityLabelKey(
  code: string | undefined
): string | undefined {
  if (!code) return undefined
  return ACADEMIC_IDENTITY_LABEL_KEYS[code as AcademicIdentityCode]
}
