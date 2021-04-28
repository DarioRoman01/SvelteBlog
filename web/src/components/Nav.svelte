<script lang="ts">
  import { goto, url } from '@roxi/routify';
  import { AppBar, Button, Icon, MaterialApp } from 'svelte-materialify';
  import { mdiPlusCircle, mdiArrowRightCircle, mdiAccountCircle  } from '@mdi/js';

  async function logout() {
    const res = await fetch("http://localhost:1323/logout", {
      method: "post",
      credentials: "include"
    });
    if (res.ok) {
      $goto("./index");
    }
  }
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
    <Button class="primary-color" on:click={$goto(`./profile`)}>
      <Icon path={mdiAccountCircle}/>
    </Button>
    <Button class="primary-color" on:click={$goto("./create-expense")}>
      <Icon path={mdiPlusCircle}/>
    </Button>
    <Button class="ml-3 mr-1 primary-color" on:click={logout}>
      <Icon path={mdiArrowRightCircle}/>
    </Button>
  </AppBar>
</MaterialApp>