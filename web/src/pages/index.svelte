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
					Welcome to .Blog
				</h1>
			</div>
			<div class="d-flex justify-center mt-4">
				<h3 class="text-h3">
					✨The best site to share your ideas✨
				</h3>
			</div>
			<div class="d-flex justify-center ">
				<div class="mt-5 mr-4">
					<Button class="primary-color" on:click={$redirect("./login")}>
						Login
					</Button>
				</div>
				<div class="mt-5 ml-4">
					<Button class="primary-color" on:click={$redirect("./register")}>
						Register
					</Button>
				</div>
			</div>
		</div>
	</div>
</MaterialApp>
