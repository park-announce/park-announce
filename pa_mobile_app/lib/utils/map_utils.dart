import 'package:flutter/material.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/map_controller.dart';

void animatedMapMove(LatLng destLocation, double destZoom, MapController mapController, TickerProvider provider) {
  // Create some tweens. These serve to split up the transition from one location to another.
  // In our case, we want to split the transition be<tween> our current map center and the destination.
  const String startedId = 'AnimatedMapController#MoveStarted';
  const String inProgressId = 'AnimatedMapController#MoveInProgress';
  const String finishedId = 'AnimatedMapController#MoveFinished';

  final camera = mapController.camera;
  final latTween = Tween<double>(begin: camera.center.latitude, end: destLocation.latitude);
  final lngTween = Tween<double>(begin: camera.center.longitude, end: destLocation.longitude);
  final zoomTween = Tween<double>(begin: camera.zoom, end: destZoom);

  // Create a animation controller that has a duration and a TickerProvider.
  final controller = AnimationController(duration: const Duration(milliseconds: 500), vsync: provider);
  // The animation determines what path the animation will take. You can try different Curves values, although I found
  // fastOutSlowIn to be my favorite.
  final Animation<double> animation = CurvedAnimation(parent: controller, curve: Curves.fastOutSlowIn);

  // Note this method of encoding the target destination is a workaround.
  // When proper animated movement is supported (see #1263) we should be able
  // to detect an appropriate animated movement event which contains the
  // target zoom/center.
  final startIdWithTarget = '$startedId#${destLocation.latitude},${destLocation.longitude},$destZoom';
  bool hasTriggeredMove = false;

  controller.addListener(() {
    final String id;
    if (animation.value == 1.0) {
      id = finishedId;
    } else if (!hasTriggeredMove) {
      id = startIdWithTarget;
    } else {
      id = inProgressId;
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
