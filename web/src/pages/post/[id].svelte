<script lang="ts">
  import {
    Card,
    CardText,
    Button, 
    MaterialApp,
    ProgressCircular 
  } from "svelte-materialify";
  import Nav from "../../components/Nav.svelte";
  import { api } from "../../requests/users";
  import type { Post } from "../../requests/posts"
  import { params } from "@roxi/routify";
  import { onMount } from "svelte";

  let post: Post;

  onMount(async () => {
    post = await api<Post>(`/posts/${$params["id"]}`)
  });
</script>

<MaterialApp>
  <Nav isLoggedIn={true}/>
  {#if post}  
    <div class="d-flex justify-center mt-4 mb-4">
      <Card style="max-width:800px;">
        <CardText>
          <div>{post.creator}</div>
          <div class="text--primary text-h4">{post.title}</div>
          <div class="text--primary">
            {post.body}
          </div>
        </CardText>
      </Card>
    </div>
  {:else}
    <div class="d-flex justify-center">
      <ProgressCircular size={50} indeterminate color="pink lighten-3" />
    </div>
  {/if}
</MaterialApp>