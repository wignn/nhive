import '../repositories/novel_repository.dart';
import '../entities/chapter.dart';

class GetChapterUseCase {
  final NovelRepository repository;

  GetChapterUseCase(this.repository);

  Future<Chapter> execute(String id) async {
    return await repository.getChapter(id);
  }
}
