import '../repositories/novel_repository.dart';
import '../entities/chapter.dart';

class GetChapterUseCase {
  final NovelRepository repository;

  GetChapterUseCase(this.repository);

  Future<Chapter> execute(String slug, int number) async {
    return repository.getChapterDetail(slug, number);
  }
}
