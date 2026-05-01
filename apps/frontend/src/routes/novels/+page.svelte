<script lang="ts">
  import NovelCard from "$lib/components/novel/NovelCard.svelte";
  import { Filter, SearchX } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import { goto } from "$app/navigation";
  import type { PageData } from './$types';
  import { page } from "$app/stores";
  import { navigating } from "$app/stores";
  import Skeletons from "$lib/components/ui-bits/Skeletons.svelte";

  let { data } = $props<{ data: PageData }>();

  const GENRES = [
    "Fantasy","Action","Romance","Adventure","Sci-Fi","Mystery","Horror",
    "Comedy","Drama","Slice of Life","Martial Arts","Isekai","Wuxia","Xianxia",
  ];
  const SORTS = [
    { value: "updated", label: "Recently updated" },
    { value: "views", label: "Most viewed" },
    { value: "title", label: "Title A–Z" },
    { value: "chapters", label: "Most chapters" },
  ];

  let genre = $derived(data.genre);
  let sort = $derived(data.sort);
  let novels = $derived(data.novels);
  let total = $derived(data.total);
  let coverBase = $derived(data.coverBase);

  function updateParams(newGenre: string, newSort: string) {
    const params = new URLSearchParams();
    if (newGenre) params.set("genre", newGenre);
    params.set("sort", newSort);
    goto(`?${params.toString()}`, { keepFocus: true });
  }

  let loading = $derived($navigating != null);
</script>

<svelte:head>
  <title>Browse novels — NovelHive</title>
</svelte:head>

<main class="mx-auto w-full max-w-[1600px] px-4 py-8 md:px-6 2xl:px-8">
  <div class="mb-6 flex flex-wrap items-end justify-between gap-3">
    <div>
      <div class="text-xs font-bold uppercase tracking-widest text-brand">Library</div>
      <h1 class="mt-1 font-display text-3xl font-extrabold md:text-4xl">Browse novels</h1>
    </div>
    <div class="flex items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-2">
      <Filter class="h-4 w-4 text-muted-foreground" />
      <select
        value={sort}
        onchange={(e) => updateParams(genre, (e.target as HTMLSelectElement).value)}
        class="bg-transparent text-sm font-medium focus:outline-none"
      >
        {#each SORTS as s}
          <option value={s.value} class="bg-background">{s.label}</option>
        {/each}
      </select>
    </div>
  </div>

  <!-- genre chips -->
  <div class="mb-6 -mx-4 overflow-x-auto px-4 scrollbar-none md:mx-0 md:px-0">
    <div class="flex gap-2">
      <button
        onclick={() => updateParams("", sort)}
        class={cn(
          "shrink-0 rounded-full border px-3.5 py-1.5 text-xs font-semibold transition",
          !genre
            ? "border-transparent gradient-brand text-white shadow-glow"
            : "border-white/10 bg-white/5 text-muted-foreground hover:bg-white/10 hover:text-foreground",
        )}
      >
        All
      </button>
      {#each GENRES as g}
        <button
          onclick={() => updateParams(g, sort)}
          class={cn(
            "shrink-0 rounded-full border px-3.5 py-1.5 text-xs font-semibold transition",
            genre === g
              ? "border-transparent gradient-brand text-white shadow-glow"
              : "border-white/10 bg-white/5 text-muted-foreground hover:bg-white/10 hover:text-foreground",
          )}
        >
          {g}
        </button>
      {/each}
    </div>
  </div>

  <p class="mb-4 text-xs text-muted-foreground">
    {loading ? "Loading…" : `Showing ${novels.length} of ${total} novels${genre ? ` in ${genre}` : ""}`}
  </p>

  {#if loading}
    <Skeletons count={12} />
  {:else if novels.length === 0}
    <div class="grid place-items-center rounded-2xl border border-dashed border-white/10 bg-white/[0.02] py-20 text-center text-muted-foreground">
      <SearchX class="mb-3 h-10 w-10 opacity-60" />
      <p>No novels found{genre ? ` in "${genre}"` : ""}</p>
    </div>
  {:else}
    <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:gap-5 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-8">
      {#each novels as n, i}
        <NovelCard novel={n} {coverBase} index={i} />
      {/each}
    </div>
  {/if}
</main>
