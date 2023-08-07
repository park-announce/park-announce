import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/*

















*/
final ThemeData old = ThemeData(
  elevatedButtonTheme: const ElevatedButtonThemeData(
    style: ButtonStyle(
      foregroundColor: MaterialStatePropertyAll(Colors.white),
      backgroundColor: MaterialStatePropertyAll(Colors.black),
      padding: MaterialStatePropertyAll(EdgeInsets.all(8)),
    ),
  ),
  scaffoldBackgroundColor: Colors.white,
  appBarTheme: const AppBarTheme(backgroundColor: Colors.white, foregroundColor: Colors.black),
  hintColor: Colors.grey,
  primaryColor: Colors.white,
  useMaterial3: true,
  disabledColor: Colors.grey,
  backgroundColor: Colors.black,
  buttonTheme: const ButtonThemeData(buttonColor: Colors.black),
  textTheme: TextTheme(
    bodySmall: GoogleFonts.poppins(fontSize: 12, color: Colors.black),
    bodyMedium: GoogleFonts.poppins(fontSize: 15, color: Colors.black, decorationColor: Colors.white),
  ),
  colorScheme: const ColorScheme(
      brightness: Brightness.light,
      primary: Colors.white,
      onPrimary: Colors.white,
      secondary: Colors.black,
      onSecondary: Colors.black,
      error: Colors.red,
      onError: Colors.red,
      background: Colors.black,
      onBackground: Colors.black,
      surface: Colors.grey,
      onSurface: Colors.grey),
);
