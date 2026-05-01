<script lang="ts">
  import { buildCoverUrl, cn } from "$lib/utils";
  import { BookOpen, FileText, Globe, Loader2, Pencil, Plus, Shield, Tags, Trash2, Upload, User, Users, X } from "lucide-svelte";

  import { toast } from "svelte-sonner";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  
  let genres = $state<any[]>([]);
  let fGenreName = $state("");
  
  type Tab = "novels" | "chapters" | "genres" | "users";

  let user = $derived($page.data.user);
  let loadingAuth = false; // Layout handles auth loading already

  let tab = $state<Tab>("novels");
  let novels = $state<any[]>([]);
  let chapters = $state<any[]>([]);
  let users = $state<any[]>([]);
  let coverBase = $state<string | undefined>();
  let busy = $state(false);
  let showModal = $state(false);
  let editing = $state<any>(null);

  async function rpc(action: string, payload: any = {}) {
    const res = await fetch('/api/admin/rpc', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action, payload })
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Action failed');
    return data;
  }

  async function reload() {
    if (!user || user.role !== "admin") return;
    busy = true;
    try {
      if (tab === "novels") {
        const d = await rpc("listNovels");
        novels = d.novels || [];
        if (d.cover_base_url) coverBase = d.cover_base_url;
      } else if (tab === "chapters") {
        const d = await rpc("listChapters");
        chapters = (d.chapters || []).map((c: any) => ({ ...(c.Chapter || c), novel_title: c.NovelTitle || c.novel_title }));
      } else if (tab === "genres") {
        const d = await rpc("listGenres");
        genres = d.genres || [];
      } else if (tab === "users") {
        const d = await rpc("listUsers");
        users = d.users || [];
      }
    } catch (e: any) {
      toast.error(e?.message || "Failed to load");
    } finally {
      busy = false;
    }
  }

  $effect(() => {
    if (user?.role === "admin") reload();
  });

  $effect(() => {
    if (user?.role === "admin" && novels.length === 0 && tab !== "novels") {
      rpc("listNovels").then((d: any) => {
        novels = d.novels || [];
        if (d.cover_base_url) coverBase = d.cover_base_url;
      }).catch(() => {});
    }
    if (user?.role === "admin" && genres.length === 0) {
      rpc("listGenres").then((d: any) => {
        genres = d.genres || [];
      }).catch(() => {});
    }
  });

  if (user && user.role !== "admin") {
    goto("/");
  }

  const stats = $derived([
    { label: "Novels", value: novels.length, icon: BookOpen },
    { label: "Chapters", value: chapters.length, icon: FileText },
    { label: "Genres", value: genres.length, icon: Tags },
    { label: "Users", value: users.length, icon: Users },
  ]);

  // Form states
  let fTitle = $state("");
  let fAuthor = $state("");
  let fSynopsis = $state("");
  let fStatus = $state("ongoing");
  let fCoverUrl = $state("");
  let fCoverPreview = $state("");
  let fGenres = $state<string[]>([]);
  let uploading = $state(false);

  // Chapter form states
  let fNovelId = $state("");
  let fNumber = $state(1);
  let fTitleCh = $state("");
  let fContent = $state("");

  function openNovelForm(n?: any) {
    editing = n;
    const initialPreview = n ? buildCoverUrl(n.cover_url, coverBase) || "" : "";
    fTitle = n?.title || "";
    fAuthor = n?.author || "";
    fSynopsis = n?.synopsis || "";
    fStatus = n?.status || "ongoing";
    fCoverUrl = n?.cover_url || "";
    fCoverPreview = initialPreview;
    fGenres = (n?.genres || []).map((g: any) => (typeof g === "string" ? g : g.name));
    showModal = true;
  }

  function openChapterForm(c?: any) {
    editing = c;
    fNovelId = c?.novel_id || novels[0]?.id || "";
    fNumber = c?.number || 1;
    fTitleCh = c?.title || "";
    fContent = c?.content || "";
    showModal = true;
  }

  async function uploadCover(e: Event) {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;
    uploading = true;
    try {
      const fd = new FormData();
      fd.append("image", file);
      const res = await fetch("/api/admin/upload", { method: "POST", body: fd });
      if (!res.ok) throw new Error("Upload failed");
      const r = await res.json();
      fCoverUrl = r.path;
      fCoverPreview = `${r.base_url}/${r.path}`;
      toast.success("Cover uploaded");
    } catch (e: any) {
      toast.error(e?.message || "Upload failed");
    } finally {
      uploading = false;
    }
  }

  async function saveNovel() {
    const data = {
      title: fTitle, author: fAuthor, synopsis: fSynopsis, status: fStatus,
      cover_url: fCoverUrl, genres: fGenres,
    };
    try {
      if (editing) await rpc("updateNovel", { id: editing.id, data });
      else await rpc("createNovel", data);
      toast.success(editing ? "Updated" : "Created");
      showModal = false;
      reload();
    } catch (e: any) { toast.error(e?.message); }
  }

  async function saveChapter() {
    const data = {
      novel_id: fNovelId, number: fNumber, title: fTitleCh, content: fContent,
    };
    try {
      if (editing) await rpc("updateChapter", { id: editing.id, ...data });
      else await rpc("createChapter", data);
      toast.success(editing ? "Updated" : "Created");
      showModal = false;
      reload();
    } catch (e: any) { toast.error(e?.message); }
  }
