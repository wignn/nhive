import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/app/app.dart';

void main() {
  testWidgets('App initialization test', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(const MyApp());

    // Verify that the app title is present.
    expect(find.text('NovelHive'), findsWidgets);
  });
}
