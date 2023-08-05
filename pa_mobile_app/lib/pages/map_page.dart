import 'dart:async';
import 'dart:convert' as convert;

import 'package:flutter/material.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/intective_flag.dart';
import 'package:pa_mobile_app/external/map_controller.dart';
import 'package:pa_mobile_app/external/map_events.dart';
import 'package:pa_mobile_app/external/marker_layer.dart';
import 'package:pa_mobile_app/external/options.dart';
import 'package:pa_mobile_app/external/tile_layer.dart';
import 'package:pa_mobile_app/external/widget.dart';
import 'package:pa_mobile_app/models/socket_request_models.dart';
import 'package:pa_mobile_app/models/socket_response_models.dart';
import 'package:pa_mobile_app/utils/location_utils.dart' as location_utils;
import 'package:pa_mobile_app/utils/map_utils.dart' as map_utils;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

typedef ConvertSocketResponseFunction = dynamic Function(Map<String, dynamic> json);

class MapPage extends StatefulWidget {
  const MapPage({
    Key? key,
  }) : super(key: key);

  @override
  State createState() {
    return _MapPageState();
  }
}

class _MapPageState extends State<MapPage> with TickerProviderStateMixin {
  Timer? countdownTimer;
  late WebSocketChannel _channel;

  final MapController mapController = MapController();
  // Enable pinchZoom and doubleTapZoomBy by default
  int flags = InteractiveFlag.pinchZoom | InteractiveFlag.doubleTapZoom | InteractiveFlag.drag;
  LatLng userLocation = const LatLng(41.09322441408693, 28.998664108745555);
  LatLng mapCenterPosition = const LatLng(41.09322441408693, 28.998664108745555);

  late List<Location> locations = [];
  int remainingTime = 0;
  Map<String, ConvertSocketResponseFunction> socketResponseProcessors = <String, ConvertSocketResponseFunction>{};

