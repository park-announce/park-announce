import 'package:flutter/material.dart';
import 'package:pa_mobile_app/components/pa_button.dart';
import 'package:pa_mobile_app/components/pa_login_button.dart';

class ThemeTestPage extends StatefulWidget {
  const ThemeTestPage({super.key});

  @override
  State<ThemeTestPage> createState() => _ThemeTestPageState();
}

class _ThemeTestPageState extends State<ThemeTestPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: Container(
        child: Column(
          children: [
            PaLoginButton(onPressedFunction: () {}, text: 'Pa Login Button'),
            PaButton(onPressedFunction: () {}, text: 'Pa Button'),
            MaterialButton(onPressed: () {}, child: Text('Material Button')),
            ElevatedButton(onPressed: () {}, child: Text('Elevated Button')),
          ],
        ),
      ),
    );
  }
}
