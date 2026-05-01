<script lang="ts">
  import { ArrowLeft, ArrowRight, Heart, List, MessageSquare, Settings, Type, X } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import type { PageData, ActionData } from './$types';
  import { page } from "$app/stores";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";

  let { data, form } = $props<{ data: PageData, form: ActionData }>();

  let chapter = $derived(data.chapter);
  let comments = $derived(data.comments);
  let user = $derived($page.data.user);
  let slug = $derived($page.params.slug);

  type Theme = "dark" | "light" | "sepia";
  let theme = $state<Theme>("dark");
  let fontSize = $state(18);
  let openSettings = $state(false);
  let progress = $state(0);
  let text = $state("");
  let posting = $state(false);

  $effect(() => {
    let timeout: ReturnType<typeof setTimeout>;
    let lastSaved = -1;

    function handleScroll() {
      const h = document.documentElement.scrollHeight - window.innerHeight;
      progress = h > 0 ? Math.min(100, (window.scrollY / h) * 100) : 0;

      clearTimeout(timeout);
      timeout = setTimeout(() => {
        if (user && chapter && Math.round(progress) !== lastSaved) {
          lastSaved = Math.round(progress);
          const fd = new FormData();
          fd.append('novelId', chapter.novel_id);
          fd.append('scrollPosition', lastSaved.toString());
          fetch(`/novel/${slug}/${chapter.number}?/saveProgress`, {
            method: 'POST',
            body: fd,
            headers: { 'x-sveltekit-action': 'true' }
          }).catch(() => {});
        }
      }, 1000);
    }

    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => {
      window.removeEventListener("scroll", handleScroll);
      clearTimeout(timeout);
    };
  });

  let themeClass = $derived(theme === "light" ? "reader-theme-light" : theme === "sepia" ? "reader-theme-sepia" : "");
</script>

<svelte:head>
  <title>{chapter ? `Ch. ${chapter.number} - ${chapter.novel_title}` : 'Chapter'} — NovelHive</title>
</svelte:head>

