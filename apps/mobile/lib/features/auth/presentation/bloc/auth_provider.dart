import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:mobile/features/auth/domain/entities/user.dart';
import 'package:mobile/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:mobile/core/storage/secure_storage.dart';

enum AuthStatus { initial, loading, authenticated, unauthenticated, error }

class AuthProvider extends ChangeNotifier {
  final AuthRepositoryImpl _repository;
  final SecureStorage _storage;

  AuthProvider(this._repository, this._storage);

  AuthStatus _status = AuthStatus.initial;
  AuthStatus get status => _status;

  User? _user;
  User? get user => _user;

  String? _error;
  String? get error => _error;

  Future<void> checkAuth() async {
    final token = await _storage.getToken();
    if (token != null && token.isNotEmpty) {
      try {
        _user = await _repository.getMe();
        _status = AuthStatus.authenticated;
      } catch (e) {
        _status = AuthStatus.unauthenticated;
      }
    } else {
      _status = AuthStatus.unauthenticated;
    }
    notifyListeners();
  }

  Future<bool> login(String email, String password) async {
    _status = AuthStatus.loading;
    _error = null;
    notifyListeners();

    try {
      _user = await _repository.login(email, password);
      _status = AuthStatus.authenticated;
      notifyListeners();
      return true;
    } catch (e) {
      _error = _parseError(e);
      _status = AuthStatus.error;
      notifyListeners();
      return false;
    }
  }

  Future<bool> register(String username, String email, String password) async {
    _status = AuthStatus.loading;
    _error = null;
    notifyListeners();

    try {
      await _repository.register(username, email, password);
      // Auto-login after registration
      return await login(email, password);
    } catch (e) {
      _error = _parseError(e);
      _status = AuthStatus.error;
      notifyListeners();
      return false;
    }
  }

  Future<void> logout() async {
    await _storage.clearAll();
    _user = null;
    _status = AuthStatus.unauthenticated;
    notifyListeners();
  }

  String _parseError(dynamic e) {
    final msg = e.toString();
    if (msg.contains('401')) return 'Invalid email or password';
    if (msg.contains('409')) return 'Account already exists';
    if (msg.contains('connection')) return 'No internet connection';
    return 'Something went wrong. Please try again.';
  }
}
