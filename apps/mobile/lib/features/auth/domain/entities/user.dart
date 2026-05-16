import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String id;
  final String username;
  final String email;
  final String role;
  final String avatarUrl;

  const User({
    required this.id,
    required this.username,
    required this.email,
    this.role = 'reader',
    this.avatarUrl = '',
  });

  User copyWith({
    String? id,
    String? username,
    String? email,
    String? role,
    String? avatarUrl,
  }) {
    return User(
      id: id ?? this.id,
      username: username ?? this.username,
      email: email ?? this.email,
      role: role ?? this.role,
      avatarUrl: avatarUrl ?? this.avatarUrl,
    );
  }

  @override
  List<Object?> get props => [id, username, email, role, avatarUrl];
}
