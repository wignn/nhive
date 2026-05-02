import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:mobile/app/theme/app_theme.dart';
import 'package:mobile/features/novels/presentation/bloc/novel_provider.dart';
import 'package:mobile/features/novels/presentation/widgets/novel_card.dart';

class SearchPage extends StatefulWidget {
  const SearchPage({super.key});

  @override
  State<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends State<SearchPage> {
  final _searchController = TextEditingController();
  final _focusNode = FocusNode();
  Timer? _debounce;

  @override
  void dispose() {
    _searchController.dispose();
    _focusNode.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onSearchChanged(String query) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      context.read<NovelProvider>().search(query);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            // ─── Search Bar ───
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
              child: Container(
                decoration: BoxDecoration(
                  color: AppTheme.surface,
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(color: AppTheme.border),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.15),
                      blurRadius: 10,
                      offset: const Offset(0, 4),
                    ),
                  ],
                ),
                child: TextField(
                  controller: _searchController,
                  focusNode: _focusNode,
                  onChanged: _onSearchChanged,
                  style: const TextStyle(color: AppTheme.foreground, fontSize: 16),
                  decoration: InputDecoration(
                    hintText: 'Search novels, authors...',
                    hintStyle: const TextStyle(color: AppTheme.muted),
                    prefixIcon: const Icon(Icons.search, color: AppTheme.muted),
                    suffixIcon: _searchController.text.isNotEmpty
                        ? IconButton(
                            onPressed: () {
                              _searchController.clear();
                              context.read<NovelProvider>().clearSearch();
                              setState(() {});
                            },
                            icon: const Icon(Icons.close, color: AppTheme.muted, size: 20),
                          )
                        : null,
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                  ),
                ),
              ),
            ),

            // ─── Results ───
            Expanded(
              child: Consumer<NovelProvider>(
                builder: (context, provider, _) {
                  if (provider.searchQuery.isEmpty) {
                    return _buildEmptyState();
                  }

                  if (provider.isSearching) {
                    return const Center(
                      child: CircularProgressIndicator(color: AppTheme.primary),
                    );
                  }

                  if (provider.searchResults.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(Icons.search_off, size: 64, color: AppTheme.muted.withOpacity(0.4)),
                          const SizedBox(height: 16),
                          Text(
                            'No results for "${provider.searchQuery}"',
                            style: const TextStyle(color: AppTheme.muted, fontSize: 16),
                          ),
                          const SizedBox(height: 8),
                          const Text(
                            'Try different keywords',
                            style: TextStyle(color: AppTheme.muted, fontSize: 13),
                          ),
                        ],
                      ),
                    );
                  }

                  return GridView.builder(
                    padding: const EdgeInsets.all(16),
                    gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: 3,
                      mainAxisSpacing: 16,
                      crossAxisSpacing: 12,
                      childAspectRatio: 0.52,
                    ),
                    itemCount: provider.searchResults.length,
                    itemBuilder: (context, index) {
                      final novel = provider.searchResults[index];
                      return NovelCard(
                        novel: novel,
                        onTap: () => Navigator.pushNamed(context, '/detail', arguments: novel.slug),
                      );
                    },
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
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
            child: const Icon(Icons.search, size: 40, color: AppTheme.primary),
          ),
          const SizedBox(height: 20),
          const Text(
            'Discover Novels',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 8),
          const Text(
            'Search by title, author, or genre',
            style: TextStyle(color: AppTheme.muted, fontSize: 14),
          ),
        ],
      ),
    );
  }
}
