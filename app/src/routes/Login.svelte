<script>

import { userStore, authenticate } from '../store/user'

const urlSearchParams = new URLSearchParams(window.location.search);
const params = Object.fromEntries(urlSearchParams.entries());

let email = params['email'] ?? "";
let account = params['account'] ?? "";
let password = "";
let warningMessage = "";

function handleSubmit() {
    authenticate(account, email, password).then(resp => {
        if (resp.status !== 200) {
            warningMessage = "Invalid credentials, please make sure your account name is correct";
        } else {
            return resp.json().then(user => {
                $userStore = user;

                window.location = "/";
            });
        }
    });
}

</script>

<main class="form-signin">
    <form on:submit|preventDefault={handleSubmit}>
        <h1 class="h3 mb-3 fw-normal">Please sign in</h1>

        {#if warningMessage !== ""}
            <div class="row">
                <div class="col-12">
                    <div class="alert alert-warning" role="alert">
                        {warningMessage}
                    </div>
                </div>
            </div>
        {/if}

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