<script lang="ts">
  import {
    MaterialApp,
    Icon,
    Container,
    Row,
    Col,
    Button,
    ProgressCircular,
  } from "svelte-materialify";
  import type { Profile } from "../../requests/profile";
  import { ToggleFollow } from "../../requests/profile"
  import PostCard from "../../components/PostCard.svelte";
  import type { PaginatedPosts } from "../../requests/posts";
  import Nav from "../../components/Nav.svelte";
  import { api } from "../../requests/users";
  import { mdiAccountCircle } from "@mdi/js";
  import { params } from "@roxi/routify";
  import { onMount } from "svelte";
  import { redirect } from "@roxi/routify";

  let profile: Profile;
  let posts: PaginatedPosts;
  let currentUser: Profile;
  let followState: boolean;
  let clicked: boolean;
  clicked = false;

  onMount(async () => {
    profile = await api<Profile>(`/profile/${$params["username"]}`);
    followState = profile.followState;
    currentUser = await api<Profile>("/me");
    posts = await api<PaginatedPosts>(`/profile/${profile.userID}/posts?limit=10`);
  });

  const handleFollow = () => {
    clicked = true;
    const follow = ToggleFollow(profile.userID);
    follow
      .then(() => {
        clicked = false
        followState = !followState
        profile.followers = followState ? profile.followers + 1 : profile.followers - 1
      })
      .catch(err => console.log(err))
  }
</script>

<MaterialApp>
  <Nav isLoggedIn={true}/>
  {#if profile && currentUser}
    <div class="d-flex justify-center align-center">
      <Container>
        <Row>
          <Col class="align-self-center">
            <div class="pl-2">
              <Icon size="120px" path={mdiAccountCircle} />
              {#if profile.userID === currentUser.userID}
                <Button class="pink lighten-3" on:click={$redirect(`/profile/${currentUser.username}/update`)}>
                  edit profile
                </Button>
              {:else}
                <Button class="pink lighten-3" disabled={clicked} on:click={handleFollow}>
                  {followState ? "unfollow" : "follow"}
                </Button>
              {/if}
            </div>
          </Col>
          <Col class="align-self-center">
            <div class="d-flex justify-space-between align-center ">
              <p class="text-h5">{profile.username}</p>
              <p class="text-h5">Followers: {profile.followers}</p>
              <p class="text-h5">Posted: {profile.posted}</p>
            </div>
          </Col>
        </Row>
      </Container>
    </div>
    {#if posts}
      <Container class="d-flex justify-center align-center">
        <Row>
          {#each posts.posts as post}
            <Col>
              <PostCard post={post}/>
            </Col>
          {/each}
        </Row>
      </Container>
    {:else}
      <div class="d-flex justify-center">
        <ProgressCircular size={50} indeterminate color="pink lighten-3" />
      </div>
    {/if}
  {:else}
    <div class="d-flex justify-center">
      <ProgressCircular size={50} indeterminate color="pink lighten-3" />
    </div>
  {/if}
</MaterialApp>
