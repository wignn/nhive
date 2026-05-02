import 'package:equatable/equatable.dart';

class Chapter extends Equatable {
  final String id;
  final int number;
  final String title;
  final String? content;
  final int wordCount;
  final String? createdAt;

  const Chapter({
    required this.id,
    this.number = 0,
    required this.title,
    this.content,
    this.wordCount = 0,
    this.createdAt,
  });

  @override
  List<Object?> get props => [id, number, title, content, wordCount, createdAt];
}
