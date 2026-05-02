import '../repositories/novel_repository.dart';
import '../entities/novel.dart';

class GetNovelsUseCase {
  final NovelRepository repository;

  GetNovelsUseCase(this.repository);

  Future<List<Novel>> execute() async {
    return await repository.getNovels();
  }
}
