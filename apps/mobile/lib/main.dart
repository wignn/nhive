import 'package:flutter/material.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:nhive/core/di/service_locator.dart';
import 'app/app.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize Firebase
  await Firebase.initializeApp();
  
  // Initialize Service Locator
  sl.init();
  
  // Start Notification Service
  await sl.notificationService.init();
  
  runApp(const MyApp());
}
