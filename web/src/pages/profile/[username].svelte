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
  import Nav from "../../components/Nav.svelte";
  import { api } from "../../requests/users";
  import { mdiAccountCircle } from "@mdi/js";
  import { params } from "@roxi/routify";
  import { onMount } from "svelte";

  let profile: Profile;
  onMount(async () => {
    profile = await api<Profile>(`/profile/${$params["username"]}`);
  });
</script>

<MaterialApp>
  <Nav isLoggedIn={true}/>
  {#if profile}
  <div class="d-flex justify-center align-center">
    <Container>
      <Row>
        <Col class="align-self-center">
          <div class="pl-2">
            <Icon size="120px" path={mdiAccountCircle} />
            <Button class="pink lighten-3">follow</Button>
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
  {:else}
    <div class="d-flex justify-center">
      <ProgressCircular size={50} indeterminate color="pink lighten-3" />
    </div>
  {/if}
</MaterialApp>
