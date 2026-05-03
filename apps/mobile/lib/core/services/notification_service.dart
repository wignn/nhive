import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:nhive/core/network/dio_client.dart';
import 'package:nhive/core/storage/secure_storage.dart';
import 'dart:io';

class NotificationService {
  final FirebaseMessaging _fcm = FirebaseMessaging.instance;
  final DioClient _dio;
  final SecureStorage _storage;

  NotificationService(this._dio, this._storage);

  Future<void> init() async {
    // 1. Request permission
    NotificationSettings settings = await _fcm.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    if (settings.authorizationStatus == AuthorizationStatus.authorized) {
      print('User granted permission');
      
      // 2. Get token
      String? token = await _fcm.getToken();
      if (token != null) {
        await _registerToken(token);
      }

      // 3. Listen for token refresh
      _fcm.onTokenRefresh.listen((newToken) {
        _registerToken(newToken);
      });

      // 4. Handle background messages
      FirebaseMessaging.onBackgroundMessage(_firebaseMessagingBackgroundHandler);

      // 5. Handle foreground messages
      FirebaseMessaging.onMessage.listen((RemoteMessage message) {
        print('Received foreground message: ${message.notification?.title}');
        // Here you can show a local notification using flutter_local_notifications if desired
      });
    }
  }

  Future<void> _registerToken(String token) async {
    // Only register if we have a JWT token (user is logged in)
    final jwtToken = await _storage.getToken();
    if (jwtToken == null || jwtToken.isEmpty) return;

    try {
      await _dio.post('/notifications/fcm-token', data: {'token': token});
      print('FCM Token registered successfully');
    } catch (e) {
      print('Failed to register FCM Token: $e');
    }
  }

  // Mandatory static handler for background messages
  static Future<void> _firebaseMessagingBackgroundHandler(RemoteMessage message) async {
    print('Handling background message: ${message.messageId}');
  }
}
