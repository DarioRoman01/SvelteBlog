<script lang="ts">
  import { MaterialApp } from "svelte-materialify";
  import Posts from "../components/Posts.svelte";
  import Nav from "../components/Nav.svelte";
</script>

<MaterialApp>
  <Nav isLoggedIn={true}/>
  <div class="mt-4">
    <Posts />
  </div>
</MaterialApp>