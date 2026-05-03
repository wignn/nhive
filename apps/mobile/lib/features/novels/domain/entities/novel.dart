import 'package:equatable/equatable.dart';

class Genre extends Equatable {
  final int id;
  final String name;
  final String slug;

  const Genre({required this.id, required this.name, required this.slug});

  factory Genre.fromJson(Map<String, dynamic> json) {
    return Genre(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      slug: json['slug'] ?? '',
    );
  }

  @override
  List<Object?> get props => [id, name, slug];
}

class Novel extends Equatable {
  final String id;
  final String title;
  final String slug;
  final String? description;
  final String? coverUrl;
  final String? author;
  final List<Genre> genres;
  final String? status;
  final int views;
  final int totalChapters;
  final String? updatedAt;
  final String? createdAt;

  const Novel({
    required this.id,
    required this.title,
    required this.slug,
    this.description,
    this.coverUrl,
    this.author,
    this.genres = const [],
    this.status,
    this.views = 0,
    this.totalChapters = 0,
    this.updatedAt,
    this.createdAt,
  });

  @override
  List<Object?> get props => [
    id,
    title,
    slug,
    description,
    coverUrl,
    author,
    genres,
    status,
    views,
    totalChapters,
    updatedAt,
    createdAt,
  ];
}
