import 'package:mobile/core/constants/api_constants.dart';
import '../../domain/entities/novel.dart';

class NovelModel extends Novel {
  const NovelModel({
    required super.id,
    required super.title,
    required super.slug,
    super.description,
    super.coverUrl,
    super.author,
    super.genres,
    super.status,
    super.views,
    super.totalChapters,
    super.updatedAt,
    super.createdAt,
  });

  factory NovelModel.fromJson(Map<String, dynamic> json) {
    final rawCover = json['cover_url'] as String?;
    final resolvedCover = ApiConstants.resolveCoverUrl(rawCover);

    final genresRaw = json['genres'];
    final genres = <Genre>[];
    if (genresRaw != null && genresRaw is List) {
      for (final g in genresRaw) {
        if (g is Map<String, dynamic>) {
          genres.add(Genre.fromJson(g));
        } else if (g is String) {
          genres.add(Genre(id: 0, name: g, slug: g.toLowerCase()));
        }
      }
    }

    return NovelModel(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      slug: json['slug'] ?? '',
      description: json['synopsis'] ?? json['description'] ?? '',
      coverUrl: resolvedCover,
      author: json['author'] ?? '',
      genres: genres,
      status: json['status'] ?? '',
      views: json['views'] ?? 0,
      totalChapters: json['total_chapters'] ?? 0,
      updatedAt: json['updated_at'],
      createdAt: json['created_at'],
    );
  }
}
