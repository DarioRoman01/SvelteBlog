<script lang="ts">
  import {
    TextField, 
    Button, 
    MaterialApp, 
    Alert 
  } from 'svelte-materialify';
  import { forgotPassword } from "../requests/users";
  import Nav from "../components/Nav.svelte"

  let email: string;
  let error: Error;
  let send: boolean;

  const handleClick = () => {
    const result = forgotPassword(email);
    result.then(() => send = true).catch((err: Error) => error = err)
  }
</script>

<MaterialApp>
  <Nav isLoggedIn={false}/>
  <div class="d-flex justify-center mt-4">
    <div class="d-flex flex-column">
      {#if !send}   
        <h3 class="text-h3 mb-6">Change password</h3>
        <div style="width: 700px;" class="mb-4 mt-3">
          <TextField color="pink lighten-3" bind:value={email}>
            email
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
      {:else}
        <h2 class="text-h2">We send you an email</h2>
        <h4 class="text-h4">please check your email</h4>
      {/if}
    </div>
  </div>
</MaterialApp>