const url = process.env.API_URL
const headers = {
  Accept: "application/json",
  "Content-Type": "application/json",
}

// custom error class for users signup and register requests
export class CustomError extends Error {
  constructor(public field: string, public message: string) {
    super(message);
    this.field = field;

  }
}

// user data that is retrieved from the server.
export interface User {
  id: number;
  createdAt: string;
  updatedAt: string;
  email: string;
}

// api call with generics for GET request
export async function api<T>(url: string): Promise<T> {
  const response = await fetch(url, {
    credentials: 'include',
    headers: headers,
  })
  if (!response.ok) {
    throw new Error("something wrong happend :(")
  }
  return await response.json();
}

// handle login request
export const login =  async ({email, password}): Promise<User> => {
  const res = await fetch(`${url}/login`, {
    method: "post",
    body: JSON.stringify({
      email: email,
      password: password,
    }),
    headers: headers,
    credentials: "include",
  });
  
  if (!res.ok) {
    const err: CustomError = await res.json();
    throw new CustomError(err.field, err.message);
  }

  return await res.json();
}

// send register request to server
export const register = async (
  email: string, 
  phoneNumber: number, 
  password: string,
  passwordConfirmation: string
): Promise<User> => {
  const res = await fetch(`${url}/signup`, {
    method: "POST",
    body: JSON.stringify({
      email: email,
      phoneNumber: phoneNumber,
      password: password,
      passwordConfirmation: passwordConfirmation
    }),
    headers: headers,
    credentials: "include",
  });
 
  if (!res.ok) {
    const err: CustomError = await res.json();
    throw new CustomError(err.field, err.message)
  }

  return await res.json();
}

// handle users forgot password request
export const forgotPassword = async (email: string) => {
  const res = await fetch(`${url}/forgot-password`, {
    method: "post",
    body: JSON.stringify({
      email: email
    }),
    headers: headers,
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.message)
  }

  return await res.json();
}

// handles users change password request
export const changePassword = async (
  token: string, 
  newPassword: string
  ): Promise<User> => {
  const res = await fetch(`${url}/change-password`, {
    method: "POST",
    body: JSON.stringify({
      token: token,
      newPassword: newPassword
    }),
    headers: headers,
    credentials: "include",
  });

  if (!res.ok) {
    const err: CustomError = await res.json();
    throw new CustomError(err.field, err.message);
  }

  return await res.json();
}

export const verify = async (token: string) => {
  const res = await fetch(`${url}/verify`, {
    method: "POST",
    body: JSON.stringify({
      token: token
    }),
    headers: headers,
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err)
  }

  return await res.json();
}

export const logout = async () => {
  const res = await fetch(`${url}/logout`, {
    method: "post",
    credentials: "include"
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err)
  }

  return await res.json();
}

export const emailRegex = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;