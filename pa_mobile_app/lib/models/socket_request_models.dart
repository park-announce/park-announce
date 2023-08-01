import 'package:uuid/uuid.dart';

const String kReserveParkLocation = 'reserve_park_location';
const String kCreateParkLocation = 'create_park_location';
const String kGetLocationsNearby = 'get_locations_nearby';

class SocketRequestMessage<T> {
  final String operation;
  final String transactionId = const Uuid().v4();
  final int timeOut = 5;
  final T data;

  SocketRequestMessage(this.operation, this.data);

  Map<String, dynamic> toJson() => {'operation': operation, 'transaction_id': transactionId, 'timeout': 5, 'data': data};
}

class ReserveParkRequest {
  final String id;

  ReserveParkRequest(this.id);
  Map<String, dynamic> toJson() => {'id': id};
}

class CreateParkLocationRequest {
  final double longitude;
  final double latitude;
  final int duration;

  CreateParkLocationRequest(this.longitude, this.latitude, this.duration);
  Map<String, dynamic> toJson() => {'longitude': longitude, 'latitude': latitude, 'duration': duration};
}

class GetLocationNearbyRequest {
  final double longitude;
  final double latitude;
  final double distance;
  final List<int> locationTypes = [0, 1, 2];
  final List<int> vehicleTypes = [0, 1, 2];
  final int count;

  GetLocationNearbyRequest(this.longitude, this.latitude, this.distance, this.count);
  Map<String, dynamic> toJson() =>
      {'longitude': longitude, 'latitude': latitude, 'distance': distance, 'count': count, "location_types": locationTypes, "vehicle_types": vehicleTypes};
}
