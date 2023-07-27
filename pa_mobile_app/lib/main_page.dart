import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class MainPage extends StatelessWidget {
  const MainPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF132555),
      body: FutureBuilder<String?>(
          future: _getDisplayName(),
          builder: (context, snapshot) {
            if (snapshot.hasData && snapshot.data != null) {
              return SafeArea(
                child: Center(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Text(
                        'Welcome,',
                        style: TextStyle(color: Colors.white),
                      ),
                      Text(
                        snapshot.data!,
                        style: TextStyle(color: Colors.white),
                      ),
                    ],
                  ),
                ),
              );
            } else {
              return Text('No Data');
            }
          }),
    );
  }

  Future<String?> _getDisplayName() async {
    final SharedPreferences pref = await SharedPreferences.getInstance();
    return pref.getString('Name');
  }
}
