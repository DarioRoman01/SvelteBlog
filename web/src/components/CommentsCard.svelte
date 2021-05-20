<script lang="ts">
    import {
    Card,
    CardText,
    CardActions,
    Button,
    MaterialApp,
    Icon
  } from "svelte-materialify";
  import type { Comment } from "../requests/comments"
  import { mdiTrashCan } from "@mdi/js"
  import { createEventDispatcher } from "svelte";

  export let comment: Comment;
  export let currentUserId: number;
  const dispactch = createEventDispatcher();
  const handleDelete = () => {
    dispactch("delete");
  }
</script>

<MaterialApp>
    <Card raised style="width:50em;">
      <div class="mt-4 mb-4 d-flex justify-space-between">
      <CardText>
        <div>{comment.creator}</div>
        <div class="text--primary">
          {comment.body}
        </div>
      </CardText>
      {#if currentUserId === comment.creatorId}
        <CardActions>
          <Button icon class="pink-text">
            <Icon path={mdiTrashCan} on:click={() => handleDelete}/>
          </Button>
        </CardActions>
      {/if}
      </div>
    </Card>
</MaterialApp>