import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/firebase_options.dart';
import 'package:pa_mobile_app/pages/on_boarding_page.dart';
import 'package:pa_mobile_app/pages/register_mail_page.dart';
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
    return MaterialApp(
      title: 'flutter_map Demo',
      theme: ThemeData(
        scaffoldBackgroundColor: Colors.black,
        appBarTheme: const AppBarTheme(backgroundColor: Colors.black, foregroundColor: Colors.white),
        hintColor: Colors.grey,
        primaryColor: Colors.white,
        useMaterial3: true,
        disabledColor: Colors.grey,
        backgroundColor: Colors.white,
        buttonTheme: const ButtonThemeData(buttonColor: Colors.black),
        textTheme: TextTheme(
          bodySmall: GoogleFonts.poppins(fontSize: 12, color: Colors.black),
          bodyMedium: GoogleFonts.poppins(fontSize: 15, color: Colors.black, decorationColor: Colors.white),
        ),
        //colorSchemeSeed: const Color(0xFF8dea88),
      ),
      home: const OnBoardingPage(),
    );
  }
}
