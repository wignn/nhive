<script lang="ts">
  import { BookOpen, CheckCircle2, ClipboardList, Library, XCircle } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import type { PageData } from './$types';

  let { data } = $props<{ data: PageData }>();
  let entries = $derived(data.entries);

  const TABS = [
    { key: "reading", label: "Reading", icon: BookOpen },
    { key: "completed", label: "Completed", icon: CheckCircle2 },
    { key: "plan_to_read", label: "Plan to read", icon: ClipboardList },
    { key: "dropped", label: "Dropped", icon: XCircle },
  ] as const;

  let tab = $state<typeof TABS[number]["key"]>("reading");
  let filtered = $derived(entries.filter((e) => e.status === tab));
</script>

<svelte:head>
  <title>My Library — NovelHive</title>
</svelte:head>

<main class="mx-auto w-full max-w-6xl px-4 py-10 md:px-6 2xl:px-8">
  <div class="mb-6">
    <div class="text-xs font-bold uppercase tracking-widest text-brand">Personal</div>
    <h1 class="mt-1 font-display text-3xl font-extrabold md:text-4xl">My library</h1>
    <p class="mt-2 text-sm text-muted-foreground">Track everything you're reading, completed, or planning to start.</p>
  </div>

  <div class="mb-6 grid grid-cols-2 gap-2 rounded-2xl border border-white/10 bg-white/[0.03] p-1.5 sm:grid-cols-4">
    {#each TABS as t}
      {@const count = entries.filter((e) => e.status === t.key).length}
      {@const Icon = t.icon}
      <button
        onclick={() => tab = t.key}
        class={cn(
          "flex items-center justify-center gap-2 rounded-xl px-3 py-2 text-xs font-bold transition sm:text-sm",
          tab === t.key ? "gradient-brand text-white shadow-glow" : "text-muted-foreground hover:text-foreground",
        )}
      >
        <Icon class="h-4 w-4" /> {t.label}
        <span class={cn(
          "rounded-md px-1.5 py-0.5 text-[10px]",
          tab === t.key ? "bg-white/20" : "bg-white/5",
        )}>{count}</span>
      </button>
    {/each}
  </div>

  {#if filtered.length === 0}
    <div class="grid place-items-center rounded-2xl border border-dashed border-white/10 bg-white/[0.02] py-16 text-center">
      <Library class="mb-3 h-10 w-10 text-muted-foreground/60" />
      <p class="text-sm text-muted-foreground">Nothing here yet.</p>
      <a href="/novels" class="mt-4 rounded-xl gradient-brand px-4 py-2 text-sm font-bold text-white shadow-glow">Browse novels</a>
    </div>
  {:else}
    <ul class="space-y-3">
      {#each filtered as e}
        {@const progress = e.progress ?? 0}
        {@const total = e.total ?? null}
        {@const pct = total ? Math.min(100, (progress / total) * 100) : 0}
        
        <li>
          <a
            href={e.novel_slug ? `/novel/${e.novel_slug}` : '#'}
            class="block rounded-2xl border border-white/5 bg-card/60 p-4 transition hover:border-brand/40 hover:bg-card"
          >
            <div class="flex flex-wrap items-center justify-between gap-2">
              <div>
                <div class="font-display text-base font-extrabold">{e.novel_title || "Unknown title"}</div>
                <div class="text-xs text-muted-foreground">
                  Chapter {progress} {total ? `of ${total}` : ""}
                </div>
              </div>
              {#if total}
                <div class="text-xs font-bold text-brand">{Math.round(pct)}%</div>
              {/if}
            </div>
            <div class="mt-3 h-1.5 overflow-hidden rounded-full bg-white/5">
              <div class="h-full gradient-brand" style="width: {pct}%"></div>
            </div>
          </a>
        </li>
      {/each}
    </ul>
  {/if}
</main>
