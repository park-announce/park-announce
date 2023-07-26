import 'package:flutter/material.dart';

typedef OnPress = void Function();

class PaLoginButton extends StatelessWidget {
  final Color backColor;
  final Color textColor;
  final Widget child;
  final OnPress onPressedFunction;

  const PaLoginButton({super.key, required this.backColor, required this.textColor, required this.onPressedFunction, required this.child});

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 30),
        child: MaterialButton(
            textColor: textColor,
            onPressed: () {
              onPressedFunction();
            },
            child: Container(
                height: 40,
                padding: const EdgeInsets.symmetric(horizontal: 10),
                width: double.infinity,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  border: Border.all(width: 1, color: Colors.grey),
                  color: backColor,
                  borderRadius: BorderRadius.circular(30),
                ),
                child: child)),
      ),
    );
  }
}