class SocketMessage {
  final String operation;
  final String transactionId;
  final SocketData data;

  SocketMessage(this.operation, this.transactionId, this.data);

  Map<String, dynamic> toJson() => {'operation': operation, 'transaction_id': transactionId, 'data': data};
}

class SocketData {
  final double longitude;
  final double latitude;
  final double distance;
  final int count;

  SocketData(this.longitude, this.latitude, this.distance, this.count);
  Map<String, dynamic> toJson() => {'longitude': longitude, 'latitude': latitude, 'distance': distance, 'count': count};
}

class NearestLocationsResponse {
  final String operation;
  final String transactionId;
  final Data data;

  NearestLocationsResponse(
    this.operation,
    this.transactionId,
    this.data,
  );

  factory NearestLocationsResponse.fromJson(Map<String, dynamic> json) =>
      NearestLocationsResponse(json["operation"] as String, json["transaction_id"] as String, Data.fromJson(json["data"] as Map<String, dynamic>));
}

class Data {
  final int duration;
  final List<Location> locations;

  Data(
    this.duration,
    this.locations,
  );

  factory Data.fromJson(Map<String, dynamic> json) => Data(
      json["duration"] as int,
      List<Location>.from((json["locations"] as List).map((x) {
        return Location.fromJson(x as Map<String, dynamic>);
      })));
}

class Location {
  final double? distanceTo;
  final String? id;
  final double? latitude;
  final double? longitude;

  Location({
    this.distanceTo,
    this.id,
    this.latitude,
    this.longitude,
  });

  factory Location.fromJson(Map<String, dynamic> json) {
    return Location(
        distanceTo: json["distance_to"] as double, id: json["id"] as String, latitude: json["latitude"] as double, longitude: json["longitude"] as double);
  }
}
