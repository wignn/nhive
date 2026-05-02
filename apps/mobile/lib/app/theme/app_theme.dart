import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppTheme {
  // Colors from Web OKLCH approximations
  static const Color background = Color(0xFF0F111A);
  static const Color surface = Color(0xFF181A26);
  static const Color primary = Color(0xFFD41A2E);
  static const Color secondary = Color(0xFF232533);
  static const Color accent = Color(0xFFD41A2E);
  static const Color foreground = Color(0xFFF8F8F8);
  static const Color muted = Color(0xFF949FB1);
  static const Color border = Color(0x33949FB1);

  static const LinearGradient brandGradient = LinearGradient(
    colors: [Color(0xFFD41A2E), Color(0xFFE52E3D)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  static ThemeData get darkTheme {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      scaffoldBackgroundColor: background,
      colorScheme: const ColorScheme.dark(
        primary: primary,
        secondary: secondary,
        surface: surface,
        onPrimary: Colors.white,
        onSurface: foreground,
      ),
      textTheme: GoogleFonts.interTextTheme(
        const TextTheme(
          displayLarge: TextStyle(fontWeight: FontWeight.bold, color: foreground),
          displayMedium: TextStyle(fontWeight: FontWeight.bold, color: foreground),
          displaySmall: TextStyle(fontWeight: FontWeight.bold, color: foreground),
          headlineMedium: TextStyle(fontWeight: FontWeight.bold, color: foreground),
          titleLarge: TextStyle(fontWeight: FontWeight.w600, color: foreground),
          bodyLarge: TextStyle(color: foreground),
          bodyMedium: TextStyle(color: foreground),
        ),
      ).copyWith(
        displayLarge: GoogleFonts.sora(fontWeight: FontWeight.bold, color: foreground),
        displayMedium: GoogleFonts.sora(fontWeight: FontWeight.bold, color: foreground),
        displaySmall: GoogleFonts.sora(fontWeight: FontWeight.bold, color: foreground),
        headlineMedium: GoogleFonts.sora(fontWeight: FontWeight.bold, color: foreground),
      ),
      cardTheme: CardThemeData(
        color: surface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
          side: const BorderSide(color: border),
        ),
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: false,
        titleTextStyle: TextStyle(
          fontSize: 20,
          fontWeight: FontWeight.bold,
          color: foreground,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: primary,
          foregroundColor: Colors.white,
          minimumSize: const Size.fromHeight(50),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          textStyle: const TextStyle(fontWeight: FontWeight.bold),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: secondary,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide.none,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: border),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: primary),
        ),
        hintStyle: const TextStyle(color: muted),
      ),
    );
  }
}
