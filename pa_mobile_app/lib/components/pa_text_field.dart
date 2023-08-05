import 'package:flutter/material.dart';

typedef TextValueChanged<T> = void Function(T value);

class PaTextField extends StatefulWidget {
  const PaTextField(
      {super.key,
      required this.enabled,
      required this.controller,
      required this.keyboardType,
      this.isPassword = false,
      required this.hintText,
      this.focus = false});
  final bool enabled;
  final TextEditingController controller;
  final TextInputType keyboardType;
  final bool isPassword;
  final String hintText;
  final bool focus;
  @override
  State<PaTextField> createState() => _PaTextFieldState();
}

class _PaTextFieldState extends State<PaTextField> {
  final FocusNode focusNode = FocusNode();

  @override
  void initState() {
    super.initState();
    if (widget.focus) {
      focusNode.requestFocus();
    }
  }

  bool _hidePassword = true;

  @override
  Widget build(BuildContext context) {
    return TextField(
      focusNode: focusNode,
      style: Theme.of(context).textTheme.bodyMedium,
      keyboardType: widget.keyboardType,
      enabled: widget.enabled,
      controller: widget.controller,
      obscureText: widget.isPassword && _hidePassword,
      autocorrect: false,
      decoration: InputDecoration(
        suffixIcon: widget.isPassword
            ? IconButton(
                icon: Icon(
                  _hidePassword ? Icons.visibility_off : Icons.visibility,
                  color: Colors.grey,
                ),
                onPressed: () {
                  setState(() {
                    _hidePassword = !_hidePassword;
                  });
                },
              )
            : null,
        fillColor: Colors.yellow,
        focusColor: Colors.green,
        focusedBorder: OutlineInputBorder(
          borderRadius: const BorderRadius.all(Radius.circular(30)),
          borderSide: BorderSide(width: 1, color: Theme.of(context).primaryColor),
        ),
        contentPadding: const EdgeInsets.fromLTRB(20, 15, 20, 15),
        hintText: widget.hintText,
        hintStyle: Theme.of(context).textTheme.bodyMedium!.copyWith(color: Theme.of(context).hintColor),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(30),
          borderSide: const BorderSide(width: 1, color: Colors.red),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(30),
          borderSide: BorderSide(width: 0.5, color: Colors.grey.shade500),
        ),
      ),
    );
  }
}
