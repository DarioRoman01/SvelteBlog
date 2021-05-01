<script lang="ts">
  import {TextField, Button, MaterialApp, Alert } from 'svelte-materialify';
  import Nav from "../components/Nav.svelte";
  import { redirect } from "@roxi/routify";
  import { register, emailRegex, CustomError } from "../requests/users";

  let phoneNumber: string;
  let email: string;
  let password: string;
  let passwordConfirmation: string;
  let error: CustomError;
  let show: boolean;
  let success: boolean;
  success = false;

  const emailRules = [
    (v: string) => !!v || 'Required',
    (v: string) => {
      const pattern = emailRegex;
      return pattern.test(v) || 'Invalid e-mail.';
    },
  ];

  const usernameRules = [
    (v: string) => !!v || 'Required',
    (v: string) => v.length >= 3 || 'must be at least 3 characters',
  ];

  const passwordRules = [
    (v: string) => !!v || 'Required',
    (v: string) => v.length >= 4 || 'must be at least 4 characters',
  
  ];

  const handleRegister = () => {
    const user = register(email, phoneNumber, password, passwordConfirmation);
    user.then(() => success = true).catch((err: CustomError) => {
      error = err
      show = true;
    });
  };
</script>

<MaterialApp>
  {#if !success}
    <Nav isLoggedIn={false}/>
    <div class="d-flex justify-center mt-4">
      <div class="d-flex flex-column">
        <h3 class="text-h3 mb-6">Register</h3>
        <div style="width: 700px;" class="mb-4 mt-3">
          <TextField
            color="pink lighten-3"
            error={error && error.field === "email" ? true : false}
            bind:value={email} 
            rules={emailRules}
          >
            email
          </TextField>
        </div>
        <div style="width: 700px;" class="mb-4 mt-3">
          <TextField
            color="pink lighten-3" 
            error={error && error.field === "phoneNumber" ? true : false}
            bind:value={phoneNumber}
            rules={usernameRules}>
            phone number
          </TextField>
        </div>
        <div style="width: 700px;" class="mb-4 mt-3">
          <TextField 
            color="pink lighten-3"
            type="password"
            error={error && error.field === "password" ? true : false} 
            rules={passwordRules} 
            bind:value={password}
          >
            password
          </TextField>
        </div>
        <div style="width: 700px;" class="mb-4 mt-3">
          <TextField
            color="pink lighten-3"
            type="password"
            error={error && error.field === "password" ? true : false} 
            rules={passwordRules} 
            bind:value={passwordConfirmation}
          >
            password confirmation
          </TextField>
        </div>
        <div style="align-self: center;" class="mt-4">
          <Button size="large" class="pink lighten-3" on:click={handleRegister}>
            Submit
          </Button>
        </div>
        {#if error}
          <Alert class="pink accent-3 mt-4" dismissible={true} bind:visible={show}>
            {error.message}
          </Alert>
        {/if}
      </div>
    </div>
  {:else}
    <div class="d-flex justify-center">
      <div>
        <h2 class="text-h2">
          We send you an email
        </h2>
        <h3 class="text-h3">
          Please check your email to verify your account
        </h3>
      </div>
    </div>
  {/if}
</MaterialApp>