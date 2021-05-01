<script lang="ts">
  import {
    ProgressCircular,
    MaterialApp,
    Icon,
    Button,
  } from "svelte-materialify";
  import { onMount } from "svelte";
  import PostCard from "./PostCard.svelte";
  import type { PaginatedPosts } from "../requests/posts";
  import { api } from "../requests/users";
  import { mdiPlusCircle } from "@mdi/js";

  let posts: PaginatedPosts;
  onMount(async () => {
    posts = await api<PaginatedPosts>("/posts?limit=10")
  });
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
        <Button class="pink lighten-3">
          <Icon path={mdiPlusCircle}/>
        </Button>
      </div>
    {/if}
  {:else}
    <ProgressCircular size={50} indeterminate color="pink lighten-3" />
  {/if}
</MaterialApp>