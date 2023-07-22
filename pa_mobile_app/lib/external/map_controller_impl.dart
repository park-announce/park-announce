import 'dart:async';
import 'dart:math';

import 'package:flutter/widgets.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/camera.dart';
import 'package:pa_mobile_app/external/camera_fit.dart';
import 'package:pa_mobile_app/external/center_zoom.dart';
import 'package:pa_mobile_app/external/fit_bounds_options.dart';
import 'package:pa_mobile_app/external/internal_controller.dart';
import 'package:pa_mobile_app/external/latlng_bounds.dart';
import 'package:pa_mobile_app/external/map_controller.dart';
import 'package:pa_mobile_app/external/map_events.dart';
import 'package:pa_mobile_app/external/move_and_rotate_result.dart';
import 'package:pa_mobile_app/external/point_extensions.dart';

/// Implements [MapController] whilst exposing methods for internal use which
/// should not be visible to the user (e.g. for setting the current camera or
/// linking the internal controller).
class MapControllerImpl implements MapController {
  late FlutterMapInternalController _internalController;
  final _mapEventStreamController = StreamController<MapEvent>.broadcast();

  MapControllerImpl();

  set internalController(FlutterMapInternalController internalController) {
    _internalController = internalController;
  }

  StreamSink<MapEvent> get mapEventSink => _mapEventStreamController.sink;

  @override
  Stream<MapEvent> get mapEventStream => _mapEventStreamController.stream;

  @override
  bool move(
    LatLng center,
    double zoom, {
    Offset offset = Offset.zero,
    String? id,
  }) =>
      _internalController.move(
        center,
        zoom,
        offset: offset,
        hasGesture: false,
        source: MapEventSource.mapController,
        id: id,
      );

  @override
  bool rotate(double degree, {String? id}) => _internalController.rotate(
        degree,
        hasGesture: false,
        source: MapEventSource.mapController,
        id: id,
      );

  @override
  MoveAndRotateResult rotateAroundPoint(
    double degree, {
    Point<double>? point,
    Offset? offset,
    String? id,
  }) =>
      _internalController.rotateAroundPoint(
        degree,
        point: point,
        offset: offset,
        hasGesture: false,
        source: MapEventSource.mapController,
        id: id,
      );

  @override
  MoveAndRotateResult moveAndRotate(
    LatLng center,
    double zoom,
    double degree, {
    String? id,
  }) =>
      _internalController.moveAndRotate(
        center,
        zoom,
        degree,
        offset: Offset.zero,
        hasGesture: false,
        source: MapEventSource.mapController,
        id: id,
      );

  @override
  bool fitCamera(CameraFit cameraFit) => _internalController.fitCamera(
        cameraFit,
        offset: Offset.zero,
      );

  @override
  MapCamera get camera => _internalController.camera;

  @override
  @Deprecated(
    'Prefer `fitCamera` with a CameraFit.bounds() or CameraFit.insideBounds() instead. '
    'This method has been changed to use the new `CameraFit` classes which allows different kinds of fit. '
    'This method is deprecated since v6.',
  )
  bool fitBounds(
    LatLngBounds bounds, {
    FitBoundsOptions options = const FitBoundsOptions(padding: EdgeInsets.all(12)),
  }) =>
      fitCamera(
        options.inside
            ? CameraFit.insideBounds(
                bounds: bounds,
                padding: options.padding,
                maxZoom: options.maxZoom,
                forceIntegerZoomLevel: options.forceIntegerZoomLevel,
              )
            : CameraFit.bounds(
                bounds: bounds,
                padding: options.padding,
                maxZoom: options.maxZoom,
                forceIntegerZoomLevel: options.forceIntegerZoomLevel,
              ),
      );

  @override
  @Deprecated(
    'Prefer `CameraFit.bounds(bounds: bounds).fit(controller.camera)` or `CameraFit.insideBounds(bounds: bounds).fit(controller.camera)`. '
    'This method is replaced by applying a CameraFit to the MapCamera. '
    'This method is deprecated since v6.',
  )
  CenterZoom centerZoomFitBounds(
    LatLngBounds bounds, {
    FitBoundsOptions options = const FitBoundsOptions(padding: EdgeInsets.all(12)),
  }) {
    final cameraFit = options.inside
        ? CameraFit.insideBounds(
            bounds: bounds,
            padding: options.padding,
            maxZoom: options.maxZoom,
            forceIntegerZoomLevel: options.forceIntegerZoomLevel,
          )
        : CameraFit.bounds(
            bounds: bounds,
            padding: options.padding,
            maxZoom: options.maxZoom,
            forceIntegerZoomLevel: options.forceIntegerZoomLevel,
          );

    final fittedState = cameraFit.fit(camera);
    return CenterZoom(
      center: fittedState.center,
      zoom: fittedState.zoom,
    );
  }

  @override
  @Deprecated(
    'Prefer `controller.camera.pointToLatLng()`. '
    'This method is now accessible via the camera. '
    'This method is deprecated since v6.',
  )
  LatLng pointToLatLng(Point<num> screenPoint) => camera.pointToLatLng(screenPoint);

  @override
  @Deprecated(
    'Prefer `controller.camera.latLngToScreenPoint()`. '
    'This method is now accessible via the camera. '
    'This method is deprecated since v6.',
  )
  Point<double> latLngToScreenPoint(LatLng mapCoordinate) => camera.latLngToScreenPoint(mapCoordinate);

  @override
  @Deprecated(
    'Prefer `controller.camera.rotatePoint()`. '
    'This method is now accessible via the camera. '
    'This method is deprecated since v6.',
  )
  Point<double> rotatePoint(
    Point mapCenter,
    Point point, {
    bool counterRotation = true,
  }) =>
      camera.rotatePoint(
        mapCenter.toDoublePoint(),
        point.toDoublePoint(),
        counterRotation: counterRotation,
      );

  @override
  @Deprecated(
    'Prefer `controller.camera.center`. '
    'This getter is now accessible via the camera. '
    'This getter is deprecated since v6.',
  )
  LatLng get center => camera.center;

  @override
  @Deprecated(
    'Prefer `controller.camera.visibleBounds`. '
    'This getter is now accessible via the camera. '
    'This getter is deprecated since v6.',
  )
  LatLngBounds? get bounds => camera.visibleBounds;

  @override
  @Deprecated(
    'Prefer `controller.camera.zoom`. '
    'This getter is now accessible via the camera. '
    'This getter is deprecated since v6.',
  )
  double get zoom => camera.zoom;

  @override
  @Deprecated(
    'Prefer `controller.camera.rotation`. '
    'This getter is now accessible via the camera. '
    'This getter is deprecated since v6.',
  )
  double get rotation => camera.rotation;

  @override
  void dispose() {
    _mapEventStreamController.close();
  }
}
