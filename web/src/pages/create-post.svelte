<script lang="ts">
  import {
    TextField,
    Button,
    MaterialApp,
    Alert,
    Textarea
  } from "svelte-materialify";
  import Nav from "../components/Nav.svelte";
  import { redirect } from "@roxi/routify";
  import { CreatePost } from "../requests/posts";

  let title: string;  
  let body: string;
  let error: Error;

  const handleClick = () => {
    const post = CreatePost(title, body);
    post.then(() => $redirect("/home"))
    .catch((err) => error = err);
  };
</script>

<MaterialApp>
  <Nav isLoggedIn={true} />
  <div class="d-flex justify-center mt-4">
    <div class="d-flex flex-column">
      <h3 class="text-h3 mb-6">Add expense</h3>
      <div style="width: 700px;" class="mb-4 mt-3">
        <TextField color="pink lighten-3" bind:value={title}>
          title
        </TextField>
      </div>
      <div style="width: 700px;">
        <Textarea color="pink lighten-3" bind:value={body}>
          body
        </Textarea>
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