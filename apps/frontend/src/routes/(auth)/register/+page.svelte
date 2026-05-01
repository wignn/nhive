<script lang="ts">
  import { Loader2, UserPlus } from "lucide-svelte";

  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import { goto, invalidateAll } from "$app/navigation";
  import { cn } from "$lib/utils";

  let busy = $state(false);
  let username = $state("");
  let email = $state("");
  let password = $state("");
  let confirm = $state("");

  let strength = $derived.by(() => {
    if (!password) return 0;
    if (password.length < 6) return 1;
    if (password.length < 10) return 2;
    return /[A-Z]/.test(password) && /\d/.test(password) ? 4 : 3;
  });

  function validate(e: Event) {
    if (password !== confirm) {
      e.preventDefault();
      toast.error("Passwords don't match");
      return false;
    }
    if (password.length < 6) {
      e.preventDefault();
      toast.error("Password must be ≥ 6 characters");
      return false;
    }
    if (username.length < 3) {
      e.preventDefault();
      toast.error("Username must be ≥ 3 characters");
      return false;
    }
    return true;
  }
</script>

<svelte:head>
  <title>Create account — NovelHive</title>
</svelte:head>

<main class="mx-auto grid min-h-[80vh] w-full max-w-md place-items-center px-4 py-10">
  <div class="w-full rounded-3xl border border-white/10 bg-card/70 p-7 shadow-elevated backdrop-blur">
    <div class="mb-6 text-center">
      <img src="/logo.ico" alt="NovelHive Logo" class="mx-auto h-12 w-12" />
      <h1 class="mt-4 font-display text-2xl font-extrabold">Join NovelHive</h1>
      <p class="mt-1 text-sm text-muted-foreground">Track stories, sync progress, drop comments.</p>
    </div>

    <form 
      method="POST" 
      class="space-y-3"
      onsubmit={validate}
      use:enhance={() => {
        busy = true;
        return async ({ result }) => {
          busy = false;
          if (result.type === 'redirect') {
            toast.success("Welcome to NovelHive");
            await invalidateAll();
            goto(result.location);
          } else if (result.type === 'failure') {
            toast.error(result.data?.error || "Registration failed");
          }
        };
      }}
    >
      <div>
        <label for="username" class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Username</label>
        <input
          id="username" required type="text" name="username" bind:value={username} placeholder="coolreader42" class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
      </div>
      <div>
        <label for="email" class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Email</label>
        <input
          id="email" required type="email" name="email" bind:value={email} placeholder="you@example.com" class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
      </div>
      <div>
        <label for="password" class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Password</label>
        <input
          id="password" required type="password" name="password" bind:value={password} placeholder="••••••••" class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
      </div>
      
      {#if password}
        <div>
          <div class="h-1.5 overflow-hidden rounded-full bg-white/5">
            <div
              class={cn(
                "h-full transition-all",
                strength <= 1 ? "bg-destructive" :
                strength === 2 ? "bg-warning" :
                strength === 3 ? "bg-success" : "gradient-brand",
              )}
              style="width: {Math.min(strength * 25, 100)}%"
            ></div>
          </div>
          <div class="mt-1 text-[11px] text-muted-foreground">
            {strength <= 1 ? "Weak" : strength === 2 ? "Fair" : strength === 3 ? "Good" : "Strong"} password
          </div>
        </div>
      {/if}

      <div>
        <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-muted-foreground">Confirm password</label>
        <input required type="password" bind:value={confirm} placeholder="••••••••" class="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2.5 text-sm focus:border-brand/60 focus:outline-none focus:ring-2 focus:ring-brand/30" />
      </div>

      <button
        disabled={busy}
        class="mt-2 inline-flex w-full items-center justify-center gap-2 rounded-xl gradient-brand py-2.5 text-sm font-bold text-white shadow-glow disabled:opacity-60"
      >
        {#if busy}
          <Loader2 class="h-4 w-4 animate-spin" />
          Creating…
        {:else}
          <UserPlus class="h-4 w-4" />
          Create account
        {/if}
      </button>
    </form>

    <p class="mt-6 text-center text-sm text-muted-foreground">
      Already have an account? <a href="/login" class="font-bold gradient-text">Sign in</a>
    </p>
  </div>
</main>