{#if !chapter}
  <div class="mx-auto max-w-3xl px-4 py-20 text-center">Chapter not found.</div>
{:else}
  <div class={cn("min-h-screen", themeClass)}>
    <div class="sticky top-0 z-40 glass-strong">
      <div class="mx-auto flex h-14 w-full max-w-4xl items-center gap-3 px-4">
        <a href={`/novel/${slug}`} class="inline-flex items-center gap-1.5 text-sm font-semibold text-muted-foreground hover:text-foreground">
          <ArrowLeft class="h-4 w-4" /> Back
        </a>
        <div class="mx-auto truncate text-sm font-semibold text-muted-foreground">
          <span class="text-brand">Ch. {chapter.number}</span> · {chapter.novel_title}
        </div>
        <a href={`/novel/${slug}`} class="rounded-lg p-2 hover:bg-white/5">
          <List class="h-4 w-4" />
        </a>
        <button onclick={() => openSettings = !openSettings} class="rounded-lg p-2 hover:bg-white/5">
          {#if openSettings}
            <X class="h-4 w-4" />
          {:else}
            <Settings class="h-4 w-4" />
          {/if}
        </button>
      </div>
      <div class="h-0.5 w-full bg-white/5">
        <div class="h-full gradient-brand transition-[width]" style="width: {progress}%"></div>
      </div>
    </div>

    <article class="mx-auto max-w-3xl px-4 pb-24 pt-10">
      <div class="text-center">
        <div class="text-xs font-bold uppercase tracking-widest text-brand">Chapter {chapter.number}</div>
        <h1 class="mt-2 font-display text-3xl font-extrabold md:text-4xl">{chapter.title}</h1>
        <div class="mt-2 text-xs text-muted-foreground">{chapter.word_count} words</div>
        <div class="mx-auto mt-6 h-px w-32 bg-gradient-to-r from-transparent via-brand/60 to-transparent"></div>
      </div>

      <div
        class="reader-prose mt-10 leading-relaxed text-foreground"
        style="font-size: {fontSize}px; font-family: var(--font-serif)"
      >
        {#each String(chapter.content || "").split(/\n\n+/) as p}
          <p>{p}</p>
        {/each}
      </div>

      <!-- nav -->
      <div class="mt-12 flex items-center justify-between gap-3 border-t border-white/5 pt-6">
        {#if chapter.has_prev}
          <a
            href={`/novel/${slug}/${chapter.prev_number}`}
            class="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold hover:bg-white/10"
          >
            <ArrowLeft class="h-4 w-4" /> Previous
          </a>
        {:else}
          <span></span>
        {/if}
        <a href={`/novel/${slug}`} class="text-sm font-semibold text-muted-foreground hover:text-foreground">
          All chapters
        </a>
        {#if chapter.has_next}
          <a
            href={`/novel/${slug}/${chapter.next_number}`}
            class="inline-flex items-center gap-2 rounded-xl gradient-brand px-4 py-2 text-sm font-bold text-white shadow-glow hover:opacity-90"
          >
            Next <ArrowRight class="h-4 w-4" />
          </a>
        {:else}
          <span></span>
        {/if}
      </div>

      <!-- comments -->
      <section class="mt-16">
        <h2 class="mb-4 inline-flex items-center gap-2 font-display text-xl font-extrabold">
          <MessageSquare class="h-5 w-5 text-brand" /> Comments ({comments.length})
        </h2>
        <div class="rounded-2xl border border-white/10 bg-white/[0.03] p-4">
          <form
            method="POST"
            action="?/postComment"
            use:enhance={() => {
              posting = true;
              return async ({ result, update }) => {
                posting = false;
                if (result.type === 'success') {
                  toast.success("Comment posted");
                  text = "";
                  invalidateAll();
                } else if (result.type === 'failure') {
                  toast.error((result.data?.error as string) || "Failed to post comment");
                }
                update({ reset: false });
              };
            }}
          >
            <input type="hidden" name="chapterId" value={chapter.id} />
            <textarea
              disabled={!user}
              name="content"
              bind:value={text}
              placeholder={user ? "Share your thoughts on this chapter…" : "Sign in to comment"}
              class="min-h-20 w-full resize-y rounded-lg border border-white/10 bg-background/40 p-3 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30"
            ></textarea>
            <div class="mt-2 flex justify-end">
              <button
                disabled={!user || posting || !text.trim()}
                class="rounded-lg gradient-brand px-4 py-2 text-sm font-bold text-white shadow-glow disabled:opacity-50"
              >
                {posting ? "Posting…" : "Post"}
              </button>
            </div>
          </form>
        </div>

        <div class="mt-4 divide-y divide-white/5">
          {#each comments as c}
            <div class="flex gap-3 py-4">
              <div class="grid h-10 w-10 shrink-0 place-items-center rounded-full gradient-brand text-sm font-bold text-white">
                {c.username?.[0]?.toUpperCase() || "?"}
              </div>
              <div class="flex-1">
                <div class="text-sm font-bold">{c.username}</div>
                <p class="mt-1 text-sm text-muted-foreground">{c.content}</p>
                <div class="mt-2 flex items-center gap-3 text-xs text-muted-foreground">
                  <button class="inline-flex items-center gap-1 hover:text-foreground">
                    <Heart class="h-3.5 w-3.5" /> {c.likes || 0}
                  </button>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </section>
    </article>

    <!-- Floating settings -->
    {#if openSettings}
      <div class="fixed bottom-6 right-6 z-50 w-72 rounded-2xl border border-white/10 bg-popover/95 p-4 shadow-elevated backdrop-blur fade-up">
        <div class="mb-3 text-xs font-bold uppercase tracking-widest text-muted-foreground">Theme</div>
        <div class="mb-4 grid grid-cols-3 gap-2">
          {#each ["dark", "sepia", "light"] as t}
            <button
              onclick={() => theme = t as Theme}
              class={cn(
                "rounded-lg border px-3 py-2 text-xs font-bold capitalize",
                theme === t ? "border-brand bg-brand/15 text-brand" : "border-white/10 bg-white/5",
              )}
            >
              {t}
            </button>
          {/each}
        </div>
        <div class="mb-2 inline-flex items-center gap-1 text-xs font-bold uppercase tracking-widest text-muted-foreground">
          <Type class="h-3 w-3" /> Font size
        </div>
        <div class="flex items-center gap-2 rounded-lg border border-white/10 bg-white/5 p-1">
          <button onclick={() => fontSize = Math.max(14, fontSize - 2)} class="flex-1 rounded-md py-1 text-sm font-bold hover:bg-white/10">A−</button>
          <span class="w-12 text-center text-sm font-semibold">{fontSize}px</span>
          <button onclick={() => fontSize = Math.min(28, fontSize + 2)} class="flex-1 rounded-md py-1 text-sm font-bold hover:bg-white/10">A+</button>
        </div>
      </div>
    {/if}
  </div>
{/if}
