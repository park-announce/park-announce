import 'dart:convert';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';
import 'package:pa_mobile_app/models/api_error_response.dart';
import 'package:pa_mobile_app/models/check_api_token_response.dart';
import 'package:uuid/uuid.dart';

const String kBaseUri = 'http://192.168.0.17:8000';

void checkApi() {
  http.post(Uri.parse('$kBaseUri/google/oauth2/token')).then(
    (value) {
      debugPrint(value.statusCode.toString());
      debugPrint(value.body);
    },
  );
}

Future<dynamic> register(String apiToken) async {
  final Response response = await call(
    '$kBaseUri/google/oauth2/register',
    <String, String>{
      'Content-Type': 'application/json; charset=UTF-8',
    },
    jsonEncode(<String, String>{'token': apiToken, 'client_type': Platform.isIOS ? 'ios' : 'android'}),
  );
  if (response.statusCode >= 200 && response.statusCode < 300) {
    final tokenResponse = CheckApiTokenResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
    return tokenResponse;
  } else {
    final errorResponse = ApiErrorResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
    return errorResponse;
  }
}

Future<dynamic> checkApiToken(String apiToken) async {
  final Response response = await call(
    '$kBaseUri/google/oauth2/token',
    <String, String>{
      'Content-Type': 'application/json; charset=UTF-8',
    },
    jsonEncode(<String, String>{'token': apiToken, 'client_type': Platform.isIOS ? 'ios' : 'android'}),
  );
  if (response.statusCode >= 200 && response.statusCode < 300) {
    final tokenResponse = CheckApiTokenResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
    return tokenResponse;
  } else {
    final errorResponse = ApiErrorResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
    return errorResponse;
  }
}

Future<Response> call(String url, Map<String, String> headers, String body) async {
  final String guid = Uuid().v4();
  print('Request for $guid: Url: $url, Body: $body');
  final Response response = await http.post(Uri.parse(url), headers: headers, body: body);
  print('Response for $guid: ${response.statusCode} ${response.body}');
  return response;
}
