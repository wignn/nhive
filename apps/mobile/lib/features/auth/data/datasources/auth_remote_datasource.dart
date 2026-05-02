import 'package:mobile/core/network/dio_client.dart';
import 'package:mobile/core/constants/api_constants.dart';

abstract class AuthRemoteDataSource {
  Future<Map<String, dynamic>> login(String email, String password);
  Future<Map<String, dynamic>> register(String username, String email, String password);
  Future<Map<String, dynamic>> getMe();
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final DioClient _client;

  AuthRemoteDataSourceImpl(this._client);

  @override
  Future<Map<String, dynamic>> login(String email, String password) async {
    final response = await _client.post(ApiConstants.login, data: {
      'email': email,
      'password': password,
    });
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> register(String username, String email, String password) async {
    final response = await _client.post(ApiConstants.register, data: {
      'username': username,
      'email': email,
      'password': password,
    });
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> getMe() async {
    final response = await _client.get(ApiConstants.me);
    return response.data;
  }
}
