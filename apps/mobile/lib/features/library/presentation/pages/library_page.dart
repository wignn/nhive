import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:mobile/app/theme/app_theme.dart';
import 'package:mobile/features/library/presentation/bloc/library_provider.dart';
import 'package:mobile/features/auth/presentation/bloc/auth_provider.dart';
import 'package:mobile/features/novels/presentation/widgets/novel_card.dart';

class LibraryPage extends StatefulWidget {
  const LibraryPage({super.key});

  @override
  State<LibraryPage> createState() => _LibraryPageState();
}

class _LibraryPageState extends State<LibraryPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final auth = context.read<AuthProvider>();
      if (auth.status == AuthStatus.authenticated) {
        context.read<LibraryProvider>().loadLibrary();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Consumer2<AuthProvider, LibraryProvider>(
      builder: (context, auth, library, _) {
        if (auth.status != AuthStatus.authenticated) {
          return _buildLoginPrompt(context);
        }

        return Scaffold(
          appBar: AppBar(
            title: const Text('My Library'),
            actions: [
              if (library.bookmarks.isNotEmpty)
                Padding(
                  padding: const EdgeInsets.only(right: 16),
                  child: Center(
                    child: Container(
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                      decoration: BoxDecoration(
                        color: AppTheme.primary.withOpacity(0.15),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Text(
                        '${library.bookmarks.length}',
                        style: const TextStyle(
                          color: AppTheme.primary,
                          fontWeight: FontWeight.bold,
                          fontSize: 13,
                        ),
                      ),
                    ),
                  ),
                ),
            ],
          ),
          body: RefreshIndicator(
            color: AppTheme.primary,
            onRefresh: () => library.loadLibrary(),
            child: library.isLoading
                ? const Center(child: CircularProgressIndicator(color: AppTheme.primary))
                : library.bookmarks.isEmpty
                    ? _buildEmptyLibrary()
                    : GridView.builder(
                        padding: const EdgeInsets.all(16),
                        physics: const AlwaysScrollableScrollPhysics(),
                        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 3,
                          mainAxisSpacing: 16,
                          crossAxisSpacing: 12,
                          childAspectRatio: 0.52,
                        ),
                        itemCount: library.bookmarks.length,
                        itemBuilder: (context, index) {
                          final novel = library.bookmarks[index];
                          return NovelCard(
                            novel: novel,
                            onTap: () => Navigator.pushNamed(context, '/detail', arguments: novel.slug),
                          );
                        },
                      ),
          ),
        );
      },
    );
  }

  Widget _buildLoginPrompt(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('My Library')),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 80,
                height: 80,
                decoration: BoxDecoration(
                  color: AppTheme.primary.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: const Icon(Icons.bookmark_outline, size: 40, color: AppTheme.primary),
              ),
              const SizedBox(height: 20),
              const Text(
                'Your Library',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 8),
              const Text(
                'Sign in to save and manage\nyour bookmarked novels.',
                textAlign: TextAlign.center,
                style: TextStyle(color: AppTheme.muted, fontSize: 14, height: 1.5),
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: () => Navigator.pushNamed(context, '/login'),
                child: const Text('Sign In'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildEmptyLibrary() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.collections_bookmark_outlined, size: 64, color: AppTheme.muted.withOpacity(0.4)),
          const SizedBox(height: 16),
          const Text(
            'No bookmarks yet',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600),
          ),
          const SizedBox(height: 8),
          const Text(
            'Novels you bookmark will appear here',
            style: TextStyle(color: AppTheme.muted, fontSize: 14),
          ),
        ],
      ),
    );
  }
}
