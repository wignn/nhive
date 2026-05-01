<script lang="ts">
  import { buildCoverUrl, formatCount, cn } from "$lib/utils";
  import { BookmarkCheck, BookmarkPlus, BookOpen, Eye, Loader2, PlayCircle } from "lucide-svelte";
  import type { PageData, ActionData } from './$types';
  import { page } from "$app/stores";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";

  let { data, form } = $props<{ data: PageData, form: ActionData }>();

  let novel = $derived(data.novel);
  let chapters = $derived(data.chapters);
  let coverBase = $derived(data.coverBase);
  let libEntry = $derived(data.libEntry);
  let user = $derived($page.data.user);

  let cover = $derived(buildCoverUrl(novel?.cover_url, coverBase));
  let inLibrary = $derived(
    form?.action === 'removed' ? false :
    form?.action === 'added' ? true :
    !!libEntry
  );

  let libBusy = $state(false);
</script>

<svelte:head>
  <title>{novel?.title || 'Novel'} — NovelHive</title>
</svelte:head>

<main class="mx-auto w-full max-w-[1400px] px-4 py-8 md:px-6 2xl:px-8">
  <!-- Backdrop -->
  <div class="relative mb-8 overflow-hidden rounded-3xl border border-white/5">
    {#if cover}
      <div
        class="absolute inset-0 -z-10 scale-110 opacity-40 blur-2xl"
        style="background-image: url({cover}); background-size: cover; background-position: center;"
        aria-hidden="true"
      ></div>
    {/if}
    <div class="absolute inset-0 -z-10 bg-gradient-to-t from-background via-background/85 to-background/40" aria-hidden="true"></div>

    <div class="grid gap-6 p-5 md:grid-cols-[260px_1fr] md:gap-8 md:p-8">
      <!-- Cover -->
      <div class="mx-auto w-44 md:mx-0 md:w-full">
        <div class="cover-frame rounded-2xl shadow-elevated ring-1 ring-white/10">
          {#if cover}
            <img src={cover} alt={novel.title} class="absolute inset-0 h-full w-full object-cover" />
          {:else}
            <div class="absolute inset-0 grid place-items-center text-muted-foreground/40">
              <BookOpen class="h-14 w-14" />
            </div>
          {/if}
        </div>
      </div>

      <!-- Info -->
      <div>
        <div class="flex flex-wrap items-center gap-2 text-xs">
          {#if novel.status}
            <span class={cn(
              "rounded-md px-2 py-0.5 font-bold uppercase tracking-wider",
              novel.status === "completed" ? "bg-success/20 text-success" :
              novel.status === "hiatus" ? "bg-warning/20 text-warning" :
              "bg-brand/20 text-brand",
            )}>
              {novel.status}
            </span>
          {/if}
          {#if novel.genres}
            {#each novel.genres as g}
              <span class="rounded-md border border-white/10 bg-white/5 px-2 py-0.5 font-semibold text-muted-foreground">
                {g.name}
              </span>
            {/each}
          {/if}
        </div>

        <h1 class="mt-3 font-display text-3xl font-extrabold leading-tight md:text-5xl">
          {novel.title}
        </h1>
        <div class="mt-1 text-sm text-muted-foreground">
          by <span class="font-semibold text-foreground">{novel.author}</span>
        </div>

        <p class="mt-4 max-w-2xl text-sm text-muted-foreground md:text-base">{novel.synopsis}</p>

        <div class="mt-5 flex flex-wrap items-center gap-3">
          {#if chapters.length > 0}
            {@const nextCh = libEntry?.progress > 0 ? libEntry.progress : chapters[0].number}
            <a
              href={`/novel/${novel.slug}/${nextCh}`}
              class="inline-flex items-center gap-2 rounded-xl gradient-brand px-5 py-2.5 text-sm font-bold text-white shadow-glow hover:opacity-90"
            >
              <PlayCircle class="h-4 w-4" />
              {libEntry?.progress > 0 ? `Continue ch. ${libEntry.progress}` : "Start reading"}
            </a>
          {/if}
          <form 
            method="POST" 
            action="?/toggleLibrary"
            use:enhance={() => {
              if (!user) {
                toast.info("Sign in to manage your library");
                return ({ update }) => update({ reset: false }); // Cancel submission
              }
              libBusy = true;
              return async ({ result, update }) => {
                libBusy = false;
                if (result.type === 'success') {
                  toast.success(result.data?.action === 'added' ? "Added to library" : "Removed from library");
                } else if (result.type === 'failure') {
                  toast.error((result.data?.error as string) || "Library update failed");
                }
                update({ reset: false });
              };
            }}
          >
            <input type="hidden" name="novelId" value={novel.id} />
            <input type="hidden" name="action" value={inLibrary ? "remove" : "add"} />
            <button
              disabled={libBusy}
              class={cn(
                "inline-flex items-center gap-2 rounded-xl border px-5 py-2.5 text-sm font-semibold transition",
                inLibrary ? "border-success/50 bg-success/10 text-success" : "border-white/10 bg-white/5 hover:bg-white/10",
              )}
            >
              {#if libBusy}
                <Loader2 class="h-4 w-4 animate-spin" />
              {:else if inLibrary}
                <BookmarkCheck class="h-4 w-4" />
              {:else}
                <BookmarkPlus class="h-4 w-4" />
              {/if}
              {inLibrary ? "In library" : "Add to library"}
            </button>
          </form>
        </div>

        <div class="mt-6 grid max-w-md grid-cols-3 gap-3">
          <div class="rounded-xl border border-white/5 bg-white/[0.03] p-3">
            <div class="font-display text-lg font-extrabold">{novel.total_chapters ?? "—"}</div>
            <div class="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Chapters</div>
          </div>
          <div class="rounded-xl border border-white/5 bg-white/[0.03] p-3">
            <div class="font-display text-lg font-extrabold">{typeof novel.view_count === "number" ? formatCount(novel.view_count) : "—"}</div>
            <div class="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Views</div>
          </div>
          <div class="rounded-xl border border-white/5 bg-white/[0.03] p-3">
            <div class="font-display text-lg font-extrabold capitalize">{novel.status ?? "—"}</div>
            <div class="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Status</div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Chapters -->
  <section>
    <div class="mb-4 flex items-end justify-between">
      <h2 class="font-display text-2xl font-extrabold">Chapters</h2>
      <span class="text-xs text-muted-foreground">{chapters.length} total</span>
    </div>
    {#if chapters.length === 0}
      <div class="rounded-2xl border border-dashed border-white/10 bg-white/[0.02] py-12 text-center text-muted-foreground">
        No chapters published yet.
      </div>
    {:else}
      <ul class="overflow-hidden rounded-2xl border border-white/5 bg-card/60">
        {#each chapters as c, i}
          <li>
            <a
              href={`/novel/${novel.slug}/${c.number}`}
              class={cn(
                "flex items-center gap-4 px-4 py-3 transition hover:bg-white/5 md:px-6",
                i !== chapters.length - 1 && "border-b border-white/5",
              )}
            >
              <span class="grid h-9 w-12 shrink-0 place-items-center rounded-md bg-white/5 text-xs font-bold text-muted-foreground">
                {c.number}
              </span>
              <span class={cn("flex-1 truncate text-sm font-semibold", libEntry?.progress >= c.number && "text-muted-foreground/60")}>
                {c.title}
              </span>
              <span class="hidden text-xs text-muted-foreground md:inline">{c.word_count} words</span>
            </a>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</main>
