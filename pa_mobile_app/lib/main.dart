import 'package:flutter/material.dart';
import 'package:pa_mobile_app/external/interactive_test_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'flutter_map Demo',
        theme: ThemeData(
          useMaterial3: true,
          colorSchemeSeed: const Color(0xFF8dea88),
        ),
        home: const InteractiveTestPage());
  }
}
