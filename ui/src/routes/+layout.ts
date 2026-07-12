// Single-page app: everything renders in the browser, the Go backend
// serves the static build. Do not turn SSR on - the backend has no
// Node runtime.
export const ssr = false;
export const prerender = false;
