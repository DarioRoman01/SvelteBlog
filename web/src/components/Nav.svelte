<script lang="ts">
  import { goto, redirect, url } from "@roxi/routify";
  import { logout, api } from "../requests/users";
  import type { Profile } from "../requests/profile";
  import {
    AppBar,
    Button,
    Icon,
    MaterialApp,
  } from "svelte-materialify";
  import {
    mdiPlusCircle,
    mdiArrowRightCircle,
    mdiAccountCircle,
    mdiHome,
  } from "@mdi/js";
  import { onMount } from "svelte";

  export let isLoggedIn: boolean;
  let user: Profile;

  const handleLogout = () => {
    const out = logout();
    out.then(() => $redirect("./index"));
  };

  if (isLoggedIn) {
    onMount(async () => {
      user = await api<Profile>("/me");
    });
  }
</script>

<MaterialApp>
  <AppBar class="pink darken-3 p4">
    <span slot="title">
      <Button
        fab
        size="small"
        class="pink lighten-3"
        on:click={$goto(isLoggedIn ? "/home" : "index")}
      >
        <Icon path={mdiHome} />
      </Button>
    </span>
    <div style="flex-grow:1" />
    {#if isLoggedIn}
      <Button
        fab
        size="small"
        class="pink lighten-3"
        on:click={$goto(`/profile/${user.username}`)}
      >
        <Icon path={mdiAccountCircle} />
      </Button>
      <Button
        fab
        size="small"
        class="ml-3 pink lighten-3"
        on:click={$goto("/create-expense")}
      >
        <Icon path={mdiPlusCircle} />
      </Button>
      <Button
        fab
        size="small"
        class="ml-3 mr-1 pink lighten-3"
        on:click={handleLogout}
      >
        <Icon path={mdiArrowRightCircle} />
      </Button>
    {:else}
      <Button class="pink lighten-3" on:click={$goto("./login")}>login</Button>
      <Button class="ml-3 mr-1 pink lighten-3" on:click={$goto("./register")}>
        register
      </Button>
    {/if}
  </AppBar>
</MaterialApp>
