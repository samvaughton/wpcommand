<script>
    import {Router, Link, Route} from "svelte-routing";
    import Dashboard from "./routes/Dashboard.svelte";
    import Header from "./layout/Header.svelte";
    import Config from "./routes/Config.svelte";
    import SiteDetails from "./routes/SiteDetails.svelte";
    import Login from "./routes/Login.svelte";
    import {userStore} from './store/user'
    import CommandJobs from "./routes/CommandJobs.svelte";
    import CommandJob from "./routes/CommandJob.svelte";
    import Blueprints from "./routes/Blueprints.svelte";
    import Blueprint from "./routes/Blueprint.svelte";
    import BlueprintObject from "./routes/BlueprintObject.svelte";
    import Accounts from "./routes/Accounts.svelte";
    import Account from "./routes/Account.svelte";
    import User from "./routes/User.svelte";

    export let url = "";

    if (!$userStore) {
        url = "/login"
    }
</script>

<Router url="{url}">
    <Header />
    <div class="container-xl">
        <Route path="/login"><Login /></Route>
        <Route path="/"><Dashboard /></Route>
        <Route path="/config"><Config /></Route>
        <Route path="/logs"><CommandJobs /></Route>
        <Route path="/logs/:uuid" let:params><CommandJob uuid="{params.uuid}" /></Route>
        <Route path="/sites/:key" let:params><SiteDetails key="{params.key}" /></Route>

        <Route path="/blueprints/:bpUuid/object/:objUuid/revision/:revId" let:params>
            <BlueprintObject blueprintUuid="{params.bpUuid}" objectUuid="{params.objUuid}" revisionId="{params.revId}" />
        </Route>

        <Route path="/blueprints/:uuid" let:params><Blueprint uuid="{params.uuid}" /></Route>
        <Route path="/blueprints"><Blueprints /></Route>

        <Route path="/accounts/:accUuid/users/:userUuid" let:params><User accUuid="{params.accUuid}" userUuid="{params.userUuid}" /></Route>
        <Route path="/accounts/:accUuid" let:params><Account uuid="{params.accUuid}" /></Route>
        <Route path="/accounts" let:params><Accounts /></Route>
    </div>
</Router>