  @override
  void initState() {
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) {
      SharedPreferences.getInstance().then((value) {
        _channel = WebSocketChannel.connect(Uri.parse('ws://192.168.0.17:8001/socket/connect?Authorization=${value.getString('Token')}'));
        _channel.stream.listen((event) {
          _processWebSocketMessage(event);
        });

        location_utils.getLocation().then((location) => {_setLocation(location)});
        final SocketRequestMessage<GetLocationsNearbyRequest> message = SocketRequestMessage<GetLocationsNearbyRequest>(
            kGetLocationsNearby, GetLocationsNearbyRequest(userLocation.longitude, userLocation.latitude, 5000, 10));

        _sendSocketMessage(message);
      });
    });
    super.initState();
    location_utils.getLocation().then((value) => {_setLocation(value)});
    countdownTimer = Timer.periodic(const Duration(seconds: 1), (_) {
      setState(() {
        remainingTime = remainingTime - 1;
        locations = remainingTime > 0 ? locations : [];
      });
    });
    socketResponseProcessors[kCreateParkLocation] = (Map<String, dynamic> json) {
      final ReserveParkLocationResponse _ = ReserveParkLocationResponse.fromJson(json);
    };

    socketResponseProcessors[kGetLocationsNearby] = (Map<String, dynamic> json) {
      final NearestLocationsResponse response = NearestLocationsResponse.fromJson(json);
      int index = 1;
      for (final element in response.locations) {
        element.index = index++;
      }
      if (response.locations.isEmpty) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
          content: Center(child: Text('No Data')),
        ));
      }
      setState(() {
        locations = response.locations;
        remainingTime = response.duration;
      });
    };

    socketResponseProcessors[kReserveParkLocation] = (Map<String, dynamic> json) {
      final CreateParkLocationResponse _ = CreateParkLocationResponse.fromJson(json);
    };
  }

  void _setLocation(LatLng location) {
    setState(() {
      userLocation = location;
    });
    map_utils.animatedMapMove(location, 15, mapController, this);
  }

  void _processWebSocketMessage(dynamic event) {
    final Map<String, dynamic> jsonData = convert.jsonDecode(event as String) as Map<String, dynamic>;
    print('Incoming socket message: $jsonData');
    final callback = socketResponseProcessors[jsonData['operation']] as ConvertSocketResponseFunction;
    callback(jsonData['data'] as Map<String, dynamic>);
  }

  void _sendSocketMessage<T>(SocketRequestMessage<T> message) {
    final String messageJson = convert.jsonEncode(message);
    print('Outgoing socket message: $messageJson');

    _channel.sink.add(messageJson);
  }

  void onMapEvent(MapEvent mapEvent) {
    if (mapEvent is MapEventTap) {
      final SocketRequestMessage<CreateParkLocationRequest> message =
          SocketRequestMessage(kCreateParkLocation, CreateParkLocationRequest(mapEvent.tapPosition.longitude, mapEvent.tapPosition.latitude, 30));
      _sendSocketMessage(message);
    }
    if (mapEvent is MapEventMoveEnd) {
      setState(() {
        mapCenterPosition = mapEvent.camera.center;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF132555),
      body: FutureBuilder<UserInfo>(
          future: _getDisplayName(),
          builder: (context, snapshot) {
            if (snapshot.hasData && snapshot.data != null) {
              return SafeArea(
                child: Center(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      const Text('Welcome,', style: TextStyle(color: Colors.white)),
                      Text(
                        snapshot.data!.userName,
                        style: const TextStyle(color: Colors.white),
                      ),
                      Text(snapshot.data!.eMail, style: const TextStyle(color: Colors.white)),
                      const SizedBox(height: 20),
                      Expanded(
                        child: Stack(
                          fit: StackFit.expand,
                          children: [
                            FlutterMap(
                              mapController: mapController,
                              options: MapOptions(
                                onMapEvent: onMapEvent,
                                initialCenter: const LatLng(41, 29),
                                initialZoom: 15,
                                interactionOptions: InteractionOptions(
                                  flags: flags,
                                ),
                              ),
                              children: [
                                TileLayer(
                                  urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
                                  userAgentPackageName: 'dev.fleaflet.flutter_map.example',
                                ),
                                MarkerLayer(
                                  markers: _getMarkers(),
                                )
                              ],
                            ),
                            Positioned(
                              bottom: 30,
                              child: SizedBox(
                                width: MediaQuery.of(context).size.width,
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    CircleAvatar(
                                      backgroundColor: Colors.greenAccent.shade700,
                                      child: const Icon(
                                        Icons.add,
                                        color: Colors.white,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                            Positioned(
                              bottom: 30,
                              left: 15,
                              right: 15,
                              child: Column(children: locations.isNotEmpty ? _getLocations() : _getSearchButton()),
                            )
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              );
            } else {
              return const CircularProgressIndicator();
            }
          }),
    );
  }

  List<Widget> _getSearchButton() {
    return [
      ElevatedButton(
        onPressed: () {
          final SocketRequestMessage<GetLocationsNearbyRequest> message = SocketRequestMessage<GetLocationsNearbyRequest>(
              kGetLocationsNearby, GetLocationsNearbyRequest(mapCenterPosition.longitude, mapCenterPosition.latitude, 5000, 10));
          _sendSocketMessage(message);
        },
        style: TextButton.styleFrom(backgroundColor: Colors.green.shade400, foregroundColor: Colors.white),
        child: const Text('Bu bölgede ara'),
      )
    ];
  }

  List<Widget> _getLocations() {
    final List<Widget> locationWidgets = locations.map((entity) {
      return Container(
          padding: const EdgeInsets.symmetric(vertical: 3, horizontal: 15),
          decoration: const BoxDecoration(color: Colors.grey),
          width: MediaQuery.of(context).size.width,
          child: GestureDetector(
            onTap: () {
              //_openMap(e.latitude!, e.longitude!);
            },
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                CircleAvatar(
                  backgroundColor: Colors.green,
                  child: Text(
                    entity.index.toString(),
                    style: const TextStyle(color: Colors.white),
                  ),
                ),
                Text('${entity.distanceTo!.toStringAsFixed(2)}m'),
                SizedBox(
                  width: 75,
                  child: MaterialButton(
                    onPressed: () {
                      final SocketRequestMessage<ReserveParkRequest> message =
                          SocketRequestMessage<ReserveParkRequest>(kReserveParkLocation, ReserveParkRequest(entity.id!));
                      _sendSocketMessage<ReserveParkRequest>(message);
                      print(entity.latitude);
                      print(entity.longitude);
                    },
                    color: Colors.green,
                    child: const Text('Accept', style: TextStyle(fontSize: 10)),
                  ),
                ),
                SizedBox(
                  width: 75,
                  child: MaterialButton(
                    onPressed: () {
                      locations = locations.where((a) => a.id != entity.id).toList();
                      setState(() {});
                    },
                    color: Colors.red,
                    textColor: Colors.white,
                    child: const Text('Reject', style: TextStyle(fontSize: 10)),
                  ),
                ),
              ],
            ),
          ));
    }).toList();
    locationWidgets.insert(
        0,
        Container(
          padding: const EdgeInsets.symmetric(vertical: 3, horizontal: 15),
          decoration: const BoxDecoration(color: Colors.grey),
          width: MediaQuery.of(context).size.width,
          child: Center(
            child: CircleAvatar(
              backgroundColor: Colors.green,
              child: Text(
                remainingTime.toString(),
                style: const TextStyle(color: Colors.white),
              ),
            ),
          ),
        ));

    return locationWidgets;
  }

  Future<void> _openMap(double latitude, double longitude) async {
    final mapSchema = 'comgooglemaps://?q=@$latitude,$longitude';
    final uri = Uri.parse(mapSchema);
    debugPrint(mapSchema);
    if (await canLaunchUrl(uri)) {
      await launchUrl(uri);
    } else {
      throw 'Could not launch $mapSchema';
    }
  }

  List<Marker> _getMarkers() {
    return locations.map((e) {
      return Marker(
        width: 30,
        height: 30,
        point: LatLng(e.latitude!, e.longitude!),
        builder: (ctx) => CircleAvatar(
          backgroundColor: Colors.green,
          child: Text(
            e.index.toString(),
            style: const TextStyle(color: Colors.white),
          ),
        ),
      );
    }).toList();
  }

  Future<UserInfo> _getDisplayName() async {
    final SharedPreferences pref = await SharedPreferences.getInstance();
    return UserInfo(pref.getString('Name')!, pref.getString('IdToken')!, pref.getString('Email')!);
  }
}

class UserInfo {
  final String userName;
  final String eMail;
  final String idToken;

  UserInfo(this.userName, this.idToken, this.eMail);
}
