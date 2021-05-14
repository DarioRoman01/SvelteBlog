<script lang="ts">
  import {
    ProgressCircular,
    MaterialApp,
    Icon,
    Button,
  } from "svelte-materialify";
  import { onMount } from "svelte";
  import PostCard from "./PostCard.svelte";
  import type { PaginatedPosts, Post } from "../requests/posts";
  import { api } from "../requests/users";
  import { mdiPlusCircle } from "@mdi/js";

  let lastPost: Post;
  let posts: PaginatedPosts;
  onMount(async () => {
    posts = await api<PaginatedPosts>("/posts?limit=10")
    lastPost = posts.posts[posts.posts.length - 1];
  });

  const loadMore = async (lastPost: Post) => {
    const newPosts = await api<PaginatedPosts>(`/posts?limit=10&cursor=${lastPost.createdAt}`);
    newPosts.posts.forEach((post) => posts.posts.push(post));
    posts.hasMore = newPosts.hasMore;
    lastPost = posts.posts[posts.posts.length - 1];
  }
</script>

<MaterialApp>
  {#if posts}
    {#each posts.posts as post}
      <div class="d-flex flex-column">
        <PostCard post={post}/>
      </div>
    {/each}
    {#if posts.hasMore == true}
      <div class="d-flex justify-center mt-4 mb-3">
        <Button class="pink lighten-3" on:click={() => loadMore(lastPost)}>
          <Icon path={mdiPlusCircle}/>
        </Button>
      </div>
    {/if}
  {:else}
    <ProgressCircular size={50} indeterminate color="pink lighten-3" />
  {/if}
</MaterialApp>