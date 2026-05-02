import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:mobile/app/theme/app_theme.dart';
import 'package:mobile/features/novels/presentation/bloc/novel_provider.dart';

class ChapterReaderPage extends StatefulWidget {
  final String slug;
  final int chapterNumber;
  final String novelTitle;

  const ChapterReaderPage({
    super.key,
    required this.slug,
    required this.chapterNumber,
    this.novelTitle = '',
  });

  @override
  State<ChapterReaderPage> createState() => _ChapterReaderPageState();
}

class _ChapterReaderPageState extends State<ChapterReaderPage> {
  String _theme = 'dark';
  double _fontSize = 18.0;
  late int _currentChapterNumber;

  @override
  void initState() {
    super.initState();
    _currentChapterNumber = widget.chapterNumber;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadChapter();
    });
  }

  void _loadChapter() {
    context.read<NovelProvider>().loadChapter(widget.slug, _currentChapterNumber);
  }

  void _goToChapter(int number) {
    setState(() => _currentChapterNumber = number);
    context.read<NovelProvider>().loadChapter(widget.slug, number);
  }

  @override
  Widget build(BuildContext context) {
    Color bg;
    Color fg;
    switch (_theme) {
      case 'sepia':
        bg = const Color(0xFFF4ECD8);
        fg = const Color(0xFF5B4636);
        break;
      case 'light':
        bg = Colors.white;
        fg = const Color(0xFF222222);
        break;
      default:
        bg = AppTheme.background;
        fg = AppTheme.foreground;
    }

    return Consumer<NovelProvider>(
      builder: (context, provider, _) {
        final chapter = provider.currentChapter;
        final chapters = provider.chapters;
        final isLoading = provider.isLoadingChapter;
        final error = provider.chapterError;

        final hasPrev = _currentChapterNumber > 1;
        final hasNext = chapters.isNotEmpty && _currentChapterNumber < chapters.length;

        return Scaffold(
          backgroundColor: bg,
          appBar: AppBar(
            backgroundColor: bg.withOpacity(0.9),
            foregroundColor: fg,
            elevation: 0,
            title: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  chapter != null ? 'Chapter $_currentChapterNumber' : 'Loading...',
                  style: TextStyle(fontSize: 15, color: fg, fontWeight: FontWeight.w600),
                ),
                Text(
                  widget.novelTitle,
                  style: TextStyle(fontSize: 11, color: fg.withOpacity(0.5)),
                ),
              ],
            ),
            actions: [
              IconButton(
                onPressed: () => _showSettings(context),
                icon: Icon(Icons.tune_rounded, color: fg.withOpacity(0.7)),
              ),
            ],
          ),
          body: isLoading
              ? Center(
                  child: CircularProgressIndicator(color: AppTheme.primary),
                )
              : error != null
                  ? Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(Icons.error_outline, size: 48, color: fg.withOpacity(0.4)),
                          const SizedBox(height: 12),
                          Text('Failed to load chapter', style: TextStyle(color: fg.withOpacity(0.5))),
                          const SizedBox(height: 12),
                          ElevatedButton(
                            onPressed: _loadChapter,
                            child: const Text('Retry'),
                          ),
                        ],
                      ),
                    )
                  : SingleChildScrollView(
                      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
                      child: Column(
                        children: [
                          const SizedBox(height: 12),
                          // Chapter title
                          Text(
                            chapter?.title ?? '',
                            textAlign: TextAlign.center,
                            style: GoogleFonts.sora(
                              fontSize: 24,
                              fontWeight: FontWeight.bold,
                              color: fg,
                            ),
                          ),
                          const SizedBox(height: 8),
                          if (chapter != null && chapter.wordCount > 0)
                            Text(
                              '${chapter.wordCount} words',
                              style: TextStyle(color: fg.withOpacity(0.4), fontSize: 12),
                            ),
                          const SizedBox(height: 32),
                          // Chapter content
                          Text(
                            chapter?.content ?? 'No content available.',
                            style: GoogleFonts.sourceSerif4(
                              fontSize: _fontSize,
                              height: 1.85,
                              color: fg,
                            ),
                          ),
                          const SizedBox(height: 48),
                          // Navigation
                          Container(
                            padding: const EdgeInsets.all(16),
                            decoration: BoxDecoration(
                              color: (fg == AppTheme.foreground
                                  ? AppTheme.surface
                                  : fg.withOpacity(0.05)),
                              borderRadius: BorderRadius.circular(16),
                            ),
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                TextButton.icon(
                                  onPressed: hasPrev ? () => _goToChapter(_currentChapterNumber - 1) : null,
                                  icon: Icon(Icons.arrow_back_ios, size: 16, color: hasPrev ? AppTheme.primary : fg.withOpacity(0.2)),
                                  label: Text(
                                    'Previous',
                                    style: TextStyle(color: hasPrev ? AppTheme.primary : fg.withOpacity(0.2)),
                                  ),
                                ),
                                Container(
                                  padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
                                  decoration: BoxDecoration(
                                    color: AppTheme.primary.withOpacity(0.15),
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                  child: Text(
                                    '$_currentChapterNumber / ${chapters.length}',
                                    style: const TextStyle(
                                      color: AppTheme.primary,
                                      fontWeight: FontWeight.bold,
                                      fontSize: 13,
                                    ),
                                  ),
                                ),
                                TextButton.icon(
                                  onPressed: hasNext ? () => _goToChapter(_currentChapterNumber + 1) : null,
                                  icon: Text(
                                    'Next',
                                    style: TextStyle(color: hasNext ? AppTheme.primary : fg.withOpacity(0.2)),
                                  ),
                                  label: Icon(Icons.arrow_forward_ios, size: 16, color: hasNext ? AppTheme.primary : fg.withOpacity(0.2)),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(height: 40),
                        ],
                      ),
                    ),
        );
      },
    );
  }

  void _showSettings(BuildContext context) {
    showModalBottomSheet(
      context: context,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setModalState) {
            return Padding(
              padding: const EdgeInsets.all(24.0),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Center(
                    child: Container(
                      width: 40,
                      height: 4,
                      decoration: BoxDecoration(
                        color: AppTheme.muted.withOpacity(0.3),
                        borderRadius: BorderRadius.circular(2),
                      ),
                    ),
                  ),
                  const SizedBox(height: 20),
                  const Text('Reader Settings', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 24),
                  const Text('Theme', style: TextStyle(color: AppTheme.muted, fontSize: 13)),
                  const SizedBox(height: 12),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _themeButton(setModalState, 'dark', const Color(0xFF0F111A), 'Dark'),
                      _themeButton(setModalState, 'sepia', const Color(0xFFF4ECD8), 'Sepia'),
                      _themeButton(setModalState, 'light', Colors.white, 'Light'),
                    ],
                  ),
                  const SizedBox(height: 24),
                  const Text('Font Size', style: TextStyle(color: AppTheme.muted, fontSize: 13)),
                  const SizedBox(height: 8),
                  Row(
                    children: [
                      const Text('A', style: TextStyle(fontSize: 14)),
                      Expanded(
                        child: Slider(
                          value: _fontSize,
                          min: 14,
                          max: 32,
                          activeColor: AppTheme.primary,
                          inactiveColor: AppTheme.border,
                          onChanged: (val) {
                            setState(() => _fontSize = val);
                            setModalState(() {});
                          },
                        ),
                      ),
                      Text('A', style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
                      const SizedBox(width: 8),
                      Text('${_fontSize.toInt()}', style: const TextStyle(color: AppTheme.muted)),
                    ],
                  ),
                  const SizedBox(height: 16),
                ],
              ),
            );
          },
        );
      },
    );
  }

  Widget _themeButton(StateSetter setModalState, String t, Color color, String label) {
    bool isSelected = _theme == t;
    return GestureDetector(
      onTap: () {
        setState(() => _theme = t);
        setModalState(() {});
      },
      child: Column(
        children: [
          AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            width: 56,
            height: 56,
            decoration: BoxDecoration(
              color: color,
              shape: BoxShape.circle,
              border: Border.all(
                color: isSelected ? AppTheme.primary : AppTheme.border,
                width: isSelected ? 3 : 1,
              ),
              boxShadow: isSelected
                  ? [BoxShadow(color: AppTheme.primary.withOpacity(0.3), blurRadius: 8)]
                  : [],
            ),
            child: isSelected ? const Icon(Icons.check, color: AppTheme.primary) : null,
          ),
          const SizedBox(height: 6),
          Text(label, style: TextStyle(
            fontSize: 12,
            color: isSelected ? AppTheme.primary : AppTheme.muted,
            fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
          )),
        ],
      ),
    );
  }
}
