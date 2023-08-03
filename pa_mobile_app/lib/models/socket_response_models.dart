class SocketResponseMessage<T> {
  final String operation;
  final String transactionId;
  final T data;

  SocketResponseMessage(
    this.operation,
    this.transactionId,
    this.data,
  );

  factory SocketResponseMessage.fromJson(Map<String, dynamic> json, T dta) {
    return SocketResponseMessage(json["operation"] as String, json['transction_id'] as String, dta);
  }
}

class CreateParkLocationResponse {
  final int count;

  CreateParkLocationResponse(this.count);
  factory CreateParkLocationResponse.fromJson(Map<String, dynamic> json) {
    return CreateParkLocationResponse(json["count"] as int);
  }
}

class NearestLocationsResponse {
  late int duration;
  final List<Location> locations;
  NearestLocationsResponse(this.duration, this.locations);
  factory NearestLocationsResponse.fromJson(Map<String, dynamic> json) {
    return NearestLocationsResponse(
      json['duration'] as int,
      List<Location>.from(
        (json["locations"] as List).map(
          (x) {
            return Location.fromJson(x as Map<String, dynamic>);
          },
        ),
      ),
    );
  }
}

class Location {
  final double? distanceTo;
  final String? id;
  final double? latitude;
  final double? longitude;
  late final int index;
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

class ReserveParkLocationResponse {
  final int count;

  ReserveParkLocationResponse(this.count);
  factory ReserveParkLocationResponse.fromJson(Map<String, dynamic> json) {
    return ReserveParkLocationResponse(json["count"] as int);
  }
}
