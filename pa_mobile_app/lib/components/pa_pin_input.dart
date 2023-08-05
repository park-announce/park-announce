import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:pinput/pinput.dart';

typedef PinValueChanged = void Function(String value);

class PaPinInput<T> extends StatelessWidget {
  PaPinInput(this.enabled, this.controller, this.length,
      {super.key, required this.changed, required this.completed, required this.requestFocus, required this.keyboardType});
  final bool enabled;
  final TextEditingController controller;
  final int length;
  final PinValueChanged changed;
  final PinValueChanged completed;
  final bool requestFocus;
  final TextInputType keyboardType;
  final focusNode = FocusNode();

  @override
  Widget build(BuildContext context) {
    const errorColor = Color.fromRGBO(255, 234, 238, 1);
    const fillColor = Color.fromRGBO(222, 231, 240, .57);
    if (requestFocus) {
      focusNode.requestFocus();
    }
    final defaultPinTheme = PinTheme(
      width: 56,
      height: 60,
      textStyle: GoogleFonts.poppins(
        fontSize: 22,
        color: const Color.fromRGBO(30, 60, 87, 1),
      ),
      decoration: BoxDecoration(
        color: fillColor,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.transparent),
      ),
    );

    return SizedBox(
      height: 68,
      child: Pinput(
        keyboardType: keyboardType,
        onChanged: (value) {
          changed(value);
        },
        onCompleted: (value) {
          completed(value);
        },
        enabled: enabled,
        length: length,
        controller: controller,
        focusNode: focusNode,
        defaultPinTheme: defaultPinTheme,
        separatorBuilder: (index) => const SizedBox(width: 8),
        focusedPinTheme: defaultPinTheme.copyWith(
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(8),
            boxShadow: const [
              BoxShadow(
                color: Color.fromRGBO(0, 0, 0, 0.05999999865889549),
                offset: Offset(0, 3),
                blurRadius: 16,
              )
            ],
          ),
        ),
        errorPinTheme: defaultPinTheme.copyWith(
          decoration: BoxDecoration(
            color: errorColor,
            borderRadius: BorderRadius.circular(8),
          ),
        ),
        // onClipboardFound: (value) {
        //   debugPrint('onClipboardFound: $value');
        //   controller.setText(value);
        // },
      ),
    );
  }
}
