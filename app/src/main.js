import App from './App.svelte';
import {takeover} from "./fetch/fetchm";
import {Debug} from "./util/debug";

// replaces fetch with fetchm that allows middleware
takeover();

window.wpcmd = Debug;

let app = new App({
	target: document.body
});

export default app;