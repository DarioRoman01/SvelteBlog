<script lang="ts">
  import {
    MaterialApp,
    Button,
    ProgressCircular,
    Icon,
    Overlay,
  } from "svelte-materialify";
  import type { Comment, PaginatedComments } from "../requests/comments";
  import CreateCommentModal from "./CreateCommentModal.svelte";
  import CommentsCard from "./CommentsCard.svelte";
  import { mdiPlusCircle } from "@mdi/js";
  import { api } from "../requests/users";
  import { createEventDispatcher, onMount } from "svelte";

  const dispactch = createEventDispatcher();
  let comments: PaginatedComments;
  let lastComment: Comment;

  export let id: number;
  export let showModal: boolean;
  export let currentUserId: number;

  onMount(async () => {
    comments = await api<PaginatedComments>(`/posts/${id}/comments?limit=10`);
    lastComment = comments.comments[comments.comments.length - 1];
  });

  const refetch = async () => {
    const newComments = await api<PaginatedComments>(
      `/posts/${id}/comments?limit=10`
    );
    comments.comments = newComments.comments;
    showModal = false;
  };

  const loadMore = async (lastComment: Comment) => {
    const newComments = await api<PaginatedComments>(
      `/posts/${id}/comments?limit=10&cursor=${lastComment.createdAt}`
    );

    newComments.comments.forEach((comment) => comments.comments.push(comment));
    comments.hasMore = newComments.hasMore;
    lastComment = comments.comments[comments.comments.length - 1];
  };
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
          <CommentsCard {comment} {currentUserId} on:delete={() => refetch()} />
        </div>
      {/each}
      {#if comments.hasMore === true}
        <div class="d-flex justify-center mt-4 mb-3">
          <Button class="pink lighten-3" on:click={() => loadMore(lastComment)}>
            <Icon path={mdiPlusCircle} />
          </Button>
        </div>
      {/if}
    {/if}
    <Overlay active={showModal}>
      <CreateCommentModal
        postId={id}
        on:create={() => refetch()}
        on:close={() => {
          showModal = false;
          dispactch("closeModal");
        }}
      />
    </Overlay>
  {:else}
    <ProgressCircular size={50} indeterminate color="pink lighten-3" />
  {/if}
</MaterialApp>
