import 'package:mobile/core/network/dio_client.dart';
import 'package:mobile/core/constants/api_constants.dart';

abstract class NovelRemoteDataSource {
  Future<Map<String, dynamic>> getNovels({int page = 1, int pageSize = 18, String sort = 'updated'});
  Future<Map<String, dynamic>> getNovelDetail(String slug);
  Future<Map<String, dynamic>> getChapters(String slug);
  Future<Map<String, dynamic>> getChapterDetail(String slug, int number);
  Future<Map<String, dynamic>> searchNovels(String query);
}

class NovelRemoteDataSourceImpl implements NovelRemoteDataSource {
  final DioClient _client;

  NovelRemoteDataSourceImpl(this._client);

  @override
  Future<Map<String, dynamic>> getNovels({int page = 1, int pageSize = 18, String sort = 'updated'}) async {
    final response = await _client.get(ApiConstants.novels, queryParameters: {
      'page': page,
      'page_size': pageSize,
      'sort': sort,
    });
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> getNovelDetail(String slug) async {
    final response = await _client.get('${ApiConstants.novels}/$slug');
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> getChapters(String slug) async {
    final response = await _client.get('${ApiConstants.novels}/$slug/chapters');
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> getChapterDetail(String slug, int number) async {
    final response = await _client.get('${ApiConstants.novels}/$slug/chapters/$number');
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> searchNovels(String query) async {
    final response = await _client.get(ApiConstants.search, queryParameters: {'q': query});
    return response.data;
  }
}
