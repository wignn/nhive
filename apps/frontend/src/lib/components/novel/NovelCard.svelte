<script lang="ts">
  import { buildCoverUrl, formatCount } from "$lib/utils";
  import { BookOpen, Eye } from "lucide-svelte";

  export type Novel = {
    id: string;
    slug: string;
    title: string;
    author?: string;
    synopsis?: string;
    cover_url?: string;
    status?: string;
    total_chapters?: number;
    view_count?: number;
    genres?: { id: string; name: string }[];
  };

  let { novel, coverBase, index = 0 } = $props<{ novel: Novel; coverBase?: string; index?: number }>();
  let cover = $derived(buildCoverUrl(novel.cover_url, coverBase));
</script>

<a
  href={`/novel/${novel.slug}`}
  class="group relative flex flex-col fade-up"
  style={`animation-delay: ${Math.min(index * 40, 400)}ms`}
>
  <div class="cover-frame rounded-xl ring-1 ring-white/5 hover-lift">
    {#if cover}
      <img
        src={cover}
        alt={novel.title}
        loading="lazy"
        class="absolute inset-0 h-full w-full object-cover transition-transform duration-700 group-hover:scale-110"
      />
    {:else}
      <div class="absolute inset-0 flex items-center justify-center text-muted-foreground/40">
        <BookOpen class="h-10 w-10" />
      </div>
    {/if}

    {#if novel.status}
      <span
        class={`absolute left-2 top-2 z-10 rounded-md px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider backdrop-blur ${
          novel.status === "completed"
            ? "bg-success/85 text-background"
            : novel.status === "hiatus"
              ? "bg-warning/85 text-background"
              : "bg-brand/85 text-white"
        }`}
      >
        {novel.status}
      </span>
    {/if}

    {#if novel.genres?.[0]}
      <span class="absolute right-2 top-2 z-10 rounded-md bg-black/55 px-2 py-0.5 text-[10px] font-semibold text-white/90 backdrop-blur">
        {novel.genres[0].name}
      </span>
    {/if}

    <div class="absolute bottom-0 left-0 right-0 z-10 p-3">
      <div class="line-clamp-2 text-sm font-bold leading-tight text-white drop-shadow">
        {novel.title}
      </div>
      <div class="mt-1 flex items-center justify-between text-[11px] text-white/70">
        <span class="line-clamp-1">{novel.author || "Unknown"}</span>
        <span class="flex items-center gap-2">
          {#if typeof novel.total_chapters === "number"}
            <span class="flex items-center gap-1">
              <BookOpen class="h-3 w-3" /> {novel.total_chapters}
            </span>
          {/if}
          {#if typeof novel.view_count === "number"}
            <span class="flex items-center gap-1">
              <Eye class="h-3 w-3" /> {formatCount(novel.view_count)}
            </span>
          {/if}
        </span>
      </div>
    </div>
  </div>
</a>
