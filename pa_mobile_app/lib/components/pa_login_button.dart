import 'package:flutter/material.dart';

typedef OnPress = void Function();

class PaLoginButton extends StatelessWidget {
  final String text;
  final OnPress onPressedFunction;

  const PaLoginButton({super.key, required this.onPressedFunction, required this.text});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        const SizedBox(width: 20),
        Expanded(
          child: MaterialButton(
              padding: EdgeInsets.zero,
              onPressed: () {
                onPressedFunction();
              },
              child: Container(
                  height: 40,
                  padding: const EdgeInsets.symmetric(horizontal: 30),
                  width: double.infinity,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    border: Border.all(width: 1, color: Colors.grey),
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(30),
                  ),
                  child: Text(
                    text,
                    style: TextStyle(color: Colors.black),
                  ))),
        ),
        const SizedBox(width: 20),
      ],
    );
  }
}
