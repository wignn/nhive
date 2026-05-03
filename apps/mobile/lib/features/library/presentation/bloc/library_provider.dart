import 'package:flutter/foundation.dart';
import 'package:nhive/core/network/dio_client.dart';
import 'package:nhive/core/constants/api_constants.dart';
import 'package:nhive/features/novels/domain/entities/novel.dart';
import 'package:nhive/features/novels/data/models/novel_model.dart';

class LibraryProvider extends ChangeNotifier {
  final DioClient _client;

  LibraryProvider(this._client);

  List<Novel> _bookmarks = [];
  List<Novel> get bookmarks => _bookmarks;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _error;
  String? get error => _error;

  Set<String> _bookmarkedIds = {};

  bool isBookmarked(String novelId) => _bookmarkedIds.contains(novelId);

  void clear() {
    _bookmarks = [];
    _bookmarkedIds = {};
    _error = null;
    _isLoading = false;
    notifyListeners();
  }

  Future<void> loadLibrary() async {
    if (_isLoading) return;

    // Reset state but keep old data to prevent flickering if possible
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _client.get(ApiConstants.library);
      final data = response.data;

      final List novelsData = data is List
          ? data
          : (data['entries'] as List? ?? []);
      final coverBaseUrl = data is Map ? data['cover_base_url'] as String? : null;

      _bookmarks = novelsData.map<Novel>((j) {
        // Map Gateway's flat library entry format to what NovelModel.fromJson expects
        final Map<String, dynamic> novelJson = {
          'id': j['novel_id'] ?? j['id'] ?? '',
          'title': j['novel_title'] ?? j['title'] ?? 'Unknown',
          'slug': j['novel_slug'] ?? j['slug'] ?? '',
          'cover_url': j['cover_url'],
          'author': j['author'] ?? '',
          'total_chapters': j['total'] ?? j['total_chapters'] ?? 0,
        };
        return NovelModel.fromJson(novelJson, coverBaseUrl: coverBaseUrl);
      }).toList();

      _bookmarkedIds = _bookmarks.map((n) => n.id).toSet();
    } catch (e) {
      _error = e.toString();
      print('Library Load Error: $e');
      // If unauthorized (401), we just clear bookmarks silently
      if (e.toString().contains('401')) {
        _bookmarks = [];
        _bookmarkedIds = {};
      }
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> addBookmark(String novelId) async {
    try {
      await _client.post('${ApiConstants.library}/$novelId');
      _bookmarkedIds.add(novelId);
      notifyListeners();
      await loadLibrary();
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> removeBookmark(String novelId) async {
    try {
      await _client.delete('${ApiConstants.library}/$novelId');
      _bookmarkedIds.remove(novelId);
      _bookmarks.removeWhere((n) => n.id == novelId);
      notifyListeners();
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<void> toggleBookmark(String novelId) async {
    if (isBookmarked(novelId)) {
      await removeBookmark(novelId);
    } else {
      await addBookmark(novelId);
    }
  }
}
