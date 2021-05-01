<script lang="ts">
  import {
    Card,
    CardText,
    CardActions,
    Button,
    Icon,
    MaterialApp,
  } from "svelte-materialify";
  import type { Post } from "../requests/posts";
  import { likePost } from "../requests/posts";
  import { mdiHeart } from "@mdi/js";
  import { redirect } from "@roxi/routify";

  export let post: Post;
  let stateValue = post.stateValue;

  const handleLike = () => {
    let value: 0 | 1;
    stateValue === 1 ? value = 0 : value = 1;
    const like = likePost(post.id, value);
    like.then(() => stateValue = value);
  }
</script>

<MaterialApp>
  <div class="d-flex justify-center mt-4 mb-4">
    <Card raised style="width:800px;">
      <CardText>
        <div>{post.creator}</div>
        <div class="text--primary text-h4">{post.title}</div>
        <div class="text--primary">
          {post.body}
        </div>
      </CardText>
      <CardActions>
        <Button 
          icon 
          class={stateValue === 1 ? "pink-text" : "grey-text"}
          on:click={handleLike}
        >
          <Icon path={mdiHeart}/>
        </Button>
        <Button text class="pink-text" on:click={$redirect(`/post/${post.id}`)}>
          see More
        </Button>
      </CardActions>
    </Card>
  </div>
</MaterialApp>
