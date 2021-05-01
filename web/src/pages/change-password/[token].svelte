<script lang="ts">
  import Nav from "../../components/Nav.svelte";
  import { changePassword } from "../../requests/users";
  import { params, redirect } from "@roxi/routify";
  import { TextField, Button, MaterialApp, Alert } from "svelte-materialify";

  let newPassword: string;
  let newPasswordConfirmation: string;
  let error: Error;

  const handleClick = () => {
    const user = changePassword(
      $params["token"],
      newPassword,
      newPasswordConfirmation
    );
    user.then(() => $redirect("../home")).catch((err: Error) => (error = err));
  };
</script>

<MaterialApp>
  <Nav isLoggedIn={false} />
  <div class="d-flex justify-center mt-4">
    <div class="d-flex flex-column">
      <h3 class="text-h3 mb-6">Login</h3>
      <div style="width: 700px;" class="mb-4 mt-3">
        <TextField color="pink lighten-3" bind:value={newPassword}>
          new password
        </TextField>
      </div>
      <div style="width: 700px;" class="mb-4 mt-3">
        <TextField color="pink lighten-3" bind:value={newPasswordConfirmation}>
          new password confirmation
        </TextField>
      </div>
      <div style="align-self: center;" class="mt-4">
        <Button size="large" class="pink lighten-3" on:click={handleClick}>
          Submit
        </Button>
      </div>
      {#if error}
        <Alert class="pink accent-3">
          {error.message}
        </Alert>
      {/if}
    </div>
  </div>
</MaterialApp>
