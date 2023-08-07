import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

typedef OnPress = void Function();

class PaLoginButton extends StatelessWidget {
  final String text;
  final OnPress onPressedFunction;
  final Widget? icon;
  const PaLoginButton({super.key, required this.onPressedFunction, required this.text, this.icon = null});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: ElevatedButton(
        onPressed: () {
          onPressedFunction();
        },
        child: Container(
          padding: EdgeInsets.symmetric(horizontal: 10),
          child: Row(
            children: [
              Expanded(flex: 1, child: icon != null ? icon! : Container()),
              Expanded(
                flex: 8,
                child: Center(child: Text(text, style: Theme.of(context).textTheme.labelMedium)),
              ),
              Expanded(flex: 1, child: Container()),
            ],
          ),
        ),
      ),
    );
  }
}
