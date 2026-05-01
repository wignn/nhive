<script lang="ts">
  import { Sparkles } from "lucide-svelte";

</script>

<footer class="mt-24 border-t border-white/5 bg-background/60">
  <div class="mx-auto max-w-[1600px] px-4 py-12 md:px-6 2xl:px-8">
    <div class="grid grid-cols-2 gap-8 md:grid-cols-4">
      <div class="col-span-2">
        <a href="/" class="flex items-center gap-2 font-display text-lg font-extrabold">
          <img src="/logo.ico" alt="NovelHive Logo" class="h-8 w-8" />
          <span class="gradient-text">NovelHive</span>
        </a>
        <p class="mt-3 max-w-sm text-sm text-muted-foreground">
          Premium destination for translated web novels and original fiction. Read beautifully, anywhere.
        </p>
      </div>
      <div>
        <h4 class="mb-3 text-xs font-bold uppercase tracking-wider text-muted-foreground">Explore</h4>
        <ul class="space-y-2 text-sm">
          <li><a href="/novels" class="hover:text-foreground text-muted-foreground">Browse</a></li>
          <li><a href="/search" class="hover:text-foreground text-muted-foreground">Search</a></li>
          <li><a href="/library" class="hover:text-foreground text-muted-foreground">My Library</a></li>
        </ul>
      </div>
      <div>
        <h4 class="mb-3 text-xs font-bold uppercase tracking-wider text-muted-foreground">Genres</h4>
        <ul class="space-y-2 text-sm">
          <li><a href="/novels?genre=Fantasy" class="hover:text-foreground text-muted-foreground">Fantasy</a></li>
          <li><a href="/novels?genre=Romance" class="hover:text-foreground text-muted-foreground">Romance</a></li>
          <li><a href="/novels?genre=Isekai" class="hover:text-foreground text-muted-foreground">Isekai</a></li>
          <li><a href="/novels?genre=Action" class="hover:text-foreground text-muted-foreground">Action</a></li>
        </ul>
      </div>
    </div>
    <div class="mt-10 flex flex-col items-center justify-between gap-2 border-t border-white/5 pt-6 text-xs text-muted-foreground md:flex-row">
      <span>© {new Date().getFullYear()} NovelHive. Crafted for readers.</span>
      <span>Built with SvelteKit.</span>
    </div>
  </div>
</footer>
