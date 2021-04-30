<script lang="ts">
  import {
    Card,
    CardText,
    CardActions,
    Button,
    MaterialApp,
  } from "svelte-materialify";
  import type { Post } from "../requests/posts";
  export let post: Post;
</script>

<MaterialApp>
  <div class="d-flex justify-center mt-4 mb-4">
    <Card style="max-width:800px;">
      <CardText>
        <div>{post.creator}</div>
        <div class="text--primary text-h4">{post.title}</div>
        <div class="text--primary">
          {post.body}
        </div>
      </CardText>
      <CardActions>
        <Button text class="primary-text">see More</Button>
      </CardActions>
    </Card>
  </div>
</MaterialApp>
