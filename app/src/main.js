import App from './App.svelte';
import {takeover} from "./fetch/fetchm";

// replaces fetch with fetchm that allows middleware
takeover();

let app = new App({
	target: document.body
});

export default app;