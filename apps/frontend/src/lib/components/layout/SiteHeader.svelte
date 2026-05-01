<script lang="ts">
  import { Search, BookOpen, Library, Shield, LogOut, Menu, X, ChevronDown } from "lucide-svelte";

  import { page } from "$app/stores";
  import { cn } from "$lib/utils";
  import { goto } from "$app/navigation";

  const NAV = [
    { to: "/", label: "Home" },
    { to: "/novels", label: "Browse" },
    { to: "/search", label: "Search" },
  ];

  let q = $state("");
  let suggestions = $state<any[]>([]);
  let showSugg = $state(false);
  let openMenu = $state(false);
  let openMobile = $state(false);
  
  let user = $derived($page.data.user);
  let path = $derived($page.url.pathname);
  let isReader = $derived(/\/novel\/[^/]+\/\d+/.test(path));

  let ref: HTMLElement | undefined = $state();

  function onClickOutside(e: MouseEvent) {
    if (ref && !ref.contains(e.target as Node)) {
      openMenu = false;
    }
  }

  $effect(() => {
    document.addEventListener("mousedown", onClickOutside);
    return () => document.removeEventListener("mousedown", onClickOutside);
  });

  $effect(() => {
    if (q.length < 2) {
      suggestions = [];
      return;
    }
    const t = setTimeout(async () => {
      try {
        // Implement auto-complete fetch later
        // const res = await fetch(`/api/autocomplete?q=${q}`);
        // const d = await res.json();
        // suggestions = d.suggestions || [];
        showSugg = true;
      } catch {
        suggestions = [];
      }
    }, 180);
    return () => clearTimeout(t);
  });

  function submit(e: Event) {
    e.preventDefault();
    if (!q.trim()) return;
    goto(`/search?q=${encodeURIComponent(q)}`);
    showSugg = false;
    openMobile = false;
  }

  async function logout() {
    // Implement logout later
    await fetch('/api/auth/logout', { method: 'POST' });
    window.location.href = '/';
  }
</script>

