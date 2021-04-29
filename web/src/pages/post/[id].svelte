<script lang="ts">
  import {
    Card,
    CardText,
    Button, 
    MaterialApp 
  } from "svelte-materialify";
  import { api } from "../../requests/users";
  import type { Post } from "../../requests/posts"
  import { params } from "@roxi/routify";
  import { onMount } from "svelte";

  let post: Post;

  onMount(async () => {
    post = await api<Post>(`/posts/${$params["id"]}`)
  })
</script>

<MaterialApp>
  <div class="d-flex justify-center mt-4 mb-4">
    <Card style="max-width:800px;">
      <CardText>
        <div>{post.creator.username}</div>
        <div class="text--primary text-h4">{post.title}</div>
        <div class="text--primary">
          {post.body}
        </div>
      </CardText>
    </Card>
  </div>
</MaterialApp>