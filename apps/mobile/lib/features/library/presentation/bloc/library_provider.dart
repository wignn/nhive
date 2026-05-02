import 'package:flutter/foundation.dart';
import 'package:mobile/core/network/dio_client.dart';
import 'package:mobile/core/constants/api_constants.dart';
import 'package:mobile/features/novels/domain/entities/novel.dart';
import 'package:mobile/features/novels/data/models/novel_model.dart';

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

  Future<void> loadLibrary() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _client.get(ApiConstants.library);
      final data = response.data;
      final coverBaseUrl = data['cover_base_url'] as String?;
      final novels = data['novels'] as List? ?? data['bookmarks'] as List? ?? [];
      _bookmarks = novels.map<Novel>((j) => NovelModel.fromJson(j, coverBaseUrl: coverBaseUrl)).toList();
      _bookmarkedIds = _bookmarks.map((n) => n.id).toSet();
    } catch (e) {
      _error = e.toString();
      _bookmarks = [];
    }

    _isLoading = false;
    notifyListeners();
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
