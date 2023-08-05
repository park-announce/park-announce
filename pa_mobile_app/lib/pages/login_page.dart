import 'package:flutter/material.dart';
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
              Row(
                children: [
                  Expanded(
                    child: MaterialButton(
                      padding: EdgeInsets.zero,
                      textColor: Colors.white,
                      onPressed: () {},
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
                          child: Text('Login', style: Theme.of(context).textTheme.bodyMedium)),
                    ),
                  )
                ],
              )
            ],
          ),
        ),
      ),
    );
  }
}
