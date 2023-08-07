import 'dart:async';

import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/constants/theme_constant.dart';
import 'package:pa_mobile_app/firebase_options.dart';
import 'package:pa_mobile_app/pages/on_boarding_page.dart';
import 'package:pa_mobile_app/theme_test_page.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await Firebase.initializeApp(
    options: DefaultFirebaseOptions.currentPlatform,
  );
  SharedPreferences.getInstance().then((value) {
    value.remove('Email');
    value.remove('IdToken');
    value.remove('Name');
    runApp(const MyApp());
  });
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
    return SafeArea(
      child: MediaQuery(
        data: const MediaQueryData(textScaleFactor: 1),
        child: MaterialApp(
          title: 'flutter_map Demo',
          theme: ThemeData(
            scaffoldBackgroundColor: Colors.white,
            appBarTheme: const AppBarTheme(foregroundColor: Colors.white, backgroundColor: Colors.black),
            primaryColor: Colors.white,
            brightness: Brightness.light,
            useMaterial3: true,
            colorScheme: const ColorScheme(
                brightness: Brightness.light,
                primary: Colors.black,
                onPrimary: Colors.black,
                secondary: Colors.white,
                onSecondary: Colors.white,
                error: Colors.red,
                onError: Colors.red,
                background: Colors.black,
                onBackground: Colors.black,
                surface: Colors.grey,
                onSurface: Colors.grey),
            textTheme: TextTheme(
              labelMedium: GoogleFonts.poppins(color: Colors.white, fontSize: 12),
              labelLarge: GoogleFonts.poppins(color: Colors.white, fontSize: 15),
              bodySmall: GoogleFonts.poppins(color: Colors.black, fontSize: 12),
              bodyMedium: GoogleFonts.poppins(color: Colors.black, fontSize: 15),
            ),
            elevatedButtonTheme: const ElevatedButtonThemeData(
              style: ButtonStyle(
                foregroundColor: MaterialStatePropertyAll(Colors.white),
                backgroundColor: MaterialStatePropertyAll(Colors.black),
                padding: MaterialStatePropertyAll(EdgeInsets.all(8)),
              ),
            ),
          ),
          home: const OnBoardingPage(),
        ),
      ),
    );
  }
}
