import 'dart:math';

import 'package:flutter/foundation.dart';
import 'package:flutter/rendering.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/camera.dart';
import 'package:pa_mobile_app/external/camera_fit.dart';
import 'package:pa_mobile_app/external/flutter_map_interactive_viewer.dart';
import 'package:pa_mobile_app/external/map_controller_impl.dart';
import 'package:pa_mobile_app/external/map_events.dart';
import 'package:pa_mobile_app/external/move_and_rotate_result.dart';
import 'package:pa_mobile_app/external/options.dart';
import 'package:pa_mobile_app/external/point_extensions.dart';
import 'package:pa_mobile_app/external/position.dart';
import 'package:pa_mobile_app/external/positioned_tap_detector_2.dart';

// This controller is for internal use. All updates to the state should be done
// by calling methods of this class to ensure consistency.
class FlutterMapInternalController extends ValueNotifier<_InternalState> {
  late final FlutterMapInteractiveViewerState _interactiveViewerState;
  late MapControllerImpl _mapControllerImpl;

  FlutterMapInternalController(MapOptions options)
      : super(
          _InternalState(
            options: options,
            camera: MapCamera.initialCamera(options),
          ),
        );

  // Link the viewer state with the controller. This should be done once when
  // the FlutterMapInteractiveViewerState is initialized.
  set interactiveViewerState(
    FlutterMapInteractiveViewerState interactiveViewerState,
  ) =>
      _interactiveViewerState = interactiveViewerState;

  MapOptions get options => value.options;

  MapCamera get camera => value.camera;

  void linkMapController(MapControllerImpl mapControllerImpl) {
    _mapControllerImpl = mapControllerImpl;
    _mapControllerImpl.internalController = this;
  }

  /// This setter should only be called in this class or within tests. Changes
  /// to the [FlutterMapInternalState] should be done via methods in this class.
  @visibleForTesting
  @override
  // ignore: library_private_types_in_public_api
  set value(_InternalState value) => super.value = value;

  // Note: All named parameters are required to prevent inconsistent default
  // values since this method can be called by MapController which declares
  // defaults.
  bool move(
    LatLng newCenter,
    double newZoom, {
    required Offset offset,
    required bool hasGesture,
    required MapEventSource source,
    required String? id,
  }) {
    // Algorithm thanks to https://github.com/tlserver/flutter_map_location_marker
    if (offset != Offset.zero) {
      final newPoint = camera.project(newCenter, newZoom);
      newCenter = camera.unproject(
        camera.rotatePoint(
          newPoint,
          newPoint - Point(offset.dx, offset.dy),
        ),
        newZoom,
      );
    }

    MapCamera? newCamera = camera.withPosition(
      center: newCenter,
      zoom: camera.clampZoom(newZoom),
    );

    newCamera = options.cameraConstraint.constrain(newCamera);
    if (newCamera == null || (newCamera.center == camera.center && newCamera.zoom == camera.zoom)) {
      return false;
    }

    final oldCamera = camera;
    value = value.withMapCamera(newCamera);

    final movementEvent = MapEventWithMove.fromSource(
      oldCamera: oldCamera,
      camera: camera,
      hasGesture: hasGesture,
      source: source,
      id: id,
    );
    if (movementEvent != null) _emitMapEvent(movementEvent);

    options.onPositionChanged?.call(
      MapPosition(
        center: newCenter,
        bounds: camera.visibleBounds,
        zoom: newZoom,
        hasGesture: hasGesture,
      ),
      hasGesture,
    );

    return true;
  }

  // Note: All named parameters are required to prevent inconsistent default
  // values since this method can be called by MapController which declares
  // defaults.
  bool rotate(
    double newRotation, {
    required bool hasGesture,
    required MapEventSource source,
    required String? id,
  }) {
    if (newRotation != camera.rotation) {
      final newCamera = options.cameraConstraint.constrain(
        camera.withRotation(newRotation),
      );
      if (newCamera == null) return false;

      final oldCamera = camera;

      // Update camera then emit events and callbacks
      value = value.withMapCamera(newCamera);

      _emitMapEvent(
        MapEventRotate(
          id: id,
          source: source,
          oldCamera: oldCamera,
          camera: camera,
        ),
      );
      return true;
    }

    return false;
  }

