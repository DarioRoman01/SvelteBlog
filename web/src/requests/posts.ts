export interface Post {
  id: number;
  createdAt: string;
  updatedAt: string;
  creatorId: number;
  stateValue: 0 | 1;
  title: string;
  body: string;
  likes: string;
}

export interface PaginatedPosts {
  posts: Post[];
  hasMore: boolean;
}

const url = process.env.API_URL
const headers = {
  Accept: "application/json",
  "Content-Type": "application/json",
}


export const CreatePost = async (
  title: string,
  body: string
): Promise<Post> => {
  const res = await fetch(`${url}/posts`, {
    method: "POST",
    body: JSON.stringify({
      title: title,
      body: body,
    }),
    headers: headers,
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err);
  }

  return await res.json();
};

export const likePost = async (id: number, value: 0 | 1) => {
  const res = await fetch(`${url}/posts/${id}/like`, {
    method: "POST",
    body: JSON.stringify({
      value: value
    }),
    headers: headers,
    credentials: "include",
  });
  
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err);
  }

  return await res.json();
}