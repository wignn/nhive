import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:mobile/app/theme/app_theme.dart';
import 'package:mobile/features/auth/presentation/bloc/auth_provider.dart';

class ProfilePage extends StatelessWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<AuthProvider>(
      builder: (context, auth, _) {
        if (auth.status == AuthStatus.authenticated && auth.user != null) {
          return _buildAuthenticatedView(context, auth);
        }
        return _buildUnauthenticatedView(context);
      },
    );
  }

  Widget _buildUnauthenticatedView(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 100,
                height: 100,
                decoration: BoxDecoration(
                  color: AppTheme.primary.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(25),
                ),
                child: const Icon(Icons.person_outline, size: 50, color: AppTheme.primary),
              ),
              const SizedBox(height: 24),
              const Text(
                'Welcome to NovelHive',
                style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 8),
              const Text(
                'Sign in to save your reading progress,\nbookmark novels, and more.',
                textAlign: TextAlign.center,
                style: TextStyle(color: AppTheme.muted, fontSize: 14, height: 1.5),
              ),
              const SizedBox(height: 32),
              ElevatedButton(
                onPressed: () => Navigator.pushNamed(context, '/login'),
                child: const Text('Sign In'),
              ),
              const SizedBox(height: 12),
              OutlinedButton(
                onPressed: () => Navigator.pushNamed(context, '/register'),
                style: OutlinedButton.styleFrom(
                  side: const BorderSide(color: AppTheme.primary),
                  minimumSize: const Size.fromHeight(50),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
                child: const Text('Create Account', style: TextStyle(color: AppTheme.primary)),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildAuthenticatedView(BuildContext context, AuthProvider auth) {
    final user = auth.user!;
    return Scaffold(
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: Column(
            children: [
              const SizedBox(height: 24),
              // Avatar
              Container(
                width: 90,
                height: 90,
                decoration: BoxDecoration(
                  gradient: AppTheme.brandGradient,
                  borderRadius: BorderRadius.circular(22),
                  boxShadow: [
                    BoxShadow(
                      color: AppTheme.primary.withOpacity(0.3),
                      blurRadius: 20,
                      offset: const Offset(0, 8),
                    ),
                  ],
                ),
                child: Center(
                  child: Text(
                    user.username.isNotEmpty ? user.username[0].toUpperCase() : '?',
                    style: const TextStyle(
                      fontSize: 36,
                      fontWeight: FontWeight.bold,
                      color: Colors.white,
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 16),
              Text(
                user.username,
                style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 4),
              Text(
                user.email,
                style: const TextStyle(color: AppTheme.muted, fontSize: 14),
              ),
              const SizedBox(height: 8),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                decoration: BoxDecoration(
                  color: AppTheme.primary.withOpacity(0.15),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Text(
                  user.role.toUpperCase(),
                  style: const TextStyle(
                    color: AppTheme.primary,
                    fontSize: 11,
                    fontWeight: FontWeight.bold,
                    letterSpacing: 1,
                  ),
                ),
              ),
              const SizedBox(height: 32),

              // Settings List
              _buildSettingsTile(
                icon: Icons.dark_mode_outlined,
                title: 'Appearance',
                subtitle: 'Dark mode',
                onTap: () {},
              ),
              _buildSettingsTile(
                icon: Icons.notifications_outlined,
                title: 'Notifications',
                subtitle: 'Manage alerts',
                onTap: () {},
              ),
              _buildSettingsTile(
                icon: Icons.download_outlined,
                title: 'Downloads',
                subtitle: 'Offline reading',
                onTap: () {},
              ),
              _buildSettingsTile(
                icon: Icons.info_outline,
                title: 'About',
                subtitle: 'Version 1.0.0',
                onTap: () {},
              ),
              const SizedBox(height: 24),

              // Logout
              SizedBox(
                width: double.infinity,
                child: OutlinedButton.icon(
                  onPressed: () async {
                    final confirm = await showDialog<bool>(
                      context: context,
                      builder: (ctx) => AlertDialog(
                        backgroundColor: AppTheme.surface,
                        title: const Text('Logout'),
                        content: const Text('Are you sure you want to logout?'),
                        actions: [
                          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Cancel')),
                          TextButton(
                            onPressed: () => Navigator.pop(ctx, true),
                            child: const Text('Logout', style: TextStyle(color: AppTheme.primary)),
                          ),
                        ],
                      ),
                    );
                    if (confirm == true && context.mounted) {
                      auth.logout();
                    }
                  },
                  icon: const Icon(Icons.logout, color: AppTheme.primary),
                  label: const Text('Sign Out', style: TextStyle(color: AppTheme.primary)),
                  style: OutlinedButton.styleFrom(
                    side: const BorderSide(color: AppTheme.border),
                    minimumSize: const Size.fromHeight(50),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildSettingsTile({
    required IconData icon,
    required String title,
    required String subtitle,
    required VoidCallback onTap,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      decoration: BoxDecoration(
        color: AppTheme.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppTheme.border),
      ),
      child: ListTile(
        onTap: onTap,
        leading: Container(
          width: 40,
          height: 40,
          decoration: BoxDecoration(
            color: AppTheme.primary.withOpacity(0.1),
            borderRadius: BorderRadius.circular(10),
          ),
          child: Icon(icon, color: AppTheme.primary, size: 20),
        ),
        title: Text(title, style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 15)),
        subtitle: Text(subtitle, style: const TextStyle(color: AppTheme.muted, fontSize: 12)),
        trailing: const Icon(Icons.chevron_right, color: AppTheme.muted),
      ),
    );
  }
}
