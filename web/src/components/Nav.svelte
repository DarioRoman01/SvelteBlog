<script lang="ts">
  import { goto, url } from '@roxi/routify';
  import { logout, api, User } from "../requests/users"
  import { AppBar, Button, Icon, MaterialApp } from 'svelte-materialify';
  import { mdiPlusCircle, mdiArrowRightCircle, mdiAccountCircle  } from '@mdi/js';

  let loggedIn: boolean;
  loggedIn = false;

  const handleLogout = () => {
    const out = logout();
    out.then(() => loggedIn = false);
  }

  const me = api<User>(`${process.env.API_URL}/me`);
  me.then(() => loggedIn = true);

</script>

<MaterialApp>
  <AppBar>
    <span slot="title">
      <a href={$url("./index")} 
        class="text-decoration-none"
      >
        .Blog
      </a>
    </span>
    <div style="flex-grow:1"/>
    {#if loggedIn}
      <Button class="primary-color" on:click={$goto(`./profile`)}>
        <Icon path={mdiAccountCircle}/>
      </Button>
      <Button class="primary-color" on:click={$goto("./create-expense")}>
        <Icon path={mdiPlusCircle}/>
      </Button>
      <Button class="ml-3 mr-1 primary-color" on:click={handleLogout}>
        <Icon path={mdiArrowRightCircle}/>
      </Button>
    {:else}
      <Button class="primary-color" on:click={$goto("./login")}>
        login
      </Button>
      <Button class="ml-3 mr-1 primary-color" on:click={$goto("./register")}>
        register
      </Button>
    {/if}
  </AppBar>
</MaterialApp>