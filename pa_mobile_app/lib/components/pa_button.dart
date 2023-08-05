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
    return MaterialButton(
      padding: EdgeInsets.zero,
      onPressed: () {
        onPressedFunction();
      },
      child: Container(
          height: 40,
          //padding: const EdgeInsets.symmetric(horizontal: 30),
          width: double.infinity,
          alignment: Alignment.center,
          decoration: BoxDecoration(
            border: Border.all(width: 1, color: Colors.grey),
            color: Theme.of(context).colorScheme.secondary,
            borderRadius: BorderRadius.circular(30),
          ),
          child: Text(text, style: Theme.of(context).textTheme.bodyMedium!.copyWith(color: Theme.of(context).colorScheme.primary))),
    );
  }
}
