import 'package:flutter/material.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/interactive_test_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  LatLng position = const LatLng(41, 29);

  @override
  void initState() {
    super.initState();
  }

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
