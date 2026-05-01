<script lang="ts">
  let { count = 12 } = $props<{ count?: number }>();
</script>

<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:gap-5 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-8">
  {#each Array(count) as _}
    <div class="cover-frame rounded-xl bg-white/5 ring-1 ring-white/5">
      <div class="absolute inset-0 animate-pulse bg-gradient-to-br from-white/5 to-transparent"></div>
    </div>
  {/each}
</div>
