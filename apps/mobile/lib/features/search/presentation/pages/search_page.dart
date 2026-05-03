import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:nhive/app/theme/app_theme.dart';
import 'package:nhive/features/novels/presentation/bloc/novel_provider.dart';
import 'package:nhive/features/novels/presentation/widgets/novel_card.dart';
import 'package:nhive/features/novels/domain/entities/novel.dart';

class SearchPage extends StatefulWidget {
  const SearchPage({super.key});

  @override
  State<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends State<SearchPage> {
  final _searchController = TextEditingController();
  final _focusNode = FocusNode();
  Timer? _debounce;
  String? _selectedGenre;

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
                  style: const TextStyle(
                    color: AppTheme.foreground,
                    fontSize: 16,
                  ),
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
                            icon: const Icon(
                              Icons.close,
                              color: AppTheme.muted,
                              size: 20,
                            ),
                          )
                        : null,
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 14,
                    ),
                  ),
                ),
              ),
            ),

            // ─── Genre Filters ───
            Consumer<NovelProvider>(
              builder: (context, provider, _) {
                final allGenres = <String>{};
                for (var n in provider.novels) {
                  for (var g in n.genres) {
                    allGenres.add(g.name);
                  }
                }
                final genresList = allGenres.toList()..sort();

                if (genresList.isEmpty) return const SizedBox.shrink();

                return SizedBox(
                  height: 40,
                  child: ListView.separated(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    scrollDirection: Axis.horizontal,
                    itemCount: genresList.length,
                    separatorBuilder: (_, __) => const SizedBox(width: 8),
                    itemBuilder: (context, index) {
                      final genre = genresList[index];
                      final isSelected = _selectedGenre == genre;
                      return ChoiceChip(
                        label: Text(genre),
                        selected: isSelected,
                        onSelected: (selected) {
                          setState(() {
                            _selectedGenre = selected ? genre : null;
                          });
                        },
                        selectedColor: AppTheme.primary.withOpacity(0.2),
                        backgroundColor: AppTheme.surface,
                        labelStyle: TextStyle(
                          color: isSelected ? AppTheme.primary : AppTheme.muted,
                          fontWeight: isSelected
                              ? FontWeight.bold
                              : FontWeight.normal,
                        ),
                        side: BorderSide(
                          color: isSelected
                              ? AppTheme.primary
                              : AppTheme.border,
                        ),
                      );
                    },
                  ),
                );
              },
            ),
            const SizedBox(height: 8),

            // ─── Results ───
            Expanded(
              child: Consumer<NovelProvider>(
                builder: (context, provider, _) {
                  List<Novel> displayNovels = provider.searchQuery.isEmpty
                      ? provider.novels
                      : provider.searchResults;

                  if (_selectedGenre != null) {
                    displayNovels = displayNovels
                        .where(
                          (n) => n.genres.any((g) => g.name == _selectedGenre),
                        )
                        .toList();
                  }

                  if (provider.isSearching) {
                    return const Center(
                      child: CircularProgressIndicator(color: AppTheme.primary),
                    );
                  }

                  if (displayNovels.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.search_off,
                            size: 64,
                            color: AppTheme.muted.withOpacity(0.4),
                          ),
                          const SizedBox(height: 16),
                          Text(
                            'No novels found',
                            style: const TextStyle(
                              color: AppTheme.muted,
                              fontSize: 16,
                            ),
                          ),
                        ],
                      ),
                    );
                  }

                  return GridView.builder(
                    padding: const EdgeInsets.all(16),
                    gridDelegate:
                        const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 3,
                          mainAxisSpacing: 16,
                          crossAxisSpacing: 12,
                          childAspectRatio: 0.52,
                        ),
                    itemCount: displayNovels.length,
                    itemBuilder: (context, index) {
                      final novel = displayNovels[index];
                      return NovelCard(
                        novel: novel,
                        onTap: () => Navigator.pushNamed(
                          context,
                          '/detail',
                          arguments: novel.slug,
                        ),
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
}
