<script lang="ts">
  import type { PaginatedComments } from "../requests/comments";
  import CommentsCard from "./CommentsCard.svelte"
  import { MaterialApp, Button, ProgressCircular, Icon } from "svelte-materialify";
  import { mdiPlusCircle } from "@mdi/js"
  import { api } from "../requests/users";
  import { onMount } from "svelte";

  let comments: PaginatedComments;
  export let id: number;

  onMount(async () => {
    comments = await api<PaginatedComments>(`/posts/${id}/comments?limit=10`);
  });
</script>

<MaterialApp>
  {#if comments}
    {#if comments.comments === null}
      <div class="mt-4">
        <p>no comments yet :(</p>
      </div>
    {:else}
      {#each comments.comments as comment}
        <div class="d-flex flex-column">
          <CommentsCard comment={comment}/>
        </div>
      {/each}
      {#if comments.hasMore === true}
        <div class="d-flex justify-center mt-4 mb-3">
          <Button class="pink lighten-3">
            <Icon path={mdiPlusCircle}/>
          </Button>
        </div>
      {/if}
    {/if}
  {:else}
    <ProgressCircular size={50} indeterminate color="pink lighten-3" />
  {/if}
</MaterialApp>
