// To parse this JSON data, do
//
//     final checkApiTokenResponse = checkApiTokenResponseFromMap(jsonString);

class CheckApiTokenResponse {
  String token;

  CheckApiTokenResponse({required this.token});

  factory CheckApiTokenResponse.fromJson(Map<String, dynamic> json) => CheckApiTokenResponse(
        token: json["token"] as String,
      );
}
