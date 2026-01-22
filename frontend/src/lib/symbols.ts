import type {InjectionKey, Ref} from 'vue'

export interface FormRule {
  required?: boolean
  message?: string
  pattern?: RegExp
  validator?: (value: any) => boolean | string | Promise<boolean | string | void>
  trigger?: 'change' | 'blur'
}

export type FormRules = Record<string, FormRule | FormRule[]>

export interface FormContext {
  model?: Record<string, any>
  rules?: FormRules
  errors: Record<string, string>
  addField: (field: string) => void
  removeField: (field: string) => void
}

export interface FormItemContext {
  validationError: Ref<string | undefined>
}

export const FormContextKey = Symbol('FormContextKey') as InjectionKey<FormContext>
export const FormItemContextKey = Symbol('FormItemContextKey') as InjectionKey<FormItemContext>
