import '../entities/user.dart';

abstract class AuthRepository {
  Future<User> login(String email, String password);
  Future<void> register(String username, String email, String password);
  Future<User> getMe();
  Future<User> uploadAvatar(String filePath);
}
