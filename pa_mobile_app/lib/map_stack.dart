import 'package:flutter/material.dart';
import 'package:pa_mobile_app/pages/map_page.dart';

class MapStack extends StatelessWidget {
  const MapStack({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      fit: StackFit.expand,
      children: [
        const MapPage(),
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
          child: Container(), // _getLeftbuttons(context),
        ),
        Positioned(
          right: 10,
          child: Container(), // _getRightButtons(context),
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
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    Text(
                      '750 Metre',
                      style: TextStyle(color: Colors.white, fontSize: Theme.of(context).textTheme.bodyMedium!.fontSize),
                    ),
                    SizedBox(width: 5),
                    Text(
                      '|',
                      style: TextStyle(color: Colors.white, fontSize: Theme.of(context).textTheme.bodyMedium!.fontSize),
                    ),
                    SizedBox(width: 5),
                    Text(
                      'Müsait park yeri: 20',
                      style: TextStyle(color: Colors.white, fontSize: Theme.of(context).textTheme.bodyMedium!.fontSize),
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
}