{#if !isReader}
<header class="sticky top-0 z-50 glass-strong">
  <div class="mx-auto flex h-16 w-full max-w-[1600px] items-center gap-4 px-4 md:px-6 2xl:px-8">
    <a href="/" class="flex items-center gap-2 font-display text-lg font-extrabold tracking-tight">
      <img src="/logo.ico" alt="NovelHive Logo" class="h-8 w-8" />
      <span class="gradient-text">NovelHive</span>
    </a>

    <nav class="hidden items-center gap-1 md:flex">
      {#each NAV as n}
        {@const active = path === n.to || (n.to !== "/" && path.startsWith(n.to))}
        <a
          href={n.to}
          class={cn(
            "rounded-lg px-3 py-1.5 text-sm font-medium transition",
            active ? "text-foreground bg-white/5" : "text-muted-foreground hover:text-foreground"
          )}
        >
          {n.label}
        </a>
      {/each}
    </nav>

    <!-- Search -->
    <form onsubmit={submit} class="relative ml-auto hidden flex-1 max-w-sm md:block">
      <div class="flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-3 py-1.5 focus-within:border-brand/60 focus-within:ring-2 focus-within:ring-brand/30 transition">
        <Search class="h-4 w-4 text-muted-foreground" />
        <input
          bind:value={q}
          onfocus={() => { if (suggestions.length) showSugg = true; }}
          placeholder="Search novels, authors…"
          class="w-full bg-transparent text-sm placeholder:text-muted-foreground focus:outline-none"
        />
        {#if q}
          <button type="button" onclick={() => q = ""} class="text-muted-foreground hover:text-foreground">
            <X class="h-3.5 w-3.5" />
          </button>
        {/if}
      </div>
      {#if showSugg && suggestions.length > 0}
        <div class="absolute left-0 right-0 top-full mt-2 overflow-hidden rounded-xl border border-white/10 bg-popover/95 shadow-elevated backdrop-blur">
          {#each suggestions.slice(0, 6) as s}
            <a
              href={`/novel/${s.slug}`}
              onclick={() => { showSugg = false; q = ""; }}
              class="flex items-center gap-2 px-3 py-2 text-sm hover:bg-white/5"
            >
              <BookOpen class="h-4 w-4 text-muted-foreground" />
              <span class="line-clamp-1">{s.title}</span>
            </a>
          {/each}
        </div>
      {/if}
    </form>

    <!-- User area -->
    <div class="relative ml-auto flex items-center gap-2 md:ml-0" bind:this={ref}>
      {#if user}
        <button
          onclick={() => openMenu = !openMenu}
          class="hidden items-center gap-2 rounded-full border border-white/10 bg-white/5 px-2 py-1 pr-3 text-sm hover:bg-white/10 md:flex"
        >
          <span class="grid h-7 w-7 place-items-center rounded-full gradient-brand text-xs font-bold text-white">
            {user.username[0]?.toUpperCase()}
          </span>
          <span class="font-semibold">{user.username}</span>
          <ChevronDown class="h-3.5 w-3.5 text-muted-foreground" />
        </button>
        {#if openMenu}
          <div class="absolute right-0 top-full mt-2 w-56 overflow-hidden rounded-xl border border-white/10 bg-popover/95 shadow-elevated backdrop-blur">
            <div class="border-b border-white/5 px-3 py-2.5">
              <div class="text-sm font-semibold">{user.username}</div>
              <div class="text-[11px] uppercase tracking-wider text-muted-foreground">{user.role}</div>
            </div>
            <a
              href="/library"
              onclick={() => openMenu = false}
              class="flex items-center gap-2 px-3 py-2 text-sm hover:bg-white/5"
            >
              <Library class="h-4 w-4" /> My Library
            </a>
            {#if user.role === "admin"}
              <a
                href="/admin"
                onclick={() => openMenu = false}
                class="flex items-center gap-2 px-3 py-2 text-sm hover:bg-white/5"
              >
                <Shield class="h-4 w-4" /> Admin Panel
              </a>
            {/if}
            <button
              onclick={() => { logout(); openMenu = false; }}
              class="flex w-full items-center gap-2 px-3 py-2 text-sm text-destructive hover:bg-destructive/10"
            >
              <LogOut class="h-4 w-4" /> Sign out
            </button>
          </div>
        {/if}
      {:else}
        <div class="hidden items-center gap-2 md:flex">
          <a href="/login" class="rounded-lg px-3 py-1.5 text-sm font-medium text-muted-foreground hover:text-foreground">
            Sign in
          </a>
          <a
            href="/register"
            class="rounded-lg gradient-brand px-3 py-1.5 text-sm font-bold text-white shadow-glow transition hover:opacity-90"
          >
            Sign up
          </a>
        </div>
      {/if}
      <button onclick={() => openMobile = !openMobile} class="md:hidden rounded-lg p-2 text-foreground hover:bg-white/5">
        {#if openMobile}
          <X class="h-5 w-5" />
        {:else}
          <Menu class="h-5 w-5" />
        {/if}
      </button>
    </div>
  </div>

  <!-- Mobile drawer -->
  {#if openMobile}
    <div class="border-t border-white/5 bg-background/95 px-4 py-4 backdrop-blur md:hidden">
      <form onsubmit={submit} class="mb-3 flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-3 py-2">
        <Search class="h-4 w-4 text-muted-foreground" />
        <input
          bind:value={q}
          placeholder="Search…"
          class="w-full bg-transparent text-sm focus:outline-none"
        />
      </form>
      {#each NAV as n}
        <a
          href={n.to}
          onclick={() => openMobile = false}
          class="block rounded-lg px-3 py-2 text-sm font-semibold text-muted-foreground hover:bg-white/5 hover:text-foreground"
        >
          {n.label}
        </a>
      {/each}
      <div class="my-2 h-px bg-white/5"></div>
      {#if user}
        <a href="/library" onclick={() => openMobile = false} class="block rounded-lg px-3 py-2 text-sm font-semibold hover:bg-white/5">
          My Library
        </a>
        {#if user.role === "admin"}
          <a href="/admin" onclick={() => openMobile = false} class="block rounded-lg px-3 py-2 text-sm font-semibold hover:bg-white/5">
            Admin Panel
          </a>
        {/if}
        <button
          onclick={() => { logout(); openMobile = false; }}
          class="block w-full rounded-lg px-3 py-2 text-left text-sm font-semibold text-destructive hover:bg-destructive/10"
        >
          Sign out
        </button>
      {:else}
        <div class="grid grid-cols-2 gap-2">
          <a href="/login" onclick={() => openMobile = false} class="rounded-lg border border-white/10 px-3 py-2 text-center text-sm font-semibold">
            Sign in
          </a>
          <a href="/register" onclick={() => openMobile = false} class="rounded-lg gradient-brand px-3 py-2 text-center text-sm font-bold text-white">
            Sign up
          </a>
        </div>
      {/if}
    </div>
  {/if}
</header>
{/if}
