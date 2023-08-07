import 'package:flutter/material.dart';
import 'package:pa_mobile_app/components/pa_button.dart';
import 'package:pa_mobile_app/components/pa_text_field.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key, required this.email, required this.firstName, required this.lastName});
  final String email;
  final String firstName;
  final String lastName;
  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _firstNameController = TextEditingController();
  final TextEditingController _lastNameController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _passwordAgainController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _emailController.text = widget.email;
    _firstNameController.text = widget.firstName;
    _lastNameController.text = widget.lastName;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: Padding(
        padding: const EdgeInsets.all(18),
        child: Column(
          children: [
            PaTextField(
              enabled: false,
              hintText: 'E Mail',
              controller: _emailController,
              keyboardType: TextInputType.text,
              isPassword: false,
            ),
            const SizedBox(height: 20),
            PaTextField(
              focus: true,
              enabled: true,
              hintText: 'First Name',
              controller: _firstNameController,
              keyboardType: TextInputType.text,
              isPassword: false,
            ),
            const SizedBox(height: 20),
            PaTextField(
              enabled: true,
              hintText: 'Last Name',
              controller: _lastNameController,
              keyboardType: TextInputType.text,
              isPassword: false,
            ),
            const SizedBox(height: 20),
            PaTextField(
              enabled: true,
              hintText: 'Password',
              controller: _passwordController,
              keyboardType: TextInputType.text,
              isPassword: true,
            ),
            const SizedBox(height: 20),
            PaTextField(
              enabled: true,
              hintText: 'Password Again',
              controller: _passwordAgainController,
              keyboardType: TextInputType.text,
              isPassword: true,
            ),
            const SizedBox(height: 20),
            Container(width: double.infinity, child: PaButton(text: 'Register', onPressedFunction: () {})),
          ],
        ),
      ),
    );
  }
}
