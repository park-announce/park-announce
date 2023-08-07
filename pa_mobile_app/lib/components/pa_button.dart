import 'package:flutter/material.dart';

typedef OnPress = void Function();

class PaButton extends StatelessWidget {
  final OnPress onPressedFunction;
  final String text;
  const PaButton({
    super.key,
    required this.onPressedFunction,
    required this.text,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: () {
        onPressedFunction();
      },
      child: Text(text, style: Theme.of(context).textTheme.labelLarge),
    );
  }
}
