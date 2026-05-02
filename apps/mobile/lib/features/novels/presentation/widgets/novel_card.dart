import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';
import 'package:mobile/app/theme/app_theme.dart';
import 'package:mobile/features/novels/domain/entities/novel.dart';

class NovelCard extends StatelessWidget {
  final Novel novel;
  final VoidCallback? onTap;
  final bool showAuthor;

  const NovelCard({super.key, required this.novel, this.onTap, this.showAuthor = true});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(12),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.3),
                    blurRadius: 10,
                    offset: const Offset(0, 4),
                  ),
                ],
              ),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(12),
                child: _buildCover(),
              ),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            novel.title,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 13,
              color: AppTheme.foreground,
            ),
          ),
          if (showAuthor) ...[
            const SizedBox(height: 2),
            Text(
              novel.author ?? 'Unknown Author',
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(
                fontSize: 11,
                color: AppTheme.muted,
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildCover() {
    final url = novel.coverUrl ?? '';
    if (url.isEmpty) {
      return _buildPlaceholder();
    }
    return CachedNetworkImage(
      imageUrl: url,
      fit: BoxFit.cover,
      width: double.infinity,
      height: double.infinity,
      placeholder: (context, url) => Shimmer.fromColors(
        baseColor: AppTheme.secondary,
        highlightColor: AppTheme.surface,
        child: Container(color: AppTheme.secondary),
      ),
      errorWidget: (context, url, error) => _buildPlaceholder(),
    );
  }

  Widget _buildPlaceholder() {
    return Container(
      color: AppTheme.secondary,
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.auto_stories_rounded, color: AppTheme.muted, size: 32),
            const SizedBox(height: 4),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8),
              child: Text(
                novel.title,
                maxLines: 2,
                textAlign: TextAlign.center,
                overflow: TextOverflow.ellipsis,
                style: const TextStyle(color: AppTheme.muted, fontSize: 10),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
