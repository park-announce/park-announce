import 'package:flutter/material.dart';
import 'package:pa_mobile_app/components/pa_text_field.dart';
import 'package:pa_mobile_app/pages/map_page.dart';
import 'package:pa_mobile_app/utils/navigation_utils.dart' as nav_utils;

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key, required this.email});
  final String email;
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
            MaterialButton(
              padding: EdgeInsets.zero,
              textColor: Colors.white,
              onPressed: () {
                nav_utils.navigate(context, const MapPage());
              },
              child: Container(
                  height: 40,
                  //padding: const EdgeInsets.symmetric(horizontal: 30),
                  width: double.infinity,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    border: Border.all(width: 1, color: Colors.grey),
                    color: Theme.of(context).primaryColor,
                    borderRadius: BorderRadius.circular(30),
                  ),
                  child: const Text('Register')),
            )
          ],
        ),
      ),
    );
  }
}
