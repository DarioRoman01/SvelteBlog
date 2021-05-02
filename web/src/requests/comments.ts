export interface Comment {
  creatorId: number;
  postId: number;
  createdAt: string;
  updatedAt: string;
  body: string;
  creator: string;
}

export interface PaginatedComments {
  comments: Comment[];
  hasMore: boolean;
}

const url = process.env.API_URL;
const headers = {
  Accept: "application/json",
  "Content-Type": "application/json",
};

export const addComment = async (
  id: number,
  body: string
): Promise<Comment> => {
  const res = await fetch(`${url}/posts/${id}/comments`, {
    method: "POST",
    body: JSON.stringify({
      body: body,
    }),
    headers: headers,
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err)
  }

  return await res.json();
};

export const deleteComment = async (id: number) => {
  const res = await fetch(`${url}/posts/${id}/comments`, {
    method: "DELETE",
    headers: headers,
    credentials: "include"
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err)
  }

  return await res.json();
}
