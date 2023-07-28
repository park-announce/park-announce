import 'package:flutter/material.dart';
import 'package:pa_mobile_app/external/interactive_test_page.dart';

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
                padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                decoration: BoxDecoration(
                  color: Colors.grey.shade500,
                  borderRadius: BorderRadius.circular(30),
                ),
                child: const Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    Text(
                      '750 Metre',
                      style: TextStyle(color: Colors.white, fontSize: 12),
                    ),
                    SizedBox(width: 5),
                    Text(
                      '|',
                      style: TextStyle(color: Colors.white, fontSize: 12),
                    ),
                    SizedBox(width: 5),
                    Text(
                      'Müsait park yeri: 20',
                      style: TextStyle(color: Colors.white, fontSize: 12),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 30),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(
                width: MediaQuery.of(context).size.width / 2,
                child: ElevatedButton(
                  onPressed: () {},
                  style: TextButton.styleFrom(backgroundColor: Colors.green.shade400, foregroundColor: Colors.white),
                  child: const Text('Bu bölgede ara'),
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
