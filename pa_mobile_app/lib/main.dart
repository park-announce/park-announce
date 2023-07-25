import 'package:flutter/material.dart';
import 'package:latlong2/latlong.dart';
import 'package:pa_mobile_app/external/interactive_test_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  LatLng position = const LatLng(41, 29);

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'flutter_map Demo',
        theme: ThemeData(
          useMaterial3: true,
          colorSchemeSeed: const Color(0xFF8dea88),
        ),
        home: SafeArea(
          child: Scaffold(
            body: MapStack(),
          ),
        ));
  }
}

class MapStack extends StatelessWidget {
  const MapStack({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      fit: StackFit.expand,
      children: [
        const InteractiveTestPage(),
        Positioned(
          top: 30,
          child: _getTopButtons(context),
        ),
        Positioned(
          bottom: 30,
          child: SizedBox(
            width: MediaQuery.of(context).size.width,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                GestureDetector(
                  onTap: () {
                    showModalBottomSheet(
                      context: context,
                      builder: (context) => SafeArea(
                        child: SingleChildScrollView(
                          child: Wrap(
                              children: List.generate(10, (index) => index)
                                  .map(
                                    (e) => TextButton(
                                      style: TextButton.styleFrom(fixedSize: Size(double.infinity, 3), padding: EdgeInsets.only(left: 1)),
                                      onPressed: () {
                                        debugPrint(e.toString());
                                      },
                                      child: SizedBox(
                                        width: MediaQuery.of(context).size.width,
                                        child: Row(
                                          mainAxisAlignment: MainAxisAlignment.center,
                                          children: [
                                            Text(e.toString(), style: TextStyle(fontSize: 12, color: Colors.blueGrey)),
                                          ],
                                        ),
                                      ),
                                    ),
                                  )
                                  .toList()),
                        ),
                      ),
                    );
                  },
                  child: CircleAvatar(
                    backgroundColor: Colors.greenAccent.shade700,
                    child: const Icon(
                      Icons.add,
                      color: Colors.white,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
        Positioned(
          left: 10,
          child: _getLeftbuttons(context),
        ),
        Positioned(
          right: 10,
          child: _getRightButtons(context),
        )
      ],
    );
  }

  SizedBox _getTopButtons(BuildContext context) {
    return SizedBox(
      width: MediaQuery.of(context).size.width,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                decoration: BoxDecoration(
                  color: Colors.grey.shade500,
                  borderRadius: BorderRadius.circular(30),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    Text(
                      '750 Metre',
                      style: const TextStyle(color: Colors.white, fontSize: 12),
                    ),
                    SizedBox(width: 5),
                    Text(
                      '|',
                      style: const TextStyle(color: Colors.white, fontSize: 12),
                    ),
                    SizedBox(width: 5),
                    Text(
                      'Müsait park yeri: 20',
                      style: const TextStyle(color: Colors.white, fontSize: 12),
                    ),
                  ],
                ),
              ),
            ],
          ),
          SizedBox(height: 30),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(
                width: MediaQuery.of(context).size.width / 2,
                child: ElevatedButton(
                  onPressed: () {},
                  child: Text('Bu bölgede ara'),
                  style: TextButton.styleFrom(backgroundColor: Colors.green.shade400, foregroundColor: Colors.white),
                ),
              ),
            ],
          )
        ],
      ),
    );
  }

  SizedBox _getRightButtons(BuildContext context) {
    return SizedBox(
      height: MediaQuery.of(context).size.height,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          GestureDetector(
              onTap: () {},
              child: const CircleAvatar(
                backgroundColor: Colors.black,
                child: Icon(
                  Icons.navigation,
                  color: Colors.white,
                ),
              )),
        ],
      ),
    );
  }

  Widget _getLeftbuttons(BuildContext context) {
    return SizedBox(
      height: MediaQuery.of(context).size.height,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          GestureDetector(
              onTap: () {},
              child: const CircleAvatar(
                backgroundColor: Colors.red,
                child: Icon(
                  Icons.sos_rounded,
                  color: Colors.white,
                ),
              )),
          const SizedBox(height: 10),
          GestureDetector(
              onTap: () {},
              child: const CircleAvatar(
                backgroundColor: Colors.blue,
                child: Icon(
                  Icons.local_parking,
                  color: Colors.white,
                ),
              )),
          const SizedBox(height: 10),
        ],
      ),
    );
  }
}
