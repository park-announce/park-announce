import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/firebase_options.dart';
import 'package:pa_mobile_app/login_page.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(
    options: DefaultFirebaseOptions.currentPlatform,
  );
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
          textTheme: const TextTheme(bodySmall: TextStyle(fontSize: 10)),
          //colorSchemeSeed: const Color(0xFF8dea88),
        ),
        home: const SafeArea(
          child: Scaffold(
            body: LoginPage(),
          ),
        ));
  }
}
