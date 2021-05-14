<script lang="ts">
  import {
    MaterialApp,
    TextField,
    Button,
    Icon,
    Alert
  } from "svelte-materialify";
  import { createEventDispatcher } from "svelte";
  import { addComment } from "../requests/comments";
  import { mdiClose } from "@mdi/js";

  const dispatch = createEventDispatcher();

  let body: string;
  let error: Error;
  export let postId: number;
  let show: boolean

  const handleCreate = () => {
    const newComment = addComment(postId, body);
    newComment
      .then(() => dispatch("create"))
      .catch((err) => {
        error = err
        show = true
      });
  }

  const handleClose = () => {
    dispatch("close");
  };
</script>

<MaterialApp>
  <div class="pt-4 pb-4 pl-4 pr-4">
    <div class="d-flex justify-center mt-4">
      <div class="d-flex flex-column">
        <div class="d-flex justify-space-between">
          <h3 class="text-h3 mb-6">comment</h3>
          <Button icon class="pink lighten-3" on:click={handleClose}>
            <Icon path={mdiClose} />
          </Button>
        </div>
        <div style="width: 700px;" class="mb-4 mt-3">
          <TextField color="pink lighten-3" bind:value={body} />
        </div>
        <div style="align-self: center;" class="mt-4">
          <Button size="large" class="pink lighten-3" on:click={handleCreate}>
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