import { getContext, setContext, type Component } from "svelte";

export type ApiState = {
  name: ApiName;
  authed: boolean;
  authUrl: string;
  loading: boolean;
  err: string | undefined;
  icon: Component;
};

export type AppState = {
  spotify: ApiState;
  soundcloud: ApiState;
};

export type ApiName = "spotify" | "soundcloud";

export type AuthEvent = {
  name: ApiName;
  err: string | undefined;
  ok: boolean;
};