  // Note: All named parameters are required to prevent inconsistent default
  // values since this method can be called by MapController which declares
  // defaults.
  MoveAndRotateResult rotateAroundPoint(
    double degree, {
    required Point<double>? point,
    required Offset? offset,
    required bool hasGesture,
    required MapEventSource source,
    required String? id,
  }) {
    if (point != null && offset != null) {
      throw ArgumentError('Only one of `point` or `offset` may be non-null');
    }
    if (point == null && offset == null) {
      throw ArgumentError('One of `point` or `offset` must be non-null');
    }

    if (degree == camera.rotation) {
      return MoveAndRotateResult(false, false);
    }

    if (offset == Offset.zero) {
      return MoveAndRotateResult(
        true,
        rotate(
          degree,
          hasGesture: hasGesture,
          source: source,
          id: id,
        ),
      );
    }

    final rotationDiff = degree - camera.rotation;
    final rotationCenter =
        camera.project(camera.center) + (point != null ? (point - (camera.nonRotatedSize / 2.0)) : Point(offset!.dx, offset.dy)).rotate(camera.rotationRad);

    return MoveAndRotateResult(
      move(
        camera.unproject(
          rotationCenter + (camera.project(camera.center) - rotationCenter).rotate(degToRadian(rotationDiff)),
        ),
        camera.zoom,
        offset: Offset.zero,
        hasGesture: hasGesture,
        source: source,
        id: id,
      ),
      rotate(
        camera.rotation + rotationDiff,
        hasGesture: hasGesture,
        source: source,
        id: id,
      ),
    );
  }

  // Note: All named parameters are required to prevent inconsistent default
  // values since this method can be called by MapController which declares
  // defaults.
  MoveAndRotateResult moveAndRotate(
    LatLng newCenter,
    double newZoom,
    double newRotation, {
    required Offset offset,
    required bool hasGesture,
    required MapEventSource source,
    required String? id,
  }) =>
      MoveAndRotateResult(
        move(
          newCenter,
          newZoom,
          offset: offset,
          hasGesture: hasGesture,
          source: source,
          id: id,
        ),
        rotate(newRotation, id: id, source: source, hasGesture: hasGesture),
      );

  // Note: All named parameters are required to prevent inconsistent default
  // values since this method can be called by MapController which declares
  // defaults.
  bool fitCamera(
    CameraFit cameraFit, {
    required Offset offset,
  }) {
    final fitted = cameraFit.fit(camera);

    return move(
      fitted.center,
      fitted.zoom,
      offset: offset,
      hasGesture: false,
      source: MapEventSource.fitCamera,
      id: null,
    );
  }

  bool setNonRotatedSizeWithoutEmittingEvent(
    Point<double> nonRotatedSize,
  ) {
    if (nonRotatedSize != MapCamera.kImpossibleSize && nonRotatedSize != camera.nonRotatedSize) {
      value = value.withMapCamera(camera.withNonRotatedSize(nonRotatedSize));
      return true;
    }

    return false;
  }

  void setOptions(MapOptions newOptions) {
    assert(
      newOptions != value.options,
      'Should not update options unless they change',
    );

    final newCamera = camera.withOptions(newOptions);

    assert(
      newOptions.cameraConstraint.constrain(newCamera) == newCamera,
      'MapCamera is no longer within the cameraConstraint after an option change.',
    );

    if (options.interactionOptions != newOptions.interactionOptions) {
      _interactiveViewerState.updateGestures(
        options.interactionOptions,
        newOptions.interactionOptions,
      );
    }

    value = _InternalState(
      options: newOptions,
      camera: newCamera,
    );
  }

  // To be called when a gesture that causes movement starts.
  void moveStarted(MapEventSource source) {
    _emitMapEvent(
      MapEventMoveStart(
        camera: camera,
        source: source,
      ),
    );
  }

