import 'package:email_validator/email_validator.dart';
import 'package:flutter/material.dart';
import 'package:pa_mobile_app/components/pa_button.dart';
import 'package:pa_mobile_app/components/pa_pin_input.dart';
import 'package:pa_mobile_app/components/pa_text_field.dart';
import 'package:pa_mobile_app/pages/register_page.dart';
import 'package:pa_mobile_app/utils/navigation_utils.dart' as nav_utils;

class RegisterMailPage extends StatefulWidget {
  const RegisterMailPage({super.key});

  @override
  State<RegisterMailPage> createState() => _RegisterMailPageState();
}

class _RegisterMailPageState extends State<RegisterMailPage> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _pinController = TextEditingController();

  final int kPinLength = 6;
  bool focusPin = false;
  PageStateStatus _pageStateStatus = PageStateStatus.initial;
  late String email;
  late String pin;

  @override
  void initState() {
    _emailController.addListener(() {
      setState(() {
        email = _emailController.text;
      });
    });
    _pinController.addListener(() {
      setState(() {
        pin = _pinController.text;
      });
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      appBar: AppBar(),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.start,
              children: [
                PaTextField(
                  hintText: 'E Mail',
                  enabled: _pageStateStatus == PageStateStatus.initial,
                  controller: _emailController,
                  keyboardType: TextInputType.emailAddress,
                ),
                ConstrainedBox(
                  constraints: BoxConstraints.tightForFinite(height: _pageStateStatus == PageStateStatus.initial ? 0 : double.infinity),
                  child: Column(
                    children: [
                      const SizedBox(height: 30),
                      PaPinInput(
                        _pageStateStatus == PageStateStatus.otpSent,
                        _pinController,
                        kPinLength,
                        keyboardType: TextInputType.number,
                        changed: (String value) {
                          _pinController.text = value.toString();
                          setState(() {});
                        },
                        completed: (String value) {
                          _checkOtp(value.toString(), context).then((result) {
                            if (result) {
                              nav_utils.navigate(context, RegisterPage(email: _emailController.text, firstName: '', lastName: ''), onReturn: () {
                                _pinController.text = "";
                                _pageStateStatus = PageStateStatus.initial;
                              });
                            }
                          });
                        },
                        requestFocus: focusPin,
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 30),
                Row(
                  children: [Expanded(child: _getButton(context))],
                )
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _getButton(BuildContext context) {
    bool isEnabled = false;
    String buttonText = '';
    if (_pageStateStatus == PageStateStatus.initial) {
      isEnabled = _emailController.value.text.isNotEmpty && EmailValidator.validate(_emailController.value.text);
      buttonText = 'Send Otp';
    } else if (_pageStateStatus == PageStateStatus.otpSent) {
      isEnabled = _pinController.value.text.isNotEmpty && _pinController.value.text.length == kPinLength;
      buttonText = 'Check Otp';
    }

    final Color decorationColor = isEnabled ? Theme.of(context).primaryColor : Colors.grey;

    VoidCallback onPressed = () {};
    if (_pageStateStatus == PageStateStatus.initial) {
      onPressed = () async {
        final bool sendOtpResult = await _sendOtp();
        if (sendOtpResult) {
          setState(() {
            focusPin = true;
            _pageStateStatus = PageStateStatus.otpSent;
          });
        }
      };
    } else {
      onPressed = () {
        _checkOtp(_pinController.text, context);
      };
    }
    return PaButton(
      text: buttonText,
      onPressedFunction: () {
        onPressed();
      },
    );
  }
}

enum PageStateStatus { initial, otpSent }

Future<bool> _sendOtp() async {
  return true;
}

Future<bool> _checkOtp(String value, BuildContext context) async {
  if (value != '123456') {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Invalid Otp', style: Theme.of(context).textTheme.bodyMedium)));
    return false;
  }
  print(value);
  return true;
}