</script>

<svelte:head>
  <title>Admin Panel — NovelHive</title>
</svelte:head>

{#if !user || user.role !== "admin"}
  <main class="p-10 text-center text-sm text-muted-foreground">Checking permissions…</main>
{:else}
  <div class="mx-auto grid w-full max-w-[1600px] grid-cols-1 gap-6 px-4 py-6 md:grid-cols-[240px_1fr] md:px-6 2xl:px-8 2xl:gap-8">
    <!-- Sidebar -->
    <aside class="rounded-2xl border border-white/5 bg-card/60 p-4 md:sticky md:top-20 md:h-fit">
      <a href="/" class="mb-4 flex items-center gap-2">
        <img src="/logo.ico" alt="NovelHive Logo" class="h-8 w-8" />
        <div>
          <div class="font-display text-sm font-extrabold gradient-text">NovelHive</div>
          <div class="text-[10px] uppercase tracking-widest text-muted-foreground">Admin</div>
        </div>
      </a>
      <div class="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Manage</div>
      <nav class="mt-2 space-y-1">
        {#each ["novels","chapters","genres","users"] as t}
          <button
            onclick={() => tab = t as Tab}
            class={cn(
              "flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm font-semibold capitalize transition",
              tab === t ? "bg-brand/15 text-brand" : "text-muted-foreground hover:bg-white/5 hover:text-foreground",
            )}
          >
            {#if t === "novels"}<BookOpen class="h-4 w-4" />{/if}
            {#if t === "chapters"}<FileText class="h-4 w-4" />{/if}
            {#if t === "genres"}<Tags class="h-4 w-4" />{/if}
            {#if t === "users"}<Users class="h-4 w-4" />{/if}
            {t}
          </button>
        {/each}
      </nav>
      <div class="mt-4 text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Shortcuts</div>
      <a href="/" class="mt-2 flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-semibold text-muted-foreground hover:bg-white/5 hover:text-foreground">
        <Globe class="h-4 w-4" /> View site
      </a>
    </aside>

    <!-- Main -->
    <section>
      <!-- Stats -->
      <div class="mb-6 grid grid-cols-2 gap-3 md:grid-cols-4">
        {#each stats as s}
          <div class="rounded-2xl border border-white/5 bg-card/60 p-4">
            <div class="flex items-center justify-between">
              <div>
                <div class="font-display text-2xl font-extrabold">{s.value}</div>
                <div class="text-xs text-muted-foreground">{s.label}</div>
              </div>
              <div class="grid h-9 w-9 place-items-center rounded-lg gradient-brand/30 bg-brand/10 text-brand">
                <svelte:component this={s.icon} class="h-4 w-4" />
              </div>
            </div>
          </div>
        {/each}
      </div>

      <!-- Toolbar -->
      <div class="mb-4 flex items-center justify-between">
        <h2 class="font-display text-xl font-extrabold capitalize">{tab}</h2>
        {#if tab === "novels" || tab === "chapters"}
          <button
            onclick={() => tab === "novels" ? openNovelForm() : openChapterForm()}
            class="inline-flex items-center gap-2 rounded-xl gradient-brand px-4 py-2 text-sm font-bold text-white shadow-glow"
          >
            <Plus class="h-4 w-4" /> New {tab === "novels" ? "novel" : "chapter"}
          </button>
        {/if}
      </div>

      {#if busy}
        <div class="rounded-2xl border border-white/5 bg-card/60 p-12 text-center text-sm text-muted-foreground">
          <Loader2 class="mx-auto mb-2 h-5 w-5 animate-spin" /> Loading…
        </div>
      {:else}
        {#if tab === "novels"}
          <div class="overflow-hidden rounded-2xl border border-white/5 bg-card/60">
            <table class="w-full text-left">
              <thead>
                <tr class="bg-white/[0.02]">
                  {#each ["Cover","Title","Author","Status","Ch.",""] as h}
                    <th class="px-3 py-2 text-[10px] font-bold uppercase tracking-widest text-muted-foreground">{h}</th>
                  {/each}
                </tr>
              </thead>
              <tbody>
                {#each novels as n}
                  <tr class="border-t border-white/5">
                    <td class="p-3">
                      <div class="h-14 w-10 overflow-hidden rounded-md bg-white/5">
                        {#if n.cover_url}
                          <img src={buildCoverUrl(n.cover_url, coverBase) || ""} class="h-full w-full object-cover" alt="Cover" />
                        {/if}
                      </div>
                    </td>
                    <td class="p-3 font-semibold">{n.title}</td>
                    <td class="p-3 text-sm text-muted-foreground">{n.author}</td>
                    <td class="p-3">
                      <span class={cn(
                        "rounded-md px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider",
                        n.status === "completed" ? "bg-success/20 text-success" :
                        n.status === "hiatus" ? "bg-warning/20 text-warning" :
                        "bg-brand/20 text-brand",
                      )}>{n.status}</span>
                    </td>
                    <td class="p-3 text-sm">{n.total_chapters ?? 0}</td>
                    <td class="p-3 text-right">
                      <div class="inline-flex items-center gap-2">
                        <button onclick={() => openNovelForm(n)} class="rounded-lg border border-white/10 bg-white/5 p-1.5 hover:bg-white/10" title="Edit"><Pencil class="h-3.5 w-3.5" /></button>
                        <button onclick={async () => {
                          if (!confirm("Delete this novel and its chapters?")) return;
                          try { await rpc("deleteNovel", { id: n.id }); toast.success("Deleted"); reload(); }
                          catch (e: any) { toast.error(e?.message); }
                        }} class="rounded-lg border border-destructive/30 bg-destructive/10 p-1.5 text-destructive hover:bg-destructive/20" title="Delete"><Trash2 class="h-3.5 w-3.5" /></button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}

        {#if tab === "chapters"}
          <div class="overflow-hidden rounded-2xl border border-white/5 bg-card/60">
            <table class="w-full text-left">
              <thead>
                <tr class="bg-white/[0.02]">
                  {#each ["Novel","#","Title","Words",""] as h}
                    <th class="px-3 py-2 text-[10px] font-bold uppercase tracking-widest text-muted-foreground">{h}</th>
                  {/each}
                </tr>
              </thead>
              <tbody>
                {#each chapters as c}
                  <tr class="border-t border-white/5">
                    <td class="p-3 text-sm text-muted-foreground">{c.novel_title}</td>
                    <td class="p-3 font-mono text-sm">{c.number}</td>
                    <td class="p-3 font-semibold">{c.title}</td>
                    <td class="p-3 text-sm">{c.word_count}</td>
                    <td class="p-3 text-right">
                      <div class="inline-flex items-center gap-2">
                        <button onclick={() => openChapterForm(c)} class="rounded-lg border border-white/10 bg-white/5 p-1.5 hover:bg-white/10" title="Edit"><Pencil class="h-3.5 w-3.5" /></button>
                        <button onclick={async () => {
                          if (!confirm("Delete this chapter?")) return;
                          try { await rpc("deleteChapter", { id: c.id }); toast.success("Deleted"); reload(); }
                          catch (e: any) { toast.error(e?.message); }
                        }} class="rounded-lg border border-destructive/30 bg-destructive/10 p-1.5 text-destructive hover:bg-destructive/20" title="Delete"><Trash2 class="h-3.5 w-3.5" /></button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}

        {#if tab === "genres"}
          <div class="mb-4">
            <form class="flex gap-2" onsubmit={(e) => {
              e.preventDefault();
              if (!fGenreName) return;
              rpc("createGenre", { name: fGenreName }).then(() => {
                fGenreName = "";
                toast.success("Genre created");
                reload();
              }).catch((err: any) => toast.error(err.message));
            }}>
              <input bind:value={fGenreName} placeholder="New genre name..." class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
              <button type="submit" class="rounded-xl gradient-brand px-4 py-2 text-sm font-bold text-white shadow-glow whitespace-nowrap">Create Genre</button>
            </form>
          </div>
          <div class="flex flex-wrap gap-2 rounded-2xl border border-white/5 bg-card/60 p-4">
            {#each genres as g}
              <span class="inline-flex items-center gap-1.5 rounded-md border border-white/10 bg-white/5 pl-3 pr-1 py-1 text-xs font-semibold">
                {g.name}
                <button onclick={() => {
                  if (!confirm(`Delete genre: ${g.name}?`)) return;
                  rpc("deleteGenre", { id: g.id }).then(() => {
                    toast.success("Deleted");
                    reload();
                  }).catch((err: any) => toast.error(err.message));
                }} class="rounded px-1.5 py-1 text-muted-foreground hover:bg-destructive/20 hover:text-destructive"><X class="h-3.5 w-3.5"/></button>
              </span>
            {/each}
          </div>
        {/if}

        {#if tab === "users"}
          <div class="overflow-hidden rounded-2xl border border-white/5 bg-card/60">
            <table class="w-full text-left">
              <thead>
                <tr class="bg-white/[0.02]">
                  {#each ["Username","Email","Role",""] as h}
                    <th class="px-3 py-2 text-[10px] font-bold uppercase tracking-widest text-muted-foreground">{h}</th>
                  {/each}
                </tr>
              </thead>
              <tbody>
                {#each users as u}
                  <tr class="border-t border-white/5">
                    <td class="p-3 font-semibold">{u.username}</td>
                    <td class="p-3 text-sm text-muted-foreground">{u.email}</td>
                    <td class="p-3">
                      <span class={cn(
                        "rounded-md px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider",
                        u.role === "admin" ? "bg-brand/20 text-brand" : "bg-white/10 text-muted-foreground",
                      )}>{u.role}</span>
                    </td>
                    <td class="p-3 text-right">
                      <button
                        onclick={async () => {
                          const newRole = u.role === "admin" ? "reader" : "admin";
                          try { await rpc("setUserRole", { id: u.id, role: newRole }); toast.success(`Role: ${newRole}`); reload(); }
                          catch (e: any) { toast.error(e?.message); }
                        }}
                        class="inline-flex items-center gap-1.5 rounded-lg border border-white/10 bg-white/5 px-3 py-1.5 text-xs font-semibold hover:bg-white/10"
                      >
                        {#if u.role === "admin"}<User class="h-3.5 w-3.5" /> Make reader{:else}<Shield class="h-3.5 w-3.5" /> Make admin{/if}
                      </button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      {/if}
    </section>

    <!-- MODAL -->
    {#if showModal}
      <div class="fixed inset-0 z-[60] grid place-items-center bg-black/70 backdrop-blur-sm p-4" onclick={() => showModal = false} role="presentation">
        <div onclick={(e) => e.stopPropagation()} class="w-full max-w-2xl max-h-[90vh] overflow-auto rounded-3xl border border-white/10 bg-card p-6 shadow-elevated" role="presentation">
          <div class="mb-4 flex items-center justify-end">
            <button onclick={() => showModal = false} class="rounded-lg p-1.5 hover:bg-white/5"><X class="h-4 w-4" /></button>
          </div>
          
          {#if tab === "novels"}
            <h2 class="mb-4 font-display text-xl font-extrabold">{editing ? "Edit novel" : "Create novel"}</h2>
            <div class="grid gap-3 md:grid-cols-2">
              <div>
                <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Title</label>
                <input bind:value={fTitle} class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
              </div>
              <div>
                <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Author</label>
                <input bind:value={fAuthor} class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
              </div>
            </div>
            <div class="mt-3">
              <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Cover image</label>
              <div class="flex items-center gap-3">
                <div class="h-24 w-16 shrink-0 overflow-hidden rounded-lg border border-white/10 bg-white/5">
                  {#if fCoverPreview}
                    <img src={fCoverPreview} class="h-full w-full object-cover" alt="Preview" />
                  {/if}
                </div>
                <label class="flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-xl border border-dashed border-white/15 bg-white/[0.02] px-3 py-4 text-xs font-semibold text-muted-foreground hover:bg-white/5">
                  <Upload class="h-4 w-4" /> {uploading ? "Uploading…" : "Upload cover"}
                  <input type="file" accept="image/*" disabled={uploading} onchange={uploadCover} class="hidden" />
                </label>
              </div>
            </div>
            <div class="mt-3">
              <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Synopsis</label>
              <textarea
                bind:value={fSynopsis}
                class="min-h-24 w-full resize-y rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30"
              ></textarea>
            </div>
            <div class="mt-3">
              <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Status</label>
              <select
                bind:value={fStatus}
                class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:outline-none"
              >
                <option class="bg-background" value="ongoing">Ongoing</option>
                <option class="bg-background" value="completed">Completed</option>
                <option class="bg-background" value="hiatus">Hiatus</option>
              </select>
            </div>
            <div class="mt-3">
              <label class="mb-2 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Genres</label>
              <div class="flex flex-wrap gap-1.5">
                {#each genres as gObj}
                  {@const g = gObj.name}
                  {@const on = fGenres.includes(g)}
                  <button
                    type="button"
                    onclick={() => fGenres = on ? fGenres.filter(x => x !== g) : [...fGenres, g]}
                    class={cn(
                      "rounded-md border px-2 py-1 text-[11px] font-semibold transition",
                      on ? "border-brand bg-brand/15 text-brand" : "border-white/10 bg-white/5 text-muted-foreground hover:text-foreground",
                    )}
                  >{g}</button>
                {/each}
              </div>
            </div>
            <div class="mt-5 flex justify-end gap-2">
              <button onclick={() => showModal = false} class="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold">Cancel</button>
              <button
                onclick={saveNovel}
                class="rounded-xl gradient-brand px-4 py-2 text-sm font-bold text-white shadow-glow"
              >{editing ? "Save changes" : "Create"}</button>
            </div>
          {:else}
            <h2 class="mb-4 font-display text-xl font-extrabold">{editing ? "Edit chapter" : "Create chapter"}</h2>
            {#if !editing}
              <div class="mb-3">
                <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Novel</label>
                <select
                  bind:value={fNovelId}
                  class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:outline-none"
                >
                  {#each novels as n}
                    <option value={n.id} class="bg-background">{n.title}</option>
                  {/each}
                </select>
              </div>
            {/if}
            <div class="grid grid-cols-[100px_1fr] gap-3">
              <div>
                <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Number</label>
                <input type="number" bind:value={fNumber} class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
              </div>
              <div>
                <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Title</label>
                <input bind:value={fTitleCh} class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
              </div>
            </div>
            <div class="mt-3">
              <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Content</label>
              <textarea
                bind:value={fContent}
                class="min-h-64 w-full resize-y rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30 font-mono"
              ></textarea>
            </div>
            <div class="mt-5 flex justify-end gap-2">
              <button onclick={() => showModal = false} class="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold">Cancel</button>
              <button
                onclick={saveChapter}
                class="rounded-xl gradient-brand px-4 py-2 text-sm font-bold text-white shadow-glow"
              >{editing ? "Save changes" : "Publish"}</button>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
{/if}
