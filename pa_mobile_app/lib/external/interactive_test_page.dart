import 'package:flutter/material.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/intective_flag.dart';
import 'package:pa_mobile_app/external/map_events.dart';
import 'package:pa_mobile_app/external/options.dart';
import 'package:pa_mobile_app/external/tile_layer.dart';
import 'package:pa_mobile_app/external/widget.dart';

class InteractiveTestPage extends StatefulWidget {
  const InteractiveTestPage({Key? key}) : super(key: key);

  @override
  State createState() {
    return _InteractiveTestPageState();
  }
}

class _InteractiveTestPageState extends State<InteractiveTestPage> {
  // Enable pinchZoom and doubleTapZoomBy by default
  int flags = InteractiveFlag.pinchZoom | InteractiveFlag.doubleTapZoom;

  MapEvent? _latestEvent;

  @override
  void initState() {
    super.initState();
  }

  void onMapEvent(MapEvent mapEvent) {
    if (mapEvent is! MapEventMove && mapEvent is! MapEventRotate) {
      // do not flood console with move and rotate events
      debugPrint(_eventName(mapEvent));
    }

    setState(() {
      _latestEvent = mapEvent;
    });
  }

  void updateFlags(int flag) {
    if (InteractiveFlag.hasFlag(flags, flag)) {
      // remove flag from flags
      flags &= ~flag;
    } else {
      // add flag to flags
      flags |= flag;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Test out Interactive flags!')),
      body: Padding(
        padding: const EdgeInsets.all(8),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: <Widget>[
                MaterialButton(
                  color: InteractiveFlag.hasDrag(flags) ? Colors.greenAccent : Colors.redAccent,
                  onPressed: () {
                    setState(() {
                      updateFlags(InteractiveFlag.drag);
                    });
                  },
                  child: const Text('Drag'),
                ),
                MaterialButton(
                  color: InteractiveFlag.hasFlingAnimation(flags) ? Colors.greenAccent : Colors.redAccent,
                  onPressed: () {
                    setState(() {
                      updateFlags(InteractiveFlag.flingAnimation);
                    });
                  },
                  child: const Text('Fling'),
                ),
                MaterialButton(
                  color: InteractiveFlag.hasPinchMove(flags) ? Colors.greenAccent : Colors.redAccent,
                  onPressed: () {
                    setState(() {
                      updateFlags(InteractiveFlag.pinchMove);
                    });
                  },
                  child: const Text('Pinch move'),
                ),
              ],
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: <Widget>[
                MaterialButton(
                  color: InteractiveFlag.hasDoubleTapZoom(flags) ? Colors.greenAccent : Colors.redAccent,
                  onPressed: () {
                    setState(() {
                      updateFlags(InteractiveFlag.doubleTapZoom);
                    });
                  },
                  child: const Text('Double tap zoom'),
                ),
                MaterialButton(
                  color: InteractiveFlag.hasRotate(flags) ? Colors.greenAccent : Colors.redAccent,
                  onPressed: () {
                    setState(() {
                      updateFlags(InteractiveFlag.rotate);
                    });
                  },
                  child: const Text('Rotate'),
                ),
                MaterialButton(
                  color: InteractiveFlag.hasPinchZoom(flags) ? Colors.greenAccent : Colors.redAccent,
                  onPressed: () {
                    setState(() {
                      updateFlags(InteractiveFlag.pinchZoom);
                    });
                  },
                  child: const Text('Pinch zoom'),
                ),
              ],
            ),
            Padding(
              padding: const EdgeInsets.only(top: 8, bottom: 8),
              child: Center(
                child: Text(
                  'Current event: ${_eventName(_latestEvent)}\nSource: ${_latestEvent?.source.name ?? "none"}',
                  textAlign: TextAlign.center,
                ),
              ),
            ),
            Flexible(
              child: FlutterMap(
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
                ],
              ),
            ),
          ],
        ),
      ),
    );
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
