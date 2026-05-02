import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String id;
  final String username;
  final String email;
  final String role;

  const User({
    required this.id,
    required this.username,
    required this.email,
    this.role = 'reader',
  });

  @override
  List<Object?> get props => [id, username, email, role];
}
