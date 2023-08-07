import 'dart:async';

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
  late int remainingTime = 0;
  late Timer remainingTimer;
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
  void dispose() {
    remainingTimer.cancel();
    _emailController.dispose();
    _pinController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      appBar: AppBar(
        title: const Text(
          'Sign Up',
        ),
      ),
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
                  focus: true,
                ),
                ConstrainedBox(
                  constraints: BoxConstraints.tightForFinite(height: _pageStateStatus == PageStateStatus.initial ? 0 : double.infinity),
                  child: Column(
                    children: [
                      const SizedBox(height: 30),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.center,
                        mainAxisAlignment: MainAxisAlignment.start,
                        children: [
                          Container(
                            height: 100,
                            width: 100,
                            padding: EdgeInsets.all(10),
                            decoration: BoxDecoration(
                                shape: BoxShape.circle,
                                border: Border.all(
                                  width: 3,
                                  color: Theme.of(context).colorScheme.primary,
                                )),
                            child: Center(child: Text(remainingTime.toString(), style: TextStyle(fontSize: 30))),
                          ),
                          SizedBox(height: 30),
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
    } else if (_pageStateStatus == PageStateStatus.otpTimedOut) {
      isEnabled = _pinController.value.text.isNotEmpty && _pinController.value.text.length == kPinLength;
      buttonText = 'ReSend Otp';
    }

    VoidCallback onPressed = () {};
    if (_pageStateStatus == PageStateStatus.otpTimedOut) {
      onPressed = () async {
        final bool sendOtpResult = await _sendOtp();
        if (sendOtpResult) {
          setState(() {
            focusPin = true;
            remainingTime = 15;
            _pageStateStatus = PageStateStatus.otpSent;
          });
        }
      };
    } else if (_pageStateStatus == PageStateStatus.initial) {
      onPressed = () async {
        final bool sendOtpResult = await _sendOtp();
        if (sendOtpResult) {
          remainingTimer = Timer.periodic(Duration(seconds: 1), (timer) {
            setState(() {
              if (remainingTime > 0) {
                remainingTime = remainingTime - 1;
              } else {
                _pageStateStatus = PageStateStatus.otpTimedOut;
              }
            });
          });

          setState(() {
            focusPin = true;
            remainingTime = 15;
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
}

enum PageStateStatus { initial, otpSent, otpTimedOut }
