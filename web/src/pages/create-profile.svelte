<script lang="ts">
  import {
    TextField,
    Button,
    MaterialApp,
    Alert,
    Textarea,
  } from "svelte-materialify";
  import type { CustomError } from "../requests/users";
  import Nav from "../components/Nav.svelte";
  import { CreateProfile } from "../requests/profile";
  import { redirect } from "@roxi/routify";

  let username: string;
  let biography: string;
  let error: CustomError;
  let show: boolean;

  const usernameRules = [
    (v: string) => !!v || "Required",
    (v: string) => v.length >= 3 || "must be at least 3 characters",
  ];

  const handleCreate = () => {
    const profile = CreateProfile(username, biography);
    profile.then(() => $redirect("./home")).catch((err) => (error = err));
  };
</script>

<MaterialApp>
  <Nav isLoggedIn={true} />
  <div class="d-flex justify-center mt-4">
    <div class="d-flex flex-column">
      <h3 class="text-h3 mb-6">Login</h3>
      <div style="width: 700px;" class="mb-4 mt-3">
        <TextField
          color="pink lighten-3"
          bind:value={username}
          rules={usernameRules}
        >
          username
        </TextField>
      </div>
      <div style="width: 700px;">
        <Textarea color="pink lighten-3" bind:value={biography}>
          biography
        </Textarea>
      </div>
      <div style="align-self: center;" class="mt-4">
        <Button 
          size="large" 
          class="pink lighten-3" 
          on:click={handleCreate}
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
</MaterialApp>
