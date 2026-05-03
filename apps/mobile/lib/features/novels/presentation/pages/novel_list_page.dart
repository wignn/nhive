import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shimmer/shimmer.dart';
import 'package:nhive/app/theme/app_theme.dart';
import 'package:nhive/features/novels/presentation/widgets/novel_card.dart';
import 'package:nhive/features/novels/presentation/bloc/novel_provider.dart';
import 'package:nhive/features/auth/presentation/bloc/auth_provider.dart';
import 'package:nhive/features/library/presentation/pages/library_page.dart';
import 'package:nhive/features/search/presentation/pages/search_page.dart';
import 'package:nhive/features/auth/presentation/pages/profile_page.dart';

class NovelListPage extends StatefulWidget {
  const NovelListPage({super.key});

  @override
  State<NovelListPage> createState() => _NovelListPageState();
}

class _NovelListPageState extends State<NovelListPage> {
  int _currentIndex = 0;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<NovelProvider>().loadNovels(refresh: true);
      context.read<AuthProvider>().checkAuth();
    });
  }

  @override
  Widget build(BuildContext context) {
    final pages = [
      const _HomeView(),
      const LibraryPage(),
      const SearchPage(),
      const ProfilePage(),
    ];

    return Scaffold(
      body: IndexedStack(index: _currentIndex, children: pages),
      bottomNavigationBar: Container(
        decoration: BoxDecoration(
          color: AppTheme.background,
          border: const Border(top: BorderSide(color: AppTheme.border)),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.3),
              blurRadius: 20,
              offset: const Offset(0, -5),
            ),
          ],
        ),
        child: BottomNavigationBar(
          currentIndex: _currentIndex,
          onTap: (index) => setState(() => _currentIndex = index),
          type: BottomNavigationBarType.fixed,
          backgroundColor: Colors.transparent,
          elevation: 0,
          selectedItemColor: AppTheme.primary,
          unselectedItemColor: AppTheme.muted,
          selectedFontSize: 12,
          unselectedFontSize: 12,
          items: const [
            BottomNavigationBarItem(
              icon: Icon(Icons.home_outlined),
              activeIcon: Icon(Icons.home),
              label: 'Home',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.bookmark_outline),
              activeIcon: Icon(Icons.bookmark),
              label: 'Library',
            ),
            BottomNavigationBarItem(icon: Icon(Icons.search), label: 'Search'),
            BottomNavigationBarItem(
              icon: Icon(Icons.person_outline),
              activeIcon: Icon(Icons.person),
              label: 'Profile',
            ),
          ],
        ),
      ),
    );
  }
}

class _HomeView extends StatelessWidget {
  const _HomeView();

