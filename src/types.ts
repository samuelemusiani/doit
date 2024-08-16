export interface Expiration {
  DoesExpire: boolean
  Date: string
}

export interface Todo {
  ID: number
  Title: string
  Description: string
  StateID: number
  PriorityID: number
  ColorID: number
  Expiration: Expiration
}

export interface TodoState {
  ID: number
  State: string
}

export interface TodoPriority {
  ID: number
  Priority: string
}

export interface TodoColor {
  ID: number
  Hex: string
}

export interface Options {
  States: TodoState[]
  Priorities: TodoPriority[]
  Colors: TodoColor[]
}
