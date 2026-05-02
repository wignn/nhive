import 'package:mobile/features/novels/domain/repositories/novel_repository.dart';
import 'package:mobile/features/novels/data/datasources/novel_remote_datasource.dart';
import 'package:mobile/features/novels/domain/entities/novel.dart';
import 'package:mobile/features/novels/domain/entities/chapter.dart';
import 'package:mobile/features/novels/data/models/novel_model.dart';
import 'package:mobile/features/novels/data/models/chapter_model.dart';

class NovelRepositoryImpl implements NovelRepository {
  final NovelRemoteDataSource remoteDataSource;

  NovelRepositoryImpl(this.remoteDataSource);

  @override
  Future<List<Novel>> getNovels({int page = 1, int pageSize = 18, String sort = 'updated'}) async {
    final data = await remoteDataSource.getNovels(page: page, pageSize: pageSize, sort: sort);
    final novels = data['novels'] as List? ?? [];
    return novels.map((j) => NovelModel.fromJson(j)).toList();
  }

  @override
  Future<Novel> getNovelDetail(String slug) async {
    final data = await remoteDataSource.getNovelDetail(slug);
    final novelJson = data['novel'] ?? data;
    return NovelModel.fromJson(novelJson);
  }

  @override
  Future<List<Chapter>> getChapters(String slug) async {
    final data = await remoteDataSource.getChapters(slug);
    final chapters = data['chapters'] as List? ?? [];
    return chapters.map((j) => ChapterModel.fromJson(j)).toList();
  }

  @override
  Future<Chapter> getChapterDetail(String slug, int number) async {
    final data = await remoteDataSource.getChapterDetail(slug, number);
    final chapterJson = data['chapter'] ?? data;
    return ChapterModel.fromJson(chapterJson);
  }

  @override
  Future<List<Novel>> searchNovels(String query) async {
    final data = await remoteDataSource.searchNovels(query);
    final novels = data['novels'] as List? ?? [];
    return novels.map((j) => NovelModel.fromJson(j)).toList();
  }
}
