// To parse this JSON data, do
//
//     final checkApiTokenResponse = checkApiTokenResponseFromMap(jsonString);

import 'dart:convert';

String apiErrorResponseToMap(ApiErrorResponse data) => json.encode(data.toMap());

class ApiErrorResponse {
  bool issuccess;
  String message;
  String code;
  String stack;

  ApiErrorResponse({required this.issuccess, required this.message, required this.code, required this.stack});

  factory ApiErrorResponse.fromJson(Map<String, dynamic> json) =>
      ApiErrorResponse(issuccess: json["issuccess"] as bool, message: json["message"] as String, code: json["code"] as String, stack: json["stack"] as String);

  Map<String, dynamic> toMap() => {
        "issuccess": issuccess,
        "message": message,
        "code": code,
        "stack": stack,
      };
}
