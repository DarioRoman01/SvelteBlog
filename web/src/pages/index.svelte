<script lang="ts">
    import { MaterialApp, Button } from "svelte-materialify"
		import { redirect } from "@roxi/routify";
		import Nav from "../components/Nav.svelte"
</script>

<MaterialApp>
	<Nav isLoggedIn={false}/>
	<div class="d-flex flex-column mt-6">
		<div class="">
			<div class="d-flex justify-center">
				<h1 class="text-h1">
					Welcome to SvelteBlog
				</h1>
			</div>
			<div class="d-flex justify-center mt-4">
				<h3 class="text-h3">
					✨The best site to share your ideas✨
				</h3>
			</div>
		</div>
	</div>
</MaterialApp>
