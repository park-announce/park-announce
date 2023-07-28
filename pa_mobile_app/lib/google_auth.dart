import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';
import 'package:pa_mobile_app/models/api_error_response.dart';
import 'package:pa_mobile_app/models/check_api_token_response.dart';

const String kBaseUri = 'http://192.168.0.17:8000';

void checkApi() {
  http.post(Uri.parse('$kBaseUri/google/oauth2/token')).then(
    (value) {
      debugPrint(value.statusCode.toString());
      debugPrint(value.body);
    },
  );
}

Future<dynamic> checkApiToken(String apiToken) async {
  final Response response = await http.post(
    Uri.parse('$kBaseUri/google/oauth2/token'),
    headers: <String, String>{
      'Content-Type': 'application/json; charset=UTF-8',
    },
    body: jsonEncode(<String, String>{
      'token': apiToken,
    }),
  );
  debugPrint(response.statusCode.toString());
  debugPrint(response.body);

  if (response.statusCode >= 200 && response.statusCode < 300) {
    final tokenResponse = CheckApiTokenResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
    return tokenResponse;
  } else {
    final errorResponse = ApiErrorResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
    return errorResponse;
  }
}
