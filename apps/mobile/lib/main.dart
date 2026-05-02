import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:mobile/core/di/service_locator.dart';
import 'app/app.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  sl.init();
  runApp(const MyApp());
}
