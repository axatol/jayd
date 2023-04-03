export const config = {
  auth0: {
    enabled: import.meta.env.VITE_AUTH0_ENABLED === "true",
    domain: import.meta.env.VITE_AUTH0_DOMAIN,
    clientId: import.meta.env.VITE_AUTH0_CLIENT_ID,
    callbackUri: window.location.origin,
    logoutUri: window.location.origin,
  },
  api: {
    baseUrl: import.meta.env.VITE_API_URL,
    audience: import.meta.env.VITE_API_AUDIENCE,
    scope: ["read:youtube_metadata", "create:youtube_download"].join(" "),
  },
};
