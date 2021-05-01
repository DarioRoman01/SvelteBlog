<script lang="ts">
  import { onMount } from "svelte";
  import { params, redirect } from "@roxi/routify";
  import { MaterialApp, ProgressCircular } from "svelte-materialify";
  import { verify } from "../../requests/users";

  let error: Error;

  onMount(() => {
    const verification = verify($params["token"]);
    verification.then(() => $redirect("../create-profile"))
    .catch((err: Error) => error =  err);
  });
</script>

<MaterialApp>
  {#if error}
    <div class="mt-4 d-flex justify-center">
      <h1 class="text-h1">Invalid token</h1>
    </div>
  {:else}
    <ProgressCircular size={50} indeterminate color="pink lighten-3" />
  {/if}
</MaterialApp>