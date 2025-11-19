import { type Component } from "svelte";
import type { SvelteMap } from "svelte/reactivity";

export type GenerationContext = SvelteMap<string, GenerationState>

export interface ThemeContext {
  toggle: () => void;
}

export type AppContext = {
  spotify: ApiState;
  soundcloud: ApiState;
};

export type ApiState = {
  name: ApiName;
  authed: boolean;
  authUrl: string;
  loading: boolean;
  err: string | undefined;
  icon: Component;
};

export type ApiName = "spotify" | "soundcloud";

export type AuthEvent = {
  name: ApiName;
  err: string | undefined;
  ok: boolean;
};

export type MatchProgress = {
  forId: string;
  total: number;
  progress: number;
  status: GenerateStatus;
};

export type GenerationState = {
  status: GenerateStatus;
  total: number;
  complete: number;
  error: string | undefined;
};

export type GenerateStatus = "running" | "complete" | "error" | "idle";
