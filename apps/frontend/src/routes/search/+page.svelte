<script lang="ts">
  import { Search, SearchX } from "lucide-svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  
  let initial = $derived($page.url.searchParams.get("q") || "");
  let q = $state("");
  let results = $state<any[]>([]);
  let suggestions = $state<any[]>([]);
  let loading = $state(false);
  let searched = $state(false);

  $effect(() => {
    q = initial;
    if (initial) {
      doSearch(initial);
    } else {
      results = [];
      searched = false;
    }
  });

  async function searchApi(val: string) {
    const res = await fetch(`/api/search?q=${encodeURIComponent(val)}`);
    if (!res.ok) throw new Error("Search failed");
    return res.json();
  }

  async function autocompleteApi(val: string) {
    const res = await fetch(`/api/search/autocomplete?q=${encodeURIComponent(val)}`);
    if (!res.ok) throw new Error("Autocomplete failed");
    return res.json();
  }

  async function doSearch(val: string) {
    if (!val.trim()) return;
    loading = true;
    searched = true;
    suggestions = [];
    try {
      const d: any = await searchApi(val);
      results = d?.hits || d?.novels || [];
    } catch {
      results = [];
    } finally {
      loading = false;
    }
  }

  async function onInput(val: string) {
    q = val;
    if (val.length >= 2) {
      try {
        const d: any = await autocompleteApi(val);
        suggestions = d?.suggestions || [];
      } catch {
        suggestions = [];
      }
    } else {
      suggestions = [];
    }
  }

  function submit(e: Event) {
    e.preventDefault();
    goto(`?q=${encodeURIComponent(q)}`, { keepFocus: true });
  }
</script>

<svelte:head>
  <title>Search — NovelHive</title>
</svelte:head>

<main class="mx-auto w-full max-w-3xl px-4 py-10 md:px-6">
  <div class="mb-2 text-xs font-bold uppercase tracking-widest text-brand">Find</div>
  <h1 class="mb-6 font-display text-3xl font-extrabold md:text-4xl">Search the library</h1>

  <form onsubmit={submit} class="relative">
    <div class="flex items-center gap-3 rounded-2xl border border-white/10 bg-white/5 px-4 py-3 focus-within:border-brand/60 focus-within:ring-2 focus-within:ring-brand/30">
      <Search class="h-5 w-5 text-muted-foreground" />
      <input
        value={q}
        oninput={(e) => onInput(e.currentTarget.value)}
        placeholder="Title, author, or keyword…"
        class="w-full bg-transparent text-base placeholder:text-muted-foreground focus:outline-none"
        autofocus
      />
      <button type="submit" class="rounded-lg gradient-brand px-4 py-2 text-sm font-bold text-white shadow-glow">
        Search
      </button>
    </div>
    {#if suggestions.length > 0}
      <div class="absolute left-0 right-0 top-full mt-2 overflow-hidden rounded-xl border border-white/10 bg-popover/95 shadow-elevated backdrop-blur z-50">
        {#each suggestions.slice(0, 6) as s}
          <a
            href={`/novel/${s.slug}`}
            class="block px-4 py-2 text-sm hover:bg-white/5"
            onclick={() => suggestions = []}
          >
            {s.title}
          </a>
        {/each}
      </div>
    {/if}
  </form>

  <div class="mt-8">
    {#if loading}
      <div class="space-y-3">
        {#each Array(4) as _, i}
          <div class="h-4 w-full max-w-full animate-pulse rounded bg-white/5" style="width: {70 + i * 5}%"></div>
        {/each}
      </div>
    {/if}

    {#if !loading && searched && results.length === 0}
      <div class="grid place-items-center rounded-2xl border border-dashed border-white/10 bg-white/[0.02] py-20 text-center text-muted-foreground">
        <SearchX class="mb-3 h-10 w-10 opacity-60" />
        <p>No results for &ldquo;{q}&rdquo;</p>
      </div>
    {/if}

    {#if !loading && results.length > 0}
      <div class="space-y-3">
        <p class="text-xs text-muted-foreground">{results.length} result{results.length > 1 ? "s" : ""}</p>
        {#each results as n}
          <a
            href={`/novel/${n.slug}`}
            class="group block rounded-2xl border border-white/5 bg-card/60 p-4 transition hover:border-brand/40 hover:bg-card"
          >
            <div class="font-display text-lg font-bold group-hover:gradient-text">{n.title}</div>
            <div class="mt-1 text-xs text-muted-foreground">by {n.author}</div>
            {#if n.synopsis}
              <p class="mt-2 line-clamp-2 text-sm text-muted-foreground">{n.synopsis}</p>
            {/if}
            {#if n.genres?.length > 0}
              <div class="mt-3 flex flex-wrap gap-1.5">
                {#each n.genres as g}
                  <span class="rounded-md border border-white/10 bg-white/5 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
                    {g.name}
                  </span>
                {/each}
              </div>
            {/if}
          </a>
        {/each}
      </div>
    {/if}
  </div>
</main>
