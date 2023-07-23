import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:geolocator/geolocator.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/intective_flag.dart';
import 'package:pa_mobile_app/external/map_controller.dart';
import 'package:pa_mobile_app/external/map_events.dart';
import 'package:pa_mobile_app/external/marker_layer.dart';
import 'package:pa_mobile_app/external/options.dart';
import 'package:pa_mobile_app/external/tile_layer.dart';
import 'package:pa_mobile_app/external/widget.dart';
import 'package:url_launcher/url_launcher.dart';

class InteractiveTestPage extends StatefulWidget {
  const InteractiveTestPage({
    Key? key,
  }) : super(key: key);

  @override
  State createState() {
    return _InteractiveTestPageState();
  }
}

class _InteractiveTestPageState extends State<InteractiveTestPage> with TickerProviderStateMixin {
  final MapController mapController = MapController();
  // Enable pinchZoom and doubleTapZoomBy by default
  int flags = InteractiveFlag.pinchZoom | InteractiveFlag.doubleTapZoom | InteractiveFlag.drag;
  LatLng userLocation = LatLng(41, 29);
  static const _startedId = 'AnimatedMapController#MoveStarted';
  static const _inProgressId = 'AnimatedMapController#MoveInProgress';
  static const _finishedId = 'AnimatedMapController#MoveFinished';

  @override
  void initState() {
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) {
      _getLocation().then((value) => {_setLocation(value)});
    });
    super.initState();
  }

  void _setLocation(LatLng location) {
    userLocation = location;
    _animatedMapMove(location, 15);
  }

  void _animatedMapMove(LatLng destLocation, double destZoom) {
    // Create some tweens. These serve to split up the transition from one location to another.
    // In our case, we want to split the transition be<tween> our current map center and the destination.
    final camera = mapController.camera;
    final latTween = Tween<double>(begin: camera.center.latitude, end: destLocation.latitude);
    final lngTween = Tween<double>(begin: camera.center.longitude, end: destLocation.longitude);
    final zoomTween = Tween<double>(begin: camera.zoom, end: destZoom);

    // Create a animation controller that has a duration and a TickerProvider.
    final controller = AnimationController(duration: const Duration(milliseconds: 500), vsync: this);
    // The animation determines what path the animation will take. You can try different Curves values, although I found
    // fastOutSlowIn to be my favorite.
    final Animation<double> animation = CurvedAnimation(parent: controller, curve: Curves.fastOutSlowIn);

    // Note this method of encoding the target destination is a workaround.
    // When proper animated movement is supported (see #1263) we should be able
    // to detect an appropriate animated movement event which contains the
    // target zoom/center.
    final startIdWithTarget = '$_startedId#${destLocation.latitude},${destLocation.longitude},$destZoom';
    bool hasTriggeredMove = false;

    controller.addListener(() {
      final String id;
      if (animation.value == 1.0) {
        id = _finishedId;
      } else if (!hasTriggeredMove) {
        id = startIdWithTarget;
      } else {
        id = _inProgressId;
      }

      hasTriggeredMove |= mapController.move(
        LatLng(latTween.evaluate(animation), lngTween.evaluate(animation)),
        zoomTween.evaluate(animation),
        id: id,
      );
    });

    animation.addStatusListener((status) {
      if (status == AnimationStatus.completed) {
        controller.dispose();
      } else if (status == AnimationStatus.dismissed) {
        controller.dispose();
      }
    });

    controller.forward();
  }

  void onMapEvent(MapEvent mapEvent) {
    if (mapEvent is! MapEventMove && mapEvent is! MapEventRotate) {
      // do not flood console with move and rotate events
      debugPrint(_eventName(mapEvent));
    }

    setState(() {});
  }

  Future<LatLng> _getLocation() async {
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

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.blueGrey,
      ),
      body: Padding(
        padding: const EdgeInsets.all(8),
        child: Column(
          children: [
            Flexible(
              child: FlutterMap(
                mapController: mapController,
                options: MapOptions(
                  onMapEvent: onMapEvent,
                  initialCenter: const LatLng(41, 29),
                  initialZoom: 15,
                  keepAlive: false,
                  interactionOptions: InteractionOptions(
                    flags: flags,
                  ),
                ),
                children: [
                  TileLayer(
                    urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
                    userAgentPackageName: 'dev.fleaflet.flutter_map.example',
                  ),
                  MarkerLayer(markers: [
                    Marker(
                      width: 80,
                      height: 80,
                      point: userLocation,
                      builder: (ctx) => GestureDetector(
                          onTap: () {
                            _openMap(41.08, 29);
                          },
                          child: const FaIcon(FontAwesomeIcons.car)),
                    ),
                    Marker(
                      width: 80,
                      height: 80,
                      point: const LatLng(41.08, 29),
                      builder: (ctx) => GestureDetector(
                          onTap: () {
                            _openMap(41.08, 29);
                          },
                          child: const FaIcon(FontAwesomeIcons.squareParking)),
                    ),
                    Marker(
                      width: 80,
                      height: 80,
                      point: const LatLng(41.083, 29),
                      builder: (ctx) => GestureDetector(
                          onTap: () {
                            _openMap(41.083, 29);
                          },
                          child: const FaIcon(FontAwesomeIcons.squareParking)),
                    ),
                    Marker(
                      width: 80,
                      height: 80,
                      point: const LatLng(41.09, 29.001),
                      builder: (ctx) => GestureDetector(
                          onTap: () {
                            _openMap(41.09, 29.001);
                          },
                          child: const FaIcon(FontAwesomeIcons.squareParking)),
                    ),
                  ])
                ],
              ),
            ),
          ],
        ),
      ),
    );
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

  String _eventName(MapEvent? event) {
    switch (event) {
      case MapEventTap():
        return 'MapEventTap';
      case MapEventSecondaryTap():
        return 'MapEventSecondaryTap';
      case MapEventLongPress():
        return 'MapEventLongPress';
      case MapEventMove():
        return 'MapEventMove';
      case MapEventMoveStart():
        return 'MapEventMoveStart';
      case MapEventMoveEnd():
        return 'MapEventMoveEnd';
      case MapEventFlingAnimation():
        return 'MapEventFlingAnimation';
      case MapEventFlingAnimationNotStarted():
        return 'MapEventFlingAnimationNotStarted';
      case MapEventFlingAnimationStart():
        return 'MapEventFlingAnimationStart';
      case MapEventFlingAnimationEnd():
        return 'MapEventFlingAnimationEnd';
      case MapEventDoubleTapZoom():
        return 'MapEventDoubleTapZoom';
      case MapEventScrollWheelZoom():
        return 'MapEventScrollWheelZoom';
      case MapEventDoubleTapZoomStart():
        return 'MapEventDoubleTapZoomStart';
      case MapEventDoubleTapZoomEnd():
        return 'MapEventDoubleTapZoomEnd';
      case MapEventRotate():
        return 'MapEventRotate';
      case MapEventRotateStart():
        return 'MapEventRotateStart';
      case MapEventRotateEnd():
        return 'MapEventRotateEnd';
      case MapEventNonRotatedSizeChange():
        return 'MapEventNonRotatedSizeChange';
      case null:
        return 'null';
      default:
        return 'Unknown';
    }
  }
}
