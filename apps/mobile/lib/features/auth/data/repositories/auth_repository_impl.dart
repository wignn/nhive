import 'package:mobile/core/storage/secure_storage.dart';
import 'package:mobile/features/auth/domain/entities/user.dart';
import 'package:mobile/features/auth/domain/repositories/auth_repository.dart';
import 'package:mobile/features/auth/data/datasources/auth_remote_datasource.dart';
import 'package:mobile/features/auth/data/models/user_model.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource remoteDataSource;
  final SecureStorage secureStorage;

  AuthRepositoryImpl(this.remoteDataSource, this.secureStorage);

  @override
  Future<User> login(String email, String password) async {
    final data = await remoteDataSource.login(email, password);
    final accessToken = data['access_token'];
    final refreshToken = data['refresh_token'];
    
    if (accessToken != null) {
      await secureStorage.saveToken(accessToken);
    }
    if (refreshToken != null) {
      await secureStorage.saveRefreshToken(refreshToken);
    }

    final userJson = data['user'];
    if (userJson != null) {
      return UserModel.fromJson(userJson);
    } else {
      return await getMe();
    }
  }

  @override
  Future<void> register(String username, String email, String password) async {
    await remoteDataSource.register(username, email, password);
  }

  @override
  Future<User> getMe() async {
    final data = await remoteDataSource.getMe();
    final user = UserModel.fromJson(data);
    return user;
  }
}
