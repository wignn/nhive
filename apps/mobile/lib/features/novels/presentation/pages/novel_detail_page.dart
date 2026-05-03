import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:nhive/app/theme/app_theme.dart';
import 'package:nhive/features/novels/presentation/bloc/novel_provider.dart';
import 'package:nhive/features/library/presentation/bloc/library_provider.dart';
import 'package:nhive/features/auth/presentation/bloc/auth_provider.dart';

class NovelDetailPage extends StatefulWidget {
  final String slug;
  const NovelDetailPage({super.key, required this.slug});

  @override
  State<NovelDetailPage> createState() => _NovelDetailPageState();
}

class _NovelDetailPageState extends State<NovelDetailPage> {
  bool _synopsisExpanded = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<NovelProvider>().loadNovelDetail(widget.slug);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<NovelProvider>(
      builder: (context, provider, _) {
        if (provider.isLoadingDetail) {
          return Scaffold(
            body: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const CircularProgressIndicator(color: AppTheme.primary),
                  const SizedBox(height: 16),
                  Text('Loading...', style: TextStyle(color: AppTheme.muted)),
                ],
              ),
            ),
          );
        }

        if (provider.detailError != null || provider.selectedNovel == null) {
          return Scaffold(
            appBar: AppBar(),
            body: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(
                    Icons.error_outline,
                    size: 48,
                    color: AppTheme.muted,
                  ),
                  const SizedBox(height: 12),
                  const Text(
                    'Failed to load novel',
                    style: TextStyle(color: AppTheme.muted),
                  ),
                  const SizedBox(height: 12),
                  ElevatedButton(
                    onPressed: () => provider.loadNovelDetail(widget.slug),
                    child: const Text('Retry'),
                  ),
                ],
              ),
            ),
          );
        }

        final novel = provider.selectedNovel!;
        final chapters = provider.chapters;

        return Scaffold(
          body: CustomScrollView(
            slivers: [
              // ─── Hero Cover ───
              SliverAppBar(
                expandedHeight: 380,
                pinned: true,
                stretch: true,
                flexibleSpace: FlexibleSpaceBar(
                  background: Stack(
                    fit: StackFit.expand,
                    children: [
                      if (novel.coverUrl != null && novel.coverUrl!.isNotEmpty)
                        CachedNetworkImage(
                          imageUrl: novel.coverUrl!,
                          fit: BoxFit.cover,
                          placeholder: (_, __) =>
                              Container(color: AppTheme.secondary),
                          errorWidget: (_, __, ___) => Container(
                            color: AppTheme.secondary,
                            child: const Icon(
                              Icons.auto_stories,
                              size: 64,
                              color: AppTheme.muted,
                            ),
                          ),
                        )
                      else
                        Container(
                          color: AppTheme.secondary,
                          child: const Icon(
                            Icons.auto_stories,
                            size: 64,
                            color: AppTheme.muted,
                          ),
                        ),
                      // Gradient overlay
                      Container(
                        decoration: const BoxDecoration(
                          gradient: LinearGradient(
                            begin: Alignment.topCenter,
                            end: Alignment.bottomCenter,
                            colors: [
                              Colors.transparent,
                              Color(0x99000000),
                              AppTheme.background,
                            ],
                            stops: [0.3, 0.7, 1.0],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),

              // ─── Novel Info ───
              SliverToBoxAdapter(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Status & Views row
                      Row(
                        children: [
                          _buildStatusBadge(novel.status ?? 'Unknown'),
                          const SizedBox(width: 12),
                          if (novel.views > 0) ...[
                            const Icon(
                              Icons.remove_red_eye_outlined,
                              size: 14,
                              color: AppTheme.muted,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              _formatViews(novel.views),
                              style: const TextStyle(
                                color: AppTheme.muted,
                                fontSize: 12,
                              ),
                            ),
                          ],
                          if (chapters.isNotEmpty) ...[
                            const SizedBox(width: 12),
                            const Icon(
                              Icons.menu_book_outlined,
                              size: 14,
                              color: AppTheme.muted,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              '${chapters.length} chapters',
                              style: const TextStyle(
                                color: AppTheme.muted,
                                fontSize: 12,
                              ),
                            ),
                          ],
                        ],
                      ),
                      const SizedBox(height: 16),

                      // Title
                      Text(
                        novel.title,
                        style: const TextStyle(
                          fontSize: 26,
                          fontWeight: FontWeight.w800,
                          height: 1.2,
                        ),
                      ),
                      const SizedBox(height: 8),

                      // Author
                      if (novel.author != null && novel.author!.isNotEmpty)
                        Row(
                          children: [
                            const Icon(
                              Icons.person_outline,
                              size: 16,
                              color: AppTheme.primary,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              novel.author!,
                              style: const TextStyle(
                                color: AppTheme.primary,
                                fontSize: 15,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ],
                        ),
                      const SizedBox(height: 16),

                      // Genres
                      if (novel.genres.isNotEmpty)
                        Wrap(
                          spacing: 8,
                          runSpacing: 8,
                          children: novel.genres
                              .map(
                                (g) => Container(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 12,
                                    vertical: 6,
                                  ),
                                  decoration: BoxDecoration(
                                    color: AppTheme.secondary,
                                    borderRadius: BorderRadius.circular(20),
                                    border: Border.all(color: AppTheme.border),
                                  ),
                                  child: Text(
                                    g.name,
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: AppTheme.foreground,
                                    ),
                                  ),
                                ),
                              )
                              .toList(),
                        ),
                      const SizedBox(height: 24),

                      // Synopsis
                      const Text(
                        'Synopsis',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 8),
                      GestureDetector(
                        onTap: () => setState(
                          () => _synopsisExpanded = !_synopsisExpanded,
                        ),
                        child: AnimatedCrossFade(
                          firstChild: Text(
                            novel.description ?? 'No description available.',
                            maxLines: 4,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                              color: AppTheme.muted,
                              fontSize: 14,
                              height: 1.6,
                            ),
                          ),
                          secondChild: Text(
                            novel.description ?? 'No description available.',
                            style: const TextStyle(
                              color: AppTheme.muted,
                              fontSize: 14,
                              height: 1.6,
                            ),
                          ),
                          crossFadeState: _synopsisExpanded
                              ? CrossFadeState.showSecond
                              : CrossFadeState.showFirst,
                          duration: const Duration(milliseconds: 300),
                        ),
                      ),
                      if ((novel.description?.length ?? 0) > 150)
                        TextButton(
                          onPressed: () => setState(
                            () => _synopsisExpanded = !_synopsisExpanded,
                          ),
                          child: Text(
                            _synopsisExpanded ? 'Show less' : 'Read more',
                            style: const TextStyle(
                              color: AppTheme.primary,
                              fontSize: 13,
                            ),
                          ),
                        ),

                      const SizedBox(height: 24),

                      // Chapter List
                      if (chapters.isNotEmpty) ...[
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              'Chapters (${chapters.length})',
                              style: const TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 12),
                      ],
                    ],
                  ),
                ),
              ),

              // ─── Chapter List ───
              if (chapters.isNotEmpty)
                SliverPadding(
                  padding: const EdgeInsets.symmetric(horizontal: 20),
                  sliver: SliverList(
                    delegate: SliverChildBuilderDelegate((context, index) {
                      final chapter = chapters[index];
                      return Container(
                        margin: const EdgeInsets.only(bottom: 8),
                        decoration: BoxDecoration(
                          color: AppTheme.surface,
                          borderRadius: BorderRadius.circular(12),
                          border: Border.all(color: AppTheme.border),
                        ),
                        child: ListTile(
                          onTap: () {
                            Navigator.pushNamed(
                              context,
                              '/reader',
                              arguments: {
                                'slug': novel.slug,
                                'chapterNumber': chapter.number,
                                'novelTitle': novel.title,
                              },
                            );
                          },
                          contentPadding: const EdgeInsets.symmetric(
                            horizontal: 16,
                            vertical: 4,
                          ),
                          leading: Container(
                            width: 40,
                            height: 40,
                            decoration: BoxDecoration(
                              color: AppTheme.primary.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(10),
                            ),
                            child: Center(
                              child: Text(
                                '${chapter.number}',
                                style: const TextStyle(
                                  color: AppTheme.primary,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 14,
                                ),
                              ),
                            ),
                          ),
                          title: Text(
                            chapter.title,
                            style: const TextStyle(
                              fontWeight: FontWeight.w600,
                              fontSize: 14,
                            ),
                          ),
                          subtitle: chapter.wordCount > 0
                              ? Text(
                                  '${chapter.wordCount} words',
                                  style: const TextStyle(
                                    color: AppTheme.muted,
                                    fontSize: 12,
                                  ),
                                )
                              : null,
                          trailing: const Icon(
                            Icons.chevron_right,
                            color: AppTheme.muted,
                          ),
                        ),
                      );
                    }, childCount: chapters.length),
                  ),
                )
              else if (!provider.isLoadingDetail)
                const SliverToBoxAdapter(
                  child: Padding(
                    padding: EdgeInsets.all(32),
                    child: Center(
                      child: Column(
                        children: [
                          Icon(
                            Icons.hourglass_empty,
                            size: 40,
                            color: AppTheme.muted,
                          ),
                          SizedBox(height: 8),
                          Text(
                            'No chapters yet',
                            style: TextStyle(color: AppTheme.muted),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),

              const SliverToBoxAdapter(child: SizedBox(height: 100)),
            ],
          ),

          // ─── Bottom Action Bar ───
          bottomSheet: Container(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 16),
            decoration: BoxDecoration(
              color: AppTheme.background.withOpacity(0.95),
              border: const Border(top: BorderSide(color: AppTheme.border)),
            ),
            child: SafeArea(
              top: false,
              child: Row(
                children: [
                  // Bookmark button
                  Consumer<LibraryProvider>(
                    builder: (context, libProvider, _) {
                      final isBookmarked = libProvider.isBookmarked(novel.id);
                      return SizedBox(
                        width: 56,
                        height: 50,
                        child: OutlinedButton(
                          onPressed: () {
                            final auth = context.read<AuthProvider>();
                            if (auth.status != AuthStatus.authenticated) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                  content: const Text(
                                    'Please login to bookmark',
                                  ),
                                  action: SnackBarAction(
                                    label: 'Login',
                                    onPressed: () =>
                                        Navigator.pushNamed(context, '/login'),
                                  ),
                                ),
                              );
                              return;
                            }
                            libProvider.toggleBookmark(novel.id);
                          },
                          style: OutlinedButton.styleFrom(
                            side: BorderSide(
                              color: isBookmarked
                                  ? AppTheme.primary
                                  : AppTheme.border,
                            ),
                            padding: EdgeInsets.zero,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                          child: Icon(
                            isBookmarked
                                ? Icons.bookmark
                                : Icons.bookmark_border,
                            color: isBookmarked
                                ? AppTheme.primary
                                : AppTheme.muted,
                          ),
                        ),
                      );
                    },
                  ),
                  const SizedBox(width: 12),
                  // Read Now button
                  Expanded(
                    child: ElevatedButton(
                      onPressed: chapters.isNotEmpty
                          ? () {
                              Navigator.pushNamed(
                                context,
                                '/reader',
                                arguments: {
                                  'slug': novel.slug,
                                  'chapterNumber': chapters.first.number,
                                  'novelTitle': novel.title,
                                },
                              );
                            }
                          : null,
                      style: ElevatedButton.styleFrom(
                        minimumSize: const Size.fromHeight(50),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Icon(Icons.auto_stories, size: 20),
                          const SizedBox(width: 8),
                          Text(
                            chapters.isNotEmpty
                                ? 'Read Now'
                                : 'No Chapters Yet',
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }

  Widget _buildStatusBadge(String status) {
    Color color;
    switch (status.toLowerCase()) {
      case 'ongoing':
        color = Colors.greenAccent;
        break;
      case 'completed':
        color = Colors.blueAccent;
        break;
      case 'hiatus':
        color = Colors.orangeAccent;
        break;
      default:
        color = AppTheme.muted;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: color.withOpacity(0.15),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: color.withOpacity(0.4)),
      ),
      child: Text(
        status[0].toUpperCase() + status.substring(1),
        style: TextStyle(
          color: color,
          fontWeight: FontWeight.bold,
          fontSize: 11,
        ),
      ),
    );
  }

  String _formatViews(int views) {
    if (views >= 1000000) return '${(views / 1000000).toStringAsFixed(1)}M';
    if (views >= 1000) return '${(views / 1000).toStringAsFixed(1)}k';
    return '$views';
  }
}
