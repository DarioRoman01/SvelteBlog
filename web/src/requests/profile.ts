const url = process.env.API_URL
const headers =  {
  Accept: "application/json",
  "Content-Type": "application/json",
} 
export interface Profile {
  userID: number;
  createdAt: string;
  updatedAt: string;
  username: string;
  biography: string;
  followers: number;
  posted: number;
}

export const CreateProfile = async (
  username: string,
  biography: string
): Promise<Profile> => {
  const res = await fetch(`${url}/profile`, {
    method: "POST",
    body: JSON.stringify({
      username: username,
      biography: biography,
    }),
    headers:headers,
    credentials: "include"
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err);
  }

  return await res.json();
};

export const updateProfile = async (
  id: number,
  username: string,
  biography: string
): Promise<Profile> => {
  const res = await fetch(`${url}/profile/${id}`, {
    method: "PATCH",
    body: JSON.stringify({
      username: username,
      biography: biography
    }),
    headers: headers,
    credentials: "include"
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err);
  }

  return await res.json();
}

export const ToggleFollow = async (profileId: number) => {
  const res = await fetch(`${url}/profile/${profileId}/follow`, {
    method: "POST",
    headers: headers,
    credentials: "include"
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err)
  }

  return await res.json();
}
