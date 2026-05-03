import 'package:nhive/core/storage/secure_storage.dart';
import 'package:nhive/core/network/auth_interceptor.dart';
import 'package:nhive/core/network/dio_client.dart';
import 'package:nhive/core/services/notification_service.dart';
import 'package:nhive/features/auth/data/datasources/auth_remote_datasource.dart';
import 'package:nhive/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:nhive/features/novels/data/datasources/novel_remote_datasource.dart';
import 'package:nhive/features/novels/data/repositories/novel_repository_impl.dart';
import 'package:nhive/features/novels/presentation/bloc/novel_provider.dart';
import 'package:nhive/features/auth/presentation/bloc/auth_provider.dart';
import 'package:nhive/features/library/presentation/bloc/library_provider.dart';

class ServiceLocator {
  static final ServiceLocator _instance = ServiceLocator._internal();
  factory ServiceLocator() => _instance;
  ServiceLocator._internal();

  late final SecureStorage secureStorage;
  late final AuthInterceptor authInterceptor;
  late final DioClient dioClient;
  late final NotificationService notificationService;
  late final AuthRemoteDataSourceImpl authRemoteDataSource;
  late final AuthRepositoryImpl authRepository;
  late final NovelRemoteDataSourceImpl novelRemoteDataSource;
  late final NovelRepositoryImpl novelRepository;
  late final NovelProvider novelProvider;
  late final AuthProvider authProvider;
  late final LibraryProvider libraryProvider;

  void init() {
    secureStorage = SecureStorage();
    authInterceptor = AuthInterceptor(secureStorage);
    dioClient = DioClient(authInterceptor);
    notificationService = NotificationService(dioClient, secureStorage);

    // Auth
    authRemoteDataSource = AuthRemoteDataSourceImpl(dioClient);
    authRepository = AuthRepositoryImpl(authRemoteDataSource, secureStorage);
    authProvider = AuthProvider(authRepository, secureStorage);

    // Novels
    novelRemoteDataSource = NovelRemoteDataSourceImpl(dioClient);
    novelRepository = NovelRepositoryImpl(novelRemoteDataSource);
    novelProvider = NovelProvider(novelRepository);

    // Library
    libraryProvider = LibraryProvider(dioClient);
  }
}

final sl = ServiceLocator();
