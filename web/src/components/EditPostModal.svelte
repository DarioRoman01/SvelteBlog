<script lang="ts">
  import {
    MaterialApp,
    TextField,
    Button,
    Icon,
    Alert,
  } from "svelte-materialify";
  import { createEventDispatcher } from "svelte";
  import type { Post } from "../requests/posts";
  import { updatePost } from "../requests/posts";
  import { mdiClose } from "@mdi/js";

  const dispatch = createEventDispatcher();
  export let post: Post;
  let error: Error;
  let show: boolean;

  const handleUpdate = () => {
    const update = updatePost(post.id, post.title, post.body);
    update
      .then((post) => dispatch("update", { post: post }))
      .catch((err) => {
        error = err;
      });
  };

  const handleClose = () => {
    dispatch("close");
  };
</script>

<MaterialApp>
  <div class="pt-4 pb-4 pl-4 pr-4">
    <div class="d-flex justify-center mt-4">
      <div class="d-flex flex-column">
        <div class="d-flex justify-space-between">
          <h3 class="text-h3 mb-6">Edit post</h3>
          <Button icon class="pink lighten-3" on:click={handleClose}>
            <Icon path={mdiClose} />
          </Button>
        </div>
        <div color="pink lighten-3" style="width: 700px;" class="mb-4 mt-3">
          <TextField bind:value={post.title} color="pink lighten-3" />
        </div>
        <div style="width: 700px;">
          <TextField color="pink lighten-3" bind:value={post.body} />
        </div>
        <div style="align-self: center;" class="mt-4">
          <Button size="large" class="pink lighten-3" on:click={handleUpdate}>
            Submit
          </Button>
        </div>
        {#if error}
          <Alert
            class="pink accent-3 mt-4"
            dismissible={true}
            bind:visible={show}
          >
            {error.message}
          </Alert>
        {/if}
      </div>
    </div>
  </div>
</MaterialApp>
