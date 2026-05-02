import '../repositories/novel_repository.dart';
import '../entities/novel.dart';

class GetNovelDetailUseCase {
  final NovelRepository repository;

  GetNovelDetailUseCase(this.repository);

  Future<Novel> execute(String id) async {
    return await repository.getNovelDetail(id);
  }
}
