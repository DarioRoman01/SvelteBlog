<script lang="ts">
  import { MaterialApp } from "svelte-materialify";
  import { api } from "../requests/users";
  import type { PaginatedComments } from "../requests/comments";
  import { onMount } from "svelte";

  let comments: PaginatedComments;
  export let id: number;

  onMount(async () => {
    comments = await api<PaginatedComments>(`/posts/${id}/comments`);
  });
</script>

<MaterialApp />
