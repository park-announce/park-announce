import 'package:flutter/rendering.dart';
import 'package:http/http.dart';
import 'package:http/retry.dart';
import 'package:pa_mobile_app/external/tile_coordinates.dart';
import 'package:pa_mobile_app/external/tile_layer.dart';
import 'package:pa_mobile_app/external/tile_provider/base_tile_provider.dart';
import 'package:pa_mobile_app/external/tile_provider/network_image_provider.dart';

/// [TileProvider] to fetch tiles from the network
///
/// By default, a [RetryClient] is used to retry failed requests. 'dart:http'
/// or 'dart:io' might be needed to override this.
///
/// On the web, the 'User-Agent' header cannot be changed as specified in
/// [TileLayer.tileProvider]'s documentation, due to a Dart/browser limitation.
class NetworkTileProvider extends TileProvider {
  /// [TileProvider] to fetch tiles from the network
  ///
  /// By default, a [RetryClient] is used to retry failed requests. 'dart:http'
  /// or 'dart:io' might be needed to override this.
  ///
  /// On the web, the 'User-Agent' header cannot be changed as specified in
  /// [TileLayer.tileProvider]'s documentation, due to a Dart/browser limitation.
  NetworkTileProvider({
    super.headers = const {},
    BaseClient? httpClient,
  }) : httpClient = httpClient ?? RetryClient(Client());

  /// The HTTP client used to make network requests for tiles
  final BaseClient httpClient;

  @override
  ImageProvider getImage(TileCoordinates coordinates, TileLayer options) => FlutterMapNetworkImageProvider(
        url: getTileUrl(coordinates, options),
        fallbackUrl: getTileFallbackUrl(coordinates, options),
        headers: headers,
        httpClient: httpClient,
      );
}