  @override
  Widget build(BuildContext context) {
    return Consumer<NovelProvider>(
      builder: (context, provider, _) {
        return RefreshIndicator(
          color: AppTheme.primary,
          onRefresh: () => provider.loadNovels(refresh: true),
          child: CustomScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            slivers: [
              SliverAppBar(
                pinned: true,
                backgroundColor: AppTheme.background,
                elevation: 0,
                title: ShaderMask(
                  shaderCallback: (bounds) =>
                      AppTheme.brandGradient.createShader(bounds),
                  child: const Text(
                    'NovelHive',
                    style: TextStyle(fontWeight: FontWeight.w800, fontSize: 24),
                  ),
                ),
                actions: [
                  IconButton(
                    onPressed: () {},
                    icon: const Icon(Icons.notifications_outlined),
                  ),
                ],
              ),

              // ─── Trending Section ───
              const SliverToBoxAdapter(
                child: Padding(
                  padding: EdgeInsets.fromLTRB(16, 16, 16, 12),
                  child: Row(
                    children: [
                      Icon(
                        Icons.local_fire_department,
                        color: AppTheme.primary,
                        size: 22,
                      ),
                      SizedBox(width: 8),
                      Text(
                        'Trending Now',
                        style: TextStyle(
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ),

              SliverToBoxAdapter(
                child: SizedBox(
                  height: 260,
                  child: provider.isLoadingNovels && provider.novels.isEmpty
                      ? _buildTrendingShimmer()
                      : ListView.separated(
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          scrollDirection: Axis.horizontal,
                          itemCount: provider.novels.take(5).length,
                          separatorBuilder: (_, __) =>
                              const SizedBox(width: 12),
                          itemBuilder: (context, index) {
                            final novel = provider.novels[index];
                            return SizedBox(
                              width: 130,
                              child: NovelCard(
                                novel: novel,
                                onTap: () => Navigator.pushNamed(
                                  context,
                                  '/detail',
                                  arguments: novel.slug,
                                ),
                              ),
                            );
                          },
                        ),
                ),
              ),

              // ─── Recently Updated Section ───
              const SliverToBoxAdapter(
                child: Padding(
                  padding: EdgeInsets.fromLTRB(16, 24, 16, 12),
                  child: Row(
                    children: [
                      Icon(Icons.update, color: AppTheme.primary, size: 20),
                      SizedBox(width: 8),
                      Text(
                        'Recently Updated',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ),

              if (provider.isLoadingNovels && provider.novels.isEmpty)
                SliverPadding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  sliver: SliverGrid(
                    gridDelegate:
                        const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 3,
                          mainAxisSpacing: 16,
                          crossAxisSpacing: 12,
                          childAspectRatio: 0.52,
                        ),
                    delegate: SliverChildBuilderDelegate(
                      (context, index) => _buildCardShimmer(),
                      childCount: 6,
                    ),
                  ),
                )
              else if (provider.novelsError != null && provider.novels.isEmpty)
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.all(32),
                    child: Column(
                      children: [
                        const Icon(
                          Icons.cloud_off,
                          size: 48,
                          color: AppTheme.muted,
                        ),
                        const SizedBox(height: 12),
                        Text(
                          'Failed to load novels',
                          style: TextStyle(color: AppTheme.muted, fontSize: 16),
                        ),
                        const SizedBox(height: 12),
                        ElevatedButton.icon(
                          onPressed: () => provider.loadNovels(refresh: true),
                          icon: const Icon(Icons.refresh, size: 18),
                          label: const Text('Retry'),
                        ),
                      ],
                    ),
                  ),
                )
              else
                SliverPadding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  sliver: SliverGrid(
                    gridDelegate:
                        const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 3,
                          mainAxisSpacing: 16,
                          crossAxisSpacing: 12,
                          childAspectRatio: 0.52,
                        ),
                    delegate: SliverChildBuilderDelegate((context, index) {
                      final novel = provider.novels[index];
                      return NovelCard(
                        novel: novel,
                        onTap: () => Navigator.pushNamed(
                          context,
                          '/detail',
                          arguments: novel.slug,
                        ),
                      );
                    }, childCount: provider.novels.length),
                  ),
                ),

              const SliverToBoxAdapter(child: SizedBox(height: 32)),
            ],
          ),
        );
      },
    );
  }

  Widget _buildTrendingShimmer() {
    return ListView.separated(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      scrollDirection: Axis.horizontal,
      itemCount: 4,
      separatorBuilder: (_, __) => const SizedBox(width: 12),
      itemBuilder: (_, __) => SizedBox(
        width: 130,
        child: Shimmer.fromColors(
          baseColor: AppTheme.secondary,
          highlightColor: AppTheme.surface,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Container(
                  decoration: BoxDecoration(
                    color: AppTheme.secondary,
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
              ),
              const SizedBox(height: 8),
              Container(
                height: 14,
                width: 100,
                decoration: BoxDecoration(
                  color: AppTheme.secondary,
                  borderRadius: BorderRadius.circular(4),
                ),
              ),
              const SizedBox(height: 4),
              Container(
                height: 10,
                width: 60,
                decoration: BoxDecoration(
                  color: AppTheme.secondary,
                  borderRadius: BorderRadius.circular(4),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildCardShimmer() {
    return Shimmer.fromColors(
      baseColor: AppTheme.secondary,
      highlightColor: AppTheme.surface,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                color: AppTheme.secondary,
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
          const SizedBox(height: 8),
          Container(
            height: 14,
            width: double.infinity,
            decoration: BoxDecoration(
              color: AppTheme.secondary,
              borderRadius: BorderRadius.circular(4),
            ),
          ),
          const SizedBox(height: 4),
          Container(
            height: 10,
            width: 60,
            decoration: BoxDecoration(
              color: AppTheme.secondary,
              borderRadius: BorderRadius.circular(4),
            ),
          ),
        ],
      ),
    );
  }
}
