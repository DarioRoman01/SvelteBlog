<script lang="ts">
  import {
    Card,
    CardText,
    CardActions,
    Button,
    Icon,
    MaterialApp,
    ProgressCircular,
    Overlay,
  } from "svelte-materialify";
  import Nav from "../../components/Nav.svelte";
  import Comments from "../../components/Comments.svelte";
  import { api } from "../../requests/users";
  import type { Profile } from "../../requests/profile";
  import type { Post } from "../../requests/posts";
  import { likePost, deletePost } from "../../requests/posts";
  import { params, redirect } from "@roxi/routify";
  import EditPostModal from "../../components/EditPostModal.svelte";
  import { onMount } from "svelte";
  import { mdiHeart, mdiPen, mdiTrashCan } from "@mdi/js";

  let post: Post;
  let stateValue: 0 | 1;
  let clicked: boolean;
  clicked = false;
  let currentUser: Profile;
  let active = false;
  let showModal: boolean;
  showModal = false;

  onMount(async () => {
    post = await api<Post>(`/posts/${$params["id"]}`);
    stateValue = post.stateValue;
    currentUser = await api<Profile>("/me");
  });

  const handleLike = () => {
    clicked = true;
    let value: 0 | 1;
    stateValue === 1 ? (value = 0) : (value = 1);
    const like = likePost(post.id, value);
    like.then(() => {
      stateValue = value;
      clicked = false;
    });
  };

  const handleDelete = () => {
    const deleted = deletePost(post.id);
    deleted.then(() => $redirect("/home"));
  };

  const handleUpdate = (event: CustomEvent) => {
    active = false
    post = event.detail.post;
  };
</script>

<MaterialApp>
  <Nav isLoggedIn={true} />
  {#if post && currentUser}
    <div class="d-flex justify-center mt-4 mb-4">
      <Card raised style="width:50em;">
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
            disabled={clicked}
            class={stateValue === 1 ? "pink-text" : "grey-text"}
            on:click={handleLike}
          >
            <Icon path={mdiHeart} />
          </Button>
          {#if currentUser.userID === post.creatorId}
            <Button icon class="pink-text" on:click={() => (active = true)}>
              <Icon path={mdiPen} />
            </Button>
            <Button icon class="pink-text" on:click={() => handleDelete()}>
              <Icon path={mdiTrashCan} />
            </Button>
          {/if}
        </CardActions>
      </Card>
      <Overlay {active}>
        <EditPostModal
          {post}
          on:close={(_e) => (active = false)}
          on:update={(e) => handleUpdate(e)}
        />
      </Overlay>
    </div>
    <div class="d-flex justify-center">
      <div style="width: 50em;" class="d-flex justify-space-between mt-5 mr-6 ml-6">
        <div>
          <p>Comments:</p>
        </div>
        <div>
          <Button class="pink lighten-3" on:click={() => showModal = true}>
            add comment
          </Button>
        </div>        
      </div>
    </div>
    <div class="d-flex justify-center mt-4 mb-4">
      <Comments 
        on:closeModal={() => showModal = false}
        currentUserId={currentUser.userID} 
        showModal={showModal} 
        id={post.id}
      />
    </div>
  {:else}
    <div class="d-flex justify-center">
      <ProgressCircular size={50} indeterminate color="pink lighten-3" />
    </div>
  {/if}
</MaterialApp>
