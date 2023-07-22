import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/bounds.dart';
import 'package:pa_mobile_app/external/crs.dart';
import 'package:pa_mobile_app/external/latlng_bounds.dart';
import 'package:pa_mobile_app/external/point_extensions.dart';
import 'package:pa_mobile_app/external/tile_bounds_at_zoom.dart';
import 'package:pa_mobile_app/external/tile_range.dart';

abstract class TileBounds {
  final Crs crs;
  final double _tileSize;
  final LatLngBounds? _latLngBounds;

  const TileBounds._(
    this.crs,
    this._tileSize,
    this._latLngBounds,
  );

  factory TileBounds({
    required Crs crs,
    required double tileSize,
    LatLngBounds? latLngBounds,
  }) {
    if (crs.infinite && latLngBounds == null) {
      return InfiniteTileBounds._(crs, tileSize, latLngBounds);
    } else if (crs.wrapLat == null && crs.wrapLng == null) {
      return DiscreteTileBounds._(crs, tileSize, latLngBounds);
    } else {
      return WrappedTileBounds._(crs, tileSize, latLngBounds);
    }
  }

  TileBoundsAtZoom atZoom(int zoom);

  // Returns true if these bounds may no longer be valid for the given
  // parameters.
  bool shouldReplace(
    Crs crs,
    double tileSize,
    LatLngBounds? latLngBounds,
  ) =>
      (crs != this.crs || tileSize != _tileSize || latLngBounds != _latLngBounds);
}

class InfiniteTileBounds extends TileBounds {
  const InfiniteTileBounds._(
    Crs crs,
    double tileSize,
    LatLngBounds? latLngBounds,
  ) : super._(crs, tileSize, latLngBounds);

  @override
  TileBoundsAtZoom atZoom(int zoom) => const InfiniteTileBoundsAtZoom();
}

class DiscreteTileBounds extends TileBounds {
  final Map<int, TileBoundsAtZoom> _tileBoundsAtZoomCache = {};

  DiscreteTileBounds._(
    Crs crs,
    double tileSize,
    LatLngBounds? latLngBounds,
  ) : super._(crs, tileSize, latLngBounds);

  @override
  TileBoundsAtZoom atZoom(int zoom) {
    return _tileBoundsAtZoomCache.putIfAbsent(zoom, () => _tileBoundsAtZoomImpl(zoom));
  }

  TileBoundsAtZoom _tileBoundsAtZoomImpl(int zoom) {
    final zoomDouble = zoom.toDouble();

    final Bounds<num> pixelBounds;
    if (_latLngBounds == null) {
      pixelBounds = crs.getProjectedBounds(zoomDouble)!;
    } else {
      pixelBounds = Bounds(
        crs.latLngToPoint(_latLngBounds!.southWest, zoomDouble).floor(),
        crs.latLngToPoint(_latLngBounds!.northEast, zoomDouble).ceil(),
      );
    }

    return DiscreteTileBoundsAtZoom(
      DiscreteTileRange.fromPixelBounds(
        zoom: zoom,
        tileSize: _tileSize,
        pixelBounds: pixelBounds,
      ),
    );
  }
}

class WrappedTileBounds extends TileBounds {
  final Map<int, WrappedTileBoundsAtZoom> _tileBoundsAtZoomCache = {};

  WrappedTileBounds._(
    Crs crs,
    double tileSize,
    LatLngBounds? latLngBounds,
  ) : super._(crs, tileSize, latLngBounds);

  @override
  WrappedTileBoundsAtZoom atZoom(int zoom) {
    return _tileBoundsAtZoomCache.putIfAbsent(zoom, () => _tileBoundsAtZoomImpl(zoom));
  }

  WrappedTileBoundsAtZoom _tileBoundsAtZoomImpl(int zoom) {
    final zoomDouble = zoom.toDouble();

    final Bounds<num> pixelBounds;
    if (_latLngBounds == null) {
      pixelBounds = crs.getProjectedBounds(zoomDouble)!;
    } else {
      pixelBounds = Bounds(
        crs.latLngToPoint(_latLngBounds!.southWest, zoomDouble).floor(),
        crs.latLngToPoint(_latLngBounds!.northEast, zoomDouble).ceil(),
      );
    }

    final tzDouble = zoom.toDouble();

    (int, int)? wrapX;
    if (crs.wrapLng case final wrapLng?) {
      final wrapXMin = (crs.latLngToPoint(LatLng(0, wrapLng.$1), tzDouble).x / _tileSize).floor();
      final wrapXMax = (crs.latLngToPoint(LatLng(0, wrapLng.$2), tzDouble).x / _tileSize).ceil();
      wrapX = (wrapXMin, wrapXMax - 1);
    }

    (int, int)? wrapY;
    if (crs.wrapLat case final wrapLat?) {
      final wrapYMin = (crs.latLngToPoint(LatLng(wrapLat.$1, 0), tzDouble).y / _tileSize).floor();
      final wrapYMax = (crs.latLngToPoint(LatLng(wrapLat.$2, 0), tzDouble).y / _tileSize).ceil();
      wrapY = (wrapYMin, wrapYMax - 1);
    }

    return WrappedTileBoundsAtZoom(
      tileRange: DiscreteTileRange.fromPixelBounds(
        zoom: zoom,
        tileSize: _tileSize,
        pixelBounds: pixelBounds,
      ),
      wrappedAxisIsAlwaysInBounds: _latLngBounds == null,
      wrapX: wrapX,
      wrapY: wrapY,
    );
  }
}
