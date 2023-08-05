// To parse this JSON data, do
//
//     final checkApiTokenResponse = checkApiTokenResponseFromMap(jsonString);

class PreRegisterGoogleResponse {
  String guid;

  PreRegisterGoogleResponse({required this.guid});

  factory PreRegisterGoogleResponse.fromJson(Map<String, dynamic> json) => PreRegisterGoogleResponse(
        guid: json["guid"] as String,
      );
}
