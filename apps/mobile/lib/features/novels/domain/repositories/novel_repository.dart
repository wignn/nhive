import '../entities/novel.dart';
import '../entities/chapter.dart';

abstract class NovelRepository {
  Future<List<Novel>> getNovels({int page = 1, int pageSize = 18, String sort = 'updated'});
  Future<Novel> getNovelDetail(String slug);
  Future<List<Chapter>> getChapters(String slug);
  Future<Chapter> getChapterDetail(String slug, int number);
  Future<List<Novel>> searchNovels(String query);
}
