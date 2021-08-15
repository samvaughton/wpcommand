<script>

import { userStore, authenticate } from '../store/user'

let email = "sam@rentivo.com"
let account = "admin"
let password = "password"

function handleSubmit() {
    authenticate(account, email, password).then(user => {
        $userStore = user;

        window.location = "/";
    });
}

</script>

<main class="form-signin">
    <form on:submit|preventDefault={handleSubmit}>
        <h1 class="h3 mb-3 fw-normal">Please sign in</h1>

        <div class="form-floating mb-2">
            <input type="text" class="form-control" bind:value={account} id="account" placeholder="Account">
            <label for="email">Account</label>
        </div>
        <div class="form-floating mb-2">
            <input type="email" class="form-control" bind:value={email} id="email" placeholder="Email">
            <label for="email">Email address</label>
        </div>
        <div class="form-floating mb-2">
            <input type="password" class="form-control" bind:value={password} id="password" placeholder="Password">
            <label for="password">Password</label>
        </div>

        <button class="w-100 btn btn-lg btn-primary" type="submit">Sign in</button>
    </form>
</main>