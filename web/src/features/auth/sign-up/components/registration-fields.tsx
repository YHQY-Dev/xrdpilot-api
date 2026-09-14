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
import type { UseFormReturn } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import {
  ACADEMIC_IDENTITIES,
  type RegisterFormValues,
} from '@/features/auth/lib/registration-profile'

export function RegistrationProfileFields(props: {
  form: UseFormReturn<RegisterFormValues>
}) {
  const { t } = useTranslation()

  return (
    <>
      <FormField
        control={props.form.control}
        name='realName'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Real name')}</FormLabel>
            <FormControl>
              <Input placeholder={t('Please enter your real name')} {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={props.form.control}
        name='organization'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Research organization')}</FormLabel>
            <FormControl>
              <Input
                placeholder={t('Please enter your research organization')}
                {...field}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={props.form.control}
        name='academicIdentity'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Academic identity')}</FormLabel>
            <Select
              items={ACADEMIC_IDENTITIES.map((item) => ({
                value: item.code,
                label: t(item.labelKey),
              }))}
              onValueChange={(value) => {
                if (value !== null) {
                  field.onChange(value)
                }
              }}
              value={field.value || null}
            >
              <FormControl>
                <SelectTrigger className='w-full'>
                  <SelectValue placeholder={t('Select academic identity')} />
                </SelectTrigger>
              </FormControl>
              <SelectContent alignItemWithTrigger={false}>
                <SelectGroup>
                  {ACADEMIC_IDENTITIES.map((item) => (
                    <SelectItem key={item.code} value={item.code}>
                      {t(item.labelKey)}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={props.form.control}
        name='supervisorName'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Supervisor name')}</FormLabel>
            <FormControl>
              <Input
                placeholder={t('Please enter the supervisor name')}
                {...field}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={props.form.control}
        name='researchDirection'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Research direction')}</FormLabel>
            <FormControl>
              <Input
                placeholder={t('Please enter your research direction')}
                {...field}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={props.form.control}
        name='usagePurpose'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Intended use of the program')}</FormLabel>
            <FormControl>
              <Textarea
                placeholder={t('Please enter the intended use of the program')}
                rows={4}
                {...field}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    </>
  )
}
