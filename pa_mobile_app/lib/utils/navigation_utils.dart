import 'package:flutter/material.dart';

typedef OnReturn = void Function();

void navigate(BuildContext ctx, Widget navigationPage, {OnReturn? onReturn}) {
  Navigator.of(ctx).push(MaterialPageRoute<dynamic>(builder: (context) => navigationPage));
  if (onReturn != null) {
    onReturn();
  }
}
