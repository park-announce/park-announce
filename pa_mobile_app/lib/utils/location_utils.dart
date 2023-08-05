import 'package:geolocator/geolocator.dart';
import 'package:latlong2/latlong.dart';

Future<LatLng> getLocation() async {
  final bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
  if (serviceEnabled) {
    LocationPermission permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
      if (permission == LocationPermission.always || permission == LocationPermission.whileInUse) {
        final Position position = await Geolocator.getCurrentPosition(desiredAccuracy: LocationAccuracy.high);
        return LatLng(position.latitude, position.longitude);
      }
    } else {
      final Position position = await Geolocator.getCurrentPosition(desiredAccuracy: LocationAccuracy.high);
      return LatLng(position.latitude, position.longitude);
    }
  }
  return const LatLng(41, 29);
}
