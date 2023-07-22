import 'package:flutter/material.dart';
import 'package:pa_mobile_app/external/camera.dart';
import 'package:pa_mobile_app/external/map_controller.dart';
import 'package:pa_mobile_app/external/options.dart';

/// Allows descendents of [FlutterMap] to access the [MapCamera], [MapOptions]
/// and [MapController]. Those classes provide of/maybeOf methods for users to
/// use, those methods call the relevant methods provided by this class.
///
/// Using an [InheritedModel] means dependent widgets will only rebuild when
/// the aspect they reference is updated.
class FlutterMapInheritedModel extends InheritedModel<_FlutterMapAspect> {
  final FlutterMapData data;

  FlutterMapInheritedModel({
    super.key,
    required MapCamera camera,
    required MapController controller,
    required MapOptions options,
    required super.child,
  }) : data = FlutterMapData(
          camera: camera,
          controller: controller,
          options: options,
        );

  static FlutterMapData? _maybeOf(
    BuildContext context, [
    _FlutterMapAspect? aspect,
  ]) =>
      InheritedModel.inheritFrom<FlutterMapInheritedModel>(context, aspect: aspect)?.data;

  static MapCamera? maybeCameraOf(BuildContext context) => _maybeOf(context, _FlutterMapAspect.camera)?.camera;

  static MapController? maybeControllerOf(BuildContext context) => _maybeOf(context, _FlutterMapAspect.controller)?.controller;

  static MapOptions? maybeOptionsOf(BuildContext context) => _maybeOf(context, _FlutterMapAspect.options)?.options;

  @override
  bool updateShouldNotify(FlutterMapInheritedModel oldWidget) => data != oldWidget.data;

  @override
  bool updateShouldNotifyDependent(
    covariant FlutterMapInheritedModel oldWidget,
    Set<Object> dependencies,
  ) {
    for (final dependency in dependencies) {
      if (dependency is _FlutterMapAspect) {
        switch (dependency) {
          case _FlutterMapAspect.camera:
            if (data.camera != oldWidget.data.camera) return true;
          case _FlutterMapAspect.controller:
            if (data.controller != oldWidget.data.controller) return true;
          case _FlutterMapAspect.options:
            if (data.options != oldWidget.data.options) return true;
        }
      }
    }

    return false;
  }
}

class FlutterMapData {
  final MapCamera camera;
  final MapController controller;
  final MapOptions options;

  const FlutterMapData({
    required this.camera,
    required this.controller,
    required this.options,
  });
}

enum _FlutterMapAspect { camera, controller, options }
