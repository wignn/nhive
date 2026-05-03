class PaginationModel<T> {
  final List<T> items;
  final int total;
  final int page;

  PaginationModel({
    required this.items,
    required this.total,
    required this.page,
  });
}
