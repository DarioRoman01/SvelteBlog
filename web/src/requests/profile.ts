import { CustomError } from "./users";

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
  followState: boolean;
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

export const updateProfile = async (profile: Profile): Promise<Profile> => {
  const res = await fetch(`${url}/profile/${profile.userID}`, {
    method: "PATCH",
    body: JSON.stringify({
      username: profile.username,
      biography: profile.biography
    }),
    headers: headers,
    credentials: "include"
  });

  if (!res.ok) {
    const err: CustomError = await res.json();
    throw new CustomError(err.field, err.message);
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
