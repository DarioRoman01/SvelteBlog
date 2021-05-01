<script lang="ts">
  import {
    TextField,
    Button,
    MaterialApp,
    Alert,
    Textarea,
    ProgressCircular,
  } from "svelte-materialify";
  import { onMount } from "svelte";
  import Nav from "../../../components/Nav.svelte";
  import { api } from "../../../requests/users";
  import type { CustomError } from "../../../requests/users";
  import type { Profile } from "../../../requests/profile";
  import { updateProfile } from "../../../requests/profile";
  import { redirect } from "@roxi/routify";

  let profile: Profile;
  let error: CustomError;
  let show: boolean;

  onMount(async () => {
    profile = await api<Profile>("/me");
  });

  const handleUpdate = () => {
    const updated = updateProfile(profile);
    updated
      .then(() => $redirect(`/profile/${profile.username}`))
      .catch((err) => {
        error = err
        show = true;
      });
  }
</script>

<MaterialApp>
  <Nav isLoggedIn={true} />
  {#if profile}
    <div class="d-flex justify-center mt-4">
      <div class="d-flex flex-column">
        <h3 class="text-h3 mb-6">edit profile</h3>
        <div style="width: 700px;" class="mb-4 mt-3">
          <TextField
            color="pink lighten-3"
            bind:value={profile.username}
          >
            username
          </TextField>
        </div>
        <div style="width: 700px;">
          <Textarea color="pink lighten-3" bind:value={profile.biography}>
            biography
          </Textarea>
        </div>
        <div style="align-self: center;" class="mt-4">
          <Button 
            size="large" 
            class="pink lighten-3" 
            on:click={handleUpdate}
          >
            Submit
          </Button>
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
  {:else}
    <div class="d-flex justify-center">
      <ProgressCircular size={50} indeterminate color="pink lighten-3"/>
    </div>
  {/if}
</MaterialApp>