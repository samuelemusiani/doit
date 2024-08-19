import { COLORS_URL, LOGIN_URL, NOTES_URL, PRIORITIES_URL, STATES_URL, USERS_URL } from '@/consts'
import type { TodoColor, TodoPriority, TodoState, Todo } from '@/types'
import type { User } from '@/types'

export async function getCurrentUser(): Promise<User> {
  return fetch(LOGIN_URL, {
    credentials: 'include'
  })
    .then(async (res) => {
      if (!res.ok) {
        throw new Error(await res.text())
      }
      return (await res.json()) as User
    })
    .catch((err) => {
      throw new Error(`Could not get current user: ${err}`)
    })
}

export async function getStates(): Promise<TodoState[]> {
  return fetch(STATES_URL)
    .then((res) => {
      return res.json()
    })
    .then((states) => {
      return states as TodoState[]
    })
    .catch((err) => {
      throw new Error(`Could not get states: ${err}`)
    })
}

export async function getPriorities(): Promise<TodoPriority[]> {
  return fetch(PRIORITIES_URL)
    .then((res) => {
      return res.json()
    })
    .then((priorities) => {
      return priorities as TodoPriority[]
    })
    .catch((err) => {
      throw new Error(`Could not get priorities: ${err}`)
    })
}

export async function getColors(): Promise<TodoColor[]> {
  return fetch(COLORS_URL)
    .then((res) => {
      return res.json()
    })
    .then((colors) => {
      return colors as TodoColor[]
    })
    .catch((err) => {
      throw new Error(`Could not get colors: ${err}`)
    })
}

export async function fetchNotes(): Promise<Todo[]> {
  return fetch(NOTES_URL, {
    credentials: 'include'
  })
    .then((res) => {
      return res.json()
    })
    .then((notes) => {
      return notes as Todo[]
    })
    .catch((err) => {
      throw new Error(`Could not get notes: ${err}`)
    })
}

export async function addTodo(todo: Todo): Promise<any> {
  return fetch(NOTES_URL, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(todo)
  })
    .then((res) => {
      return res.text()
    })
    .catch((err) => {
      throw new Error(`Could not add todo: ${err}`)
    })
}

export async function updateTodo(todo: Todo): Promise<any> {
  return fetch(NOTES_URL + '/' + todo.ID, {
    method: 'PUT',
    credentials: 'include',
    body: JSON.stringify(todo)
  })
    .then((res) => {
      return res.text()
    })
    .catch((err) => {
      throw new Error(`Could not update todo: ${err}`)
    })
}

export async function deleteTodo(id: number): Promise<any> {
  return fetch(NOTES_URL + '/' + id, {
    method: 'DELETE',
    credentials: 'include'
  })
    .then((res) => {
      return res.text()
    })
    .catch((err) => {
      throw new Error(`Could not delete todo: ${err}`)
    })
}

export async function login(user: string, password: string): Promise<any> {
  return fetch(LOGIN_URL, {
    method: 'POST',
    body: JSON.stringify({
      username: user,
      password: password
    }),
    // Used to set coockies; DOTO Check if this should be in production
    credentials: 'include'
  })
    .then(async (res) => {
      let t = await res.text()
      if (!res.ok) {
        throw new Error(t)
      }
      return t
    })
    .catch((err) => {
      throw new Error(`Could not login: ${err}`)
    })
}

export async function logout(): Promise<any> {
  return fetch(LOGIN_URL, {
    method: 'DELETE',
    credentials: 'include'
  })
    .then((res) => {
      return res.text()
    })
    .catch((err) => {
      throw new Error(`Could not logout: ${err}`)
    })
}

export async function getUsers(): Promise<User[]> {
  return fetch(USERS_URL, {
    credentials: 'include'
  })
    .then((res) => {
      return res.json()
    })
    .then((users) => {
      return users as User[]
    })
    .catch((err) => {
      throw new Error(`Could not get users: ${err}`)
    })
}

export async function getUser(id: number): Promise<User> {
  return fetch(USERS_URL + `/${id}`, {
    credentials: 'include'
  })
    .then((res) => {
      return res.json()
    })
    .then((user) => {
      return user as User
    })
    .catch((err) => {
      throw new Error(`Could not get user: ${err}`)
    })
}

export async function modifyUser(user: User): Promise<User> {
  return fetch(USERS_URL + `/${user.ID}`, {
    method: 'PUT',
    credentials: 'include',
    body: JSON.stringify(user)
  })
    .then((res) => {
      return res.json()
    })
    .catch((err) => {
      throw new Error(`Could not update user: ${err}`)
    })
}

export async function addUser(user: User): Promise<User> {
  return fetch(USERS_URL, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(user)
  })
    .then(async (res) => {
      if (!res.ok) {
        throw new Error(await res.text())
      }

      return await res.json()
    })
    .catch((err) => {
      throw new Error(`Could not create user: ${err}`)
    })
}

export async function deleteUser(id: number): Promise<string> {
  return fetch(USERS_URL + `/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  })
    .then(async (res) => {
      let t = await res.text()
      if (!res.ok) {
        throw new Error(t)
      }

      return t
    })
    .catch((err) => {
      throw new Error(`Could not delete user: ${err}`)
    })
}

export async function isLoggedIn(): Promise<boolean> {
  return fetch(LOGIN_URL, {
    credentials: 'include'
  }).then((res) => {
    if (res.ok) {
      return true
    } else if (res.status == 401) {
      return false
    } else {
      throw new Error(`Could not check if user is logged in: ${res.status}`)
    }
  })
}
