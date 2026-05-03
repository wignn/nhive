import 'package:flutter/foundation.dart';
import 'package:nhive/features/novels/domain/entities/novel.dart';
import 'package:nhive/features/novels/domain/entities/chapter.dart';
import 'package:nhive/features/novels/domain/repositories/novel_repository.dart';

class NovelProvider extends ChangeNotifier {
  final NovelRepository _repository;

  NovelProvider(this._repository);

  // ─── Home novels ───
  List<Novel> _novels = [];
  List<Novel> get novels => _novels;

  bool _isLoadingNovels = false;
  bool get isLoadingNovels => _isLoadingNovels;

  String? _novelsError;
  String? get novelsError => _novelsError;

  int _currentPage = 1;
  bool _hasMore = true;
  bool get hasMore => _hasMore;

  Future<void> loadNovels({bool refresh = false}) async {
    if (_isLoadingNovels) return;
    if (refresh) {
      _currentPage = 1;
      _hasMore = true;
      _novels = [];
    }
    if (!_hasMore) return;

    _isLoadingNovels = true;
    _novelsError = null;
    notifyListeners();

    try {
      final result = await _repository.getNovels(
        page: _currentPage,
        pageSize: 18,
      );
      if (refresh) {
        _novels = result;
      } else {
        _novels = [..._novels, ...result];
      }
      _hasMore = result.length >= 18;
      _currentPage++;
    } catch (e) {
      _novelsError = e.toString();
    }

    _isLoadingNovels = false;
    notifyListeners();
  }

  // ─── Novel Detail ───
  Novel? _selectedNovel;
  Novel? get selectedNovel => _selectedNovel;

  bool _isLoadingDetail = false;
  bool get isLoadingDetail => _isLoadingDetail;

  String? _detailError;
  String? get detailError => _detailError;

  List<Chapter> _chapters = [];
  List<Chapter> get chapters => _chapters;

  bool _isLoadingChapters = false;
  bool get isLoadingChapters => _isLoadingChapters;

  Future<void> loadNovelDetail(String slug) async {
    _isLoadingDetail = true;
    _detailError = null;
    _selectedNovel = null;
    _chapters = [];
    notifyListeners();

    try {
      final results = await Future.wait([
        _repository.getNovelDetail(slug),
        _repository.getChapters(slug),
      ]);
      _selectedNovel = results[0] as Novel;
      _chapters = results[1] as List<Chapter>;
    } catch (e) {
      _detailError = e.toString();
    }

    _isLoadingDetail = false;
    notifyListeners();
  }

  // ─── Chapter Reader ───
  Chapter? _currentChapter;
  Chapter? get currentChapter => _currentChapter;

  bool _isLoadingChapter = false;
  bool get isLoadingChapter => _isLoadingChapter;

  String? _chapterError;
  String? get chapterError => _chapterError;

  String? _currentNovelSlug;
  String? get currentNovelSlug => _currentNovelSlug;

  Future<void> loadChapter(String slug, int number) async {
    _isLoadingChapter = true;
    _chapterError = null;
    _currentNovelSlug = slug;
    notifyListeners();

    try {
      _currentChapter = await _repository.getChapterDetail(slug, number);
    } catch (e) {
      _chapterError = e.toString();
    }

    _isLoadingChapter = false;
    notifyListeners();
  }

  // ─── Search ───
  List<Novel> _searchResults = [];
  List<Novel> get searchResults => _searchResults;

  bool _isSearching = false;
  bool get isSearching => _isSearching;

  String _searchQuery = '';
  String get searchQuery => _searchQuery;

  Future<void> search(String query) async {
    if (query.trim().isEmpty) {
      _searchResults = [];
      _searchQuery = '';
      notifyListeners();
      return;
    }

    _isSearching = true;
    _searchQuery = query;
    notifyListeners();

    try {
      _searchResults = await _repository.searchNovels(query);
    } catch (e) {
      _searchResults = [];
    }

    _isSearching = false;
    notifyListeners();
  }

  void clearSearch() {
    _searchResults = [];
    _searchQuery = '';
    notifyListeners();
  }
}
