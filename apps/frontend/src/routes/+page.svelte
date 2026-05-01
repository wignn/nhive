<script lang="ts">
  import NovelCard from "$lib/components/novel/NovelCard.svelte";
  import { buildCoverUrl, formatCount } from "$lib/utils";
  import { ArrowRight, BookOpen, Eye, Flame, Sparkles, TrendingUp } from "lucide-svelte";
  import type { PageData } from './$types';

  let { data } = $props<{ data: PageData }>();

  let novels = $derived(data.recent);
  let trending = $derived(data.trending);
  let coverBase = $derived(data.coverBase);

  let hero = $derived(trending[0] || novels[0]);
  let heroCover = $derived(hero ? buildCoverUrl(hero.cover_url, coverBase) : null);
</script>

<svelte:head>
  <title>NovelHive — Discover your next obsession</title>
</svelte:head>

<main class="mx-auto w-full max-w-[1600px] px-4 md:px-6 2xl:px-8">
  <!-- HERO -->
  <section class="relative isolate overflow-hidden rounded-3xl border border-white/5 bg-surface/40 p-6 md:p-10 mt-6">
    <div class="absolute -top-32 right-0 h-96 w-96 shimmer-bg opacity-80" aria-hidden="true"></div>
    <div class="absolute -bottom-32 -left-20 h-80 w-80 shimmer-bg opacity-60" aria-hidden="true"></div>

    <div class="relative grid gap-8 md:grid-cols-[1.1fr_.9fr] md:items-center">
      <div class="fade-up">
        <div class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs font-semibold text-muted-foreground">
          <span class="h-1.5 w-1.5 rounded-full bg-brand pulse-dot"></span> Welcome back to NovelHive
        </div>
        <h1 class="mt-4 font-display text-4xl font-extrabold leading-tight text-balance md:text-6xl">
          Worlds within <span class="gradient-text">words</span>.<br />
          Find your next favorite chapter.
        </h1>
        <p class="mt-4 max-w-lg text-base text-muted-foreground md:text-lg">
          Hand-picked translations, original fiction and trending hits — wrapped in a reader you'll actually love using.
        </p>
        <div class="mt-6 flex flex-wrap items-center gap-3">
          <a href="/novels" class="inline-flex items-center gap-2 rounded-xl gradient-brand px-5 py-3 text-sm font-bold text-white shadow-glow transition hover:opacity-90">
            Start reading <ArrowRight class="h-4 w-4" />
          </a>
          <a href="/search" class="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-5 py-3 text-sm font-semibold text-foreground hover:bg-white/10">
            <Sparkles class="h-4 w-4 text-brand" /> Discover
          </a>
        </div>
      </div>

      <!-- Featured floating cover -->
      {#if hero}
        <a
          href={`/novel/${hero.slug}`}
          class="relative mx-auto w-full max-w-md fade-up"
          style="animation-delay: 120ms"
        >
          <div class="absolute inset-0 -z-10 translate-x-6 translate-y-6 rounded-2xl gradient-brand opacity-40 blur-2xl" aria-hidden="true"></div>
          <div class="grid grid-cols-[140px_1fr] gap-4 overflow-hidden rounded-2xl border border-white/10 bg-card/70 p-3 backdrop-blur md:grid-cols-[180px_1fr]">
            <div class="cover-frame rounded-xl ring-1 ring-white/10">
              {#if heroCover}
                <img src={heroCover} alt={hero.title} class="absolute inset-0 h-full w-full object-cover" />
              {:else}
                <div class="absolute inset-0 grid place-items-center text-muted-foreground/40"><BookOpen class="h-8 w-8" /></div>
              {/if}
            </div>
            <div class="flex flex-col justify-between p-1">
              <div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-brand">Editor's pick</div>
                <div class="mt-1 line-clamp-2 font-display text-xl font-extrabold leading-tight md:text-2xl">{hero.title}</div>
                <div class="mt-1 text-xs text-muted-foreground">by {hero.author || "Unknown"}</div>
                <p class="mt-3 line-clamp-3 text-sm text-muted-foreground">{hero.synopsis || "A story waiting to be told."}</p>
              </div>
              <div class="mt-3 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                {#if hero.genres?.[0]}
                  <span class="rounded-md border border-white/10 bg-white/5 px-2 py-0.5">{hero.genres[0].name}</span>
                {/if}
                {#if typeof hero.total_chapters === "number"}
                  <span class="flex items-center gap-1"><BookOpen class="h-3 w-3" />{hero.total_chapters} ch</span>
                {/if}
                {#if typeof hero.view_count === "number"}
                  <span class="flex items-center gap-1"><Eye class="h-3 w-3" />{formatCount(hero.view_count)}</span>
                {/if}
              </div>
            </div>
          </div>
        </a>
      {/if}
    </div>
  </section>

  <!-- TRENDING marquee strip -->
  {#if trending.length > 0}
    <section class="mt-10">
      <div class="flex items-end justify-between gap-3">
        <div>
          <div class="inline-flex items-center gap-2 text-xs font-bold uppercase tracking-widest text-brand">
            <Flame class="h-4 w-4" /> Trending now
          </div>
          <div class="mt-1 font-display text-xl font-extrabold md:text-2xl">Most read this week</div>
        </div>
        <a href="/novels" class="group inline-flex items-center gap-1 text-sm font-semibold text-muted-foreground hover:text-foreground">
          View all <ArrowRight class="h-3.5 w-3.5 transition group-hover:translate-x-0.5" />
        </a>
      </div>
      <div class="-mx-4 mt-3 overflow-hidden md:-mx-6">
        <div class="scrollbar-none flex snap-x snap-mandatory gap-3 overflow-x-auto px-4 pb-2 md:gap-5 md:px-6">
          {#each trending as n, i}
            <div class="snap-start shrink-0 w-[140px] sm:w-[160px] md:w-[180px]">
              <NovelCard novel={n} {coverBase} index={i} />
            </div>
          {/each}
        </div>
      </div>
    </section>
  {/if}

  <!-- RECENT GRID -->
  <section class="mt-12">
    <div class="flex items-end justify-between gap-3">
      <div>
        <div class="inline-flex items-center gap-2 text-xs font-bold uppercase tracking-widest text-brand">
          <TrendingUp class="h-4 w-4" /> Recently updated
        </div>
        <div class="mt-1 font-display text-xl font-extrabold md:text-2xl">Fresh chapters from authors you love</div>
      </div>
      <a href="/novels" class="group inline-flex items-center gap-1 text-sm font-semibold text-muted-foreground hover:text-foreground">
        View all <ArrowRight class="h-3.5 w-3.5 transition group-hover:translate-x-0.5" />
      </a>
    </div>
    
    <div class="mt-3 grid grid-cols-2 gap-3 sm:grid-cols-3 md:gap-5 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-8">
      {#each novels as n, i}
        <NovelCard novel={n} {coverBase} index={i} />
      {/each}
    </div>
  </section>
</main>
