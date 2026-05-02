import '../../domain/entities/chapter.dart';

class ChapterModel extends Chapter {
  const ChapterModel({
    required super.id,
    super.number,
    required super.title,
    super.content,
    super.wordCount,
    super.createdAt,
  });

  factory ChapterModel.fromJson(Map<String, dynamic> json) {
    return ChapterModel(
      id: json['id'] ?? '',
      number: json['number'] ?? 0,
      title: json['title'] ?? '',
      content: json['content'] ?? '',
      wordCount: json['word_count'] ?? 0,
      createdAt: json['created_at'],
    );
  }
}