  // To be called when an ongoing drag movement updates.
  void dragUpdated(MapEventSource source, Offset offset) {
    final oldCenterPt = camera.project(camera.center);

    final newCenterPt = oldCenterPt + offset.toPoint();
    final newCenter = camera.unproject(newCenterPt);

    move(
      newCenter,
      camera.zoom,
      offset: Offset.zero,
      hasGesture: true,
      source: source,
      id: null,
    );
  }

  // To be called when a drag gesture ends.
  void moveEnded(MapEventSource source) {
    _emitMapEvent(
      MapEventMoveEnd(
        camera: camera,
        source: source,
      ),
    );
  }

  // To be called when a rotation gesture starts.
  void rotateStarted(MapEventSource source) {
    _emitMapEvent(
      MapEventRotateStart(
        camera: camera,
        source: source,
      ),
    );
  }

  // To be called when a rotation gesture ends.
  void rotateEnded(MapEventSource source) {
    _emitMapEvent(
      MapEventRotateEnd(
        camera: camera,
        source: source,
      ),
    );
  }

  // To be called when a fling gesture starts.
  void flingStarted(MapEventSource source) {
    _emitMapEvent(
      MapEventFlingAnimationStart(
        camera: camera,
        source: MapEventSource.flingAnimationController,
      ),
    );
  }

  // To be called when a fling gesture ends.
  void flingEnded(MapEventSource source) {
    _emitMapEvent(
      MapEventFlingAnimationEnd(
        camera: camera,
        source: source,
      ),
    );
  }

  // To be called when a fling gesture does not start.
  void flingNotStarted(MapEventSource source) {
    _emitMapEvent(
      MapEventFlingAnimationNotStarted(
        camera: camera,
        source: source,
      ),
    );
  }

  // To be called when a double tap zoom starts.
  void doubleTapZoomStarted(MapEventSource source) {
    _emitMapEvent(
      MapEventDoubleTapZoomStart(
        camera: camera,
        source: source,
      ),
    );
  }

  // To be called when a double tap zoom ends.
  void doubleTapZoomEnded(MapEventSource source) {
    _emitMapEvent(
      MapEventDoubleTapZoomEnd(
        camera: camera,
        source: source,
      ),
    );
  }

  void tapped(
    MapEventSource source,
    TapPosition tapPosition,
    LatLng position,
  ) {
    options.onTap?.call(tapPosition, position);
    _emitMapEvent(
      MapEventTap(
        tapPosition: position,
        camera: camera,
        source: source,
      ),
    );
  }

  void secondaryTapped(
    MapEventSource source,
    TapPosition tapPosition,
    LatLng position,
  ) {
    options.onSecondaryTap?.call(tapPosition, position);
    _emitMapEvent(
      MapEventSecondaryTap(
        tapPosition: position,
        camera: camera,
        source: source,
      ),
    );
  }

  void longPressed(
    MapEventSource source,
    TapPosition tapPosition,
    LatLng position,
  ) {
    options.onLongPress?.call(tapPosition, position);
    _emitMapEvent(
      MapEventLongPress(
        tapPosition: position,
        camera: camera,
        source: MapEventSource.longPress,
      ),
    );
  }

  // To be called when the map's size constraints change.
  void nonRotatedSizeChange(
    MapEventSource source,
    MapCamera oldCamera,
    MapCamera newCamera,
  ) {
    _emitMapEvent(
      MapEventNonRotatedSizeChange(
        source: MapEventSource.nonRotatedSizeChange,
        oldCamera: oldCamera,
        camera: newCamera,
      ),
    );
  }

  void _emitMapEvent(MapEvent event) {
    if (event.source == MapEventSource.mapController && event is MapEventMove) {
      _interactiveViewerState.interruptAnimatedMovement(event);
    }

    options.onMapEvent?.call(event);

    _mapControllerImpl.mapEventSink.add(event);
  }
}

class _InternalState {
  final MapCamera camera;
  final MapOptions options;

  const _InternalState({
    required this.options,
    required this.camera,
  });

  _InternalState withMapCamera(MapCamera camera) => _InternalState(
        options: options,
        camera: camera,
      );
}
