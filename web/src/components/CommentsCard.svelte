<script lang="ts">
    import {
    Card,
    CardText,
    CardActions,
    Button,
    MaterialApp,
  } from "svelte-materialify";
  import type { Comment } from "../requests/comments"

  export let comment: Comment;
</script>

<MaterialApp>
  <div class="d-flex justify-center mt-4 mb-4">
    <Card style="max-width:800px;">
      <CardText>
        <div>{comment.creator}</div>
        <div class="text--primary">
          {comment.body}
        </div>
      </CardText>
      <CardActions>
        <Button text class="primary-text">see More</Button>
      </CardActions>
    </Card>
  </div>
</MaterialApp>