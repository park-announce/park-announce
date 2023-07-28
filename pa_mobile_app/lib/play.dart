import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

const String kBaseUri = 'http://192.168.0.17:8000';

void checkApi() {
  http.post(Uri.parse('$kBaseUri/google/oauth2/token')).then(
    (value) {
      debugPrint(value.statusCode.toString());
      debugPrint(value.body);
    },
  );
}
