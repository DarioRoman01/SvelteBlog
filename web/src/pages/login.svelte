<script lang="ts">
  import { login } from "../requests/users";
  import Nav from "../components/Nav.svelte";
  import { redirect, url } from "@roxi/routify";
  import { TextField, Button, MaterialApp, Alert } from "svelte-materialify";

  let email: string;
  let password: string;
  let error: Error;
  let show: boolean;

  const usernameRules = [(v: any) => !!v || "Required"];
  const passwordRules = [
    (v: any) => !!v || "Required",
    (v: any) => v.length >= 4 || "must be 4 characters",
  ];
  const handleLogin = () => {
    const user = login({ email, password });
    user
      .then(() => $redirect("./home"))
      .catch((err: Error) => {
        error = err;
        show = true;
      });
  };
</script>

<MaterialApp>
  <Nav isLoggedIn={false} />
  <div class="d-flex justify-center mt-4">
    <div class="d-flex flex-column">
      <h3 class="text-h3 mb-6">Login</h3>
      <div style="width: 700px;" class="mb-4 mt-3">
        <TextField
          color="pink lighten-3"
          bind:value={email}
          rules={usernameRules}
        >
          email
        </TextField>
      </div>
      <div style="width: 700px;">
        <TextField
          type="password"
          color="pink lighten-3"
          bind:value={password}
          rules={passwordRules}
        >
          password
        </TextField>
      </div>
      <div style="align-self: center;" class="mt-4">
        <Button size="large" class="pink lighten-3" on:click={handleLogin}>
          Submit
        </Button>
        <a href={$url("./forgot-password")} class="text-decoration-none ml-3">
          forgot password?
        </a>
      </div>
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
</MaterialApp>
