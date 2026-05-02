class ApiConstants {
  static const String baseUrl = 'https://api-novel.wign.cloud';
  static const String internalKey = 'silviavoniliade';
  static const String coverBaseUrl = 'https://pub-f47f8acab5f64bc28ad7f0095c8a9124.r2.dev';

  // Auth
  static const String login = '/api/v1/auth/login';
  static const String register = '/api/v1/auth/register';
  static const String me = '/api/v1/auth/me';

  // Novels
  static const String novels = '/api/v1/novels';
  static const String library = '/api/v1/library';
  static const String progress = '/api/v1/progress';
  static const String search = '/api/v1/search';
  static const String genres = '/api/v1/genres';

  /// Resolves a potentially relative cover URL to a full URL.
  static String resolveCoverUrl(String? rawCoverUrl) {
    if (rawCoverUrl == null || rawCoverUrl.isEmpty) return '';
    if (rawCoverUrl.startsWith('http') || rawCoverUrl.startsWith('data:')) return rawCoverUrl;
    try {
      final base = Uri.parse(coverBaseUrl);
      final resolved = base.resolve(rawCoverUrl);
      return resolved.toString();
    } catch (_) {
      final base = coverBaseUrl.replaceFirst(RegExp(r'/$'), '');
      final path = rawCoverUrl.replaceFirst(RegExp(r'^/'), '');
      return '$base/$path';
    }
  }
}
