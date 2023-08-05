import 'package:flutter/material.dart';
import 'package:pa_mobile_app/components/pa_button.dart';
import 'package:pa_mobile_app/components/pa_text_field.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(18),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              PaTextField(enabled: true, controller: _emailController, keyboardType: TextInputType.emailAddress, hintText: 'E Mail'),
              const SizedBox(height: 30),
              PaTextField(enabled: true, controller: _passwordController, keyboardType: TextInputType.text, hintText: 'Password', isPassword: true),
              const SizedBox(height: 30),
              PaButton(text: 'Login', onPressedFunction: () {})
            ],
          ),
        ),
      ),
    );
  }
}
