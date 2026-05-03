import 'package:flutter_test/flutter_test.dart';
import 'package:nhive/core/di/service_locator.dart';
import 'package:nhive/features/auth/presentation/pages/login_page.dart';

void main() {
  testWidgets('App initialization test', (WidgetTester tester) async {
    TestWidgetsFlutterBinding.ensureInitialized();
    sl.init();



    expect(find.byType(LoginPage), findsOneWidget);
  });
}
