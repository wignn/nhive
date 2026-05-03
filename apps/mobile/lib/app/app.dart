import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:nhive/core/di/service_locator.dart';
import 'package:nhive/features/novels/presentation/pages/novel_list_page.dart';
import 'package:nhive/features/auth/presentation/pages/login_page.dart';
import 'package:nhive/features/auth/presentation/pages/register_page.dart';
import 'package:nhive/features/novels/presentation/pages/novel_detail_page.dart';
import 'package:nhive/features/novels/presentation/pages/chapter_reader_page.dart';
import 'package:nhive/features/notifications/presentation/pages/notifications_page.dart';
import 'package:nhive/features/auth/presentation/bloc/auth_provider.dart';
import 'package:nhive/features/library/presentation/bloc/library_provider.dart';
import 'package:nhive/app/theme/app_theme.dart';

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: sl.novelProvider),
        ChangeNotifierProvider.value(value: sl.libraryProvider),
        ChangeNotifierProxyProvider<LibraryProvider, AuthProvider>(
          create: (_) => sl.authProvider,
          update: (_, library, auth) {
            auth!.updateLibraryProvider(library);
            return auth;
          },
        ),
      ],
      child: MaterialApp(
        title: 'NovelHive',
        debugShowCheckedModeBanner: false,
        theme: AppTheme.darkTheme,
        initialRoute: '/',
        onGenerateRoute: (settings) {
          switch (settings.name) {
            case '/':
              return MaterialPageRoute(builder: (_) => const NovelListPage());
            case '/login':
              return MaterialPageRoute(builder: (_) => const LoginPage());
            case '/register':
              return MaterialPageRoute(builder: (_) => const RegisterPage());
            case '/detail':
              final slug = settings.arguments as String? ?? '';
              return MaterialPageRoute(
                builder: (_) => NovelDetailPage(slug: slug),
              );
            case '/reader':
              final args = settings.arguments as Map<String, dynamic>? ?? {};
              return MaterialPageRoute(
                builder: (_) => ChapterReaderPage(
                  slug: args['slug'] as String? ?? '',
                  chapterNumber: args['chapterNumber'] as int? ?? 1,
                  novelTitle: args['novelTitle'] as String? ?? '',
                ),
              );
            case '/notifications':
              return MaterialPageRoute(builder: (_) => const NotificationsPage());
            default:
              return MaterialPageRoute(builder: (_) => const NovelListPage());
          }
        },
      ),
    );
  }
}
