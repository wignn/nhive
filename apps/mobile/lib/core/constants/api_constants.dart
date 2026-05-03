class ApiConstants {
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8000',
  );

  static const String internalKey = String.fromEnvironment(
    'INTERNAL_KEY',
    defaultValue: '',
  );

  static const String apiKey = String.fromEnvironment(
    'API_KEY',
    defaultValue: '',
  );

  static const String login = '/api/v1/auth/login';
  static const String register = '/api/v1/auth/register';
  static const String me = '/api/v1/auth/me';

  // Novels
  static const String novels = '/api/v1/novels';
  static const String library = '/api/v1/library';
  static const String progress = '/api/v1/progress';
  static const String search = '/api/v1/search';
  static const String genres = '/api/v1/genres';
  static const String notifications = '/api/v1/notifications';
}
