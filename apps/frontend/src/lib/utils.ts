import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function classNames(...classes: (string | undefined | null | false)[]) {
  return classes.filter(Boolean).join(" ");
}

export function buildCoverUrl(coverUrl?: string | null, base?: string | null): string | null {
  if (!coverUrl) return null;
  if (coverUrl.startsWith("http")) return coverUrl;
  if (!base) return coverUrl;
  return `${base.replace(/\/$/, "")}/${coverUrl.replace(/^\//, "")}`;
}

export function formatCount(n?: number) {
  if (!n && n !== 0) return "—";
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
  if (n >= 1_000) return `${(n / 1000).toFixed(1)}k`;
  return String(n);
}
